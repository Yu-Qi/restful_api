package inmemory

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/Yu-Qi/restful_api/domain"
	"github.com/Yu-Qi/restful_api/domain/model"
	"github.com/Yu-Qi/restful_api/pkg/code"
	"github.com/Yu-Qi/restful_api/pkg/lock"
)

var (
	lockWaitSecond = 5
)

type inMemoryTaskRepo struct {
	StorageMap sync.Map // thread-safe map
	CreateLock sync.Mutex
	// WriteRowLock key: task id, value: sync.Mutex
	WriteRowLock *lock.LockMap
}

// NewInMemoryTaskRepo will create an object that represent the task.Repository interface
func NewInMemoryTaskRepo() domain.TaskRepository {
	return &inMemoryTaskRepo{
		StorageMap:   sync.Map{},
		CreateLock:   sync.Mutex{},
		WriteRowLock: lock.NewLockMap(lockWaitSecond),
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
	i.CreateLock.Lock()
	defer i.CreateLock.Unlock()

	id := getLen(&i.StorageMap) + 1
	i.StorageMap.Store(id, &model.Task{
		Id:     id,
		Name:   task.Name,
		Status: task.Status,
	})
	return nil
}

// UpdateTask will update a task
func (i *inMemoryTaskRepo) UpdateTask(ctx context.Context, params *domain.UpdateTaskParams) *code.CustomError {
	if params.ID == 0 {
		return code.NewCustomError(code.ParamIncorrect, http.StatusBadRequest, fmt.Errorf("id is required"))
	}

	err := i.WriteRowLock.Lock(params.ID)
	if err != nil {
		return code.NewCustomError(code.Timeout, http.StatusInternalServerError, err)
	}
	defer i.WriteRowLock.Unlock(params.ID)

	modelTask, ok := i.StorageMap.Load(params.ID)
	if !ok {
		return code.NewCustomError(code.NotFound, http.StatusNotFound, fmt.Errorf("task not found"))
	}

	if params.Name != nil {
		modelTask.(*model.Task).Name = *params.Name
	}
	if params.Status != nil {
		modelTask.(*model.Task).Status = *params.Status
	}

	i.StorageMap.Store(params.ID, modelTask)
	return nil
}

// DeleteTask will delete a task
func (i *inMemoryTaskRepo) DeleteTask(ctx context.Context, id int) *code.CustomError {
	err := i.WriteRowLock.Lock(id)
	if err != nil {
		return code.NewCustomError(code.Timeout, http.StatusInternalServerError, err)
	}
	defer i.WriteRowLock.Unlock(id)

	_, existing := i.StorageMap.LoadAndDelete(id)
	if !existing {
		return code.NewCustomError(code.NotFound, http.StatusNotFound, fmt.Errorf("task not found"))
	}

	return nil
}
