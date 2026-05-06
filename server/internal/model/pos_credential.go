package model

import "time"

type POSCredential struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	BusinessID   uint       `gorm:"not null;uniqueIndex" json:"business_id"`
	APIKey       string     `gorm:"type:varchar(64);uniqueIndex;not null" json:"api_key"`
	Secret       string     `gorm:"type:varchar(128);not null" json:"-"`
	Status       string     `gorm:"type:varchar(20);default:active" json:"status"`
	LastCalledAt *time.Time `json:"last_called_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func (POSCredential) TableName() string { return "pos_credentials" }
