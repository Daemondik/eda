package security

import (
	"eda/models"
	"eda/utils/token"
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"net/http"
	"regexp"
)

func GetUserIdByJWTOrOauth(c *gin.Context) (uint, error) {
	var userId uint
	var err error

	// Сначала пытаемся получить куки
	cookie, err := c.Request.Cookie("access_token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			// Если куки нет, пытаемся извлечь ID пользователя из JWT
			userId, err = token.ExtractTokenID(c)
			if err != nil {
				return 0, err
			}
		} else {
			// Если произошла другая ошибка, возвращаем её
			return 0, err
		}
	} else if cookie.Value != "" {
		// Если куки есть, используем её для получения ID пользователя
		client := oauth2.NewClient(c.Request.Context(), oauth2.StaticTokenSource(&oauth2.Token{AccessToken: cookie.Value}))
		response, err := client.Get("https://www.googleapis.com/userinfo/v2/me")
		if err != nil {
			return 0, err
		}
		defer response.Body.Close()

		userValue, err := models.RedisClient.Get(cookie.Value).Uint64()
		if err != nil {
			return 0, err
		}
		userId = uint(userValue)
	}

	return userId, nil
}

func IsValidRussianPhoneNumber(phone string) bool {
	// 7XXXXXXXXXX
	pattern := `^7\d{10}$`

	r := regexp.MustCompile(pattern)

	return r.MatchString(phone)
}
