package service

import (
	"context"
	"testing"

	"github.com/blazee5/imageChecker/internal/domain"
	"github.com/blazee5/imageChecker/internal/repository"
	"github.com/blazee5/imageChecker/internal/repository/mocks"
	"github.com/blazee5/imageChecker/lib/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_CheckImage(t *testing.T) {
	t.Parallel()

	type fields struct {
		cacheRepo  *mocks.CacheRepository
		dockerRepo *mocks.DockerRepository
	}

	tests := []struct {
		name     string
		input    domain.CheckImageRequest
		response bool
		mockFunc func(f *fields)
		wantErr  bool
	}{
		{
			name: "success get from cache - image exists",
			input: domain.CheckImageRequest{
				Image:     "nginx",
				IsPrivate: false,
				Username:  "",
				Password:  "",
			},
			mockFunc: func(f *fields) {
				f.cacheRepo.On("GetByImage", mock.Anything, "library/nginx").
					Return(true, nil)
			},
			response: true,
			wantErr:  false,
		},
		{
			name: "success get from repository - image exists",
			input: domain.CheckImageRequest{
				Image:     "bitnami/nginx",
				IsPrivate: false,
				Username:  "",
				Password:  "",
			},
			mockFunc: func(f *fields) {
				f.cacheRepo.On("GetByImage", mock.Anything, "bitnami/nginx").
					Return(false, assert.AnError)
				f.dockerRepo.On("GetExists", mock.Anything, "index.docker.io", "bitnami/nginx", "latest", "", "").
					Return(true, nil)
				f.cacheRepo.On("SetByImage", mock.Anything, "bitnami/nginx", true).
					Return(nil)
			},
			response: true,
			wantErr:  false,
		},
		{
			name: "success get from repository - image does not exist",
			input: domain.CheckImageRequest{
				Image:     "nonexistent:latest",
				IsPrivate: false,
				Username:  "",
				Password:  "",
			},
			mockFunc: func(f *fields) {
				f.cacheRepo.On("GetByImage", mock.Anything, "library/nonexistent").
					Return(false, assert.AnError)
				f.dockerRepo.On("GetExists", mock.Anything, "index.docker.io", "library/nonexistent", "latest", "", "").
					Return(false, nil)
				f.cacheRepo.On("SetByImage", mock.Anything, "library/nonexistent", false).
					Return(nil)
			},
			response: false,
			wantErr:  false,
		},
		{
			name: "error from repository",
			input: domain.CheckImageRequest{
				Image:     "test/image",
				IsPrivate: false,
				Username:  "",
				Password:  "",
			},
			mockFunc: func(f *fields) {
				f.cacheRepo.On("GetByImage", mock.Anything, "test/image").
					Return(false, assert.AnError)
				f.dockerRepo.On("GetExists", mock.Anything, "index.docker.io", "test/image", "latest", "", "").
					Return(false, assert.AnError)
			},
			response: false,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fields{
				cacheRepo:  mocks.NewCacheRepository(t),
				dockerRepo: mocks.NewDockerRepository(t),
			}
			tt.mockFunc(&f)

			log := logger.NewLogger()
			svc := &ImageService{
				log: log,
				repo: &repository.Repository{
					CacheRepository:  f.cacheRepo,
					DockerRepository: f.dockerRepo,
				},
			}

			result, err := svc.CheckImage(context.TODO(), tt.input)

			assert.Equal(t, tt.wantErr, err != nil)
			assert.EqualValues(t, tt.response, result)

			f.cacheRepo.AssertExpectations(t)
			f.dockerRepo.AssertExpectations(t)
		})
	}
}
