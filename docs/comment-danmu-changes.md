# 评论弹幕系统 - 变更记录

## 实现时间
2026-04-02

## 功能概述
在文章详情页新增实时评论弹幕功能，用户可以发送评论，评论以弹幕形式从右向左飘过。同时支持历史评论在进入页面时飘过。

---

## 后端变更

### 1. 数据库模型
**文件**: `frontend-api/internal/model/model.go`

新增 Comment 模型：
```go
type Comment struct {
    gorm.Model
    PostSlug  string `gorm:"size:200;index"`
    UserID    int64  `gorm:"index"`
    Nickname  string `gorm:"size:50"`
    Content   string `gorm:"type:text"`
    IP        string `gorm:"size:50;index"`
    DeviceID  string `gorm:"size:64;index"`
    UserAgent string `gorm:"size:500"`
    Status    int    `gorm:"default:1;index"` // 1=正常 0=待审核 -1=违规
}
```

### 2. 数据库迁移
**文件**: `frontend-api/migrations/004_create_comments.sql`

创建评论表 SQL 迁移文件，包含完整索引设计。

### 3. Repository层
**文件**: `frontend-api/internal/repository/repository.go`

新增方法：
- `CreateComment(comment *model.Comment) error`
- `GetCommentsByPostSlug(postSlug string, limit int) ([]model.Comment, error)`
- `IncrCommentRateLimit(ctx, key, window) (int64, error)` - Redis滑动窗口计数
- `GetCommentRateLimit(ctx, key) (int64, error)`

同时将 `AutoMigrate` 改为只迁移基础表，评论表使用手动迁移。

### 4. Service层
**文件**: `frontend-api/internal/service/service_refactor.go`

新增方法：
- `CreateComment(comment *model.Comment) error`
- `GetCommentsByPostSlug(postSlug string, limit int) ([]model.Comment, error)`
- `CheckCommentRateLimit(ctx, ip, deviceID, userID)` - 频率限制检查
- `ContainsSensitiveWords(content string) bool` - 敏感词检测

**频率限制规则**：
- 匿名用户：IP每分钟3条 + DeviceID每分钟5条
- 登录用户：UserID每10秒1条

### 5. Handler层（新增）
**文件**: `frontend-api/internal/handler/comment.go`

新增 `CommentHandler`：
- `GetComments(c *gin.Context)` - 获取文章评论
- `CreateComment(c *gin.Context)` - 创建评论（含敏感词检测+频率限制）
- `HandleWebSocket(c *gin.Context)` - WebSocket处理

**WebSocket功能**：
- 支持多客户端订阅同一文章slug
- 新评论实时广播到所有订阅者
- 客户端可发送subscribe消息切换订阅

### 6. 路由注册
**文件**: `frontend-api/cmd/server/main.go`

新增路由：
```
GET  /api/v1/comments/:postSlug  - 获取评论
POST /api/v1/comments            - 创建评论
GET  /ws/comments/:postSlug      - WebSocket订阅
```

### 7. 依赖
**文件**: `frontend-api/go.mod`

新增依赖：`github.com/gorilla/websocket v1.5.1`

### 8. 假数据填充
**文件**: `frontend-api/cmd/seed_comments/main.go`

生成每篇文章100-200条假评论，贴合文章内容。

---

## 前端变更

### 1. API接口
**文件**: `blog-frontend/src/api/index.js`

新增接口：
```javascript
getComments(postSlug)   // 获取评论
createComment(data)     // 创建评论
```

### 2. 弹幕图层组件（新增）
**文件**: `blog-frontend/src/components/common/DanmuLayer.vue`

功能：
- 全屏固定定位，pointer-events: none
- 支持1-10条轨道
- 12秒动画从右向左飘过
- 轨道高度60px起，每轨道50px间隔
- 随机背景配色（10种协调颜色）
- 字体颜色可自定义
- 支持字体大小调节
- 队列机制，轨道占满时排队等待

### 3. 弹幕控制面板（新增）
**文件**: `blog-frontend/src/components/common/DanmuControl.vue`

功能：
- 弹幕开关按钮
- 密度滑动条（1-10）
- 字体大小滑动条（12-24px）
- 字体颜色选择器
- 设置保存到localStorage

### 4. 弹幕层集成到App.vue
**文件**: `blog-frontend/src/App.vue`

- DanmuLayer 组件置于 Header 之上，避免遮挡
- 通过 provide/inject 共享弹幕状态和方法
- 支持动态调节密度、字体大小、颜色

### 5. 文章页集成
**文件**: `blog-frontend/src/views/Post.vue`

- WebSocket连接监听新评论
- 获取历史评论并以弹幕形式显示
- 评论提交后即时显示到弹幕层
- 评论表单：昵称+内容

### 6. 工具箱集成
**文件**: `blog-frontend/src/components/common/Toolbox.vue`

- 集成 DanmuControl 组件
- 转发弹幕设置变更到 App.vue

---

## 测试用例

### 1. 评论基础功能
- [ ] 发布匿名评论成功
- [ ] 评论显示在弹幕中
- [ ] 评论出现在评论列表中
- [ ] 评论持久化到数据库

### 2. 弹幕显示
- [ ] 进入文章页，历史评论以弹幕形式飘过
- [ ] 弹幕随机颜色显示
- [ ] 弹幕不显示昵称，只显示内容
- [ ] 弹幕速度12秒从右到左
- [ ] 弹幕不被header遮挡

### 3. 弹幕控制
- [ ] 弹幕开关按钮生效
- [ ] 密度调节生效（轨道数量变化）
- [ ] 字体大小调节生效
- [ ] 字体颜色选择器生效
- [ ] 设置保存到localStorage，刷新后保留

### 4. 实时推送
- [ ] 打开两个浏览器窗口进入同一文章
- [ ] 在一个窗口发评论，另一个窗口实时收到弹幕

### 5. 频率限制
- [ ] 匿名用户IP每分钟超过3条被拒绝
- [ ] 匿名用户DeviceID每分钟超过5条被拒绝
- [ ] 提示"评论太频繁，请稍后再试"

### 6. 敏感词过滤
- [ ] 提交包含敏感词的评论被拒绝
- [ ] 提示"评论包含不当内容，请修改后重试"

### 7. WebSocket
- [ ] WebSocket连接成功
- [ ] 断线自动重连（3秒后）
- [ ] 多标签页同时订阅

### 8. 性能
- [ ] 200条历史评论加载无明显卡顿
- [ ] 弹幕队列机制正常（轨道占满时排队）

---

## 待优化项

1. [ ] 历史弹幕只显示一次，刷新后不再显示
2. [ ] 评论审核功能（后台）
3. [ ] DFA敏感词算法优化
4. [ ] WebSocket心跳检测
5. [ ] 弹幕重复检测（同一人短时间内不重复）
