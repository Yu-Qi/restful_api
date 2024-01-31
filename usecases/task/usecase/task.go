package usecase

import (
	"context"

	"github.com/Yu-Qi/restful_api/domain"
	"github.com/Yu-Qi/restful_api/pkg/code"
)

// GetTasks get all tasks
func GetTasks(ctx context.Context) ([]*domain.Task, *code.CustomError) {
	tasks, customErr := taskRepo.GetTasks(ctx)
	if customErr != nil {
		return nil, customErr
	}

	return tasks, nil
}
