package models

import (
	"time"
	"gorm.io/gorm"
)

// Material 物料信息
type Material struct {
	ID           uint           `json:"id" gorm:"primarykey"`
	Code         string         `json:"code" gorm:"uniqueIndex;size:50;not null"`
	Name         string         `json:"name" gorm:"size:100;not null"`
	Type         string         `json:"type" gorm:"size:50"`           // 改为 Type，与服务层一致
	Unit         string         `json:"unit" gorm:"size:20"`
	Price        float64        `json:"price" gorm:"type:decimal(10,2)"`
	MinStock     int            `json:"min_stock" gorm:"default:0"`
	MaxStock     int            `json:"max_stock" gorm:"default:0"`
	CurrentStock int            `json:"current_stock" gorm:"default:0"`
	Description  string         `json:"description" gorm:"size:500"`    // 添加描述字段
	Status       int            `json:"status" gorm:"default:1"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

// MaterialTransaction 物料出入库记录
type MaterialTransaction struct {
	ID                uint           `json:"id" gorm:"primarykey"`
	MaterialID        uint           `json:"material_id"`
	Material          Material       `json:"material" gorm:"foreignKey:MaterialID"`
	Type              string         `json:"type" gorm:"size:20;not null"` // in:入库 out:出库
	Quantity          int            `json:"quantity" gorm:"not null"`
	Price             float64        `json:"price" gorm:"type:decimal(10,2);default:0"` // 添加单价字段
	TotalAmount       float64        `json:"total_amount" gorm:"type:decimal(12,2);default:0"` // 添加总金额字段
	Supplier          string         `json:"supplier" gorm:"size:100"`      // 添加供应商字段
	ProductionOrderID *uint          `json:"production_order_id"`           // 添加生产工单ID字段
	Remark            string         `json:"remark" gorm:"size:500"`        // 改名为 Remark，与服务层一致
	OperatorID        uint           `json:"operator_id"`
	Operator          User           `json:"operator" gorm:"foreignKey:OperatorID"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName 指定表名
func (Material) TableName() string {
	return "materials"
}

func (MaterialTransaction) TableName() string {
	return "material_transactions"
}