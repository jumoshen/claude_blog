#!/bin/bash
# ===========================================
# Claude Blog 生产环境部署脚本
# 使用方式: ./deploy.sh [选项]
# ===========================================

set -e

# 配置
REGISTRY=${REGISTRY:-docker.io}
PROJECT_PREFIX=jumoshen
COMPOSE_FILE="docker-compose.yml"
ENV_FILE=".env"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查必要命令
check_commands() {
    if ! command -v docker &> /dev/null; then
        log_error "Docker 未安装，请先安装 Docker"
        exit 1
    fi

    if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
        log_error "Docker Compose 未安装"
        exit 1
    fi
}

# 检查环境变量文件
check_env() {
    if [ ! -f "$ENV_FILE" ]; then
        log_warn "环境变量文件 $ENV_FILE 不存在，尝试创建示例..."
        cat > "$ENV_FILE.example" << 'EOF'
# 数据库配置
DB_HOST=mysql
DB_PORT=3306
DB_USER=blog
DB_PASSWORD=your_password_here
DB_NAME=blog

# MySQL Root 密码
MYSQL_ROOT_PASSWORD=root_password_here

# Redis
REDIS_HOST=redis
REDIS_PORT=6379

# JWT
JWT_SECRET=your_jwt_secret_here

# MinIO
MINIO_ENDPOINT=minio:9000
MINIO_ACCESS_KEY=your_access_key
MINIO_SECRET_KEY=your_secret_key
MINIO_BUCKET=uploads

# 镜像仓库（可选）
REGISTRY=docker.io
EOF
        log_info "已创建 $ENV_FILE.example，请复制为 $ENV_FILE 并填入实际配置"
        exit 1
    fi
}

# 加载环境变量
load_env() {
    if [ -f "$ENV_FILE" ]; then
        export $(grep -v '^#' "$ENV_FILE" | xargs)
        log_info "环境变量已加载"
    fi
}

# 拉取最新镜像
pull_images() {
    log_info "正在拉取最新镜像..."

    docker compose -f "$COMPOSE_FILE" pull

    log_info "镜像拉取完成"
}

# 构建并推送镜像（可选）
build_and_push() {
    log_info "正在构建镜像..."

    # C端前端
    log_info "构建 C端前端..."
    docker build -t ${REGISTRY}/${PROJECT_PREFIX}/claude_blog_frontend:latest ../blog-frontend

    # C端 API
    log_info "构建 C端 API..."
    docker build -t ${REGISTRY}/${PROJECT_PREFIX}/claude_blog_api:latest ../blog-frontend-api

    # B端前端
    log_info "构建 B端前端..."
    docker build -t ${REGISTRY}/${PROJECT_PREFIX}/claude_blog_admin_frontend:latest ../claude-blog-admin/blog-admin-frontend

    # B端后端
    log_info "构建 B端后端..."
    docker build -t ${REGISTRY}/${PROJECT_PREFIX}/claude_blog_admin_backend:latest ../claude-blog-admin/blog-admin-backend

    log_info "推送镜像到仓库..."
    docker push ${REGISTRY}/${PROJECT_PREFIX}/claude_blog_frontend:latest
    docker push ${REGISTRY}/${PROJECT_PREFIX}/claude_blog_api:latest
    docker push ${REGISTRY}/${PROJECT_PREFIX}/claude_blog_admin_frontend:latest
    docker push ${REGISTRY}/${PROJECT_PREFIX}/claude_blog_admin_backend:latest

    log_info "镜像构建并推送完成"
}

# 停止服务
stop_services() {
    log_info "正在停止服务..."
    docker compose -f "$COMPOSE_FILE" down
    log_info "服务已停止"
}

# 启动服务
start_services() {
    log_info "正在启动服务..."
    docker compose -f "$COMPOSE_FILE" up -d
    log_info "服务启动完成"
}

# 查看日志
show_logs() {
    docker compose -f "$COMPOSE_FILE" logs -f --tail=100
}

# 查看状态
show_status() {
    docker compose -f "$COMPOSE_FILE" ps
}

# 清理未使用的镜像
cleanup() {
    log_info "正在清理未使用的镜像..."
    docker image prune -f
    log_info "清理完成"
}

# 显示帮助
show_help() {
    echo "Claude Blog 部署脚本"
    echo ""
    echo "用法: ./deploy.sh [命令]"
    echo ""
    echo "命令:"
    echo "  start       启动所有服务"
    echo "  stop        停止所有服务"
    echo "  restart     重启所有服务"
    echo "  pull        拉取最新镜像并启动"
    echo "  build       构建并推送镜像（需要配置镜像仓库）"
    echo "  logs        查看日志"
    echo "  status      查看服务状态"
    echo "  cleanup     清理未使用的镜像"
    echo "  help        显示帮助"
    echo ""
    echo "示例:"
    echo "  ./deploy.sh start     # 启动服务"
    echo "  ./deploy.sh pull      # 拉取新镜像并重启"
    echo "  ./deploy.sh logs      # 查看日志"
}

# 主逻辑
main() {
    check_commands

    case "${1:-help}" in
        start)
            check_env
            load_env
            start_services
            show_status
            ;;
        stop)
            stop_services
            ;;
        restart)
            check_env
            load_env
            stop_services
            start_services
            show_status
            ;;
        pull)
            check_env
            load_env
            pull_images
            stop_services
            start_services
            show_status
            ;;
        build)
            build_and_push
            ;;
        logs)
            show_logs
            ;;
        status)
            show_status
            ;;
        cleanup)
            cleanup
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            log_error "未知命令: $1"
            show_help
            exit 1
            ;;
    esac
}

main "$@"
