package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"markdown-blog/internal/logger"
	"markdown-blog/internal/pkg/response"
	"markdown-blog/internal/service"
)

type TagHandler struct {
	svc *service.Service
	log *logger.Logger
}

func NewTagHandler(svc *service.Service, log *logger.Logger) *TagHandler {
	return &TagHandler{
		svc: svc,
		log: log,
	}
}

// ListTags 标签列表（分页）
func (h *TagHandler) ListTags(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	tags, total, err := h.svc.ListTags(page, pageSize)
	if err != nil {
		h.log.Error("Failed to list tags: %v", err)
		response.InternalError(c, "Failed to load tags")
		return
	}

	response.Success(c, gin.H{
		"list":      tags,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// ListAllTags 全部标签（下拉框用）
func (h *TagHandler) ListAllTags(c *gin.Context) {
	tags, err := h.svc.ListAllTags()
	if err != nil {
		h.log.Error("Failed to list all tags: %v", err)
		response.InternalError(c, "Failed to load tags")
		return
	}

	response.Success(c, tags)
}

// GetTag 标签详情
func (h *TagHandler) GetTag(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid tag ID")
		return
	}

	tag, err := h.svc.GetTagByID(uint(id))
	if err != nil {
		response.NotFound(c, "Tag not found")
		return
	}

	response.Success(c, tag)
}

// CreateTag 创建标签
func (h *TagHandler) CreateTag(c *gin.Context) {
	var req struct {
		Name  string `json:"name" binding:"required"`
		Slug  string `json:"slug" binding:"required"`
		Color string `json:"color"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	tag, err := h.svc.CreateTag(req.Name, req.Slug, req.Color)
	if err != nil {
		h.log.Error("Failed to create tag: %v", err)
		response.InternalError(c, "Failed to create tag")
		return
	}

	response.Success(c, tag)
}

// UpdateTag 更新标签
func (h *TagHandler) UpdateTag(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid tag ID")
		return
	}

	var req struct {
		Name  string `json:"name" binding:"required"`
		Slug  string `json:"slug" binding:"required"`
		Color string `json:"color"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	tag, err := h.svc.UpdateTag(uint(id), req.Name, req.Slug, req.Color)
	if err != nil {
		h.log.Error("Failed to update tag: %v", err)
		response.InternalError(c, "Failed to update tag")
		return
	}

	response.Success(c, tag)
}

// DeleteTag 删除标签
func (h *TagHandler) DeleteTag(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid tag ID")
		return
	}

	if err := h.svc.DeleteTag(uint(id)); err != nil {
		h.log.Error("Failed to delete tag: %v", err)
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "Tag deleted"})
}
