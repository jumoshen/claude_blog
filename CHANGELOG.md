# 更新日志

## v1.6.0 (2026-04-10)

### C端功能
- [x] T1-T2: 文章分类/标签管理（前端已实现）
- [x] T3-T9: 后端API已完成（文章搜索、置顶/推荐、定时发布、目录、上下篇导航、阅读时间估算、热门/相关文章推荐）
- [x] T10-T13: 用户注册/登录、评论、点赞、收藏（前端已实现）
- [x] T14-T16: 文章分享海报、打赏/捐赠、RSS订阅（前端已实现）
- [x] T23: 键盘快捷键导航（j/k 上下篇）
- [x] T24: 代码块复制按钮
- [x] T25: 图片懒加载
- [x] T26: 页面切换平滑过渡动画
- [x] T27: 阅读进度保存（localStorage）
- [x] T28: PWA 支持
- [x] T21: 暗黑模式
- [ ] T22: 多语言支持

### B端功能
- [x] T19: 管理后台操作日志
- [x] T20: 管理后台 RBAC 权限管理
- [x] 分类管理 - CRUD + 层级分类支持
- [x] 标签管理 - CRUD + 颜色选择器
- [x] 文章编辑器插入代码块 - 支持 20 种编程语言（golang/python/javascript/shell 等）
- [x] AI 助手 - 润色/解释/扩展/续写/生成标题

### 修复
- [x] 标签颜色不显示问题
- [x] 标签计数为空问题
- [x] 分类筛选 SQL 错误
- [x] 文章编辑器下拉选择失焦问题
- [x] 暗黑模式下文章详情页和 Footer 样式

---

## v1.5.0 (2026-04-09)

### 镜像构建优化

- [x] C端 (frontend-api, blog-frontend) 全面切换到 GitHub Actions 构建
- [x] 使用 GitHub Container Registry (ghcr.io) 托管镜像
- [x] 解决服务器 1.7GB 内存 Docker 构建 OOM 问题

### 部署修复

- [x] docker-compose.yml B端改用 ghcr.io 预构建镜像
- [x] deploy-all.sh 改用 GitHub Container Registry 拉取所有镜像
- [x] 修复 config-admin.yaml 缺少 AI 配置导致 AI 接口返回空内容
- [x] deploy-all.sh 添加 `docker compose pull` 步骤确保使用最新镜像

### 文档

- [x] 添加 AI 功能说明到产品文档

---

## v1.4.0 (2026-04-08)

### AI 助手功能
- [x] B端文章编辑器深度集成 AI 能力
- [x] 支持选中文字后直接 AI 操作（润色/解释/扩展）
- [x] AI 续写功能 - 在光标位置续写内容
- [x] AI 生成标题 - 根据文章内容自动生成吸引人的标题
- [x] Token 节省策略 - 选中操作只发送选中文案，续写只发送最后 2000 字符
- [x] 预览模式 - AI 结果以原文/结果对比展示，用户确认后才应用
- [x] SSE 流式输出 - 实时展示 AI 生成进度

### 后端接口
- [x] `POST /api/admin/ai/process_selection` - 处理选中文字（润色/解释/扩展）
- [x] `POST /api/admin/ai/continue_writing` - 续写内容
- [x] `POST /api/admin/ai/generate_title` - 生成标题
- [x] MiniMax 模型支持

### 数据库
- [x] 新增 `ai_generations` 表 - 存储 AI 创作历史

---

## v1.3.0 (2026-04-01)

### SEO优化
- [x] 添加 sitemap.xml 百度/Google收录
- [x] robots.txt 正确配置 Sitemap 地址
- [x] 百度站点验证 (jumoshen.cn)

### 图片上传修复
- [x] B端上传图片URL改为 `https://jumoshen.cn/uploads/xxx`
- [x] C端nginx添加 `/uploads/` 代理到MinIO
- [x] B端新增 `public_base_url` 配置支持独立域名

### 部署修复
- [x] deploy-all.sh 添加B端仓库正确remote检查
- [x] 修复服务器B端仓库git remote指向错误问题

### 幸运抽奖工具
- [x] C端新增可收起工具箱（hover展开）
- [x] 骰子工具 - 随机跳转文章
- [x] 翻牌工具 - 3张扑克牌，3D翻转动画
- [x] 牌背面/正面使用设计图片，随机选取
- [x] 三种主题适配（像素风/可爱风/Q版）

---

## [进行中] 项目重构 (2026-03-24)

### 环境配置
- **MySQL**: `127.0.0.1:3306` (root/123456) - 宿主机 Docker
- **Redis**: `127.0.0.1:6379` - 宿主机 Docker
- docker-compose 中不内置 MySQL/Redis，使用宿主机已有的服务

### 后端重构
- [ ] 前后端分离架构
- [ ] JWT 认证 + Redis 缓存 token 实现登录态管理与踢人功能
- [ ] 代码风格规范 + 模块拆分 (handler/service/repository)
- [ ] Swagger 文档生成 (make swagger)
- [ ] 配置外置化 (config.yaml)

### 前端重构
- [ ] Vue.js 框架
- [ ] 组件化拆分
- [ ] 用户信息本地缓存 (localStorage)

### DevOps
- [x] 前后端分别 Dockerfile
- [x] docker-compose 统一管理
- [ ] 支持 enable_proxy 代理功能

### 测试
- [x] Docker 部署验证通过
- [x] 后端单元测试覆盖 (config 92.3%, response 100%, logger 84.7%, middleware 95.0%, service 34.4%)
- [ ] 前后端联调测试

---

## v1.2.0 (2026-03-24) - markdown-blog

### Bug 修复
- **数据库初始化修复** - 修复每次启动服务时 DropTable 导致数据被清空的问题，改为只做 AutoMigrate 迁移
- **首页侧边栏** - 标签云和分类改为瀑布流展示在首页右侧

### 新增功能
- **首页侧边栏** - 文章列表右侧展示标签云和分类
- **SQL 文件** - 添加 `db/v1.0.0_init.sql` 数据库初始化脚本

---

## v1.1.0 (2026-03-24) - markdown-blog

### 新增功能
- **标签云页面** - 访问 `/tags` 查看所有标签，标签颜色根据文章数量动态生成
- **标签筛选** - 点击任意标签可筛选该标签下的所有文章，URL: `/tags?name=xxx`
- **首页标签可点击** - 文章列表中的标签可直接点击跳转
- **文章页标签可点击** - 文章详情页的标签可点击跳转

### 功能修复
- **归档页面年份统计** - 修复"共 X 个年份"显示错误的问题

### 样式优化
- **全新配色方案** - 从紫色改为清爽的蓝绿色系 (Cyan/Teal)
- **首页美化** - 添加动画渐变背景、浮动粒子效果、玻璃态卡片
- **归档页美化** - 时间线样式优化、卡片悬停动效
- **关于页美化** - 头像脉冲动画、玻璃态内容卡片
- **登录页美化** - 统一视觉风格
- **文章页美化** - 阅读进度条、代码块样式优化

### 日志系统
- 日志输出改为 JSON 格式
- 日志写入到 `logs/blog.log` 文件
- Handler 中的日志改用项目 logger

---

## v1.0.0 (2026-03-23) - markdown-blog

### 初始功能
- 博客首页 - 展示文章列表
- 文章详情页 - Markdown 渲染、代码高亮、阅读量统计
- 归档页面 - 按年份整理文章
- 关于页面 - 个人介绍
- GitHub OAuth 登录认证
- Hugo 内容解析 - 自动解析 Markdown 文件的 front matter
- MySQL 数据存储
- 响应式设计 (Tailwind CSS)
