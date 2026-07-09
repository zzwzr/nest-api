package member

import (
	"nest-api/internal/utils"
	"nest-api/pkg/logger"
	"nest-api/pkg/response"
	"nest-api/pkg/validator"

	"github.com/gin-gonic/gin"
)

type Handler struct{}

func (*Handler) List(c *gin.Context) {
	var params ListRequest
	if err := validator.Bind(c, &params); err != nil {
		response.Fail(c, err)
		return
	}

	list, err := (Service{}).List(c.Request.Context(), utils.GetUserID(c), params)
	if err != nil {
		logger.Error("member list failed", err)
		response.Fail(c, err)
		return
	}
	response.Success(c, list)
}

func (*Handler) Invite(c *gin.Context) {
	var params InviteRequest
	if err := validator.Bind(c, &params); err != nil {
		response.Fail(c, err)
		return
	}

	if err := (Service{}).Invite(c.Request.Context(), utils.GetUserID(c), params); err != nil {
		logger.Error("member invite failed", err)
		response.Fail(c, err)
		return
	}
	response.Success(c, nil)
}

func (*Handler) Update(c *gin.Context) {
	var params UpdateRequest
	if err := validator.Bind(c, &params); err != nil {
		response.Fail(c, err)
		return
	}

	if err := (Service{}).Update(c.Request.Context(), utils.GetUserID(c), params); err != nil {
		logger.Error("member update failed", err)
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
		logger.Error("member delete failed", err)
		response.Fail(c, err)
		return
	}
	response.Success(c, nil)
}
