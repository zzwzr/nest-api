package auth

import (
	"nest-api/internal/runtime"
	"nest-api/internal/utils"
	"nest-api/pkg/logger"
	"nest-api/pkg/response"
	"nest-api/pkg/validator"

	"github.com/gin-gonic/gin"
)

type Handler struct{}

func (*Handler) Register(c *gin.Context) {
	var params RegisterRequest
	if err := validator.Bind(c, &params); err != nil {
		response.Fail(c, err)
		return
	}

	data, err := (Service{}).Register(c.Request.Context(), params)
	if err != nil {
		logger.Error("register failed", err)
		response.Fail(c, err)
		return
	}

	response.Success(c, data)
}

func (*Handler) Login(c *gin.Context) {
	var params LoginRequest
	if err := validator.Bind(c, &params); err != nil {
		response.Fail(c, err)
		return
	}

	data, err := (Service{}).Login(c.Request.Context(), params)
	if err != nil {
		logger.Error("login failed", err)
		response.Fail(c, err)
		return
	}

	response.Success(c, data)
}

func (*Handler) Me(c *gin.Context) {
	userID := utils.GetUserID(c)
	data, err := (Service{}).Me(c.Request.Context(), userID)
	if err != nil {
		response.Fail(c, err)
		return
	}

	response.Success(c, data)
}

func (*Handler) Site(c *gin.Context) {
	response.Success(c, SiteResponse{
		Installed: runtime.IsInstalled(),
		SiteURL:   runtime.SiteURL(),
	})
}
