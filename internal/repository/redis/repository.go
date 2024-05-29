package redis

import (
	"context"
	"github.com/blazee5/imageChecker/lib/docker"
	"github.com/redis/go-redis/v9"
	"time"
)

type Repository struct {
	rdb *redis.Client
}

func NewRepository(rdb *redis.Client) *Repository {
	return &Repository{rdb: rdb}
}

func (repo *Repository) GetExists(ctx context.Context, registry, repository, tag, username, password string) (bool, error) {
	exists, err := repo.rdb.Get(ctx, repository).Bool()

	if err == nil {
		return exists, nil
	}

	exists, err = docker.CheckImage(ctx, registry, repository, tag, username, password)

	if err != nil {
		return false, err
	}

	err = repo.rdb.Set(ctx, repository, exists, 24*time.Hour).Err()

	if err != nil {
		return false, err
	}

	return exists, nil
}
