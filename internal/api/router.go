package api

import (
	"git-gud-bot/internal/api/handler"
	"git-gud-bot/internal/api/middleware"

	"github.com/gin-gonic/gin"
)

type Router struct {
	handler    *handler.ReviewHandler
	middleware *middleware.AuthMiddleware
}

func NewRouter(handler *handler.ReviewHandler, middleware *middleware.AuthMiddleware) *Router {
	return &Router{
		handler:    handler,
		middleware: middleware,
	}
}

func (r *Router) Setup(engine *gin.Engine) {
	// Health check endpoint (no auth required)
	engine.GET("/health", r.healthCheck)

	// API v1 group
	v1 := engine.Group("/api/v1")
	{
		// Public routes (no auth required)
		public := v1.Group("/")
		{
			public.GET("/status", r.getAPIStatus)
		}

		// Protected routes (auth required)
		protected := v1.Group("/")
		protected.Use(r.middleware.Authenticate())
		{
			// Review endpoints
			reviews := protected.Group("/reviews")
			{
				reviews.POST("/", r.handler.CreateReview)
				reviews.GET("/", r.handler.GetReviews)
				reviews.GET("/:id", r.handler.GetReview)
			}

			// Analysis endpoints
			analysis := protected.Group("/analysis")
			{
				analysis.POST("/analyze", r.analyzeCode)
				analysis.GET("/metrics", r.getMetrics)
				analysis.GET("/reports", r.getReports)
			}

			// GitHub webhook endpoints
			webhooks := protected.Group("/webhooks")
			{
				webhooks.POST("/github", r.handleGithubWebhook)
			}

			// User management endpoints
			users := protected.Group("/users")
			{
				users.GET("/me", r.getCurrentUser)
				users.PUT("/me", r.updateCurrentUser)
			}
		}
	}
}

// Health check handler
func (r *Router) healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
		"time":   "",
	})
}

// API status handler
func (r *Router) getAPIStatus(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "operational",
		"version": "v1.0.0",
	})
}

// Code analysis handler
func (r *Router) analyzeCode(c *gin.Context) {
	// TODO: Implement code analysis endpoint
	c.JSON(501, gin.H{"message": "Not implemented yet"})
}

// Metrics handler
func (r *Router) getMetrics(c *gin.Context) {
	// TODO: Implement metrics endpoint
	c.JSON(501, gin.H{"message": "Not implemented yet"})
}

// Reports handler
func (r *Router) getReports(c *gin.Context) {
	// TODO: Implement reports endpoint
	c.JSON(501, gin.H{"message": "Not implemented yet"})
}

// GitHub webhook handler
func (r *Router) handleGithubWebhook(c *gin.Context) {
	// TODO: Implement GitHub webhook handler
	c.JSON(501, gin.H{"message": "Not implemented yet"})
}

// User handlers
func (r *Router) getCurrentUser(c *gin.Context) {
	// TODO: Implement get current user endpoint
	c.JSON(501, gin.H{"message": "Not implemented yet"})
}

func (r *Router) updateCurrentUser(c *gin.Context) {
	// TODO: Implement update current user endpoint
	c.JSON(501, gin.H{"message": "Not implemented yet"})
}
