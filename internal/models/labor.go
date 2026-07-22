package models

import "time"

type Labor struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Labor     string    `gorm:"size:100;not null" json:"labor"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

