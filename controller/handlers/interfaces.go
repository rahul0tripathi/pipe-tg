package handlers

import (
	"context"

	"github.com/rahul0tripathi/pipetg/entity"
)

type AuthFlowSvc interface {
	RequestNewCode(ctx context.Context) error
	SubmitCode(ctx context.Context, code string) error
}

type Scraper interface {
	Run(ctx context.Context) ([]entity.PipeMessage, error)
}
