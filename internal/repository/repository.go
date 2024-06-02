package repository

import (
	"context"
	"github.com/blazee5/imageChecker/internal/config"
	"github.com/blazee5/imageChecker/internal/domain"
	"github.com/blazee5/imageChecker/internal/repository/docker"
	"github.com/blazee5/imageChecker/internal/repository/nomad"
	redisRepo "github.com/blazee5/imageChecker/internal/repository/redis"
	"github.com/hashicorp/nomad/api"
	"github.com/redis/go-redis/v9"
	"log/slog"
)

type Repository struct {
	CacheRepository
	DockerRepository
	JobRepository
}

//go:generate go run github.com/vektra/mockery/v2@v2.43.2 --name=CacheRepository
type CacheRepository interface {
	GetByImage(ctx context.Context, image string) (bool, error)
	SetByImage(ctx context.Context, image string, exists bool) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.43.2 --name=DockerRepository
type DockerRepository interface {
	GetExists(ctx context.Context, registry, repository, tag, username, password string) (bool, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.43.2 --name=JobRepository
type JobRepository interface {
	CreateJob(ctx context.Context, input domain.CreateJobRequest) error
}

func NewRepository(log *slog.Logger, cfg *config.Config, rdb *redis.Client, client *api.Client) *Repository {
	return &Repository{
		CacheRepository:  redisRepo.NewImageRepository(log, rdb),
		DockerRepository: docker.NewRepository(log, cfg.Timeout),
		JobRepository:    nomad.NewJobRepository(log, client),
	}
}
