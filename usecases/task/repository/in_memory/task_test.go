package inmemory

import (
	"context"
	"testing"

	"github.com/Yu-Qi/restful_api/domain"
	"github.com/Yu-Qi/restful_api/domain/seed"
	"github.com/stretchr/testify/suite"
)

type getTaskSuite struct {
	suite.Suite
	taskRepo domain.TaskRepository
}

func TestGetTaskSuite(t *testing.T) {
	suite.Run(t, new(getTaskSuite))
}

func (s *getTaskSuite) SetupSuite() {
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
