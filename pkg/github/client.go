package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Client struct {
	httpClient *http.Client
	token      string
	baseURL    string
}

type PullRequest struct {
	Title       string `json:"title"`
	Description string `json:"body"`
	State       string `json:"state"`
	Head        struct {
		SHA string `json:"sha"`
	} `json:"head"`
	Files []File `json:"files"`
}

type File struct {
	Name        string `json:"filename"`
	Status      string `json:"status"`
	Additions   int    `json:"additions"`
	Deletions   int    `json:"deletions"`
	Changes     int    `json:"changes"`
	ContentsURL string `json:"contents_url"`
	PatchURL    string `json:"patch_url"`
}

func NewClient(token string) *Client {
	return &Client{
		httpClient: &http.Client{},
		token:      token,
		baseURL:    "https://api.github.com",
	}
}

func (c *Client) GetPullRequest(ctx context.Context, owner, repo string, number int) (*PullRequest, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/pulls/%d", c.baseURL, owner, repo, number)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API error: %s, status: %d", string(body), resp.StatusCode)
	}

	var pr PullRequest
	if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Get PR files
	files, err := c.GetPullRequestFiles(ctx, owner, repo, number)
	if err != nil {
		return nil, fmt.Errorf("failed to get PR files: %w", err)
	}
	pr.Files = files

	return &pr, nil
}

func (c *Client) GetPullRequestFiles(ctx context.Context, owner, repo string, number int) ([]File, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/pulls/%d/files", c.baseURL, owner, repo, number)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API error: %s, status: %d", string(body), resp.StatusCode)
	}

	var files []File
	if err := json.NewDecoder(resp.Body).Decode(&files); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return files, nil
}

func (c *Client) CreateReviewComment(ctx context.Context, owner, repo string, number int, comment *ReviewComment) error {
	url := fmt.Sprintf("%s/repos/%s/%s/pulls/%d/comments", c.baseURL, owner, repo, number)

	jsonBody, err := json.Marshal(comment)
	if err != nil {
		return fmt.Errorf("failed to marshal comment: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(jsonBody)))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		responseBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("GitHub API error: %s, status: %d", string(responseBody), resp.StatusCode)
	}

	return nil
}

type ReviewComment struct {
	Body     string `json:"body"`
	Path     string `json:"path"`
	Position int    `json:"position"`
	CommitID string `json:"commit_id"`
}

func (c *Client) setHeaders(req *http.Request) {
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Content-Type", "application/json")
}
