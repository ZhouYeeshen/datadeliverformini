package model

import "time"

type Alert struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	BusinessID     uint      `gorm:"not null;index" json:"business_id"`
	BusinessName   string    `gorm:"type:varchar(200)" json:"business_name"`
	AlertType      string    `gorm:"type:varchar(30);not null" json:"alert_type"`
	AlertMonth     string    `gorm:"type:varchar(7);not null" json:"alert_month"`
	ReferenceValue float64   `gorm:"type:decimal(15,2)" json:"reference_value"`
	CurrentValue   float64   `gorm:"type:decimal(15,2)" json:"current_value"`
	ChangePercent  float64   `gorm:"type:decimal(10,2)" json:"change_percent"`
	IsRead         bool      `gorm:"default:false" json:"is_read"`
	CreatedAt      time.Time `json:"created_at"`
}

func (Alert) TableName() string { return "alerts" }
