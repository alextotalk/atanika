package app

import (
	"context"
	"errors"
	"github.com/alextotalk/atanika/internal/config"
	"github.com/alextotalk/atanika/internal/server"
	_ "github.com/lib/pq"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// Run initializes whole application.
func Run() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log.Info(
		"starting url-shortener",
		slog.String("env", cfg.Env),
		slog.String("version", "123"),
	)
	log.Debug("debug messages are enabled")

	// HTTP Server
	srv := server.NewServer()

	go func() {
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Error occurred while running http server: %s\n", err.Error())
		}
	}()

	slog.Info("Server started on port: %s", cfg.Http.Port)

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		slog.Error("Failed to stop server: %v", err)
	}

	//if err := mongoClient.Disconnect(context.Background()); err != nil {
	//	slog.Error(err.Error())
	//}

}
func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
