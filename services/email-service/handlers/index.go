package handlers

import (
	"email_service/models"
	"encoding/json"
	"log"
)

func Route(raw_data []byte) {
	msg := models.MessageForm{}
	json.Unmarshal(raw_data, &msg)
	log.Printf("Received a message from: %s", msg.From)
	if msg.Task == "OTP"{
		SendOTPMsg(msg.Content)
	} else if msg.Task == "Normal" {
		SendMsg(msg.Content)
	}
}