#!/bin/bash

# 步骤2: 服务器部署（在云服务器上执行）

set -e

echo "========================================="
echo "  步骤2: 服务器部署"
echo "========================================="

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

APP_NAME="plms"

# 检查 Docker
if ! command -v docker &> /dev/null; then
    echo -e "${RED}错误: Docker 未安装${NC}"
    echo "请先安装 Docker: curl -fsSL https://get.docker.com | sh"
    exit 1
fi

echo -e "${GREEN}✓ Docker 已安装${NC}"

# 检查部署包
if [ ! -f "plms-deploy.tar.gz" ]; then
    echo -e "${RED}错误: 未找到部署包 plms-deploy.tar.gz${NC}"
    echo "请先上传部署包到当前目录"
    exit 1
fi

# 创建部署目录
DEPLOY_DIR="/opt/plms"
mkdir -p ${DEPLOY_DIR}

echo ""
echo "[1/4] 解压部署包..."
tar xzvf plms-deploy.tar.gz -C ${DEPLOY_DIR}/
cd ${DEPLOY_DIR}

echo ""
echo "[2/4] 检查环境变量..."
if [ ! -f .env ]; then
    if [ -f .env.example ]; then
        cp .env.example .env
        echo -e "${YELLOW}已创建 .env 文件，请编辑配置:${NC}"
        echo "  vi ${DEPLOY_DIR}/.env"
        echo ""
        echo "配置项:"
        cat .env
        exit 1
    fi
fi

echo -e "${GREEN}✓ 环境变量已配置${NC}"

# 加载环境变量
set -a
source .env
set +a

echo ""
echo "[3/4] 生成 SSL 证书（如不存在）..."
if [ ! -f nginx/ssl/cert.pem ]; then
    echo "生成自签名证书..."
    mkdir -p nginx/ssl
    openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
        -keyout nginx/ssl/key.pem \
        -out nginx/ssl/cert.pem \
        -subj "/C=CN/ST=Beijing/O=PLMS/CN=localhost" 2>/dev/null || true
fi

echo ""
echo "[4/4] 构建并启动容器..."

# 停止旧容器
echo "停止旧容器..."
docker stop ${APP_NAME}-app 2>/dev/null || true
docker rm ${APP_NAME}-app 2>/dev/null || true
docker stop ${APP_NAME}-nginx 2>/dev/null || true
docker rm ${APP_NAME}-nginx 2>/dev/null || true

# 构建应用镜像
echo "构建应用镜像..."
docker build -f Dockerfile -t ${APP_NAME}:latest .

# 启动应用容器（仅监听 127.0.0.1，不对外暴露）
echo "启动应用容器..."
docker run -d \
    --name ${APP_NAME}-app \
    --restart unless-stopped \
    -p 127.0.0.1:8081:8080 \
    --env-file .env \
    ${APP_NAME}:latest

sleep 3

# 准备 Nginx 配置（指向宿主机 8081 端口）
cat nginx/conf.d/plms-docker.conf | sed 's/host.docker.internal/172.17.0.1/g' > /tmp/plms-nginx.conf

# 启动 Nginx 容器
echo "启动 Nginx 容器..."
docker run -d \
    --name ${APP_NAME}-nginx \
    --restart unless-stopped \
    -p 8080:443 \
    -v "${DEPLOY_DIR}/nginx/nginx.conf:/etc/nginx/nginx.conf:ro" \
    -v "/tmp/plms-nginx.conf:/etc/nginx/conf.d/default.conf:ro" \
    -v "${DEPLOY_DIR}/nginx/ssl:/etc/nginx/ssl:ro" \
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
echo "  查看应用日志: docker logs -f ${APP_NAME}-app"
echo "  查看Nginx日志: docker logs -f ${APP_NAME}-nginx"
echo "  重启应用: docker restart ${APP_NAME}-app"
echo "  重启Nginx: docker restart ${APP_NAME}-nginx"
echo "  停止所有: docker stop ${APP_NAME}-app ${APP_NAME}-nginx"
echo ""
