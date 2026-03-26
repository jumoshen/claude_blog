# TODO List

## C端 (blog-frontend)

### [ ] 1. Logo 与多风格切换
- [x] 添加三个 logo：像素巨魔、可爱巨魔、Q版巨魔 (已放到 public/)
- [x] Header 右侧添加风格切换下拉菜单
- [x] 三种风格配色和样式调整 (store 已实现)
- [ ] Docker 部署待验证

### [ ] 2. Tags 点击跳转
- [x] Tag 点击跳转到文章列表页
- [x] 文章列表页按 tag 筛选
- [ ] 待验证

### [ ] 3. Archives 样式优化
- [x] 改进 Archives 区块配色和样式
- [ ] 待验证

### [ ] 4. About 页面修复
- [x] About.vue 添加空状态提示
- [ ] 待验证 (需要创建 content/about/index.md)

## B端 (claude-blog-admin)

### [ ] 5. Dashboard 访问趋势图
- [x] 支持按分钟/小时/天查看
- [x] 折线图展示
- [ ] 待验证

### [ ] 6. 访问统计页面趋势图
- [x] 同 Dashboard，支持分钟/小时/天
- [ ] 待验证

### [ ] 7. 用户管理字段完善
- [x] User model 添加 GitHub 相关字段
- [x] Users.vue 显示 GitHub 登录名
- [ ] 待验证

### [ ] 8. 清理旧 docker-compose
- [ ] claude-blog-admin 目录下旧 docker-compose.yml

## 已完成

- [x] B端图片上传到 OSS (MinIO)
- [x] v-md-editor 图片上传功能
- [x] Markdown 预览和代码高亮

## 阻塞项

- Docker Hub 网络问题，构建反复超时
