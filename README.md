# Claude Blog

博客系统，包含 C 端（前后端）和 B 端管理后台。

## 项目结构

```
├── frontend-api/       # C端 API 服务 (Gin)
├── blog-frontend/      # C端前端 (Vue 3)
└── claude-blog-admin/   # B端管理后台 (Vue 3 + Gin)
```

## C 端

| 项目 | 说明 | 端口 |
|------|------|------|
| [frontend-api](./frontend-api/) | API 服务 | 8080 |
| [blog-frontend](./blog-frontend/) | 前端页面 | 3000 |

### 快速启动

```bash
# 启动后端
cd frontend-api
go mod tidy
go run ./cmd/server

# 启动前端
cd blog-frontend
npm install
npm run dev
```

## B 端

| 项目 | 说明 | 端口 |
|------|------|------|
| [claude-blog-admin](./claude-blog-admin/) | 管理后台 | 3001 |

### AI 助手功能

B 端文章编辑器深度集成 AI 能力（MiniMax 模型）：

- **选中操作**：选中文字后可直接 AI 润色/解释/扩展
- **AI 续写**：在光标位置续写内容
- **AI 生成标题**：根据文章内容自动生成吸引人的标题
- **Token 节省**：选中操作只发送选中文案，续写只发送最后 2000 字符
- **预览模式**：AI 结果以原文/结果对比展示，确认后应用

### Docker Compose (统一部署)

所有服务统一管理在根目录 `docker-compose.yml`：

```bash
docker-compose up --build
```

| 服务 | 端口 |
|------|------|
| MySQL | 3307 |
| Redis | 6380 |
| C端 API | 8080 |
| C端 前端 | 3000 |
| B端 后端 | 8082 |
| B端 前端 | 3001 |

## 技术栈

### C 端
- 后端: Gin + GORM + MySQL + Redis
- 前端: Vue 3 + 动态组件

### B 端
- 后端: Gin + GORM + MySQL
- 前端: Vue 3 + Ant Design Vue 4 + ECharts
