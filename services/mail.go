package services

import (
	"gopkg.in/gomail.v2"
)

type EmailParams struct {
	To               string
	From             string
	Subject          string
	Code             string
	PasswordResetURL string
	SMTPConfig       struct {
		Host     string
		Port     int
		UserName string
		Password string
	}
}

type IMail interface {
	Send2FAOTP() error
	SendResetPasswordEmail() error
}

func Send2FAOTP(EP EmailParams) error {
	message := gomail.NewMessage()
	message.SetHeader("From", EP.From)
	message.SetHeader("To", EP.To)
	message.SetHeader("Subject", EP.Subject)
	message.SetBody("text/plain", "Your One Time Password : "+EP.Code)

	dialer := gomail.NewDialer(
		EP.SMTPConfig.Host,
		EP.SMTPConfig.Port,
		EP.SMTPConfig.UserName,
		EP.SMTPConfig.Password)
	err := dialer.DialAndSend(message)
	return err
}

func SendResetPasswordEmail(EP EmailParams) error {
	message := gomail.NewMessage()
	message.SetHeader("From", EP.From)
	message.SetHeader("To", EP.To)
	message.SetHeader("Subject", EP.Subject)
	message.SetBody("text/plain", "Click the link to reset your password :"+EP.PasswordResetURL)

	dialer := gomail.NewDialer(
		EP.SMTPConfig.Host,
		EP.SMTPConfig.Port,
		EP.SMTPConfig.UserName,
		EP.SMTPConfig.Password)
	err := dialer.DialAndSend(message)
	return err
}
