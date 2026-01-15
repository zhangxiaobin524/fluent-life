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

// 删除所有绕口令
func (h *AdminHandler) DeleteAllTongueTwisters(c *gin.Context) {
	if err := h.db.Where("1 = 1").Delete(&models.TongueTwister{}).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "删除所有绕口令失败")
		return
	}
	response.Success(c, nil, "所有绕口令删除成功")
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

// 删除帖子 (支持批量删除)
func (h *AdminHandler) DeletePost(c *gin.Context) {
	var req struct {
		IDs []string `json:"ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误，需要提供帖子ID列表")
		return
	}

	tx := h.db.Begin()
	if tx.Error != nil {
		response.Error(c, http.StatusInternalServerError, "删除失败")
		return
	}

	// 先删除相关的点赞记录
	if err := tx.Where("post_id IN ?", req.IDs).Delete(&models.PostLike{}).Error; err != nil {
		tx.Rollback()
		response.Error(c, http.StatusInternalServerError, "删除相关点赞失败")
		return
	}

	// 再删除相关的评论记录
	if err := tx.Where("post_id IN ?", req.IDs).Delete(&models.Comment{}).Error; err != nil {
		tx.Rollback()
		response.Error(c, http.StatusInternalServerError, "删除相关评论失败")
		return
	}

	// 再删除相关的收藏记录
	if err := tx.Where("post_id IN ?", req.IDs).Delete(&models.PostCollection{}).Error; err != nil {
		tx.Rollback()
		response.Error(c, http.StatusInternalServerError, "删除相关收藏失败")
		return
	}

	// 最后删除帖子
	if err := tx.Where("id IN ?", req.IDs).Delete(&models.Post{}).Error; err != nil {
		tx.Rollback()
		response.Error(c, http.StatusInternalServerError, "删除帖子失败")
		return
	}

	tx.Commit()
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

// 删除房间 (支持批量删除)
func (h *AdminHandler) DeleteRoom(c *gin.Context) {
	var req struct {
		IDs []string `json:"ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误，需要提供房间ID列表")
		return
	}

	tx := h.db.Begin()
	if tx.Error != nil {
		response.Error(c, http.StatusInternalServerError, "删除失败")
		return
	}

	// 先删除相关的房间成员记录
	if err := tx.Where("room_id IN ?", req.IDs).Delete(&models.PracticeRoomMember{}).Error; err != nil {
		tx.Rollback()
		response.Error(c, http.StatusInternalServerError, "删除相关房间成员失败")
		return
	}

	// 最后删除房间
	if err := tx.Where("id IN ?", req.IDs).Delete(&models.PracticeRoom{}).Error; err != nil {
		tx.Rollback()
		response.Error(c, http.StatusInternalServerError, "删除房间失败")
		return
	}

	tx.Commit()
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

	// 按用户ID筛选
	if userID := c.Query("user_id"); userID != "" {
		query = query.Where("user_id = ?", userID)
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

// ========== 绕口令管理 ==========

// 获取绕口令列表
func (h *AdminHandler) GetTongueTwisters(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	offset := (page - 1) * pageSize

	var tongueTwisters []models.TongueTwister
	var total int64

	query := h.db.Model(&models.TongueTwister{})
	
	// 搜索
	if keyword := c.Query("keyword"); keyword != "" {
		query = query.Where("title LIKE ? OR content LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 按难度筛选
	if level := c.Query("level"); level != "" {
		query = query.Where("level = ?", level)
	}

	// 按状态筛选
	if isActive := c.Query("is_active"); isActive != "" {
		active := isActive == "true"
		query = query.Where("is_active = ?", active)
	}

	query.Count(&total)
	
	if err := query.Offset(offset).Limit(pageSize).Order("level ASC, \"order\" ASC, created_at DESC").Find(&tongueTwisters).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "查询失败")
		return
	}

	response.Success(c, gin.H{
		"tongue_twisters": tongueTwisters,
		"total":           total,
		"page":            page,
		"page_size":       pageSize,
	}, "获取成功")
}

// 获取绕口令详情
func (h *AdminHandler) GetTongueTwister(c *gin.Context) {
	id := c.Param("id")
	
	var tongueTwister models.TongueTwister
	if err := h.db.Where("id = ?", id).First(&tongueTwister).Error; err != nil {
		response.Error(c, http.StatusNotFound, "绕口令不存在")
		return
	}

	response.Success(c, tongueTwister, "获取成功")
}

// 创建绕口令
func (h *AdminHandler) CreateTongueTwister(c *gin.Context) {
	var req models.TongueTwister
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := h.db.Create(&req).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "创建失败")
		return
	}

	response.Success(c, req, "创建成功")
}

// 批量创建绕口令
func (h *AdminHandler) BatchCreateTongueTwisters(c *gin.Context) {
	var req []models.TongueTwister
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if len(req) == 0 {
		response.Error(c, http.StatusBadRequest, "请求体不能为空")
		return
	}

	tx := h.db.Begin()
	if tx.Error != nil {
		response.Error(c, http.StatusInternalServerError, "创建失败")
		return
	}

	for _, twister := range req {
		if err := tx.Create(&twister).Error; err != nil {
			tx.Rollback()
			response.Error(c, http.StatusInternalServerError, "批量创建失败: "+err.Error())
			return
		}
	}

	tx.Commit()
	response.Success(c, nil, "批量创建成功")
}

// 更新绕口令
func (h *AdminHandler) UpdateTongueTwister(c *gin.Context) {
	id := c.Param("id")
	
	var tongueTwister models.TongueTwister
	if err := h.db.Where("id = ?", id).First(&tongueTwister).Error; err != nil {
		response.Error(c, http.StatusNotFound, "绕口令不存在")
		return
	}

	var req models.TongueTwister
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	// 更新字段
	tongueTwister.Title = req.Title
	tongueTwister.Content = req.Content
	tongueTwister.Tips = req.Tips
	tongueTwister.Level = req.Level
	tongueTwister.Order = req.Order
	tongueTwister.IsActive = req.IsActive

	if err := h.db.Save(&tongueTwister).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "更新失败")
		return
	}

	response.Success(c, tongueTwister, "更新成功")
}

// 删除绕口令（支持批量删除）
func (h *AdminHandler) DeleteTongueTwister(c *gin.Context) {
	var req struct {
		IDs []string `json:"ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误，需要提供绕口令ID列表")
		return
	}

	if err := h.db.Where("id IN ?", req.IDs).Delete(&models.TongueTwister{}).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "删除失败")
		return
	}

	response.Success(c, nil, "删除成功")
}

// 清理绕口令（删除空白和重复的）
func (h *AdminHandler) CleanTongueTwisters(c *gin.Context) {
	tx := h.db.Begin()
	if tx.Error != nil {
		response.Error(c, http.StatusInternalServerError, "清理失败")
		return
	}

	// 1. 删除空白绕口令
	// 标题为空或只包含空格，或者内容为空或只包含空格
	deleteBlankResult := tx.Where("TRIM(title) = '' OR TRIM(content) = ''").Delete(&models.TongueTwister{})
	if deleteBlankResult.Error != nil {
		tx.Rollback()
		response.Error(c, http.StatusInternalServerError, "删除空白绕口令失败: "+deleteBlankResult.Error.Error())
		return
	}

	// 2. 删除重复绕口令
	// 查找所有重复的 content，并保留每个组合中 ID 最小的一个
	var duplicateTwisters []models.TongueTwister
	// 使用子查询找到每个重复组中最小的 ID
	err := tx.Raw(`
		DELETE FROM tongue_twisters
		WHERE id IN (
			SELECT id FROM (
				SELECT
					id,
					ROW_NUMBER() OVER(PARTITION BY content ORDER BY created_at) as rn
				FROM
					tongue_twisters
			) AS sub
			WHERE sub.rn > 1
		) RETURNING *;
	`).Scan(&duplicateTwisters).Error

	if err != nil {
		tx.Rollback()
		response.Error(c, http.StatusInternalServerError, "删除重复绕口令失败: "+err.Error())
		return
	}

	tx.Commit()
	response.Success(c, gin.H{
		"deleted_blank_count": deleteBlankResult.RowsAffected,
		"deleted_duplicate_count": len(duplicateTwisters),
	}, "绕口令清理成功")
}

// ========== 每日朗诵文案管理 ==========

// 获取每日朗诵文案列表
func (h *AdminHandler) GetDailyExpressions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	offset := (page - 1) * pageSize

	var expressions []models.DailyExpression
	var total int64

	query := h.db.Model(&models.DailyExpression{})
	
	// 搜索
	if keyword := c.Query("keyword"); keyword != "" {
		query = query.Where("title LIKE ? OR content LIKE ? OR source LIKE ?", 
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 按状态筛选
	if isActive := c.Query("is_active"); isActive != "" {
		active := isActive == "true"
		query = query.Where("is_active = ?", active)
	}

	query.Count(&total)
	
	if err := query.Offset(offset).Limit(pageSize).Order("date DESC, created_at DESC").Find(&expressions).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "查询失败")
		return
	}

	response.Success(c, gin.H{
		"expressions": expressions,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
	}, "获取成功")
}

// 获取每日朗诵文案详情
func (h *AdminHandler) GetDailyExpression(c *gin.Context) {
	id := c.Param("id")
	
	var expression models.DailyExpression
	if err := h.db.Where("id = ?", id).First(&expression).Error; err != nil {
		response.Error(c, http.StatusNotFound, "文案不存在")
		return
	}

	response.Success(c, expression, "获取成功")
}

// 创建每日朗诵文案
func (h *AdminHandler) CreateDailyExpression(c *gin.Context) {
	var req models.DailyExpression
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := h.db.Create(&req).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "创建失败")
		return
	}

	response.Success(c, req, "创建成功")
}

// 更新每日朗诵文案
func (h *AdminHandler) UpdateDailyExpression(c *gin.Context) {
	id := c.Param("id")
	
	var expression models.DailyExpression
	if err := h.db.Where("id = ?", id).First(&expression).Error; err != nil {
		response.Error(c, http.StatusNotFound, "文案不存在")
		return
	}

	var req models.DailyExpression
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	// 更新字段
	expression.Title = req.Title
	expression.Content = req.Content
	expression.Tips = req.Tips
	expression.Source = req.Source
	expression.Date = req.Date
	expression.IsActive = req.IsActive

	if err := h.db.Save(&expression).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "更新失败")
		return
	}

	response.Success(c, expression, "更新成功")
}

// 删除每日朗诵文案（支持批量删除）
func (h *AdminHandler) DeleteDailyExpression(c *gin.Context) {
	var req struct {
		IDs []string `json:"ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误，需要提供文案ID列表")
		return
	}

	if err := h.db.Where("id IN ?", req.IDs).Delete(&models.DailyExpression{}).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "删除失败")
		return
	}

	response.Success(c, nil, "删除成功")
}

// 批量创建每日朗诵文案
func (h *AdminHandler) BatchCreateDailyExpressions(c *gin.Context) {
	var req []models.DailyExpression
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if len(req) == 0 {
		response.Error(c, http.StatusBadRequest, "请求体不能为空")
		return
	}

	tx := h.db.Begin()
	if tx.Error != nil {
		response.Error(c, http.StatusInternalServerError, "创建失败")
		return
	}

	for _, expression := range req {
		if err := tx.Create(&expression).Error; err != nil {
			tx.Rollback()
			response.Error(c, http.StatusInternalServerError, "批量创建失败: "+err.Error())
			return
		}
	}

	tx.Commit()
	response.Success(c, nil, "批量创建成功")
}

// TestRoute 用于测试路由是否正常工作
func (h *AdminHandler) TestRoute(c *gin.Context) {
	response.Success(c, nil, "Test route works!")
}