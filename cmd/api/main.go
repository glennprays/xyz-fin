package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/glennprays/xyz-fin/config"
	"github.com/glennprays/xyz-fin/internal/app/database"
	"github.com/glennprays/xyz-fin/internal/app/handler"
	"github.com/glennprays/xyz-fin/internal/app/middleware"
	"github.com/glennprays/xyz-fin/internal/app/repository"
	"github.com/glennprays/xyz-fin/internal/app/router"
	"github.com/glennprays/xyz-fin/internal/app/service"
	"github.com/glennprays/xyz-fin/internal/app/usecase"
	"github.com/glennprays/xyz-fin/pkg/auth"
	"github.com/glennprays/xyz-fin/pkg/hasher"
)

func main() {
	log.Println("starting application...")

	log.Println("loading configuration...")
	cfg := config.LoadConfig()

	log.Println("initializing database connection...")
	dbConfig := database.Config{
		DSN: fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Jakarta",
			cfg.DBHost,
			cfg.DBPort,
			cfg.DBUser,
			cfg.DBPassword,
			cfg.DBName,
		),
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		ConnMaxLifetime: time.Hour,
		LogQueries:      true,
	}

	db, err := database.NewConnection(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer func() {
		log.Println("closing database connection...")
		if err := db.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	log.Println("initializing JWT Manager...")
	jwtManager := auth.NewJWTManager(
		cfg.JWTAccessSecret,
		cfg.JWTRefreshSecret,
		cfg.AppName,
		cfg.JWTAccessTokenDurationMinutes,
		cfg.JWTRefreshTokenDurationMinutes,
	)

	log.Println("initializing password hasher...")
	argonHasher := hasher.NewArgon2IDHasher()

	log.Println("Initializing repositories...")
	consumerRepo := repository.NewConsumerRepository(db)
	consumerLimitRepo := repository.NewConsumerLimitRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)

	log.Println("initializing services...")
	transactionService := service.NewTransactionService()

	log.Println("initializing usecases...")
	consumerUsecase := usecase.NewConsumerUsecase(consumerRepo, jwtManager, *argonHasher)
	transactionUsecase := usecase.NewTransactionUsecase(
		db,
		transactionService,
		transactionRepo,
		consumerRepo,
		consumerLimitRepo,
	)
	consumerLimitUsecase := usecase.NewConsumerLimitUsecase(consumerRepo, consumerLimitRepo)

	log.Println("initializing middleware...")
	authMiddleware := middleware.NewAuthMiddleware(jwtManager)

	log.Println("initializing handlers...")
	consumerHandler := handler.NewConsumerHandler(consumerUsecase)
	transactionHandler := handler.NewTransactionHandler(transactionUsecase)
	consumerLimitHandler := handler.NewConsumerLimitHandler(consumerLimitUsecase)

	log.Println("setting up router...")
	routerEngine := router.SetupRouter(
		authMiddleware,
		consumerHandler,
		transactionHandler,
		consumerLimitHandler,
	)

	log.Println("setting up HTTP server...")
	server := &http.Server{
		Addr:         ":" + cfg.AppPort,
		Handler:      routerEngine,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("starting server on port %s...", cfg.AppPort)

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
		log.Println("HTTP server stopped.")
	}()
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutdown signal received, initiating graceful shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("server exiting gracefully...")
}
