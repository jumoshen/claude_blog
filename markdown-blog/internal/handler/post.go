package handler

import (
	"html/template"

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

// ListPosts returns all posts
// @Summary Get all posts
// @Description Returns a list of all blog posts
// @Tags posts
// @Produce json
// @Success 200 {object} response.Response{data=[]PostInfo}
// @Router /api/v1/posts [get]
func (h *PostHandler) ListPosts(c *gin.Context) {
	filterTag := c.Query("tag")

	posts, err := h.svc.ListPosts()
	if err != nil {
		h.log.Error("Failed to list posts: %v", err)
		response.InternalError(c, "Failed to load posts")
		return
	}

	// Filter by tag if specified
	if filterTag != "" {
		var filtered []service.PostInfo
		for _, p := range posts {
			for _, tag := range p.Tags {
				if tag == filterTag {
					filtered = append(filtered, p)
					break
				}
			}
		}
		posts = filtered
	}

	h.log.Info("ListPosts: count=%d, filter=%s", len(posts), filterTag)
	response.Success(c, posts)
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
