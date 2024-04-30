package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"net/http"
	"os"
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

//func GoogleLogin(c *gin.Context) {
//	HandleLogin(c, oauthConfGl, oauthStateStringGl)
//}
//
//func HandleLogin(c *gin.Context, oauthConf *oauth2.Config, oauthStateString string) {
//	URL, err := url.Parse(oauthConf.Endpoint.AuthURL)
//	if err != nil {
//		logger.Log.Error("Parse: " + err.Error())
//	}
//	logger.Log.Info(URL.String())
//	parameters := url.Values{}
//	parameters.Add("client_id", oauthConf.ClientID)
//	parameters.Add("scope", strings.Join(oauthConf.Scopes, " "))
//	parameters.Add("redirect_uri", oauthConf.RedirectURL)
//	parameters.Add("response_type", "code")
//	parameters.Add("state", oauthStateString)
//	URL.RawQuery = parameters.Encode()
//	url := URL.String()
//	logger.Log.Info(url)
//	c.JSON(http.StatusOK, gin.H{"data": url})
//}

func GoogleLogin(c *gin.Context) {
	url := oauthConfGl.AuthCodeURL(oauthStateStringGl)
	c.JSON(http.StatusOK, gin.H{"url": url})
}

func CallBackFromGoogle(c *gin.Context) {
	state := c.Query("state")
	if state != oauthStateStringGl {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateStringGl, state)
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	code := c.Query("code")
	token, err := oauthConfGl.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Printf("code exchange failed: %s\n", err.Error())
		return
	}

	// Set a cookie with the access token
	expiration := time.Now().Add(24 * time.Hour) // Adjust expiration as needed
	cookie := http.Cookie{Name: "access_token", Value: token.AccessToken, Expires: expiration}
	http.SetCookie(c.Writer, &cookie)
}

func Profile(c *gin.Context) {
	// Read access token from cookie
	cookie, err := c.Request.Cookie("access_token")
	if err != nil {
		fmt.Println("Access token cookie not found.")
		return
	}

	// Use access token to fetch user info
	client := oauth2.NewClient(oauth2.NoContext, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: cookie.Value}))
	response, err := client.Get("https://www.googleapis.com/userinfo/v2/me")
	if err != nil {
		fmt.Printf("failed getting user info: %s\n", err.Error())
		return
	}
	defer response.Body.Close()

	// Display user info
	c.String(http.StatusOK, "User Info:\n"+
		"Status: %s\n"+
		"Headers: %v\n", response.Status, response.Header)

	// You may want to parse and display the response body here
}

//func CallBackFromGoogle(c *gin.Context) {
//	logger.Log.Info("Callback-gl..")
//
//	state := c.Request.FormValue("state")
//	logger.Log.Info(state)
//	if state != oauthStateStringGl {
//		logger.Log.Info("invalid oauth state, expected " + oauthStateStringGl + ", got " + state + "\n")
//		return
//	}
//
//	code := c.Request.FormValue("code")
//	logger.Log.Info(code)
//
//	if code == "" {
//		logger.Log.Warn("Code not found..")
//		c.JSON(http.StatusBadRequest, gin.H{"message": "Code Not Found to provide AccessToken.."})
//		reason := c.Request.FormValue("error_reason")
//		if reason == "user_denied" {
//			c.JSON(http.StatusBadRequest, gin.H{"message": "User has denied Permission"})
//		}
//	} else {
//		token, err := oauthConfGl.Exchange(oauth2.NoContext, code)
//		if err != nil {
//			logger.Log.Error("oauthConfGl.Exchange() failed with " + err.Error() + "\n")
//			return
//		}
//		logger.Log.Info("TOKEN>> AccessToken>> " + token.AccessToken)
//		logger.Log.Info("TOKEN>> Expiration Time>> " + token.Expiry.String())
//		logger.Log.Info("TOKEN>> RefreshToken>> " + token.RefreshToken)
//
//		resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + url.QueryEscape(token.AccessToken))
//		if err != nil {
//			logger.Log.Error("Get: " + err.Error() + "\n")
//			return
//		}
//		defer resp.Body.Close()
//
//		response, err := io.ReadAll(resp.Body)
//		if err != nil {
//			logger.Log.Error("ReadAll: " + err.Error() + "\n")
//			return
//		}
//
//		logger.Log.Info("parseResponseBody: " + string(response) + "\n")
//
//		var profile struct {
//			ID    string `json:"id"`
//			Email string `json:"email"`
//		}
//		if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
//			c.AbortWithError(http.StatusBadRequest, err)
//			return
//		}
//
//		c.JSON(http.StatusOK, gin.H{
//			"ID":    profile.ID,
//			"Email": profile.Email,
//		})
//
//		return
//	}
//}
