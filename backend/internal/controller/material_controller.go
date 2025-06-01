package controller

import (
	"net/http"
	"strconv"

	"mes-system/internal/service"
	"mes-system/pkg/response"

	"github.com/gin-gonic/gin"
)

// MaterialController 物料控制器
type MaterialController struct {
	materialService *service.MaterialService
}

// NewMaterialController 创建物料控制器实例
func NewMaterialController(materialService *service.MaterialService) *MaterialController {
	return &MaterialController{
		materialService: materialService,
	}
}

// CreateMaterial 创建物料
// @Summary 创建物料
// @Description 创建新的物料信息
// @Tags 物料管理
// @Accept json
// @Produce json
// @Param material body service.MaterialRequest true "物料信息"
// @Success 200 {object} response.Response{data=service.MaterialResponse}
// @Failure 400 {object} response.Response
// @Router /api/materials [post]
func (c *MaterialController) CreateMaterial(ctx *gin.Context) {
	var req service.MaterialRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	material, err := c.materialService.CreateMaterial(&req)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessWithMessage(ctx, "创建物料成功", material)
}

// GetMaterial 获取物料详情
// @Summary 获取物料详情
// @Description 根据ID获取物料详细信息
// @Tags 物料管理
// @Accept json
// @Produce json
// @Param id path int true "物料ID"
// @Success 200 {object} response.Response{data=service.MaterialResponse}
// @Failure 400 {object} response.Response
// @Router /api/materials/{id} [get]
func (c *MaterialController) GetMaterial(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "无效的物料ID")
		return
	}

	material, err := c.materialService.GetMaterial(uint(id))
	if err != nil {
		response.Error(ctx, http.StatusNotFound, err.Error())
		return
	}

	response.SuccessWithMessage(ctx, "获取物料详情成功", material)
}

// GetMaterialList 获取物料列表
// @Summary 获取物料列表
// @Description 分页获取物料列表，支持按类型和关键词筛选
// @Tags 物料管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param type query string false "物料类型"
// @Param keyword query string false "搜索关键词"
// @Success 200 {object} response.Response{data=response.PageResponse}
// @Router /api/materials [get]
func (c *MaterialController) GetMaterialList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	materialType := ctx.Query("type")
	keyword := ctx.Query("keyword")

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	materials, total, err := c.materialService.GetMaterialList(page, pageSize, materialType, keyword)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithPage(ctx, materials, total, page, pageSize, "获取物料列表成功")
}

// UpdateMaterial 更新物料
// @Summary 更新物料
// @Description 更新物料信息
// @Tags 物料管理
// @Accept json
// @Produce json
// @Param id path int true "物料ID"
// @Param material body service.MaterialRequest true "物料信息"
// @Success 200 {object} response.Response{data=service.MaterialResponse}
// @Failure 400 {object} response.Response
// @Router /api/materials/{id} [put]
func (c *MaterialController) UpdateMaterial(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "无效的物料ID")
		return
	}

	var req service.MaterialRequest
	if err = ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	material, err := c.materialService.UpdateMaterial(uint(id), &req)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessWithMessage(ctx, "更新物料成功", material)
}

// DeleteMaterial 删除物料
// @Summary 删除物料
// @Description 删除物料信息
// @Tags 物料管理
// @Accept json
// @Produce json
// @Param id path int true "物料ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/materials/{id} [delete]
func (c *MaterialController) DeleteMaterial(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "无效的物料ID")
		return
	}

	if err := c.materialService.DeleteMaterial(uint(id)); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessWithMessage(ctx, "删除物料成功", nil)
}

// CreateTransaction 创建物料交易
// @Summary 创建物料交易
// @Description 创建物料入库或出库记录
// @Tags 物料管理
// @Accept json
// @Produce json
// @Param transaction body service.MaterialTransactionRequest true "交易信息"
// @Success 200 {object} response.Response{data=service.MaterialTransactionResponse}
// @Failure 400 {object} response.Response
// @Router /api/materials/transactions [post]
func (c *MaterialController) CreateTransaction(ctx *gin.Context) {
	var req service.MaterialTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	transaction, err := c.materialService.CreateTransaction(&req)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessWithMessage(ctx, "创建交易记录成功", transaction)
}

// GetTransactionList 获取物料交易列表
// @Summary 获取物料交易列表
// @Description 分页获取物料交易记录，支持按物料和交易类型筛选
// @Tags 物料管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param material_id query int false "物料ID"
// @Param type query string false "交易类型(in/out)"
// @Success 200 {object} response.Response{data=response.PageResponse}
// @Router /api/materials/transactions [get]
func (c *MaterialController) GetTransactionList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	materialIDStr := ctx.Query("material_id")
	transactionType := ctx.Query("type")

	var materialID uint
	if materialIDStr != "" {
		id, err := strconv.ParseUint(materialIDStr, 10, 32)
		if err != nil {
			response.Error(ctx, http.StatusBadRequest, "无效的物料ID")
			return
		}
		materialID = uint(id)
	}

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	transactions, total, err := c.materialService.GetTransactionList(page, pageSize, materialID, transactionType)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithPage(ctx, transactions, total, page, pageSize, "获取交易列表成功")
}

// GetLowStockMaterials 获取低库存物料
// @Summary 获取低库存物料
// @Description 获取库存低于最小库存的物料列表
// @Tags 物料管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]service.MaterialResponse}
// @Router /api/materials/low-stock [get]
func (c *MaterialController) GetLowStockMaterials(ctx *gin.Context) {
	materials, err := c.materialService.GetLowStockMaterials()
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithMessage(ctx, "获取低库存物料成功", materials)
}

// GetMaterialTypes 获取物料类型
// @Summary 获取物料类型
// @Description 获取所有物料类型列表
// @Tags 物料管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]string}
// @Router /api/materials/types [get]
func (c *MaterialController) GetMaterialTypes(ctx *gin.Context) {
	types, err := c.materialService.GetMaterialTypes()
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithMessage(ctx, "获取物料类型成功", types)
}
