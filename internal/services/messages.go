package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gotd/td/tg"
	gotg "github.com/gotd/td/tg"
	"github.com/rahul0tripathi/pipetg/entity"
	wrapper "github.com/rahul0tripathi/pipetg/internal/integrations/tg"
	"github.com/rahul0tripathi/pipetg/pkg/log"
)

var (
	_backOffDuration = time.Second * 2
)

type Scraper struct {
	tg     *wrapper.Client
	window time.Duration
}

func NewScraper(c *wrapper.Client, window time.Duration) *Scraper {
	fmt.Println(window.String())
	return &Scraper{tg: c, window: window}
}

func (s *Scraper) scrape(
	ctx context.Context,
	api *gotg.Client,
	chans map[int64]*tg.Channel,
) ([]entity.PipeMessage, error) {
	logger := log.GetLogger(ctx)
	endTime := time.Now()
	startTime := endTime.Add(-(s.window))

	var collector []entity.PipeMessage
	for _, channel := range chans {
		messages, err := s.scrapeChannel(ctx, api, channel, startTime, endTime)
		if err != nil {
			logger.Error(
				"failed to scrape channel",
				log.Err(err),
				log.Str("title", channel.Title),
			)
			continue
		}
		collector = append(collector, messages...)
	}

	return collector, nil
}

func (s *Scraper) scrapeChannel(
	ctx context.Context,
	api *gotg.Client,
	channel *tg.Channel,
	startTime, endTime time.Time,
) ([]entity.PipeMessage, error) {
	var messages []entity.PipeMessage

	req := &tg.MessagesGetHistoryRequest{
		Peer: &tg.InputPeerChannel{
			ChannelID:  channel.ID,
			AccessHash: channel.AccessHash,
		},
		OffsetDate: int(endTime.Unix()),
		OffsetID:   0,
		AddOffset:  0,
		Limit:      100,
		MaxID:      0,
		MinID:      0,
	}

	for {
		batch, err := s.fetchMessageBatch(ctx, api, req)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch message batch: %w", err)
		}

		if batch == nil || len(batch.Messages) == 0 {
			break
		}

		reachedStartTime := false
		newMessages, shouldContinue := s.processMessageBatch(
			batch.Messages,
			channel.Title,
			startTime,
			req,
			&reachedStartTime,
		)
		messages = append(messages, newMessages...)

		if reachedStartTime || len(batch.Messages) < 100 || !shouldContinue {
			break
		}

		time.Sleep(_backOffDuration)
	}

	return messages, nil
}

func (s *Scraper) fetchMessageBatch(
	ctx context.Context,
	api *gotg.Client,
	req *tg.MessagesGetHistoryRequest,
) (*tg.MessagesChannelMessages, error) {
	messages, err := api.MessagesGetHistory(ctx, req)
	if err != nil {
		return nil, err
	}

	msgObj, ok := messages.(*tg.MessagesChannelMessages)
	if !ok {
		return nil, errors.New("type assertion failed: expected tg.MessagesChannelMessages")
	}

	return msgObj, nil
}

func (s *Scraper) processMessageBatch(
	messages []tg.MessageClass,
	channelTitle string,
	startTime time.Time,
	req *tg.MessagesGetHistoryRequest,
	reachedStartTime *bool,
) ([]entity.PipeMessage, bool) {
	var batch []entity.PipeMessage

	for _, msg := range messages {
		message, ok := msg.(*tg.Message)
		if !ok || message.Message == "" {
			continue
		}

		msgTime := time.Unix(int64(message.Date), 0)
		if msgTime.Before(startTime) {
			*reachedStartTime = true
			return batch, false
		}

		batch = append(batch, entity.PipeMessage{
			ID:        fmt.Sprintf("%d", message.ID),
			Channel:   channelTitle,
			Message:   message.Message,
			Timestamp: message.Date,
		})

		req.OffsetID = message.ID
		req.OffsetDate = message.Date
	}

	return batch, true
}

func (s *Scraper) channels(
	ctx context.Context,
	api *gotg.Client,
	id int64,
	accessHash int64,
) (map[int64]*tg.Channel, error) {
	_dialogs, err := api.MessagesGetDialogs(ctx, &tg.MessagesGetDialogsRequest{
		OffsetPeer: &tg.InputPeerUser{
			UserID:     id,
			AccessHash: accessHash,
		},
		Limit: 100,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get dialogs: %w", err)
	}

	dialogs, ok := _dialogs.(*tg.MessagesDialogs)
	if !ok {
		return nil, errors.New("channels: tg.MessagesDialogs assertion failed")
	}

	channels := make(map[int64]*tg.Channel)
	for _, chat := range dialogs.Chats {
		if channel, ok := chat.(*tg.Channel); ok {
			channels[channel.ID] = channel
		}
	}

	return channels, nil
}

func (s *Scraper) self(ctx context.Context, api *gotg.Client) (*tg.User, error) {
	_self, err := api.UsersGetUsers(ctx, []tg.InputUserClass{
		&tg.InputUserSelf{},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get self user: %w", err)
	}

	self, ok := _self[0].(*tg.User)
	if !ok {
		return nil, errors.New("self: tg.User assertion failed")
	}

	return self, nil
}

func (s *Scraper) Run(ctx context.Context) ([]entity.PipeMessage, error) {
	client, err := s.tg.GetTgConnFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	usr, err := s.self(ctx, client.API())
	if err != nil {
		return nil, err
	}

	channels, err := s.channels(ctx, client.API(), usr.ID, usr.AccessHash)
	if err != nil {
		return nil, err
	}

	return s.scrape(ctx, client.API(), channels)
}
