---
title: "使用 GitHub Actions 构建 Docker 镜像：解决服务器资源限制问题"
date: 2026-04-09
tags: ["GitHub Actions", "Docker", "CI/CD", "DevOps"]
categories: ["技术实践"]
---

## 背景

在部署博客项目时遇到了一个典型问题：**生产服务器内存只有 1.7GB**，使用 Docker 在服务器上构建镜像时，Go 编译器经常被 OOM Kill 导致构建失败。

传统的解决方案是使用更高配置的服务器，但这意味着更高的成本。有没有更优雅的解决方案？

答案是：**将镜像构建放到 GitHub Actions 上完成，利用 GitHub 提供的计算资源，服务器只负责运行容器**。

## 解决方案

### 1. 使用 GitHub Container Registry (ghcr.io)

GitHub Packages 提供了免费的容器镜像托管服务，每个开源项目都有充足的存储空间。

```yaml
# 登录 GitHub Container Registry
- uses: docker/login-action@v3
  with:
    registry: ghcr.io
    username: ${{ github.actor }}
    password: ${{ secrets.GITHUB_TOKEN }}
```

### 2. 配置 GitHub Actions 工作流

在项目中创建 `.github/workflows/docker-build.yml`：

```yaml
name: Build and Push Docker Images

on:
  push:
    branches: [main, feature/**]

jobs:
  build:
    runs-on: ubuntu-latest

    permissions:
      contents: read
      packages: write

    steps:
      - uses: actions/checkout@v4

      - uses: docker/setup-buildx-action@v3

      - uses: docker/login-action@v3
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - uses: docker/build-push-action@v5
        with:
          context: ./frontend-api
          push: true
          tags: |
            ghcr.io/${{ github.repository_owner }}/claude-blog/frontend-api:${{ github.sha }}
            ghcr.io/${{ github.repository_owner }}/claude-blog/frontend-api:latest
```

### 3. 修改 docker-compose.yml 使用 ghcr.io 镜像

```yaml
# 之前：本地构建
frontend-api:
  build:
    context: ./frontend-api
    dockerfile: Dockerfile
  image: jumoshen/claude_blog_api:latest

# 之后：使用 GitHub 构建的镜像
frontend-api:
  image: ghcr.io/jumoshen/claude-blog/frontend-api:latest
```

### 4. 更新部署脚本

```bash
# 拉取 GitHub 构建的镜像，而不是本地构建
docker pull ghcr.io/jumoshen/claude-blog/frontend-api:latest
docker pull ghcr.io/jumoshen/claude-blog/blog-frontend:latest
```

## 完整工作流

```
┌─────────────────────────────────────────────────────────────┐
│                        开发者                                │
│                          │                                  │
│                          ▼                                  │
│                    git push                                  │
│                          │                                  │
└──────────────────────────┼──────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                    GitHub Actions                           │
│                                                              │
│   ┌──────────────┐    ┌──────────────┐    ┌──────────────┐ │
│   │ checkout code│───▶│ build image  │───▶│ push to ghcr │ │
│   └──────────────┘    └──────────────┘    └──────────────┘ │
│                                                              │
│   GitHub 提供 Ubuntu 最新版 + 7GB 内存                       │
└─────────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                    GitHub Container Registry                │
│                                                              │
│   ghcr.io/jumoshen/claude-blog/frontend-api:latest         │
│   ghcr.io/jumoshen/claude-blog/blog-frontend:latest        │
└─────────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                      生产服务器                              │
│                                                              │
│   docker pull ghcr.io/...  # 只拉取，不构建                 │
│   docker compose up -d      # 启动容器                       │
│                                                              │
│   服务器只需 1GB 内存！                                      │
└─────────────────────────────────────────────────────────────┘
```

## 优势

1. **零构建成本**：GitHub Actions 提供免费计算资源
2. **构建速度快**：GitHub 提供高性能构建机器，比服务器快 5-10 倍
3. **服务器资源利用率低**：1GB 内存足以运行容器
4. **构建版本可追溯**：每次 commit 都有对应的镜像版本
5. **自动化程度高**：push 代码自动构建，无需人工干预

## 注意事项

1. **隐私保护**：如果镜像包含敏感信息，建议设置仓库为 private
2. **网络速度**：服务器拉取镜像速度取决于网络质量
3. **构建缓存**：合理利用 Docker 层缓存可以加速构建

## 结语

将 Docker 镜像构建放到 CI/CD 平台是云原生时代的最佳实践，不仅解决了服务器资源限制问题，还提升了部署的自动化程度和可维护性。
