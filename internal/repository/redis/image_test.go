package redis_test

import (
	"context"
	redisRepo "github.com/blazee5/imageChecker/internal/repository/redis"
	"github.com/blazee5/imageChecker/lib/logger"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"testing"
)

func setupRedisContainer(ctx context.Context, t *testing.T) (testcontainers.Container, string) {
	req := testcontainers.ContainerRequest{
		Image:        "redis:latest",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}

	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	require.NoError(t, err)

	endpoint, err := redisC.Endpoint(ctx, "")

	require.NoError(t, err)

	return redisC, endpoint
}

func setupImageRepo(ctx context.Context, t *testing.T) (*redisRepo.ImageRepository, testcontainers.Container) {
	container, connURI := setupRedisContainer(ctx, t)

	client := redis.NewClient(&redis.Options{
		Addr: connURI,
	})

	log := logger.NewLogger()

	imageRepo := redisRepo.NewImageRepository(log, client)

	return imageRepo, container
}

func TestImageRepository_GetByImage(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	repo, container := setupImageRepo(ctx, t)

	err := repo.SetByImage(ctx, "nginx", true)

	require.NoError(t, err)

	defer func(container testcontainers.Container, ctx context.Context) {
		err := container.Terminate(ctx)
		if err != nil {
			t.Fatalf("could not terminate redis container: %v", err.Error())
		}
	}(container, ctx)

	tests := []struct {
		name     string
		image    string
		response bool
		wantErr  bool
	}{
		{
			name:     "success",
			image:    "nginx",
			response: true,
			wantErr:  false,
		},
		{
			name:     "image not exists",
			image:    "test",
			response: false,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exists, err := repo.GetByImage(ctx, tt.image)

			assert.Equal(t, tt.wantErr, err != nil)

			assert.Equal(t, tt.response, exists)
		})
	}
}

func TestImageRepository_SetByImage(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	repo, container := setupImageRepo(ctx, t)

	defer func(container testcontainers.Container, ctx context.Context) {
		err := container.Terminate(ctx)
		if err != nil {
			t.Fatalf("could not terminate redis container: %v", err.Error())
		}
	}(container, ctx)

	tests := []struct {
		name    string
		image   string
		exists  bool
		wantErr bool
	}{
		{
			name:    "success",
			image:   "nginx",
			exists:  true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.SetByImage(ctx, tt.image, tt.exists)

			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
