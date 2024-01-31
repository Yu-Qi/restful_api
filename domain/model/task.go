package model

import "github.com/Yu-Qi/restful_api/domain"

// Task represents a task entity for repository
type Task struct {
	Id     int
	Name   string
	Status domain.TaskStatus
}
