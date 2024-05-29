package service

import (
	"context"
	"github.com/blazee5/imageChecker/internal/domain"
	"github.com/blazee5/imageChecker/lib/docker"
	"log/slog"
)

type Repository interface {
	GetByImage(ctx context.Context, image string) (bool, error)
	SetImage(ctx context.Context, image string, exists bool) error
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

	cacheExists, err := s.repo.GetByImage(ctx, repository)

	if err == nil {
		return cacheExists, nil
	}

	exists, err := docker.CheckImage(ctx, registry, repository, tag, input.Username, input.Password)

	if err != nil {
		return false, err
	}

	err = s.repo.SetImage(ctx, repository, exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}
