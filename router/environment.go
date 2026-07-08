package router

import (
	"nest-api/app/environment"
	"nest-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Environment(rg *gin.RouterGroup) {
	h := &environment.Handler{}
	r := rg.Group("/v1")
	r.Use(middleware.JWTAuthMiddleware())
	{
		r.GET("/environments", h.List)
		r.POST("/environments", h.Create)
		r.PUT("/environments", h.Update)
		r.DELETE("/environments", h.Delete)
	}
}
