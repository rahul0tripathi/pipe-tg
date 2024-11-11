package integrations

import (
	"context"

	"github.com/rahul0tripathi/pipetg/entity"
	"github.com/rahul0tripathi/pipetg/pkg/log"
)

type MessageSink struct{}

func NewDummyMessageSink() *MessageSink {
	return nil
}

func (m *MessageSink) Collect(ctx context.Context, messages []entity.PipeMessage) error {
	logger := log.GetLogger(ctx)
	logger.Info("dummy sink", log.Int("message", len(messages)))
	return nil
}
