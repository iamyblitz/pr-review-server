package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/iamyblitz/pr-reviewer-service/internal/model"
	"github.com/iamyblitz/pr-reviewer-service/internal/service"
)

type TeamMemberDTO struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type TeamDTO struct {
	TeamName string          `json:"team_name"`
	Members  []TeamMemberDTO `json:"members"`
}

type Handler struct {
	svc *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

// POST /team/add
func (h *Handler) AddTeam(w http.ResponseWriter, r *http.Request) {
	var dto TeamDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	members := make([]model.User, 0, len(dto.Members))
	for _, m := range dto.Members {
		members = append(members, model.User{
			ID:       m.UserID,
			Username: m.Username,
			TeamName: dto.TeamName,
			IsActive: m.IsActive,
		})
	}

	team, err := h.svc.CreateTeam(dto.TeamName, members)
	if err != nil {
		if errors.Is(err, service.ErrTeamExists) {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"error": map[string]any{
					"code":    "TEAM_EXISTS",
					"message": "team_name already exists",
				},
			})
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	respMembers := make([]TeamMemberDTO, 0, len(team.Members))
	for _, m := range team.Members {
		respMembers = append(respMembers, TeamMemberDTO{
			UserID:   m.ID,
			Username: m.Username,
			IsActive: m.IsActive,
		})
	}

	resp := map[string]any{
		"team": TeamDTO{
			TeamName: team.Name,
			Members:  respMembers,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *Handler) GetTeam(w http.ResponseWriter, r *http.Request) {
	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		http.Error(w, "team_name is required", http.StatusBadRequest)
		return
	}

	team, err := h.svc.GetTeam(teamName)
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

	members := make([]TeamMemberDTO, 0, len(team.Members))
	for _, m := range team.Members {
		members = append(members, TeamMemberDTO{
			UserID:   m.ID,
			Username: m.Username,
			IsActive: m.IsActive,
		})
	}

	resp := TeamDTO{
		TeamName: team.Name,
		Members:  members,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
