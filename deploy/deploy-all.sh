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

# GitHub Container Registry (用于 B 端)
GHCR_REGISTRY=ghcr.io
GHCR_OWNER=jumoshen

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
    # B端仓库是独立的，需要单独拉取并检查remote是否正确
    log "   检查 B端仓库..."
    ssh_cmd "cd $SERVER_DEPLOY_PATH/claude-blog-admin && git remote set-url origin git@github.com:jumoshen/claude-blog-admin.git 2>/dev/null || true && git pull origin main"

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

    # Step 3: 登录 GitHub Container Registry
    log_step "3. 登录 GitHub Container Registry..."
    ssh_cmd "echo ${GITHUB_TOKEN} | docker login ${GHCR_REGISTRY} -u ${GHCR_OWNER} --password-stdin"

    # Step 4: 拉取 GitHub 构建的镜像
    log_step "4. 拉取镜像..."

    # 确保网络存在
    ssh_cmd "docker network create blogdeploy_blog-network 2>/dev/null || true"

    # C端镜像
    log "   拉取 C端 API..."
    ssh_cmd "docker pull ${GHCR_REGISTRY}/${GHCR_OWNER}/claude-blog/frontend-api:latest"
    log "   拉取 C端 Frontend..."
    ssh_cmd "docker pull ${GHCR_REGISTRY}/${GHCR_OWNER}/claude-blog/blog-frontend:latest"

    # B端镜像
    log "   拉取 B端 Backend..."
    ssh_cmd "docker pull ${GHCR_REGISTRY}/${GHCR_OWNER}/claude-blog-admin/admin-backend:latest"
    log "   拉取 B端 Frontend..."
    ssh_cmd "docker pull ${GHCR_REGISTRY}/${GHCR_OWNER}/claude-blog-admin/admin-frontend:latest"

    # Step 6: 停止并删除所有旧服务
    log_step "4. 停止并删除旧服务..."
    ssh_cmd "cd $SERVER_DEPLOY_PATH && docker compose down 2>/dev/null || true"
    ssh_cmd "docker ps -a --format '{{.Names}}' | grep -E '^cc-blog|^nginx-proxy' | xargs -r docker stop 2>/dev/null || true"
    ssh_cmd "docker ps -a --format '{{.Names}}' | grep -E '^cc-blog|^nginx-proxy' | xargs -r docker rm 2>/dev/null || true"

    # Step 7: 启动新服务
    log_step "5. 启动服务..."
    ssh_cmd "cd $SERVER_DEPLOY_PATH && docker compose up -d"

    # Step 8: 等待服务启动
    log_step "6. 等待服务就绪..."
    sleep 5

    # Step 9: 确保 MySQL 用户存在
    log_step "7. 初始化数据库..."
    ssh_cmd "docker exec cc-blog-mysql mysql -u root -pBlog2024Secure -e \"CREATE USER IF NOT EXISTS 'blog'@'%' IDENTIFIED BY 'Blog2024Secure'; GRANT ALL PRIVILEGES ON blog.* TO 'blog'@'%'; FLUSH PRIVILEGES;\" 2>/dev/null || true"

    # Step 10: 检查状态
    log_step "8. 检查服务状态..."
    echo ""
    ssh_cmd "docker ps --format 'table {{.Names}}\t{{.Status}}' | grep -E 'cc-blog|nginx'"
    echo ""

    # Step 11: 测试访问
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
