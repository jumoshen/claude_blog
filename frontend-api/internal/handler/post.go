package handler

import (
	"html/template"
	"strconv"

	"github.com/gin-gonic/gin"

	"markdown-blog/internal/logger"
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
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} response.Response{data=[]PostInfo}
// @Router /api/v1/posts [get]
func (h *PostHandler) ListPosts(c *gin.Context) {
	filterTag := c.Query("tag")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	posts, total, err := h.svc.ListPostsPaginated(filterTag, page, pageSize)
	if err != nil {
		h.log.Error("Failed to list posts: %v", err)
		response.InternalError(c, "Failed to load posts")
		return
	}

	h.log.Info("ListPosts: page=%d, page_size=%d, total=%d, tag=%s", page, pageSize, total, filterTag)
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
