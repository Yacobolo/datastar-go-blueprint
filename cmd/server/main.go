package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yacobolo/datastar-go-starter-kit/internal/config"
	"github.com/yacobolo/datastar-go-starter-kit/internal/platform/router"
	"github.com/yacobolo/datastar-go-starter-kit/internal/store"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/sessions"
	embeddednats "github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"golang.org/x/sync/errgroup"
)

func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: config.Global.LogLevel,
	}))
	slog.SetDefault(logger)

	if err := run(ctx); err != nil && err != http.ErrServerClosed {
		slog.Error("error running server", "error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {

	addr := fmt.Sprintf("%s:%s", config.Global.Host, config.Global.Port)
	slog.Info("server started", "addr", addr)
	defer slog.Info("server shutdown complete")

	eg, egctx := errgroup.WithContext(ctx)

	r := chi.NewMux()
	r.Use(
		middleware.Logger,
		middleware.Recoverer,
	)

	sessionStore := sessions.NewCookieStore([]byte(config.Global.SessionSecret))
	sessionStore.MaxAge(86400 * 30)
	sessionStore.Options.Path = "/"
	sessionStore.Options.HttpOnly = true
	sessionStore.Options.Secure = false
	sessionStore.Options.SameSite = http.SameSiteLaxMode

	// Start embedded NATS server
	natsOpts := &embeddednats.Options{
		Host:      "localhost",
		Port:      4222,
		JetStream: true,
	}
	ns, err := embeddednats.NewServer(natsOpts)
	if err != nil {
		return fmt.Errorf("failed to start NATS: %w", err)
	}
	go ns.Start()
	if !ns.ReadyForConnections(4 * time.Second) {
		return fmt.Errorf("NATS not ready")
	}
	defer ns.Shutdown()

	slog.Info("NATS server started", "url", ns.ClientURL())

	// Connect to NATS
	nc, err := nats.Connect(ns.ClientURL())
	if err != nil {
		return fmt.Errorf("failed to connect to NATS: %w", err)
	}
	defer nc.Close()

	// Initialize database with new store package
	dbStore, err := store.Open(config.Global.DBPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer dbStore.Close()

	slog.Info("database initialized", "path", config.Global.DBPath)

	// Setup routes with NATS connection
	if err := router.SetupRoutes(egctx, r, sessionStore, dbStore.Queries(), nc); err != nil {
		return fmt.Errorf("error setting up routes: %w", err)
	}

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
		BaseContext: func(l net.Listener) context.Context {
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

		slog.Debug("shutting down server...")

		if err := srv.Shutdown(shutdownCtx); err != nil {
			slog.Error("error during shutdown", "error", err)
			return err
		}

		return nil
	})

	return eg.Wait()
}
