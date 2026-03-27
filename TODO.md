# TODO List

## C端 (blog-frontend)

### [x] 1. Logo 与多风格切换
- [x] 添加三个 logo：像素巨魔、可爱巨魔、Q版巨魔 (已放到 public/)
- [x] Header 右侧添加风格切换下拉菜单
- [x] 三种风格配色和样式调整 (store 已实现)
- [x] Docker 部署配置正确，待实际部署验证

### [x] 2. Tags 点击跳转
- [x] Tag 点击跳转到文章列表页
- [x] 文章列表页按 tag 筛选
- [x] 已验证

### [x] 3. Archives 样式优化
- [x] 改进 Archives 区块配色和样式
- [x] 已验证

### [x] 4. About 页面修复
- [x] About.vue 添加空状态提示
- [x] 已验证 (需要创建 content/about/index.md 自定义内容)

## B端 (claude-blog-admin)

### [x] 5. Dashboard 访问趋势图
- [x] 支持按分钟/小时/天查看
- [x] 折线图展示
- [x] 已验证

### [x] 6. 访问统计页面趋势图
- [x] 同 Dashboard，支持分钟/小时/天
- [x] 已验证

### [x] 7. 用户管理字段完善
- [x] User model 添加 GitHub 相关字段
- [x] Users.vue 显示 GitHub 登录名
- [x] 已验证

### [x] 8. 清理旧 docker-compose
- [x] claude-blog-admin 目录下无 docker-compose.yml (已清理)

## 已完成

- [x] B端图片上传到 OSS (MinIO)
- [x] v-md-editor 图片上传功能
- [x] Markdown 预览和代码高亮

## 阻塞项

- Docker Hub 网络问题，构建反复超时
