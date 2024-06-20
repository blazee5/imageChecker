package service

import (
	"context"
	"github.com/blazee5/imageChecker/internal/domain"
	"github.com/blazee5/imageChecker/internal/repository"
	"log/slog"
)

type Service struct {
	Image
	Job
}

//go:generate go run github.com/vektra/mockery/v2@v2.43.2 --name=Image
type Image interface {
	CheckImage(ctx context.Context, input domain.CheckImageRequest) (bool, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.43.2 --name=Job
type Job interface {
	CreateJob(ctx context.Context, input domain.CreateJobRequest) error
}

func NewService(log *slog.Logger, repo *repository.Repository) *Service {
	return &Service{
		Image: NewImageService(log, repo),
		Job:   NewJobService(log, repo),
	}
}
