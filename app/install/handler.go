package install

import (
	"nest-api/pkg/logger"
	"nest-api/pkg/response"
	"nest-api/pkg/validator"

	"github.com/gin-gonic/gin"
)

type Handler struct{}

func (*Handler) Status(c *gin.Context) {
	data := (Service{}).Status()
	response.Success(c, data)
}

func (*Handler) TestDatabase(c *gin.Context) {
	var params TestDatabaseRequest
	if err := validator.Bind(c, &params); err != nil {
		response.Fail(c, err)
		return
	}

	data, err := (Service{}).TestDatabase(params)
	if err != nil {
		logger.Error("install test database failed", err)
		response.Fail(c, err)
		return
	}

	response.Success(c, data)
}

func (*Handler) Install(c *gin.Context) {
	var params InstallRequest
	if err := validator.Bind(c, &params); err != nil {
		response.Fail(c, err)
		return
	}

	data, err := (Service{}).Install(c.Request.Context(), params)
	if err != nil {
		logger.Error("install failed", err)
		response.Fail(c, err)
		return
	}

	response.Success(c, data)
}
