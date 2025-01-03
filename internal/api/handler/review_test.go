package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"git-gud-bot/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockReviewService is a mock implementation of the ReviewService
type MockReviewService struct {
	mock.Mock
}

func (m *MockReviewService) CreateReview(review *model.Review) error {
	args := m.Called(review)
	return args.Error(0)
}

func (m *MockReviewService) GetReview(id string) (*model.Review, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Review), args.Error(1)
}

func (m *MockReviewService) GetReviews() ([]*model.Review, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Review), args.Error(1)
}

func setupTestRouter(h *ReviewHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Setup routes
	router.POST("/reviews", h.CreateReview)
	router.GET("/reviews", h.GetReviews)
	router.GET("/reviews/:id", h.GetReview)

	return router
}

func TestCreateReview(t *testing.T) {
	mockService := new(MockReviewService)
	handler := NewReviewHandler(mockService)
	router := setupTestRouter(handler)

	review := &model.Review{
		PRNumber: 123,
		RepoName: "test-repo",
		Status:   "pending",
		Feedback: "Test feedback",
	}

	mockService.On("CreateReview", mock.AnythingOfType("*model.Review")).Return(nil)

	body, _ := json.Marshal(review)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/reviews", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}

func TestGetReview(t *testing.T) {
	mockService := new(MockReviewService)
	handler := NewReviewHandler(mockService)
	router := setupTestRouter(handler)

	review := &model.Review{
		ID:       "test-id",
		PRNumber: 123,
		RepoName: "test-repo",
	}

	mockService.On("GetReview", "test-id").Return(review, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/reviews/test-id", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)

	var response map[string]*model.Review
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, review.ID, response["review"].ID)
}

func TestGetReviews(t *testing.T) {
	mockService := new(MockReviewService)
	handler := NewReviewHandler(mockService)
	router := setupTestRouter(handler)

	reviews := []*model.Review{
		{
			ID:       "test-id-1",
			PRNumber: 123,
			RepoName: "test-repo-1",
		},
		{
			ID:       "test-id-2",
			PRNumber: 124,
			RepoName: "test-repo-2",
		},
	}

	mockService.On("GetReviews").Return(reviews, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/reviews", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response["reviews"].([]interface{}), 2)
}
