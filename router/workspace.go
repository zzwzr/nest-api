package router

import (
	"nest-api/app/workspace"
	"nest-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Workspace(rg *gin.RouterGroup) {
	h := &workspace.Handler{}
	r := rg.Group("/v1")
	r.Use(middleware.JWTAuthMiddleware())
	{
		r.POST("/workspaces", h.Create)
		r.GET("/workspaces", h.List)
		r.GET("/workspaces/item", h.Get)
		r.PUT("/workspaces", h.Update)
		r.DELETE("/workspaces", h.Delete)
		r.PUT("/workspaces/transfer", h.TransferOwner)
	}
}
