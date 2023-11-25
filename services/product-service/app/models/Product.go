package models

import (
	guuid "github.com/google/uuid"
)

type Product struct {
	ID        guuid.UUID `gorm:"primaryKey" json:"-"`
	OwnerId   string     `json:"ownerId"`
	Name      string     `json:"name"`
	Price     float32    `json:"price"`
	Picture   string     `json:"picture"`
	Detail    string     `json:"detail"`
	CreatedAt int64      `gorm:"autoCreateTime" json:"-" `
	UpdatedAt int64      `gorm:"autoUpdateTime:milli" json:"-"`
}

type ProductForm struct {
	UserId  string  `json:"userId"`
	Name    string  `json:"name"`
	Price   float32 `json:"price"`
	Picture string  `json:"picture"`
	Detail  string  `json:"detail"`
}
