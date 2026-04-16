package handler

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"markdown-blog/internal/pkg/response"
)

// CaptchaHandler 验证码处理
type CaptchaHandler struct {
	captchas map[string]*Captcha // id -> captcha
}

// Captcha 验证码结构
type Captcha struct {
	ID     string
	Code   string
	Expire time.Time
}

// NewCaptchaHandler 创建验证码handler
func NewCaptchaHandler() *CaptchaHandler {
	h := &CaptchaHandler{
		captchas: make(map[string]*Captcha),
	}
	// 启动清理过期验证码的goroutine
	go h.cleanExpired()
	return h
}

// cleanExpired 清理过期的验证码
func (h *CaptchaHandler) cleanExpired() {
	ticker := time.NewTicker(5 * time.Minute)
	for range ticker.C {
		now := time.Now()
		for id, c := range h.captchas {
			if now.After(c.Expire) {
				delete(h.captchas, id)
			}
		}
	}
}

// characters 用于生成验证码的字符集
const characters = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

// Generate 生成验证码
func (h *CaptchaHandler) Generate(c *gin.Context) {
	// 生成4位验证码
	code := generateCode(4)

	// 生成唯一ID
	id := generateID()

	// 保存验证码
	h.captchas[id] = &Captcha{
		ID:     id,
		Code:   code,
		Expire: time.Now().Add(5 * time.Minute),
	}

	// 生成SVG图片
	image := generateSVG(code)

	response.Success(c, gin.H{
		"id":    id,
		"image": "data:image/svg+xml;base64," + image,
	})
}

// Validate 验证验证码
func (h *CaptchaHandler) Validate(id, code string) bool {
	c, ok := h.captchas[id]
	if !ok {
		return false
	}

	// 检查是否过期
	if time.Now().After(c.Expire) {
		delete(h.captchas, id)
		return false
	}

	// 不区分大小写
	if strings.EqualFold(c.Code, code) {
		delete(h.captchas, id) // 验证成功后删除
		return true
	}

	return false
}

// generateID 生成唯一ID
func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// generateCode 生成随机验证码
func generateCode(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	result := make([]byte, length)
	for i := range result {
		result[i] = characters[int(b[i])%len(characters)]
	}
	return string(result)
}

// generateSVG 生成SVG验证码图片
func generateSVG(code string) string {
	width := 120
	height := 48

	// 简单的SVG生成
	svg := `<svg xmlns="http://www.w3.org/2000/svg" width="` + intToStr(width) + `" height="` + intToStr(height) + `">`
	svg += `<rect fill="#f5f5f5" width="100%" height="100%"/>`

	// 绘制干扰线和点
	for i := 0; i < 8; i++ {
		x1 := int(randomInt(0, width))
		y1 := int(randomInt(0, height))
		x2 := int(randomInt(0, width))
		y2 := int(randomInt(0, height))
		svg += `<line x1="` + intToStr(x1) + `" y1="` + intToStr(y1) + `" x2="` + intToStr(x2) + `" y2="` + intToStr(y2) + `" stroke="#ddd" stroke-width="1" opacity="0.5"/>`
	}

	// 绘制文字
	fontSize := 24
	x := 15
	for _, ch := range code {
		offsetY := 30 + int(randomInt(0, 10))
		rotate := -15 + int(randomInt(0, 30))
		svg += `<text x="` + intToStr(x) + `" y="` + intToStr(offsetY) + `" font-family="Arial,Helvetica" font-size="` + intToStr(fontSize) + `" fill="#333" transform="rotate(` + intToStr(rotate) + `,` + intToStr(x+10) + `,` + intToStr(offsetY) + `)">` + string(ch) + `</text>`
		x += 25
	}

	svg += `</svg>`

	// Base64编码
	return base64.StdEncoding.EncodeToString([]byte(svg))
}

// randomInt 生成随机整数
func randomInt(min, max int) int64 {
	b := make([]byte, 8)
	rand.Read(b)
	n := int64(0)
	for _, byte := range b {
		n = n*256 + int64(byte)
	}
	return n%int64(max-min) + int64(min)
}

// intToStr 简单的int转string
func intToStr(n int) string {
	if n == 0 {
		return "0"
	}
	result := ""
	for n > 0 {
		result = string(rune('0'+n%10)) + result
		n /= 10
	}
	return result
}
