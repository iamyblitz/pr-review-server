package repo

import (
	"errors"
	"sync"

	"github.com/iamyblitz/pr-reviewer-service/internal/model"
)

var (
	ErrNotFound   = errors.New("not found")
	ErrTeamExists = errors.New("team already exists")
)

type MemoryRepo struct {
	mu sync.RWMutex

	teams     map[string]*model.Team        // по team_name
	users     map[string]*model.User        // по user_id
	prs       map[string]*model.PullRequest // по pull_request_id
	reviewers map[string][]string           // по pull_request_id
}

func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{
		teams:     make(map[string]*model.Team),
		users:     make(map[string]*model.User),
		prs:       make(map[string]*model.PullRequest),
		reviewers: make(map[string][]string),
	}
}

func (m *MemoryRepo) CreateTeam(team *model.Team) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.teams[team.Name]; exists {
		return ErrTeamExists
	}
	m.teams[team.Name] = team

	for _, u := range team.Members {
		user := u // копия
		user.TeamName = team.Name
		m.users[user.ID] = &user
	}

	return nil

}

func (m *MemoryRepo) GetTeam(name string) (*model.Team, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	team, ok := m.teams[name]
	if !ok {
		return nil, ErrNotFound
	}

	copyTeam := *team
	copyMembers := make([]model.User, len(team.Members))
	copy(copyMembers, team.Members)
	copyTeam.Members = copyMembers

	return &copyTeam, nil
}
func (m *MemoryRepo) SetUserActive(userID string, isActive bool) (*model.User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	u, ok := m.users[userID]
	if !ok {
		return nil, ErrNotFound
	}

	u.IsActive = isActive

	copyUser := *u
	return &copyUser, nil
}

func (m *MemoryRepo) GetUserByID(userID string) (*model.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	u, ok := m.users[userID]
	if !ok {
		return nil, ErrNotFound
	}

	copyUser := *u
	return &copyUser, nil
}

func (m *MemoryRepo) CreatePullRequest(pr *model.PullRequest) error { panic("not implemented") }
func (m *MemoryRepo) GetPullRequestByID(id string) (*model.PullRequest, error) {
	panic("not implemented")
}
func (m *MemoryRepo) UpdatePullRequest(pr *model.PullRequest) error      { panic("not implemented") }
func (m *MemoryRepo) SetReviewers(prID string, reviewers []string) error { panic("not implemented") }
func (m *MemoryRepo) GetReviewers(prID string) ([]string, error)         { panic("not implemented") }
func (m *MemoryRepo) GetPullRequestsByReviewer(userID string) ([]model.PullRequest, error) {
	panic("not implemented")
}
