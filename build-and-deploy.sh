#!/bin/bash

# PLMS 本地编译 + Docker 部署脚本

set -e

echo "========================================="
echo "  PLMS 本地编译 + Docker 部署"
echo "========================================="

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

APP_NAME="plms"
IMAGE_NAME="plms"

# 检查环境变量
if [ ! -f .env ]; then
    echo -e "${YELLOW}创建 .env 文件...${NC}"
    cat > .env << 'EOF'
APP_NAME=PLMS
APP_VERSION=1.0.0
APP_PORT=8080
APP_ENV=production
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=
DB_NAME=plms
JWT_SECRET=change-this-secret-key
EOF
    echo -e "${RED}请先编辑 .env 文件配置数据库信息${NC}"
    exit 1
fi

source .env

echo -e "${GREEN}✓ 环境变量已加载${NC}"

# 步骤1: 本地编译 Go 应用
echo ""
echo "步骤 1/4: 本地编译 Go 应用..."

# 检查是否有 Go
if ! command -v go &> /dev/null; then
    echo -e "${RED}错误: 未找到 Go，请先安装 Go${NC}"
    echo "下载地址: https://go.dev/dl/"
    exit 1
fi

go version

# 下载依赖并编译
export GOPROXY=https://goproxy.cn,direct
export GO111MODULE=on
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64

go mod download
go build -ldflags="-w -s" -o server ./cmd/server/main.go

echo -e "${GREEN}✓ 编译完成${NC}"

# 步骤2: 生成 SSL 证书
echo ""
echo "步骤 2/4: 检查 SSL 证书..."
if [ ! -f nginx/ssl/cert.pem ]; then
    echo -e "${YELLOW}生成自签名 SSL 证书...${NC}"
    mkdir -p nginx/ssl
    openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
        -keyout nginx/ssl/key.pem \
        -out nginx/ssl/cert.pem \
        -subj "/C=CN/ST=Beijing/O=PLMS/CN=localhost" 2>/dev/null || true
fi

# 步骤3: 构建 Docker 镜像
echo ""
echo "步骤 3/4: 构建 Docker 镜像..."
docker build -f Dockerfile.local -t ${IMAGE_NAME}:latest .

# 清理编译文件
rm -f server

# 步骤4: 启动容器
echo ""
echo "步骤 4/4: 启动容器..."

# 停止旧容器
docker stop ${APP_NAME}-app 2>/dev/null || true
docker rm ${APP_NAME}-app 2>/dev/null || true
docker stop ${APP_NAME}-nginx 2>/dev/null || true
docker rm ${APP_NAME}-nginx 2>/dev/null || true

# 启动应用容器
docker run -d \
    --name ${APP_NAME}-app \
    --restart unless-stopped \
    -p 127.0.0.1:8081:8080 \
    --env-file .env \
    ${IMAGE_NAME}:latest

sleep 2

# 准备 Nginx 配置
cat nginx/conf.d/plms-docker.conf | sed 's/host.docker.internal/172.17.0.1/g' > /tmp/plms-nginx.conf

# 启动 Nginx 容器
docker run -d \
    --name ${APP_NAME}-nginx \
    --restart unless-stopped \
    -p 8080:443 \
    -v "$(pwd)/nginx/nginx.conf:/etc/nginx/nginx.conf:ro" \
    -v "/tmp/plms-nginx.conf:/etc/nginx/conf.d/default.conf:ro" \
    -v "$(pwd)/nginx/ssl:/etc/nginx/ssl:ro" \
    nginx:alpine

echo ""
echo "========================================="
echo -e "${GREEN}部署完成!${NC}"
echo "========================================="
echo ""
docker ps --filter "name=${APP_NAME}"
echo ""
SERVER_IP=$(curl -s ifconfig.me 2>/dev/null || hostname -I | awk '{print $1}')
echo "访问地址: https://${SERVER_IP}:8080"
echo ""
echo "常用命令:"
echo "  查看日志: docker logs -f ${APP_NAME}-app"
echo "  停止服务: docker stop ${APP_NAME}-app ${APP_NAME}-nginx"
