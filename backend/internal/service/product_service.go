package service

import (
	"errors"
	"mes-system/internal/models"

	"gorm.io/gorm"
)

// ProductService 产品管理服务
type ProductService struct {
	db *gorm.DB
}

// NewProductService 创建产品管理服务实例
func NewProductService(db *gorm.DB) *ProductService {
	return &ProductService{db: db}
}

// CreateProductRequest 创建产品请求
type CreateProductRequest struct {
	Code        string  `json:"code" binding:"required,max=50"`
	Name        string  `json:"name" binding:"required,max=100"`
	Description string  `json:"description"`
	Unit        string  `json:"unit" binding:"required,max=20"`
	Price       float64 `json:"price" binding:"min=0"`
}

// UpdateProductRequest 更新产品请求
type UpdateProductRequest struct {
	Name        *string  `json:"name,omitempty" binding:"omitempty,max=100"`
	Description *string  `json:"description,omitempty"`
	Unit        *string  `json:"unit,omitempty" binding:"omitempty,max=20"`
	Price       *float64 `json:"price,omitempty" binding:"omitempty,min=0"`
	Status      *int     `json:"status,omitempty" binding:"omitempty,oneof=0 1"`
}

// ProductListResponse 产品列表响应
type ProductListResponse struct {
	List     []models.Product `json:"list"`
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	PageSize int              `json:"page_size"`
}

// CreateProduct 创建产品
func (s *ProductService) CreateProduct(req *CreateProductRequest) (*models.Product, error) {
	// 检查产品编码是否已存在
	var count int64
	s.db.Model(&models.Product{}).Where("code = ?", req.Code).Count(&count)
	if count > 0 {
		return nil, errors.New("产品编码已存在")
	}

	// 创建产品
	product := models.Product{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		Unit:        req.Unit,
		Price:       req.Price,
		Status:      1,
	}

	err := s.db.Create(&product).Error
	if err != nil {
		return nil, err
	}

	return &product, nil
}

// GetProductByID 根据ID获取产品
func (s *ProductService) GetProductByID(id uint) (*models.Product, error) {
	var product models.Product
	err := s.db.First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// GetProductList 获取产品列表
func (s *ProductService) GetProductList(page, pageSize int, keyword string, status *int) (*ProductListResponse, error) {
	var products []models.Product
	var total int64

	query := s.db.Model(&models.Product{})

	// 状态筛选
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	// 关键词搜索
	if keyword != "" {
		query = query.Where("code LIKE ? OR name LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 获取总数
	query.Count(&total)

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&products).Error
	if err != nil {
		return nil, err
	}

	return &ProductListResponse{
		List:     products,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// UpdateProduct 更新产品
func (s *ProductService) UpdateProduct(id uint, req *UpdateProductRequest) (*models.Product, error) {
	var product models.Product
	err := s.db.First(&product, id).Error
	if err != nil {
		return nil, err
	}

	// 更新字段
	updateData := make(map[string]interface{})

	if req.Name != nil {
		updateData["name"] = *req.Name
	}

	if req.Description != nil {
		updateData["description"] = *req.Description
	}

	if req.Unit != nil {
		updateData["unit"] = *req.Unit
	}

	if req.Price != nil {
		updateData["price"] = *req.Price
	}

	if req.Status != nil {
		updateData["status"] = *req.Status
	}

	// 执行更新
	err = s.db.Model(&product).Updates(updateData).Error
	if err != nil {
		return nil, err
	}

	// 重新查询更新后的数据
	err = s.db.First(&product, id).Error
	if err != nil {
		return nil, err
	}

	return &product, nil
}

// DeleteProduct 删除产品
func (s *ProductService) DeleteProduct(id uint) error {
	// 检查是否有关联的生产工单
	var count int64
	s.db.Model(&models.ProductionOrder{}).Where("product_id = ?", id).Count(&count)
	if count > 0 {
		return errors.New("该产品存在关联的生产工单，不能删除")
	}

	return s.db.Delete(&models.Product{}, id).Error
}

// GetAllProducts 获取所有启用的产品（用于下拉选择）
func (s *ProductService) GetAllProducts() ([]models.Product, error) {
	var products []models.Product
	err := s.db.Where("status = ?", 1).Order("name").Find(&products).Error
	return products, err
}