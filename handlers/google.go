package handlers

import (
	"eda/logger"
	"eda/models"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

var (
	oauthConfGl = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:8080/api/callback-gl",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
	oauthStateStringGl = "random"
)

func GoogleLogin(c *gin.Context) {
	authCodeURL := oauthConfGl.AuthCodeURL(oauthStateStringGl)
	c.JSON(http.StatusOK, gin.H{"url": authCodeURL})
}

func CallBackFromGoogle(c *gin.Context) {
	state := c.Query("state")
	if state != oauthStateStringGl {
		logger.Log.Error("invalid oauth state, expected '" + oauthStateStringGl + "', got '" + state + "'\n")
		return
	}

	code := c.Query("code")
	token, err := oauthConfGl.Exchange(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "code exchange failed:" + err.Error()})
		return
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + url.QueryEscape(token.AccessToken))
	if err != nil {
		logger.Log.Error("Get: " + err.Error() + "\n")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	defer resp.Body.Close()

	var profile struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		logger.Log.Error("Decode: " + err.Error() + "\n")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var u models.User

	u, err = models.GetUserByEmail(profile.Email)
	if err != nil {
		logger.Log.Error("User Exist: " + err.Error() + "\n")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if u.ID == 0 {
		u.Email = profile.Email
		u.IsActive = true
		u, err = u.SaveUser()
		if err != nil {
			logger.Log.Error("Saving User: " + err.Error() + "\n")
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
	}

	// Set a cookie with the access token
	expiration := time.Now().Add(24 * time.Hour) // Adjust expiration as needed
	cookie := http.Cookie{Name: "access_token", Value: token.AccessToken, Expires: expiration, Path: "/", Secure: true, HttpOnly: true}
	http.SetCookie(c.Writer, &cookie)

	status := models.RedisClient.Set(token.AccessToken, strconv.Itoa(int(u.ID)), expiration.Sub(time.Now()))
	logger.Log.Info("Setting redis key: " + status.String() + "\n")
}
