package models

import (
	"time"
	"gorm.io/gorm"
)

// Equipment 设备信息
type Equipment struct {
	ID           uint           `json:"id" gorm:"primarykey"`
	Code         string         `json:"code" gorm:"uniqueIndex;size:50;not null"`
	Name         string         `json:"name" gorm:"size:100;not null"`
	Type         string         `json:"type" gorm:"size:50"`
	Model        string         `json:"model" gorm:"size:50"`
	Manufacturer string         `json:"manufacturer" gorm:"size:100"`
	PurchaseDate *time.Time     `json:"purchase_date"`
	WarrantyDate *time.Time     `json:"warranty_date"`
	Location     string         `json:"location" gorm:"size:100"`
	Status       string         `json:"status" gorm:"size:20;default:'running';not null"` // running, stopped, maintenance, fault
	Description  string         `json:"description" gorm:"type:text"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

// MaintenanceRecord 维护记录
type MaintenanceRecord struct {
	ID              uint           `json:"id" gorm:"primarykey"`
	EquipmentID     uint           `json:"equipment_id" gorm:"not null"`
	Equipment       Equipment      `json:"equipment" gorm:"foreignKey:EquipmentID"`
	MaintainerID    uint           `json:"maintainer_id" gorm:"not null"`
	Maintainer      User           `json:"maintainer" gorm:"foreignKey:MaintainerID"`
	Type            string         `json:"type" gorm:"size:20;not null"` // preventive, corrective, emergency
	Description     string         `json:"description" gorm:"type:text;not null"`
	StartTime       time.Time      `json:"start_time" gorm:"not null"`
	EndTime         *time.Time     `json:"end_time"`
	Cost            float64        `json:"cost" gorm:"type:decimal(10,2);default:0"`
	PartsReplaced   string         `json:"parts_replaced" gorm:"type:text"`
	Result          string         `json:"result" gorm:"type:text"`
	NextMaintenance *time.Time     `json:"next_maintenance"`
	Remark          string         `json:"remark" gorm:"type:text"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName 指定表名
func (Equipment) TableName() string {
	return "equipment"
}

func (MaintenanceRecord) TableName() string {
	return "maintenance_records"
}