package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"lims/internal/config"
)

// Run starts the console server and shuts down gracefully on context cancel.
func Run(ctx context.Context, cfg config.Config) error {
	app, err := NewApp(cfg)
	if err != nil {
		return err
	}
	httpServer := &http.Server{
		Addr:              cfg.Addr,
		Handler:           app.Handler,
		ReadHeaderTimeout: 5 * time.Second,
	}
	errorChannel := make(chan error, 1)
	go func() {
		err := httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errorChannel <- err
		}
		close(errorChannel)
	}()
	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("graceful shutdown: %w", err)
		}
		return nil
	case err := <-errorChannel:
		return err
	}
}
