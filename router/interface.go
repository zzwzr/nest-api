package router

import (
	"nest-api/app/interfaces"
	"nest-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Interface(rg *gin.RouterGroup) {
	h := &interfaces.Handler{}
	r := rg.Group("/v1")
	r.Use(middleware.JWTAuthMiddleware())
	{
		r.GET("/interfaces", h.List)
		r.GET("/interfaces/detail", h.Detail)
		r.POST("/interfaces", h.Create)
		r.PUT("/interfaces", h.Update)
		r.PUT("/interfaces/reorder", h.Reorder)
		r.DELETE("/interfaces", h.Delete)
	}
}
