package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"time"
)

type Repository struct {
	log *slog.Logger
	rdb *redis.Client
}

func NewRepository(log *slog.Logger, rdb *redis.Client) *Repository {
	return &Repository{log: log, rdb: rdb}
}

func (repo *Repository) GetByImage(ctx context.Context, image string) (bool, error) {
	exists, err := repo.rdb.Get(ctx, image).Bool()

	if err != nil {
		repo.log.Error("error while get image exists in redis", "error", err)

		return exists, err
	}

	return exists, nil
}

func (repo *Repository) SetByImage(ctx context.Context, image string, exists bool) error {
	err := repo.rdb.Set(ctx, image, exists, 24*time.Hour).Err()

	if err != nil {
		repo.log.Error("error while set image exists in redis", "error", err)

		return err
	}

	return nil
}
