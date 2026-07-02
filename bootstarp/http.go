package bootstrap

import (
	"fmt"

	"nest-api/configs"
	"nest-api/internal/middleware"
	"nest-api/internal/utils"
	"nest-api/pkg/response"
	"nest-api/router"

	"github.com/gin-gonic/gin"
)

func InitHTTP() {
	gin.SetMode(utils.String("GIN_MODE", gin.DebugMode))

	r := gin.New()

	_ = r.SetTrustedProxies(nil)

	r.Use(
		gin.Logger(),
		middleware.CORS(),
		middleware.Recovery(),
	)

	// 开启405错误
	r.HandleMethodNotAllowed = true

	r.NoRoute(func(c *gin.Context) {
		response.NotFound(c)
	})

	r.NoMethod(func(c *gin.Context) {
		response.MethodNotAllowed(c)
	})

	router.Register(r)

	addr := fmt.Sprintf(
		"%s:%d",
		configs.App.Server.Host,
		configs.App.Server.Port,
	)

	r.Run(addr)
}
