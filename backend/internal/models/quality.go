package models

import (
	"time"
	"gorm.io/gorm"
)

// QualityStandard 质量标准
type QualityStandard struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	ProductID   uint           `json:"product_id" gorm:"not null"`
	Product     Product        `json:"product" gorm:"foreignKey:ProductID"`
	Name        string         `json:"name" gorm:"size:100;not null"`
	Type        string         `json:"type" gorm:"size:50;not null"`
	MinValue    float64        `json:"min_value"`
	MaxValue    float64        `json:"max_value"`
	TargetValue float64        `json:"target_value"`
	Unit        string         `json:"unit" gorm:"size:20"`
	Description string         `json:"description" gorm:"type:text"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// QualityInspection 质量检验记录
type QualityInspection struct {
	ID                uint           `json:"id" gorm:"primarykey"`
	ProductionOrderID uint           `json:"production_order_id" gorm:"not null"`
	ProductionOrder   ProductionOrder `json:"production_order" gorm:"foreignKey:ProductionOrderID"`
	QualityStandardID uint           `json:"quality_standard_id" gorm:"not null"`
	QualityStandard   QualityStandard `json:"quality_standard" gorm:"foreignKey:QualityStandardID"`
	InspectorID       uint           `json:"inspector_id" gorm:"not null"`
	Inspector         User           `json:"inspector" gorm:"foreignKey:InspectorID"`
	ActualValue       float64        `json:"actual_value"`
	Result            string         `json:"result" gorm:"size:20;not null"` // pass, fail
	Remark            string         `json:"remark" gorm:"type:text"`
	InspectionTime    time.Time      `json:"inspection_time" gorm:"not null"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName 指定表名
func (QualityStandard) TableName() string {
	return "quality_standards"
}

func (QualityInspection) TableName() string {
	return "quality_inspections"
}