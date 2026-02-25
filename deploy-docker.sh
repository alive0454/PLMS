#!/bin/bash

# PLMS Docker 部署脚本（不使用 docker-compose）

set -e

echo "========================================="
echo "  PLMS Docker 部署脚本"
echo "========================================="

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

APP_NAME="plms"
NETWORK_NAME="plms-network"

# 检查 Docker
if ! command -v docker &> /dev/null; then
    echo -e "${RED}错误: Docker 未安装${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Docker 已安装${NC}"

# 检查环境变量文件
if [ ! -f .env ]; then
    echo -e "${YELLOW}警告: 未找到 .env 文件，从模板创建...${NC}"
    cp .env.example .env
    echo -e "${RED}请编辑 .env 文件配置数据库和 JWT 密钥后再运行此脚本${NC}"
    exit 1
fi

# 加载环境变量
set -a
source .env
set +a

echo -e "${GREEN}✓ 环境变量已加载${NC}"

# 检查 SSL 证书
if [ ! -f nginx/ssl/cert.pem ] || [ ! -f nginx/ssl/key.pem ]; then
    echo -e "${YELLOW}警告: 未找到 SSL 证书，生成自签名证书...${NC}"
    mkdir -p nginx/ssl
    openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
        -keyout nginx/ssl/key.pem \
        -out nginx/ssl/cert.pem \
        -subj "/C=CN/ST=Beijing/L=Beijing/O=PLMS/CN=localhost"
    echo -e "${GREEN}✓ 自签名证书已生成${NC}"
fi

# 创建网络
echo "创建 Docker 网络..."
docker network create ${NETWORK_NAME} 2>/dev/null || true

# 停止并删除旧容器
echo "停止旧容器..."
docker stop ${APP_NAME}-app 2>/dev/null || true
docker rm ${APP_NAME}-app 2>/dev/null || true
docker stop ${APP_NAME}-nginx 2>/dev/null || true
docker rm ${APP_NAME}-nginx 2>/dev/null || true

# 构建应用镜像
echo "构建应用镜像..."
docker build -t ${APP_NAME}:latest .

# 启动 Go 应用容器
echo "启动 Go 应用容器..."
docker run -d \
    --name ${APP_NAME}-app \
    --network ${NETWORK_NAME} \
    --restart unless-stopped \
    -p 8081:8080 \
    -e APP_NAME="${APP_NAME}" \
    -e APP_VERSION="1.0.0" \
    -e APP_PORT="8080" \
    -e APP_ENV="production" \
    -e DB_HOST="${DB_HOST}" \
    -e DB_PORT="${DB_PORT}" \
    -e DB_USER="${DB_USER}" \
    -e DB_PASSWORD="${DB_PASSWORD}" \
    -e DB_NAME="${DB_NAME}" \
    -e JWT_SECRET="${JWT_SECRET}" \
    --health-cmd="wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1" \
    --health-interval=30s \
    --health-timeout=10s \
    --health-retries=3 \
    --health-start-period=40s \
    ${APP_NAME}:latest

# 启动 Nginx 容器
echo "启动 Nginx 容器..."
docker run -d \
    --name ${APP_NAME}-nginx \
    --network ${NETWORK_NAME} \
    --restart unless-stopped \
    -p 8080:443 \
    -v $(pwd)/nginx/nginx.conf:/etc/nginx/nginx.conf:ro \
    -v $(pwd)/nginx/conf.d/plms-docker.conf:/etc/nginx/conf.d/default.conf:ro \
    -v $(pwd)/nginx/ssl:/etc/nginx/ssl:ro \
    -v $(pwd)/nginx/logs:/var/log/nginx \
    nginx:alpine

# 等待服务启动
echo "等待服务启动..."
sleep 5

# 检查状态
echo ""
echo "容器状态:"
docker ps --filter "name=${APP_NAME}"

echo ""
echo "========================================="
echo -e "${GREEN}部署完成!${NC}"
echo "========================================="
echo ""
echo "访问地址:"
echo "  HTTPS: https://$(curl -s ifconfig.me 2>/dev/null || echo 'your-server-ip'):8080"
echo ""
echo "查看日志:"
echo "  应用日志: docker logs -f ${APP_NAME}-app"
echo "  Nginx日志: docker logs -f ${APP_NAME}-nginx"
echo ""
echo "常用命令:"
echo "  停止: docker stop ${APP_NAME}-app ${APP_NAME}-nginx"
echo "  重启: docker restart ${APP_NAME}-app ${APP_NAME}-nginx"
echo "  删除: docker rm -f ${APP_NAME}-app ${APP_NAME}-nginx"
echo ""
