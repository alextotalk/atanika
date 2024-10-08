package app

import (
	"context"
	"errors"
	"github.com/alextotalk/atanika/internal/config"
	"github.com/alextotalk/atanika/internal/handler"
	"github.com/alextotalk/atanika/internal/server"
	"github.com/alextotalk/atanika/internal/service"
	"github.com/alextotalk/atanika/internal/storage"
	"github.com/alextotalk/atanika/internal/storage/pg"

	_ "github.com/jackc/pgx/v5"
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
		"starting atanika",
		slog.String("env", cfg.Env),
	)
	log.Debug("debug messages are enabled")

	// Initialize database
	pgDB, err := pg.New(pg.Config{
		Host:     cfg.PgHost,
		Port:     cfg.PgPort,
		Username: cfg.PgUser,
		DBName:   cfg.PgName,
		SSLMode:  cfg.SSLMode,
		Password: cfg.PgPassword,
	})
	if err != nil {
		log.Error("Failed to initialize database: %s", err)
		os.Exit(1)
	}
	defer pgDB.Close()

	// NewRepository creates a new repository using the provided pgDB connection.
	// It returns the created repository.
	repos := storage.NewRepository(pgDB)

	services := service.NewService(repos)

	handlers := handler.NewHandler(services)

	// Initialize HTTP server
	srv := server.NewServer(handlers.NewRouter()) // Створення Server з екземпляром Echo

	// Start HTTP server
	go func() {
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Error occurred while running http server: %s\n", err.Error())
		}
	}()

	slog.Info("Server started on port:", cfg.Http.Port, "successful")

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	// Stop the server gracefully
	if err := srv.Stop(ctx); err != nil {
		log.Error("Failed to stop server: %v", err)
	}

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
