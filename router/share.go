package router

import (
	"nest-api/app/share"
	"nest-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Share(rg *gin.RouterGroup) {
	h := &share.Handler{}
	r := rg.Group("/v1")

	r.GET("/share/preview", h.Preview)
	r.POST("/share/content", h.AccessContent)
	r.POST("/share/interface", h.AccessDetail)

	r.Use(middleware.JWTAuthMiddleware())
	{
		r.GET("/project/shares", h.List)
		r.GET("/project/shares/detail", h.Get)
		r.POST("/project/shares", h.Create)
		r.PUT("/project/shares", h.Update)
		r.DELETE("/project/shares", h.Delete)
	}
}
