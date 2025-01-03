package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"git-gud-bot/internal/api"
	"git-gud-bot/internal/api/handler"
	"git-gud-bot/internal/api/middleware"
	"git-gud-bot/internal/config"
	"git-gud-bot/internal/repository/postgres"
	"git-gud-bot/internal/service"
	"git-gud-bot/pkg/analyzer"
	"git-gud-bot/pkg/github"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize configuration
	cfg := config.New()

	// Initialize router
	engine := gin.Default()

	// Initialize dependencies
	githubClient := github.NewClient(cfg.GithubToken)
	codeAnalyzer := analyzer.NewCodeAnalyzer(githubClient)

	// Initialize services and repositories
	reviewRepo := postgres.NewReviewRepository(cfg.DB)
	reviewService := service.NewReviewService(reviewRepo, githubClient, codeAnalyzer)
	reviewHandler := handler.NewReviewHandler(reviewService)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware()

	// Setup router with all dependencies
	router := api.NewRouter(reviewHandler, authMiddleware)
	router.Setup(engine)

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: engine,
	}

	// Server run context
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Start server
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v\n", err)
		}
	}()

	// Wait for interrupt signal
	<-ctx.Done()
	stop()
	log.Println("Shutting down gracefully, press Ctrl+C again to force")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
