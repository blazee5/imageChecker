package service

import (
	"context"
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
	return s.repo.JobRepository.CreateJob(ctx, input)
}
