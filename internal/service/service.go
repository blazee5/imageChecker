package service

import (
	"context"
	"github.com/blazee5/imageChecker/internal/domain"
	"github.com/blazee5/imageChecker/lib/docker"
	"log/slog"
)

//go:generate go run github.com/vektra/mockery/v2@v2.43.2 --name=Repository
type Repository interface {
	GetExists(ctx context.Context, registry, repository, tag, username, password string) (bool, error)
}

type Service struct {
	log  *slog.Logger
	repo Repository
}

func NewService(log *slog.Logger, repo Repository) *Service {
	return &Service{log: log, repo: repo}
}

func (s *Service) CheckImage(ctx context.Context, input domain.CheckImageRequest) (bool, error) {
	registry, repository, tag := docker.ParseDockerImage(input.Image)

	exists, err := s.repo.GetExists(ctx, registry, repository, tag, input.Username, input.Password)

	if err == nil {
		return exists, nil
	}

	return exists, nil
}
