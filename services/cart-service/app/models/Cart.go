package models

import (
	guuid "github.com/google/uuid"
)

type Cart struct {
	ID        guuid.UUID `gorm:"primaryKey" json:"-"`
	UserId    string     `json:"ownerId"`
	ProductId string     `json:"name"`
	CreatedAt int64      `gorm:"autoCreateTime" json:"-" `
}

type AddCartRequest struct {
	UserId    string
	ProductId string
}

type CartResponse struct {
	TotalPrice float32
	Items      []Product
}
