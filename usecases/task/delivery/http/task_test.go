package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Yu-Qi/restful_api/domain/seed"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"

	_taskRepo "github.com/Yu-Qi/restful_api/usecases/task/repository/in_memory"
	_taskUsecase "github.com/Yu-Qi/restful_api/usecases/task/usecase"
)

type task struct {
	Name   string `json:"name"`
	Status int    `json:"status"`
}

// Get /v1/tasks
func TestGetTaskSuite(t *testing.T) {
	suite.Run(t, new(getTaskSuite))
}

type getTaskSuite struct {
	suite.Suite
	Router *gin.Engine
	Url    string
	Ctx    context.Context
}

func (s *getTaskSuite) SetupSuite() {
	s.Router = gin.Default()
	NewTaskHandler(s.Router.Group(""))

	s.Url = "/v1/tasks"
}

func (s *getTaskSuite) SetupTest() {
	taskRepo := _taskRepo.NewInMemoryTaskRepo()
	_taskUsecase.Init(_taskUsecase.InitParam{
		TaskRepo: taskRepo,
	})

	s.Ctx = context.Background()

	for _, task := range seed.Tasks() {
		customErr := _taskUsecase.CreateTask(s.Ctx, task)
		s.Nil(customErr)
	}
}

func (s *getTaskSuite) TestSuccess() {
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", s.Url, nil)
	s.NoError(err)
	s.Router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)
	var response struct {
		Code int    `json:"code"`
		Data []task `json:"data"`
	}

	err = json.Unmarshal(w.Body.Bytes(), &response)
	s.Nil(err)
	s.Equal(0, response.Code)
	s.Equal(len(seed.Tasks()), len(response.Data))

	for _, actualTask := range response.Data {
		for _, expectedTask := range seed.Tasks() {
			if actualTask.Name == expectedTask.Name {
				s.Equal(int(expectedTask.Status), actualTask.Status)
			}
		}
	}
}

// Post /v1/tasks
func TestCreateTaskSuite(t *testing.T) {
	suite.Run(t, new(createTaskSuite))
}

type createTaskSuite struct {
	suite.Suite
	Router *gin.Engine
	Url    string
	Ctx    context.Context
}

func (s *createTaskSuite) SetupSuite() {
	s.Router = gin.Default()
	NewTaskHandler(s.Router.Group(""))

	s.Url = "/v1/tasks"
}

func (s *createTaskSuite) SetupTest() {
	taskRepo := _taskRepo.NewInMemoryTaskRepo()
	_taskUsecase.Init(_taskUsecase.InitParam{
		TaskRepo: taskRepo,
	})

	s.Ctx = context.Background()
}

func (s *createTaskSuite) TestSuccess() {
	body := map[string]interface{}{
		"name":   "test",
		"status": 1,
	}

	w := httptest.NewRecorder()
	jsonStr, err := json.Marshal(body)
	s.NoError(err)
	req, err := http.NewRequest("POST", s.Url, bytes.NewBuffer(jsonStr))
	s.NoError(err)
	s.Router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)
	var response struct {
		Code int `json:"code"`
	}

	err = json.Unmarshal(w.Body.Bytes(), &response)
	s.Nil(err)
	s.Equal(0, response.Code)
	actualTask, customErr := _taskUsecase.GetTasks(s.Ctx)
	s.Nil(customErr)
	s.Equal(1, len(actualTask))
}

func (s *createTaskSuite) TestStatusInvalid() {
	body := map[string]interface{}{
		"name":   "test",
		"status": 2,
	}

	w := httptest.NewRecorder()
	jsonStr, err := json.Marshal(body)
	s.NoError(err)
	req, err := http.NewRequest("POST", s.Url, bytes.NewBuffer(jsonStr))
	s.NoError(err)
	s.Router.ServeHTTP(w, req)
	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *createTaskSuite) TestParamIncorrect() {
	body := map[string]interface{}{
		"name": "test",
	}

	w := httptest.NewRecorder()
	jsonStr, err := json.Marshal(body)
	s.NoError(err)
	req, err := http.NewRequest("POST", s.Url, bytes.NewBuffer(jsonStr))
	s.NoError(err)
	s.Router.ServeHTTP(w, req)
	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *createTaskSuite) TestParamTypeIncorrect() {
	body := map[string]interface{}{
		"name":   "test",
		"status": "1",
	}

	w := httptest.NewRecorder()
	jsonStr, err := json.Marshal(body)
	s.NoError(err)
	req, err := http.NewRequest("POST", s.Url, bytes.NewBuffer(jsonStr))
	s.NoError(err)
	s.Router.ServeHTTP(w, req)
	s.Equal(http.StatusBadRequest, w.Code)
}

// PUT /v1/tasks/:id
func TestUpdateTaskSuite(t *testing.T) {
	suite.Run(t, new(updateTaskSuite))
}

type updateTaskSuite struct {
	suite.Suite
	Router    *gin.Engine
	Url       string
	UrlFormat string
	Ctx       context.Context
}

func (s *updateTaskSuite) SetupSuite() {
	s.Router = gin.Default()
	NewTaskHandler(s.Router.Group(""))

	s.Url = "/v1/tasks/:id"
	s.UrlFormat = "/v1/tasks/%d"
}

func (s *updateTaskSuite) SetupTest() {
	taskRepo := _taskRepo.NewInMemoryTaskRepo()
	_taskUsecase.Init(_taskUsecase.InitParam{
		TaskRepo: taskRepo,
	})

	s.Ctx = context.Background()

	for _, task := range seed.Tasks() {
		customErr := _taskUsecase.CreateTask(s.Ctx, task)
		s.Nil(customErr)
	}
}

func (s *updateTaskSuite) TestSuccessWithCompleted() {
	body := map[string]interface{}{
		"name":   "test",
		"status": 1,
	}

	w := httptest.NewRecorder()
	jsonStr, err := json.Marshal(body)
	s.NoError(err)
	req, err := http.NewRequest("PUT", fmt.Sprintf(s.UrlFormat, 1), bytes.NewBuffer(jsonStr))
	s.NoError(err)
	s.Router.ServeHTTP(w, req)
	s.Equal(http.StatusOK, w.Code)

	var response struct {
		Code int `json:"code"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	s.Nil(err)
	s.Equal(0, response.Code)
	actualTask, customErr := _taskUsecase.GetTasks(s.Ctx)
	s.Nil(customErr)
	for _, task := range actualTask {
		if task.ID == 1 {
			s.Equal("test", task.Name)
			s.Equal(1, int(task.Status))
		}
	}
}

func (s *updateTaskSuite) TestSuccessWithIncompleted() {
	body := map[string]interface{}{
		"name":   "test",
		"status": 0,
	}

	w := httptest.NewRecorder()
	jsonStr, err := json.Marshal(body)
	s.NoError(err)
	req, err := http.NewRequest("PUT", fmt.Sprintf(s.UrlFormat, 1), bytes.NewBuffer(jsonStr))
	s.NoError(err)
	s.Router.ServeHTTP(w, req)
	s.Equal(http.StatusOK, w.Code)

	var response struct {
		Code int `json:"code"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	s.Nil(err)
	s.Equal(0, response.Code)
	actualTask, customErr := _taskUsecase.GetTasks(s.Ctx)
	s.Nil(customErr)
	for _, task := range actualTask {
		if task.ID == 1 {
			s.Equal("test", task.Name)
			s.Equal(0, int(task.Status))
		}
	}
}

func (s *updateTaskSuite) TestNotFound() {
	body := map[string]interface{}{
		"name":   "test",
		"status": 1,
	}

	w := httptest.NewRecorder()
	jsonStr, err := json.Marshal(body)
	s.NoError(err)
	req, err := http.NewRequest("PUT", fmt.Sprintf(s.UrlFormat, 6), bytes.NewBuffer(jsonStr))
	s.NoError(err)
	s.Router.ServeHTTP(w, req)
	s.Equal(http.StatusNotFound, w.Code)
}

func (s *updateTaskSuite) TestParamTypeIncorrect() {
	body := map[string]interface{}{
		"name":   "test",
		"status": "1",
	}

	w := httptest.NewRecorder()
	jsonStr, err := json.Marshal(body)
	s.NoError(err)
	req, err := http.NewRequest("PUT", fmt.Sprintf(s.UrlFormat, 6), bytes.NewBuffer(jsonStr))
	s.NoError(err)
	s.Router.ServeHTTP(w, req)
	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *updateTaskSuite) TestStatusInvalid() {
	body := map[string]interface{}{
		"name":   "test",
		"status": 2,
	}

	w := httptest.NewRecorder()
	jsonStr, err := json.Marshal(body)
	s.NoError(err)
	req, err := http.NewRequest("PUT", fmt.Sprintf(s.UrlFormat, 6), bytes.NewBuffer(jsonStr))
	s.NoError(err)
	s.Router.ServeHTTP(w, req)
	s.Equal(http.StatusBadRequest, w.Code)
}

// DELETE /v1/tasks/:id
func TestDeleteTaskSuite(t *testing.T) {
	suite.Run(t, new(deleteTaskSuite))
}

type deleteTaskSuite struct {
	suite.Suite
	Router    *gin.Engine
	Url       string
	UrlFormat string
	Ctx       context.Context
}

func (s *deleteTaskSuite) SetupSuite() {
	s.Router = gin.Default()
	NewTaskHandler(s.Router.Group(""))

	s.Url = "/v1/tasks/:id"
	s.UrlFormat = "/v1/tasks/%v"
}

func (s *deleteTaskSuite) SetupTest() {
	taskRepo := _taskRepo.NewInMemoryTaskRepo()
	_taskUsecase.Init(_taskUsecase.InitParam{
		TaskRepo: taskRepo,
	})

	s.Ctx = context.Background()

	for _, task := range seed.Tasks() {
		customErr := _taskUsecase.CreateTask(s.Ctx, task)
		s.Nil(customErr)
	}
}

func (s *deleteTaskSuite) TestSuccess() {
	w := httptest.NewRecorder()
	req, err := http.NewRequest("DELETE", fmt.Sprintf(s.UrlFormat, 1), nil)
	s.NoError(err)
	s.Router.ServeHTTP(w, req)
	s.Equal(http.StatusOK, w.Code)

	var response struct {
		Code int `json:"code"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	s.Nil(err)
	s.Equal(0, response.Code)
	actualTask, customErr := _taskUsecase.GetTasks(s.Ctx)
	s.Nil(customErr)
	s.Equal(len(seed.Tasks())-1, len(actualTask))
}

func (s *deleteTaskSuite) TestNotFound() {
	w := httptest.NewRecorder()
	req, err := http.NewRequest("DELETE", fmt.Sprintf(s.UrlFormat, 6), nil)
	s.NoError(err)
	s.Router.ServeHTTP(w, req)
	s.Equal(http.StatusNotFound, w.Code)
}

func (s *deleteTaskSuite) TestPathParamIncorrect() {
	w := httptest.NewRecorder()
	req, err := http.NewRequest("DELETE", fmt.Sprintf(s.UrlFormat, "abc"), nil)
	s.NoError(err)
	s.Router.ServeHTTP(w, req)
	s.Equal(http.StatusBadRequest, w.Code)
}
