#!/bin/bash

# PLMS 二进制直接部署（不使用 Docker）

set -e

echo "========================================="
echo "  PLMS 二进制部署"
echo "========================================="

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

APP_NAME="plms"
DEPLOY_DIR="/opt/plms"

# 检查环境变量
if [ ! -f .env ]; then
    echo -e "${RED}错误: 未找到 .env 文件${NC}"
    exit 1
fi

# 检查是否为模板值
if grep -q "your-db-host\|your-password\|your-secret" .env; then
    echo -e "${YELLOW}警告: .env 包含模板值，请编辑修改${NC}"
    exit 1
fi

echo -e "${GREEN}✓ 环境变量已配置${NC}"

# 创建目录
mkdir -p ${DEPLOY_DIR}
cp plms-server ${DEPLOY_DIR}/
cp .env ${DEPLOY_DIR}/
cp -r nginx ${DEPLOY_DIR}/

# 生成证书
if [ ! -f nginx/ssl/cert.pem ]; then
    echo "生成自签名证书..."
    mkdir -p nginx/ssl
    openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
        -keyout nginx/ssl/key.pem \
        -out nginx/ssl/cert.pem \
        -subj "/C=CN/ST=Beijing/O=PLMS/CN=localhost" 2>/dev/null || true
fi

# 创建 systemd 服务
cat > /etc/systemd/system/${APP_NAME}.service << EOF
[Unit]
Description=PLMS Application
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=${DEPLOY_DIR}
EnvironmentFile=${DEPLOY_DIR}/.env
ExecStart=${DEPLOY_DIR}/plms-server
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

# 配置 Nginx
cat > /etc/nginx/conf.d/plms.conf << EOF
server {
    listen 8080 ssl http2;
    server_name _;

    ssl_certificate ${DEPLOY_DIR}/nginx/ssl/cert.pem;
    ssl_certificate_key ${DEPLOY_DIR}/nginx/ssl/key.pem;
    ssl_protocols TLSv1.2 TLSv1.3;

    location / {
        proxy_pass http://127.0.0.1:8081;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
    }
}
EOF

# 修改 .env 端口为 8081（内部端口）
sed -i 's/APP_PORT=8080/APP_PORT=8081/' ${DEPLOY_DIR}/.env

# 启动服务
systemctl daemon-reload
systemctl enable plms
systemctl restart plms
systemctl restart nginx

echo ""
echo "========================================="
echo -e "${GREEN}部署完成!${NC}"
echo "========================================="
echo ""
echo "访问: https://\$(curl -s ifconfig.me):8080"
echo ""
echo "命令:"
echo "  systemctl status plms"
echo "  systemctl restart plms"
echo "  journalctl -u plms -f"
