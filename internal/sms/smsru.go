package sms

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const StatusOk = "OK"
const StatusError = "ERROR"

type Response struct {
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
	SMS        map[string]Data
	Balance    float64 `json:"balance"`
}

type Data struct {
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
	SMSId      string `json:"sms_id"`
	StatusText string `json:"status_text"`
	Cost       string `json:"cost"`
	SMSCount   int    `json:"sms"`
}

// SMSruSender реализует SMSSender для отправки сообщений через sms.ru.
type SMSruSender struct {
	APIKey string
}

func NewSMSruSender() *SMSruSender {
	return &SMSruSender{}
}

func (s *SMSruSender) SendSMSCode(phoneNumber string, code string) error {
	url := fmt.Sprintf("https://sms.ru/sms/send?api_id=%s&to=%s&msg=%s&json=1&test=%s", os.Getenv("SMS_API_KEY"), phoneNumber, code, os.Getenv("SMS_IS_TEST"))

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Обработка ответа от SMS-сервиса.
	var smsResponse Response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &smsResponse)
	if err != nil {
		return err
	}

	if smsResponse.Status == StatusError {
		return fmt.Errorf("SMS response error: %s", smsResponse.SMS[phoneNumber].StatusText)
	}

	return nil
}
