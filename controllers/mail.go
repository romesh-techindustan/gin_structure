package controllers

import (
	"fmt"

	"gopkg.in/gomail.v2"
)


func SendEmail(to ,code  string) {
	m := gomail.NewMessage()
	m.SetHeader("From", "romeshkhaba@gmail.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", "OTP for two factor authentication")
	m.SetBody("text/plain", code)

	d := gomail.NewDialer("smtp.gmail.com", 587, "romeshkhaba@gmail.com", "xjlwvvlqfkzqhbvi")

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
	fmt.Println("Email Sent successfully")
}
