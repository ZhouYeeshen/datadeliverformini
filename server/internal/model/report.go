package model

import (
	"time"

	"gorm.io/gorm"
)

type Report struct {
	ID                     uint           `gorm:"primaryKey" json:"id"`
	BusinessID             uint           `gorm:"not null;index" json:"business_id"`
	BusinessName           string         `gorm:"->" json:"business_name,omitempty"`
	ReportMonth            string         `gorm:"type:varchar(7);not null;index" json:"report_month"`
	RestaurantRevenue      float64        `gorm:"type:decimal(15,2);default:0" json:"restaurant_revenue"`
	RetailRevenue          float64        `gorm:"type:decimal(15,2);default:0" json:"retail_revenue"`
	AccommodationRevenue   float64        `gorm:"type:decimal(15,2);default:0" json:"accommodation_revenue"`
	TobaccoAlcoholRevenue  float64        `gorm:"type:decimal(15,2);default:0" json:"tobacco_alcohol_revenue"`
	OtherRevenue           float64        `gorm:"type:decimal(15,2);default:0" json:"other_revenue"`
	TotalRevenue           float64        `gorm:"type:decimal(15,2);default:0" json:"total_revenue"`
	SourceType             string         `gorm:"type:varchar(20);not null" json:"source_type"`
	PhotoURLs              string         `gorm:"type:text" json:"photo_urls"`
	Status                 string         `gorm:"type:varchar(20);default:submitted" json:"status"`
	SubmittedBy            uint           `gorm:"default:0" json:"submitted_by"`
	SubmittedAt            time.Time      `json:"submitted_at"`
	CreatedAt              time.Time      `json:"created_at"`
	UpdatedAt              time.Time      `json:"updated_at"`
	DeletedAt              gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Report) TableName() string { return "reports" }
