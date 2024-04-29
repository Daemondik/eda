package middlewares

import (
	"eda/models"
	"eda/utils/token"
	"github.com/gin-gonic/gin"
	"net/http"
)

func JwtAuthMiddleware(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := token.ValidToken(c)
		if err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}
		userId, err := token.ExtractTokenID(c)
		if err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}
		if userRole, _ := models.GetUserRole(userId); role != userRole {
			c.String(http.StatusForbidden, "Wrong role")
			c.Abort()
			return
		}

		c.Next()
	}
}
