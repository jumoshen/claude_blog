package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"markdown-blog/internal/config"
	"markdown-blog/internal/model"
	"markdown-blog/internal/repository"
)

type Service struct {
	repo *repository.Repository
	cfg  *config.Config
	md   goldmark.Markdown
}

func New(cfg *config.Config, repo *repository.Repository) *Service {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
		goldmark.WithRendererOptions(html.WithHardWraps()),
	)

	return &Service{
		repo: repo,
		cfg:  cfg,
		md:   md,
	}
}

func NewWithRedis(cfg *config.Config, redisClient *redis.Client) (*Service, error) {
	repo, err := repository.New(cfg)
	if err != nil {
		return nil, err
	}

	return New(cfg, repo), nil
}

// toPostInfoList converts posts to PostInfo list using normalized tables for tags/categories
func (s *Service) toPostInfoList(posts []model.Post) []PostInfo {
	if len(posts) == 0 {
		return []PostInfo{}
	}

	postIDs := make([]uint, len(posts))
	for i, p := range posts {
		postIDs[i] = p.ID
	}

	tagMap, _ := s.repo.GetPostTags(postIDs)
	catMap, _ := s.repo.GetPostCategories(postIDs)

	result := make([]PostInfo, 0, len(posts))
	for _, p := range posts {
		result = append(result, PostInfo{
			Slug:        p.Slug,
			Title:       p.Title,
			Date:        p.Date,
			Tags:        tagMap[p.ID],
			Categories:  catMap[p.ID],
			Summary:     p.Summary,
			Views:       p.Views,
			IsPinned:    p.IsPinned,
			IsFeatured:  p.IsFeatured,
			ScheduledAt: p.ScheduledAt,
			ReadingTime: CalculateReadingTime(p.Content),
		})
	}
	return result
}

type PostInfo struct {
	Slug        string     `json:"slug"`
	Title       string     `json:"title"`
	Date        time.Time  `json:"date"`
	Tags        []string   `json:"tags"`
	Categories  []string   `json:"categories"`
	Summary     string     `json:"summary"`
	Views       int64      `json:"views"`
	IsPinned    bool       `json:"is_pinned"`
	IsFeatured  bool       `json:"is_featured"`
	ScheduledAt *time.Time `json:"scheduled_at,omitempty"`
	ReadingTime int        `json:"reading_time"` // 阅读时间（分钟）
}

type UserInfo struct {
	ID        uint
	GitHubID  int64
	Login     string
	Name      string
	AvatarURL string
	Email     string
}

func (s *Service) GetRepo() *repository.Repository {
	return s.repo
}

// ListPosts returns all posts
func (s *Service) ListPosts() ([]PostInfo, error) {
	posts, err := s.repo.ListPosts()
	if err != nil {
		return nil, err
	}

	if len(posts) == 0 {
		return []PostInfo{}, nil
	}

	// 获取所有 post IDs
	postIDs := make([]uint, len(posts))
	for i, p := range posts {
		postIDs[i] = p.ID
	}

	// 从 normalized 表获取 tags 和 categories
	tagMap, _ := s.repo.GetPostTags(postIDs)
	catMap, _ := s.repo.GetPostCategories(postIDs)

	result := make([]PostInfo, 0, len(posts))
	for _, p := range posts {
		result = append(result, PostInfo{
			Slug:        p.Slug,
			Title:       p.Title,
			Date:        p.Date,
			Tags:        tagMap[p.ID],
			Categories:  catMap[p.ID],
			Summary:     p.Summary,
			Views:       p.Views,
			IsPinned:    p.IsPinned,
			IsFeatured:  p.IsFeatured,
			ScheduledAt: p.ScheduledAt,
		})
	}
	return result, nil
}

// ListPostsPaginated 分页获取已发布文章
func (s *Service) ListPostsPaginated(tag string, category string, page, pageSize int) ([]PostInfo, int64, error) {
	posts, total, err := s.repo.ListPostsPaginated(tag, category, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	if len(posts) == 0 {
		return []PostInfo{}, total, nil
	}

	// 获取所有 post IDs
	postIDs := make([]uint, len(posts))
	for i, p := range posts {
		postIDs[i] = p.ID
	}

	// 从 normalized 表获取 tags 和 categories
	tagMap, _ := s.repo.GetPostTags(postIDs)
	catMap, _ := s.repo.GetPostCategories(postIDs)

	result := make([]PostInfo, 0, len(posts))
	for _, p := range posts {
		result = append(result, PostInfo{
			Slug:        p.Slug,
			Title:       p.Title,
			Date:        p.Date,
			Tags:        tagMap[p.ID],
			Categories:  catMap[p.ID],
			Summary:     p.Summary,
			Views:       p.Views,
			IsPinned:    p.IsPinned,
			IsFeatured:  p.IsFeatured,
			ScheduledAt: p.ScheduledAt,
			ReadingTime: CalculateReadingTime(p.Content),
		})
	}
	return result, total, nil
}

// GetPost returns a single post by slug
func (s *Service) GetPost(slug string) (*PostInfo, string, error) {
	post, err := s.repo.GetPostBySlug(slug)
	if err != nil {
		return nil, "", err
	}

	// Increment views
	s.repo.GetDB().Model(post).Update("views", gorm.Expr("views + 1"))

	// 从 normalized 表获取 tags 和 categories
	tagMap, _ := s.repo.GetPostTags([]uint{post.ID})
	catMap, _ := s.repo.GetPostCategories([]uint{post.ID})

	info := &PostInfo{
		Slug:        post.Slug,
		Title:       post.Title,
		Date:        post.Date,
		Tags:        tagMap[post.ID],
		Categories:  catMap[post.ID],
		Summary:     post.Summary,
		Views:       post.Views + 1,
		IsPinned:    post.IsPinned,
		IsFeatured:  post.IsFeatured,
		ReadingTime: CalculateReadingTime(post.Content),
		ScheduledAt: post.ScheduledAt,
	}

	return info, post.Content, nil
}

// ParseContentDir parses all markdown files in the content directory
func (s *Service) ParseContentDir() error {
	postsPath := filepath.Join(s.cfg.Content.Path, "post")
	return filepath.Walk(postsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".md") {
			return s.parseAndSavePost(path)
		}
		return nil
	})
}

func (s *Service) parseAndSavePost(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	slug := strings.TrimSuffix(filepath.Base(path), ".md")
	title, dateStr, tagsRaw, categoriesRaw, summary, bodyLines := parseFrontMatter(string(data))

	// Parse tags
	tags := parseArrayField(tagsRaw)
	categories := parseArrayField(categoriesRaw)

	if summary == "" && len(bodyLines) > 0 {
		sumLines := bodyLines
		if len(sumLines) > 3 {
			sumLines = sumLines[:3]
		}
		summary = strings.TrimSpace(strings.Join(sumLines, " "))
		summary = truncateString(summary, 200)
	}

	date, _ := time.Parse("2006-01-02T15:04:05Z07:00", dateStr)
	if date.IsZero() {
		date, _ = time.Parse("2006-01-02", dateStr)
	}
	if date.IsZero() {
		date = time.Now()
	}

	body := strings.Join(bodyLines, "\n")
	var buf bytes.Buffer
	if err := s.md.Convert([]byte(body), &buf); err != nil {
		return err
	}
	htmlContent := buf.String()

	if title == "" {
		title = slug
	}

	post := model.Post{
		Slug:       slug,
		Title:      title,
		Date:       date,
		Tags:       tags,
		Categories: categories,
		Summary:    summary,
		Content:    htmlContent,
	}

	return s.repo.UpsertPost(&post)
}

func parseFrontMatter(content string) (title, dateStr, tagsRaw, categoriesRaw, summary string, bodyLines []string) {
	lines := strings.Split(content, "\n")
	var inFrontMatter bool
	var frontMatterEnded bool

	for _, line := range lines {
		if strings.HasPrefix(line, "---") {
			if !inFrontMatter {
				inFrontMatter = true
				continue
			} else {
				frontMatterEnded = true
				continue
			}
		}
		if !frontMatterEnded {
			if strings.HasPrefix(line, "title:") {
				title = strings.TrimSpace(strings.TrimPrefix(line, "title:"))
				title = strings.Trim(title, "\"")
			} else if strings.HasPrefix(line, "date:") {
				dateStr = strings.TrimSpace(strings.TrimPrefix(line, "date:"))
			} else if strings.HasPrefix(line, "tags:") {
				tagsRaw = strings.TrimSpace(strings.TrimPrefix(line, "tags:"))
			} else if strings.HasPrefix(line, "categories:") {
				categoriesRaw = strings.TrimSpace(strings.TrimPrefix(line, "categories:"))
			}
		} else {
			if strings.Contains(line, "<!--more-->") {
				parts := strings.Split(content, "<!--more-->")
				rawSummary := parts[0]
				if idx := strings.LastIndex(rawSummary, "---"); idx >= 0 {
					rawSummary = rawSummary[idx+3:]
				}
				summary = strings.TrimSpace(rawSummary)
				summary = truncateString(summary, 200)
				continue
			}
			bodyLines = append(bodyLines, line)
		}
	}
	return
}

// parseArrayField parses ["tag1", "tag2"] or [tag1, tag2] into "tag1,tag2"
func parseArrayField(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	re := regexp.MustCompile(`"([^"]+)"`)
	matches := re.FindAllStringSubmatch(raw, -1)
	if len(matches) > 0 {
		var result []string
		for _, m := range matches {
			result = append(result, m[1])
		}
		return strings.Join(result, ",")
	}
	return strings.Trim(strings.Trim(raw, "[]"), " ")
}

func truncateString(s string, n int) string {
	if utf8.RuneCountInString(s) <= n {
		return s
	}
	runes := []rune(s)
	return string(runes[:n]) + "..."
}

// GetArchives returns posts grouped by year
func (s *Service) GetArchives() (map[string][]PostInfo, error) {
	posts, err := s.ListPosts()
	if err != nil {
		return nil, err
	}

	archives := make(map[string][]PostInfo)
	for _, p := range posts {
		year := p.Date.Format("2006")
		archives[year] = append(archives[year], p)
	}

	years := make([]string, 0, len(archives))
	for y := range archives {
		years = append(years, y)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(years)))

	result := make(map[string][]PostInfo)
	for _, y := range years {
		result[y] = archives[y]
	}

	return result, nil
}

// GetAboutContent returns the about page content
func (s *Service) GetAboutContent() (string, error) {
	aboutPath := filepath.Join(s.cfg.Content.Path, "about", "index.md")
	return s.getContent(aboutPath)
}

func (s *Service) getContent(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	content := stripFrontMatter(string(data))

	var buf bytes.Buffer
	if err := s.md.Convert([]byte(content), &buf); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func stripFrontMatter(content string) string {
	lines := strings.Split(content, "\n")
	if len(lines) < 3 || lines[0] != "---" {
		return content
	}

	for i := 1; i < len(lines); i++ {
		if lines[i] == "---" {
			return strings.Join(lines[i+1:], "\n")
		}
	}
	return content
}

// ServeMarkdownFile serves a markdown file as HTML
func (s *Service) ServeMarkdownFile(filename string, w io.Writer) error {
	path := filepath.Join(s.cfg.Content.Path, filename)
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := s.md.Convert(data, &buf); err != nil {
		return err
	}

	_, err = w.Write(buf.Bytes())
	return err
}

// GetAllTags returns all tags with post count
func (s *Service) GetAllTags() ([]repository.TagWithCount, error) {
	return s.repo.GetAllTags()
}

// GetAllCategories returns all categories with post count
func (s *Service) GetAllCategories() (map[string]int, error) {
	return s.repo.GetAllCategories()
}

// GetAllPostsForSitemap returns all posts for sitemap generation
func (s *Service) GetAllPostsForSitemap() ([]PostInfo, error) {
	return s.ListPosts()
}

// User operations
func (s *Service) SaveUser(user *model.User) error {
	return s.repo.GetDB().Save(user).Error
}

func (s *Service) GetUserByGitHubID(id int64) (*model.User, error) {
	return s.repo.GetUserByGitHubID(id)
}

func (s *Service) GetUserByID(id uint) (*model.User, error) {
	return s.repo.GetUserByID(id)
}

// Logout adds the token to blacklist
func (s *Service) Logout(ctx context.Context, jti string, exp time.Time) error {
	ttl := time.Until(exp)
	if ttl <= 0 {
		return nil
	}
	return s.repo.AddToBlacklist(ctx, jti, ttl)
}

// RecordVisit 记录访问日志（防刷）
func (s *Service) RecordVisit(slug, ip, userAgent string, userID int64) {
	// 异步检查并记录，不阻塞请求
	go func() {
		ctx := context.Background()

		// 防刷检查
		canRecord, _ := s.repo.CanRecordVisit(ctx, slug, ip, userID)
		if !canRecord {
			return // 1小时内已记录过，跳过
		}

		visit := &model.Visit{
			PostSlug:  slug,
			IP:        ip,
			UserAgent: userAgent,
			UserID:    userID,
		}
		s.repo.CreateVisit(visit)
	}()
}

// Comment operations
func (s *Service) CreateComment(comment *model.Comment) error {
	return s.repo.CreateComment(comment)
}

func (s *Service) GetCommentsByPostSlug(postSlug string, limit int) ([]model.Comment, error) {
	return s.repo.GetCommentsByPostSlug(postSlug, limit)
}

// CheckCommentRateLimit 检查评论频率限制
// 匿名用户: IP每分钟3条 + DeviceID每分钟5条
// 登录用户: UserID每10秒1条
func (s *Service) CheckCommentRateLimit(ctx context.Context, ip, deviceID string, userID int64) (bool, error) {
	if userID > 0 {
		// 登录用户：每10秒1条
		key := fmt.Sprintf("comment:rate:user:%d", userID)
		count, err := s.repo.IncrCommentRateLimit(ctx, key, 10*time.Second)
		if err != nil {
			return true, err
		}
		return count <= 1, nil
	}

	// 匿名用户：IP每分钟3条
	ipKey := fmt.Sprintf("comment:rate:ip:%s", ip)
	ipCount, err := s.repo.IncrCommentRateLimit(ctx, ipKey, time.Minute)
	if err != nil {
		return true, err
	}
	if ipCount > 3 {
		return false, nil
	}

	// DeviceID每分钟5条
	if deviceID != "" {
		deviceKey := fmt.Sprintf("comment:rate:device:%s", deviceID)
		deviceCount, err := s.repo.IncrCommentRateLimit(ctx, deviceKey, time.Minute)
		if err != nil {
			return true, err
		}
		return deviceCount <= 5, nil
	}

	return true, nil
}

// ContainsSensitiveWords 检测敏感词
func (s *Service) ContainsSensitiveWords(content string) bool {
	words, err := s.GetSensitiveWords()
	if err != nil {
		// 如果获取失败，使用内存检测作为fallback
		return s.checkSensitiveWordsFallback(content)
	}

	if len(words) == 0 {
		return false
	}

	for _, word := range words {
		if containsString(content, word) {
			return true
		}
	}
	return false
}

// GetSensitiveWords 获取敏感词列表（带缓存）
func (s *Service) GetSensitiveWords() ([]string, error) {
	ctx := context.Background()

	// 先从缓存获取
	words, err := s.repo.GetSensitiveWordsCache(ctx)
	if err == nil && words != nil {
		return words, nil
	}

	// 缓存没有，从数据库获取
	dbWords, err := s.repo.GetAllSensitiveWords()
	if err != nil {
		return nil, err
	}

	words = make([]string, 0, len(dbWords))
	for _, w := range dbWords {
		words = append(words, w.Word)
	}

	// 更新缓存
	if len(words) > 0 {
		s.repo.SetSensitiveWordsCache(ctx, words)
	}

	return words, nil
}

// InvalidateSensitiveWordsCache 使敏感词缓存失效
func (s *Service) InvalidateSensitiveWordsCache() error {
	return s.repo.InvalidateSensitiveWordsCache(context.Background())
}

// checkSensitiveWordsFallback 敏感词检测fallback（内存检测）
func (s *Service) checkSensitiveWordsFallback(content string) bool {
	fallbackWords := []string{
		"习近平", "毛主席", "周恩来", "刘少奇", "朱德", "邓小平",
		"傻逼", "傻b", "sb", "操", "艹", "妈逼", "妈b", "mb",
		"fuck", "shit", "色情", "赌博", "赌场", "澳门赌场",
	}

	for _, word := range fallbackWords {
		if containsString(content, word) {
			return true
		}
	}
	return false
}

func containsString(s, substr string) bool {
	sLower := toLower(s)
	substrLower := toLower(substr)

	for i := 0; i <= len(sLower)-len(substrLower); i++ {
		if sLower[i:i+len(substrLower)] == substrLower {
			return true
		}
	}
	return false
}

func toLower(s string) string {
	var result []byte
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		result = append(result, c)
	}
	return string(result)
}

// CategoryInfo 分类信息
type CategoryInfo struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Sort        int    `json:"sort"`
	PostCount   int64  `json:"post_count"`
	CreatedAt   string `json:"created_at"`
}

// Category operations
func (s *Service) ListCategories(page, pageSize int) ([]CategoryInfo, int64, error) {
	categories, total, err := s.repo.ListCategories(page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	result := make([]CategoryInfo, 0, len(categories))
	for _, c := range categories {
		count, _ := s.repo.CountPostsByCategory(c.ID)
		result = append(result, CategoryInfo{
			ID:          c.ID,
			Name:        c.Name,
			Slug:        c.Slug,
			Description: c.Description,
			Sort:        c.Sort,
			PostCount:   count,
			CreatedAt:   c.CreatedAt.Format("2006-01-02"),
		})
	}
	return result, total, nil
}

func (s *Service) ListAllCategories() ([]CategoryInfo, error) {
	categories, err := s.repo.ListAllCategories()
	if err != nil {
		return nil, err
	}

	result := make([]CategoryInfo, 0, len(categories))
	for _, c := range categories {
		count, _ := s.repo.CountPostsByCategory(c.ID)
		result = append(result, CategoryInfo{
			ID:          c.ID,
			Name:        c.Name,
			Slug:        c.Slug,
			Description: c.Description,
			Sort:        c.Sort,
			PostCount:   count,
			CreatedAt:   c.CreatedAt.Format("2006-01-02"),
		})
	}
	return result, nil
}

func (s *Service) GetCategoryByID(id uint) (*CategoryInfo, error) {
	category, err := s.repo.GetCategoryByID(id)
	if err != nil {
		return nil, err
	}
	count, _ := s.repo.CountPostsByCategory(category.ID)
	return &CategoryInfo{
		ID:          category.ID,
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		Sort:        category.Sort,
		PostCount:   count,
		CreatedAt:   category.CreatedAt.Format("2006-01-02"),
	}, nil
}

func (s *Service) CreateCategory(name, slug, description string, sort int) (*CategoryInfo, error) {
	category := &model.Category{
		Name:        name,
		Slug:        slug,
		Description: description,
		Sort:        sort,
	}
	if err := s.repo.CreateCategory(category); err != nil {
		return nil, err
	}
	return &CategoryInfo{
		ID:          category.ID,
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		Sort:        category.Sort,
		PostCount:   0,
		CreatedAt:   category.CreatedAt.Format("2006-01-02"),
	}, nil
}

func (s *Service) UpdateCategory(id uint, name, slug, description string, sort int) (*CategoryInfo, error) {
	category, err := s.repo.GetCategoryByID(id)
	if err != nil {
		return nil, err
	}
	category.Name = name
	category.Slug = slug
	category.Description = description
	category.Sort = sort
	if err := s.repo.UpdateCategory(category); err != nil {
		return nil, err
	}
	count, _ := s.repo.CountPostsByCategory(category.ID)
	return &CategoryInfo{
		ID:          category.ID,
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		Sort:        category.Sort,
		PostCount:   count,
		CreatedAt:   category.CreatedAt.Format("2006-01-02"),
	}, nil
}

func (s *Service) DeleteCategory(id uint) error {
	count, err := s.repo.CountPostsByCategory(id)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("该分类下有 %d 篇文章，无法删除", count)
	}
	return s.repo.DeleteCategory(id)
}

// TagInfo 标签信息
type TagInfo struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	Color     string `json:"color"`
	PostCount int64  `json:"post_count"`
	CreatedAt string `json:"created_at"`
}

func (s *Service) ListTags(page, pageSize int) ([]TagInfo, int64, error) {
	tags, total, err := s.repo.ListTags(page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	result := make([]TagInfo, 0, len(tags))
	for _, t := range tags {
		count, _ := s.repo.CountPostsByTag(t.ID)
		result = append(result, TagInfo{
			ID:        t.ID,
			Name:      t.Name,
			Slug:      t.Slug,
			Color:     t.Color,
			PostCount: count,
			CreatedAt: t.CreatedAt.Format("2006-01-02"),
		})
	}
	return result, total, nil
}

func (s *Service) ListAllTags() ([]TagInfo, error) {
	tags, err := s.repo.ListAllTags()
	if err != nil {
		return nil, err
	}

	result := make([]TagInfo, 0, len(tags))
	for _, t := range tags {
		count, _ := s.repo.CountPostsByTag(t.ID)
		result = append(result, TagInfo{
			ID:        t.ID,
			Name:      t.Name,
			Slug:      t.Slug,
			Color:     t.Color,
			PostCount: count,
			CreatedAt: t.CreatedAt.Format("2006-01-02"),
		})
	}
	return result, nil
}

func (s *Service) GetTagByID(id uint) (*TagInfo, error) {
	tag, err := s.repo.GetTagByID(id)
	if err != nil {
		return nil, err
	}
	count, _ := s.repo.CountPostsByTag(tag.ID)
	return &TagInfo{
		ID:        tag.ID,
		Name:      tag.Name,
		Slug:      tag.Slug,
		Color:     tag.Color,
		PostCount: count,
		CreatedAt: tag.CreatedAt.Format("2006-01-02"),
	}, nil
}

func (s *Service) CreateTag(name, slug, color string) (*TagInfo, error) {
	tag := &model.Tag{
		Name:  name,
		Slug:  slug,
		Color: color,
	}
	if err := s.repo.CreateTag(tag); err != nil {
		return nil, err
	}
	return &TagInfo{
		ID:        tag.ID,
		Name:      tag.Name,
		Slug:      tag.Slug,
		Color:     tag.Color,
		PostCount: 0,
		CreatedAt: tag.CreatedAt.Format("2006-01-02"),
	}, nil
}

func (s *Service) UpdateTag(id uint, name, slug, color string) (*TagInfo, error) {
	tag, err := s.repo.GetTagByID(id)
	if err != nil {
		return nil, err
	}
	tag.Name = name
	tag.Slug = slug
	tag.Color = color
	if err := s.repo.UpdateTag(tag); err != nil {
		return nil, err
	}
	count, _ := s.repo.CountPostsByTag(tag.ID)
	return &TagInfo{
		ID:        tag.ID,
		Name:      tag.Name,
		Slug:      tag.Slug,
		Color:     tag.Color,
		PostCount: count,
		CreatedAt: tag.CreatedAt.Format("2006-01-02"),
	}, nil
}

func (s *Service) DeleteTag(id uint) error {
	count, err := s.repo.CountPostsByTag(id)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("该标签被 %d 篇文章使用，无法删除", count)
	}
	return s.repo.DeleteTag(id)
}

// SearchPosts 搜索文章
func (s *Service) SearchPosts(keyword string, page, pageSize int) ([]PostInfo, int64, error) {
	posts, total, err := s.repo.SearchPosts(keyword, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	return s.toPostInfoList(posts), total, nil
}

// ListFeaturedPosts 获取推荐文章
func (s *Service) ListFeaturedPosts() ([]PostInfo, error) {
	posts, err := s.repo.ListPostsFeatured()
	if err != nil {
		return nil, err
	}

	return s.toPostInfoList(posts), nil
}

// SetPostPin 设置/取消置顶
func (s *Service) SetPostPin(id uint, pinned bool) error {
	return s.repo.SetPostPinned(id, pinned)
}

// SetPostFeature 设置/取消推荐
func (s *Service) SetPostFeature(id uint, featured bool) error {
	return s.repo.SetPostFeatured(id, featured)
}

// SchedulePost 设置定时发布
func (s *Service) SchedulePost(id uint, scheduledAt time.Time) error {
	return s.repo.SchedulePost(id, scheduledAt)
}

// PublishScheduledPosts 发布已到时的定时文章
func (s *Service) PublishScheduledPosts() (int64, error) {
	return s.repo.PublishScheduledPosts()
}

// GetPostByID 根据ID获取文章
func (s *Service) GetPostByID(id uint) (*PostInfo, error) {
	post, err := s.repo.GetPostByID(id)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, nil
	}

	infoList := s.toPostInfoList([]model.Post{*post})
	if len(infoList) == 0 {
		return nil, nil
	}
	return &infoList[0], nil
}

// TocItem 目录项
type TocItem struct {
	Level int    `json:"level"`
	Text  string `json:"text"`
	ID    string `json:"id"`
}

// ExtractTOC 从 markdown HTML 中提取目录
func (s *Service) ExtractTOC(htmlContent string) []TocItem {
	var toc []TocItem
	re := regexp.MustCompile(`<h([1-6])[^>]*id="([^"]+)"[^>]*>([^<]+)</h[1-6]>`)
	matches := re.FindAllStringSubmatch(htmlContent, -1)

	for _, m := range matches {
		if len(m) == 4 {
			level, _ := strconv.Atoi(m[1])
			toc = append(toc, TocItem{
				Level: level,
				ID:    m[2],
				Text:  strings.TrimSpace(m[3]),
			})
		}
	}
	return toc
}

// PostNavigation 文章导航信息
type PostNavigation struct {
	Prev *PostNavItem `json:"prev,omitempty"`
	Next *PostNavItem `json:"next,omitempty"`
}

// PostNavItem 导航项
type PostNavItem struct {
	Slug  string `json:"slug"`
	Title string `json:"title"`
}

// CalculateReadingTime 计算阅读时间（分钟）
// 中文：约400字/分钟，英文：约200词/分钟
func CalculateReadingTime(content string) int {
	// 移除 HTML 标签
	re := regexp.MustCompile(`<[^>]+>`)
	plainText := re.ReplaceAllString(content, "")

	// 统计中文字符数
	chineseCount := 0
	// 统计英文单词数
	englishWordCount := 0

	runes := []rune(plainText)
	inEnglish := false
	englishWord := ""

	for _, r := range runes {
		if r >= 0x4e00 && r <= 0x9fff {
			// 中文字符
			chineseCount++
			if inEnglish && englishWord != "" {
				englishWordCount++
			}
			inEnglish = false
			englishWord = ""
		} else if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			// 英文字母
			inEnglish = true
			englishWord += string(r)
		} else {
			if inEnglish && englishWord != "" {
				englishWordCount++
			}
			inEnglish = false
			englishWord = ""
		}
	}

	// 如果末尾是英文，加一次
	if inEnglish && englishWord != "" {
		englishWordCount++
	}

	// 计算总阅读时间（分钟）
	// 中文速度：400字/分钟，英文速度：200词/分钟
	chineseTime := float64(chineseCount) / 400.0
	englishTime := float64(englishWordCount) / 200.0

	totalMinutes := chineseTime + englishTime
	if totalMinutes < 1 {
		return 1
	}
	return int(totalMinutes)
}

// GetPostNavigation 获取文章导航（上一篇、下一篇）
func (s *Service) GetPostNavigation(slug string) (*PostNavigation, error) {
	post, err := s.repo.GetPostBySlug(slug)
	if err != nil {
		return nil, err
	}

	nav := &PostNavigation{}

	// 获取上一篇
	prevPost, err := s.repo.GetPrevPost(post.Date, post.Slug)
	if err == nil && prevPost != nil {
		nav.Prev = &PostNavItem{
			Slug:  prevPost.Slug,
			Title: prevPost.Title,
		}
	}

	// 获取下一篇
	nextPost, err := s.repo.GetNextPost(post.Date, post.Slug)
	if err == nil && nextPost != nil {
		nav.Next = &PostNavItem{
			Slug:  nextPost.Slug,
			Title: nextPost.Title,
		}
	}

	return nav, nil
}

// ListPopularPosts 获取热门文章
func (s *Service) ListPopularPosts(limit int) ([]PostInfo, error) {
	if limit <= 0 {
		limit = 10
	}
	posts, err := s.repo.ListPopularPosts(limit)
	if err != nil {
		return nil, err
	}

	return s.toPostInfoList(posts), nil
}

// ListRelatedPosts 获取相关文章
func (s *Service) ListRelatedPosts(slug string, limit int) ([]PostInfo, error) {
	if limit <= 0 {
		limit = 5
	}

	// 获取当前文章
	currentPost, err := s.repo.GetPostBySlug(slug)
	if err != nil {
		return nil, err
	}

	// 从 normalized 表获取当前文章的 tags
	tagMap, _ := s.repo.GetPostTags([]uint{currentPost.ID})
	currentTags := tagMap[currentPost.ID]

	posts, err := s.repo.ListRelatedPosts(slug, currentTags, limit)
	if err != nil {
		return nil, err
	}

	return s.toPostInfoList(posts), nil
}

// BlogUserInfo 博客用户信息
type BlogUserInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

// RegisterBlogUser 注册博客用户
func (s *Service) RegisterBlogUser(username, email, password string) (*BlogUserInfo, error) {
	// 检查用户名是否存在
	if _, err := s.repo.GetBlogUserByUsername(username); err == nil {
		return nil, fmt.Errorf("用户名已存在")
	}

	// 检查邮箱是否存在
	if _, err := s.repo.GetBlogUserByEmail(email); err == nil {
		return nil, fmt.Errorf("邮箱已被注册")
	}

	// 密码哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败")
	}

	user := &model.BlogUser{
		Username:     username,
		Email:       email,
		PasswordHash: string(hashedPassword),
		Nickname:    username,
	}

	if err := s.repo.CreateBlogUser(user); err != nil {
		return nil, err
	}

	return &BlogUserInfo{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Nickname: user.Nickname,
		Avatar:   user.AvatarURL,
	}, nil
}

// LoginBlogUser 登录博客用户
func (s *Service) LoginBlogUser(username, password string) (*BlogUserInfo, string, error) {
	user, err := s.repo.GetBlogUserByUsername(username)
	if err != nil {
		return nil, "", fmt.Errorf("用户名或密码错误")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, "", fmt.Errorf("用户名或密码错误")
	}

	// 检查用户状态
	if user.Status != 1 {
		return nil, "", fmt.Errorf("用户已被禁用")
	}

	// 生成 JWT token
	token, err := s.generateBlogUserToken(user)
	if err != nil {
		return nil, "", fmt.Errorf("生成token失败")
	}

	return &BlogUserInfo{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Nickname: user.Nickname,
		Avatar:   user.AvatarURL,
	}, token, nil
}

// generateBlogUserToken 生成博客用户 JWT token
func (s *Service) generateBlogUserToken(user *model.BlogUser) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
		"exp":      now.Add(time.Duration(s.cfg.JWT.Expiration) * time.Second).Unix(),
		"iat":      now.Unix(),
		"iss":      s.cfg.JWT.Issuer,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWT.Secret))
}

// GetBlogUserByID 根据ID获取用户
func (s *Service) GetBlogUserByID(id uint) (*BlogUserInfo, error) {
	user, err := s.repo.GetBlogUserByID(id)
	if err != nil {
		return nil, err
	}
	return &BlogUserInfo{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Nickname: user.Nickname,
		Avatar:   user.AvatarURL,
	}, nil
}

// LikePost 点赞/取消点赞文章
func (s *Service) LikePost(postSlug string, userID uint) (bool, int64, error) {
	// 获取文章
	post, err := s.repo.GetPostBySlug(postSlug)
	if err != nil {
		return false, 0, err
	}

	// 检查是否已点赞
	like, err := s.repo.GetPostLike(post.ID, userID)
	if err != nil {
		return false, 0, err
	}

	liked := true
	if like != nil {
		// 已点赞，取消
		if err := s.repo.RemovePostLike(post.ID, userID); err != nil {
			return false, 0, err
		}
		liked = false
	} else {
		// 未点赞，添加
		if err := s.repo.AddPostLike(post.ID, userID); err != nil {
			return false, 0, err
		}
		liked = true
	}

	// 获取最新点赞数
	count, err := s.repo.CountPostLikes(post.ID)
	if err != nil {
		return liked, 0, err
	}

	return liked, count, nil
}

// GetPostLikeCount 获取文章点赞数
func (s *Service) GetPostLikeCount(postSlug string) (int64, error) {
	post, err := s.repo.GetPostBySlug(postSlug)
	if err != nil {
		return 0, err
	}
	return s.repo.CountPostLikes(post.ID)
}

// FavoritePost 收藏/取消收藏文章
func (s *Service) FavoritePost(postSlug string, userID uint) (bool, error) {
	post, err := s.repo.GetPostBySlug(postSlug)
	if err != nil {
		return false, err
	}

	fav, err := s.repo.GetPostFavorite(post.ID, userID)
	if err != nil {
		return false, err
	}

	favorited := true
	if fav != nil {
		if err := s.repo.RemovePostFavorite(post.ID, userID); err != nil {
			return false, err
		}
		favorited = false
	} else {
		if err := s.repo.AddPostFavorite(post.ID, userID); err != nil {
			return false, err
		}
		favorited = true
	}

	return favorited, nil
}

// ListMyFavorites 获取我的收藏列表
func (s *Service) ListMyFavorites(userID uint) ([]PostInfo, error) {
	posts, err := s.repo.ListUserFavorites(userID)
	if err != nil {
		return nil, err
	}

	return s.toPostInfoList(posts), nil
}

// HasPassword 检查文章是否有密码保护
func (s *Service) HasPassword(slug string) (bool, error) {
	return s.repo.HasPassword(slug)
}

// VerifyPostPassword 验证文章密码
func (s *Service) VerifyPostPassword(slug, password string) (bool, error) {
	post, err := s.repo.GetPostBySlug(slug)
	if err != nil {
		return false, err
	}

	if post.PasswordHash == "" {
		return true, nil // 没有密码，直接通过
	}

	if err := bcrypt.CompareHashAndPassword([]byte(post.PasswordHash), []byte(password)); err != nil {
		return false, nil
	}
	return true, nil
}

// SetPostPassword 设置文章密码
func (s *Service) SetPostPassword(slug, password string) error {
	post, err := s.repo.GetPostBySlug(slug)
	if err != nil {
		return err
	}

	if password == "" {
		post.PasswordHash = ""
	} else {
		hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		post.PasswordHash = string(hashed)
	}

	return s.repo.UpdatePost(post)
}

// AdminLogInfo 操作日志信息
type AdminLogInfo struct {
	ID         uint   `json:"id"`
	UserID     uint   `json:"user_id"`
	Username   string `json:"username"`
	Action     string `json:"action"`
	TargetType string `json:"target_type"`
	TargetID   *uint  `json:"target_id"`
	TargetName string `json:"target_name"`
	Details    string `json:"details"`
	IP         string `json:"ip"`
	CreatedAt  string `json:"created_at"`
}

// RecordAdminLog 记录管理后台操作
func (s *Service) RecordAdminLog(userID uint, username, action, targetType string, targetID *uint, targetName, details, ip, userAgent string) error {
	log := &model.AdminLog{
		UserID:     userID,
		Username:   username,
		Action:     action,
		TargetType: targetType,
		TargetID:   targetID,
		TargetName: targetName,
		Details:    details,
		IP:         ip,
		UserAgent:  userAgent,
	}
	return s.repo.CreateAdminLog(log)
}

// ListAdminLogs 获取操作日志
func (s *Service) ListAdminLogs(page, pageSize int) ([]AdminLogInfo, int64, error) {
	logs, total, err := s.repo.ListAdminLogs(page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	result := make([]AdminLogInfo, 0, len(logs))
	for _, l := range logs {
		result = append(result, AdminLogInfo{
			ID:         l.ID,
			UserID:     l.UserID,
			Username:   l.Username,
			Action:     l.Action,
			TargetType: l.TargetType,
			TargetID:   l.TargetID,
			TargetName: l.TargetName,
			Details:    l.Details,
			IP:         l.IP,
			CreatedAt:  l.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return result, total, nil
}
