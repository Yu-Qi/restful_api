package main

import (
	"context"
	"os"

	"github.com/Yu-Qi/restful_api/pkg/api/middleware"
	_taskHttpDelivery "github.com/Yu-Qi/restful_api/usecases/task/delivery/http"
	_taskUsecase "github.com/Yu-Qi/restful_api/usecases/task/usecase"

	_taskRepo "github.com/Yu-Qi/restful_api/usecases/task/repository/in_memory"

	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()
	r := gin.New()
	r.Use(
		middleware.HandlePanic,
	)

	registerV1API(r, ctx)
	appPort := os.Getenv("APP_PORT")
	_ = r.Run(":" + appPort)

}

func registerV1API(r *gin.Engine, ctx context.Context) {

	// task
	taskRepo := _taskRepo.NewInMemoryTaskRepo()
	_taskUsecase.Init(_taskUsecase.InitParam{
		TaskRepo: taskRepo,
	})
	_taskHttpDelivery.NewTaskHandler(r.Group(""))
}
