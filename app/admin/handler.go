package admin

import (
	"strconv"

	"nest-api/pkg/logger"
	"nest-api/pkg/response"
	"nest-api/pkg/validator"

	"github.com/gin-gonic/gin"
)

type Handler struct{}

func (*Handler) ListUsers(c *gin.Context) {
	data, err := (Service{}).ListUsers(c.Request.Context())
	if err != nil {
		logger.Error("admin list users failed", err)
		response.Fail(c, err)
		return
	}

	response.Success(c, data)
}

func (*Handler) ListWorkspaces(c *gin.Context) {
	data, err := (Service{}).ListWorkspaces(c.Request.Context())
	if err != nil {
		logger.Error("admin list workspaces failed", err)
		response.Fail(c, err)
		return
	}

	response.Success(c, data)
}

func (*Handler) TransferWorkspace(c *gin.Context) {
	workspaceID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || workspaceID <= 0 {
		response.Fail(c, err)
		return
	}

	var params TransferWorkspaceRequest
	if err := validator.Bind(c, &params); err != nil {
		response.Fail(c, err)
		return
	}

	if err := (Service{}).TransferWorkspace(c.Request.Context(), workspaceID, params.OwnerID); err != nil {
		logger.Error("admin transfer workspace failed", err)
		response.Fail(c, err)
		return
	}

	response.Success(c, nil)
}
