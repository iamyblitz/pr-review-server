package service

import (
	"github.com/iamyblitz/pr-reviewer-service/internal/repo"
)

type Service struct {
	repo repo.Repository
}

func NewService(r repo.Repository) *Service {
	return &Service{repo: r}
}
