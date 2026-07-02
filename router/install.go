package router

import (
	"nest-api/app/install"

	"github.com/gin-gonic/gin"
)

func Install(rg *gin.RouterGroup) {
	h := &install.Handler{}
	r := rg.Group("/v1/install")
	{
		r.GET("/status", h.Status)
		r.POST("/test-database", h.TestDatabase)
		r.POST("", h.Install)
	}
}
