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
		var err error

		// Сначала проверяем наличие куки
		cookie, err := c.Request.Cookie("access_token")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				// Если куки нет, проверяем токен
				err = token.ValidToken(c)
				if err != nil {
					c.AbortWithError(http.StatusUnauthorized, err)
					return
				}

				userId, err = token.ExtractTokenID(c)
				if err != nil {
					c.AbortWithError(http.StatusUnauthorized, err)
					return
				}
			} else {
				// Если произошла другая ошибка, возвращаем 500
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
		} else {
			// Если куки есть, используем её значение для аутентификации
			client := oauth2.NewClient(c.Request.Context(), oauth2.StaticTokenSource(&oauth2.Token{AccessToken: cookie.Value}))
			response, err := client.Get("https://www.googleapis.com/userinfo/v2/me")
			if err != nil {
				c.AbortWithError(http.StatusUnauthorized, err)
				return
			}
			defer response.Body.Close()

			userValue, err := models.RedisClient.Get(cookie.Value).Uint64()
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			userId = uint(userValue)
			logger.Log.Info("User Id: " + strconv.Itoa(int(userId)))
		}

		// Проверка роли пользователя после успешной аутентификации
		userRole, err := models.GetUserRoleById(userId)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if role != userRole {
			c.AbortWithError(http.StatusForbidden, errors.New("wrong role"))
			return
		}

		c.Next()
	}
}
