package routes

import (
	"mes-system/internal/controller"
	"mes-system/internal/middleware"
	"mes-system/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// Controllers 控制器集合
type Controllers struct {
	User       *controller.UserController
	Production *controller.ProductionController
	Product    *controller.ProductController
	Material   *controller.MaterialController
	Quality    *controller.QualityController
	Equipment  *controller.EquipmentController
}

// SetupRoutes 设置所有路由
func SetupRoutes(r *gin.Engine, controllers *Controllers, jwtConfig *jwt.JWTConfig) {
	// API版本分组
	v1 := r.Group("/api/v1")

	// 设置用户相关路由（无需认证）
	setupUserRoutes(v1, controllers.User)

	// 需要认证的路由组
	auth := v1.Group("")
	auth.Use(middleware.AuthMiddleware(jwtConfig))
	{
		// 设置生产管理路由
		setupProductionRoutes(auth, controllers.Production)

		// 设置产品管理路由
		setupProductRoutes(auth, controllers.Product)

		// 设置物料管理路由
		setupMaterialRoutes(auth, controllers.Material)

		// 设置质量管理路由
		setupQualityRoutes(auth, controllers.Quality)

		// 设置设备管理路由
		setupEquipmentRoutes(auth, controllers.Equipment)

		// 设置需要认证的用户路由
		setupAuthUserRoutes(auth, controllers.User)
	}
}

// setupUserRoutes 设置用户相关路由（无需认证）
func setupUserRoutes(rg *gin.RouterGroup, ctrl *controller.UserController) {
	userGroup := rg.Group("/users")
	{
		userGroup.POST("/login", ctrl.Login)       // 用户登录
		userGroup.POST("/register", ctrl.Register) // 用户注册
	}
}

// setupAuthUserRoutes 设置需要认证的用户路由
func setupAuthUserRoutes(rg *gin.RouterGroup, ctrl *controller.UserController) {
	userGroup := rg.Group("/users")
	{
		userGroup.GET("/profile", ctrl.GetProfile)                    // 获取用户信息
		userGroup.PUT("/profile", ctrl.UpdateProfile)                 // 更新用户信息
		userGroup.PUT("/password", ctrl.ChangePassword)               // 修改密码
		userGroup.POST("/refresh", ctrl.RefreshToken)                 // 刷新令牌
		userGroup.GET("/list", middleware.RoleMiddleware("admin"), ctrl.GetUserList) // 获取用户列表（仅管理员）
	}
}

// setupProductionRoutes 设置生产管理路由
func setupProductionRoutes(rg *gin.RouterGroup, ctrl *controller.ProductionController) {
	productionGroup := rg.Group("/production")
	{
		// 生产工单管理
		productionGroup.POST("/orders", ctrl.CreateProductionOrder)           // 创建生产工单
		productionGroup.GET("/orders/:id", ctrl.GetProductionOrder)           // 获取生产工单详情
		productionGroup.GET("/orders", ctrl.GetProductionOrderList)           // 获取生产工单列表
		productionGroup.PUT("/orders/:id", ctrl.UpdateProductionOrder)        // 更新生产工单
		productionGroup.DELETE("/orders/:id", ctrl.DeleteProductionOrder)     // 删除生产工单
		productionGroup.GET("/statistics", ctrl.GetProductionStatistics)      // 获取生产统计
	}
}

// setupProductRoutes 设置产品管理路由
func setupProductRoutes(rg *gin.RouterGroup, ctrl *controller.ProductController) {
	productGroup := rg.Group("/products")
	{
		productGroup.POST("", ctrl.CreateProduct)        // 创建产品
		productGroup.GET("/:id", ctrl.GetProduct)        // 获取产品详情
		productGroup.GET("", ctrl.GetProductList)        // 获取产品列表
		productGroup.PUT("/:id", ctrl.UpdateProduct)     // 更新产品
		productGroup.DELETE("/:id", ctrl.DeleteProduct)  // 删除产品
		productGroup.GET("/all", ctrl.GetAllProducts)    // 获取所有产品（用于下拉选择）
	}
}

// setupMaterialRoutes 设置物料管理路由
func setupMaterialRoutes(rg *gin.RouterGroup, ctrl *controller.MaterialController) {
	materialGroup := rg.Group("/materials")
	{
		// 物料信息管理
		materialGroup.POST("", ctrl.CreateMaterial)                    // 创建物料
		materialGroup.GET("/:id", ctrl.GetMaterial)                    // 获取物料详情
		materialGroup.GET("", ctrl.GetMaterialList)                    // 获取物料列表
		materialGroup.PUT("/:id", ctrl.UpdateMaterial)                 // 更新物料
		materialGroup.DELETE("/:id", ctrl.DeleteMaterial)              // 删除物料

		// 物料交易管理
		materialGroup.POST("/transactions", ctrl.CreateTransaction)     // 创建物料交易
		materialGroup.GET("/transactions", ctrl.GetTransactionList)     // 获取交易列表

		// 库存管理
		materialGroup.GET("/low-stock", ctrl.GetLowStockMaterials)     // 获取低库存物料
		materialGroup.GET("/types", ctrl.GetMaterialTypes)             // 获取物料类型
	}
}

// setupQualityRoutes 设置质量管理路由
func setupQualityRoutes(rg *gin.RouterGroup, ctrl *controller.QualityController) {
	qualityGroup := rg.Group("/quality")
	{
		// 质量标准管理
		qualityGroup.POST("/standards", ctrl.CreateQualityStandard)        // 创建质量标准
		qualityGroup.GET("/standards/:id", ctrl.GetQualityStandard)        // 获取质量标准详情
		qualityGroup.GET("/standards", ctrl.GetQualityStandardList)        // 获取质量标准列表
		qualityGroup.PUT("/standards/:id", ctrl.UpdateQualityStandard)     // 更新质量标准
		qualityGroup.DELETE("/standards/:id", ctrl.DeleteQualityStandard)  // 删除质量标准

		// 质量检测管理
		qualityGroup.POST("/inspections", ctrl.CreateQualityInspection)     // 创建质量检测
		qualityGroup.GET("/inspections/:id", ctrl.GetQualityInspection)     // 获取质量检测详情
		qualityGroup.GET("/inspections", ctrl.GetQualityInspectionList)     // 获取质量检测列表
		qualityGroup.PUT("/inspections/:id", ctrl.UpdateQualityInspection)  // 更新质量检测
		qualityGroup.DELETE("/inspections/:id", ctrl.DeleteQualityInspection) // 删除质量检测

		// 质量统计
		qualityGroup.GET("/statistics", ctrl.GetQualityStatistics)          // 获取质量统计
	}
}

// setupEquipmentRoutes 设置设备管理路由
func setupEquipmentRoutes(rg *gin.RouterGroup, ctrl *controller.EquipmentController) {
	equipmentGroup := rg.Group("/equipment")
	{
		// 设备信息管理
		equipmentGroup.POST("", ctrl.CreateEquipment)                    // 创建设备
		equipmentGroup.GET("/:id", ctrl.GetEquipment)                    // 获取设备详情
		equipmentGroup.GET("", ctrl.GetEquipmentList)                    // 获取设备列表
		equipmentGroup.PUT("/:id", ctrl.UpdateEquipment)                 // 更新设备
		equipmentGroup.DELETE("/:id", ctrl.DeleteEquipment)              // 删除设备

		// 维护记录管理
		equipmentGroup.POST("/maintenance", ctrl.CreateMaintenanceRecord)   // 创建维护记录
		equipmentGroup.GET("/maintenance/:id", ctrl.GetMaintenanceRecord)   // 获取维护记录详情
		equipmentGroup.GET("/maintenance", ctrl.GetMaintenanceRecordList)   // 获取维护记录列表
		equipmentGroup.PUT("/maintenance/:id", ctrl.UpdateMaintenanceRecord) // 更新维护记录
		equipmentGroup.DELETE("/maintenance/:id", ctrl.DeleteMaintenanceRecord) // 删除维护记录

		// 设备统计
		equipmentGroup.GET("/statistics", ctrl.GetEquipmentStatistics)      // 获取设备统计
		equipmentGroup.GET("/upcoming-maintenance", ctrl.GetUpcomingMaintenances) // 获取即将维护的设备
	}
}