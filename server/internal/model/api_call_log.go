package model

import "time"

type APICallLog struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	CredentialID uint      `gorm:"index" json:"credential_id"`
	Method       string    `gorm:"type:varchar(10)" json:"method"`
	Path         string    `gorm:"type:varchar(200)" json:"path"`
	RequestBody  string    `gorm:"type:text" json:"request_body"`
	ResponseCode int       `json:"response_code"`
	IPAddress    string    `gorm:"type:varchar(50)" json:"ip_address"`
	DurationMs   int64     `json:"duration_ms"`
	CreatedAt    time.Time `json:"created_at"`
}

func (APICallLog) TableName() string { return "api_call_logs" }
