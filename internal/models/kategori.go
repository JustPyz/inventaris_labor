package models

import "time"

type Kategori struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Kategori      string    `gorm:"size:100;not null" json:"kategori"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

