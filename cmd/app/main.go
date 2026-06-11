package main

import (
	application_auth "Order-Management-System/internal/application/auth"
	application_order "Order-Management-System/internal/application/order"
	application_product "Order-Management-System/internal/application/product"
	health "Order-Management-System/internal/health"
	"Order-Management-System/internal/infractructure/cache"
	"Order-Management-System/internal/infractructure/db"
	"Order-Management-System/internal/infractructure/security"
	postgres_storage "Order-Management-System/internal/infractructure/storage/postgres"
	transport_http "Order-Management-System/internal/transport/http"
	transport_http_auth "Order-Management-System/internal/transport/http/auth"
	transport_health "Order-Management-System/internal/transport/http/health"
	"Order-Management-System/internal/transport/http/middleware"
	transport_http_order "Order-Management-System/internal/transport/http/order"
	transport_http_product "Order-Management-System/internal/transport/http/product"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint:     true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		logger.Fatal("DB_DSN is empty")
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		logger.Fatal("JWT_SECRET is empty")
	}

	jwtIssuer := os.Getenv("JWT_ISSUER")
	if jwtIssuer == "" {
		logger.Fatal("JWT_ISSUER is empty")
	}

	jwtAccessTTLRaw := os.Getenv("JWT_ACCESS_TTL")
	if jwtAccessTTLRaw == "" {
		logger.Fatal("JWT_ACCESS_TTL is empty")
	}

	jwtAccessTTL, err := time.ParseDuration(jwtAccessTTLRaw)
	if err != nil {
		logger.Fatalf("invalid JWT_ACCESS_TTL: %v", err)
	}

	ctx := context.Background()

	pool, err := db.NewPool(ctx, dsn)
	if err != nil {
		logger.Fatalf("db connection failed: %v", err)
	}
	defer pool.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,

		DialTimeout:  100 * time.Millisecond,
		ReadTimeout:  100 * time.Millisecond,
		WriteTimeout: 100 * time.Millisecond,

		MaxRetries:      0,
		MinRetryBackoff: -1,
		MaxRetryBackoff: -1,
	})

	postgresChecker := health.NewPostgresChecker(pool)
	redisChecker := health.NewRedisChecker(rdb)
	healthChecker := health.NewChecker(redisChecker, postgresChecker)
	healthHandler := transport_health.NewHealthHandler(healthChecker)

	redisCache := cache.NewRedisCache(rdb)
	loggingCache := cache.NewLoggingCache(redisCache, logger)

	userRepo := postgres_storage.NewPostgresUserRepository(pool)
	passwordHasher := security.NewBcryptPasswordHasher()

	tokenManager, err := security.NewJWTManager(
		jwtSecret,
		jwtAccessTTL,
		jwtIssuer,
	)
	if err != nil {
		logger.Fatalf("jwt manager init failed: %v", err)
	}

	authService := application_auth.NewAuthService(
		userRepo,
		passwordHasher,
		tokenManager,
	)

	authHandler := transport_http_auth.NewAuthHandler(authService, logger)

	authMiddleware := middleware.NewAuthMiddleware(tokenManager)

	productRepo := postgres_storage.NewPostgresProductRepository(pool)
	productUseCase := application_product.NewUseCase(productRepo, loggingCache)

	orderRepo := postgres_storage.NewPostgresOrderRepository(pool)
	orderUseCase := application_order.NewUseCase(productUseCase, orderRepo)

	orderHandler := transport_http_order.NewOrderHandler(orderUseCase, logger)
	productHandler := transport_http_product.NewProductHandler(productUseCase, logger)
	router := transport_http.NewRouter(orderHandler, productHandler, healthHandler, authHandler, authMiddleware, logger)

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
