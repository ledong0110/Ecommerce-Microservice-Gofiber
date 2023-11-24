package models

type MessageForm struct {
	From string `json:"from"`
	To string  `json:"to"`
	Task string `json:"task"`
	Content map[string]string `json:"content"`
}