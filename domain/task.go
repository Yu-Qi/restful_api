//go:generate go-enum
package domain

import (
	"context"

	"github.com/Yu-Qi/restful_api/pkg/code"
)

// TaskStatus description:0 represents an incomplete task, while 1 represents a completed task
// ENUM(incomplete,completed)
type TaskStatus int

// Task represents a task entity
type Task struct {
	ID     int        `json:"-"`
	Name   string     `json:"name"`
	Status TaskStatus `json:"status"`
}

type TaskRepository interface {
	GetTasks(context.Context) ([]*Task, *code.CustomError)
	CreateTask(ctx context.Context, task *Task) *code.CustomError
	UpdateTask(ctx context.Context, params *UpdateTaskParams) *code.CustomError
	DeleteTask(ctx context.Context, id int) *code.CustomError
}
type UpdateTaskParams struct {
	ID     int
	Name   *string
	Status *TaskStatus
}
