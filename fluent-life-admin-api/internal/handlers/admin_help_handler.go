package handlers

import (
	"net/http"
	"strings"

	"fluent-life-admin-api/internal/models"
	"fluent-life-admin-api/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminHelpHandler struct {
	db *gorm.DB
}

func NewAdminHelpHandler(db *gorm.DB) *AdminHelpHandler {
	return &AdminHelpHandler{db: db}
}

// -------- categories --------

// GetHelpCategories 获取帮助分类列表（管理员）
// GET /api/v1/admin/help-categories?with_articles=true
func (h *AdminHelpHandler) GetHelpCategories(c *gin.Context) {
	withArticles := strings.ToLower(c.Query("with_articles")) == "true"

	var categories []models.HelpCategory
	q := h.db.Model(&models.HelpCategory{}).Order("`order` ASC")
	if withArticles {
		q = q.Preload("Articles", func(db *gorm.DB) *gorm.DB {
			return db.Order("`order` ASC")
		})
	}

	if err := q.Find(&categories).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "获取帮助分类失败")
		return
	}
	response.Success(c, gin.H{"categories": categories, "total": len(categories)}, "获取成功")
}

// CreateHelpCategory 创建帮助分类（管理员）
// POST /api/v1/admin/help-categories
func (h *AdminHelpHandler) CreateHelpCategory(c *gin.Context) {
	var req struct {
		Name  string `json:"name" binding:"required"`
		Order int    `json:"order"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}
	name := strings.TrimSpace(req.Name)
	if name == "" {
		response.Error(c, http.StatusBadRequest, "name 不能为空")
		return
	}

	cat := models.HelpCategory{Name: name, Order: req.Order}
	if err := h.db.Create(&cat).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "创建帮助分类失败")
		return
	}
	response.Success(c, cat, "创建成功")
}

// UpdateHelpCategory 更新帮助分类（管理员）
// PUT /api/v1/admin/help-categories/:id
func (h *AdminHelpHandler) UpdateHelpCategory(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Name  *string `json:"name"`
		Order *int    `json:"order"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	var cat models.HelpCategory
	if err := h.db.Where("id = ?", id).First(&cat).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Error(c, http.StatusNotFound, "帮助分类不存在")
			return
		}
		response.Error(c, http.StatusInternalServerError, "获取帮助分类失败")
		return
	}

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			response.Error(c, http.StatusBadRequest, "name 不能为空")
			return
		}
		cat.Name = name
	}
	if req.Order != nil {
		cat.Order = *req.Order
	}

	if err := h.db.Save(&cat).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "更新帮助分类失败")
		return
	}
	response.Success(c, cat, "更新成功")
}

// DeleteHelpCategory 删除帮助分类（管理员）
// DELETE /api/v1/admin/help-categories/:id
func (h *AdminHelpHandler) DeleteHelpCategory(c *gin.Context) {
	id := c.Param("id")

	// 先把该分类下文章的 category_id 置空/阻止删除？这里选择级联删除文章（简单）
	h.db.Where("category_id = ?", id).Delete(&models.HelpArticle{})

	if err := h.db.Delete(&models.HelpCategory{}, "id = ?", id).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "删除帮助分类失败")
		return
	}
	response.Success(c, nil, "删除成功")
}

// -------- articles --------

// GetHelpArticles 获取帮助文章列表（管理员）
// GET /api/v1/admin/help-articles?category_id=...&q=...
func (h *AdminHelpHandler) GetHelpArticles(c *gin.Context) {
	categoryID := strings.TrimSpace(c.Query("category_id"))
	search := strings.TrimSpace(c.Query("q"))

	dbq := h.db.Model(&models.HelpArticle{}).Order("`order` ASC, created_at DESC")
	if categoryID != "" {
		dbq = dbq.Where("category_id = ?", categoryID)
	}
	if search != "" {
		like := "%" + search + "%"
		dbq = dbq.Where("question ILIKE ? OR answer ILIKE ?", like, like)
	}

	var articles []models.HelpArticle
	if err := dbq.Find(&articles).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "获取帮助文章失败")
		return
	}
	response.Success(c, gin.H{"articles": articles, "total": len(articles)}, "获取成功")
}

// CreateHelpArticle 创建帮助文章（管理员）
// POST /api/v1/admin/help-articles
func (h *AdminHelpHandler) CreateHelpArticle(c *gin.Context) {
	var req struct {
		CategoryID string `json:"category_id" binding:"required"`
		Question   string `json:"question" binding:"required"`
		Answer     string `json:"answer" binding:"required"`
		Order      int    `json:"order"`
		IsActive   *bool  `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}
	question := strings.TrimSpace(req.Question)
	answer := strings.TrimSpace(req.Answer)
	if question == "" || answer == "" {
		response.Error(c, http.StatusBadRequest, "question/answer 不能为空")
		return
	}

	// 确保分类存在
	var cat models.HelpCategory
	if err := h.db.Where("id = ?", req.CategoryID).First(&cat).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Error(c, http.StatusBadRequest, "category_id 无效")
			return
		}
		response.Error(c, http.StatusInternalServerError, "创建帮助文章失败")
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	article := models.HelpArticle{
		CategoryID: cat.ID,
		Question:   question,
		Answer:     answer,
		Order:      req.Order,
		IsActive:   isActive,
	}

	if err := h.db.Create(&article).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "创建帮助文章失败")
		return
	}
	response.Success(c, article, "创建成功")
}

// UpdateHelpArticle 更新帮助文章（管理员）
// PUT /api/v1/admin/help-articles/:id
func (h *AdminHelpHandler) UpdateHelpArticle(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		CategoryID *string `json:"category_id"`
		Question   *string `json:"question"`
		Answer     *string `json:"answer"`
		Order      *int    `json:"order"`
		IsActive   *bool   `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	var article models.HelpArticle
	if err := h.db.Where("id = ?", id).First(&article).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Error(c, http.StatusNotFound, "帮助文章不存在")
			return
		}
		response.Error(c, http.StatusInternalServerError, "获取帮助文章失败")
		return
	}

	if req.CategoryID != nil {
		var cat models.HelpCategory
		if err := h.db.Where("id = ?", *req.CategoryID).First(&cat).Error; err != nil {
			response.Error(c, http.StatusBadRequest, "category_id 无效")
			return
		}
		article.CategoryID = cat.ID
	}
	if req.Question != nil {
		q := strings.TrimSpace(*req.Question)
		if q == "" {
			response.Error(c, http.StatusBadRequest, "question 不能为空")
			return
		}
		article.Question = q
	}
	if req.Answer != nil {
		a := strings.TrimSpace(*req.Answer)
		if a == "" {
			response.Error(c, http.StatusBadRequest, "answer 不能为空")
			return
		}
		article.Answer = a
	}
	if req.Order != nil {
		article.Order = *req.Order
	}
	if req.IsActive != nil {
		article.IsActive = *req.IsActive
	}

	if err := h.db.Save(&article).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "更新帮助文章失败")
		return
	}
	response.Success(c, article, "更新成功")
}

// DeleteHelpArticle 删除帮助文章（管理员）
// DELETE /api/v1/admin/help-articles/:id
func (h *AdminHelpHandler) DeleteHelpArticle(c *gin.Context) {
	id := c.Param("id")

	if err := h.db.Delete(&models.HelpArticle{}, "id = ?", id).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "删除帮助文章失败")
		return
	}
	response.Success(c, nil, "删除成功")
}
