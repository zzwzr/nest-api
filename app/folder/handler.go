package folder

import (
	"nest-api/internal/utils"
	"nest-api/pkg/logger"
	"nest-api/pkg/response"
	"nest-api/pkg/validator"

	"github.com/gin-gonic/gin"
)

type Handler struct{}

func (*Handler) Tree(c *gin.Context) {
	var params ProjectScopeRequest
	if err := validator.Bind(c, &params); err != nil {
		response.Fail(c, err)
		return
	}

	list, err := (Service{}).Tree(c.Request.Context(), utils.GetUserID(c), params)
	if err != nil {
		logger.Error("folder tree failed", err)
		response.Fail(c, err)
		return
	}
	response.Success(c, list)
}

func (*Handler) Create(c *gin.Context) {
	var params CreateRequest
	if err := validator.Bind(c, &params); err != nil {
		response.Fail(c, err)
		return
	}

	id, err := (Service{}).Create(c.Request.Context(), utils.GetUserID(c), params)
	if err != nil {
		logger.Error("folder create failed", err)
		response.Fail(c, err)
		return
	}
	response.Success(c, gin.H{"id": id})
}

func (*Handler) Update(c *gin.Context) {
	var params UpdateRequest
	if err := validator.Bind(c, &params); err != nil {
		response.Fail(c, err)
		return
	}

	if err := (Service{}).Update(c.Request.Context(), utils.GetUserID(c), params); err != nil {
		logger.Error("folder update failed", err)
		response.Fail(c, err)
		return
	}
	response.Success(c, nil)
}

func (*Handler) Delete(c *gin.Context) {
	var params DeleteRequest
	if err := validator.Bind(c, &params); err != nil {
		response.Fail(c, err)
		return
	}

	if err := (Service{}).Delete(c.Request.Context(), utils.GetUserID(c), params); err != nil {
		logger.Error("folder delete failed", err)
		response.Fail(c, err)
		return
	}
	response.Success(c, nil)
}
