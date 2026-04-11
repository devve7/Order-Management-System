package main

import (
	application_order "Order-Management-System/internal/application/order"
	"Order-Management-System/internal/infractructure/postgres"
	inmemory_product_service "Order-Management-System/internal/infractructure/product/inmemory"
	postgres_storage "Order-Management-System/internal/infractructure/storage/postgres"
	transport_http "Order-Management-System/internal/transport/http"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint:     true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	if err := godotenv.Load(); err != nil {
		logger.Fatalf("failed to load .env: %v", err)
	}

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		logger.Fatal("DB_DSN is empty")
	}

	ctx := context.Background()

	pool, err := postgres.NewPool(ctx, dsn)
	if err != nil {
		logger.Fatalf("db connection failed: %v", err)
	}
	defer pool.Close()

	productService := inmemory_product_service.NewInMemoryService()
	productService.AddProduct(1, "Iphone", 1000, 1000)

	repo := postgres_storage.NewPostgresOrderRepository(pool)
	usecase := application_order.NewUseCase(productService, repo)

	handler := transport_http.NewHandler(usecase, logger)
	router := transport_http.NewRouter(handler, logger)

	server := transport_http.NewServer(":9091", router, logger)
	go func() {
		if err := server.Start(); err != nil {
			logger.Errorf("server error: %v", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	sig := <-quit
	logger.Printf("shutdown signal received: %v", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Errorf("graceful shutdown failed: %v", err)
	} else {
		logger.Println("server stopped gracefully")
	}
}
