package models

import "time"

type Role struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Role      string    `gorm:"size:100;not null" json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
