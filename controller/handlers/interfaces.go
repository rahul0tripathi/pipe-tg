package handlers

import (
	"context"
)

type AuthFlowSvc interface {
	RequestNewCode(ctx context.Context) error
	SubmitCode(ctx context.Context, code string) error
}

type MessageLogger interface {
	All(ctx context.Context) (interface{}, error)
}
