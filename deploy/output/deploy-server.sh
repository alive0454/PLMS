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
IMAGE_NAME="plms"

# 检查 Docker
if ! command -v docker &> /dev/null; then
    echo -e "${RED}错误: Docker 未安装${NC}"
    echo "请先安装 Docker: curl -fsSL https://get.docker.com | sh"
    exit 1
fi

echo -e "${GREEN}✓ Docker 已安装${NC}"

# 检查环境变量
if [ ! -f .env ]; then
    echo -e "${RED}错误: 未找到 .env 文件${NC}"
    exit 1
fi

# 检查是否为模板值
if grep -q "your-db-host\|your-password\|your-secret" .env; then
    echo -e "${YELLOW}警告: .env 包含模板值，请编辑修改:${NC}"
    echo "  vi .env"
    echo ""
    echo "需要修改的配置:"
    grep -E "DB_HOST|DB_PASSWORD|JWT_SECRET" .env | grep -E "your-|example|change"
    exit 1
fi

echo -e "${GREEN}✓ 环境变量已配置${NC}"

# 加载环境变量
set -a
source .env
set +a

echo ""
echo "[1/5] 生成 SSL 证书（如不存在）..."
if [ ! -f nginx/ssl/cert.pem ]; then
    echo "生成自签名证书..."
    mkdir -p nginx/ssl
    openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
        -keyout nginx/ssl/key.pem \
        -out nginx/ssl/cert.pem \
        -subj "/C=CN/ST=Beijing/O=PLMS/CN=localhost" 2>/dev/null || true
fi

echo ""
echo "[2/5] 停止旧容器..."
docker stop ${APP_NAME}-app 2>/dev/null || true
docker rm ${APP_NAME}-app 2>/dev/null || true
docker stop ${APP_NAME}-nginx 2>/dev/null || true
docker rm ${APP_NAME}-nginx 2>/dev/null || true

echo ""
echo "[3/5] 构建应用镜像..."
docker build -f Dockerfile -t ${IMAGE_NAME}:latest .

echo ""
echo "[4/5] 启动应用容器..."
# 绑定到 0.0.0.0:8081，让 Nginx 容器可以访问
docker run -d \
    --name ${APP_NAME}-app \
    --restart unless-stopped \
    -p 0.0.0.0:8081:8080 \
    --env-file .env \
    ${IMAGE_NAME}:latest

# 等待应用启动
sleep 3

# 检查应用是否健康
if ! curl -sf http://127.0.0.1:8081/health > /dev/null 2>&1; then
    echo -e "${YELLOW}应用启动中，等待 5 秒...${NC}"
    sleep 5
fi

# 检查应用是否真的启动
if ! curl -sf http://127.0.0.1:8081/health > /dev/null 2>&1; then
    echo -e "${RED}应用启动失败，查看日志:${NC}"
    docker logs ${APP_NAME}-app
    exit 1
fi

echo -e "${GREEN}✓ 应用启动成功${NC}"

echo ""
echo "[5/5] 启动 Nginx 容器..."

# 创建 Nginx 配置（使用 host.docker.internal 或 172.17.0.1 访问应用）
cat > /tmp/plms-nginx.conf << 'EOF'
server {
    listen 443 ssl;
    server_name _;

    ssl_certificate /etc/nginx/ssl/cert.pem;
    ssl_certificate_key /etc/nginx/ssl/key.pem;
    ssl_protocols TLSv1.2 TLSv1.3;

    # 日志
    access_log /var/log/nginx/plms-access.log;
    error_log /var/log/nginx/plms-error.log;

    # 反向代理到宿主机的 8081 端口
    location / {
        proxy_pass http://172.17.0.1:8081;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        proxy_http_version 1.1;
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    # 健康检查
    location /health {
        proxy_pass http://172.17.0.1:8081/health;
        access_log off;
    }
}
EOF

# 启动 Nginx
docker run -d \
    --name ${APP_NAME}-nginx \
    --restart unless-stopped \
    -p 8080:443 \
    -v "$(pwd)/nginx/nginx.conf:/etc/nginx/nginx.conf:ro" \
    -v "/tmp/plms-nginx.conf:/etc/nginx/conf.d/default.conf:ro" \
    -v "$(pwd)/nginx/ssl:/etc/nginx/ssl:ro" \
    nginx:alpine

sleep 2

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
echo ""
