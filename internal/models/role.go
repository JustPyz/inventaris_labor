package models

type Role struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Role string `gorm:"size:100;not null" json:"role"`
}
