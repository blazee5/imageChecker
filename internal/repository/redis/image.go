package redis

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

const dayDuration = 24 * time.Hour

type ImageRepository struct {
	log *slog.Logger
	rdb *redis.Client
}

func NewImageRepository(log *slog.Logger, rdb *redis.Client) *ImageRepository {
	return &ImageRepository{log: log, rdb: rdb}
}

func (repo *ImageRepository) GetByImage(ctx context.Context, image string) (bool, error) {
	exists, err := repo.rdb.Get(ctx, image).Bool()

	if err != nil {
		repo.log.Error("error while get image exists in redis", "error", err)

		return exists, fmt.Errorf("error while get image in redis: %w", err)
	}

	return exists, nil
}

func (repo *ImageRepository) SetByImage(ctx context.Context, image string, exists bool) error {
	err := repo.rdb.Set(ctx, image, exists, dayDuration).Err()

	if err != nil {
		repo.log.Error("error while set image exists in redis", "error", err)

		return fmt.Errorf("error while set image exists in redis: %w", err)
	}

	return nil
}
