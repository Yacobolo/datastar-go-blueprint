// Package main is the entry point for the Datastar Go Blueprint server.
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yacobolo/datastar-go-blueprint/internal/app"
	"github.com/yacobolo/datastar-go-blueprint/internal/config"
	"github.com/yacobolo/datastar-go-blueprint/internal/platform/router"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/sync/errgroup"
)

func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: config.Global.LogLevel,
	}))
	slog.SetDefault(logger)

	if err := run(ctx, logger); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("error running server", "error", err)
		cancel()
		os.Exit(1)
	}
	cancel()
}

func run(ctx context.Context, logger *slog.Logger) error {

	addr := fmt.Sprintf("%s:%s", config.Global.Host, config.Global.Port)
	logger.Info("server started", "addr", addr)
	defer logger.Info("server shutdown complete")

	eg, egctx := errgroup.WithContext(ctx)

	r := chi.NewMux()
	r.Use(
		middleware.Logger,
		middleware.Recoverer,
	)

	// Wire up application dependencies
	application, err := app.New(config.Global, logger)
	if err != nil {
		return fmt.Errorf("failed to initialize app: %w", err)
	}
	defer func() { _ = application.Close() }()

	// Setup routes with App container
	if err := router.SetupRoutes(egctx, r, application); err != nil {
		return fmt.Errorf("error setting up routes: %w", err)
	}

	srv := &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
		BaseContext: func(_ net.Listener) context.Context {
			return egctx
		},
		ErrorLog: slog.NewLogLogger(
			slog.Default().Handler(),
			slog.LevelError,
		),
	}

	eg.Go(func() error {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("server error: %w", err)
		}
		return nil
	})

	eg.Go(func() error {
		<-egctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		logger.Debug("shutting down server...")

		if err := srv.Shutdown(shutdownCtx); err != nil {
			logger.Error("error during shutdown", "error", err)
			return err
		}

		return nil
	})

	return eg.Wait()
}
