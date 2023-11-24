package utils

import (
	"fmt"
	"email_service/models"
	"strings"
)


func BuildMessage(mail models.Mail) string {
	
    msg := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
    msg += fmt.Sprintf("From: %s\r\n", mail.Sender)
    msg += fmt.Sprintf("To: %s\r\n", strings.Join(mail.To, ";"))
    msg += fmt.Sprintf("Subject: %s\r\n", mail.Subject)
    msg += fmt.Sprintf("\r\n%s\r\n", mail.Body)

    return msg
}