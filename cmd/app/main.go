package main

import (
	application_order "Order-Management-System/internal/application/order"
	inmemory_product_service "Order-Management-System/internal/infractructure/product/inmemory"
	inmemory_storage "Order-Management-System/internal/infractructure/storage/inmemory"
	transport_http "Order-Management-System/internal/transport/http"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	productService := inmemory_product_service.NewInMemoryService()
	productService.AddProduct(1, "Iphone", 1000, 1000)
	repo := inmemory_storage.NewInmemoryOrderRepository()
	usecase := application_order.NewUseCase(productService, repo)

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint:     true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	handler := transport_http.NewHandler(usecase)
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
