package inmemory

import (
	"context"
	"sync"

	"github.com/Yu-Qi/restful_api/domain"
	"github.com/Yu-Qi/restful_api/domain/model"
	"github.com/Yu-Qi/restful_api/pkg/code"
)

type inMemoryTaskRepo struct {
	StorageMap sync.Map
	Mutex      sync.Mutex
}

// NewInMemoryTaskRepo will create an object that represent the task.Repository interface
func NewInMemoryTaskRepo() domain.TaskRepository {
	return &inMemoryTaskRepo{
		StorageMap: sync.Map{},
		Mutex:      sync.Mutex{},
	}
}

// GetTasks will get all tasks
func (i *inMemoryTaskRepo) GetTasks(ctx context.Context) ([]*domain.Task, *code.CustomError) {
	var tasks []*domain.Task

	i.StorageMap.Range(func(key, value interface{}) bool {
		modelTask, ok := value.(*model.Task)
		if !ok {
			// skip
			return true
		}
		tasks = append(tasks, &domain.Task{
			Name:   modelTask.Name,
			Status: modelTask.Status,
		})
		return true
	})

	return tasks, nil
}

// CreateTask will create a task
func (i *inMemoryTaskRepo) CreateTask(ctx context.Context, task *domain.Task) *code.CustomError {
	i.Mutex.Lock()
	defer i.Mutex.Unlock()

	id := getLen(&i.StorageMap) + 1
	i.StorageMap.Store(id, &model.Task{
		Id:     id,
		Name:   task.Name,
		Status: task.Status,
	})
	return nil
}
