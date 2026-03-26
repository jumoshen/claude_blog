# 博客系统产品文档

## 1. 产品概述

博客系统分为 C端（展示端）和 B端（管理端）两部分。

| 端 | 说明 |
|---|------|
| C端前台 | 博客文章展示 |
| B端管理后台 | 运营管理 |

## 2. C端功能

### 2.1 文章列表
- 展示已发布的文章列表
- 显示文章标题、日期、标签、浏览数

### 2.2 文章详情
- Markdown 内容渲染（支持代码高亮）
- 显示发布日期、浏览数、标签
- 访问日志自动记录

### 2.3 访问日志
- 自动记录每次文章访问
- 记录 IP、User-Agent、访问时间

## 3. B端功能

### 3.1 认证模块
- 管理端用户登录/登出
- JWT 认证

### 3.2 C端用户管理
- 查看用户列表
- 查看用户详情
- 编辑用户信息（昵称、GitHub信息等）

### 3.3 管理端用户管理
- 管理员 CRUD
- 密码修改

### 3.4 文章管理
- **列表**：分页展示、状态筛选（草稿/已发布/已下架）
- **创建/编辑**：Markdown 编辑器、代码高亮预览、图片上传
- **发布/下架**：一键切换文章状态
- **软删除**：保留数据，仅标记删除
- **Markdown 导入**：批量导入 Markdown 文件
- **模版下载**：下载标准 Markdown 模版
- **数据初始化**：从老数据（study_course）批量导入文章

### 3.5 访问统计
- **概览**：总访问量、今日访问量、热门文章
- **访问记录**：按时间查询详细访问日志
- **文章排行**：按浏览数排序的文章列表
- ECharts 可视化图表展示

### 3.6 文件服务
- 图片上传到 MinIO 对象存储
- 支持 jpg、png、gif、webp、svg 格式
- 最大 10MB
- 公开访问

## 4. 技术架构

### 4.1 项目结构

```
blog/
├── frontend-api/          # C端 API (Go Gin)
├── blog-frontend/        # C端前台 (Vue 3 + Element Plus)
├── claude-blog-admin/     # B端管理后台
│   ├── blog-admin-backend/   # 后端 (Go Gin + GORM)
│   └── blog-admin-frontend/  # 前端 (Vue 3 + Ant Design Vue)
└── file-service/         # MinIO 文件服务
```

### 4.2 部署架构

```
┌─────────────────────────────────────┐
│         Docker Compose              │
├───────┬───────┬───────┬─────────────┤
│ MySQL │ Redis │ MinIO │   Nginx    │
├───────┴───────┴───────┴─────────────┤
│  frontend-api  │  admin-backend    │
└────────────────┴───────────────────┘
```

### 4.3 服务端口

> 端口配置参见 docker-compose.yml

## 5. 待实现功能

- [ ] 文章分类管理（增删改查分类）
- [ ] 文章标签管理（增删改查标签）
- [ ] 文章搜索（标题、内容关键字搜索）
- [ ] 文章置顶/推荐
- [ ] 用户注册/登录（C端）
- [ ] 用户评论功能
- [ ] 文章点赞功能
- [ ] 文章收藏功能
- [ ] 定时发布文章
- [ ] 文章分享功能（生成海报）
- [ ] 文章打赏/捐赠
- [ ] RSS 订阅
- [ ] SEO 优化（sitemap、robots.txt）
- [ ] 文章密码保护
- [ ] 文章阅读量真实统计（防刷）
- [ ] 管理后台操作日志
- [ ] 管理后台角色权限管理 (RBAC)
- [ ] 多语言支持
- [ ] 暗黑模式

## 6. 部署指南

### 7.1 Docker Compose 部署

```bash
docker compose up -d
```

### 7.2 服务说明

**MySQL:** 密码和数据库名在 docker-compose.yml 中配置

**MinIO Console:** 账号密码在 MinIO 容器环境变量中配置

**admin-backend 环境变量:**
- `MINIO_ENDPOINT`: MinIO 服务地址
- `MINIO_ACCESS_KEY`: MinIO Access Key
- `MINIO_SECRET_KEY`: MinIO Secret Key
- `MINIO_BUCKET`: 存储桶名
- `BASE_URL`: 图片访问基础URL
