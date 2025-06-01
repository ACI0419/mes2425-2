package controller

import (
	"mes-system/internal/service"
	"mes-system/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ProductionController 生产管理控制器
type ProductionController struct {
	productionService *service.ProductionService
	productService    *service.ProductService
}

// NewProductionController 创建生产管理控制器实例
func NewProductionController(productionService *service.ProductionService, productService *service.ProductService) *ProductionController {
	return &ProductionController{
		productionService: productionService,
		productService:    productService,
	}
}

// CreateProductionOrder 创建生产工单
func (ctrl *ProductionController) CreateProductionOrder(c *gin.Context) {
	var req service.CreateProductionOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "用户未登录")
		return
	}

	order, err := ctrl.productionService.CreateProductionOrder(&req, userID.(uint))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "生产工单创建成功", order)
}

// GetProductionOrder 获取生产工单详情
func (ctrl *ProductionController) GetProductionOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的工单ID")
		return
	}

	order, err := ctrl.productionService.GetProductionOrderByID(uint(id))
	if err != nil {
		response.NotFound(c, "生产工单不存在")
		return
	}

	response.Success(c, order)
}

// GetProductionOrderList 获取生产工单列表
func (ctrl *ProductionController) GetProductionOrderList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")
	keyword := c.Query("keyword")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	resp, err := ctrl.productionService.GetProductionOrderList(page, pageSize, status, keyword)
	if err != nil {
		response.BadRequest(c, "获取生产工单列表失败")
		return
	}

	response.Success(c, resp)
}

// UpdateProductionOrder 更新生产工单
func (ctrl *ProductionController) UpdateProductionOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的工单ID")
		return
	}

	var req service.UpdateProductionOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	order, err := ctrl.productionService.UpdateProductionOrder(uint(id), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "生产工单更新成功", order)
}

// DeleteProductionOrder 删除生产工单
func (ctrl *ProductionController) DeleteProductionOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的工单ID")
		return
	}

	err = ctrl.productionService.DeleteProductionOrder(uint(id))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "生产工单删除成功", nil)
}

// GetProductionStatistics 获取生产统计数据
func (ctrl *ProductionController) GetProductionStatistics(c *gin.Context) {
	stats, err := ctrl.productionService.GetProductionStatistics()
	if err != nil {
		response.BadRequest(c, "获取生产统计数据失败")
		return
	}

	response.Success(c, stats)
}