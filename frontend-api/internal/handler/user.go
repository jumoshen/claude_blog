package handler

import (
	"github.com/gin-gonic/gin"

	"markdown-blog/internal/logger"
	"markdown-blog/internal/pkg/response"
	"markdown-blog/internal/service"
)

type UserHandler struct {
	svc *service.Service
	log *logger.Logger
}

func NewUserHandler(svc *service.Service, log *logger.Logger) *UserHandler {
	return &UserHandler{
		svc: svc,
		log: log,
	}
}

// Register 用户注册
// @Summary Register a new user
// @Description Register a new blog user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Register request"
// @Success 200 {object} response.Response
// @Router /api/v1/auth/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required,min=3,max=50"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	user, err := h.svc.RegisterBlogUser(req.Username, req.Email, req.Password)
	if err != nil {
		h.log.Error("Failed to register user: %v", err)
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, user)
}

// Login 用户登录
// @Summary Login
// @Description Login with username, password and captcha
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login request"
// @Success 200 {object} response.Response
// @Router /api/v1/auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req struct {
		Username  string `json:"username" binding:"required"`
		Password  string `json:"password" binding:"required"`
		Captcha   string `json:"captcha" binding:"required"`
		CaptchaID string `json:"captcha_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	// 验证验证码
	captchaHandler := c.MustGet("captchaHandler").(*CaptchaHandler)
	if !captchaHandler.Validate(req.CaptchaID, req.Captcha) {
		response.BadRequest(c, "验证码错误")
		return
	}

	user, token, err := h.svc.LoginBlogUser(req.Username, req.Password)
	if err != nil {
		h.log.Error("Failed to login: %v", err)
		response.Unauthorized(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"user":  user,
		"token": token,
	})
}

// Me 获取当前用户信息
// @Summary Get current user
// @Description Get current logged in user info
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/v1/auth/me [get]
func (h *UserHandler) Me(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "Not logged in")
		return
	}

	user, err := h.svc.GetBlogUserByID(userID.(uint))
	if err != nil {
		response.NotFound(c, "User not found")
		return
	}

	response.Success(c, user)
}

// Request types for documentation
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
