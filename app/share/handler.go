package share

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
		logger.Error("share list failed", err)
		response.Fail(c, err)
		return
	}
	response.Success(c, list)
}

func (*Handler) Get(c *gin.Context) {
	var params GetRequest
	if err := validator.Bind(c, &params); err != nil {
		response.Fail(c, err)
		return
	}

	item, err := (Service{}).Get(c.Request.Context(), utils.GetUserID(c), params)
	if err != nil {
		logger.Error("share get failed", err)
		response.Fail(c, err)
		return
	}
	response.Success(c, item)
}

func (*Handler) Create(c *gin.Context) {
	var params CreateRequest
	if err := validator.Bind(c, &params); err != nil {
		response.Fail(c, err)
		return
	}

	item, err := (Service{}).Create(c.Request.Context(), utils.GetUserID(c), params)
	if err != nil {
		logger.Error("share create failed", err)
		response.Fail(c, err)
		return
	}
	response.Success(c, item)
}

func (*Handler) Update(c *gin.Context) {
	var params UpdateRequest
	if err := validator.Bind(c, &params); err != nil {
		response.Fail(c, err)
		return
	}

	item, err := (Service{}).Update(c.Request.Context(), utils.GetUserID(c), params)
	if err != nil {
		logger.Error("share update failed", err)
		response.Fail(c, err)
		return
	}
	response.Success(c, item)
}

func (*Handler) Delete(c *gin.Context) {
	var params DeleteRequest
	if err := validator.Bind(c, &params); err != nil {
		response.Fail(c, err)
		return
	}

	if err := (Service{}).Delete(c.Request.Context(), utils.GetUserID(c), params); err != nil {
		logger.Error("share delete failed", err)
		response.Fail(c, err)
		return
	}
	response.Success(c, nil)
}

func (*Handler) Preview(c *gin.Context) {
	var params PreviewRequest
	if err := validator.Bind(c, &params); err != nil {
		response.Fail(c, err)
		return
	}

	data, err := (Service{}).Preview(c.Request.Context(), params)
	if err != nil {
		logger.Error("share preview failed", err)
		response.Fail(c, err)
		return
	}
	response.Success(c, data)
}

func (*Handler) AccessContent(c *gin.Context) {
	var params AccessRequest
	if err := validator.Bind(c, &params); err != nil {
		response.Fail(c, err)
		return
	}

	data, err := (Service{}).AccessContent(c.Request.Context(), params)
	if err != nil {
		logger.Error("share access content failed", err)
		response.Fail(c, err)
		return
	}
	response.Success(c, data)
}

func (*Handler) AccessDetail(c *gin.Context) {
	var params AccessDetailRequest
	if err := validator.Bind(c, &params); err != nil {
		response.Fail(c, err)
		return
	}

	data, err := (Service{}).AccessDetail(c.Request.Context(), params)
	if err != nil {
		logger.Error("share access detail failed", err)
		response.Fail(c, err)
		return
	}
	response.Success(c, data)
}
