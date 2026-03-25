# Frontend API (C端API服务)

纯 API 服务，为前端博客项目提供数据支持。

## 技术栈

| 层级 | 技术 |
|------|------|
| 后端框架 | Gin + GORM |
| 数据库 | MySQL + Redis |
| 认证 | JWT |

## 项目结构

```
frontend-api/
├── cmd/server/             # 入口
├── internal/
│   ├── config/             # 配置加载
│   ├── handler/            # HTTP处理器
│   ├── logger/             # 日志
│   ├── middleware/          # 中间件
│   ├── model/              # 数据模型
│   ├── pkg/
│   │   ├── jwt/            # JWT工具
│   │   └── response/       # 统一响应
│   ├── repository/         # 数据访问
│   └── service/            # 业务逻辑
├── db/                     # 数据库迁移SQL
├── config.yaml             # 配置文件
└── Makefile
```

## 快速启动

```bash
# 修改 config.yaml 配置数据库和Redis

go mod tidy
go run ./cmd/server
```

## API

| Method | Path | 说明 |
|--------|------|------|
| GET | /api/v1/site | 站点信息 |
| GET | /api/v1/posts | 文章列表 |
| GET | /api/v1/posts/:slug | 文章详情 |
| GET | /api/v1/archives | 归档列表 |
| GET | /api/v1/about | 关于页面 |
| GET | /api/v1/tags | 标签列表 |
| GET | /api/v1/categories | 分类列表 |
| GET | /api/v1/auth/login | 登录信息 |
| POST | /api/v1/auth/callback | GitHub OAuth回调 |
| POST | /api/v1/auth/logout | 登出 |
| GET | /api/v1/auth/me | 当前用户信息 |
| POST | /api/v1/admin/refresh | 刷新内容 (需认证) |
