package middlewares

import (
	"eda/logger"
	"eda/models"
	"eda/utils/token"
	"errors"
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
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				c.String(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
				return
			}
			c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		if cookie.Value != "" {
			client := oauth2.NewClient(c.Request.Context(), oauth2.StaticTokenSource(&oauth2.Token{AccessToken: cookie.Value}))
			response, err := client.Get("https://www.googleapis.com/userinfo/v2/me")
			if err != nil {
				c.String(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
				return
			}
			defer response.Body.Close()

			userValue, err := models.RedisClient.Get(cookie.Value).Uint64()
			if err != nil {
				c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
				return
			}
			userId = uint(userValue)

			logger.Log.Info("User Id: " + strconv.Itoa(int(userId)))
		} else {
			err = token.ValidToken(c)
			if err != nil {
				c.String(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
				return
			}

			userId, err = token.ExtractTokenID(c)
			if err != nil {
				c.String(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
				return
			}
		}

		userRole, err := models.GetUserRoleById(userId)
		if err != nil {
			c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		if role != userRole {
			c.String(http.StatusForbidden, http.StatusText(http.StatusForbidden))
			return
		}

		c.Next()
	}
}
