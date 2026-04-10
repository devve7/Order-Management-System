package main

import (
	application_order "Order-Management-System/internal/application/order"
	domain_order "Order-Management-System/internal/domain/order"
	inmemory_catalog "Order-Management-System/internal/infractructure/catalog/inmemory"
	inmemory_storage "Order-Management-System/internal/infractructure/storage/inmemory"
	transport_http "Order-Management-System/internal/transport/http"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

func fieldCatalog(catalog *inmemory_catalog.InMemoryOrderCatalog) {
	catalog.AddProduct(domain_order.Product{
		Name:  "Iphone 13",
		Price: 1000,
		ID:    1,
	})
	catalog.AddProduct(domain_order.Product{
		Name:  "Iphone 14",
		Price: 1200,
		ID:    2,
	})
	catalog.AddProduct(domain_order.Product{
		Name:  "Iphone 15",
		Price: 1600,
		ID:    3,
	})
}

func main() {
	catalog := inmemory_catalog.NewInMemoryOrderCatalog()
	fieldCatalog(catalog)
	factory := domain_order.NewOrderItemFactory(catalog)

	repo := inmemory_storage.NewInmemoryOrderRepository()
	usecase := application_order.NewUseCase(factory, repo)

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
