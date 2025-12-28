package handlers

import (
	"net/http"
	"strconv"

	"fluent-life-admin-api/internal/models"
	"fluent-life-admin-api/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminHandler struct {
	db *gorm.DB
}

func NewAdminHandler(db *gorm.DB) *AdminHandler {
	return &AdminHandler{db: db}
}

// 管理员登录（简化版，实际应该使用JWT）
func (h *AdminHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	// 简单的管理员验证（实际应该从数据库验证）
	if req.Username == "admin" && req.Password == "admin123" {
		response.Success(c, gin.H{
			"token": "admin_token_12345",
			"username": "admin",
		}, "登录成功")
		return
	}

	response.Error(c, http.StatusUnauthorized, "用户名或密码错误")
}

// 获取用户列表
func (h *AdminHandler) GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	offset := (page - 1) * pageSize

	var users []models.User
	var total int64

	query := h.db.Model(&models.User{})
	
	// 搜索
	if keyword := c.Query("keyword"); keyword != "" {
		query = query.Where("username LIKE ? OR email LIKE ? OR phone LIKE ?", 
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)
	
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&users).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "查询失败")
		return
	}

	response.Success(c, gin.H{
		"users": users,
		"total": total,
		"page": page,
		"page_size": pageSize,
	}, "获取成功")
}

// 获取用户详情
func (h *AdminHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	
	var user models.User
	if err := h.db.Where("id = ?", id).First(&user).Error; err != nil {
		response.Error(c, http.StatusNotFound, "用户不存在")
		return
	}

	response.Success(c, user, "获取成功")
}

// 删除用户
func (h *AdminHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	
	if err := h.db.Where("id = ?", id).Delete(&models.User{}).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "删除失败")
		return
	}

	response.Success(c, nil, "删除成功")
}

// 获取帖子列表
func (h *AdminHandler) GetPosts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	offset := (page - 1) * pageSize

	var posts []models.Post
	var total int64

	query := h.db.Model(&models.Post{}).Preload("User")
	
	// 搜索
	if keyword := c.Query("keyword"); keyword != "" {
		query = query.Where("content LIKE ?", "%"+keyword+"%")
	}

	query.Count(&total)
	
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&posts).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "查询失败")
		return
	}

	response.Success(c, gin.H{
		"posts": posts,
		"total": total,
		"page": page,
		"page_size": pageSize,
	}, "获取成功")
}

// 获取帖子详情
func (h *AdminHandler) GetPost(c *gin.Context) {
	id := c.Param("id")
	
	var post models.Post
	if err := h.db.Preload("User").Preload("Comments.User").Where("id = ?", id).First(&post).Error; err != nil {
		response.Error(c, http.StatusNotFound, "帖子不存在")
		return
	}

	response.Success(c, post, "获取成功")
}

// 删除帖子
func (h *AdminHandler) DeletePost(c *gin.Context) {
	id := c.Param("id")
	
	if err := h.db.Where("id = ?", id).Delete(&models.Post{}).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "删除失败")
		return
	}

	response.Success(c, nil, "删除成功")
}

// 获取房间列表
func (h *AdminHandler) GetRooms(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	offset := (page - 1) * pageSize

	var rooms []models.PracticeRoom
	var total int64

	query := h.db.Model(&models.PracticeRoom{}).Preload("User")
	
	// 搜索
	if keyword := c.Query("keyword"); keyword != "" {
		query = query.Where("title LIKE ? OR theme LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 筛选活跃状态
	if isActive := c.Query("is_active"); isActive != "" {
		active := isActive == "true"
		query = query.Where("is_active = ?", active)
	}

	query.Count(&total)
	
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&rooms).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "查询失败")
		return
	}

	response.Success(c, gin.H{
		"rooms": rooms,
		"total": total,
		"page": page,
		"page_size": pageSize,
	}, "获取成功")
}

// 获取房间详情
func (h *AdminHandler) GetRoom(c *gin.Context) {
	id := c.Param("id")
	
	var room models.PracticeRoom
	if err := h.db.Preload("User").Preload("Members.User").Where("id = ?", id).First(&room).Error; err != nil {
		response.Error(c, http.StatusNotFound, "房间不存在")
		return
	}

	response.Success(c, room, "获取成功")
}

// 删除房间
func (h *AdminHandler) DeleteRoom(c *gin.Context) {
	id := c.Param("id")
	
	if err := h.db.Where("id = ?", id).Delete(&models.PracticeRoom{}).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "删除失败")
		return
	}

	response.Success(c, nil, "删除成功")
}

// 关闭/开启房间
func (h *AdminHandler) ToggleRoom(c *gin.Context) {
	id := c.Param("id")
	
	var room models.PracticeRoom
	if err := h.db.Where("id = ?", id).First(&room).Error; err != nil {
		response.Error(c, http.StatusNotFound, "房间不存在")
		return
	}

	room.IsActive = !room.IsActive
	if err := h.db.Save(&room).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "操作失败")
		return
	}

	response.Success(c, room, "操作成功")
}

// 获取训练统计
func (h *AdminHandler) GetTrainingStats(c *gin.Context) {
	var stats struct {
		TotalRecords    int64 `json:"total_records"`
		TotalUsers      int64 `json:"total_users"`
		MeditationCount int64 `json:"meditation_count"`
		AirflowCount    int64 `json:"airflow_count"`
		ExposureCount   int64 `json:"exposure_count"`
		PracticeCount   int64 `json:"practice_count"`
	}

	h.db.Model(&models.TrainingRecord{}).Count(&stats.TotalRecords)
	h.db.Model(&models.User{}).Count(&stats.TotalUsers)
	h.db.Model(&models.TrainingRecord{}).Where("type = ?", "meditation").Count(&stats.MeditationCount)
	h.db.Model(&models.TrainingRecord{}).Where("type = ?", "airflow").Count(&stats.AirflowCount)
	h.db.Model(&models.TrainingRecord{}).Where("type = ?", "exposure").Count(&stats.ExposureCount)
	h.db.Model(&models.TrainingRecord{}).Where("type = ?", "practice").Count(&stats.PracticeCount)

	response.Success(c, stats, "获取成功")
}

// 获取训练记录列表
func (h *AdminHandler) GetTrainingRecords(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	offset := (page - 1) * pageSize

	var records []models.TrainingRecord
	var total int64

	query := h.db.Model(&models.TrainingRecord{}).Preload("User")
	
	// 按类型筛选
	if recordType := c.Query("type"); recordType != "" {
		query = query.Where("type = ?", recordType)
	}

	query.Count(&total)
	
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&records).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "查询失败")
		return
	}

	response.Success(c, gin.H{
		"records": records,
		"total": total,
		"page": page,
		"page_size": pageSize,
	}, "获取成功")
}

