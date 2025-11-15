package repo

import "github.com/iamyblitz/pr-reviewer-service/internal/model"

type MemoryRepo struct {
}

func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{}
}

func (m *MemoryRepo) CreateTeam(team *model.Team) error {
	panic("not implemented")
}

func (m *MemoryRepo) GetTeam(name string) (*model.Team, error) { panic("not implemented") }
func (m *MemoryRepo) SetUserActive(userID string, isActive bool) (*model.User, error) {
	panic("not implemented")
}
func (m *MemoryRepo) GetUserByID(userID string) (*model.User, error) { panic("not implemented") }
func (m *MemoryRepo) CreatePullRequest(pr model.PullRequest) error   { panic("not implemented") }
func (m *MemoryRepo) GetPullRequestByID(id string) (*model.PullRequest, error) {
	panic("not implemented")
}
func (m *MemoryRepo) UpdatePullRequest(pr model.PullRequest) error       { panic("not implemented") }
func (m *MemoryRepo) SetReviewers(prID string, reviewers []string) error { panic("not implemented") }
func (m *MemoryRepo) GetReviewers(prID string) ([]string, error)         { panic("not implemented") }
func (m *MemoryRepo) GetPullRequestsByReviewer(userID string) ([]model.PullRequest, error) {
	panic("not implemented")
}
