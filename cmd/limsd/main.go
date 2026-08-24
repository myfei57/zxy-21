package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"lims/internal/config"
	"lims/internal/server"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	cfg := config.LoadFromEnv()
	log.Printf("lims starting on %s with data dir %s", cfg.Addr, cfg.DataDir)
	if err := server.Run(ctx, cfg); err != nil {
		log.Fatalf("lims exited: %v", err)
	}
}
