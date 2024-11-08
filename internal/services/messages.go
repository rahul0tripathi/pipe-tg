package services

import (
	"context"
	"fmt"
	"strconv"

	"github.com/gotd/td/tg"
	"github.com/rahul0tripathi/pipetg/entity"
	wrapper "github.com/rahul0tripathi/pipetg/internal/integrations/tg"
)

type MessageLogger struct {
	tg *wrapper.Client
}

func NewMessageLogger(c *wrapper.Client) *MessageLogger {
	return &MessageLogger{tg: c}
}

func (m *MessageLogger) All(ctx context.Context) (interface{}, error) {
	out := make([]entity.PipeMessage, 0)
	client := m.tg.Raw()
	err := client.Run(ctx, func(ctx context.Context) error {
		if err := m.tg.Validate(ctx, client); err != nil {
			return err
		}
		api := client.API()
		self, err := api.UsersGetUsers(ctx, []tg.InputUserClass{
			&tg.InputUserSelf{},
		})
		if err != nil {
			return fmt.Errorf("failed to get self user: %w", err)
		}

		selfUser := self[0].(*tg.User)

		offsetPeer := &tg.InputPeerUser{
			UserID:     selfUser.ID,
			AccessHash: selfUser.AccessHash,
		}

		dialogs, err := api.MessagesGetDialogs(ctx, &tg.MessagesGetDialogsRequest{
			OffsetPeer: offsetPeer,
			OffsetID:   0,
			OffsetDate: 0,
			Limit:      100,
			Hash:       0,
		})
		if err != nil {
			return fmt.Errorf("failed to get dialogs: %w", err)
		}

		dialogsObj := dialogs.(*tg.MessagesDialogs)

		channelInfo := make(map[int64]*tg.Channel)

		for _, chat := range dialogsObj.Chats {
			if channel, ok := chat.(*tg.Channel); ok {
				channelInfo[channel.ID] = channel
			}
		}

		for _, dialog := range dialogsObj.Dialogs {
			if channel, ok := dialog.GetPeer().(*tg.PeerChannel); ok {
				channelData := channelInfo[channel.ChannelID]
				if channelData == nil {
					fmt.Printf("No access information for channel %d", channel.ChannelID)
					continue
				}

				// Now create input peer with proper access hash
				messages, err := api.MessagesGetHistory(ctx, &tg.MessagesGetHistoryRequest{
					Peer: &tg.InputPeerChannel{
						ChannelID:  channelData.ID,
						AccessHash: channelData.AccessHash, // This is crucial!
					},
					Limit: 100,
				})
				if err != nil {
					fmt.Printf("Failed to get messages for channel %s: %v", channelData.Title, err)
					continue
				}
				var logs []entity.PipeMessage
				msgObj, ok := messages.(*tg.MessagesChannelMessages)
				if !ok {
					fmt.Printf("Unexpected message type for channel %s", channelData.Title)
					continue
				}

				for _, msg := range msgObj.Messages {
					if m, ok := msg.(*tg.Message); ok && m.Message != "" {
						logs = append(logs, entity.PipeMessage{
							ID:        strconv.Itoa(m.ID),
							Channel:   channelData.Title,
							Message:   m.Message,
							Timestamp: m.Date,
						})
					}
				}
				out = append(out, logs...)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return out, nil
}
