package service

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"mes-system/internal/models"
)

// MaterialRequest 物料请求结构体
type MaterialRequest struct {
	Code        string  `json:"code" binding:"required"`         // 物料编码
	Name        string  `json:"name" binding:"required"`         // 物料名称
	Type        string  `json:"type" binding:"required"`         // 物料类型
	Unit        string  `json:"unit" binding:"required"`         // 计量单位
	Price       float64 `json:"price" binding:"min=0"`           // 单价
	MinStock    int     `json:"min_stock" binding:"min=0"`       // 最小库存
	MaxStock    int     `json:"max_stock" binding:"min=0"`       // 最大库存
	Description string  `json:"description"`                     // 描述
}

// MaterialResponse 物料响应结构体
type MaterialResponse struct {
	ID          uint      `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Unit        string    `json:"unit"`
	Price       float64   `json:"price"`
	CurrentStock int      `json:"current_stock"`
	MinStock    int       `json:"min_stock"`
	MaxStock    int       `json:"max_stock"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// MaterialTransactionRequest 物料交易请求结构体
type MaterialTransactionRequest struct {
	MaterialID    uint    `json:"material_id" binding:"required"`    // 物料ID
	Type          string  `json:"type" binding:"required"`           // 交易类型：in/out
	Quantity      int     `json:"quantity" binding:"required,min=1"` // 数量
	Price         float64 `json:"price" binding:"min=0"`             // 单价
	Supplier      string  `json:"supplier"`                          // 供应商
	ProductionOrderID *uint `json:"production_order_id"`             // 生产工单ID（出库时）
	Remark        string  `json:"remark"`                            // 备注
}

// MaterialTransactionResponse 物料交易响应结构体
type MaterialTransactionResponse struct {
	ID                uint      `json:"id"`
	MaterialID        uint      `json:"material_id"`
	MaterialCode      string    `json:"material_code"`
	MaterialName      string    `json:"material_name"`
	Type              string    `json:"type"`
	Quantity          int       `json:"quantity"`
	Price             float64   `json:"price"`
	TotalAmount       float64   `json:"total_amount"`
	Supplier          string    `json:"supplier"`
	ProductionOrderID *uint     `json:"production_order_id"`
	Remark            string    `json:"remark"`
	CreatedAt         time.Time `json:"created_at"`
}

// MaterialService 物料服务
type MaterialService struct {
	db *gorm.DB
}

// NewMaterialService 创建物料服务实例
func NewMaterialService(db *gorm.DB) *MaterialService {
	return &MaterialService{db: db}
}

// CreateMaterial 创建物料
func (s *MaterialService) CreateMaterial(req *MaterialRequest) (*MaterialResponse, error) {
	// 检查物料编码是否已存在
	if s.isMaterialCodeExists(req.Code, 0) {
		return nil, errors.New("物料编码已存在")
	}

	// 验证最大库存必须大于最小库存
	if req.MaxStock <= req.MinStock {
		return nil, errors.New("最大库存必须大于最小库存")
	}

	material := &models.Material{
		Code:        req.Code,
		Name:        req.Name,
		Type:        req.Type,
		Unit:        req.Unit,
		Price:       req.Price,
		CurrentStock: 0, // 初始库存为0
		MinStock:    req.MinStock,
		MaxStock:    req.MaxStock,
		Description: req.Description,
	}

	if err := s.db.Create(material).Error; err != nil {
		return nil, fmt.Errorf("创建物料失败: %v", err)
	}

	return s.materialToResponse(material), nil
}

// GetMaterial 获取物料详情
func (s *MaterialService) GetMaterial(id uint) (*MaterialResponse, error) {
	var material models.Material
	if err := s.db.First(&material, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("物料不存在")
		}
		return nil, fmt.Errorf("获取物料失败: %v", err)
	}

	return s.materialToResponse(&material), nil
}

// GetMaterialList 获取物料列表
func (s *MaterialService) GetMaterialList(page, pageSize int, materialType, keyword string) ([]MaterialResponse, int64, error) {
	var materials []models.Material
	var total int64

	query := s.db.Model(&models.Material{})

	// 按类型筛选
	if materialType != "" {
		query = query.Where("type = ?", materialType)
	}

	// 关键词搜索
	if keyword != "" {
		query = query.Where("code LIKE ? OR name LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("获取物料总数失败: %v", err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&materials).Error; err != nil {
		return nil, 0, fmt.Errorf("获取物料列表失败: %v", err)
	}

	var responses []MaterialResponse
	for _, material := range materials {
		responses = append(responses, *s.materialToResponse(&material))
	}

	return responses, total, nil
}

// UpdateMaterial 更新物料
func (s *MaterialService) UpdateMaterial(id uint, req *MaterialRequest) (*MaterialResponse, error) {
	var material models.Material
	if err := s.db.First(&material, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("物料不存在")
		}
		return nil, fmt.Errorf("获取物料失败: %v", err)
	}

	// 检查物料编码是否已存在（排除当前物料）
	if s.isMaterialCodeExists(req.Code, id) {
		return nil, errors.New("物料编码已存在")
	}

	// 验证最大库存必须大于最小库存
	if req.MaxStock <= req.MinStock {
		return nil, errors.New("最大库存必须大于最小库存")
	}

	// 更新物料信息
	material.Code = req.Code
	material.Name = req.Name
	material.Type = req.Type
	material.Unit = req.Unit
	material.Price = req.Price
	material.MinStock = req.MinStock
	material.MaxStock = req.MaxStock
	material.Description = req.Description

	if err := s.db.Save(&material).Error; err != nil {
		return nil, fmt.Errorf("更新物料失败: %v", err)
	}

	return s.materialToResponse(&material), nil
}

// DeleteMaterial 删除物料
func (s *MaterialService) DeleteMaterial(id uint) error {
	var material models.Material
	if err := s.db.First(&material, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("物料不存在")
		}
		return fmt.Errorf("获取物料失败: %v", err)
	}

	// 检查是否有相关的交易记录
	var count int64
	if err := s.db.Model(&models.MaterialTransaction{}).Where("material_id = ?", id).Count(&count).Error; err != nil {
		return fmt.Errorf("检查物料交易记录失败: %v", err)
	}

	if count > 0 {
		return errors.New("该物料存在交易记录，无法删除")
	}

	if err := s.db.Delete(&material).Error; err != nil {
		return fmt.Errorf("删除物料失败: %v", err)
	}

	return nil
}

// CreateTransaction 创建物料交易（入库/出库）
func (s *MaterialService) CreateTransaction(req *MaterialTransactionRequest) (*MaterialTransactionResponse, error) {
	// 验证交易类型
	if req.Type != "in" && req.Type != "out" {
		return nil, errors.New("交易类型必须是 in 或 out")
	}

	// 获取物料信息
	var material models.Material
	if err := s.db.First(&material, req.MaterialID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("物料不存在")
		}
		return nil, fmt.Errorf("获取物料失败: %v", err)
	}

	// 出库时检查库存是否充足
	if req.Type == "out" && material.CurrentStock < req.Quantity {
		return nil, errors.New("库存不足")
	}

	// 开始事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建交易记录
	transaction := &models.MaterialTransaction{
		MaterialID:        req.MaterialID,
		Type:              req.Type,
		Quantity:          req.Quantity,
		Price:             req.Price,
		TotalAmount:       float64(req.Quantity) * req.Price,
		Supplier:          req.Supplier,
		ProductionOrderID: req.ProductionOrderID,
		Remark:            req.Remark,
	}

	if err := tx.Create(transaction).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("创建交易记录失败: %v", err)
	}

	// 更新库存
	if req.Type == "in" {
		material.CurrentStock += req.Quantity
	} else {
		material.CurrentStock -= req.Quantity
	}

	if err := tx.Save(&material).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("更新库存失败: %v", err)
	}

	tx.Commit()

	return s.transactionToResponse(transaction, &material), nil
}

// GetTransactionList 获取物料交易列表
func (s *MaterialService) GetTransactionList(page, pageSize int, materialID uint, transactionType string) ([]MaterialTransactionResponse, int64, error) {
	var transactions []models.MaterialTransaction
	var total int64

	query := s.db.Model(&models.MaterialTransaction{}).Preload("Material")

	// 按物料筛选
	if materialID > 0 {
		query = query.Where("material_id = ?", materialID)
	}

	// 按交易类型筛选
	if transactionType != "" {
		query = query.Where("type = ?", transactionType)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("获取交易总数失败: %v", err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&transactions).Error; err != nil {
		return nil, 0, fmt.Errorf("获取交易列表失败: %v", err)
	}

	var responses []MaterialTransactionResponse
	for _, transaction := range transactions {
		responses = append(responses, *s.transactionToResponse(&transaction, &transaction.Material))
	}

	return responses, total, nil
}

// GetLowStockMaterials 获取低库存物料列表
func (s *MaterialService) GetLowStockMaterials() ([]MaterialResponse, error) {
	var materials []models.Material
	if err := s.db.Where("current_stock <= min_stock").Find(&materials).Error; err != nil {
		return nil, fmt.Errorf("获取低库存物料失败: %v", err)
	}

	var responses []MaterialResponse
	for _, material := range materials {
		responses = append(responses, *s.materialToResponse(&material))
	}

	return responses, nil
}

// GetMaterialTypes 获取所有物料类型
func (s *MaterialService) GetMaterialTypes() ([]string, error) {
	var types []string
	if err := s.db.Model(&models.Material{}).Distinct("type").Where("type != ''").Pluck("type", &types).Error; err != nil {
		return nil, fmt.Errorf("获取物料类型失败: %v", err)
	}
	return types, nil
}

// 辅助函数：检查物料编码是否存在
func (s *MaterialService) isMaterialCodeExists(code string, excludeID uint) bool {
	var count int64
	query := s.db.Model(&models.Material{}).Where("code = ?", code)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	query.Count(&count)
	return count > 0
}

// 辅助函数：将物料模型转换为响应结构体
func (s *MaterialService) materialToResponse(material *models.Material) *MaterialResponse {
	return &MaterialResponse{
		ID:           material.ID,
		Code:         material.Code,
		Name:         material.Name,
		Type:         material.Type,        // 确保字段名一致
		Unit:         material.Unit,
		Price:        material.Price,
		CurrentStock: material.CurrentStock,
		MinStock:     material.MinStock,
		MaxStock:     material.MaxStock,
		Description:  material.Description, // 确保字段名一致
		CreatedAt:    material.CreatedAt,
		UpdatedAt:    material.UpdatedAt,
	}
}

// 辅助函数：将交易模型转换为响应结构体
func (s *MaterialService) transactionToResponse(transaction *models.MaterialTransaction, material *models.Material) *MaterialTransactionResponse {
	return &MaterialTransactionResponse{
		ID:                transaction.ID,
		MaterialID:        transaction.MaterialID,
		MaterialCode:      material.Code,
		MaterialName:      material.Name,
		Type:              transaction.Type,
		Quantity:          transaction.Quantity,
		Price:             transaction.Price,             // 现在模型中有这个字段
		TotalAmount:       transaction.TotalAmount,       // 现在模型中有这个字段
		Supplier:          transaction.Supplier,          // 现在模型中有这个字段
		ProductionOrderID: transaction.ProductionOrderID, // 现在模型中有这个字段
		Remark:            transaction.Remark,            // 现在模型中有这个字段
		CreatedAt:         transaction.CreatedAt,
	}
}