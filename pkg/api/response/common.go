package response

import (
	"fmt"
	"net/http"
	"time"

	em "emperror.dev/errors"

	"github.com/Yu-Qi/restful_api/pkg/code"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ErrorResp is the error response struct.
type ErrorResp struct {
	Status    int    `json:"status"`
	Code      int    `json:"code"`
	RequestID string `json:"request_id"`
	Message   string `json:"message"`
	Path      string `json:"path"`
	Timestamp int64  `json:"timestamp"`
}

// AbortByAny .
func AbortByAny(c *gin.Context, anyError interface{}) {
	switch errObject := anyError.(type) {
	case validator.ValidationErrors:
		ErrorWithMsg(c, http.StatusBadRequest, code.ParamIncorrect, errObject.Error())
	case error:
		unwrappedErr := em.Unwrap(errObject)
		if unwrappedErr != nil {
			AbortByAny(c, unwrappedErr)
			return
		}
		ErrorWithMsg(c, http.StatusInternalServerError, code.InternalUnknownError, errObject.Error())
	default:
		ErrorWithMsg(c, http.StatusInternalServerError, code.InternalUnknownError, fmt.Sprintf("%v", anyError))
	}
}

// ErrorWithMsg .
func ErrorWithMsg(ctx *gin.Context, status int, code int, msg string) {
	if msg == "" {
		msg = "ERROR"
	}
	var err = &ErrorResp{
		Status:    status,
		Code:      code,
		RequestID: requestid.Get(ctx),
		Message:   msg,
		Path:      ctx.Request.RequestURI,
		Timestamp: time.Now().Unix(),
	}
	ctx.AbortWithStatusJSON(status, err)
}

// CustomError .
func CustomError(ctx *gin.Context, err *code.CustomError) {
	ErrorWithMsg(ctx, err.HttpStatus, err.Code, err.Error.Error())
}

// OKResp is the ok response struct
type OKResp struct {
	Code int         `json:"code"`
	Data interface{} `json:"data,omitempty"`
}

// OK responses code and data in JSON format.
func OK(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, &OKResp{
		Code: 0,
		Data: data,
	})
}
