package service

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/go-redis/redis/v8"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
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

type PostInfo struct {
	Slug       string    `json:"slug"`
	Title      string    `json:"title"`
	Date       time.Time `json:"date"`
	Tags       []string  `json:"tags"`
	Categories []string  `json:"categories"`
	Summary    string    `json:"summary"`
	Views      int64     `json:"views"`
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

	result := make([]PostInfo, 0, len(posts))
	for _, p := range posts {
		result = append(result, PostInfo{
			Slug:       p.Slug,
			Title:      p.Title,
			Date:       p.Date,
			Tags:       strings.Split(p.Tags, ","),
			Categories: strings.Split(p.Categories, ","),
			Summary:    p.Summary,
			Views:      p.Views,
		})
	}
	return result, nil
}

// ListPostsPaginated 分页获取已发布文章
func (s *Service) ListPostsPaginated(tag string, page, pageSize int) ([]PostInfo, int64, error) {
	posts, total, err := s.repo.ListPostsPaginated(tag, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	result := make([]PostInfo, 0, len(posts))
	for _, p := range posts {
		result = append(result, PostInfo{
			Slug:       p.Slug,
			Title:      p.Title,
			Date:       p.Date,
			Tags:       strings.Split(p.Tags, ","),
			Categories: strings.Split(p.Categories, ","),
			Summary:    p.Summary,
			Views:      p.Views,
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

	info := &PostInfo{
		Slug:       post.Slug,
		Title:      post.Title,
		Date:       post.Date,
		Tags:       strings.Split(post.Tags, ","),
		Categories: strings.Split(post.Categories, ","),
		Summary:    post.Summary,
		Views:      post.Views + 1,
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
func (s *Service) GetAllTags() (map[string]int, error) {
	return s.repo.GetAllTags()
}

// GetAllCategories returns all categories with post count
func (s *Service) GetAllCategories() (map[string]int, error) {
	return s.repo.GetAllCategories()
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

// RecordVisit 记录访问日志
func (s *Service) RecordVisit(slug, ip, userAgent string, userID int64) {
	visit := &model.Visit{
		PostSlug:  slug,
		IP:        ip,
		UserAgent: userAgent,
		UserID:    userID,
	}
	// 异步记录，不阻塞请求
	go func() {
		s.repo.CreateVisit(visit)
	}()
}
