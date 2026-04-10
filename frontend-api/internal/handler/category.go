package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"markdown-blog/internal/logger"
	"markdown-blog/internal/pkg/response"
	"markdown-blog/internal/service"
)

type CategoryHandler struct {
	svc *service.Service
	log *logger.Logger
}

func NewCategoryHandler(svc *service.Service, log *logger.Logger) *CategoryHandler {
	return &CategoryHandler{
		svc: svc,
		log: log,
	}
}

// ListCategories 分类列表（分页）
func (h *CategoryHandler) ListCategories(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	categories, total, err := h.svc.ListCategories(page, pageSize)
	if err != nil {
		h.log.Error("Failed to list categories: %v", err)
		response.InternalError(c, "Failed to load categories")
		return
	}

	response.Success(c, gin.H{
		"list":      categories,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// ListAllCategories 全部分类（下拉框用）
func (h *CategoryHandler) ListAllCategories(c *gin.Context) {
	categories, err := h.svc.ListAllCategories()
	if err != nil {
		h.log.Error("Failed to list all categories: %v", err)
		response.InternalError(c, "Failed to load categories")
		return
	}

	response.Success(c, categories)
}

// GetCategory 分类详情
func (h *CategoryHandler) GetCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid category ID")
		return
	}

	category, err := h.svc.GetCategoryByID(uint(id))
	if err != nil {
		response.NotFound(c, "Category not found")
		return
	}

	response.Success(c, category)
}

// CreateCategory 创建分类
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Slug        string `json:"slug" binding:"required"`
		Description string `json:"description"`
		Sort        int    `json:"sort"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	category, err := h.svc.CreateCategory(req.Name, req.Slug, req.Description, req.Sort)
	if err != nil {
		h.log.Error("Failed to create category: %v", err)
		response.InternalError(c, "Failed to create category")
		return
	}

	response.Success(c, category)
}

// UpdateCategory 更新分类
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid category ID")
		return
	}

	var req struct {
		Name        string `json:"name" binding:"required"`
		Slug        string `json:"slug" binding:"required"`
		Description string `json:"description"`
		Sort        int    `json:"sort"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	category, err := h.svc.UpdateCategory(uint(id), req.Name, req.Slug, req.Description, req.Sort)
	if err != nil {
		h.log.Error("Failed to update category: %v", err)
		response.InternalError(c, "Failed to update category")
		return
	}

	response.Success(c, category)
}

// DeleteCategory 删除分类
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid category ID")
		return
	}

	if err := h.svc.DeleteCategory(uint(id)); err != nil {
		h.log.Error("Failed to delete category: %v", err)
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "Category deleted"})
}
