#!/bin/bash

# PLMS 传统部署脚本（非 Docker）

set -e

echo "========================================="
echo "  PLMS 传统部署脚本"
echo "========================================="

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# 配置
APP_NAME="plms"
APP_DIR="/opt/plms"
APP_USER="plms"
GO_VERSION="1.24.0"

# 检查是否以 root 运行
if [ "$EUID" -ne 0 ]; then
    echo -e "${RED}请使用 root 权限运行: sudo ./deploy-traditional.sh${NC}"
    exit 1
fi

echo "步骤 1/6: 安装依赖..."
if command -v apt-get &> /dev/null; then
    apt-get update
    apt-get install -y git nginx curl wget
elif command -v yum &> /dev/null; then
    yum install -y git nginx curl wget
else
    echo -e "${YELLOW}请手动安装 git、nginx、curl${NC}"
fi

echo "步骤 2/6: 安装 Go..."
if ! command -v go &> /dev/null; then
    wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz -O /tmp/go.tar.gz
    tar -C /usr/local -xzf /tmp/go.tar.gz
    echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
    export PATH=$PATH:/usr/local/go/bin
    rm /tmp/go.tar.gz
fi
export PATH=$PATH:/usr/local/go/bin
go version

echo "步骤 3/6: 创建应用目录和用户..."
mkdir -p ${APP_DIR}
useradd -r -s /bin/false ${APP_USER} 2>/dev/null || true

# 复制代码到应用目录
cp -r . ${APP_DIR}/
chown -R ${APP_USER}:${APP_USER} ${APP_DIR}

echo "步骤 4/6: 编译应用..."
cd ${APP_DIR}
export GOPROXY=https://goproxy.cn,direct
export GO111MODULE=on
/usr/local/go/bin/go mod download
/usr/local/go/bin/go build -ldflags="-w -s" -o ${APP_NAME} ./cmd/server/main.go
chown ${APP_USER}:${APP_USER} ${APP_NAME}

echo "步骤 5/6: 配置环境变量..."
if [ ! -f ${APP_DIR}/.env ]; then
    echo -e "${YELLOW}请配置环境变量文件: ${APP_DIR}/.env${NC}"
    cat > ${APP_DIR}/.env.example << 'EOF'
APP_NAME=PLMS
APP_VERSION=1.0.0
APP_PORT=8080
APP_ENV=production

DB_HOST=your-db-host
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your-password
DB_NAME=plms

JWT_SECRET=your-secret-key
EOF
    cp ${APP_DIR}/.env.example ${APP_DIR}/.env
fi

echo "步骤 6/6: 创建 systemd 服务..."
cat > /etc/systemd/system/${APP_NAME}.service << EOF
[Unit]
Description=PLMS Application
After=network.target

[Service]
Type=simple
User=${APP_USER}
Group=${APP_USER}
WorkingDirectory=${APP_DIR}
EnvironmentFile=${APP_DIR}/.env
ExecStart=${APP_DIR}/${APP_NAME}
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

# 配置 Nginx
echo "配置 Nginx..."
cat > /etc/nginx/sites-available/plms << 'EOF'
server {
    listen 8080 ssl http2;
    server_name _;

    ssl_certificate /opt/plms/nginx/ssl/cert.pem;
    ssl_certificate_key /opt/plms/nginx/ssl/key.pem;

    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256;
    ssl_prefer_server_ciphers off;

    location / {
        proxy_pass http://127.0.0.1:8081;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
EOF

# 检查 SSL 证书
if [ ! -f ${APP_DIR}/nginx/ssl/cert.pem ]; then
    echo -e "${YELLOW}生成自签名 SSL 证书...${NC}"
    mkdir -p ${APP_DIR}/nginx/ssl
    openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
        -keyout ${APP_DIR}/nginx/ssl/key.pem \
        -out ${APP_DIR}/nginx/ssl/cert.pem \
        -subj "/C=CN/ST=Beijing/L=Beijing/O=PLMS/CN=localhost"
    chown -R ${APP_USER}:${APP_USER} ${APP_DIR}/nginx/ssl
fi

# 启用 Nginx 配置
if [ -d /etc/nginx/sites-enabled ]; then
    ln -sf /etc/nginx/sites-available/plms /etc/nginx/sites-enabled/plms
    rm -f /etc/nginx/sites-enabled/default
else
    # CentOS/RHEL
    cp /etc/nginx/sites-available/plms /etc/nginx/conf.d/plms.conf
fi

# 测试 Nginx 配置
nginx -t

# 启动服务
systemctl daemon-reload
systemctl enable plms
systemctl restart plms
systemctl restart nginx

# 等待服务启动
sleep 3

echo ""
echo "========================================="
echo -e "${GREEN}部署完成!${NC}"
echo "========================================="
echo ""
echo "应用目录: ${APP_DIR}"
echo "访问地址: https://$(curl -s ifconfig.me 2>/dev/null || echo 'your-server-ip'):8080"
echo ""
echo "常用命令:"
echo "  查看状态: systemctl status plms"
echo "  重启应用: systemctl restart plms"
echo "  查看日志: journalctl -u plms -f"
echo "  重启 Nginx: systemctl restart nginx"
echo ""
echo -e "${YELLOW}注意: 请修改 ${APP_DIR}/.env 中的数据库配置${NC}"
