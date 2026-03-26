#!/bin/bash
set -e

# ============================================
# Claude Blog 部署脚本 (服务器端执行)
# ============================================

DEPLOY_DIR="/home/blog-deploy"
BLOG_REPO="git@github.com:jumoshen/claude_blog.git"
ADMIN_REPO="git@github.com:jumoshen/claude-blog-admin.git"

echo "=== 1. 克隆代码 ==="
cd $DEPLOY_DIR
git clone $BLOG_REPO temp_build
git clone $ADMIN_REPO temp_admin

echo "=== 2. 修改架构为 amd64 ==="
# 修改 frontend-api 的 Dockerfile
sed -i 's/GOARCH=arm64/GOARCH=amd64/g' temp_build/frontend-api/Dockerfile 2>/dev/null || true

echo "=== 3. 构建 C 端 API ==="
cd $DEPLOY_DIR/temp_build/frontend-api
docker build -t jumoshen/claude_blog_api:latest .

echo "=== 4. 构建 B 端后端 ==="
cd $DEPLOY_DIR/temp_admin/blog-admin-backend
docker build -t jumoshen/claude_blog_admin_backend:latest .

echo "=== 5. 构建 C 端前端 ==="
cd $DEPLOY_DIR/temp_build/blog-frontend
cat > nginx.conf << 'EOF'
server {
    listen 80;
    location / {
        root /usr/share/nginx/html;
        try_files $uri $uri/ /index.html;
    }
}
EOF
cat > Dockerfile << 'EOF'
FROM nginx:alpine
COPY dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EOF
docker build -t jumoshen/claude_blog_frontend:latest .

echo "=== 6. 构建 B 端前端 ==="
cd $DEPLOY_DIR/temp_admin/blog-admin-frontend
cat > nginx.conf << 'EOF'
server {
    listen 80;
    location / {
        root /usr/share/nginx/html;
        try_files $uri $uri/ /index.html;
    }
}
EOF
cat > Dockerfile << 'EOF'
FROM nginx:alpine
COPY dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EOF
docker build -t jumoshen/claude_blog_admin_frontend:latest .

echo "=== 7. 停止旧容器 ==="
cd $DEPLOY_DIR
docker-compose down 2>/dev/null || true
docker stop cc-blog-api cc-blog-admin-backend cc-blog-web cc-blog-admin-web 2>/dev/null || true
docker rm cc-blog-api cc-blog-admin-backend cc-blog-web cc-blog-admin-web 2>/dev/null || true

echo "=== 8. 启动容器 ==="
docker-compose up -d

echo "=== 9. 连接容器到 dnmp_default 网络 ==="
docker network connect dnmp_default cc-blog-web 2>/dev/null || true
docker network connect dnmp_default cc-blog-api 2>/dev/null || true
docker network connect dnmp_default cc-blog-admin-web 2>/dev/null || true
docker network connect dnmp_default cc-blog-admin-backend 2>/dev/null || true

echo "=== 10. 重载 nginx ==="
docker exec nginx nginx -s reload 2>/dev/null || true

echo "=== 11. 清理临时文件 ==="
rm -rf $DEPLOY_DIR/temp_build $DEPLOY_DIR/temp_admin

echo "=== 部署完成 ==="
echo "C端: https://jumoshen.cn"
echo "B端: https://back.jumoshen.cn"
