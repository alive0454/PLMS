#!/bin/bash

# 使用 Let's Encrypt 获取 SSL 证书

set -e

echo "========================================="
echo "  Let's Encrypt SSL 证书申请"
echo "========================================="
echo ""

# 检查域名
read -p "请输入您的域名 (例如: plms.example.com): " DOMAIN

if [ -z "$DOMAIN" ]; then
    echo "错误: 域名不能为空"
    exit 1
fi

echo "域名: $DOMAIN"
echo ""

# 检查 certbot
if ! command -v certbot &> /dev/null; then
    echo "安装 Certbot..."
    if command -v apt-get &> /dev/null; then
        # Debian/Ubuntu
        sudo apt-get update
        sudo apt-get install -y certbot
    elif command -v yum &> /dev/null; then
        # CentOS/RHEL
        sudo yum install -y certbot
    else
        echo "错误: 无法自动安装 Certbot，请手动安装"
        exit 1
    fi
fi

# 创建目录
CERT_DIR="$(dirname "$0")/../nginx/ssl"
mkdir -p "$CERT_DIR"

# 申请证书
echo "申请证书..."
sudo certbot certonly --standalone -d "$DOMAIN" --agree-tos --no-eff-email -m "admin@$DOMAIN"

# 复制证书到项目目录
if [ -f "/etc/letsencrypt/live/$DOMAIN/fullchain.pem" ]; then
    sudo cp "/etc/letsencrypt/live/$DOMAIN/fullchain.pem" "$CERT_DIR/cert.pem"
    sudo cp "/etc/letsencrypt/live/$DOMAIN/privkey.pem" "$CERT_DIR/key.pem"
    sudo chmod 644 "$CERT_DIR/cert.pem"
    sudo chmod 600 "$CERT_DIR/key.pem"
    echo "证书已复制到 $CERT_DIR"
else
    echo "错误: 证书申请失败"
    exit 1
fi

# 创建自动续期脚本
RENEW_SCRIPT="$(dirname "$0")/renew-ssl.sh"
echo "#!/bin/bash" > "$RENEW_SCRIPT"
echo "# SSL 证书自动续期脚本" >> "$RENEW_SCRIPT"
echo "" >> "$RENEW_SCRIPT"
echo "cd \"$(dirname "$0")/..\"" >> "$RENEW_SCRIPT"
echo "" >> "$RENEW_SCRIPT"
echo "# 检测 Docker Compose 命令" >> "$RENEW_SCRIPT"
echo "if command -v docker-compose &> /dev/null; then" >> "$RENEW_SCRIPT"
echo "    DOCKER_COMPOSE=\"docker-compose\"" >> "$RENEW_SCRIPT"
echo "elif docker compose version &> /dev/null 2>&1; then" >> "$RENEW_SCRIPT"
echo "    DOCKER_COMPOSE=\"docker compose\"" >> "$RENEW_SCRIPT"
echo "else" >> "$RENEW_SCRIPT"
echo "    echo \"错误: Docker Compose 未安装\"" >> "$RENEW_SCRIPT"
echo "    exit 1" >> "$RENEW_SCRIPT"
echo "fi" >> "$RENEW_SCRIPT"
echo "" >> "$RENEW_SCRIPT"
echo "echo \"检查并续期 SSL 证书...\"" >> "$RENEW_SCRIPT"
echo "certbot renew --quiet" >> "$RENEW_SCRIPT"
echo "" >> "$RENEW_SCRIPT"
echo "# 复制新证书" >> "$RENEW_SCRIPT"
echo "if [ -f \"/etc/letsencrypt/live/$DOMAIN/fullchain.pem\" ]; then" >> "$RENEW_SCRIPT"
echo "    cp \"/etc/letsencrypt/live/$DOMAIN/fullchain.pem\" \"$CERT_DIR/cert.pem\"" >> "$RENEW_SCRIPT"
echo "    cp \"/etc/letsencrypt/live/$DOMAIN/privkey.pem\" \"$CERT_DIR/key.pem\"" >> "$RENEW_SCRIPT"
echo "    " >> "$RENEW_SCRIPT"
echo "    # 重启 Nginx" >> "$RENEW_SCRIPT"
echo "    \$DOCKER_COMPOSE restart nginx" >> "$RENEW_SCRIPT"
echo "    echo \"证书已更新并重启 Nginx\"" >> "$RENEW_SCRIPT"
echo "fi" >> "$RENEW_SCRIPT"

chmod +x "$RENEW_SCRIPT"

echo ""
echo "========================================="
echo "证书申请完成!"
echo "========================================="
echo ""
echo "续期:"
echo "  Let's Encrypt 证书有效期为 90 天"
echo "  续期脚本: $RENEW_SCRIPT"
echo ""
echo "添加定时任务自动续期:"
echo "  sudo crontab -e"
echo "  添加: 0 2 * * 0 $RENEW_SCRIPT"
echo ""
