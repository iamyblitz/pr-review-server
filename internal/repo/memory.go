package repo

import (
	"errors"
	"sync"

	"github.com/iamyblitz/pr-reviewer-service/internal/model"
)

var (
	ErrNotFound   = errors.New("not found")
	ErrTeamExists = errors.New("team already exists")
	ErrPRExists   = errors.New("pr already exists")
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

func (m *MemoryRepo) CreatePullRequest(pr *model.PullRequest) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.prs[pr.ID]; exists {
		return ErrPRExists
	}

	copyPR := *pr
	m.prs[pr.ID] = &copyPR

	reviewersCopy := make([]string, len(pr.AssignedReviewers))
	copy(reviewersCopy, pr.AssignedReviewers)
	m.reviewers[pr.ID] = reviewersCopy

	return nil
}
func (m *MemoryRepo) GetPullRequestByID(id string) (*model.PullRequest, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	pr, ok := m.prs[id]
	if !ok {
		return nil, ErrNotFound
	}

	copyPR := *pr
	reviewers := m.reviewers[id]
	reviewersCopy := make([]string, len(reviewers))
	copy(reviewersCopy, reviewers)
	copyPR.AssignedReviewers = reviewersCopy

	return &copyPR, nil
}

func (m *MemoryRepo) UpdatePullRequest(pr *model.PullRequest) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.prs[pr.ID]; !ok {
		return ErrNotFound
	}

	copyPR := *pr
	m.prs[pr.ID] = &copyPR

	reviewersCopy := make([]string, len(pr.AssignedReviewers))
	copy(reviewersCopy, pr.AssignedReviewers)
	m.reviewers[pr.ID] = reviewersCopy

	return nil
}
func (m *MemoryRepo) SetReviewers(prID string, reviewers []string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	pr, ok := m.prs[prID]
	if !ok {
		return ErrNotFound
	}

	reviewersCopy := make([]string, len(reviewers))
	copy(reviewersCopy, reviewers)
	m.reviewers[prID] = reviewersCopy

	pr.AssignedReviewers = reviewersCopy

	return nil
}

func (m *MemoryRepo) GetReviewers(prID string) ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	reviewers, ok := m.reviewers[prID]
	if !ok {
		return nil, ErrNotFound
	}

	reviewersCopy := make([]string, len(reviewers))
	copy(reviewersCopy, reviewers)

	return reviewersCopy, nil
}

func (m *MemoryRepo) GetPullRequestsByReviewer(userID string) ([]model.PullRequest, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []model.PullRequest

	for prID, pr := range m.prs {
		reviewers := m.reviewers[prID]
		for _, r := range reviewers {
			if r == userID {
				copyPR := *pr
				reviewersCopy := make([]string, len(reviewers))
				copy(reviewersCopy, reviewers)
				copyPR.AssignedReviewers = reviewersCopy
				result = append(result, copyPR)
				break
			}
		}
	}

	return result, nil
}
