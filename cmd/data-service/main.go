package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"

	"github.com/cityhunteur/data-service/internal/database"
	"github.com/cityhunteur/data-service/internal/handler"
)

var (
	logger        *zap.SugaredLogger
	migrationsDir = flag.String("migrationsDir", "migrations/", "migrations dir")
)

func main() {
	flag.Parse()

	zlog, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to initialise logger: %v", err)
	}
	defer func() { _ = zlog.Sync() }()
	logger = zlog.Sugar()

	var config database.Config
	err = envconfig.Process("", &config)
	if err != nil {
		log.Fatalf("failed to read database env config: %v", err)
	}

	db, err := database.NewFromEnv(context.Background(), &config)
	if err != nil {
		log.Fatalf("failed to initialise database: %v", err)
	}

	if err := runMigrations(&config); err != nil {
		log.Fatalf("failed to run database migrations: %v", err)
	}

	h, err := handler.NewDataHandler(context.Background(), database.NewDataDB(db))
	if err != nil {
		log.Fatalf("failed to create handler: %v", err)
	}

	// setup API routes
	router := gin.Default()
	router.GET("/v1/data", h.GetData)
	router.POST("/v1/data", h.PostData)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Infof("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown:", err)
	}

	logger.Infof("Server exiting")
}

func runMigrations(config *database.Config) error {
	dir := fmt.Sprintf("file://%s", *migrationsDir)
	m, err := migrate.New(dir, config.ConnectionURL())
	if err != nil {
		return fmt.Errorf("failed to create migrate: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to run migrate: %w", err)
	}
	srcErr, dbErr := m.Close()
	if srcErr != nil {
		return fmt.Errorf("source error: %w", srcErr)
	}
	if dbErr != nil {
		return fmt.Errorf("database error: %w", dbErr)
	}

	logger.Debugw("migrations run successfully")
	return nil
}
