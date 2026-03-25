# Markdown Blog (C端博客)

基于 Markdown 文件的静态博客系统，提供文章列表、详情、分类、标签等功能。

## 技术栈

| 层级 | 技术 |
|------|------|
| 后端框架 | Gin + GORM |
| 前端 | 原生 HTML/JS |
| 数据库 | MySQL + Redis |
| 认证 | JWT |

## 项目结构

```
markdown-blog/
├── cmd/server/             # 入口
├── internal/
│   ├── config/             # 配置加载
│   ├── handler/            # HTTP处理器
│   ├── logger/             # 日志 (lumberjack)
│   ├── middleware/          # 中间件
│   ├── model/              # 数据模型
│   ├── pkg/
│   │   ├── jwt/            # JWT工具
│   │   └── response/       # 统一响应
│   ├── repository/         # 数据访问
│   └── service/            # 业务逻辑
├── posts/                  # Markdown文章目录
├── config.yaml             # 配置文件
└── main.go
```

## 快速启动

```bash
# 修改 config.yaml 配置数据库和Redis

go mod tidy
go run ./cmd/server
```

服务启动后访问 http://localhost:8080

## 文章格式

```markdown
---
title: "文章标题"
date: 2024-01-01T10:00:00+08:00
draft: false
tags: ["tag1", "tag2"]
categories: ["Category"]
keywords: "关键词"
---

文章标题<!--more-->

正文内容...
```

## API

| Method | Path | 说明 |
|--------|------|------|
| GET | / | 首页 |
| GET | /posts | 文章列表 |
| GET | /posts/:slug | 文章详情 |
| GET | /tags | 标签列表 |
| GET | /tags/:tag | 标签下的文章 |
| GET | /categories | 分类列表 |
| GET | /categories/:category | 分类下的文章 |
| POST | /api/auth/login | 登录 |
| POST | /api/auth/logout | 登出 |
