package router

import "github.com/gin-gonic/gin"

func Register(r *gin.Engine) {
	api := r.Group("/api")
	Install(api)
	Auth(api)
	Admin(api)
	User(api)
}
