package model

import "time"

type WeChatAccount struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	OpenID     string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"openid"`
	BusinessID uint      `gorm:"not null;index" json:"business_id"`
	RealName   string    `gorm:"type:varchar(50);not null" json:"real_name"`
	Phone      string    `gorm:"type:varchar(20)" json:"phone"`
	Status     string    `gorm:"type:varchar(20);default:active" json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

func (WeChatAccount) TableName() string { return "wechat_accounts" }
