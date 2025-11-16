package service

import (
	"errors"

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
	ErrTeamExists = errors.New("team already exists")
	ErrNotFound   = errors.New("not found")
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
