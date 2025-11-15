package repo

import "github.com/iamyblitz/pr-reviewer-service/internal/model"

type Repository interface {
	// Teams
	CreateTeam(team *model.Team) error
	GetTeam(name string) (*model.Team, error)

	// Users
	SetUserActive(userID string, isActive bool) (*model.User, error)
	GetUserByID(userID string) (*model.User, error)

	// Pull Requests
	CreatePullRequest(pr *model.PullRequest) error
	GetPullRequestByID(id string) (*model.PullRequest, error)
	UpdatePullRequest(pr *model.PullRequest) error

	// Reviewers
	SetReviewers(prID string, reviewers []string) error
	GetReviewers(prID string) ([]string, error)
	GetPullRequestsByReviewer(userID string) ([]model.PullRequest, error)
}
