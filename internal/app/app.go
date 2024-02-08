package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/alextotalk/atanika/internal/config"
	"github.com/alextotalk/atanika/internal/server"
	"github.com/alextotalk/atanika/internal/storage/pg"
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

// Run initializes the entire application.
func Run() {
	// Load configuration
	cfg := config.MustLoad()

	// Set up logger
	log := setupLogger(cfg.Env)
	log.Info(
		"starting url-shortener",
		slog.String("env", cfg.Env),
		slog.String("version", "123"),
	)
	log.Debug("debug messages are enabled")

	// Initialize database
	db, err := pg.New(pg.Config{
		Host:     cfg.PgHost,
		Port:     cfg.PgPort,
		Username: cfg.PgUser,
		DBName:   cfg.PgName,
		SSLMode:  cfg.SSLMode,
		Password: cfg.PgPassword,
	})
	fmt.Printf("db: %v\n", db)
	_ = db
	if err != nil {
		log.Error("Failed to initialize database: %s", err)
		os.Exit(1)
	}

	// Initialize HTTP server
	srv := server.NewServer()

	// Start HTTP server
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

	// Stop the server gracefully
	if err := srv.Stop(ctx); err != nil {
		slog.Error("Failed to stop server: %v", err)
	}

	if err := db.Close(); err != nil {
		slog.Error("Failed to close database: %v", err)
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
