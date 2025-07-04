package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Nikola-Milovic/vyking-interview/internal/cache/memory"
	"github.com/Nikola-Milovic/vyking-interview/internal/clients"
	"github.com/Nikola-Milovic/vyking-interview/internal/config"
	"github.com/Nikola-Milovic/vyking-interview/internal/service"
	"github.com/Nikola-Milovic/vyking-interview/internal/store"
	httpTransport "github.com/Nikola-Milovic/vyking-interview/internal/transport/http"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	if err := run(); err != nil {
		slog.Error("application failed", "error", err)
		os.Exit(1)
	}
}

func run() (err error) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	slog.Info("starting vyking player activity service")

	cfg, err := config.LoadFromEnv()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	slog.Info("configuration loaded",
		"db_host", cfg.DB.Host,
		"db_name", cfg.DB.Name,
		"server_port", cfg.Server.Port,
		"cache_size", cfg.Cache.Size,
		"cache_ttl", cfg.Cache.TTL)

	slog.Info("connecting to database")
	db, err := sql.Open("mysql", cfg.DB.DSN())
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}
	slog.Info("database connection established")

	cache := memory.New(cfg.Cache.Size, cfg.Cache.TTL)
	store := store.New(db)
	countryClient := clients.NewRestCountriesClient(cache, cfg.Cache.TTL)
	svc := service.New(store, countryClient)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		Handler:      newHTTPHandler(svc),
	}
	srvErr := make(chan error, 1)
	go func() {
		slog.Info("starting HTTP server", "addr", srv.Addr)
		srvErr <- srv.ListenAndServe()
	}()

	select {
	case err = <-srvErr:
		return
	case <-ctx.Done():
		slog.Info("received interrupt signal, shutting down gracefully")
		stop()
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	slog.Info("shutting down HTTP server")
	err = srv.Shutdown(shutdownCtx)
	return
}

func newHTTPHandler(svc service.Service) http.Handler {
	mux := http.NewServeMux()

	handler := httpTransport.NewHandler(svc)
	handler.RegisterRoutes(mux)

	return mux
}
