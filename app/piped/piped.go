package piped

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/rahul0tripathi/pipetg/config"
	"github.com/rahul0tripathi/pipetg/controller"
	"github.com/rahul0tripathi/pipetg/internal/integrations"
	"github.com/rahul0tripathi/pipetg/internal/integrations/tg"
	"github.com/rahul0tripathi/pipetg/internal/services"
	"github.com/rahul0tripathi/pipetg/pkg/httpserver"
	"github.com/rahul0tripathi/pipetg/pkg/log"
)

type App struct {
	cfg     *config.Config
	client  *tg.Client
	scraper *services.Scraper
	sink    *integrations.MessageSink
	server  *httpserver.Server
}

func newApp(cfg *config.Config) (*App, error) {
	client, err := tg.NewTelegramClient(cfg.UID, cfg.AppID, cfg.AppHash, cfg.SessionConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create telegram client: %w", err)
	}

	scrapeWindow, err := time.ParseDuration(cfg.Window)
	if err != nil {
		return nil, fmt.Errorf("invalid scrape window duration: %w", err)
	}

	return &App{
		cfg:     cfg,
		client:  client,
		scraper: services.NewScraper(client, scrapeWindow),
		sink:    integrations.NewDummyMessageSink(),
	}, nil
}

func (a *App) setupHTTPServer() {
	a.server = httpserver.NewServerWithMiddlewares(fmt.Sprintf(":%s", a.cfg.Port))
	controller.Router(
		a.server.Router(),
		a.client,
		services.NewAuthFlowService(a.client),
		a.scraper,
	)
}

func (a *App) runScraperWorker(ctx context.Context, wg *sync.WaitGroup) {
	log.GetLogger(ctx).Info("starting scraper")
	wg.Add(1)
	go func() {
		defer wg.Done()
		a.scraperWorker(ctx)
	}()
}

func (a *App) scraperWorker(ctx context.Context) {
	scrapeWindow, _ := time.ParseDuration(a.cfg.Window) // Error already checked in newApp
	ticker := time.NewTicker(scrapeWindow)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := a.runScrape(ctx); err != nil {
				fmt.Printf("Error during scheduled scrape: %v\n", err)
			}
		}
	}
}

func (a *App) runScrape(ctx context.Context) error {
	return a.client.WithContext(ctx, func(ctx context.Context) error {
		messages, err := a.scraper.Run(ctx)
		if err != nil {
			return fmt.Errorf("scraper run failed: %w", err)
		}

		if err := a.sink.Collect(ctx, messages); err != nil {
			return fmt.Errorf("message collection failed: %w", err)
		}

		return nil
	})
}

func (a *App) waitForShutdown(ctx context.Context, cancel context.CancelFunc) error {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		fmt.Println("Context cancelled, shutting down")
		return nil
	case s := <-interrupt:
		fmt.Printf("Received signal: %s, shutting down\n", s.String())
		cancel()
		return nil
	}
}

func (a *App) shutdown() {
	if a.server != nil {
		a.server.Shutdown()
	}
}

func Run() error {
	cfg, err := config.NewConfigFromEnv()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	app, err := newApp(cfg)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	app.runScraperWorker(ctx, &wg)

	app.setupHTTPServer()
	app.server.Start()
	defer app.shutdown()

	if err := app.waitForShutdown(ctx, cancel); err != nil {
		return fmt.Errorf("shutdown error: %w", err)
	}

	wg.Wait()
	return nil
}
