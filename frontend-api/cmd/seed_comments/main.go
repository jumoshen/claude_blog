package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Post struct {
	ID    uint
	Slug  string
	Title string
	Tags  string
}

type Comment struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	PostSlug  string
	UserID    int64
	Nickname  string
	Content   string
	IP        string
	DeviceID  string
	UserAgent string
	Status    int
}

var nicknames = []string{
	"技术爱好者", "代码搬运工", "深夜程序员", "Bug猎人", "架构师小王",
	"学习中的小白", "Python狂人", "Go语言粉", "全栈开发者", "DevOps工程师",
	"数据库管理员", "Linux玩家", "前端切图仔", "后端CRUD", "算法工程师",
	"网络安全员", "云计算从业者", "AI研究员", "开源贡献者", "独立开发者",
	"程序媛", "技术总监", "CTO候选人", "资深架构", "十年老兵",
}

var commentTemplates = map[string][]string{
	"default": {
		"写得很好，受益匪浅",
		"终于看懂了，感谢作者",
		"收藏了，准备深入学习",
		"这个思路很清晰",
		"代码简洁有力",
		"学习到了新东西",
		"写得不错，赞一个",
		"很有帮助，谢谢分享",
		"这篇文章解决了我很久的困惑",
		"期待更多相关文章",
		"请问有配套视频吗",
		"请问这个方案有性能问题吗",
		"收藏备用",
		"跟着做了一遍，成功了",
		"有没有更简单的实现方式",
		"这个思路很巧妙",
		"学习了",
		"感谢博主",
		"写得很好",
		"666",
	},
	"expect": {
		"expect脚本确实很强大",
		"这个自动化方案很实用",
		"我用expect实现了服务器批量管理",
		"请问如何处理交互式密码输入",
		"这个教程太及时了",
		"正好需要这个",
		"学到了新技能",
		"请问支持中文输出吗",
		"这个方法比ansible简单多了",
	},
	"certbot": {
		"certbot真的很好用",
		"免费SSL终于搞定了",
		"按照步骤操作成功了",
		"请问如何自动续期",
		"这个证书有效期多久",
		"谢谢博主分享",
		"省了不少钱",
		"网站安全终于有保障了",
	},
	"goroutine": {
		"协程并发控制确实是个难点",
		"这个channel用法很巧妙",
		"终于理解select了",
		"请问如何避免协程泄露",
		"这个方案比WaitGroup好",
		"context的用法很讲究",
		"学到了",
		"请问性能如何",
		"goroutine真方便",
		"终于理解协程了",
	},
	"mysql": {
		"无限极分类终于搞懂了",
		"这个查询效率怎么样",
		"递归查询确实复杂",
		"请问支持多少层嵌套",
		"我喜欢这个设计",
		"很实用的技巧",
		"收藏了",
		"请问有性能对比吗",
		"这个方案在生产环境用过吗",
		"终于理解了",
	},
	"opencv": {
		"人脸识别很有意思",
		"训练集需要多少图片",
		"识别准确率能达到多少",
		"这个很有趣",
		"请问如何提高准确率",
		"跟着做了一遍，效果不错",
		"学习opencv的好教程",
		"级联分类器原理是什么",
		"很实用的教程",
	},
	"pprof": {
		"pprof真是个好东西",
		"终于学会性能分析了",
		"这个火焰图怎么看",
		"请问如何定位内存泄漏",
		"性能优化很有讲究",
		"很实用的profiling工具",
		"CPU占用高怎么排查",
		"感谢分享",
	},
}

func main() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/blog?charset=utf8mb4&parseTime=True"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Disable redis check
	rdb := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})

	var posts []Post
	if err := db.Find(&posts).Error; err != nil {
		panic(err)
	}

	rand.Seed(time.Now().UnixNano())

	totalComments := 0
	for _, post := range posts {
		// Random 100-200 comments per post
		count := rand.Intn(101) + 100

		for i := 0; i < count; i++ {
			comment := generateComment(post, rdb)
			if err := db.Create(&comment).Error; err != nil {
				fmt.Printf("Error inserting comment for %s: %v\n", post.Slug, err)
			}
		}
		totalComments += count
		fmt.Printf("Inserted %d comments for post: %s\n", count, post.Title)
	}

	fmt.Printf("\nTotal comments inserted: %d\n", totalComments)
}

func generateComment(post Post, rdb *redis.Client) Comment {
	// Select template based on post slug or tags
	templates := commentTemplates["default"]
	for key, t := range commentTemplates {
		if contains(post.Slug, key) || contains(post.Tags, key) {
			templates = append(templates, t...)
		}
	}

	// Add some random generic comments
	allTemplates := append(templates, commentTemplates["default"]...)

	content := allTemplates[rand.Intn(len(allTemplates))]
	nickname := nicknames[rand.Intn(len(nicknames))]

	// Generate random IP
	ip := fmt.Sprintf("192.168.%d.%d", rand.Intn(256), rand.Intn(256))

	// Generate device ID
	deviceID := fmt.Sprintf("dev_%d", rand.Int63())

	return Comment{
		PostSlug:  post.Slug,
		UserID:    0,
		Nickname:  nickname,
		Content:   content,
		IP:        ip,
		DeviceID:  deviceID,
		UserAgent: "Mozilla/5.0",
		Status:    1,
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
