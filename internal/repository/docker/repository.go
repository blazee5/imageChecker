package docker

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/blazee5/imageChecker/lib/api/docker"
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

		return false, fmt.Errorf("error while get image exists in docker api: %w", err)
	}

	return exists, nil
}
