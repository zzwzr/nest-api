package router

import (
	"nest-api/app/auth"
	"nest-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Auth(rg *gin.RouterGroup) {
	h := &auth.Handler{}
	r := rg.Group("/v1/auth")
	{
		r.GET("/site", h.Site)
		r.POST("/register", h.Register)
		r.POST("/login", h.Login)
		r.GET("/me", middleware.JWTAuthMiddleware(), h.Me)
	}
}
