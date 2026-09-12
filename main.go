package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/yogiatwork/expense-tracker/config"
	"github.com/yogiatwork/expense-tracker/middleware"
	"github.com/yogiatwork/expense-tracker/route"
	"github.com/yogiatwork/expense-tracker/service"
)

// main is the entry point
// it loads the env file
// creates a new gin router
// starts the server on port 8080
func main() {
	// create a new logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	// create a new gin router
	r := gin.New()
	r.Use(middleware.LoggingMiddleware(logger))
	r.Use(middleware.RecoveryMiddleware(logger))

	// cors configuration
	r.Use(cors.New(cors.Config{
		AllowMethods: []string{"GET", "POST"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept"},
		AllowOrigins: []string{"*"},
	}))

	// load config
	if err := config.LoadConfig(); err != nil {
		logger.Error("failed to load config", slog.String("error", err.Error()))
		os.Exit(1)
	}
	logger.Info("config loaded successfully")

	// init google sheet service
	service.InitSheetService()
	logger.Info("google sheet service initialized successfully")

	// api routes
	api := r.Group("/api/v1")

	// add more routes here
	route.HealthzRoute(api)
	route.WhatsAppRoute(api)
	route.TelegramRoute(api)
	route.GoogleSheetsRoute(api)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// create a new http server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// start the server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("listen: ", slog.String("error", err.Error()))
		}
	}()

	logger.Info("Server started on port 8080")

	<-ctx.Done()
	stop()
	logger.Info("Shutting down gracefully, press Ctrl+C again to force")

	if err := srv.Shutdown(context.Background()); err != nil {
		logger.Error("Server forced to shutdown: ", slog.String("error", err.Error()))
	}

	logger.Info("Server exiting")

}
