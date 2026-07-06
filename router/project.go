package router

import (
	"nest-api/app/project"
	"nest-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Project(rg *gin.RouterGroup) {
	h := &project.Handler{}
	r := rg.Group("/v1")
	r.Use(middleware.JWTAuthMiddleware())
	{
		r.GET("/projects", h.List)
		r.POST("/projects", h.Create)
		r.PUT("/projects", h.Update)
		r.DELETE("/projects", h.Delete)
	}
}
