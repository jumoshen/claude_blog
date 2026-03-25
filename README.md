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
| [blog-frontend](./blog-frontend/) | 前端页面 | 5173 |

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
| [claude-blog-admin](./claude-blog-admin/) | 管理后台 | 3000 |

### Docker Compose

```bash
cd claude-blog-admin
docker-compose up --build
```

## 技术栈

### C 端
- 后端: Gin + GORM + MySQL + Redis
- 前端: Vue 3 + 动态组件

### B 端
- 后端: Gin + GORM + MySQL
- 前端: Vue 3 + Ant Design Vue 4 + ECharts
