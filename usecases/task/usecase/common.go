package usecase

import (
	"github.com/Yu-Qi/restful_api/domain"
)

var (
	taskRepo domain.TaskRepository
)

// InitParam defines the parameters for initializing the service.
type InitParam struct {
	TaskRepo domain.TaskRepository
}

// Init injects implementations into the service.
func Init(param InitParam) {
	taskRepo = param.TaskRepo
}
