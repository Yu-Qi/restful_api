package http

import (
	"github.com/Yu-Qi/restful_api/pkg/api/response"
	"github.com/Yu-Qi/restful_api/usecases/task/usecase"
	"github.com/gin-gonic/gin"
)

// type getLeaderboardsParams struct {
// 	Metrics       domain.MetricsType       `form:"metrics" binding:"required"`
// 	Group         domain.GroupType         `form:"type" binding:"required"`
// 	TimeDimension domain.TimeDimensionType `form:"time_dimension" binding:"required"`
// }

// TaskHandler represent the http handler for tasks
type TaskHandler struct{}

// NewTaskHandler will initialize the tasks/ resources endpoint
func NewTaskHandler(r *gin.RouterGroup) {
	handler := &TaskHandler{}
	r.GET("/tasks", handler.GetTasks)
	// r.POST("/tasks", handler.CreateTask)
	// r.PUT("/tasks/{id}", handler.UpdateTask)
	// r.DELETE("/tasks/{id}", handler.DeleteTask)

}

// TaskHandler get all tasks
func (t *TaskHandler) GetTasks(ctx *gin.Context) {
	// TODO: filter, 分頁
	// params := getLeaderboardsParams{}
	// udonErr := util.ToGinContextExt(ctx).BindQueryOrPanic(params)
	// if udonErr != nil {
	// 	response.UdonError(ctx, udonErr)
	// 	return
	// }

	tasks, customErr := usecase.GetTasks(ctx)
	if customErr != nil {
		response.CustomError(ctx, customErr)
		return
	}

	response.OK(ctx, tasks)
}
