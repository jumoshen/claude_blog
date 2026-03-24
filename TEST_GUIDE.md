# 博客系统前后端分离测试用例

## 测试环境准备

### 1. 启动依赖服务
```bash
# 启动 MySQL 和 Redis (如果没有运行)
docker run -d --name mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=123456 mysql:8
docker run -d --name redis -p 6379:6379 redis:7-alpine
```

### 2. 启动后端服务
```bash
cd /Users/mac/code/blog/markdown-blog
go build -o blog ./cmd/server
./blog
# 服务将在 localhost:8080 启动
```

### 3. 启动前端开发服务器
```bash
cd /Users/mac/code/blog/blog-frontend
npm run dev
# 前端将在 localhost:5173 启动
```

---

## 后端 API 测试

### 1. 健康检查
```bash
curl http://localhost:8080/health
# 期望: {"status":"ok"}
```

### 2. 文章列表 (公开)
```bash
curl http://localhost:8080/api/v1/posts
# 期望: {"code":0,"message":"success","data":[...]}
```

### 3. 按标签筛选文章
```bash
curl "http://localhost:8080/api/v1/posts?tag=golang"
# 期望: 返回标签为 golang 的文章列表
```

### 4. 文章详情 (公开)
```bash
curl http://localhost:8080/api/v1/posts/hello-world
# 期望: {"code":0,"message":"success","data":{"post":{...},"content":"<html>..."}}
```

### 5. 归档列表 (公开)
```bash
curl http://localhost:8080/api/v1/archives
# 期望: {"code":0,"message":"success","data":{"2024":[...],"2023":[...]}}
```

### 6. 关于页 (公开)
```bash
curl http://localhost:8080/api/v1/about
# 期望: {"code":0,"message":"success","data":{"content":"<html>..."}}
```

### 7. 标签云 (公开)
```bash
curl http://localhost:8080/api/v1/tags
# 期望: {"code":0,"message":"success","data":{"golang":3,"rust":2,...}}
```

### 8. 分类列表 (公开)
```bash
curl http://localhost:8080/api/v1/categories
# 期望: {"code":0,"message":"success","data":{"后端":5,"前端":3,...}}
```

### 9. 登录页面信息 (公开)
```bash
curl http://localhost:8080/api/v1/auth/login
# 期望: {"code":0,"message":"success","data":{"client_id":"...","callback_url":"...","state":"..."}}
```

### 10. GitHub OAuth 回调 (需要真实 code)
```bash
# 先访问 GitHub 获取 code，然后用 code 交换 token
curl -X POST "http://localhost:8080/api/v1/auth/callback?code=YOUR_CODE"
# 期望: {"code":0,"message":"success","data":{"token":"...","token_type":"Bearer",...,"user":{...}}}
```

### 11. 获取当前用户 (需要认证)
```bash
# 先登录获取 token
TOKEN="your_jwt_token"
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/auth/me
# 期望: {"code":0,"message":"success","data":{"id":...,"login":"...","name":"...","avatar_url":"...","email":"..."}}
```

### 12. 登出 (需要认证)
```bash
TOKEN="your_jwt_token"
curl -X POST -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/auth/logout
# 期望: {"code":0,"message":"success"}
# 之后再用同一个 token 请求 /auth/me 应该返回 401
```

### 13. 刷新内容 (需要认证)
```bash
TOKEN="your_jwt_token"
curl -X POST -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/admin/refresh
# 期望: {"code":0,"message":"success","data":{"message":"Content refreshed"}}
```

### 14. 未认证访问受保护接口
```bash
curl http://localhost:8080/api/v1/auth/me
# 期望: {"code":401,"message":"missing authorization header"}

curl -H "Authorization: Bearer invalid_token" http://localhost:8080/api/v1/auth/me
# 期望: {"code":401,"message":"invalid or expired token"}
```

---

## 前端测试

### 1. 首页
- 访问 http://localhost:5173
- 期望: 显示文章列表，侧边栏显示标签云和分类
- 点击文章卡片跳转到文章详情

### 2. 文章详情页
- 访问 http://localhost:5173/post/hello-world
- 期望: 显示文章标题、元信息、标签、内容

### 3. 归档页
- 访问 http://localhost:5173/archives
- 期望: 按年份分组显示文章列表

### 4. 关于页
- 访问 http://localhost:5173/about
- 期望: 显示关于页面内容

### 5. 登录页
- 访问 http://localhost:5173/login
- 期望: 显示 "Login with GitHub" 按钮
- 点击后跳转到 GitHub 授权页

### 6. 登录后状态
- 完成 GitHub OAuth 后
- 期望: Header 显示用户头像和用户名
- 下拉菜单可登出

---

## JWT 黑名单测试 (踢出登录)

### 测试步骤:
1. 正常登录，获取 token
2. 用 token 访问 /auth/me - 应该成功
3. 调用 /auth/logout 登出
4. 再次用同一个 token 访问 /auth/me - 应该返回 401

---

## 预期结果

| 测试项 | 状态 |
|--------|------|
| 健康检查 | ✅ |
| 公开 API (posts, archives, about, tags, categories) | ✅ |
| 登录流程 (login -> callback -> token) | ✅ |
| 受保护 API (me, logout, admin/refresh) | ✅ |
| Token 黑名单 (踢出登录) | ✅ |
| 前端页面渲染 | ✅ |
| 前端登录状态管理 | ✅ |
