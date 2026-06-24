package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dmytroyunyk/MikrotikApi/config"
	"github.com/dmytroyunyk/MikrotikApi/internal/api"
	"github.com/dmytroyunyk/MikrotikApi/internal/firewall"
	"github.com/dmytroyunyk/MikrotikApi/internal/mikrotik"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	cfg, err := config.Load()
	if err != nil {
		slog.Error("error configuration", "err", err)
		os.Exit(1)
	}

	mt := mikrotik.NewClient(
		cfg.MikroTik.Host,
		cfg.MikroTik.Username,
		cfg.MikroTik.Password,
		cfg.MikroTik.Timeout,
		cfg.MikroTik.Retries,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := mt.Ping(ctx); err != nil {
		slog.Warn("router unavailable at startup", "err", err)
	} else {
		slog.Info("connection to MikroTik established", "host", cfg.MikroTik.Host)
	}

	mikrotikSvc := mikrotik.NewService(mt)
	fwService := firewall.NewServise(mt)

	fwHandler := firewall.NewHandler(fwService)

	handler := api.NewHandler(mikrotikSvc, mt, fwHandler)

	router := api.NewRouter(handler)

	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	go func() {
		slog.Info("the server is started", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("finishing work")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("completion error", "err", err)
	}
	slog.Info("the server is stopped")
}
