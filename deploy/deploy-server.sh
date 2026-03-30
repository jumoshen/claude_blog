#!/bin/bash
set -e

# ============================================
# Claude Blog 部署脚本 (服务器端执行)
# ============================================

DEPLOY_DIR="/home/blog-deploy"
BLOG_REPO="git@github.com:jumoshen/claude_blog.git"
ADMIN_REPO="git@github.com:jumoshen/claude-blog-admin.git"

echo "=== 1. 清理旧代码 ==="
cd $DEPLOY_DIR
rm -rf temp_build temp_admin 2>/dev/null || true

echo "=== 2. 克隆代码 ==="
git clone $BLOG_REPO temp_build
git clone $ADMIN_REPO temp_admin

echo "=== 3. 复制配置文件 ==="
# C端 API 配置
cp $DEPLOY_DIR/config-api.yaml temp_build/frontend-api/config.yaml
# B端后端配置
cp $DEPLOY_DIR/config-admin.yaml temp_admin/blog-admin-backend/config.yaml

echo "=== 4. 构建 C 端前端 ==="
cd $DEPLOY_DIR/temp_build/blog-frontend
# 使用生产 API 地址
sed -i 's|http://localhost:8080|https://jumoshen.cn|g' .env 2>/dev/null || true
npm install
npm run build
docker build -t jumoshen/claude_blog_frontend:latest .

echo "=== 5. 构建 C 端 API ==="
cd $DEPLOY_DIR/temp_build/frontend-api
docker build -t jumoshen/claude_blog_api:latest .

echo "=== 6. 构建 B 端前端 ==="
cd $DEPLOY_DIR/temp_admin/blog-admin-frontend
npm install
npm run build
docker build -t jumoshen/claude_blog_admin_frontend:latest .

echo "=== 7. 构建 B 端后端 ==="
cd $DEPLOY_DIR/temp_admin/blog-admin-backend
docker build -t jumoshen/claude_blog_admin_backend:latest .

echo "=== 8. 停止旧容器 ==="
docker stop cc-blog-api cc-blog-admin-backend cc-blog-web cc-blog-admin-web 2>/dev/null || true
docker rm cc-blog-api cc-blog-admin-backend cc-blog-web cc-blog-admin-web 2>/dev/null || true

echo "=== 9. 启动容器 ==="
cd $DEPLOY_DIR
docker-compose up -d

echo "=== 10. 清理临时文件 ==="
rm -rf $DEPLOY_DIR/temp_build $DEPLOY_DIR/temp_admin

echo "=== 部署完成 ==="
echo "C端: https://jumoshen.cn"
echo "B端: https://back.jumoshen.cn"
