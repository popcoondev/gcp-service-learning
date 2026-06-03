package main

import (
	"log"
	"net/http"
	"os"
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

	log.Printf("order api listening on :%s", addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}

func envOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
