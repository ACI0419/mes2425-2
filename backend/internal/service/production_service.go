package service

import (
	"errors"
	"fmt"
	"mes-system/internal/models"
	"time"

	"gorm.io/gorm"
)

// ProductionService 生产管理服务
type ProductionService struct {
	db *gorm.DB
}

// NewProductionService 创建生产管理服务实例
func NewProductionService(db *gorm.DB) *ProductionService {
	return &ProductionService{db: db}
}

// CreateProductionOrderRequest 创建生产工单请求
type CreateProductionOrderRequest struct {
	ProductID uint       `json:"product_id" binding:"required"`
	Quantity  int        `json:"quantity" binding:"required,min=1"`
	Priority  int        `json:"priority" binding:"min=1,max=5"`
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
}

// UpdateProductionOrderRequest 更新生产工单请求
type UpdateProductionOrderRequest struct {
	Quantity  *int       `json:"quantity,omitempty" binding:"omitempty,min=1"`
	Produced  *int       `json:"produced,omitempty" binding:"omitempty,min=0"`
	Status    *string    `json:"status,omitempty"`
	Priority  *int       `json:"priority,omitempty" binding:"omitempty,min=1,max=5"`
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
}

// ProductionOrderListResponse 生产工单列表响应
type ProductionOrderListResponse struct {
	List     []models.ProductionOrder `json:"list"`
	Total    int64                    `json:"total"`
	Page     int                      `json:"page"`
	PageSize int                      `json:"page_size"`
}

// CreateProductionOrder 创建生产工单
func (s *ProductionService) CreateProductionOrder(req *CreateProductionOrderRequest, createdBy uint) (*models.ProductionOrder, error) {
	// 验证产品是否存在
	var product models.Product
	err := s.db.First(&product, req.ProductID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("产品不存在")
		}
		return nil, err
	}

	// 生成工单号
	orderNo := s.generateOrderNo()

	// 设置默认优先级
	priority := req.Priority
	if priority == 0 {
		priority = 3 // 默认中等优先级
	}

	// 创建生产工单
	order := models.ProductionOrder{
		OrderNo:   orderNo,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		Produced:  0,
		Status:    "pending",
		Priority:  priority,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		CreatedBy: createdBy,
	}

	err = s.db.Create(&order).Error
	if err != nil {
		return nil, err
	}

	// 预加载关联数据
	err = s.db.Preload("Product").Preload("Creator").First(&order, order.ID).Error
	if err != nil {
		return nil, err
	}

	return &order, nil
}

// GetProductionOrderByID 根据ID获取生产工单
func (s *ProductionService) GetProductionOrderByID(id uint) (*models.ProductionOrder, error) {
	var order models.ProductionOrder
	err := s.db.Preload("Product").Preload("Creator").First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// GetProductionOrderList 获取生产工单列表
func (s *ProductionService) GetProductionOrderList(page, pageSize int, status, keyword string) (*ProductionOrderListResponse, error) {
	var orders []models.ProductionOrder
	var total int64

	query := s.db.Model(&models.ProductionOrder{})

	// 状态筛选
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 关键词搜索
	if keyword != "" {
		query = query.Where("order_no LIKE ?", "%"+keyword+"%")
	}

	// 获取总数
	query.Count(&total)

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Preload("Product").Preload("Creator").
		Order("priority DESC, created_at DESC").
		Offset(offset).Limit(pageSize).Find(&orders).Error
	if err != nil {
		return nil, err
	}

	return &ProductionOrderListResponse{
		List:     orders,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// UpdateProductionOrder 更新生产工单
func (s *ProductionService) UpdateProductionOrder(id uint, req *UpdateProductionOrderRequest) (*models.ProductionOrder, error) {
	var order models.ProductionOrder
	err := s.db.First(&order, id).Error
	if err != nil {
		return nil, err
	}

	// 检查工单状态，已完成或已取消的工单不能修改
	if order.Status == "completed" || order.Status == "cancelled" {
		return nil, errors.New("已完成或已取消的工单不能修改")
	}

	// 更新字段
	updateData := make(map[string]interface{})

	if req.Quantity != nil {
		updateData["quantity"] = *req.Quantity
	}

	if req.Produced != nil {
		// 验证已生产数量不能超过计划数量
		quantity := order.Quantity
		if req.Quantity != nil {
			quantity = *req.Quantity
		}
		if *req.Produced > quantity {
			return nil, errors.New("已生产数量不能超过计划数量")
		}
		updateData["produced"] = *req.Produced

		// 自动更新状态
		if *req.Produced == 0 {
			updateData["status"] = "pending"
		} else if *req.Produced < quantity {
			updateData["status"] = "processing"
		} else {
			updateData["status"] = "completed"
		}
	}

	if req.Status != nil {
		// 验证状态转换的合法性
		if !s.isValidStatusTransition(order.Status, *req.Status) {
			return nil, fmt.Errorf("不能从状态 %s 转换到 %s", order.Status, *req.Status)
		}
		updateData["status"] = *req.Status
	}

	if req.Priority != nil {
		updateData["priority"] = *req.Priority
	}

	if req.StartDate != nil {
		updateData["start_date"] = *req.StartDate
	}

	if req.EndDate != nil {
		updateData["end_date"] = *req.EndDate
	}

	// 执行更新
	err = s.db.Model(&order).Updates(updateData).Error
	if err != nil {
		return nil, err
	}

	// 重新查询更新后的数据
	err = s.db.Preload("Product").Preload("Creator").First(&order, id).Error
	if err != nil {
		return nil, err
	}

	return &order, nil
}

// DeleteProductionOrder 删除生产工单
func (s *ProductionService) DeleteProductionOrder(id uint) error {
	var order models.ProductionOrder
	err := s.db.First(&order, id).Error
	if err != nil {
		return err
	}

	// 检查工单状态，进行中的工单不能删除
	if order.Status == "processing" {
		return errors.New("进行中的工单不能删除")
	}

	return s.db.Delete(&order).Error
}

// GetProductionStatistics 获取生产统计数据
func (s *ProductionService) GetProductionStatistics() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 工单状态统计
	var statusStats []struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}
	err := s.db.Model(&models.ProductionOrder{}).
		Select("status, COUNT(*) as count").
		Group("status").Find(&statusStats).Error
	if err != nil {
		return nil, err
	}
	stats["status_stats"] = statusStats

	// 今日生产统计
	today := time.Now().Format("2006-01-02")
	var todayStats struct {
		TotalOrders    int64 `json:"total_orders"`
		CompletedOrders int64 `json:"completed_orders"`
		TotalProduced  int64 `json:"total_produced"`
	}

	s.db.Model(&models.ProductionOrder{}).
		Where("DATE(created_at) = ?", today).
		Count(&todayStats.TotalOrders)

	s.db.Model(&models.ProductionOrder{}).
		Where("DATE(updated_at) = ? AND status = ?", today, "completed").
		Count(&todayStats.CompletedOrders)

	s.db.Model(&models.ProductionOrder{}).
		Where("DATE(updated_at) = ?", today).
		Select("COALESCE(SUM(produced), 0)").
		Scan(&todayStats.TotalProduced)

	stats["today_stats"] = todayStats

	// 本月生产趋势
	var monthlyTrend []struct {
		Date      string `json:"date"`
		Produced  int64  `json:"produced"`
		Completed int64  `json:"completed"`
	}

	startOfMonth := time.Now().AddDate(0, 0, -time.Now().Day()+1).Format("2006-01-02")
	err = s.db.Model(&models.ProductionOrder{}).
		Select("DATE(updated_at) as date, SUM(produced) as produced, COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed").
		Where("updated_at >= ?", startOfMonth).
		Group("DATE(updated_at)").
		Order("date").
		Find(&monthlyTrend).Error
	if err != nil {
		return nil, err
	}
	stats["monthly_trend"] = monthlyTrend

	return stats, nil
}

// generateOrderNo 生成工单号
func (s *ProductionService) generateOrderNo() string {
	now := time.Now()
	prefix := fmt.Sprintf("PO%s", now.Format("20060102"))

	// 查询当天最大序号
	var count int64
	s.db.Model(&models.ProductionOrder{}).
		Where("order_no LIKE ?", prefix+"%").
		Count(&count)

	return fmt.Sprintf("%s%04d", prefix, count+1)
}

// isValidStatusTransition 验证状态转换是否合法
func (s *ProductionService) isValidStatusTransition(from, to string) bool {
	validTransitions := map[string][]string{
		"pending":    {"processing", "cancelled"},
		"processing": {"completed", "cancelled"},
		"completed":  {}, // 已完成状态不能转换
		"cancelled":  {"pending"}, // 已取消可以重新开始
	}

	allowedStates, exists := validTransitions[from]
	if !exists {
		return false
	}

	for _, state := range allowedStates {
		if state == to {
			return true
		}
	}

	return false
}