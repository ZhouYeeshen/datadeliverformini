package model

import (
	"time"
)

type Business struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"type:varchar(200);not null;index" json:"name"`
	LicenseNo    string    `gorm:"type:varchar(100);index" json:"license_no"`
	LegalPerson  string    `gorm:"type:varchar(50);not null" json:"legal_person"`
	IndustryType string    `gorm:"type:varchar(50);not null" json:"industry_type"`
	Address      string    `gorm:"type:varchar(500)" json:"address"`
	ContactPhone string    `gorm:"type:varchar(20)" json:"contact_phone"`
	Status       string    `gorm:"type:varchar(20);default:active" json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (Business) TableName() string { return "businesses" }
