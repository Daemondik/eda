package middlewares

import (
	"eda/logger"
	"eda/models"
	"eda/utils/token"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"net/http"
	"strconv"
)

func AuthMiddleware(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userId uint

		// Read access token from cookie
		cookie, err := c.Request.Cookie("access_token")
		if cookie.Value != "" {
			if err != nil {
				c.String(http.StatusUnauthorized, "Unauthorized")
				return
			}

			// Use access token to fetch user info
			client := oauth2.NewClient(oauth2.NoContext, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: cookie.Value}))
			response, err := client.Get("https://www.googleapis.com/userinfo/v2/me")
			if err != nil {
				c.String(http.StatusUnauthorized, "Unauthorized")
				c.Abort()
				return
			}
			defer response.Body.Close()

			userValue, _ := models.RedisClient.Get(cookie.Value).Uint64()
			userId = uint(userValue)

			logger.Log.Error("User Id: " + strconv.Itoa(int(userId)) + "\n")
		} else {
			err = token.ValidToken(c)

			if err != nil {
				c.String(http.StatusUnauthorized, "Unauthorized")
				c.Abort()
				return
			}

			userId, err = token.ExtractTokenID(c)
			if err != nil {
				c.String(http.StatusUnauthorized, "Unauthorized")
				c.Abort()
				return
			}
		}

		if userRole, _ := models.GetUserRoleById(userId); role != userRole {
			c.String(http.StatusForbidden, "Wrong role")
			c.Abort()
			return
		}

		c.Next()
	}
}
