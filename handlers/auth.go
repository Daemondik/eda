package handlers

import (
	"eda/logger"
	"eda/models"
	"eda/utils/security"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type SMSResponse struct {
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
	SMS        []SMS
	Balance    float64 `json:"balance"`
}

type SMS struct {
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
	SMSId      string `json:"sms_id"`
	StatusText string `json:"status_text"`
}

const SMSStatusOk = "OK"
const SMSStatusError = "ERROR"

func CurrentUser(c *gin.Context) {
	userId, err := security.GetUserIdByJWTOrOauth(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := models.GetUserByID(userId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": u})
}

type LoginInput struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := models.User{}

	if phoneValid := security.IsValidRussianPhoneNumber(input.Phone); !phoneValid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "phone should be format +7 (XXX) XXX-XX-XX"})
		return
	}
	u.Phone = input.Phone
	u.Password = input.Password

	generatedToken, err := models.LoginCheck(u.Phone, u.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "phone or password is incorrect."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": generatedToken})
}

type RegisterInput struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {

	var input RegisterInput
	var err error

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := models.User{}

	if phoneValid := security.IsValidRussianPhoneNumber(input.Phone); !phoneValid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "phone should be format +7 (XXX) XXX-XX-XX"})
		return
	}

	code := rand.Int() + rand.Int() + rand.Int() + rand.Int()

	url := fmt.Sprintf("https://sms.ru/sms/send?api_id=%s&to=%s&msg=Code: %d&json=1&test=1", os.Getenv("SMS_API_KEY"), input.Phone, code)

	resp, err := http.Get(url)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("sending error: %s", err.Error()))
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("response error: %s", err.Error()))
		return
	}

	var smsResponse SMSResponse

	err = json.Unmarshal(body, &smsResponse)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("unmarshal error: %s", err.Error()))
		return
	}

	if smsResponse.Status == SMSStatusError {
		c.String(http.StatusBadRequest, fmt.Sprintf("unmarshal error: %s", smsResponse.SMS[0].StatusText))
		return
	}

	u.Phone = input.Phone
	u.Password = input.Password
	u.IsActive = false

	_, err = u.SaveUser()

	expiration := time.Now().Add(time.Hour)
	models.RedisClient.Set(u.Phone, code, expiration.Sub(time.Now()))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}

type ConfirmSMSCodeInput struct {
	Phone string `json:"phone" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

func ConfirmSMSCode(c *gin.Context) {
	var input ConfirmSMSCodeInput
	var err error

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currentCode := models.RedisClient.Get(input.Phone)
	models.RedisClient.Del(input.Phone)

	if input.Code != currentCode.String() {
		c.String(http.StatusBadRequest, "incorrect code")
		return
	}

	u, err := models.GetUserByPhone(input.Phone)
	if err != nil {
		logger.Log.Error("User Exist: " + err.Error() + "\n")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	_, err = u.SetActive()
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "registration success"})
}
