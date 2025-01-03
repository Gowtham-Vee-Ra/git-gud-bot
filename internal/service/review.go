package service

import (
	"context"

	"git-gud-bot/internal/model"
	"git-gud-bot/internal/repository/postgres"
	"git-gud-bot/pkg/analyzer"
	"git-gud-bot/pkg/github"
)

type ReviewService struct {
	repo     *postgres.ReviewRepository
	github   *github.Client
	analyzer *analyzer.CodeAnalyzer
}

func NewReviewService(
	repo *postgres.ReviewRepository,
	github *github.Client,
	analyzer *analyzer.CodeAnalyzer,
) *ReviewService {
	return &ReviewService{
		repo:     repo,
		github:   github,
		analyzer: analyzer,
	}
}

func (s *ReviewService) CreateReview(ctx context.Context, req *model.ReviewRequest) (*model.Review, error) {
	// Fetch PR details from GitHub
	prDetails, err := s.github.GetPullRequest(ctx, req.RepoOwner, req.RepoName, req.PRNumber)
	if err != nil {
		return nil, err
	}

	// Analyze code
	analysis, err := s.analyzer.AnalyzeCode(ctx, prDetails)
	if err != nil {
		return nil, err
	}

	// Create review object
	review := &model.Review{
		PRNumber:      req.PRNumber,
		RepoOwner:     req.RepoOwner,
		RepoName:      req.RepoName,
		Status:        model.StatusPending,
		Title:         prDetails.Title,
		Description:   prDetails.Description,
		CommitHash:    req.CommitHash,
		CodeQuality:   analysis.CodeQuality,
		Performance:   analysis.Performance,
		BestPractices: analysis.BestPractices,
	}

	// Save to database
	if err := s.repo.CreateReview(ctx, review); err != nil {
		return nil, err
	}

	return review, nil
}

func (s *ReviewService) GetReview(ctx context.Context, id string) (*model.Review, error) {
	return s.repo.GetReview(ctx, id)
}

func (s *ReviewService) GetReviews(ctx context.Context) ([]*model.Review, error) {
	return s.repo.GetReviews(ctx)
}
