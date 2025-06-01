package controller

import (
	"net/http"
	"strconv"

	"mes-system/internal/service"
	"mes-system/pkg/response"

	"github.com/gin-gonic/gin"
)

// EquipmentController 设备控制器
type EquipmentController struct {
	equipmentService *service.EquipmentService
}

// NewEquipmentController 创建设备控制器实例
func NewEquipmentController(equipmentService *service.EquipmentService) *EquipmentController {
	return &EquipmentController{
		equipmentService: equipmentService,
	}
}

// CreateEquipment 创建设备
// @Summary 创建设备
// @Description 创建新的设备信息
// @Tags 设备管理
// @Accept json
// @Produce json
// @Param equipment body service.EquipmentRequest true "设备信息"
// @Success 200 {object} response.Response{data=service.EquipmentResponse}
// @Failure 400 {object} response.Response
// @Router /api/equipments [post]
func (c *EquipmentController) CreateEquipment(ctx *gin.Context) {
	var req service.EquipmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	equipment, err := c.equipmentService.CreateEquipment(&req)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessWithMessage(ctx, "创建设备成功", equipment)
}

// GetEquipment 获取设备详情
// @Summary 获取设备详情
// @Description 根据ID获取设备详细信息
// @Tags 设备管理
// @Accept json
// @Produce json
// @Param id path int true "设备ID"
// @Success 200 {object} response.Response{data=service.EquipmentResponse}
// @Failure 400 {object} response.Response
// @Router /api/equipments/{id} [get]
func (c *EquipmentController) GetEquipment(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "无效的设备ID")
		return
	}

	equipment, err := c.equipmentService.GetEquipment(uint(id))
	if err != nil {
		response.Error(ctx, http.StatusNotFound, err.Error())
		return
	}

	response.SuccessWithMessage(ctx, "获取设备详情成功", equipment)
}

// GetEquipmentList 获取设备列表
// @Summary 获取设备列表
// @Description 分页获取设备列表，支持按类型、状态和关键词筛选
// @Tags 设备管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param type query string false "设备类型"
// @Param status query string false "设备状态"
// @Param keyword query string false "搜索关键词"
// @Success 200 {object} response.Response{data=response.PageResponse}
// @Router /api/equipments [get]
func (c *EquipmentController) GetEquipmentList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	equipmentType := ctx.Query("type")
	status := ctx.Query("status")
	keyword := ctx.Query("keyword")

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	equipments, total, err := c.equipmentService.GetEquipmentList(page, pageSize, equipmentType, status, keyword)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithPage(ctx, equipments, total, page, pageSize, "获取设备列表成功")
}

// UpdateEquipment 更新设备
// @Summary 更新设备
// @Description 更新设备信息
// @Tags 设备管理
// @Accept json
// @Produce json
// @Param id path int true "设备ID"
// @Param equipment body service.EquipmentRequest true "设备信息"
// @Success 200 {object} response.Response{data=service.EquipmentResponse}
// @Failure 400 {object} response.Response
// @Router /api/equipments/{id} [put]
func (c *EquipmentController) UpdateEquipment(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "无效的设备ID")
		return
	}

	var req service.EquipmentRequest
	if err = ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	equipment, err := c.equipmentService.UpdateEquipment(uint(id), &req)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessWithMessage(ctx, "更新设备成功", equipment)
}

// DeleteEquipment 删除设备
// @Summary 删除设备
// @Description 删除设备信息
// @Tags 设备管理
// @Accept json
// @Produce json
// @Param id path int true "设备ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/equipments/{id} [delete]
func (c *EquipmentController) DeleteEquipment(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "无效的设备ID")
		return
	}

	if err := c.equipmentService.DeleteEquipment(uint(id)); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessWithMessage(ctx, "删除设备成功", nil)
}

// CreateMaintenanceRecord 创建维护记录
// @Summary 创建维护记录
// @Description 创建新的设备维护记录
// @Tags 设备管理
// @Accept json
// @Produce json
// @Param record body service.MaintenanceRecordRequest true "维护记录信息"
// @Success 200 {object} response.Response{data=service.MaintenanceRecordResponse}
// @Failure 400 {object} response.Response
// @Router /api/equipments/maintenance-records [post]
func (c *EquipmentController) CreateMaintenanceRecord(ctx *gin.Context) {
	var req service.MaintenanceRecordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	record, err := c.equipmentService.CreateMaintenanceRecord(&req)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessWithMessage(ctx, "创建维护记录成功", record)
}

// GetMaintenanceRecord 获取维护记录详情
// @Summary 获取维护记录详情
// @Description 根据ID获取维护记录详细信息
// @Tags 设备管理
// @Accept json
// @Produce json
// @Param id path int true "维护记录ID"
// @Success 200 {object} response.Response{data=service.MaintenanceRecordResponse}
// @Failure 400 {object} response.Response
// @Router /api/equipments/maintenance-records/{id} [get]
func (c *EquipmentController) GetMaintenanceRecord(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "无效的维护记录ID")
		return
	}

	record, err := c.equipmentService.GetMaintenanceRecord(uint(id))
	if err != nil {
		response.Error(ctx, http.StatusNotFound, err.Error())
		return
	}

	response.SuccessWithMessage(ctx, "获取维护记录详情成功", record)
}

// GetMaintenanceRecordList 获取维护记录列表
// @Summary 获取维护记录列表
// @Description 分页获取维护记录列表，支持多条件筛选
// @Tags 设备管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param equipment_id query int false "设备ID"
// @Param maintainer_id query int false "维护人员ID"
// @Param type query string false "维护类型"
// @Success 200 {object} response.Response{data=response.PageResponse}
// @Router /api/equipments/maintenance-records [get]
func (c *EquipmentController) GetMaintenanceRecordList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	equipmentIDStr := ctx.Query("equipment_id")
	maintainerIDStr := ctx.Query("maintainer_id")
	maintenanceType := ctx.Query("type")

	var equipmentID, maintainerID uint

	if equipmentIDStr != "" {
		id, err := strconv.ParseUint(equipmentIDStr, 10, 32)
		if err != nil {
			response.Error(ctx, http.StatusBadRequest, "无效的设备ID")
			return
		}
		equipmentID = uint(id)
	}

	if maintainerIDStr != "" {
		id, err := strconv.ParseUint(maintainerIDStr, 10, 32)
		if err != nil {
			response.Error(ctx, http.StatusBadRequest, "无效的维护人员ID")
			return
		}
		maintainerID = uint(id)
	}

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	records, total, err := c.equipmentService.GetMaintenanceRecordList(page, pageSize, equipmentID, maintainerID, maintenanceType)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithPage(ctx, records, total, page, pageSize, "获取维护记录列表成功")
}

// UpdateMaintenanceRecord 更新维护记录
// @Summary 更新维护记录
// @Description 更新维护记录信息
// @Tags 设备管理
// @Accept json
// @Produce json
// @Param id path int true "维护记录ID"
// @Param record body service.MaintenanceRecordRequest true "维护记录信息"
// @Success 200 {object} response.Response{data=service.MaintenanceRecordResponse}
// @Failure 400 {object} response.Response
// @Router /api/equipments/maintenance-records/{id} [put]
func (c *EquipmentController) UpdateMaintenanceRecord(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "无效的维护记录ID")
		return
	}

	var req service.MaintenanceRecordRequest
	if err = ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	record, err := c.equipmentService.UpdateMaintenanceRecord(uint(id), &req)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessWithMessage(ctx, "更新维护记录成功", record)
}

// GetEquipmentStatistics 获取设备统计数据
// @Summary 获取设备统计数据
// @Description 获取设备状态统计信息
// @Tags 设备管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=service.EquipmentStatistics}
// @Router /api/equipments/statistics [get]
func (c *EquipmentController) GetEquipmentStatistics(ctx *gin.Context) {
	statistics, err := c.equipmentService.GetEquipmentStatistics()
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithMessage(ctx, "获取设备统计数据成功", statistics)
}

// GetEquipmentTypes 获取设备类型
// @Summary 获取设备类型
// @Description 获取所有设备类型列表
// @Tags 设备管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]string}
// @Router /api/equipments/types [get]
func (c *EquipmentController) GetEquipmentTypes(ctx *gin.Context) {
	types, err := c.equipmentService.GetEquipmentTypes()
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithMessage(ctx, "获取设备类型成功", types)
}

// GetMaintenanceTypes 获取维护类型
// @Summary 获取维护类型
// @Description 获取所有维护类型列表
// @Tags 设备管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]string}
// @Router /api/equipments/maintenance-types [get]
func (c *EquipmentController) GetMaintenanceTypes(ctx *gin.Context) {
	types := c.equipmentService.GetMaintenanceTypes()
	response.SuccessWithMessage(ctx, "获取维护类型成功", types)
}

// GetEquipmentStatuses 获取设备状态
// @Summary 获取设备状态
// @Description 获取所有设备状态列表
// @Tags 设备管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]string}
// @Router /api/equipments/statuses [get]
func (c *EquipmentController) GetEquipmentStatuses(ctx *gin.Context) {
	statuses := c.equipmentService.GetEquipmentStatuses()
	response.SuccessWithMessage(ctx, "获取设备状态成功", statuses)
}

// GetUpcomingMaintenances 获取即将到期的维护
// @Summary 获取即将到期的维护
// @Description 获取指定天数内即将到期的维护记录
// @Tags 设备管理
// @Accept json
// @Produce json
// @Param days query int false "天数" default(7)
// @Success 200 {object} response.Response{data=[]service.MaintenanceRecordResponse}
// @Router /api/equipments/upcoming-maintenances [get]
func (c *EquipmentController) GetUpcomingMaintenances(ctx *gin.Context) {
	days, _ := strconv.Atoi(ctx.DefaultQuery("days", "7"))
	if days <= 0 {
		days = 7
	}

	maintenances, err := c.equipmentService.GetUpcomingMaintenances(days)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithMessage(ctx, "获取即将到期的维护成功", maintenances)
}

// DeleteMaintenanceRecord 删除维护记录
// @Summary 删除维护记录
// @Description 删除指定ID的维护记录
// @Tags 设备管理
// @Accept json
// @Produce json
// @Param id path int true "维护记录ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/equipment/maintenance/{id} [delete]
func (c *EquipmentController) DeleteMaintenanceRecord(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "无效的维护记录ID")
		return
	}

	err = c.equipmentService.DeleteMaintenanceRecord(uint(id))
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessWithMessage(ctx, "删除维护记录成功", nil)
}
