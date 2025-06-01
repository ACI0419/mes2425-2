package controller

import (
	"net/http"
	"strconv"
	"time"

	"mes-system/internal/service"
	"mes-system/pkg/response"

	"github.com/gin-gonic/gin"
)

// QualityController 质量控制器
type QualityController struct {
	qualityService *service.QualityService
}

// NewQualityController 创建质量控制器实例
func NewQualityController(qualityService *service.QualityService) *QualityController {
	return &QualityController{
		qualityService: qualityService,
	}
}

// CreateQualityStandard 创建质量标准
// @Summary 创建质量标准
// @Description 创建新的质量标准
// @Tags 质量管理
// @Accept json
// @Produce json
// @Param standard body service.QualityStandardRequest true "质量标准信息"
// @Success 200 {object} response.Response{data=service.QualityStandardResponse}
// @Failure 400 {object} response.Response
// @Router /api/quality/standards [post]
func (c *QualityController) CreateQualityStandard(ctx *gin.Context) {
	var req service.QualityStandardRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	standard, err := c.qualityService.CreateQualityStandard(&req)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessWithMessage(ctx, "创建质量标准成功", standard)
}

// GetQualityStandard 获取质量标准详情
// @Summary 获取质量标准详情
// @Description 根据ID获取质量标准详细信息
// @Tags 质量管理
// @Accept json
// @Produce json
// @Param id path int true "质量标准ID"
// @Success 200 {object} response.Response{data=service.QualityStandardResponse}
// @Failure 400 {object} response.Response
// @Router /api/quality/standards/{id} [get]
func (c *QualityController) GetQualityStandard(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "无效的质量标准ID")
		return
	}

	standard, err := c.qualityService.GetQualityStandard(uint(id))
	if err != nil {
		response.Error(ctx, http.StatusNotFound, err.Error())
		return
	}

	response.SuccessWithMessage(ctx, "获取质量标准详情成功", standard)
}

// GetQualityStandardList 获取质量标准列表
// @Summary 获取质量标准列表
// @Description 分页获取质量标准列表，支持按产品、类型和状态筛选
// @Tags 质量管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param product_id query int false "产品ID"
// @Param type query string false "标准类型"
// @Param is_active query bool false "是否启用"
// @Success 200 {object} response.Response{data=response.PageResponse}
// @Router /api/quality/standards [get]
func (c *QualityController) GetQualityStandardList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	productIDStr := ctx.Query("product_id")
	standardType := ctx.Query("type")
	isActiveStr := ctx.Query("is_active")

	var productID uint
	if productIDStr != "" {
		id, err := strconv.ParseUint(productIDStr, 10, 32)
		if err != nil {
			response.Error(ctx, http.StatusBadRequest, "无效的产品ID")
			return
		}
		productID = uint(id)
	}

	var isActive *bool
	if isActiveStr != "" {
		active, err := strconv.ParseBool(isActiveStr)
		if err != nil {
			response.Error(ctx, http.StatusBadRequest, "无效的状态参数")
			return
		}
		isActive = &active
	}

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	standards, total, err := c.qualityService.GetQualityStandardList(page, pageSize, productID, standardType, isActive)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithPage(ctx, standards, total, page, pageSize, "获取质量标准列表成功")
}

// UpdateQualityStandard 更新质量标准
// @Summary 更新质量标准
// @Description 更新质量标准信息
// @Tags 质量管理
// @Accept json
// @Produce json
// @Param id path int true "质量标准ID"
// @Param standard body service.QualityStandardRequest true "质量标准信息"
// @Success 200 {object} response.Response{data=service.QualityStandardResponse}
// @Failure 400 {object} response.Response
// @Router /api/quality/standards/{id} [put]
func (c *QualityController) UpdateQualityStandard(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "无效的质量标准ID")
		return
	}

	var req service.QualityStandardRequest
	if err = ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	standard, err := c.qualityService.UpdateQualityStandard(uint(id), &req)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessWithMessage(ctx, "更新质量标准成功", standard)
}

// DeleteQualityStandard 删除质量标准
// @Summary 删除质量标准
// @Description 删除质量标准
// @Tags 质量管理
// @Accept json
// @Produce json
// @Param id path int true "质量标准ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/quality/standards/{id} [delete]
func (c *QualityController) DeleteQualityStandard(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "无效的质量标准ID")
		return
	}

	if err := c.qualityService.DeleteQualityStandard(uint(id)); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessWithMessage(ctx, "删除质量标准成功", nil)
}

// CreateQualityInspection 创建质量检测记录
// @Summary 创建质量检测记录
// @Description 创建新的质量检测记录
// @Tags 质量管理
// @Accept json
// @Produce json
// @Param inspection body service.QualityInspectionRequest true "质量检测信息"
// @Success 200 {object} response.Response{data=service.QualityInspectionResponse}
// @Failure 400 {object} response.Response
// @Router /api/quality/inspections [post]
func (c *QualityController) CreateQualityInspection(ctx *gin.Context) {
	var req service.QualityInspectionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	inspection, err := c.qualityService.CreateQualityInspection(&req)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessWithMessage(ctx, "创建质量检测记录成功", inspection)
}

// GetQualityInspection 获取质量检测记录详情
// @Summary 获取质量检测记录详情
// @Description 根据ID获取质量检测记录详细信息
// @Tags 质量管理
// @Accept json
// @Produce json
// @Param id path int true "质量检测记录ID"
// @Success 200 {object} response.Response{data=service.QualityInspectionResponse}
// @Failure 400 {object} response.Response
// @Router /api/quality/inspections/{id} [get]
func (c *QualityController) GetQualityInspection(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "无效的质量检测记录ID")
		return
	}

	inspection, err := c.qualityService.GetQualityInspection(uint(id))
	if err != nil {
		response.Error(ctx, http.StatusNotFound, err.Error())
		return
	}

	response.SuccessWithMessage(ctx, "获取质量检测记录详情成功", inspection)
}

// GetQualityInspectionList 获取质量检测记录列表
// @Summary 获取质量检测记录列表
// @Description 分页获取质量检测记录列表，支持多条件筛选
// @Tags 质量管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param production_order_id query int false "生产工单ID"
// @Param quality_standard_id query int false "质量标准ID"
// @Param inspector_id query int false "检测员ID"
// @Param result query string false "检测结果(pass/fail)"
// @Success 200 {object} response.Response{data=response.PageResponse}
// @Router /api/quality/inspections [get]
func (c *QualityController) GetQualityInspectionList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	productionOrderIDStr := ctx.Query("production_order_id")
	qualityStandardIDStr := ctx.Query("quality_standard_id")
	inspectorIDStr := ctx.Query("inspector_id")
	result := ctx.Query("result")

	var productionOrderID, qualityStandardID, inspectorID uint

	if productionOrderIDStr != "" {
		id, err := strconv.ParseUint(productionOrderIDStr, 10, 32)
		if err != nil {
			response.Error(ctx, http.StatusBadRequest, "无效的生产工单ID")
			return
		}
		productionOrderID = uint(id)
	}

	if qualityStandardIDStr != "" {
		id, err := strconv.ParseUint(qualityStandardIDStr, 10, 32)
		if err != nil {
			response.Error(ctx, http.StatusBadRequest, "无效的质量标准ID")
			return
		}
		qualityStandardID = uint(id)
	}

	if inspectorIDStr != "" {
		id, err := strconv.ParseUint(inspectorIDStr, 10, 32)
		if err != nil {
			response.Error(ctx, http.StatusBadRequest, "无效的检测员ID")
			return
		}
		inspectorID = uint(id)
	}

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	inspections, total, err := c.qualityService.GetQualityInspectionList(page, pageSize, productionOrderID, qualityStandardID, inspectorID, result)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithPage(ctx, inspections, total, page, pageSize, "获取质量检测记录列表成功")
}

// GetQualityStatistics 获取质量统计数据
// @Summary 获取质量统计数据
// @Description 获取质量检测统计数据，支持时间范围和条件筛选
// @Tags 质量管理
// @Accept json
// @Produce json
// @Param start_date query string false "开始日期(YYYY-MM-DD)"
// @Param end_date query string false "结束日期(YYYY-MM-DD)"
// @Param production_order_id query int false "生产工单ID"
// @Param quality_standard_id query int false "质量标准ID"
// @Success 200 {object} response.Response{data=service.QualityStatistics}
// @Router /api/quality/statistics [get]
func (c *QualityController) GetQualityStatistics(ctx *gin.Context) {
	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")
	productionOrderIDStr := ctx.Query("production_order_id")
	qualityStandardIDStr := ctx.Query("quality_standard_id")

	var startDate, endDate *time.Time
	var productionOrderID, qualityStandardID uint

	if startDateStr != "" {
		parsedDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			response.Error(ctx, http.StatusBadRequest, "无效的开始日期格式")
			return
		}
		startDate = &parsedDate
	}

	if endDateStr != "" {
		parsedDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			response.Error(ctx, http.StatusBadRequest, "无效的结束日期格式")
			return
		}
		// 设置为当天的23:59:59
		parsedDate = parsedDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		endDate = &parsedDate
	}

	if productionOrderIDStr != "" {
		id, err := strconv.ParseUint(productionOrderIDStr, 10, 32)
		if err != nil {
			response.Error(ctx, http.StatusBadRequest, "无效的生产工单ID")
			return
		}
		productionOrderID = uint(id)
	}

	if qualityStandardIDStr != "" {
		id, err := strconv.ParseUint(qualityStandardIDStr, 10, 32)
		if err != nil {
			response.Error(ctx, http.StatusBadRequest, "无效的质量标准ID")
			return
		}
		qualityStandardID = uint(id)
	}

	statistics, err := c.qualityService.GetQualityStatistics(startDate, endDate, productionOrderID, qualityStandardID)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithMessage(ctx, "获取质量统计数据成功", statistics)
}

// GetQualityStandardTypes 获取质量标准类型
// @Summary 获取质量标准类型
// @Description 获取所有质量标准类型列表
// @Tags 质量管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]string}
// @Router /api/quality/standards/types [get]
func (c *QualityController) GetQualityStandardTypes(ctx *gin.Context) {
	types, err := c.qualityService.GetQualityStandardTypes()
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithMessage(ctx, "获取质量标准类型成功", types)
}
