package security

import (
	"eda/models"
	"eda/utils/token"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"regexp"
)

func GetUserIdByJWTOrOauth(c *gin.Context) (uint, error) {
	var userId uint
	var err error

	cookie, err := c.Request.Cookie("access_token")
	if err != nil {
		return 0, err
	}
	if cookie.Value != "" {
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
	} else {
		userId, err = token.ExtractTokenID(c)
		if err != nil {
			return 0, err
		}
	}

	return userId, nil
}

func IsValidRussianPhoneNumber(phone string) bool {
	// +7 (XXX) XXX-XX-XX
	pattern := `^\+7\s\(\d{3}\)\s\d{3}-\d{2}-\d{2}$`

	r := regexp.MustCompile(pattern)

	return r.MatchString(phone)
}
