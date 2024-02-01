package http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Yu-Qi/restful_api/domain"
	"github.com/Yu-Qi/restful_api/pkg/api/response"
	"github.com/Yu-Qi/restful_api/pkg/code"
	"github.com/Yu-Qi/restful_api/pkg/util"
	"github.com/Yu-Qi/restful_api/usecases/task/usecase"
	"github.com/gin-gonic/gin"
)

// TaskHandler represent the http handler for tasks
type TaskHandler struct{}

// NewTaskHandler will initialize the tasks/ resources endpoint
func NewTaskHandler(r *gin.RouterGroup) {
	v1 := r.Group("/v1")

	handler := &TaskHandler{}
	v1.GET("/tasks", handler.GetTasks)
	v1.POST("/tasks", handler.CreateTask)
	v1.PUT("/tasks/:id", handler.UpdateTask)
	v1.DELETE("/tasks/:id", handler.DeleteTask)
}

// TaskHandler get all tasks
func (t *TaskHandler) GetTasks(ctx *gin.Context) {
	tasks, customErr := usecase.GetTasks(ctx)
	if customErr != nil {
		response.CustomError(ctx, customErr)
		return
	}

	response.OK(ctx, tasks)
}

type createTaskParams struct {
	Name   string            `json:"name" binding:"required"`
	Status domain.TaskStatus `json:"status" binding:"required"`
}

// CreateTask create a task
func (t *TaskHandler) CreateTask(ctx *gin.Context) {
	task := createTaskParams{}
	customErr := util.ToGinContextExt(ctx).BindJson(&task)
	if customErr != nil {
		response.CustomError(ctx, customErr)
		return
	}

	if !task.Status.IsValid() {
		customErr = code.NewCustomError(code.ParamIncorrect, http.StatusBadRequest, fmt.Errorf("status is invalid"))
		response.CustomError(ctx, customErr)
		return
	}

	customErr = usecase.CreateTask(ctx, &domain.Task{
		Name:   task.Name,
		Status: task.Status,
	})
	if customErr != nil {
		response.CustomError(ctx, customErr)
		return
	}
	response.OK(ctx, task)
}

type updateTaskParams struct {
	Name   *string            `json:"name"`
	Status *domain.TaskStatus `json:"status"`
}

// UpdateTask update a task
func (t *TaskHandler) UpdateTask(ctx *gin.Context) {
	task := updateTaskParams{}
	customErr := util.ToGinContextExt(ctx).BindJson(&task)
	if customErr != nil {
		response.CustomError(ctx, customErr)
		return
	}
	taskID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		customErr = code.NewCustomError(code.ParamIncorrect, http.StatusBadRequest, err)
		response.CustomError(ctx, customErr)
		return
	}

	if !task.Status.IsValid() {
		customErr = code.NewCustomError(code.ParamIncorrect, http.StatusBadRequest, fmt.Errorf("status is invalid"))
		response.CustomError(ctx, customErr)
		return
	}

	customErr = usecase.UpdateTask(ctx, &domain.UpdateTaskParams{
		ID:     taskID,
		Name:   task.Name,
		Status: task.Status,
	})
	if customErr != nil {
		response.CustomError(ctx, customErr)
		return
	}
	response.OK(ctx, task)
}

// DeleteTask delete a task
func (t *TaskHandler) DeleteTask(ctx *gin.Context) {
	taskID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		customErr := code.NewCustomError(code.ParamIncorrect, http.StatusBadRequest, err)
		response.CustomError(ctx, customErr)
		return
	}

	customErr := usecase.DeleteTask(ctx, taskID)
	if customErr != nil {
		response.CustomError(ctx, customErr)
		return
	}
	response.OK(ctx, nil)
}
