package middleware

import (
	em "emperror.dev/emperror"
	"emperror.dev/errors"
	ee "emperror.dev/errors"

	"github.com/Yu-Qi/restful_api/pkg/api/response"
	customlog "github.com/Yu-Qi/restful_api/pkg/custom_log"
	"github.com/gin-gonic/gin"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

// HandlePanic that recovers from any panics and handles the error
func HandlePanic(c *gin.Context) {
	handleError := em.ErrorHandlerFunc(func(err error) {
		customlog.Error(err.Error())
		errTracer, ok := err.(stackTracer) // ok is false if errors doesn't implement stackTracer
		if ok {
			customlog.ErrorWithData("stack trace", errTracer.StackTrace())
		}
		response.AbortByAny(c, ee.WithStackDepth(err, 10))
	})
	defer em.HandleRecover(handleError)
	c.Next()
}
