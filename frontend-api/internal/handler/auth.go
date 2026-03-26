package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v50/github"

	"markdown-blog/internal/config"
	"markdown-blog/internal/logger"
	"markdown-blog/internal/middleware"
	"markdown-blog/internal/model"
	"markdown-blog/internal/pkg/jwt"
	"markdown-blog/internal/pkg/response"
	"markdown-blog/internal/service"
)

type AuthHandler struct {
	svc  *service.Service
	cfg  *config.Config
	log  *logger.Logger
	jwt  *jwt.JWT
}

func NewAuthHandler(svc *service.Service, cfg *config.Config, l *logger.Logger, jwtUtil *jwt.JWT) *AuthHandler {
	return &AuthHandler{
		svc:  svc,
		cfg:  cfg,
		log:  l,
		jwt:  jwtUtil,
	}
}

// LoginInfo returns the GitHub OAuth login page info
// @Summary Get login page info
// @Description Returns GitHub OAuth configuration for login
// @Tags auth
// @Produce json
// @Success 200 {object} response.Response{data=LoginInfo}
// @Router /api/v1/auth/login [get]
func (h *AuthHandler) LoginInfo(c *gin.Context) {
	state, _ := generateState()
	response.Success(c, gin.H{
		"client_id":    h.cfg.Github.ClientID,
		"callback_url": h.cfg.Github.CallbackURL,
		"state":        state,
	})
}

// Callback handles GitHub OAuth callback
// @Summary GitHub OAuth callback
// @Description Handles the OAuth callback from GitHub
// @Tags auth
// @Accept json
// @Produce json
// @Param code query string true "OAuth code"
// @Param state query string false "OAuth state"
// @Success 200 {object} response.Response{data=TokenResponse}
// @Router /api/v1/auth/callback [post]
func (h *AuthHandler) Callback(c *gin.Context) {
	h.log.Info("[AUTH] AuthCallback called, code: %s", c.Query("code"))

	code := c.Query("code")
	if code == "" {
		h.log.Warn("[AUTH] No code provided")
		response.BadRequest(c, "Missing code")
		return
	}

	token, err := h.exchangeCode(code)
	if err != nil {
		h.log.Error("[AUTH] Failed to exchange code: %v", err)
		response.InternalError(c, "Failed to get token")
		return
	}

	client := github.NewTokenClient(c.Request.Context(), token)
	user, _, err := client.Users.Get(c.Request.Context(), "")
	if err != nil {
		h.log.Error("[AUTH] Failed to get user: %v", err)
		response.InternalError(c, "Failed to get user info")
		return
	}

	emails, _, _ := client.Users.ListEmails(c.Request.Context(), nil)
	var email string
	if len(emails) > 0 {
		for _, e := range emails {
			if e.GetPrimary() && e.GetVerified() {
				email = e.GetEmail()
				break
			}
		}
	}
	if email == "" {
		email = user.GetEmail()
	}

	dbUser := &model.User{
		GitHubID:  user.GetID(),
		Login:     user.GetLogin(),
		Name:      user.GetName(),
		AvatarURL: user.GetAvatarURL(),
		Email:     email,
	}

	if dbUser.Name == "" {
		dbUser.Name = dbUser.Login
	}

	if err := h.svc.SaveUser(dbUser); err != nil {
		h.log.Warn("[AUTH] Failed to save user to DB: %v", err)
	}

	// Generate JWT token
	jwtToken, _, err := h.jwt.GenerateToken(dbUser.GitHubID, dbUser.Login, dbUser.Name, dbUser.AvatarURL, dbUser.Email)
	if err != nil {
		h.log.Error("[AUTH] Failed to generate JWT: %v", err)
		response.InternalError(c, "Failed to generate token")
		return
	}

	h.log.Info("[AUTH] User logged in successfully: %s (ID: %d)", dbUser.Login, dbUser.GitHubID)

	// Redirect to frontend login page with token
	// Frontend will store token and redirect to home
	c.Redirect(302, "https://jumoshen.cn/login?token="+jwtToken)
}

// Logout handles user logout
// @Summary Logout
// @Description Logout current user and invalidate token
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	claims := middleware.GetUserClaims(c)
	if claims == nil {
		response.Unauthorized(c, "not authenticated")
		return
	}

	// Add token to blacklist
	ctx := context.Background()
	exp := claims.ExpiresAt.Time
	jti := claims.ID

	if err := h.svc.Logout(ctx, jti, exp); err != nil {
		h.log.Warn("[AUTH] Failed to add token to blacklist: %v", err)
	}

	h.log.Info("[AUTH] User logged out: %s", claims.Login)
	response.Success(c, nil)
}

// Me returns current user info
// @Summary Get current user
// @Description Returns the currently authenticated user
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.Response{data=UserInfo}
// @Router /api/v1/auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	claims := middleware.GetUserClaims(c)
	if claims == nil {
		response.Unauthorized(c, "not authenticated")
		return
	}

	response.Success(c, gin.H{
		"id":         claims.UserID,
		"login":      claims.Login,
		"name":       claims.Name,
		"avatar_url": claims.AvatarURL,
		"email":      claims.Email,
	})
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Error       string `json:"error"`
}

func (h *AuthHandler) exchangeCode(code string) (string, error) {
	params := url.Values{}
	params.Set("client_id", h.cfg.Github.ClientID)
	params.Set("client_secret", h.cfg.Github.ClientSecret)
	params.Set("code", code)

	req, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token", strings.NewReader(params.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result tokenResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if result.Error != "" {
		return "", fmt.Errorf("oauth error: %s", result.Error)
	}

	return result.AccessToken, nil
}

func generateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}

type UserInfo struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
	Email     string `json:"email"`
}

type TokenResponse struct {
	Token     string `json:"token"`
	TokenType string `json:"token_type"`
	ExpiresIn int64  `json:"expires_in"`
	JTI       string `json:"jti"`
	User      UserInfo `json:"user"`
}

type LoginInfo struct {
	ClientID    string `json:"client_id"`
	CallbackURL string `json:"callback_url"`
	State       string `json:"state"`
}
