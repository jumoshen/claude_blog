package handler

import (
	"html/template"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"markdown-blog/internal/logger"
	"markdown-blog/internal/middleware"
	"markdown-blog/internal/pkg/response"
	"markdown-blog/internal/service"
)

type PostHandler struct {
	svc *service.Service
	cfg interface{}
	log *logger.Logger
}

func NewPostHandler(svc *service.Service, log *logger.Logger) *PostHandler {
	return &PostHandler{
		svc: svc,
		log: log,
	}
}

// ListPosts returns posts with pagination
// @Summary Get posts with pagination
// @Description Returns a paginated list of blog posts
// @Tags posts
// @Produce json
// @Param tag query string false "Filter by tag"
// @Param category query string false "Filter by category"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} response.Response{data=[]PostInfo}
// @Router /api/v1/posts [get]
func (h *PostHandler) ListPosts(c *gin.Context) {
	filterTag := c.Query("tag")
	filterCategory := c.Query("category")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	posts, total, err := h.svc.ListPostsPaginated(filterTag, filterCategory, page, pageSize)
	if err != nil {
		h.log.Error("Failed to list posts: %v", err)
		response.InternalError(c, "Failed to load posts")
		return
	}

	h.log.Info("ListPosts: page=%d, page_size=%d, total=%d, tag=%s, category=%s", page, pageSize, total, filterTag, filterCategory)
	response.Success(c, gin.H{
		"list":      posts,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetPost returns a single post by slug
// @Summary Get a post
// @Description Returns a single blog post by slug
// @Tags posts
// @Produce json
// @Param slug path string true "Post slug"
// @Success 200 {object} response.Response{data=PostDetail}
// @Router /api/v1/posts/{slug} [get]
func (h *PostHandler) GetPost(c *gin.Context) {
	slug := c.Param("slug")
	info, content, err := h.svc.GetPost(slug)
	if err != nil {
		response.NotFound(c, "Post not found")
		return
	}

	// 记录访问
	h.svc.RecordVisit(slug, c.ClientIP(), c.Request.UserAgent(), 0)

	response.Success(c, gin.H{
		"post":    info,
		"content": template.HTML(content),
	})
}

// GetArchives returns posts grouped by year
// @Summary Get archives
// @Description Returns all posts grouped by year
// @Tags posts
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/v1/archives [get]
func (h *PostHandler) GetArchives(c *gin.Context) {
	archives, err := h.svc.GetArchives()
	if err != nil {
		h.log.Error("Failed to get archives: %v", err)
		response.InternalError(c, "Failed to load archives")
		return
	}

	response.Success(c, archives)
}

// GetAbout returns the about page content
// @Summary Get about page
// @Description Returns the about page content
// @Tags posts
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/v1/about [get]
func (h *PostHandler) GetAbout(c *gin.Context) {
	content, err := h.svc.GetAboutContent()
	if err != nil {
		content = "<p>About page not found.</p>"
	}

	response.Success(c, gin.H{
		"content": template.HTML(content),
	})
}

// Refresh refreshes the content from disk
// @Summary Refresh content
// @Description Refreshes all posts from the content directory
// @Tags admin
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/v1/admin/refresh [post]
func (h *PostHandler) Refresh(c *gin.Context) {
	if err := h.svc.ParseContentDir(); err != nil {
		h.log.Error("Failed to parse content: %v", err)
		response.InternalError(c, "Failed to refresh content")
		return
	}

	h.log.Info("Content refreshed successfully")
	response.Success(c, gin.H{"message": "Content refreshed"})
}

// GetTags returns all tags with post count
// @Summary Get all tags
// @Description Returns all tags with their post counts
// @Tags posts
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/v1/tags [get]
func (h *PostHandler) GetTags(c *gin.Context) {
	tags, err := h.svc.GetAllTags()
	if err != nil {
		h.log.Error("Failed to get tags: %v", err)
		response.InternalError(c, "Failed to get tags")
		return
	}

	response.Success(c, tags)
}

// GetCategories returns all categories with post count
// @Summary Get all categories
// @Description Returns all categories with their post counts
// @Tags posts
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/v1/categories [get]
func (h *PostHandler) GetCategories(c *gin.Context) {
	categories, err := h.svc.GetAllCategories()
	if err != nil {
		h.log.Error("Failed to get categories: %v", err)
		response.InternalError(c, "Failed to get categories")
		return
	}

	response.Success(c, categories)
}

// SearchPosts 搜索文章
// @Summary Search posts
// @Description Search posts by keyword in title and content
// @Tags posts
// @Produce json
// @Param q query string true "Search keyword"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} response.Response
// @Router /api/v1/posts/search [get]
func (h *PostHandler) SearchPosts(c *gin.Context) {
	keyword := c.Query("q")
	if keyword == "" {
		response.BadRequest(c, "Search keyword is required")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	posts, total, err := h.svc.SearchPosts(keyword, page, pageSize)
	if err != nil {
		h.log.Error("Failed to search posts: %v", err)
		response.InternalError(c, "Failed to search posts")
		return
	}

	response.Success(c, gin.H{
		"list":      posts,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// ListFeaturedPosts 获取推荐文章
// @Summary Get featured posts
// @Description Returns a list of featured posts
// @Tags posts
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/v1/posts/featured [get]
func (h *PostHandler) ListFeaturedPosts(c *gin.Context) {
	posts, err := h.svc.ListFeaturedPosts()
	if err != nil {
		h.log.Error("Failed to list featured posts: %v", err)
		response.InternalError(c, "Failed to load featured posts")
		return
	}

	response.Success(c, posts)
}

// SetPostPin 设置/取消置顶
// @Summary Set post pin status
// @Description Toggle pin status of a post
// @Tags admin
// @Security BearerAuth
// @Produce json
// @Param id path int true "Post ID"
// @Param pinned query bool true "Pin status"
// @Success 200 {object} response.Response
// @Router /api/admin/posts/{id}/pin [put]
func (h *PostHandler) SetPostPin(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid post ID")
		return
	}

	pinnedStr := c.Query("pinned")
	pinned := pinnedStr == "true" || pinnedStr == "1"

	if err := h.svc.SetPostPin(uint(id), pinned); err != nil {
		h.log.Error("Failed to set post pin: %v", err)
		response.InternalError(c, "Failed to update post")
		return
	}

	response.Success(c, gin.H{"message": "Post pin status updated"})
}

// SetPostFeature 设置/取消推荐
// @Summary Set post feature status
// @Description Toggle feature status of a post
// @Tags admin
// @Security BearerAuth
// @Produce json
// @Param id path int true "Post ID"
// @Param featured query bool true "Feature status"
// @Success 200 {object} response.Response
// @Router /api/admin/posts/{id}/feature [put]
func (h *PostHandler) SetPostFeature(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid post ID")
		return
	}

	featuredStr := c.Query("featured")
	featured := featuredStr == "true" || featuredStr == "1"

	if err := h.svc.SetPostFeature(uint(id), featured); err != nil {
		h.log.Error("Failed to set post feature: %v", err)
		response.InternalError(c, "Failed to update post")
		return
	}

	response.Success(c, gin.H{"message": "Post feature status updated"})
}

// SchedulePost 设置定时发布
// @Summary Schedule a post
// @Description Set a post to be published at a scheduled time
// @Tags admin
// @Security BearerAuth
// @Produce json
// @Param id path int true "Post ID"
// @Param scheduled_at body string true "Scheduled time (RFC3339 format)"
// @Success 200 {object} response.Response
// @Router /api/admin/posts/{id}/schedule [post]
func (h *PostHandler) SchedulePost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid post ID")
		return
	}

	var req struct {
		ScheduledAt string `json:"scheduled_at" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	scheduledAt, err := time.Parse(time.RFC3339, req.ScheduledAt)
	if err != nil {
		response.BadRequest(c, "Invalid date format, use RFC3339")
		return
	}

	if err := h.svc.SchedulePost(uint(id), scheduledAt); err != nil {
		h.log.Error("Failed to schedule post: %v", err)
		response.InternalError(c, "Failed to schedule post")
		return
	}

	response.Success(c, gin.H{"message": "Post scheduled successfully"})
}

// GetTOC 获取文章目录
// @Summary Get post table of contents
// @Description Returns the table of contents for a post
// @Tags posts
// @Produce json
// @Param slug path string true "Post slug"
// @Success 200 {object} response.Response
// @Router /api/v1/posts/{slug}/toc [get]
func (h *PostHandler) GetTOC(c *gin.Context) {
	slug := c.Param("slug")
	_, content, err := h.svc.GetPost(slug)
	if err != nil {
		response.NotFound(c, "Post not found")
		return
	}

	toc := h.svc.ExtractTOC(content)
	response.Success(c, gin.H{"toc": toc})
}

// GetNavigation 获取文章导航（上一篇、下一篇）
// @Summary Get post navigation
// @Description Returns previous and next posts for navigation
// @Tags posts
// @Produce json
// @Param slug path string true "Post slug"
// @Success 200 {object} response.Response
// @Router /api/v1/posts/{slug}/navigation [get]
func (h *PostHandler) GetNavigation(c *gin.Context) {
	slug := c.Param("slug")
	nav, err := h.svc.GetPostNavigation(slug)
	if err != nil {
		response.NotFound(c, "Post not found")
		return
	}

	response.Success(c, nav)
}

// ListPopularPosts 获取热门文章
// @Summary Get popular posts
// @Description Returns a list of popular posts by views
// @Tags posts
// @Produce json
// @Param limit query int false "Number of posts to return" default(10)
// @Success 200 {object} response.Response
// @Router /api/v1/posts/popular [get]
func (h *PostHandler) ListPopularPosts(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit <= 0 || limit > 50 {
		limit = 10
	}

	posts, err := h.svc.ListPopularPosts(limit)
	if err != nil {
		h.log.Error("Failed to list popular posts: %v", err)
		response.InternalError(c, "Failed to load popular posts")
		return
	}

	response.Success(c, posts)
}

// ListRelatedPosts 获取相关文章
// @Summary Get related posts
// @Description Returns a list of related posts by tags
// @Tags posts
// @Produce json
// @Param slug path string true "Post slug"
// @Param limit query int false "Number of posts to return" default(5)
// @Success 200 {object} response.Response
// @Router /api/v1/posts/{slug}/related [get]
func (h *PostHandler) ListRelatedPosts(c *gin.Context) {
	slug := c.Param("slug")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))
	if limit <= 0 || limit > 20 {
		limit = 5
	}

	posts, err := h.svc.ListRelatedPosts(slug, limit)
	if err != nil {
		response.NotFound(c, "Post not found")
		return
	}

	response.Success(c, posts)
}

// LikePost 点赞/取消点赞文章
// @Summary Like or unlike a post
// @Description Toggle like status of a post
// @Tags posts
// @Security BearerAuth
// @Produce json
// @Param slug path string true "Post slug"
// @Success 200 {object} response.Response
// @Router /api/v1/posts/{slug}/like [post]
func (h *PostHandler) LikePost(c *gin.Context) {
	slug := c.Param("slug")

	// 获取登录用户ID
	claims := middleware.GetUserClaims(c)
	if claims == nil {
		response.Unauthorized(c, "Please login first")
		return
	}

	liked, count, err := h.svc.LikePost(slug, uint(claims.UserID))
	if err != nil {
		h.log.Error("Failed to like post: %v", err)
		response.InternalError(c, "Failed to like post")
		return
	}

	response.Success(c, gin.H{
		"liked": liked,
		"count": count,
	})
}

// GetPostLikes 获取文章点赞数
// @Summary Get post like count
// @Description Returns the like count of a post and current user's like status
// @Tags posts
// @Produce json
// @Param slug path string true "Post slug"
// @Success 200 {object} response.Response
// @Router /api/v1/posts/{slug}/likes [get]
func (h *PostHandler) GetPostLikes(c *gin.Context) {
	slug := c.Param("slug")

	count, err := h.svc.GetPostLikeCount(slug)
	if err != nil {
		response.NotFound(c, "Post not found")
		return
	}

	result := gin.H{"count": count}

	// Check if user is logged in and get their like status
	claims := middleware.GetUserClaims(c)
	if claims != nil {
		liked, err := h.svc.HasUserLikedPost(slug, uint(claims.UserID))
		if err == nil {
			result["liked"] = liked
		}
	}

	response.Success(c, result)
}

// FavoritePost 收藏/取消收藏文章
// @Summary Favorite or unfavorite a post
// @Description Toggle favorite status of a post
// @Tags posts
// @Security BearerAuth
// @Produce json
// @Param slug path string true "Post slug"
// @Success 200 {object} response.Response
// @Router /api/v1/posts/{slug}/favorite [post]
func (h *PostHandler) FavoritePost(c *gin.Context) {
	slug := c.Param("slug")

	claims := middleware.GetUserClaims(c)
	if claims == nil {
		response.Unauthorized(c, "Please login first")
		return
	}

	favorited, err := h.svc.FavoritePost(slug, uint(claims.UserID))
	if err != nil {
		h.log.Error("Failed to favorite post: %v", err)
		response.InternalError(c, "Failed to favorite post")
		return
	}

	response.Success(c, gin.H{"favorited": favorited})
}

// GetPostFavorite 获取文章收藏状态
// @Summary Get post favorite status
// @Description Returns whether the current user has favorited the post
// @Tags posts
// @Produce json
// @Param slug path string true "Post slug"
// @Success 200 {object} response.Response
// @Router /api/v1/posts/{slug}/favorite [get]
func (h *PostHandler) GetPostFavorite(c *gin.Context) {
	slug := c.Param("slug")

	claims := middleware.GetUserClaims(c)
	if claims == nil {
		response.Success(c, gin.H{"favorited": false})
		return
	}

	favorited, err := h.svc.HasUserFavoritedPost(slug, uint(claims.UserID))
	if err != nil {
		response.NotFound(c, "Post not found")
		return
	}

	response.Success(c, gin.H{"favorited": favorited})
}

// ListMyFavorites 获取我的收藏列表
// @Summary Get my favorites
// @Description Returns the current user's favorite posts
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/v1/users/me/favorites [get]
func (h *PostHandler) ListMyFavorites(c *gin.Context) {
	claims := middleware.GetUserClaims(c)
	if claims == nil {
		response.Unauthorized(c, "Please login first")
		return
	}

	posts, err := h.svc.ListMyFavorites(uint(claims.UserID))
	if err != nil {
		h.log.Error("Failed to list favorites: %v", err)
		response.InternalError(c, "Failed to load favorites")
		return
	}

	response.Success(c, posts)
}

// CheckPassword 检查文章是否有密码保护
// @Summary Check if post has password protection
// @Description Returns whether a post requires password
// @Tags posts
// @Produce json
// @Param slug path string true "Post slug"
// @Success 200 {object} response.Response
// @Router /api/v1/posts/{slug}/check [get]
func (h *PostHandler) CheckPassword(c *gin.Context) {
	slug := c.Param("slug")

	hasPassword, err := h.svc.HasPassword(slug)
	if err != nil {
		response.NotFound(c, "Post not found")
		return
	}

	response.Success(c, gin.H{"protected": hasPassword})
}

// VerifyPassword 验证文章密码
// @Summary Verify post password
// @Description Verify if the provided password is correct
// @Tags posts
// @Produce json
// @Param slug path string true "Post slug"
// @Param password body string true "Password"
// @Success 200 {object} response.Response
// @Router /api/v1/posts/{slug}/verify [post]
func (h *PostHandler) VerifyPassword(c *gin.Context) {
	slug := c.Param("slug")

	var req struct {
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Password required")
		return
	}

	valid, err := h.svc.VerifyPostPassword(slug, req.Password)
	if err != nil {
		response.NotFound(c, "Post not found")
		return
	}

	if !valid {
		response.Error(c, 403, "Invalid password")
		return
	}

	response.Success(c, gin.H{"valid": true})
}

// SetPassword 设置文章密码（管理后台）
// @Summary Set post password
// @Description Set or remove password protection for a post
// @Tags admin
// @Security BearerAuth
// @Produce json
// @Param slug path string true "Post slug"
// @Param password body string true "Password (empty to remove)"
// @Success 200 {object} response.Response
// @Router /api/admin/posts/{slug}/password [post]
func (h *PostHandler) SetPassword(c *gin.Context) {
	slug := c.Param("slug")

	var req struct {
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request")
		return
	}

	if err := h.svc.SetPostPassword(slug, req.Password); err != nil {
		h.log.Error("Failed to set password: %v", err)
		response.InternalError(c, "Failed to set password")
		return
	}

	response.Success(c, gin.H{"message": "Password set successfully"})
}

// GetSitemap generates XML sitemap
// @Summary Get sitemap
// @Description Returns XML sitemap for SEO
// @Tags posts
// @Produce xml
// @Success 200 {string} string
// @Router /api/v1/sitemap.xml [get]
func (h *PostHandler) GetSitemap(c *gin.Context) {
	posts, err := h.svc.GetAllPostsForSitemap()
	if err != nil {
		h.log.Error("Failed to get posts for sitemap: %v", err)
		response.InternalError(c, "Failed to generate sitemap")
		return
	}

	baseURL := "https://jumoshen.cn"

	xml := `<?xml version="1.0" encoding="UTF-8"?>`
	xml += `<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`

	// Homepage
	xml += `<url><loc>` + baseURL + `/</loc><changefreq>weekly</changefreq><priority>1.0</priority></url>`

	// Archives
	xml += `<url><loc>` + baseURL + `/archives</loc><changefreq>weekly</changefreq><priority>0.8</priority></url>`

	// About
	xml += `<url><loc>` + baseURL + `/about</loc><changefreq>monthly</changefreq><priority>0.5</priority></url>`

	// Posts
	for _, post := range posts {
		xml += `<url>`
		xml += `<loc>` + baseURL + `/post/` + post.Slug + `</loc>`
		xml += `<lastmod>` + post.Date.Format("2006-01-02") + `</lastmod>`
		xml += `<changefreq>monthly</changefreq>`
		xml += `<priority>0.9</priority>`
		xml += `</url>`
	}

	xml += `</urlset>`

	c.Header("Content-Type", "application/xml")
	c.String(200, xml)
}

// GetRSS generates RSS feed
// @Summary Get RSS feed
// @Description Returns RSS 2.0 feed
// @Tags posts
// @Produce xml
// @Success 200 {string} string
// @Router /api/v1/feed.xml [get]
func (h *PostHandler) GetRSS(c *gin.Context) {
	posts, err := h.svc.GetAllPostsForSitemap()
	if err != nil {
		h.log.Error("Failed to get posts for RSS: %v", err)
		response.InternalError(c, "Failed to generate RSS")
		return
	}

	baseURL := "https://jumoshen.cn"
	siteTitle := "Jumoshen"
	siteDescription := "巨魔深博客"

	var xml string
	xml = `<?xml version="1.0" encoding="UTF-8"?>`
	xml += `<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">`
	xml += `<channel>`
	xml += `<title>` + siteTitle + `</title>`
	xml += `<description>` + siteDescription + `</description>`
	xml += `<link>` + baseURL + `</link>`
	xml += `<atom:link href="` + baseURL + `/api/v1/feed.xml" rel="self" type="application/rss+xml"/>`
	xml += `<language>zh-CN</language>`

	for _, post := range posts {
		xml += `<item>`
		xml += `<title><![CDATA[` + post.Title + `]]></title>`
		xml += `<link>` + baseURL + `/post/` + post.Slug + `</link>`
		xml += `<guid isPermaLink="true">` + baseURL + `/post/` + post.Slug + `</guid>`
		xml += `<pubDate>` + post.Date.Format(time.RFC1123) + `</pubDate>`
		if post.Summary != "" {
			xml += `<description><![CDATA[` + post.Summary + `]]></description>`
		}
		xml += `</item>`
	}

	xml += `</channel>`
	xml += `</rss>`

	c.Header("Content-Type", "application/xml; charset=utf-8")
	c.String(200, xml)
}

type PostInfo struct {
	Slug       string `json:"slug"`
	Title      string `json:"title"`
	Date       string `json:"date"`
	Tags       []string `json:"tags"`
	Categories []string `json:"categories"`
	Summary    string `json:"summary"`
	Views      int64  `json:"views"`
}

type PostDetail struct {
	Post    PostInfo `json:"post"`
	Content template.HTML `json:"content"`
}
