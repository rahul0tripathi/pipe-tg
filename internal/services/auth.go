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
	return a.tg.SendCode(ctx, a.tg.Raw())
}

func (a *AuthFlowService) SubmitCode(ctx context.Context, code string) error {
	return a.tg.AuthenticateWithCode(ctx, code, a.tg.Raw())
}
