package app

import (
	"context"
	"errors"
	"github.com/alextotalk/atanika/internal/server"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Run initializes whole application.
func Run() {

	// HTTP Server
	srv := server.NewServer()

	go func() {
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			slog.Error("error occurred while running http server: %s\n", err.Error())
		}
	}()
	slog.Info("Server started")

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		slog.Error("failed to stop server: %v", err)
	}

	//if err := mongoClient.Disconnect(context.Background()); err != nil {
	//	slog.Error(err.Error())
	//}
}
