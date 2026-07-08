package router

import (
	"nest-api/app/envvariable"
	"nest-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func EnvVariable(rg *gin.RouterGroup) {
	h := &envvariable.Handler{}
	r := rg.Group("/v1")
	r.Use(middleware.JWTAuthMiddleware())
	{
		r.GET("/variables", h.List)
		r.POST("/variables", h.Create)
		r.PUT("/variables", h.Update)
		r.DELETE("/variables", h.Delete)
	}
}
