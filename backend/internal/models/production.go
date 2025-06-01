package models

import (
	"time"
	"gorm.io/gorm"
)

// ProductionOrder 生产工单
type ProductionOrder struct {
	ID           uint           `json:"id" gorm:"primarykey"`
	OrderNo      string         `json:"order_no" gorm:"uniqueIndex;size:50;not null"`
	ProductID    uint           `json:"product_id"`
	Product      Product        `json:"product" gorm:"foreignKey:ProductID"`
	Quantity     int            `json:"quantity" gorm:"not null"`
	Produced     int            `json:"produced" gorm:"default:0"`
	Status       string         `json:"status" gorm:"size:20;default:'pending'"` // pending, processing, completed, cancelled
	Priority     int            `json:"priority" gorm:"default:1"`
	StartDate    *time.Time     `json:"start_date"`
	EndDate      *time.Time     `json:"end_date"`
	CreatedBy    uint           `json:"created_by"`
	Creator      User           `json:"creator" gorm:"foreignKey:CreatedBy"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

// Product 产品信息
type Product struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	Code        string         `json:"code" gorm:"uniqueIndex;size:50;not null"`
	Name        string         `json:"name" gorm:"size:100;not null"`
	Description string         `json:"description" gorm:"type:text"`
	Unit        string         `json:"unit" gorm:"size:20"`
	Price       float64        `json:"price" gorm:"type:decimal(10,2)"`
	Status      int            `json:"status" gorm:"default:1"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName 指定表名
func (ProductionOrder) TableName() string {
	return "production_orders"
}

func (Product) TableName() string {
	return "products"
}