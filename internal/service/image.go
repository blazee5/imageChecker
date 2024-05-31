package service

import (
	"context"
	"github.com/blazee5/imageChecker/internal/domain"
	"github.com/blazee5/imageChecker/internal/repository"
	"github.com/blazee5/imageChecker/lib/docker"
	"log/slog"
)

type ImageService struct {
	log  *slog.Logger
	repo *repository.Repository
}

func NewImageService(log *slog.Logger, repo *repository.Repository) *ImageService {
	return &ImageService{log: log, repo: repo}
}

func (s *ImageService) CheckImage(ctx context.Context, input domain.CheckImageRequest) (bool, error) {
	registry, repository, tag := docker.ParseDockerImage(input.Image)

	cacheExists, err := s.repo.CacheRepository.GetByImage(ctx, repository)

	if err == nil {
		return cacheExists, nil
	}

	exists, err := s.repo.DockerRepository.GetExists(ctx, registry, repository, tag, input.Username, input.Password)

	if err != nil {
		s.log.Error("error while get image in docker api", "error", err)

		return exists, err
	}

	err = s.repo.CacheRepository.SetByImage(ctx, repository, exists)

	if err != nil {
		s.log.Error("error while set image in cache", "error", err)

		return exists, err
	}

	return exists, nil
}
