package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"markdown-blog/internal/logger"
	"markdown-blog/internal/pkg/response"
	"markdown-blog/internal/service"
)

type AdminHandler struct {
	svc *service.Service
	log *logger.Logger
}

func NewAdminHandler(svc *service.Service, log *logger.Logger) *AdminHandler {
	return &AdminHandler{
		svc: svc,
		log: log,
	}
}

// ListAdminLogs 获取操作日志
// @Summary Get admin operation logs
// @Description Returns a paginated list of admin operation logs
// @Tags admin
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} response.Response
// @Router /api/admin/logs [get]
func (h *AdminHandler) ListAdminLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	logs, total, err := h.svc.ListAdminLogs(page, pageSize)
	if err != nil {
		h.log.Error("Failed to list admin logs: %v", err)
		response.InternalError(c, "Failed to load logs")
		return
	}

	response.Success(c, gin.H{
		"list":      logs,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}
