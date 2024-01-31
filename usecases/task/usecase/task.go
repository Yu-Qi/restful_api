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

// CreateTask create a task
func CreateTask(ctx context.Context, task *domain.Task) *code.CustomError {
	customErr := taskRepo.CreateTask(ctx, task)
	if customErr != nil {
		return customErr
	}

	return nil
}

// UpdateTask update a task
func UpdateTask(ctx context.Context, params *domain.UpdateTaskParams) *code.CustomError {
	customErr := taskRepo.UpdateTask(ctx, params)
	if customErr != nil {
		return customErr
	}

	return nil
}

// DeleteTask delete a task
func DeleteTask(ctx context.Context, id int) *code.CustomError {
	customErr := taskRepo.DeleteTask(ctx, id)
	if customErr != nil {
		return customErr
	}

	return nil
}
