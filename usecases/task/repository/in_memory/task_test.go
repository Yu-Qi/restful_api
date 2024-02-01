package inmemory

import (
	"context"
	"sync"
	"testing"

	"github.com/samber/lo"

	"github.com/Yu-Qi/restful_api/domain"
	"github.com/Yu-Qi/restful_api/domain/seed"
	"github.com/Yu-Qi/restful_api/pkg/code"
	"github.com/Yu-Qi/restful_api/pkg/util"
	"github.com/stretchr/testify/suite"
)

type getTaskSuite struct {
	suite.Suite
	taskRepo domain.TaskRepository
}

func TestTaskSuite(t *testing.T) {
	suite.Run(t, new(getTaskSuite))
	suite.Run(t, new(createTaskSuite))
	suite.Run(t, new(updateTaskSuite))
	suite.Run(t, new(deleteTaskSuite))
}

func (s *getTaskSuite) SetupTest() {
	s.taskRepo = NewInMemoryTaskRepo()

	// setup data
	for _, task := range seed.Tasks() {
		s.taskRepo.CreateTask(context.Background(), task)
	}
}

func (s *getTaskSuite) TestGetTasks() {
	tasks, customErr := s.taskRepo.GetTasks(context.Background())
	s.Nil(customErr)
	s.Equal(5, len(tasks))
}

type createTaskSuite struct {
	suite.Suite
	taskRepo domain.TaskRepository
}

func (s *createTaskSuite) SetupTest() {
	s.taskRepo = NewInMemoryTaskRepo()
}

func (s *createTaskSuite) TestCreateTask() {
	ctx := context.Background()
	var workers sync.WaitGroup
	for _, task := range seed.Tasks() {
		workers.Add(1)
		go func(task *domain.Task) {
			defer workers.Done()
			customErr := s.taskRepo.CreateTask(ctx, task)
			s.Nil(customErr)
		}(task)
	}

	workers.Wait()
	actualTasks, customErr := s.taskRepo.GetTasks(ctx)
	s.Nil(customErr)
	s.Equal(len(seed.Tasks()), len(actualTasks))
}

type updateTaskSuite struct {
	suite.Suite
	taskRepo domain.TaskRepository
}

func (s *updateTaskSuite) SetupTest() {
	s.taskRepo = NewInMemoryTaskRepo()

	// setup data
	for _, task := range seed.Tasks() {
		s.taskRepo.CreateTask(context.Background(), task)
	}
}

func (s *updateTaskSuite) TestAtomicUpdateTask() {
	ctx := context.Background()
	var workers sync.WaitGroup

	newTasks := make([]*domain.Task, 0, len(seed.Tasks()))
	for _, task := range seed.Tasks() {
		newTasks = append(newTasks, &domain.Task{
			Name:   task.Name + "_new",
			Status: domain.TaskStatus(1 - int(task.Status)), // covert 0 -> 1, 1 -> 0
		})
	}

	for i, task := range seed.Tasks() {
		workers.Add(1)
		go func(task *domain.Task, i int) {
			defer workers.Done()
			customErr := s.taskRepo.UpdateTask(ctx, &domain.UpdateTaskParams{
				ID:     i + 1,
				Name:   &newTasks[i].Name,
				Status: &newTasks[i].Status,
			})
			s.Nil(customErr)

			// check
			actualTasks, customErr := s.taskRepo.GetTasks(ctx)
			s.Nil(customErr)
			for _, actualTask := range actualTasks {
				if actualTask.ID == i+1 {
					s.Equal(newTasks[i].Name, actualTask.Name)
					s.Equal(newTasks[i].Status, actualTask.Status)
				}
			}
		}(task, i)
	}

	workers.Wait()
	actualTasks, customErr := s.taskRepo.GetTasks(ctx)
	s.Nil(customErr)
	s.Equal(len(seed.Tasks()), len(actualTasks))
}

func (s *updateTaskSuite) TestUpdateSomeFields() {
	ctx := context.Background()
	var workers sync.WaitGroup

	newTasks := make([]*domain.Task, 0, len(seed.Tasks()))
	for _, task := range seed.Tasks() {
		newTasks = append(newTasks, &domain.Task{
			Name: task.Name + "_new",
		})
	}

	for i, task := range seed.Tasks() {
		workers.Add(1)
		go func(task *domain.Task, i int) {
			defer workers.Done()
			customErr := s.taskRepo.UpdateTask(ctx, &domain.UpdateTaskParams{
				ID:   i + 1,
				Name: &newTasks[i].Name,
			})
			s.Nil(customErr)
		}(task, i)
	}

	workers.Wait()
	actualTasks, customErr := s.taskRepo.GetTasks(ctx)
	s.Nil(customErr)
	s.Equal(len(seed.Tasks()), len(actualTasks))

	for _, actualTask := range actualTasks {
		s.Equal(newTasks[actualTask.ID-1].Name, actualTask.Name)         // name should not be changed
		s.Equal(seed.Tasks()[actualTask.ID-1].Status, actualTask.Status) // status should not be changed
	}
}

func (s *updateTaskSuite) TestUpdateSameRow() {
	ctx := context.Background()
	var workers sync.WaitGroup

	taskID := 1
	for i := 0; i < 10; i++ {
		workers.Add(1)
		go func() {
			defer workers.Done()
			customErr := s.taskRepo.UpdateTask(ctx, &domain.UpdateTaskParams{
				ID:   taskID,
				Name: util.Ptr(util.RandString(10)),
			})
			s.Nil(customErr)
		}()
	}

	workers.Wait()
	actualTasks, customErr := s.taskRepo.GetTasks(ctx)
	s.Nil(customErr)
	s.Equal(len(seed.Tasks()), len(actualTasks))

	for _, actualTask := range actualTasks {
		if actualTask.ID == taskID {
			s.NotEqual(seed.Tasks()[taskID-1].Name, actualTask.Name)
		}

	}
}

type deleteTaskSuite struct {
	suite.Suite
	taskRepo domain.TaskRepository
}

func (s *deleteTaskSuite) SetupTest() {
	s.taskRepo = NewInMemoryTaskRepo()

	// setup data
	for _, task := range seed.Tasks() {
		s.taskRepo.CreateTask(context.Background(), task)
	}
}

func (s *deleteTaskSuite) TestDeleteTask() {
	ctx := context.Background()
	var workers sync.WaitGroup

	deleteTaskIDs := []int{1, 2}
	for _, taskID := range deleteTaskIDs {
		workers.Add(1)
		go func(taskID int) {
			defer workers.Done()
			customErr := s.taskRepo.DeleteTask(ctx, taskID)
			s.Nil(customErr)
		}(taskID)
	}

	workers.Wait()
	actualTasks, customErr := s.taskRepo.GetTasks(ctx)
	s.Nil(customErr)
	s.Equal(len(seed.Tasks())-len(deleteTaskIDs), len(actualTasks))

	for _, actualTask := range actualTasks {
		s.False(lo.Contains(deleteTaskIDs, actualTask.ID))
	}
}

func (s *deleteTaskSuite) TestSimultaneousDeleteAndUpdate() {
	ctx := context.Background()
	var workers sync.WaitGroup

	for _, task := range seed.Tasks() {
		workers.Add(1)
		go func(task *domain.Task) {
			defer workers.Done()
			customErr := s.taskRepo.DeleteTask(ctx, task.ID)
			s.Nil(customErr)

			customErr = s.taskRepo.UpdateTask(ctx, &domain.UpdateTaskParams{
				ID:   task.ID,
				Name: util.Ptr(util.RandString(10)),
			})
			s.True(customErr == nil || customErr.Code == code.NotFound)
		}(task)
	}

	workers.Wait()
}
