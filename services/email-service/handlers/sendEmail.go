package handlers

import (
	"email_service/utils"
	"email_service/models"
	"log"
	"net/smtp"
	"strings"
	"os"
)

func SendOTPMsg(data map[string]string) {
	from := os.Getenv("MAIL_USER")
	pass := os.Getenv("MAIL_TOKEN")
	sender := "BKMarket"
	to := strings.Split(data["to"],",")
	subject := "OTP Code From BKMarket"
	body := "Here is your OTP code: <b>"+data["otp"] + "</b>"
	msgProposal := models.Mail{
		Sender:  sender,
        To:      to,
        Subject: subject,
        Body:    body,
	}
	msg := utils.BuildMessage(msgProposal)

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		"BKMarket", to, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}
}

func SendMsg(data map[string]string) {
	from := os.Getenv("MAIL_USER")
	pass := os.Getenv("MAIL_TOKEN")
	sender := "BKMarket"
	to := strings.Split(data["to"],",")
	subject := data["subject"]
	body := data["body"]
	msgProposal := models.Mail{
		Sender:  sender,
        To:      to,
        Subject: subject,
        Body:    body,
	}
	msg := utils.BuildMessage(msgProposal)

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		"BKMarket", to, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}
}