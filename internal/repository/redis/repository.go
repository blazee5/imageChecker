package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type Repository struct {
	rdb *redis.Client
}

func NewRepository(rdb *redis.Client) *Repository {
	return &Repository{rdb: rdb}
}

func (repo *Repository) GetByImage(ctx context.Context, image string) (bool, error) {
	exists, err := repo.rdb.Get(ctx, image).Bool()

	if err != nil {
		return exists, err
	}

	return exists, nil
}

func (repo *Repository) SetByImage(ctx context.Context, image string, exists bool) error {
	err := repo.rdb.Set(ctx, image, exists, 24*time.Hour).Err()

	if err != nil {
		return err
	}

	return nil
}
