package service

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"mes-system/internal/models"
)

// EquipmentRequest 设备请求结构体
type EquipmentRequest struct {
	Code         string    `json:"code" binding:"required"`         // 设备编码
	Name         string    `json:"name" binding:"required"`         // 设备名称
	Type         string    `json:"type" binding:"required"`         // 设备类型
	Model        string    `json:"model"`                           // 设备型号
	Manufacturer string    `json:"manufacturer"`                    // 制造商
	PurchaseDate *time.Time `json:"purchase_date"`                 // 采购日期
	WarrantyDate *time.Time `json:"warranty_date"`                 // 保修期至
	Location     string    `json:"location"`                        // 设备位置
	Status       string    `json:"status" binding:"required"`       // 设备状态：running/stopped/maintenance/fault
	Description  string    `json:"description"`                     // 描述
}

// EquipmentResponse 设备响应结构体
type EquipmentResponse struct {
	ID           uint       `json:"id"`
	Code         string     `json:"code"`
	Name         string     `json:"name"`
	Type         string     `json:"type"`
	Model        string     `json:"model"`
	Manufacturer string     `json:"manufacturer"`
	PurchaseDate *time.Time `json:"purchase_date"`
	WarrantyDate *time.Time `json:"warranty_date"`
	Location     string     `json:"location"`
	Status       string     `json:"status"`
	Description  string     `json:"description"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// MaintenanceRecordRequest 维护记录请求结构体
type MaintenanceRecordRequest struct {
	EquipmentID     uint       `json:"equipment_id" binding:"required"`     // 设备ID
	MaintainerID    uint       `json:"maintainer_id" binding:"required"`    // 维护人员ID
	Type            string     `json:"type" binding:"required"`             // 维护类型：preventive/corrective/emergency
	Description     string     `json:"description" binding:"required"`      // 维护描述
	StartTime       time.Time  `json:"start_time" binding:"required"`       // 开始时间
	EndTime         *time.Time `json:"end_time"`                            // 结束时间
	Cost            float64    `json:"cost" binding:"min=0"`                // 维护费用
	PartsReplaced   string     `json:"parts_replaced"`                      // 更换部件
	Result          string     `json:"result"`                              // 维护结果
	NextMaintenance *time.Time `json:"next_maintenance"`                    // 下次维护时间
	Remark          string     `json:"remark"`                              // 备注
}

// MaintenanceRecordResponse 维护记录响应结构体
type MaintenanceRecordResponse struct {
	ID              uint       `json:"id"`
	EquipmentID     uint       `json:"equipment_id"`
	EquipmentCode   string     `json:"equipment_code"`
	EquipmentName   string     `json:"equipment_name"`
	MaintainerID    uint       `json:"maintainer_id"`
	MaintainerName  string     `json:"maintainer_name"`
	Type            string     `json:"type"`
	Description     string     `json:"description"`
	StartTime       time.Time  `json:"start_time"`
	EndTime         *time.Time `json:"end_time"`
	Duration        *int       `json:"duration"` // 维护时长（分钟）
	Cost            float64    `json:"cost"`
	PartsReplaced   string     `json:"parts_replaced"`
	Result          string     `json:"result"`
	NextMaintenance *time.Time `json:"next_maintenance"`
	Remark          string     `json:"remark"`
	CreatedAt       time.Time  `json:"created_at"`
}

// EquipmentStatistics 设备统计结构体
type EquipmentStatistics struct {
	TotalEquipment     int     `json:"total_equipment"`
	RunningCount       int     `json:"running_count"`
	StoppedCount       int     `json:"stopped_count"`
	MaintenanceCount   int     `json:"maintenance_count"`
	FaultCount         int     `json:"fault_count"`
	RunningRate        float64 `json:"running_rate"`
	MaintenanceRate    float64 `json:"maintenance_rate"`
	FaultRate          float64 `json:"fault_rate"`
}

// EquipmentService 设备服务
type EquipmentService struct {
	db *gorm.DB
}

// NewEquipmentService 创建设备服务实例
func NewEquipmentService(db *gorm.DB) *EquipmentService {
	return &EquipmentService{db: db}
}

// CreateEquipment 创建设备
func (s *EquipmentService) CreateEquipment(req *EquipmentRequest) (*EquipmentResponse, error) {
	// 检查设备编码是否已存在
	if s.isEquipmentCodeExists(req.Code, 0) {
		return nil, errors.New("设备编码已存在")
	}

	// 验证设备状态
	if !s.isValidEquipmentStatus(req.Status) {
		return nil, errors.New("无效的设备状态")
	}

	equipment := &models.Equipment{
		Code:         req.Code,
		Name:         req.Name,
		Type:         req.Type,
		Model:        req.Model,
		Manufacturer: req.Manufacturer,
		PurchaseDate: req.PurchaseDate,
		WarrantyDate: req.WarrantyDate,
		Location:     req.Location,
		Status:       req.Status,
		Description:  req.Description,
	}

	if err := s.db.Create(equipment).Error; err != nil {
		return nil, fmt.Errorf("创建设备失败: %v", err)
	}

	return s.equipmentToResponse(equipment), nil
}

// GetEquipment 获取设备详情
func (s *EquipmentService) GetEquipment(id uint) (*EquipmentResponse, error) {
	var equipment models.Equipment
	if err := s.db.First(&equipment, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("设备不存在")
		}
		return nil, fmt.Errorf("获取设备失败: %v", err)
	}

	return s.equipmentToResponse(&equipment), nil
}

// GetEquipmentList 获取设备列表
func (s *EquipmentService) GetEquipmentList(page, pageSize int, equipmentType, status, keyword string) ([]EquipmentResponse, int64, error) {
	var equipments []models.Equipment
	var total int64

	query := s.db.Model(&models.Equipment{})

	// 按类型筛选
	if equipmentType != "" {
		query = query.Where("type = ?", equipmentType)
	}

	// 按状态筛选
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 关键词搜索
	if keyword != "" {
		query = query.Where("code LIKE ? OR name LIKE ? OR location LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("获取设备总数失败: %v", err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&equipments).Error; err != nil {
		return nil, 0, fmt.Errorf("获取设备列表失败: %v", err)
	}

	var responses []EquipmentResponse
	for _, equipment := range equipments {
		responses = append(responses, *s.equipmentToResponse(&equipment))
	}

	return responses, total, nil
}

// UpdateEquipment 更新设备
func (s *EquipmentService) UpdateEquipment(id uint, req *EquipmentRequest) (*EquipmentResponse, error) {
	var equipment models.Equipment
	if err := s.db.First(&equipment, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("设备不存在")
		}
		return nil, fmt.Errorf("获取设备失败: %v", err)
	}

	// 检查设备编码是否已存在（排除当前设备）
	if s.isEquipmentCodeExists(req.Code, id) {
		return nil, errors.New("设备编码已存在")
	}

	// 验证设备状态
	if !s.isValidEquipmentStatus(req.Status) {
		return nil, errors.New("无效的设备状态")
	}

	// 更新设备信息
	equipment.Code = req.Code
	equipment.Name = req.Name
	equipment.Type = req.Type
	equipment.Model = req.Model
	equipment.Manufacturer = req.Manufacturer
	equipment.PurchaseDate = req.PurchaseDate
	equipment.WarrantyDate = req.WarrantyDate
	equipment.Location = req.Location
	equipment.Status = req.Status
	equipment.Description = req.Description

	if err := s.db.Save(&equipment).Error; err != nil {
		return nil, fmt.Errorf("更新设备失败: %v", err)
	}

	return s.equipmentToResponse(&equipment), nil
}

// DeleteEquipment 删除设备
func (s *EquipmentService) DeleteEquipment(id uint) error {
	var equipment models.Equipment
	if err := s.db.First(&equipment, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("设备不存在")
		}
		return fmt.Errorf("获取设备失败: %v", err)
	}

	// 检查是否有相关的维护记录
	var count int64
	if err := s.db.Model(&models.MaintenanceRecord{}).Where("equipment_id = ?", id).Count(&count).Error; err != nil {
		return fmt.Errorf("检查设备维护记录失败: %v", err)
	}

	if count > 0 {
		return errors.New("该设备存在维护记录，无法删除")
	}

	if err := s.db.Delete(&equipment).Error; err != nil {
		return fmt.Errorf("删除设备失败: %v", err)
	}

	return nil
}

// CreateMaintenanceRecord 创建维护记录
func (s *EquipmentService) CreateMaintenanceRecord(req *MaintenanceRecordRequest) (*MaintenanceRecordResponse, error) {
	// 验证设备是否存在
	var equipment models.Equipment
	if err := s.db.First(&equipment, req.EquipmentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("设备不存在")
		}
		return nil, fmt.Errorf("验证设备失败: %v", err)
	}

	// 验证维护人员是否存在
	var maintainer models.User
	if err := s.db.First(&maintainer, req.MaintainerID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("维护人员不存在")
		}
		return nil, fmt.Errorf("验证维护人员失败: %v", err)
	}

	// 验证维护类型
	if !s.isValidMaintenanceType(req.Type) {
		return nil, errors.New("无效的维护类型")
	}

	// 验证时间
	if req.EndTime != nil && req.EndTime.Before(req.StartTime) {
		return nil, errors.New("结束时间不能早于开始时间")
	}

	maintenanceRecord := &models.MaintenanceRecord{
		EquipmentID:     req.EquipmentID,
		MaintainerID:    req.MaintainerID,
		Type:            req.Type,
		Description:     req.Description,
		StartTime:       req.StartTime,
		EndTime:         req.EndTime,
		Cost:            req.Cost,
		PartsReplaced:   req.PartsReplaced,
		Result:          req.Result,
		NextMaintenance: req.NextMaintenance,
		Remark:          req.Remark,
	}

	if err := s.db.Create(maintenanceRecord).Error; err != nil {
		return nil, fmt.Errorf("创建维护记录失败: %v", err)
	}

	return s.maintenanceRecordToResponse(maintenanceRecord, &equipment, &maintainer), nil
}

// GetMaintenanceRecord 获取维护记录详情
func (s *EquipmentService) GetMaintenanceRecord(id uint) (*MaintenanceRecordResponse, error) {
	var maintenanceRecord models.MaintenanceRecord
	if err := s.db.Preload("Equipment").Preload("Maintainer").First(&maintenanceRecord, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("维护记录不存在")
		}
		return nil, fmt.Errorf("获取维护记录失败: %v", err)
	}

	return s.maintenanceRecordToResponse(&maintenanceRecord, &maintenanceRecord.Equipment, &maintenanceRecord.Maintainer), nil
}

// GetMaintenanceRecordList 获取维护记录列表
func (s *EquipmentService) GetMaintenanceRecordList(page, pageSize int, equipmentID, maintainerID uint, maintenanceType string) ([]MaintenanceRecordResponse, int64, error) {
	var maintenanceRecords []models.MaintenanceRecord
	var total int64

	query := s.db.Model(&models.MaintenanceRecord{}).Preload("Equipment").Preload("Maintainer")

	// 按设备筛选
	if equipmentID > 0 {
		query = query.Where("equipment_id = ?", equipmentID)
	}

	// 按维护人员筛选
	if maintainerID > 0 {
		query = query.Where("maintainer_id = ?", maintainerID)
	}

	// 按维护类型筛选
	if maintenanceType != "" {
		query = query.Where("type = ?", maintenanceType)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("获取维护记录总数失败: %v", err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("start_time DESC").Find(&maintenanceRecords).Error; err != nil {
		return nil, 0, fmt.Errorf("获取维护记录列表失败: %v", err)
	}

	var responses []MaintenanceRecordResponse
	for _, record := range maintenanceRecords {
		responses = append(responses, *s.maintenanceRecordToResponse(&record, &record.Equipment, &record.Maintainer))
	}

	return responses, total, nil
}

// UpdateMaintenanceRecord 更新维护记录
func (s *EquipmentService) UpdateMaintenanceRecord(id uint, req *MaintenanceRecordRequest) (*MaintenanceRecordResponse, error) {
	var maintenanceRecord models.MaintenanceRecord
	if err := s.db.Preload("Equipment").Preload("Maintainer").First(&maintenanceRecord, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("维护记录不存在")
		}
		return nil, fmt.Errorf("获取维护记录失败: %v", err)
	}

	// 验证设备是否存在
	var equipment models.Equipment
	if err := s.db.First(&equipment, req.EquipmentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("设备不存在")
		}
		return nil, fmt.Errorf("验证设备失败: %v", err)
	}

	// 验证维护人员是否存在
	var maintainer models.User
	if err := s.db.First(&maintainer, req.MaintainerID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("维护人员不存在")
		}
		return nil, fmt.Errorf("验证维护人员失败: %v", err)
	}

	// 验证维护类型
	if !s.isValidMaintenanceType(req.Type) {
		return nil, errors.New("无效的维护类型")
	}

	// 验证时间
	if req.EndTime != nil && req.EndTime.Before(req.StartTime) {
		return nil, errors.New("结束时间不能早于开始时间")
	}

	// 更新维护记录信息
	maintenanceRecord.EquipmentID = req.EquipmentID
	maintenanceRecord.MaintainerID = req.MaintainerID
	maintenanceRecord.Type = req.Type
	maintenanceRecord.Description = req.Description
	maintenanceRecord.StartTime = req.StartTime
	maintenanceRecord.EndTime = req.EndTime
	maintenanceRecord.Cost = req.Cost
	maintenanceRecord.PartsReplaced = req.PartsReplaced
	maintenanceRecord.Result = req.Result
	maintenanceRecord.NextMaintenance = req.NextMaintenance
	maintenanceRecord.Remark = req.Remark

	if err := s.db.Save(&maintenanceRecord).Error; err != nil {
		return nil, fmt.Errorf("更新维护记录失败: %v", err)
	}

	return s.maintenanceRecordToResponse(&maintenanceRecord, &equipment, &maintainer), nil
}

// GetEquipmentStatistics 获取设备统计数据
func (s *EquipmentService) GetEquipmentStatistics() (*EquipmentStatistics, error) {
	// 获取总设备数
	var totalEquipment int64
	if err := s.db.Model(&models.Equipment{}).Count(&totalEquipment).Error; err != nil {
		return nil, fmt.Errorf("获取总设备数失败: %v", err)
	}

	// 获取各状态设备数量
	var runningCount, stoppedCount, maintenanceCount, faultCount int64

	if err := s.db.Model(&models.Equipment{}).Where("status = ?", "running").Count(&runningCount).Error; err != nil {
		return nil, fmt.Errorf("获取运行设备数失败: %v", err)
	}

	if err := s.db.Model(&models.Equipment{}).Where("status = ?", "stopped").Count(&stoppedCount).Error; err != nil {
		return nil, fmt.Errorf("获取停机设备数失败: %v", err)
	}

	if err := s.db.Model(&models.Equipment{}).Where("status = ?", "maintenance").Count(&maintenanceCount).Error; err != nil {
		return nil, fmt.Errorf("获取维护设备数失败: %v", err)
	}

	if err := s.db.Model(&models.Equipment{}).Where("status = ?", "fault").Count(&faultCount).Error; err != nil {
		return nil, fmt.Errorf("获取故障设备数失败: %v", err)
	}

	// 计算比率
	var runningRate, maintenanceRate, faultRate float64
	if totalEquipment > 0 {
		runningRate = float64(runningCount) / float64(totalEquipment) * 100
		maintenanceRate = float64(maintenanceCount) / float64(totalEquipment) * 100
		faultRate = float64(faultCount) / float64(totalEquipment) * 100
	}

	return &EquipmentStatistics{
		TotalEquipment:   int(totalEquipment),
		RunningCount:     int(runningCount),
		StoppedCount:     int(stoppedCount),
		MaintenanceCount: int(maintenanceCount),
		FaultCount:       int(faultCount),
		RunningRate:      runningRate,
		MaintenanceRate:  maintenanceRate,
		FaultRate:        faultRate,
	}, nil
}

// GetEquipmentTypes 获取所有设备类型
func (s *EquipmentService) GetEquipmentTypes() ([]string, error) {
	var types []string
	if err := s.db.Model(&models.Equipment{}).Distinct("type").Pluck("type", &types).Error; err != nil {
		return nil, fmt.Errorf("获取设备类型失败: %v", err)
	}
	return types, nil
}

// GetMaintenanceTypes 获取所有维护类型
func (s *EquipmentService) GetMaintenanceTypes() []string {
	return []string{"preventive", "corrective", "emergency"}
}

// GetEquipmentStatuses 获取所有设备状态
func (s *EquipmentService) GetEquipmentStatuses() []string {
	return []string{"running", "stopped", "maintenance", "fault"}
}

// GetUpcomingMaintenances 获取即将到期的维护
func (s *EquipmentService) GetUpcomingMaintenances(days int) ([]MaintenanceRecordResponse, error) {
	if days <= 0 {
		days = 7 // 默认7天
	}

	futureDate := time.Now().AddDate(0, 0, days)

	var maintenanceRecords []models.MaintenanceRecord
	if err := s.db.Preload("Equipment").Preload("Maintainer").
		Where("next_maintenance IS NOT NULL AND next_maintenance <= ? AND next_maintenance >= ?", futureDate, time.Now()).
		Order("next_maintenance ASC").Find(&maintenanceRecords).Error; err != nil {
		return nil, fmt.Errorf("获取即将到期的维护失败: %v", err)
	}

	var responses []MaintenanceRecordResponse
	for _, record := range maintenanceRecords {
		responses = append(responses, *s.maintenanceRecordToResponse(&record, &record.Equipment, &record.Maintainer))
	}

	return responses, nil
}

// 辅助函数：检查设备编码是否存在
func (s *EquipmentService) isEquipmentCodeExists(code string, excludeID uint) bool {
	var count int64
	query := s.db.Model(&models.Equipment{}).Where("code = ?", code)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	query.Count(&count)
	return count > 0
}

// 辅助函数：验证设备状态
func (s *EquipmentService) isValidEquipmentStatus(status string) bool {
	validStatuses := []string{"running", "stopped", "maintenance", "fault"}
	for _, validStatus := range validStatuses {
		if status == validStatus {
			return true
		}
	}
	return false
}

// 辅助函数：验证维护类型
func (s *EquipmentService) isValidMaintenanceType(maintenanceType string) bool {
	validTypes := []string{"preventive", "corrective", "emergency"}
	for _, validType := range validTypes {
		if maintenanceType == validType {
			return true
		}
	}
	return false
}

// 辅助函数：将设备模型转换为响应结构体
func (s *EquipmentService) equipmentToResponse(equipment *models.Equipment) *EquipmentResponse {
	return &EquipmentResponse{
		ID:           equipment.ID,
		Code:         equipment.Code,
		Name:         equipment.Name,
		Type:         equipment.Type,
		Model:        equipment.Model,
		Manufacturer: equipment.Manufacturer,
		PurchaseDate: equipment.PurchaseDate,
		WarrantyDate: equipment.WarrantyDate,
		Location:     equipment.Location,
		Status:       equipment.Status,
		Description:  equipment.Description,
		CreatedAt:    equipment.CreatedAt,
		UpdatedAt:    equipment.UpdatedAt,
	}
}

// 辅助函数：将维护记录模型转换为响应结构体
func (s *EquipmentService) maintenanceRecordToResponse(record *models.MaintenanceRecord, equipment *models.Equipment, maintainer *models.User) *MaintenanceRecordResponse {
	response := &MaintenanceRecordResponse{
		ID:              record.ID,
		EquipmentID:     record.EquipmentID,
		EquipmentCode:   equipment.Code,
		EquipmentName:   equipment.Name,
		MaintainerID:    record.MaintainerID,
		MaintainerName:  maintainer.Username,
		Type:            record.Type,
		Description:     record.Description,
		StartTime:       record.StartTime,
		EndTime:         record.EndTime,
		Cost:            record.Cost,
		PartsReplaced:   record.PartsReplaced,
		Result:          record.Result,
		NextMaintenance: record.NextMaintenance,
		Remark:          record.Remark,
		CreatedAt:       record.CreatedAt,
	}

	// 计算维护时长（分钟）
	if record.EndTime != nil {
		duration := int(record.EndTime.Sub(record.StartTime).Minutes())
		response.Duration = &duration
	}

	return response
}