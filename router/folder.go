package router

import (
	"nest-api/app/folder"
	"nest-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Folder(rg *gin.RouterGroup) {
	h := &folder.Handler{}
	r := rg.Group("/v1")
	r.Use(middleware.JWTAuthMiddleware())
	{
		r.GET("/folders/tree", h.Tree)
		r.POST("/folders", h.Create)
		r.PUT("/folders", h.Update)
		r.DELETE("/folders", h.Delete)
	}
}
