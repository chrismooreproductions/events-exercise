package app

import (
	"context"
	"events-exercise/config"
	"events-exercise/internal/events"
	"events-exercise/internal/logger"
	"events-exercise/internal/transport"
	"fmt"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

type app struct {
	server *http.Server
	logger *slog.Logger
}

func Load(cfg config.Config) *app {
	logger := logger.Logger

	processor := events.NewStreamProcessor()
	// Create a stream reader and read the events - the output will be on processor.Result()
	events.NewStreamReader(processor, logger)

	handler := transport.NewHandler(logger, processor)

	return &app{
		server: &http.Server{
			Addr:         fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
			Handler:      handler.Mux,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
		},
		logger: logger,
	}
}

func (a *app) Run() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	a.logger.Info("starting server", "addr", a.server.Addr)

	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Error("server error", "error", err)
		}
	}()
	<-ctx.Done()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	a.server.Shutdown(shutdownCtx)
	a.logger.Info("server stopped")
}
