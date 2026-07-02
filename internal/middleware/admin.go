package middleware

import (
	"nest-api/internal/database"
	"nest-api/internal/ent"
	"nest-api/internal/ent/user"
	"nest-api/internal/utils"
	"nest-api/pkg/response"

	"github.com/gin-gonic/gin"
)

func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := utils.GetUserID(c)
		if userID == 0 {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		u, err := database.DB.User.
			Query().
			Where(user.IDEQ(userID)).
			Only(c.Request.Context())
		if err != nil {
			if ent.IsNotFound(err) {
				response.Unauthorized(c)
			} else {
				response.Fail(c, err)
			}
			c.Abort()
			return
		}

		if !u.IsAdmin {
			response.Forbidden(c, "需要管理员权限")
			c.Abort()
			return
		}

		c.Next()
	}
}
