package service

import (
	"context"
	"testing"

	"github.com/blazee5/imageChecker/internal/domain"
	"github.com/blazee5/imageChecker/internal/service/mocks"
	"github.com/blazee5/imageChecker/lib/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_CheckImage(t *testing.T) {
	t.Parallel()

	type fields struct {
		repository *mocks.Repository
	}

	tests := []struct {
		name     string
		input    domain.CheckImageRequest
		response bool
		mockFunc func(f *fields)
		wantErr  bool
	}{
		{
			name: "success get from repository - image exists",
			input: domain.CheckImageRequest{
				Image:     "nginx",
				IsPrivate: false,
				Username:  "",
				Password:  "",
			},
			mockFunc: func(f *fields) {
				f.repository.On("GetExists", mock.Anything, "index.docker.io", "library/nginx", "latest", "", "").
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
				f.repository.On("GetExists", mock.Anything, "index.docker.io", "bitnami/nginx", "latest", "", "").
					Return(true, nil)
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
				f.repository.On("GetExists", mock.Anything, "index.docker.io", "library/nonexistent", "latest", "", "").
					Return(false, nil)
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
				f.repository.On("GetExists", mock.Anything, "index.docker.io", "test/image", "latest", "", "").
					Return(false, assert.AnError)
			},
			response: false,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fields{repository: mocks.NewRepository(t)}
			tt.mockFunc(&f)

			s := &Service{
				log:  logger.NewLogger(),
				repo: f.repository,
			}

			result, err := s.CheckImage(context.TODO(), tt.input)

			assert.Equal(t, tt.wantErr, err != nil)
			assert.EqualValues(t, tt.response, result)

			f.repository.AssertExpectations(t)
		})
	}
}
