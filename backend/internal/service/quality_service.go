package service

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"mes-system/internal/models"
)

// QualityStandardRequest 质量标准请求结构体
type QualityStandardRequest struct {
	ProductID   uint    `json:"product_id" binding:"required"`   // 产品ID
	Name        string  `json:"name" binding:"required"`         // 标准名称
	Type        string  `json:"type" binding:"required"`         // 检测类型
	MinValue    float64 `json:"min_value"`                       // 最小值
	MaxValue    float64 `json:"max_value"`                       // 最大值
	TargetValue float64 `json:"target_value"`                    // 目标值
	Unit        string  `json:"unit"`                            // 单位
	Description string  `json:"description"`                     // 描述
	IsActive    bool    `json:"is_active"`                       // 是否启用
}

// QualityStandardResponse 质量标准响应结构体
type QualityStandardResponse struct {
	ID          uint      `json:"id"`
	ProductID   uint      `json:"product_id"`
	ProductCode string    `json:"product_code"`
	ProductName string    `json:"product_name"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	MinValue    float64   `json:"min_value"`
	MaxValue    float64   `json:"max_value"`
	TargetValue float64   `json:"target_value"`
	Unit        string    `json:"unit"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// QualityInspectionRequest 质量检测请求结构体
type QualityInspectionRequest struct {
	ProductionOrderID   uint    `json:"production_order_id" binding:"required"` // 生产工单ID
	QualityStandardID   uint    `json:"quality_standard_id" binding:"required"` // 质量标准ID
	InspectorID         uint    `json:"inspector_id" binding:"required"`        // 检测员ID
	ActualValue         float64 `json:"actual_value" binding:"required"`        // 实际值
	Result              string  `json:"result" binding:"required"`              // 检测结果：pass/fail
	Remark              string  `json:"remark"`                                  // 备注
	InspectionTime      *time.Time `json:"inspection_time"`                     // 检测时间
}

// QualityInspectionResponse 质量检测响应结构体
type QualityInspectionResponse struct {
	ID                  uint      `json:"id"`
	ProductionOrderID   uint      `json:"production_order_id"`
	ProductionOrderNo   string    `json:"production_order_no"`
	QualityStandardID   uint      `json:"quality_standard_id"`
	QualityStandardName string    `json:"quality_standard_name"`
	InspectorID         uint      `json:"inspector_id"`
	InspectorName       string    `json:"inspector_name"`
	ActualValue         float64   `json:"actual_value"`
	TargetValue         float64   `json:"target_value"`
	MinValue            float64   `json:"min_value"`
	MaxValue            float64   `json:"max_value"`
	Unit                string    `json:"unit"`
	Result              string    `json:"result"`
	Remark              string    `json:"remark"`
	InspectionTime      time.Time `json:"inspection_time"`
	CreatedAt           time.Time `json:"created_at"`
}

// QualityStatistics 质量统计结构体
type QualityStatistics struct {
	TotalInspections int     `json:"total_inspections"`
	PassedCount      int     `json:"passed_count"`
	FailedCount      int     `json:"failed_count"`
	PassRate         float64 `json:"pass_rate"`
	FailRate         float64 `json:"fail_rate"`
}

// QualityService 质量服务
type QualityService struct {
	db *gorm.DB
}

// NewQualityService 创建质量服务实例
func NewQualityService(db *gorm.DB) *QualityService {
	return &QualityService{db: db}
}

// CreateQualityStandard 创建质量标准
func (s *QualityService) CreateQualityStandard(req *QualityStandardRequest) (*QualityStandardResponse, error) {
	// 验证产品是否存在
	var product models.Product
	if err := s.db.First(&product, req.ProductID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("产品不存在")
		}
		return nil, fmt.Errorf("验证产品失败: %v", err)
	}

	// 验证数值范围
	if req.MinValue >= req.MaxValue {
		return nil, errors.New("最小值必须小于最大值")
	}

	if req.TargetValue < req.MinValue || req.TargetValue > req.MaxValue {
		return nil, errors.New("目标值必须在最小值和最大值之间")
	}

	// 检查同一产品下是否已存在相同名称的质量标准
	if s.isQualityStandardNameExists(req.ProductID, req.Name, 0) {
		return nil, errors.New("该产品下已存在相同名称的质量标准")
	}

	qualityStandard := &models.QualityStandard{
		ProductID:   req.ProductID,
		Name:        req.Name,
		Type:        req.Type,
		MinValue:    req.MinValue,
		MaxValue:    req.MaxValue,
		TargetValue: req.TargetValue,
		Unit:        req.Unit,
		Description: req.Description,
		IsActive:    req.IsActive,
	}

	if err := s.db.Create(qualityStandard).Error; err != nil {
		return nil, fmt.Errorf("创建质量标准失败: %v", err)
	}

	return s.qualityStandardToResponse(qualityStandard, &product), nil
}

// GetQualityStandard 获取质量标准详情
func (s *QualityService) GetQualityStandard(id uint) (*QualityStandardResponse, error) {
	var qualityStandard models.QualityStandard
	if err := s.db.Preload("Product").First(&qualityStandard, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("质量标准不存在")
		}
		return nil, fmt.Errorf("获取质量标准失败: %v", err)
	}

	return s.qualityStandardToResponse(&qualityStandard, &qualityStandard.Product), nil
}

// GetQualityStandardList 获取质量标准列表
func (s *QualityService) GetQualityStandardList(page, pageSize int, productID uint, standardType string, isActive *bool) ([]QualityStandardResponse, int64, error) {
	var qualityStandards []models.QualityStandard
	var total int64

	query := s.db.Model(&models.QualityStandard{}).Preload("Product")

	// 按产品筛选
	if productID > 0 {
		query = query.Where("product_id = ?", productID)
	}

	// 按类型筛选
	if standardType != "" {
		query = query.Where("type = ?", standardType)
	}

	// 按状态筛选
	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("获取质量标准总数失败: %v", err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&qualityStandards).Error; err != nil {
		return nil, 0, fmt.Errorf("获取质量标准列表失败: %v", err)
	}

	var responses []QualityStandardResponse
	for _, standard := range qualityStandards {
		responses = append(responses, *s.qualityStandardToResponse(&standard, &standard.Product))
	}

	return responses, total, nil
}

// UpdateQualityStandard 更新质量标准
func (s *QualityService) UpdateQualityStandard(id uint, req *QualityStandardRequest) (*QualityStandardResponse, error) {
	var qualityStandard models.QualityStandard
	if err := s.db.Preload("Product").First(&qualityStandard, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("质量标准不存在")
		}
		return nil, fmt.Errorf("获取质量标准失败: %v", err)
	}

	// 验证产品是否存在
	var product models.Product
	if err := s.db.First(&product, req.ProductID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("产品不存在")
		}
		return nil, fmt.Errorf("验证产品失败: %v", err)
	}

	// 验证数值范围
	if req.MinValue >= req.MaxValue {
		return nil, errors.New("最小值必须小于最大值")
	}

	if req.TargetValue < req.MinValue || req.TargetValue > req.MaxValue {
		return nil, errors.New("目标值必须在最小值和最大值之间")
	}

	// 检查同一产品下是否已存在相同名称的质量标准（排除当前标准）
	if s.isQualityStandardNameExists(req.ProductID, req.Name, id) {
		return nil, errors.New("该产品下已存在相同名称的质量标准")
	}

	// 更新质量标准信息
	qualityStandard.ProductID = req.ProductID
	qualityStandard.Name = req.Name
	qualityStandard.Type = req.Type
	qualityStandard.MinValue = req.MinValue
	qualityStandard.MaxValue = req.MaxValue
	qualityStandard.TargetValue = req.TargetValue
	qualityStandard.Unit = req.Unit
	qualityStandard.Description = req.Description
	qualityStandard.IsActive = req.IsActive

	if err := s.db.Save(&qualityStandard).Error; err != nil {
		return nil, fmt.Errorf("更新质量标准失败: %v", err)
	}

	return s.qualityStandardToResponse(&qualityStandard, &product), nil
}

// DeleteQualityStandard 删除质量标准
func (s *QualityService) DeleteQualityStandard(id uint) error {
	var qualityStandard models.QualityStandard
	if err := s.db.First(&qualityStandard, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("质量标准不存在")
		}
		return fmt.Errorf("获取质量标准失败: %v", err)
	}

	// 检查是否有相关的检测记录
	var count int64
	if err := s.db.Model(&models.QualityInspection{}).Where("quality_standard_id = ?", id).Count(&count).Error; err != nil {
		return fmt.Errorf("检查质量检测记录失败: %v", err)
	}

	if count > 0 {
		return errors.New("该质量标准存在检测记录，无法删除")
	}

	if err := s.db.Delete(&qualityStandard).Error; err != nil {
		return fmt.Errorf("删除质量标准失败: %v", err)
	}

	return nil
}

// CreateQualityInspection 创建质量检测记录
func (s *QualityService) CreateQualityInspection(req *QualityInspectionRequest) (*QualityInspectionResponse, error) {
	// 验证生产工单是否存在
	var productionOrder models.ProductionOrder
	if err := s.db.First(&productionOrder, req.ProductionOrderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("生产工单不存在")
		}
		return nil, fmt.Errorf("验证生产工单失败: %v", err)
	}

	// 验证质量标准是否存在
	var qualityStandard models.QualityStandard
	if err := s.db.First(&qualityStandard, req.QualityStandardID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("质量标准不存在")
		}
		return nil, fmt.Errorf("验证质量标准失败: %v", err)
	}

	// 验证检测员是否存在
	var inspector models.User
	if err := s.db.First(&inspector, req.InspectorID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("检测员不存在")
		}
		return nil, fmt.Errorf("验证检测员失败: %v", err)
	}

	// 验证检测结果
	if req.Result != "pass" && req.Result != "fail" {
		return nil, errors.New("检测结果必须是 pass 或 fail")
	}

	// 设置检测时间
	inspectionTime := time.Now()
	if req.InspectionTime != nil {
		inspectionTime = *req.InspectionTime
	}

	qualityInspection := &models.QualityInspection{
		ProductionOrderID: req.ProductionOrderID,
		QualityStandardID: req.QualityStandardID,
		InspectorID:       req.InspectorID,
		ActualValue:       req.ActualValue,
		Result:            req.Result,
		Remark:            req.Remark,
		InspectionTime:    inspectionTime,
	}

	if err := s.db.Create(qualityInspection).Error; err != nil {
		return nil, fmt.Errorf("创建质量检测记录失败: %v", err)
	}

	return s.qualityInspectionToResponse(qualityInspection, &productionOrder, &qualityStandard, &inspector), nil
}

// GetQualityInspection 获取质量检测记录详情
func (s *QualityService) GetQualityInspection(id uint) (*QualityInspectionResponse, error) {
	var qualityInspection models.QualityInspection
	if err := s.db.Preload("ProductionOrder").Preload("QualityStandard").Preload("Inspector").First(&qualityInspection, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("质量检测记录不存在")
		}
		return nil, fmt.Errorf("获取质量检测记录失败: %v", err)
	}

	return s.qualityInspectionToResponse(&qualityInspection, &qualityInspection.ProductionOrder, &qualityInspection.QualityStandard, &qualityInspection.Inspector), nil
}

// GetQualityInspectionList 获取质量检测记录列表
func (s *QualityService) GetQualityInspectionList(page, pageSize int, productionOrderID, qualityStandardID, inspectorID uint, result string) ([]QualityInspectionResponse, int64, error) {
	var qualityInspections []models.QualityInspection
	var total int64

	query := s.db.Model(&models.QualityInspection{}).Preload("ProductionOrder").Preload("QualityStandard").Preload("Inspector")

	// 按生产工单筛选
	if productionOrderID > 0 {
		query = query.Where("production_order_id = ?", productionOrderID)
	}

	// 按质量标准筛选
	if qualityStandardID > 0 {
		query = query.Where("quality_standard_id = ?", qualityStandardID)
	}

	// 按检测员筛选
	if inspectorID > 0 {
		query = query.Where("inspector_id = ?", inspectorID)
	}

	// 按检测结果筛选
	if result != "" {
		query = query.Where("result = ?", result)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("获取质量检测记录总数失败: %v", err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("inspection_time DESC").Find(&qualityInspections).Error; err != nil {
		return nil, 0, fmt.Errorf("获取质量检测记录列表失败: %v", err)
	}

	var responses []QualityInspectionResponse
	for _, inspection := range qualityInspections {
		responses = append(responses, *s.qualityInspectionToResponse(&inspection, &inspection.ProductionOrder, &inspection.QualityStandard, &inspection.Inspector))
	}

	return responses, total, nil
}

// GetQualityStatistics 获取质量统计数据
func (s *QualityService) GetQualityStatistics(startDate, endDate *time.Time, productionOrderID, qualityStandardID uint) (*QualityStatistics, error) {
	query := s.db.Model(&models.QualityInspection{})

	// 时间范围筛选
	if startDate != nil {
		query = query.Where("inspection_time >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("inspection_time <= ?", *endDate)
	}

	// 按生产工单筛选
	if productionOrderID > 0 {
		query = query.Where("production_order_id = ?", productionOrderID)
	}

	// 按质量标准筛选
	if qualityStandardID > 0 {
		query = query.Where("quality_standard_id = ?", qualityStandardID)
	}

	// 获取总检测次数
	var totalInspections int64
	if err := query.Count(&totalInspections).Error; err != nil {
		return nil, fmt.Errorf("获取总检测次数失败: %v", err)
	}

	// 获取通过次数
	var passedCount int64
	if err := query.Where("result = ?", "pass").Count(&passedCount).Error; err != nil {
		return nil, fmt.Errorf("获取通过次数失败: %v", err)
	}

	failedCount := totalInspections - passedCount

	var passRate, failRate float64
	if totalInspections > 0 {
		passRate = float64(passedCount) / float64(totalInspections) * 100
		failRate = float64(failedCount) / float64(totalInspections) * 100
	}

	return &QualityStatistics{
		TotalInspections: int(totalInspections),
		PassedCount:      int(passedCount),
		FailedCount:      int(failedCount),
		PassRate:         passRate,
		FailRate:         failRate,
	}, nil
}

// GetQualityStandardTypes 获取所有质量标准类型
func (s *QualityService) GetQualityStandardTypes() ([]string, error) {
	var types []string
	if err := s.db.Model(&models.QualityStandard{}).Distinct("type").Pluck("type", &types).Error; err != nil {
		return nil, fmt.Errorf("获取质量标准类型失败: %v", err)
	}
	return types, nil
}

// 辅助函数：检查质量标准名称是否存在
func (s *QualityService) isQualityStandardNameExists(productID uint, name string, excludeID uint) bool {
	var count int64
	query := s.db.Model(&models.QualityStandard{}).Where("product_id = ? AND name = ?", productID, name)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	query.Count(&count)
	return count > 0
}

// 辅助函数：将质量标准模型转换为响应结构体
func (s *QualityService) qualityStandardToResponse(standard *models.QualityStandard, product *models.Product) *QualityStandardResponse {
	return &QualityStandardResponse{
		ID:          standard.ID,
		ProductID:   standard.ProductID,
		ProductCode: product.Code,
		ProductName: product.Name,
		Name:        standard.Name,
		Type:        standard.Type,
		MinValue:    standard.MinValue,
		MaxValue:    standard.MaxValue,
		TargetValue: standard.TargetValue,
		Unit:        standard.Unit,
		Description: standard.Description,
		IsActive:    standard.IsActive,
		CreatedAt:   standard.CreatedAt,
		UpdatedAt:   standard.UpdatedAt,
	}
}

// 辅助函数：将质量检测模型转换为响应结构体
func (s *QualityService) qualityInspectionToResponse(inspection *models.QualityInspection, productionOrder *models.ProductionOrder, qualityStandard *models.QualityStandard, inspector *models.User) *QualityInspectionResponse {
	return &QualityInspectionResponse{
		ID:                  inspection.ID,
		ProductionOrderID:   inspection.ProductionOrderID,
		ProductionOrderNo:   productionOrder.OrderNo,
		QualityStandardID:   inspection.QualityStandardID,
		QualityStandardName: qualityStandard.Name,
		InspectorID:         inspection.InspectorID,
		InspectorName:       inspector.Username,
		ActualValue:         inspection.ActualValue,
		TargetValue:         qualityStandard.TargetValue,
		MinValue:            qualityStandard.MinValue,
		MaxValue:            qualityStandard.MaxValue,
		Unit:                qualityStandard.Unit,
		Result:              inspection.Result,
		Remark:              inspection.Remark,
		InspectionTime:      inspection.InspectionTime,
		CreatedAt:           inspection.CreatedAt,
	}
}

// UpdateQualityInspection 更新质量检测记录
func (s *QualityService) UpdateQualityInspection(id uint, req *QualityInspectionRequest) (*QualityInspectionResponse, error) {
	var qualityInspection models.QualityInspection
	if err := s.db.Preload("ProductionOrder").Preload("QualityStandard").Preload("Inspector").First(&qualityInspection, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("质量检测记录不存在")
		}
		return nil, fmt.Errorf("获取质量检测记录失败: %v", err)
	}

	// 验证生产工单是否存在
	var productionOrder models.ProductionOrder
	if err := s.db.First(&productionOrder, req.ProductionOrderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("生产工单不存在")
		}
		return nil, fmt.Errorf("验证生产工单失败: %v", err)
	}

	// 验证质量标准是否存在
	var qualityStandard models.QualityStandard
	if err := s.db.First(&qualityStandard, req.QualityStandardID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("质量标准不存在")
		}
		return nil, fmt.Errorf("验证质量标准失败: %v", err)
	}

	// 验证检测员是否存在
	var inspector models.User
	if err := s.db.First(&inspector, req.InspectorID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("检测员不存在")
		}
		return nil, fmt.Errorf("验证检测员失败: %v", err)
	}

	// 验证检测结果
	if req.Result != "pass" && req.Result != "fail" {
		return nil, errors.New("检测结果必须是 pass 或 fail")
	}

	// 更新质量检测记录信息
	qualityInspection.ProductionOrderID = req.ProductionOrderID
	qualityInspection.QualityStandardID = req.QualityStandardID
	qualityInspection.InspectorID = req.InspectorID
	qualityInspection.ActualValue = req.ActualValue
	qualityInspection.Result = req.Result
	qualityInspection.Remark = req.Remark

	// 更新检测时间（如果提供）
	if req.InspectionTime != nil {
		qualityInspection.InspectionTime = *req.InspectionTime
	}

	if err := s.db.Save(&qualityInspection).Error; err != nil {
		return nil, fmt.Errorf("更新质量检测记录失败: %v", err)
	}

	return s.qualityInspectionToResponse(&qualityInspection, &productionOrder, &qualityStandard, &inspector), nil
}

// DeleteQualityInspection 删除质量检测记录
func (s *QualityService) DeleteQualityInspection(id uint) error {
	var qualityInspection models.QualityInspection
	if err := s.db.First(&qualityInspection, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("质量检测记录不存在")
		}
		return fmt.Errorf("获取质量检测记录失败: %v", err)
	}

	if err := s.db.Delete(&qualityInspection).Error; err != nil {
		return fmt.Errorf("删除质量检测记录失败: %v", err)
	}

	return nil
}

// ... existing code ...