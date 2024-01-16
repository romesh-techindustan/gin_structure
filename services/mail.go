package services

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

type EmailParams struct {
	To               string
	Code             string
	PasswordResetURL string
}

func Send2FAOTP(EP EmailParams) {
	m := gomail.NewMessage()
	m.SetHeader("From", "romeshkhaba@gmail.com")
	m.SetHeader("To", EP.To)
	m.SetHeader("Subject", "OTP for two factor authentication")
	m.SetBody("text/plain", "Your One Time Password : "+EP.Code)

	d := gomail.NewDialer("smtp.gmail.com", 587, "romeshkhaba@gmail.com", "xjlwvvlqfkzqhbvi")

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
	fmt.Println("Email Sent successfully")
}

func SendResetPasswordEmail(EP EmailParams) {
	m := gomail.NewMessage()
	m.SetHeader("From", "romeshkhaba@gmail.com")
	m.SetHeader("To", EP.To)
	m.SetHeader("Subject", "OTP for two factor authentication")
	m.SetBody("text/plain", "Click the link to reset your password :"+EP.PasswordResetURL)

	d := gomail.NewDialer("smtp.gmail.com", 587, "romeshkhaba@gmail.com", "xjlwvvlqfkzqhbvi")

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
	fmt.Println("Email Sent successfully")
}
