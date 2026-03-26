# Claude Blog 部署指南

## 环境要求

- CentOS 7/8 或其他 Linux 发行版
- Docker 20.10+
- Docker Compose 2.0+
- 2GB+ RAM
- 20GB+ 磁盘空间

## 目录结构

```
deploy/
├── docker-compose.yml    # 生产环境编排配置
├── deploy.sh            # 部署脚本
├── .env.example         # 环境变量示例
└── README.md            # 本文件
```

## 快速部署

### 1. 安装 Docker

```bash
# CentOS
sudo yum install -y yum-utils
sudo yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
sudo yum install docker-ce docker-ce-cli containerd.io
sudo systemctl start docker
sudo systemctl enable docker
```

### 2. 安装 Docker Compose

```bash
sudo curl -L "https://github.com/docker/compose/releases/download/v2.20.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

### 3. 配置环境变量

```bash
cd deploy
cp .env.example .env
nano .env  # 编辑配置
```

### 4. 启动服务

```bash
chmod +x deploy.sh
./deploy.sh start
```

## 部署命令

| 命令 | 说明 |
|-----|------|
| `./deploy.sh start` | 启动所有服务 |
| `./deploy.sh stop` | 停止所有服务 |
| `./deploy.sh restart` | 重启所有服务 |
| `./deploy.sh pull` | 拉取最新镜像并重启 |
| `./deploy.sh logs` | 查看日志 |
| `./deploy.sh status` | 查看服务状态 |
| `./deploy.sh cleanup` | 清理未使用的镜像 |

## 服务访问

| 服务 | 地址 |
|-----|------|
| C端前端 | http://your_domain:3000 |
| C端 API | http://your_domain:8080 |
| B端前端 | http://your_domain:3001 |
| B端后端 | http://your_domain:8082 |
| MinIO Console | http://your_domain:9001 |

## 数据持久化

所有数据通过 Docker Volume 持久化：

- `mysql-data` - MySQL 数据
- `redis-data` - Redis 数据
- `minio-data` - MinIO 文件存储

## 备份

### 备份 MySQL 数据

```bash
docker exec cc-blog-mysql mysqldump -u root -p blog > backup_$(date +%Y%m%d).sql
```

### 备份 MinIO 文件

```bash
docker exec cc-blog-minio mc alias set local http://localhost:9000 $MINIO_ACCESS_KEY $MINIO_SECRET_KEY
docker exec cc-blog-minio mc mirror local/uploads ./backups/
```

## 更新部署

### 方式一：使用预编译镜像

1. 在有源码的机器上构建镜像
2. 推送到镜像仓库
3. 在生产服务器上执行：

```bash
./deploy.sh pull
```

### 方式二：直接在服务器构建（不推荐）

```bash
# 在 deploy 目录执行
git clone https://github.com/jumoshen/claude_blog.git temp_source
cd temp_source
# 修改 Dockerfile 等
docker build -t jumoshen/claude_blog_frontend:latest ./blog-frontend
# ... 构建其他服务
cd ..
rm -rf temp_source
./deploy.sh restart
```

## 故障排查

### 查看日志

```bash
./deploy.sh logs
```

### 查看服务状态

```bash
./deploy.sh status
```

### 重启单个服务

```bash
docker compose -f docker-compose.yml restart blog-frontend
```

### 进入容器调试

```bash
docker exec -it cc-blog-api /bin/sh
```

## 安全建议

1. 修改 `.env` 中的所有密码
2. 配置防火墙，仅开放必要端口
3. 启用 HTTPS（通过 Nginx 反向代理）
4. 定期更新镜像版本
5. 配置自动备份

## 镜像构建（可选）

如果使用自己的镜像仓库：

```bash
# 设置镜像仓库
export REGISTRY=your.registry.com

# 构建并推送
./deploy.sh build
```
