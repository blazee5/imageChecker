package docker

import (
	"context"
	"github.com/blazee5/imageChecker/lib/api/docker"
	"log/slog"
	"time"
)

type Repository struct {
	log     *slog.Logger
	timeout time.Duration
}

func NewRepository(log *slog.Logger, timeout time.Duration) *Repository {
	return &Repository{log: log, timeout: timeout}
}

func (repo *Repository) GetExists(ctx context.Context, registry, repository, tag, username, password string) (bool, error) {
	exists, err := docker.CheckImage(ctx, registry, repository, tag, username, password, repo.timeout)

	if err != nil {
		repo.log.Error("error while get image exists in docker api", "error", err)

		return false, err
	}

	return exists, nil
}
