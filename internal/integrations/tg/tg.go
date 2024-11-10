package tg

import (
	"context"
	"fmt"
	"time"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
	"go.uber.org/zap"
)

type InjectedSessionStorage struct {
	config []byte
}

func NewInjectedSessionStorage(preload string) *InjectedSessionStorage {
	config := make([]byte, 0)
	if preload != "" {
		config = []byte(preload)
	}

	return &InjectedSessionStorage{config: config}
}

func (s *InjectedSessionStorage) LoadSession(_ context.Context) ([]byte, error) {
	return s.config, nil
}

func (s *InjectedSessionStorage) StoreSession(_ context.Context, data []byte) error {
	s.config = data
	return nil
}

type Client struct {
	session         *InjectedSessionStorage
	uid             string
	appID           int
	appHash         string
	pendingCodeHash string
	ctx             context.Context
	logger          *zap.Logger
}

func NewTelegramClient(
	uid string,
	appID int,
	appHash string,
	raw string,
) (*Client, error) {
	l, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}

	sessionStorage := NewInjectedSessionStorage(raw)
	sessionStorageSvc := &Client{
		uid:     uid,
		session: sessionStorage,
		appID:   appID,
		appHash: appHash,
		logger:  l,
	}

	return sessionStorageSvc, nil
}

func (c *Client) Validate(ctx context.Context, client *telegram.Client) error {
	err := client.Ping(ctx)
	if err != nil {
		return err
	}
	status, err := client.Auth().Status(ctx)
	if err != nil {
		return fmt.Errorf("failed to check auth status :%w", err)
	}

	if !status.Authorized {
		return ErrAuthExp
	}

	return nil
}

func (c *Client) SendCode(ctx context.Context, client *telegram.Client) error {
	return client.Run(ctx, func(ctx context.Context) error {
		resp, err := client.Auth().SendCode(ctx, c.uid, auth.SendCodeOptions{})
		if err != nil {
			return err
		}

		c.pendingCodeHash = resp.(*tg.AuthSentCode).PhoneCodeHash
		return nil
	})
}

func (c *Client) AuthenticateWithCode(ctx context.Context, code string, client *telegram.Client) error {
	return client.Run(ctx, func(ctx context.Context) error {
		_, err := client.Auth().SignIn(ctx, c.uid, code, c.pendingCodeHash)
		if err != nil {
			return err
		}

		c.pendingCodeHash = ""
		return nil
	})
}

func (c *Client) Raw() *telegram.Client {
	return telegram.NewClient(c.appID, c.appHash, telegram.Options{
		SessionStorage: c.session,
		MaxRetries:     5,
		DialTimeout:    time.Second * 10,
		NoUpdates:      true,
		Logger:         c.logger,
	})
}

func (c *Client) RawWithoutSession() *telegram.Client {
	return telegram.NewClient(c.appID, c.appHash, telegram.Options{
		MaxRetries:  5,
		DialTimeout: time.Second * 10,
		NoUpdates:   true,
		Logger:      c.logger,
	})
}
