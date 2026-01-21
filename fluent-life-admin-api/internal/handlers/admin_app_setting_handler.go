package handlers

import (
	"net/http"
	"strings"

	"fluent-life-admin-api/internal/models"
	"fluent-life-admin-api/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminAppSettingHandler struct {
	db *gorm.DB
}

func NewAdminAppSettingHandler(db *gorm.DB) *AdminAppSettingHandler {
	return &AdminAppSettingHandler{db: db}
}

// GetAppSettings 获取所有应用设置（管理员）
// GET /api/v1/admin/app-settings
func (h *AdminAppSettingHandler) GetAppSettings(c *gin.Context) {
	var settings []models.AppSetting
	if err := h.db.Order("key ASC").Find(&settings).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "获取应用设置失败")
		return
	}
	response.Success(c, gin.H{"settings": settings, "total": len(settings)}, "获取成功")
}

// CreateAppSetting 创建应用设置（管理员）
// POST /api/v1/admin/app-settings
func (h *AdminAppSettingHandler) CreateAppSetting(c *gin.Context) {
	var req struct {
		Key         string `json:"key" binding:"required"`
		Value       string `json:"value" binding:"required"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}
	key := strings.TrimSpace(req.Key)
	if key == "" {
		response.Error(c, http.StatusBadRequest, "key 不能为空")
		return
	}

	setting := models.AppSetting{Key: key, Value: req.Value, Description: req.Description}
	if err := h.db.Create(&setting).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "创建应用设置失败")
		return
	}
	response.Success(c, setting, "创建成功")
}

// UpdateAppSetting 更新应用设置（管理员）
// PUT /api/v1/admin/app-settings/:id
func (h *AdminAppSettingHandler) UpdateAppSetting(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Key         *string `json:"key"`
		Value       *string `json:"value"`
		Description *string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	var setting models.AppSetting
	if err := h.db.Where("id = ?", id).First(&setting).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Error(c, http.StatusNotFound, "应用设置不存在")
			return
		}
		response.Error(c, http.StatusInternalServerError, "获取应用设置失败")
		return
	}

	if req.Key != nil {
		k := strings.TrimSpace(*req.Key)
		if k == "" {
			response.Error(c, http.StatusBadRequest, "key 不能为空")
			return
		}
		setting.Key = k
	}
	if req.Value != nil {
		setting.Value = *req.Value
	}
	if req.Description != nil {
		setting.Description = *req.Description
	}

	if err := h.db.Save(&setting).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "更新应用设置失败")
		return
	}
	response.Success(c, setting, "更新成功")
}

// DeleteAppSetting 删除应用设置（管理员）
// DELETE /api/v1/admin/app-settings/:id
func (h *AdminAppSettingHandler) DeleteAppSetting(c *gin.Context) {
	id := c.Param("id")

	if err := h.db.Delete(&models.AppSetting{}, "id = ?", id).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "删除应用设置失败")
		return
	}
	response.Success(c, nil, "删除成功")
}
