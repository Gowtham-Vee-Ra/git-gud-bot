package model

import (
	"time"
)

type ReviewStatus string

const (
	StatusPending  ReviewStatus = "pending"
	StatusApproved ReviewStatus = "approved"
	StatusRejected ReviewStatus = "rejected"
	StatusNeedWork ReviewStatus = "needs_work"
)

type Review struct {
	ID            string       `json:"id"`
	PRNumber      int          `json:"pr_number"`
	RepoOwner     string       `json:"repo_owner"`
	RepoName      string       `json:"repo_name"`
	Status        ReviewStatus `json:"status"`
	Title         string       `json:"title"`
	Description   string       `json:"description"`
	Feedback      string       `json:"feedback"`
	CommitHash    string       `json:"commit_hash"`
	CodeQuality   float64      `json:"code_quality"`
	Performance   float64      `json:"performance"`
	BestPractices float64      `json:"best_practices"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}

type ReviewRequest struct {
	PRNumber   int    `json:"pr_number" binding:"required"`
	RepoOwner  string `json:"repo_owner" binding:"required"`
	RepoName   string `json:"repo_name" binding:"required"`
	CommitHash string `json:"commit_hash" binding:"required"`
}

type ReviewResponse struct {
	Review  *Review `json:"review"`
	Message string  `json:"message,omitempty"`
	Error   string  `json:"error,omitempty"`
}
