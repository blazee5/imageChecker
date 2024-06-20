package nomad

import (
	"context"
	"fmt"
	"github.com/blazee5/imageChecker/internal/domain"
	"github.com/hashicorp/nomad/api"
	"log/slog"
)

type JobRepository struct {
	log    *slog.Logger
	client *api.Client
}

func NewJobRepository(log *slog.Logger, client *api.Client) *JobRepository {
	return &JobRepository{log: log, client: client}
}

func (repo *JobRepository) CreateJob(ctx context.Context, input domain.CreateJobRequest) error {
	job := &api.Job{
		ID:   &input.Name,
		Name: &input.Name,
		TaskGroups: []*api.TaskGroup{
			{
				Name: &input.Name,
				Tasks: []*api.Task{
					{
						Name:   "task",
						Driver: "docker",
						Config: map[string]any{
							"image": input.Image,
							"auth": map[string]any{
								"username": input.Username,
								"password": input.Password,
							},
						},
					},
				},
			},
		},
	}

	opts := api.WriteOptions{}

	_, _, err := repo.client.Jobs().Register(job, opts.WithContext(ctx))

	if err != nil {
		repo.log.Error("error while register job in nomad", "error", err)

		return fmt.Errorf("error while register job in nomad: %w", err)
	}

	return nil
}
