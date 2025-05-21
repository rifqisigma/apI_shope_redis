package helper

import (
	"fmt"
	"os"

	"gopkg.in/gomail.v2"
)

func SendEmail(toEmail, message string) {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", os.Getenv("EMAIL_SENDER"))
	mailer.SetHeader("To", toEmail)
	mailer.SetHeader("Subject", "Golang redis queue test")
	mailer.SetBody("text/html", fmt.Sprintf(`this is redis queue \n %v`, message))
	dialer := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("EMAIL_SENDER"), os.Getenv("APP_PASSWORD"))

	if err := dialer.DialAndSend(mailer); err != nil {
		fmt.Println("Error sending email:", err)
	}
}
