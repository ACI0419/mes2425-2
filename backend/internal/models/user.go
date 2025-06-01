package models

import (
	"time"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	Username  string         `json:"username" gorm:"uniqueIndex;size:50;not null"`
	Password  string         `json:"-" gorm:"size:255;not null"`
	Email     string         `json:"email" gorm:"uniqueIndex;size:100"`
	RealName  string         `json:"real_name" gorm:"size:50"`
	Phone     string         `json:"phone" gorm:"size:20"`
	Role      string         `json:"role" gorm:"size:20;default:'user'"`
	Status    int            `json:"status" gorm:"default:1"` // 1:启用 0:禁用
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}