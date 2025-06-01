package main

import (
	"log"
	"mes-system/configs"
	"mes-system/internal/controller"
	"mes-system/internal/service"
	"mes-system/pkg/jwt"
	"mes-system/routes"

	"github.com/gin-gonic/gin"
)

// main 主程序入口
func main() {
	// 初始化数据库
	dbConfig := configs.GetDefaultDatabaseConfig()
	db, err := configs.InitDatabase(dbConfig)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 初始化JWT配置
	jwtConfig := jwt.GetDefaultJWTConfig()

	// 初始化服务层
	userService := service.NewUserService(db, jwtConfig)
	productionService := service.NewProductionService(db)
	productService := service.NewProductService(db)
	materialService := service.NewMaterialService(db)
	qualityService := service.NewQualityService(db)
	equipmentService := service.NewEquipmentService(db)

	// 初始化控制器层
	userController := controller.NewUserController(userService)
	// 修复：ProductionController 需要两个服务参数
	productionController := controller.NewProductionController(productionService, productService)
	productController := controller.NewProductController(productService)
	materialController := controller.NewMaterialController(materialService)
	qualityController := controller.NewQualityController(qualityService)
	equipmentController := controller.NewEquipmentController(equipmentService)

	// 创建控制器集合
	controllers := &routes.Controllers{
		User:       userController,
		Production: productionController,
		Product:    productController,
		Material:   materialController,
		Quality:    qualityController,
		Equipment:  equipmentController,
	}

	// 初始化Gin引擎
	gin.SetMode(gin.ReleaseMode) // 生产模式
	r := gin.New()

	// 添加中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(corsMiddleware()) // CORS中间件

	// 设置路由
	routes.SetupRoutes(r, controllers, jwtConfig)

	// 启动服务器
	log.Println("MES System Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// corsMiddleware CORS中间件
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
