package services

import (
	"context"

	"github.com/rahul0tripathi/pipetg/internal/integrations/tg"
)

type AuthFlowService struct {
	tg *tg.Client
}

func NewAuthFlowService(c *tg.Client) *AuthFlowService {
	return &AuthFlowService{tg: c}
}

func (a *AuthFlowService) RequestNewCode(ctx context.Context) error {
	conn, err := a.tg.GetTgConnFromCtx(ctx)
	if err != nil {
		return err
	}
	return a.tg.SendCode(ctx, conn)
}

func (a *AuthFlowService) SubmitCode(ctx context.Context, code string) error {
	conn, err := a.tg.GetTgConnFromCtx(ctx)
	if err != nil {
		return err
	}

	return a.tg.AuthenticateWithCode(ctx, code, conn)
}
