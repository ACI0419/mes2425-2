package models

import (
	"time"
	"gorm.io/gorm"
)

// QualityStandard 质量标准
type QualityStandard struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	Code        string         `json:"code" gorm:"uniqueIndex;size:50;not null"`
	Name        string         `json:"name" gorm:"size:100;not null"`
	Description string         `json:"description" gorm:"type:text"`
	ProductID   uint           `json:"product_id"`
	Product     Product        `json:"product" gorm:"foreignKey:ProductID"`
	Status      int            `json:"status" gorm:"default:1"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// QualityInspection 质量检验记录
type QualityInspection struct {
	ID               uint           `json:"id" gorm:"primarykey"`
	InspectionNo     string         `json:"inspection_no" gorm:"uniqueIndex;size:50;not null"`
	ProductionOrderID uint          `json:"production_order_id"`
	ProductionOrder  ProductionOrder `json:"production_order" gorm:"foreignKey:ProductionOrderID"`
	StandardID       uint           `json:"standard_id"`
	Standard         QualityStandard `json:"standard" gorm:"foreignKey:StandardID"`
	InspectedQty     int            `json:"inspected_qty" gorm:"not null"`
	PassedQty        int            `json:"passed_qty" gorm:"not null"`
	FailedQty        int            `json:"failed_qty" gorm:"not null"`
	Result           string         `json:"result" gorm:"size:20"` // pass, fail, pending
	Remark           string         `json:"remark" gorm:"type:text"`
	InspectorID      uint           `json:"inspector_id"`
	Inspector        User           `json:"inspector" gorm:"foreignKey:InspectorID"`
	InspectedAt      time.Time      `json:"inspected_at"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName 指定表名
func (QualityStandard) TableName() string {
	return "quality_standards"
}

func (QualityInspection) TableName() string {
	return "quality_inspections"
}