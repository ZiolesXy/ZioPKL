package models

import "time"

type AuditLog struct {
	ID        uint     `gorm:"primaryKey" json:"id"`
	AdminID   uint     `json:"admin_id"`
	AdminName string   `json:"admin_name"`
	Action    string   `json:"action"`
	Entity    string   `json:"entity"`
	EntityID  string   `gorm:"type:varchar(50)" json:"entity_id"`
	Details   string   `gorm:"type:text" json:"details"`
	IPAddress string   `json:"ip_address"`
	CreatedAt time.Time `json:"created_at"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}