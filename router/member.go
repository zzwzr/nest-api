package router

import (
	"nest-api/app/member"
	"nest-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Member(rg *gin.RouterGroup) {
	h := &member.Handler{}
	r := rg.Group("/v1")
	r.GET("/members/invite/preview", h.PreviewInvite)

	r.Use(middleware.JWTAuthMiddleware())
	{
		r.GET("/members", h.List)
		r.POST("/members", h.Invite)
		r.PUT("/members", h.Update)
		r.DELETE("/members", h.Delete)

		r.GET("/members/invite/link", h.GetInviteLink)
		r.POST("/members/invite/link/refresh", h.RefreshInviteLink)
		r.POST("/members/accept/invite", h.AcceptInvite)
	}
}
