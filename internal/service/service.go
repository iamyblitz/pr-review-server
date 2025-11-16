package service

import (
	"errors"
	"math/rand"
	"time"

	"github.com/iamyblitz/pr-reviewer-service/internal/model"
	"github.com/iamyblitz/pr-reviewer-service/internal/repo"
)

type Service struct {
	repo repo.Repository
}

func NewService(r repo.Repository) *Service {
	return &Service{repo: r}
}

var (
	ErrTeamExists  = errors.New("team already exists")
	ErrNotFound    = errors.New("not found")
	ErrPRExists    = errors.New("pr already exists")
	ErrPRMerged    = errors.New("pr already merged")
	ErrNotAssigned = errors.New("reviewer not assigned to this PR")
	ErrNoCandidate = errors.New("no active replacement candidate in team")
)

func (s *Service) CreateTeam(name string, members []model.User) (*model.Team, error) {
	team := &model.Team{
		Name:    name,
		Members: members,
	}

	err := s.repo.CreateTeam(team)
	if err != nil {
		if errors.Is(err, repo.ErrTeamExists) {
			return nil, ErrTeamExists
		}
		return nil, err
	}

	return team, nil
}

func (s *Service) GetTeam(name string) (*model.Team, error) {
	team, err := s.repo.GetTeam(name)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return team, nil
}

func (s *Service) SetUserIsActive(userID string, isActive bool) (*model.User, error) {
	user, err := s.repo.SetUserActive(userID, isActive)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return user, nil
}

func (s *Service) CreatePullRequest(id, name, authorID string) (*model.PullRequest, error) {

	author, err := s.repo.GetUserByID(authorID)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	team, err := s.repo.GetTeam(author.TeamName)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	candidates := make([]model.User, 0, len(team.Members))
	for _, m := range team.Members {
		if m.ID == authorID {
			continue
		}
		if !m.IsActive {
			continue
		}
		candidates = append(candidates, m)
	}

	reviewerIDs := chooseReviewers(candidates, 2)

	now := time.Now().UTC()

	pr := &model.PullRequest{
		ID:                id,
		Name:              name,
		AuthorID:          authorID,
		Status:            model.PRStatusOpen,
		AssignedReviewers: reviewerIDs,
		CreatedAt:         &now,
		MergedAt:          nil,
	}

	if err := s.repo.CreatePullRequest(pr); err != nil {
		if errors.Is(err, repo.ErrPRExists) {
			return nil, ErrPRExists
		}
		return nil, err
	}

	return pr, nil
}

func chooseReviewers(candidates []model.User, max int) []string {
	if len(candidates) == 0 || max <= 0 {
		return nil
	}

	rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})

	if len(candidates) > max {
		candidates = candidates[:max]
	}

	result := make([]string, 0, len(candidates))
	for _, c := range candidates {
		result = append(result, c.ID)
	}
	return result
}

func (s *Service) MergePullRequest(prID string) (*model.PullRequest, error) {
	pr, err := s.repo.GetPullRequestByID(prID)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if pr.Status == model.PRStatusMerged {
		return pr, nil
	}

	now := time.Now().UTC()
	pr.Status = model.PRStatusMerged
	pr.MergedAt = &now

	if err := s.repo.UpdatePullRequest(pr); err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return pr, nil
}

func (s *Service) GetUserReviews(userID string) ([]model.PullRequest, error) {
	_, err := s.repo.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	prs, err := s.repo.GetPullRequestsByReviewer(userID)
	if err != nil {
		return nil, err
	}

	return prs, nil
}

func (s *Service) ReassignReviewer(prID, oldUserID string) (*model.PullRequest, string, error) {
	pr, err := s.repo.GetPullRequestByID(prID)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return nil, "", ErrNotFound
		}
		return nil, "", err
	}

	if pr.Status == model.PRStatusMerged {
		return nil, "", ErrPRMerged
	}

	idx := -1
	for i, r := range pr.AssignedReviewers {
		if r == oldUserID {
			idx = i
			break
		}
	}
	if idx == -1 {
		return nil, "", ErrNotAssigned
	}

	oldUser, err := s.repo.GetUserByID(oldUserID)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return nil, "", ErrNotFound
		}
		return nil, "", err
	}

	team, err := s.repo.GetTeam(oldUser.TeamName)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return nil, "", ErrNotFound
		}
		return nil, "", err
	}

	assignedSet := make(map[string]struct{}, len(pr.AssignedReviewers))
	for _, r := range pr.AssignedReviewers {
		assignedSet[r] = struct{}{}
	}

	candidates := make([]model.User, 0, len(team.Members))
	for _, m := range team.Members {
		if !m.IsActive {
			continue
		}
		if m.ID == oldUserID {
			continue
		}
		if m.ID == pr.AuthorID {
			continue
		}
		if _, alreadyAssigned := assignedSet[m.ID]; alreadyAssigned {
			continue
		}
		candidates = append(candidates, m)
	}

	if len(candidates) == 0 {
		return nil, "", ErrNoCandidate
	}

	newIdx := rand.Intn(len(candidates))
	newReviewer := candidates[newIdx]

	pr.AssignedReviewers[idx] = newReviewer.ID

	if err := s.repo.UpdatePullRequest(pr); err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return nil, "", ErrNotFound
		}
		return nil, "", err
	}

	return pr, newReviewer.ID, nil
}
