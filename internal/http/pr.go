package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/iamyblitz/pr-reviewer-service/internal/service"
)

type CreatePRRequest struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
}

type PullRequestDTO struct {
	ID                string   `json:"pull_request_id"`
	Name              string   `json:"pull_request_name"`
	AuthorID          string   `json:"author_id"`
	Status            string   `json:"status"`
	AssignedReviewers []string `json:"assigned_reviewers"`
	CreatedAt         *string  `json:"createdAt,omitempty"`
	MergedAt          *string  `json:"mergedAt,omitempty"`
}

func (h *Handler) CreatePullRequest(w http.ResponseWriter, r *http.Request) {
	var req CreatePRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if req.PullRequestID == "" || req.PullRequestName == "" || req.AuthorID == "" {
		http.Error(w, "pull_request_id, pull_request_name and author_id are required", http.StatusBadRequest)
		return
	}

	pr, err := h.svc.CreatePullRequest(req.PullRequestID, req.PullRequestName, req.AuthorID)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"error": map[string]any{
					"code":    "NOT_FOUND",
					"message": "resource not found",
				},
			})
			return
		}
		if errors.Is(err, service.ErrPRExists) {
			w.WriteHeader(http.StatusConflict)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"error": map[string]any{
					"code":    "PR_EXISTS",
					"message": "PR id already exists",
				},
			})
			return
		}

		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	var createdAtStr *string
	if pr.CreatedAt != nil {
		s := pr.CreatedAt.Format(time.RFC3339)
		createdAtStr = &s
	}

	var mergedAtStr *string
	if pr.MergedAt != nil {
		s := pr.MergedAt.Format(time.RFC3339)
		mergedAtStr = &s
	}

	resp := map[string]any{
		"pr": PullRequestDTO{
			ID:                pr.ID,
			Name:              pr.Name,
			AuthorID:          pr.AuthorID,
			Status:            string(pr.Status),
			AssignedReviewers: pr.AssignedReviewers,
			CreatedAt:         createdAtStr,
			MergedAt:          mergedAtStr,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}

type MergePRRequest struct {
	PullRequestID string `json:"pull_request_id"`
}

func (h *Handler) MergePullRequest(w http.ResponseWriter, r *http.Request) {
	var req MergePRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if req.PullRequestID == "" {
		http.Error(w, "pull_request_id is required", http.StatusBadRequest)
		return
	}

	pr, err := h.svc.MergePullRequest(req.PullRequestID)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"error": map[string]any{
					"code":    "NOT_FOUND",
					"message": "resource not found",
				},
			})
			return
		}

		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	var createdAtStr *string
	if pr.CreatedAt != nil {
		s := pr.CreatedAt.Format(time.RFC3339)
		createdAtStr = &s
	}

	var mergedAtStr *string
	if pr.MergedAt != nil {
		s := pr.MergedAt.Format(time.RFC3339)
		mergedAtStr = &s
	}

	resp := map[string]any{
		"pr": PullRequestDTO{
			ID:                pr.ID,
			Name:              pr.Name,
			AuthorID:          pr.AuthorID,
			Status:            string(pr.Status),
			AssignedReviewers: pr.AssignedReviewers,
			CreatedAt:         createdAtStr,
			MergedAt:          mergedAtStr,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
