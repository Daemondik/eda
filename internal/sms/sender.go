package sms

type Sender interface {
	SendSMSCode(phoneNumber string, code string) error
}
