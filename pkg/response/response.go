package response

import (
	"encoding/json"
	"net/http"

	bizerr "nest-api/pkg/errors"
	"nest-api/pkg/paginator"
	requestvalidator "nest-api/pkg/validator"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Resp struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    interface{}     `json:"data"`
	Meta    *paginator.Meta `json:"meta,omitempty"`
}

type Meta struct {
	Total int `json:"total"`
	Page  int `json:"page"`
	Size  int `json:"size"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Resp{
		Code:    bizerr.CodeSuccess,
		Message: "ok",
		Data:    data,
	})
}

func Paginate(c *gin.Context, data interface{}, meta *paginator.Meta) {
	c.JSON(http.StatusOK, Resp{
		Code:    bizerr.CodeSuccess,
		Message: "ok",
		Data:    data,
		Meta:    meta,
	})
}

func Fail(c *gin.Context, err error) {
	if err == nil {
		Success(c, nil)
		return
	}

	// business error
	if e, ok := bizerr.IsBizError(err); ok {
		c.JSON(http.StatusOK, Resp{
			Code:    e.Code,
			Message: e.Message,
		})
		return
	}

	// param type error
	if isRequestError(err) {
		c.JSON(http.StatusBadRequest, Resp{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	// system error
	c.JSON(http.StatusInternalServerError, Resp{
		Code:    http.StatusInternalServerError,
		Message: "server error",
	})
}

// 401
func Unauthorized(c *gin.Context, msg ...string) {
	message := "unauthorized"

	if len(msg) > 0 {
		message = msg[0]
	}

	c.JSON(http.StatusUnauthorized, Resp{
		Code:    http.StatusUnauthorized,
		Message: message,
	})
}

// 403
func Forbidden(c *gin.Context, msg ...string) {
	message := "forbidden"

	if len(msg) > 0 {
		message = msg[0]
	}

	c.JSON(http.StatusForbidden, Resp{
		Code:    http.StatusForbidden,
		Message: message,
	})
}

// 404
func NotFound(c *gin.Context, msg ...string) {
	message := "not found"

	if len(msg) > 0 {
		message = msg[0]
	}

	c.JSON(http.StatusNotFound, Resp{
		Code:    http.StatusNotFound,
		Message: message,
	})
}

// 405
func MethodNotAllowed(c *gin.Context, msg ...string) {
	message := "method not allowed"

	if len(msg) > 0 {
		message = msg[0]
	}

	c.JSON(http.StatusMethodNotAllowed, Resp{
		Code:    http.StatusMethodNotAllowed,
		Message: message,
	})
}

func JSON(c *gin.Context, data interface{}, err error) {
	if err != nil {
		Fail(c, err)
		return
	}

	Success(c, data)
}

func isRequestError(err error) bool {
	switch err.(type) {
	case *json.UnmarshalTypeError,
		*json.SyntaxError,
		validator.ValidationErrors,
		*validator.InvalidValidationError:
		return true
	default:
		return requestvalidator.IsRequestError(err)
	}
}
