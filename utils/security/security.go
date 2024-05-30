package security

import (
	"eda/models"
	"eda/utils/token"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

func GetUserIdByJWTOrOauth(c *gin.Context) (uint, error) {
	var userId uint
	var err error

	cookie, err := c.Request.Cookie("access_token")
	if cookie.Value != "" {
		if err != nil {
			return 0, err
		}

		client := oauth2.NewClient(oauth2.NoContext, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: cookie.Value}))
		response, err := client.Get("https://www.googleapis.com/userinfo/v2/me")
		if err != nil {
			return 0, err
		}
		defer response.Body.Close()

		userValue, _ := models.RedisClient.Get(cookie.Value).Uint64()
		userId = uint(userValue)
	} else {
		userId, err = token.ExtractTokenID(c)
		if err != nil {
			return 0, err
		}
	}

	return userId, nil
}
