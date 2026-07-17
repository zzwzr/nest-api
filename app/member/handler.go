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

func (*Handler) GetInviteLink(c *gin.Context) {
	var params InviteLinkRequest
	if err := validator.Bind(c, &params); err != nil {
		response.Fail(c, err)
		return
	}

	data, err := (Service{}).GetInviteLink(c.Request.Context(), utils.GetUserID(c), params)
	if err != nil {
		logger.Error("member invite link failed", err)
		response.Fail(c, err)
		return
	}
	response.Success(c, data)
}

func (*Handler) RefreshInviteLink(c *gin.Context) {
	var params InviteLinkRequest
	if err := validator.Bind(c, &params); err != nil {
		response.Fail(c, err)
		return
	}

	data, err := (Service{}).RefreshInviteLink(c.Request.Context(), utils.GetUserID(c), params)
	if err != nil {
		logger.Error("member invite link refresh failed", err)
		response.Fail(c, err)
		return
	}
	response.Success(c, data)
}

func (*Handler) PreviewInvite(c *gin.Context) {
	var params InvitePreviewRequest
	if err := validator.Bind(c, &params); err != nil {
		response.Fail(c, err)
		return
	}

	data, err := (Service{}).PreviewInvite(c.Request.Context(), params)
	if err != nil {
		logger.Error("member invite preview failed", err)
		response.Fail(c, err)
		return
	}
	response.Success(c, data)
}

func (*Handler) AcceptInvite(c *gin.Context) {
	var params AcceptInviteRequest
	if err := validator.Bind(c, &params); err != nil {
		response.Fail(c, err)
		return
	}

	data, err := (Service{}).AcceptInvite(c.Request.Context(), utils.GetUserID(c), params)
	if err != nil {
		logger.Error("member accept invite failed", err)
		response.Fail(c, err)
		return
	}
	response.Success(c, data)
}
