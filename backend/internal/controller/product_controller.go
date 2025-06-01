package controller

import (
	"mes-system/internal/service"
	"mes-system/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ProductController 产品管理控制器
type ProductController struct {
	productService *service.ProductService
}

// NewProductController 创建产品管理控制器实例
func NewProductController(productService *service.ProductService) *ProductController {
	return &ProductController{
		productService: productService,
	}
}

// CreateProduct 创建产品
func (ctrl *ProductController) CreateProduct(c *gin.Context) {
	var req service.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	product, err := ctrl.productService.CreateProduct(&req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "产品创建成功", product)
}

// GetProduct 获取产品详情
func (ctrl *ProductController) GetProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的产品ID")
		return
	}

	product, err := ctrl.productService.GetProductByID(uint(id))
	if err != nil {
		response.NotFound(c, "产品不存在")
		return
	}

	response.Success(c, product)
}

// GetProductList 获取产品列表
func (ctrl *ProductController) GetProductList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	keyword := c.Query("keyword")

	var status *int
	if statusStr := c.Query("status"); statusStr != "" {
		if s, err := strconv.Atoi(statusStr); err == nil {
			status = &s
		}
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	resp, err := ctrl.productService.GetProductList(page, pageSize, keyword, status)
	if err != nil {
		response.BadRequest(c, "获取产品列表失败")
		return
	}

	response.Success(c, resp)
}

// UpdateProduct 更新产品
func (ctrl *ProductController) UpdateProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的产品ID")
		return
	}

	var req service.UpdateProductRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	product, err := ctrl.productService.UpdateProduct(uint(id), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "产品更新成功", product)
}

// DeleteProduct 删除产品
func (ctrl *ProductController) DeleteProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的产品ID")
		return
	}

	err = ctrl.productService.DeleteProduct(uint(id))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "产品删除成功", nil)
}

// GetAllProducts 获取所有产品（用于下拉选择）
func (ctrl *ProductController) GetAllProducts(c *gin.Context) {
	products, err := ctrl.productService.GetAllProducts()
	if err != nil {
		response.BadRequest(c, "获取产品列表失败")
		return
	}

	response.Success(c, products)
}
