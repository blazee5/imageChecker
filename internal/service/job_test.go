package service

import (
	"context"
	"github.com/blazee5/imageChecker/internal/domain"
	"github.com/blazee5/imageChecker/internal/repository"
	"github.com/blazee5/imageChecker/internal/repository/mocks"
	"github.com/blazee5/imageChecker/lib/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestJobService_CreateJob(t *testing.T) {
	t.Parallel()

	type fields struct {
		jobRepo *mocks.JobRepository
	}

	tests := []struct {
		name     string
		input    domain.CreateJobRequest
		mockFunc func(f *fields)
		wantErr  bool
	}{
		{
			name: "success create nginx",
			input: domain.CreateJobRequest{
				Name:      "test",
				Image:     "nginx",
				IsPrivate: false,
				Username:  "",
				Password:  "",
			},
			mockFunc: func(f *fields) {
				f.jobRepo.On("CreateJob", mock.Anything, mock.AnythingOfType("domain.CreateJobRequest")).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "success create private",
			input: domain.CreateJobRequest{
				Name:      "my-private-service",
				Image:     "my.registry.com/service",
				IsPrivate: true,
				Username:  "user",
				Password:  "password",
			},
			mockFunc: func(f *fields) {
				f.jobRepo.On("CreateJob", mock.Anything, mock.AnythingOfType("domain.CreateJobRequest")).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error from api",
			input: domain.CreateJobRequest{
				Name:      "test",
				Image:     "invalid image",
				IsPrivate: false,
				Username:  "",
				Password:  "",
			},
			mockFunc: func(f *fields) {
				f.jobRepo.On("CreateJob", mock.Anything, mock.AnythingOfType("domain.CreateJobRequest")).
					Return(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "missing job name",
			input: domain.CreateJobRequest{
				Image:     "nginx",
				IsPrivate: false,
				Username:  "",
				Password:  "",
			},
			mockFunc: func(f *fields) {
				f.jobRepo.On("CreateJob", mock.Anything, mock.AnythingOfType("domain.CreateJobRequest")).
					Return(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "missing image name",
			input: domain.CreateJobRequest{
				Name:      "test",
				IsPrivate: false,
				Username:  "",
				Password:  "",
			},
			mockFunc: func(f *fields) {
				f.jobRepo.On("CreateJob", mock.Anything, mock.AnythingOfType("domain.CreateJobRequest")).
					Return(assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fields{
				jobRepo: mocks.NewJobRepository(t),
			}
			tt.mockFunc(&f)

			log := logger.NewLogger()
			svc := &JobService{
				log: log,
				repo: &repository.Repository{
					JobRepository: f.jobRepo,
				},
			}

			err := svc.CreateJob(context.TODO(), tt.input)

			assert.Equal(t, tt.wantErr, err != nil)

			f.jobRepo.AssertExpectations(t)
		})
	}
}
