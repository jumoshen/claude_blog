#!/bin/bash
# ===========================================
# Claude Blog 一键部署脚本
# 使用方式: ./deploy-all.sh
# ===========================================

set -e

# 配置
DEPLOY_DIR="$(cd "$(dirname "$0")" && pwd)"
SERVER_USER=${SERVER_USER:-root}
SERVER_HOST=${SERVER_HOST:-47.120.0.121}
SERVER_DEPLOY_PATH="/home/blog-deploy/claude_blog"

# 镜像
REGISTRY=${REGISTRY:-docker.io}
PROJECT_PREFIX=jumoshen

# 颜色
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log() { echo -e "${GREEN}[$(date '+%H:%M:%S')]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }
log_step() { echo -e "${BLUE}[STEP]${NC} $1"; }

# SSH 执行远程命令
ssh_cmd() {
    ssh -o StrictHostKeyChecking=no "$SERVER_USER@$SERVER_HOST" "$1"
}

# 主流程
main() {
    log "=========================================="
    log " Claude Blog 一键部署"
    log "=========================================="

    # Step 1: Git 拉取最新代码
    log_step "1. 拉取最新代码..."
    ssh_cmd "cd $SERVER_DEPLOY_PATH && git fetch origin && git checkout main && git pull origin main"

    # Step 2: 复制配置文件 (从 deploy 目录)
    log_step "2. 配置服务..."
    # C端 API 配置
    log "   配置 C端 API..."
    ssh_cmd "cp $SERVER_DEPLOY_PATH/deploy/config-api.yaml $SERVER_DEPLOY_PATH/frontend-api/config.yaml"
    # B端 Backend 配置
    log "   配置 B端 Backend..."
    ssh_cmd "cp $SERVER_DEPLOY_PATH/deploy/config-admin.yaml $SERVER_DEPLOY_PATH/claude-blog-admin/blog-admin-backend/config.yaml"
    # Nginx 配置
    log "   配置 Nginx..."
    ssh_cmd "cp $SERVER_DEPLOY_PATH/deploy/nginx.conf $SERVER_DEPLOY_PATH/nginx.conf"

    # Step 3: 构建镜像
    log_step "3. 构建 Docker 镜像..."

    # 确保网络存在
    ssh_cmd "docker network create blogdeploy_blog-network 2>/dev/null || true"

    # 构建 C端
    log "   构建 C端 API..."
    ssh_cmd "cd $SERVER_DEPLOY_PATH/frontend-api && docker build -t ${REGISTRY}/${PROJECT_PREFIX}/claude_blog_api:latest ."
    log "   构建 C端 Frontend..."
    ssh_cmd "cd $SERVER_DEPLOY_PATH/blog-frontend && docker build -t ${REGISTRY}/${PROJECT_PREFIX}/claude_blog_frontend:latest ."

    # 构建 B端
    log "   构建 B端 Backend..."
    ssh_cmd "cd $SERVER_DEPLOY_PATH/claude-blog-admin/blog-admin-backend && docker build -t ${REGISTRY}/${PROJECT_PREFIX}/claude_blog_admin_backend:latest ."
    log "   构建 B端 Frontend..."
    ssh_cmd "cd $SERVER_DEPLOY_PATH/claude-blog-admin/blog-admin-frontend && docker build -t ${REGISTRY}/${PROJECT_PREFIX}/claude_blog_admin_frontend:latest ."

    # Step 4: 停止并删除所有旧服务
    log_step "4. 停止并删除旧服务..."
    ssh_cmd "cd $SERVER_DEPLOY_PATH && docker compose down 2>/dev/null || true"
    ssh_cmd "docker ps -a --format '{{.Names}}' | grep -E '^cc-blog|^nginx-proxy' | xargs -r docker stop 2>/dev/null || true"
    ssh_cmd "docker ps -a --format '{{.Names}}' | grep -E '^cc-blog|^nginx-proxy' | xargs -r docker rm 2>/dev/null || true"

    # Step 5: 启动新服务
    log_step "5. 启动服务..."
    ssh_cmd "cd $SERVER_DEPLOY_PATH && docker compose up -d"

    # Step 6: 等待服务启动
    log_step "6. 等待服务就绪..."
    sleep 5

    # Step 7: 确保 MySQL 用户存在
    log_step "7. 初始化数据库..."
    ssh_cmd "docker exec cc-blog-mysql mysql -u root -pBlog2024Secure -e \"CREATE USER IF NOT EXISTS 'blog'@'%' IDENTIFIED BY 'Blog2024Secure'; GRANT ALL PRIVILEGES ON blog.* TO 'blog'@'%'; FLUSH PRIVILEGES;\" 2>/dev/null || true"

    # Step 8: 检查状态
    log_step "8. 检查服务状态..."
    echo ""
    ssh_cmd "docker ps --format 'table {{.Names}}\t{{.Status}}' | grep -E 'cc-blog|nginx'"
    echo ""

    # Step 9: 测试访问
    log_step "9. 测试访问..."
    HTTP_CODE=$(ssh_cmd "curl -s -o /dev/null -w '%{http_code}' http://localhost:3000/")
    if [ "$HTTP_CODE" = "200" ]; then
        log "   C端前台: http://jumoshen.cn - OK"
    else
        log_error "   C端前台: http://jumoshen.cn - Failed (HTTP $HTTP_CODE)"
    fi

    HTTP_CODE=$(ssh_cmd "curl -s -o /dev/null -w '%{http_code}' http://localhost:3001/")
    if [ "$HTTP_CODE" = "200" ]; then
        log "   B端后台: http://back.jumoshen.cn - OK"
    else
        log_error "   B端后台: http://back.jumoshen.cn - Failed (HTTP $HTTP_CODE)"
    fi

    HTTPS_CODE=$(ssh_cmd "curl -s -o /dev/null -w '%{http_code}' https://jumoshen.cn -k")
    if [ "$HTTPS_CODE" = "200" ]; then
        log "   HTTPS: https://jumoshen.cn - OK"
    else
        log_error "   HTTPS: https://jumoshen.cn - Failed (HTTP $HTTPS_CODE)"
    fi

    log ""
    log "=========================================="
    log " 部署完成!"
    log "=========================================="
}

main "$@"
