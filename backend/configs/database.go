package configs

import (
	"fmt"
	"log"
	"mes-system/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	Charset  string
}

// GetDefaultDatabaseConfig 获取默认数据库配置
func GetDefaultDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:     "localhost",
		Port:     "3306",
		User:     "root",
		Password: "123456",
		DBName:   "mes_system",
		Charset:  "utf8mb4",
	}
}

// InitDatabase 初始化数据库连接
func InitDatabase(config *DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
		config.Charset,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %v", err)
	}

	// 自动迁移数据库表
	err = autoMigrate(db)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %v", err)
	}

	log.Println("Database connected and migrated successfully")
	return db, nil
}

// autoMigrate 自动迁移数据库表
func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.ProductionOrder{},
		&models.Material{},
		&models.MaterialTransaction{},
		&models.QualityStandard{},
		&models.QualityInspection{},
		&models.Equipment{},
		&models.MaintenanceRecord{},
	)
}