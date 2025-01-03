package handler

import (
	"net/http"

	"git-gud-bot/internal/model"
	"git-gud-bot/internal/service"

	"github.com/gin-gonic/gin"
)

type ReviewHandler struct {
	service *service.ReviewService
}

func NewReviewHandler(service *service.ReviewService) *ReviewHandler {
	return &ReviewHandler{
		service: service,
	}
}

// CreateReview handles the creation of a new code review
func (h *ReviewHandler) CreateReview(c *gin.Context) {
	var req model.ReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ReviewResponse{
			Error: "Invalid request format: " + err.Error(),
		})
		return
	}

	review, err := h.service.CreateReview(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ReviewResponse{
			Error: "Failed to create review: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, model.ReviewResponse{
		Review:  review,
		Message: "Review created successfully",
	})
}

// GetReview handles fetching a single review by ID
func (h *ReviewHandler) GetReview(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, model.ReviewResponse{
			Error: "Review ID is required",
		})
		return
	}

	review, err := h.service.GetReview(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ReviewResponse{
			Error: "Failed to fetch review: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.ReviewResponse{
		Review: review,
	})
}

// GetReviews handles fetching all reviews
func (h *ReviewHandler) GetReviews(c *gin.Context) {
	reviews, err := h.service.GetReviews(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ReviewResponse{
			Error: "Failed to fetch reviews: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"reviews": reviews,
		"count":   len(reviews),
	})
}
