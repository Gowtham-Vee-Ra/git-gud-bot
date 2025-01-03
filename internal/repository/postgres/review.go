package postgres

import (
	"context"
	"database/sql"
	"time"

	"git-gud-bot/internal/model"

	"github.com/google/uuid"
)

type ReviewRepository struct {
	db *sql.DB
}

func NewReviewRepository(db *sql.DB) *ReviewRepository {
	return &ReviewRepository{
		db: db,
	}
}

func (r *ReviewRepository) CreateReview(ctx context.Context, review *model.Review) error {
	query := `
		INSERT INTO reviews (
			id, pr_number, repo_owner, repo_name, status, title,
			description, feedback, commit_hash, code_quality,
			performance, best_practices, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`

	if review.ID == "" {
		review.ID = uuid.New().String()
	}

	now := time.Now()
	review.CreatedAt = now
	review.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		review.ID, review.PRNumber, review.RepoOwner, review.RepoName,
		review.Status, review.Title, review.Description, review.Feedback,
		review.CommitHash, review.CodeQuality, review.Performance,
		review.BestPractices, review.CreatedAt, review.UpdatedAt,
	)

	return err
}

func (r *ReviewRepository) GetReview(ctx context.Context, id string) (*model.Review, error) {
	query := `
		SELECT id, pr_number, repo_owner, repo_name, status, title,
			   description, feedback, commit_hash, code_quality,
			   performance, best_practices, created_at, updated_at
		FROM reviews
		WHERE id = $1
	`

	review := &model.Review{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&review.ID, &review.PRNumber, &review.RepoOwner, &review.RepoName,
		&review.Status, &review.Title, &review.Description, &review.Feedback,
		&review.CommitHash, &review.CodeQuality, &review.Performance,
		&review.BestPractices, &review.CreatedAt, &review.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return review, nil
}

func (r *ReviewRepository) GetReviews(ctx context.Context) ([]*model.Review, error) {
	query := `
		SELECT id, pr_number, repo_owner, repo_name, status, title,
			   description, feedback, commit_hash, code_quality,
			   performance, best_practices, created_at, updated_at
		FROM reviews
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []*model.Review
	for rows.Next() {
		review := &model.Review{}
		err := rows.Scan(
			&review.ID, &review.PRNumber, &review.RepoOwner, &review.RepoName,
			&review.Status, &review.Title, &review.Description, &review.Feedback,
			&review.CommitHash, &review.CodeQuality, &review.Performance,
			&review.BestPractices, &review.CreatedAt, &review.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}

	return reviews, nil
}
