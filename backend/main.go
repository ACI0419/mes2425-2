package main

import (
	"log"
	"mes-system/configs"
	"mes-system/internal/controller"
	"mes-system/internal/service"
	"mes-system/pkg/jwt"
	"mes-system/routes"

	// 修正Swagger导入路径
	_ "mes-system/docs"

	"github.com/gin-gonic/gin"
	files "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title MES制造执行系统 API
// @version 1.0
// @description MES制造执行系统的RESTful API文档，包含用户管理、生产管理、产品管理、物料管理、质量管理和设备管理等模块。
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type \"Bearer\" followed by a space and JWT token.

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

	// 创建Gin引擎
	r := gin.Default()

	// 添加CORS中间件
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 添加Swagger路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(files.Handler)) // 修改这里：使用 files.Handler

	// 设置路由
	routes.SetupRoutes(r, controllers, jwtConfig)

	// 启动服务器
	log.Println("Server starting on :8080")
	log.Println("Swagger UI available at: http://localhost:8080/swagger/index.html")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
