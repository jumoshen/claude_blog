package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"markdown-blog/internal/logger"
	"markdown-blog/internal/middleware"
	"markdown-blog/internal/model"
	"markdown-blog/internal/pkg/response"
	"markdown-blog/internal/service"
)

type CommentHandler struct {
	svc        *service.Service
	log        *logger.Logger
	upgrader   websocket.Upgrader
	clients    map[string]map[*websocket.Conn]bool // postSlug -> connections
	register   chan *Client
	unregister chan *Client
	broadcast  chan *BroadcastMessage
	mu         sync.RWMutex
}

type Client struct {
	conn    *websocket.Conn
	postSlug string
}

type BroadcastMessage struct {
	PostSlug string
	Data     []byte
}

// WSRateLimit WS频率限制
type WSRateLimit struct {
	mu         sync.RWMutex
	lastAccess map[string]time.Time
	interval   time.Duration
}

func NewWSRateLimit(interval time.Duration) *WSRateLimit {
	return &WSRateLimit{
		lastAccess: make(map[string]time.Time),
		interval:   interval,
	}
}

func (r *WSRateLimit) Allow(ip string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	if last, ok := r.lastAccess[ip]; ok && now.Sub(last) < r.interval {
		return false
	}
	r.lastAccess[ip] = now
	return true
}

// Global WS rate limiter - 同一IP每分钟最多连接30次
var wsRateLimiter = NewWSRateLimit(time.Minute)

func NewCommentHandler(svc *service.Service, log *logger.Logger) *CommentHandler {
	h := &CommentHandler{
		svc:     svc,
		log:     log,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		clients:    make(map[string]map[*websocket.Conn]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *BroadcastMessage, 256),
	}

	go h.run()
	return h
}

func (h *CommentHandler) run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if h.clients[client.postSlug] == nil {
				h.clients[client.postSlug] = make(map[*websocket.Conn]bool)
			}
			h.clients[client.postSlug][client.conn] = true
			h.mu.Unlock()
			h.log.Info("Client subscribed to post: %s", client.postSlug)

		case client := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.clients[client.postSlug]; ok {
				if _, exists := clients[client.conn]; exists {
					delete(clients, client.conn)
					client.conn.Close()
					if len(clients) == 0 {
						delete(h.clients, client.postSlug)
					}
				}
			}
			h.mu.Unlock()
			h.log.Info("Client unsubscribed from post: %s", client.postSlug)

		case msg := <-h.broadcast:
			h.mu.RLock()
			clients := h.clients[msg.PostSlug]
			for conn := range clients {
				err := conn.WriteMessage(websocket.TextMessage, msg.Data)
				if err != nil {
					h.log.Error("Failed to broadcast to client: %v", err)
					conn.Close()
					delete(clients, conn)
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *CommentHandler) broadcastToPost(postSlug string, data []byte) {
	h.broadcast <- &BroadcastMessage{
		PostSlug: postSlug,
		Data:     data,
	}
}

// GetComments 获取文章评论
func (h *CommentHandler) GetComments(c *gin.Context) {
	slug := c.Param("postSlug")
	if slug == "" {
		response.BadRequest(c, "missing post slug")
		return
	}

	comments, err := h.svc.GetCommentsByPostSlug(slug, 50)
	if err != nil {
		h.log.Error("Failed to get comments: %v", err)
		response.InternalError(c, "Failed to load comments")
		return
	}

	// 转换为API响应格式
	type CommentResponse struct {
		ID        uint      `json:"id"`
		Nickname  string    `json:"nickname"`
		Content   string    `json:"content"`
		CreatedAt time.Time `json:"created_at"`
	}

	result := make([]CommentResponse, 0, len(comments))
	for _, c := range comments {
		result = append(result, CommentResponse{
			ID:        c.ID,
			Nickname:  c.Nickname,
			Content:   c.Content,
			CreatedAt: c.CreatedAt,
		})
	}

	response.Success(c, result)
}

// CreateCommentRequest 创建评论请求
type CreateCommentRequest struct {
	PostSlug  string `json:"post_slug" binding:"required"`
	Nickname  string `json:"nickname" binding:"required,max=50"`
	Content   string `json:"content" binding:"required"`
	DeviceID  string `json:"device_id"`
}

// CreateComment 创建评论
func (h *CommentHandler) CreateComment(c *gin.Context) {
	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	if req.PostSlug == "" {
		response.BadRequest(c, "missing post_slug")
		return
	}

	claims := middleware.GetUserClaims(c)
	userID := int64(0)
	if claims != nil {
		userID = claims.UserID
	}

	// 匿名用户检查频率限制
	if userID == 0 {
		ip := c.ClientIP()
		deviceID := req.DeviceID

		allowed, err := h.svc.CheckCommentRateLimit(c.Request.Context(), ip, deviceID, userID)
		if err != nil {
			h.log.Error("Failed to check rate limit: %v", err)
			// 允许继续，不因为Redis错误阻止评论
		} else if !allowed {
			response.Error(c, 429, "评论太频繁，请稍后再试")
			return
		}
	}

	// 敏感词检测（昵称和内容）
	if h.svc.ContainsSensitiveWords(req.Nickname) || h.svc.ContainsSensitiveWords(req.Content) {
		response.Error(c, 400, "评论包含不当内容，请修改后重试")
		return
	}

	comment := &model.Comment{
		PostSlug:  req.PostSlug,
		UserID:    userID,
		Nickname:  req.Nickname,
		Content:   req.Content,
		IP:        c.ClientIP(),
		DeviceID:  req.DeviceID,
		UserAgent: c.Request.UserAgent(),
		Status:    1,
	}

	if err := h.svc.CreateComment(comment); err != nil {
		h.log.Error("Failed to create comment: %v", err)
		response.InternalError(c, "Failed to create comment")
		return
	}

	// 广播新评论到该文章的WebSocket客户端
	msg := map[string]interface{}{
		"type": "new_comment",
		"data": map[string]interface{}{
			"id":         comment.ID,
			"nickname":   comment.Nickname,
			"content":     comment.Content,
			"created_at": comment.CreatedAt,
		},
	}
	msgBytes, _ := json.Marshal(msg)
	h.broadcastToPost(req.PostSlug, msgBytes)

	response.Created(c, map[string]interface{}{
		"id":         comment.ID,
		"nickname":   comment.Nickname,
		"content":    comment.Content,
		"created_at": comment.CreatedAt,
	})
}

// HandleWebSocket WebSocket处理
func (h *CommentHandler) HandleWebSocket(c *gin.Context) {
	postSlug := c.Param("postSlug")
	if postSlug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing post slug"})
		return
	}

	// WS连接频率限制
	ip := c.ClientIP()
	if !wsRateLimiter.Allow(ip) {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "连接太频繁"})
		return
	}

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.log.Error("WebSocket upgrade failed: %v", err)
		return
	}

	client := &Client{
		conn:    conn,
		postSlug: postSlug,
	}

	h.register <- client

	// 处理客户端消息
	go func() {
		defer func() {
			h.unregister <- client
		}()

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					h.log.Error("WebSocket read error: %v", err)
				}
				break
			}

			// 解析客户端消息
			var msg map[string]interface{}
			if err := json.Unmarshal(message, &msg); err != nil {
				continue
			}

			// 处理订阅消息
			if msgType, ok := msg["type"].(string); ok && msgType == "subscribe" {
				if ps, ok := msg["post_slug"].(string); ok && ps != "" {
					client.postSlug = ps
					h.unregister <- client
					client = &Client{conn: conn, postSlug: ps}
					h.register <- client
				}
			}
		}
	}()
}

// CommentHub 管理所有评论WebSocket连接
type CommentHub struct {
	rooms map[string]map[*CommentClient]bool
	mu    sync.RWMutex
}

type CommentClient struct {
	conn    *websocket.Conn
	postSlug string
	send    chan []byte
}

func NewCommentHub() *CommentHub {
	return &CommentHub{
		rooms: make(map[string]map[*CommentClient]bool),
	}
}

func (h *CommentHub) Register(client *CommentClient, postSlug string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.rooms[postSlug] == nil {
		h.rooms[postSlug] = make(map[*CommentClient]bool)
	}
	h.rooms[postSlug][client] = true
}

func (h *CommentHub) Unregister(client *CommentClient, postSlug string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if clients, ok := h.rooms[postSlug]; ok {
		delete(clients, client)
		if len(clients) == 0 {
			delete(h.rooms, postSlug)
		}
	}
}

func (h *CommentHub) Broadcast(postSlug string, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if clients, ok := h.rooms[postSlug]; ok {
		for client := range clients {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(clients, client)
			}
		}
	}
}

// ServeWsComment 处理评论WebSocket连接
func ServeWsComment(hub *CommentHub, w http.ResponseWriter, r *http.Request, postSlug string) {
	conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &CommentClient{
		conn:    conn,
		postSlug: postSlug,
		send:    make(chan []byte, 256),
	}

	hub.Register(client, postSlug)

	go client.writePump()
	go client.readPump(hub)
}

func (c *CommentClient) writePump() {
	defer func() {
		c.conn.Close()
	}()

	for message := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			return
		}
	}
}

func (c *CommentClient) readPump(hub *CommentHub) {
	defer func() {
		hub.Unregister(c, c.postSlug)
		c.conn.Close()
	}()

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

// CommentHubManager 全局评论Hub管理器
var CommentHubInstance *CommentHub
var once sync.Once

func GetCommentHub() *CommentHub {
	once.Do(func() {
		CommentHubInstance = NewCommentHub()
	})
	return CommentHubInstance
}
