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
	raw             *telegram.Client
	uid             string
	pendingCodeHash string
}

func NewTelegramClient(
	uid string,
	appID int,
	appHash string,
	raw string,
	logger *zap.Logger,
) (*Client, error) {
	sessionStorage := NewInjectedSessionStorage(raw)
	sessionStorageSvc := &Client{
		uid:     uid,
		session: sessionStorage,
		raw: telegram.NewClient(appID, appHash, telegram.Options{
			SessionStorage: sessionStorage,
			MaxRetries:     5,
			DialTimeout:    time.Minute * 10,
			NoUpdates:      true,
			Logger:         logger,
		}),
	}

	return sessionStorageSvc, nil
}

func (c *Client) Validate(ctx context.Context) error {
	err := c.raw.Ping(ctx)
	if err != nil {
		return err
	}
	status, err := c.raw.Auth().Status(ctx)
	if err != nil {
		return fmt.Errorf("failed to check auth status :%w", err)
	}

	if !status.Authorized {
		return ErrAuthExp
	}

	return nil
}

func (c *Client) SendCode(ctx context.Context) error {
	return c.raw.Run(ctx, func(ctx context.Context) error {
		resp, err := c.raw.Auth().SendCode(ctx, c.uid, auth.SendCodeOptions{})
		if err != nil {
			return err
		}

		c.pendingCodeHash = resp.(*tg.AuthSentCode).PhoneCodeHash
		return nil
	})
}

func (c *Client) AuthenticateWithCode(ctx context.Context, code string) error {
	return c.raw.Run(ctx, func(ctx context.Context) error {
		_, err := c.raw.Auth().SignIn(ctx, c.uid, code, c.pendingCodeHash)
		if err != nil {
			return err
		}

		c.pendingCodeHash = ""
		return nil
	})
}

func (c *Client) Raw() *telegram.Client {
	return c.raw
}
