package models

import (
	"time"
	"gorm.io/gorm"
)

// Material 物料信息
type Material struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	Code        string         `json:"code" gorm:"uniqueIndex;size:50;not null"`
	Name        string         `json:"name" gorm:"size:100;not null"`
	Category    string         `json:"category" gorm:"size:50"`
	Unit        string         `json:"unit" gorm:"size:20"`
	Price       float64        `json:"price" gorm:"type:decimal(10,2)"`
	MinStock    int            `json:"min_stock" gorm:"default:0"`
	MaxStock    int            `json:"max_stock" gorm:"default:0"`
	CurrentStock int           `json:"current_stock" gorm:"default:0"`
	Status      int            `json:"status" gorm:"default:1"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// MaterialTransaction 物料出入库记录
type MaterialTransaction struct {
	ID         uint           `json:"id" gorm:"primarykey"`
	MaterialID uint           `json:"material_id"`
	Material   Material       `json:"material" gorm:"foreignKey:MaterialID"`
	Type       string         `json:"type" gorm:"size:20;not null"` // in:入库 out:出库
	Quantity   int            `json:"quantity" gorm:"not null"`
	Reason     string         `json:"reason" gorm:"size:100"`
	OperatorID uint           `json:"operator_id"`
	Operator   User           `json:"operator" gorm:"foreignKey:OperatorID"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName 指定表名
func (Material) TableName() string {
	return "materials"
}

func (MaterialTransaction) TableName() string {
	return "material_transactions"
}