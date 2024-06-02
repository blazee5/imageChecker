package service

import (
	"context"
	"fmt"
	"github.com/blazee5/imageChecker/internal/domain"
	"github.com/blazee5/imageChecker/internal/repository"
	"log/slog"
)

type JobService struct {
	log  *slog.Logger
	repo *repository.Repository
}

func NewJobService(log *slog.Logger, repo *repository.Repository) *JobService {
	return &JobService{log: log, repo: repo}
}

func (s *JobService) CreateJob(ctx context.Context, input domain.CreateJobRequest) error {
	err := s.repo.JobRepository.CreateJob(ctx, input)

	if err != nil {
		s.log.Error("error while create job in repo", "error", err)

		return fmt.Errorf("error while create job in repo: %w", err)
	}

	return nil
}
