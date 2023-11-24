package models

import (
	guuid "github.com/google/uuid"
)

type User struct {
	ID       guuid.UUID `gorm:"primaryKey" json:"-"`
	Fullname string     `json:"fullname"`
	Username string     `json:"username"`
	Email    string     `json:"email"`
	Role     string     `json:"role"`
	Picture  string     `json:"picture"`
	Password string     `json:"-"`
	// Sessions  []Session  `gorm:"foreignKey:UserRefer; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;" json:"-"`
	// Products  []Product  `gorm:"foreignKey:UserRefer; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;" json:"-"`
	CreatedAt int64 `gorm:"autoCreateTime" json:"-" `
	UpdatedAt int64 `gorm:"autoUpdateTime:milli" json:"-"`
}
