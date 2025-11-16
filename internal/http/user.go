package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/iamyblitz/pr-reviewer-service/internal/service"
)

type SetIsActiveRequest struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type UserDTO struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

func (h *Handler) SetUserIsActive(w http.ResponseWriter, r *http.Request) {
	var req SetIsActiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if req.UserID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	user, err := h.svc.SetUserIsActive(req.UserID, req.IsActive)
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

	resp := map[string]any{
		"user": UserDTO{
			UserID:   user.ID,
			Username: user.Username,
			TeamName: user.TeamName,
			IsActive: user.IsActive,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

type PullRequestShortDTO struct {
	ID       string `json:"pull_request_id"`
	Name     string `json:"pull_request_name"`
	AuthorID string `json:"author_id"`
	Status   string `json:"status"`
}

type GetUserReviewResponse struct {
	UserID       string                `json:"user_id"`
	PullRequests []PullRequestShortDTO `json:"pull_requests"`
}

func (h *Handler) GetUserReviews(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	prs, err := h.svc.GetUserReviews(userID)
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

	resp := GetUserReviewResponse{
		UserID:       userID,
		PullRequests: make([]PullRequestShortDTO, 0, len(prs)),
	}

	for _, pr := range prs {
		resp.PullRequests = append(resp.PullRequests, PullRequestShortDTO{
			ID:       pr.ID,
			Name:     pr.Name,
			AuthorID: pr.AuthorID,
			Status:   string(pr.Status),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
