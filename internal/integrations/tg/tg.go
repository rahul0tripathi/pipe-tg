package tg

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type TgConn struct{}

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
	s.config = make([]byte, len(data))
	copy(s.config, data)
	return nil
}

type Client struct {
	session         *InjectedSessionStorage
	uid             string
	appID           int
	appHash         string
	pendingCodeHash string
	logger          *zap.Logger
	isAuth          bool
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
	resp, err := client.Auth().SendCode(ctx, c.uid, auth.SendCodeOptions{})
	if err != nil {
		return err
	}

	c.pendingCodeHash = resp.(*tg.AuthSentCode).PhoneCodeHash
	return nil
}

func (c *Client) AuthenticateWithCode(ctx context.Context, code string, client *telegram.Client) error {
	_, err := client.Auth().SignIn(ctx, c.uid, code, c.pendingCodeHash)
	if err != nil {
		return err
	}

	c.pendingCodeHash = ""
	return nil
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

func (c *Client) WithContext(ctx context.Context, exec func(ctx context.Context) error) error {
	conn := c.Raw()
	return conn.Run(ctx, func(ctx context.Context) error {
		if err := c.Validate(ctx, conn); err != nil {
			return err
		}

		ctx = context.WithValue(ctx, TgConn{}, conn)
		return exec(ctx)
	})
}

func (c *Client) WithUncheckedEchoContext(e echo.Context, exec func(ctx context.Context) error) error {
	return c.WithUncheckedContext(e.Request().Context(), exec)
}

func (c *Client) WithUncheckedContext(ctx context.Context, exec func(ctx context.Context) error) error {
	conn := c.Raw()
	return conn.Run(ctx, func(ctx context.Context) error {

		ctx = context.WithValue(ctx, TgConn{}, conn)
		return exec(ctx)
	})
}

func (c *Client) WithEchoContext(e echo.Context, exec func(ctx context.Context) error) error {
	return c.WithContext(e.Request().Context(), exec)
}

func (c *Client) GetTgConnFromCtx(ctx context.Context) (*telegram.Client, error) {
	_rawConn := ctx.Value(TgConn{})
	if _rawConn == nil {
		return c.Raw(), nil
	}
	conn, ok := _rawConn.(*telegram.Client)
	if !ok {
		return nil, errors.New("failed to convert ctx val to *telegram.Client")
	}

	return conn, nil
}

func (c *Client) GetSessionConfig() (string, error) {
	sessionData, err := c.session.LoadSession(nil)
	if err != nil {
		return "", fmt.Errorf("failed to load session: %w", err)
	}

	return string(sessionData), nil
}
