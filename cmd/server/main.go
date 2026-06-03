package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpadapter "github.com/popcoondev/gcp-service-learning/internal/adapters/http"
	memoryrepo "github.com/popcoondev/gcp-service-learning/internal/adapters/repository/memory"
	"github.com/popcoondev/gcp-service-learning/internal/application"
)

func main() {
	addr := envOrDefault("PORT", "8080")

	repo := memoryrepo.NewOrderRepository()
	service := application.NewOrderService(repo, application.SystemClock{}, application.StaticIDGenerator{})
	handler := httpadapter.NewHandler(service)

	server := &http.Server{
		Addr:              ":" + addr,
		Handler:           handler.Routes(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	shutdownCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		<-shutdownCtx.Done()
		log.Printf("shutdown signal received")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("graceful shutdown failed: %v", err)
		}
	}()

	log.Printf("order api listening on :%s", addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}

	log.Printf("order api stopped")
}

func envOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
