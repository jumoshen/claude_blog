package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"markdown-blog/internal/config"
	"markdown-blog/internal/handler"
	"markdown-blog/internal/logger"
	"markdown-blog/internal/middleware"
	"markdown-blog/internal/pkg/jwt"
	"markdown-blog/internal/pkg/response"
	"markdown-blog/internal/repository"
	"markdown-blog/internal/service"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	l := logger.New(&logger.Config{
		Level:      cfg.Log.Level,
		Path:       cfg.Log.Path,
		Filename:   cfg.Log.Filename,
		MaxSize:    cfg.Log.MaxSize,
		MaxBackups: cfg.Log.MaxBackups,
		MaxAge:     cfg.Log.MaxAge,
		Compress:   cfg.Log.Compress,
		Console:    cfg.Log.Console,
	})
	defer l.Close()

	l.Info("Starting blog server...")

	// Initialize repository
	repo, err := repository.New(cfg)
	if err != nil {
		l.Fatal("Failed to create repository: %v", err)
	}

	// Initialize service
	svc := service.New(cfg, repo)

	// Parse content directory
	if err := svc.ParseContentDir(); err != nil {
		l.Warn("Failed to parse content directory: %v", err)
	}

	// Initialize JWT
	jwtUtil := jwt.NewWithConfig(
		cfg.JWT.Secret,
		cfg.JWT.Expiration,
		cfg.JWT.Issuer,
		repo.GetRedis(),
	)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(svc, cfg, l, jwtUtil)
	postHandler := handler.NewPostHandler(svc, l)
	commentHandler := handler.NewCommentHandler(svc, l)
	categoryHandler := handler.NewCategoryHandler(svc, l)
	tagHandler := handler.NewTagHandler(svc, l)
	userHandler := handler.NewUserHandler(svc, l)
	adminHandler := handler.NewAdminHandler(svc, l)
	roleHandler := handler.NewRoleHandler(svc, l)

	// Setup Gin
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(requestLogger(l))
	r.Use(cors())

	// Static files
	r.Static("/static", "./static")

	// API routes
	v1 := r.Group("/api/v1")
	{
		// Site info
		v1.GET("/site", func(c *gin.Context) {
			response.Success(c, gin.H{
				"title":       cfg.Site.Title,
				"author":      cfg.Site.Author,
				"description": cfg.Site.Description,
				"beian":       cfg.Site.Beian,
				"keywords":    cfg.Site.Keywords,
			})
		})

		// Public routes with visit logging
		v1.GET("/posts", postHandler.ListPosts)
		v1.GET("/posts/search", postHandler.SearchPosts)
		v1.GET("/posts/featured", postHandler.ListFeaturedPosts)
		v1.GET("/posts/popular", postHandler.ListPopularPosts)
		v1.GET("/posts/:slug/toc", postHandler.GetTOC)
		v1.GET("/posts/:slug/navigation", postHandler.GetNavigation)
		v1.GET("/posts/:slug/related", postHandler.ListRelatedPosts)
		v1.GET("/posts/:slug/likes", postHandler.GetPostLikes)
		v1.POST("/posts/:slug/like", middleware.AuthMiddleware(jwtUtil, repo), postHandler.LikePost)
		v1.POST("/posts/:slug/favorite", middleware.AuthMiddleware(jwtUtil, repo), postHandler.FavoritePost)
		v1.GET("/posts/:slug/check", postHandler.CheckPassword)
		v1.POST("/posts/:slug/verify", postHandler.VerifyPassword)
		v1.GET("/posts/:slug", postHandler.GetPost)

		// User routes
		users := v1.Group("/users")
		users.Use(middleware.AuthMiddleware(jwtUtil, repo))
		{
			users.GET("/me/favorites", postHandler.ListMyFavorites)
		}
		v1.GET("/archives", postHandler.GetArchives)
		v1.GET("/about", postHandler.GetAbout)
		v1.GET("/tags", postHandler.GetTags)
		v1.GET("/categories", postHandler.GetCategories)
		v1.GET("/sitemap.xml", postHandler.GetSitemap)
		v1.GET("/feed.xml", postHandler.GetRSS)

		// Auth routes
		auth := v1.Group("/auth")
		{
			auth.GET("/login", authHandler.LoginInfo)
			auth.GET("/callback", authHandler.Callback)
			auth.POST("/logout", middleware.AuthMiddleware(jwtUtil, repo), authHandler.Logout)
			auth.GET("/me", middleware.AuthMiddleware(jwtUtil, repo), authHandler.Me)

			// Blog user auth routes (C端)
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
		}

		// Blog user routes (C端)
		blogUsers := v1.Group("/blogusers")
		{
			blogUsers.GET("/me", middleware.AuthMiddleware(jwtUtil, repo), userHandler.Me)
		}

		// Comment routes
		comments := v1.Group("/comments")
		{
			comments.GET("/:postSlug", commentHandler.GetComments)
			comments.POST("", commentHandler.CreateComment)
		}

		// WebSocket route for comments
		r.GET("/ws/comments/:postSlug", commentHandler.HandleWebSocket)

		// Admin routes (protected)
		admin := v1.Group("/admin")
		admin.Use(middleware.AuthMiddleware(jwtUtil, repo))
		{
			admin.POST("/refresh", postHandler.Refresh)

			// Category routes
			admin.GET("/categories", categoryHandler.ListCategories)
			admin.GET("/categories/all", categoryHandler.ListAllCategories)
			admin.GET("/categories/:id", categoryHandler.GetCategory)
			admin.POST("/categories", categoryHandler.CreateCategory)
			admin.PUT("/categories/:id", categoryHandler.UpdateCategory)
			admin.DELETE("/categories/:id", categoryHandler.DeleteCategory)

			// Tag routes
			admin.GET("/tags", tagHandler.ListTags)
			admin.GET("/tags/all", tagHandler.ListAllTags)
			admin.GET("/tags/:id", tagHandler.GetTag)
			admin.POST("/tags", tagHandler.CreateTag)
			admin.PUT("/tags/:id", tagHandler.UpdateTag)
			admin.DELETE("/tags/:id", tagHandler.DeleteTag)

			// Post pin/feature routes
			admin.PUT("/posts/:slug/pin", postHandler.SetPostPin)
			admin.PUT("/posts/:slug/feature", postHandler.SetPostFeature)
			admin.POST("/posts/:slug/schedule", postHandler.SchedulePost)
			admin.POST("/posts/:slug/password", postHandler.SetPassword)

			// Admin log routes
			admin.GET("/logs", adminHandler.ListAdminLogs)

			// Role routes
			admin.GET("/roles", roleHandler.ListRoles)
			admin.GET("/roles/:id", roleHandler.GetRole)
			admin.POST("/roles", roleHandler.CreateRole)
			admin.PUT("/roles/:id", roleHandler.UpdateRole)
			admin.DELETE("/roles/:id", roleHandler.DeleteRole)
		}
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Start server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// Graceful shutdown
	go func() {
		l.Info("Server starting on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			l.Fatal("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	l.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		l.Fatal("Server forced to shutdown: %v", err)
	}

	l.Info("Server exited")
}

// requestLogger request logging middleware
func requestLogger(l *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		clientIP := c.ClientIP()

		l.Info("[HTTP] %s %s %d %v %s",
			method,
			path,
			status,
			latency,
			clientIP,
		)
	}
}

// cors middleware
func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
