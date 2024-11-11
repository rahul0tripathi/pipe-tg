package piped

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rahul0tripathi/pipetg/config"
	"github.com/rahul0tripathi/pipetg/controller"
	"github.com/rahul0tripathi/pipetg/internal/integrations/tg"
	"github.com/rahul0tripathi/pipetg/internal/services"
	"github.com/rahul0tripathi/pipetg/pkg/httpserver"
)

func Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.NewConfigFromEnv()
	if err != nil {
		return err
	}

	cli, err := tg.NewTelegramClient(cfg.UID, cfg.AppID, cfg.AppHash, cfg.SessionConfig)
	if err != nil {
		return err
	}

	scrapeWindow, err := time.ParseDuration(cfg.Window)
	if err != nil {
		return err
	}

	server := httpserver.NewServerWithMiddlewares(fmt.Sprintf(":%s", cfg.Port))
	controller.Router(
		server.Router(),
		cli,
		services.NewAuthFlowService(cli),
		services.NewScraper(cli, scrapeWindow),
	)
	server.Start()
	defer server.Shutdown()

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("ctx done")
			return nil
		case s := <-interrupt:
			fmt.Println("sig: " + s.String())
			return nil
		}
	}
}
