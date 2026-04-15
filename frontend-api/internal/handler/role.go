package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"markdown-blog/internal/logger"
	"markdown-blog/internal/pkg/response"
	"markdown-blog/internal/service"
)

type RoleHandler struct {
	svc *service.Service
	log *logger.Logger
}

func NewRoleHandler(svc *service.Service, log *logger.Logger) *RoleHandler {
	return &RoleHandler{
		svc: svc,
		log: log,
	}
}

// ListRoles 获取所有角色
// @Summary Get all roles
// @Description Returns a list of all roles
// @Tags roles
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/admin/roles [get]
func (h *RoleHandler) ListRoles(c *gin.Context) {
	roles, err := h.svc.ListRoles()
	if err != nil {
		h.log.Error("Failed to list roles: %v", err)
		response.InternalError(c, "Failed to load roles")
		return
	}

	response.Success(c, gin.H{
		"list": roles,
	})
}

// GetRole 获取角色
// @Summary Get role by ID
// @Description Returns a single role
// @Tags roles
// @Security BearerAuth
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {object} response.Response
// @Router /api/admin/roles/{id} [get]
func (h *RoleHandler) GetRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid role ID")
		return
	}

	role, err := h.svc.GetRole(uint(id))
	if err != nil {
		h.log.Error("Failed to get role: %v", err)
		response.InternalError(c, "Failed to load role")
		return
	}

	if role == nil {
		response.NotFound(c, "Role not found")
		return
	}

	response.Success(c, role)
}

// CreateRole 创建角色
// @Summary Create a new role
// @Description Creates a new role with permissions
// @Tags roles
// @Security BearerAuth
// @Produce json
// @Param name body string true "Role name"
// @Param permissions body []string true "Permissions"
// @Success 200 {object} response.Response
// @Router /api/admin/roles [post]
func (h *RoleHandler) CreateRole(c *gin.Context) {
	var req struct {
		Name        string   `json:"name" binding:"required"`
		Permissions []string `json:"permissions"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: name is required")
		return
	}

	role, err := h.svc.CreateRole(req.Name, req.Permissions)
	if err != nil {
		h.log.Error("Failed to create role: %v", err)
		response.InternalError(c, "Failed to create role")
		return
	}

	response.Success(c, role)
}

// UpdateRole 更新角色
// @Summary Update a role
// @Description Updates a role with new name and permissions
// @Tags roles
// @Security BearerAuth
// @Produce json
// @Param id path int true "Role ID"
// @Param name body string true "Role name"
// @Param permissions body []string true "Permissions"
// @Success 200 {object} response.Response
// @Router /api/admin/roles/{id} [put]
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid role ID")
		return
	}

	var req struct {
		Name        string   `json:"name" binding:"required"`
		Permissions []string `json:"permissions"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: name is required")
		return
	}

	role, err := h.svc.UpdateRole(uint(id), req.Name, req.Permissions)
	if err != nil {
		h.log.Error("Failed to update role: %v", err)
		response.InternalError(c, "Failed to update role")
		return
	}

	response.Success(c, role)
}

// DeleteRole 删除角色
// @Summary Delete a role
// @Description Deletes a role by ID
// @Tags roles
// @Security BearerAuth
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {object} response.Response
// @Router /api/admin/roles/{id} [delete]
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid role ID")
		return
	}

	if err := h.svc.DeleteRole(uint(id)); err != nil {
		h.log.Error("Failed to delete role: %v", err)
		response.InternalError(c, "Failed to delete role")
		return
	}

	response.Success(c, nil)
}
