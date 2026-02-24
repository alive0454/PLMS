#!/bin/bash

# PLMS Docker 部署脚本

set -e

echo "========================================="
echo "  PLMS Docker 部署脚本"
echo "========================================="

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 检查 Docker 和 Docker Compose
if ! command -v docker &> /dev/null; then
    echo -e "${RED}错误: Docker 未安装${NC}"
    echo "请先安装 Docker: https://docs.docker.com/engine/install/"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo -e "${RED}错误: Docker Compose 未安装${NC}"
    echo "请先安装 Docker Compose"
    exit 1
fi

echo -e "${GREEN}✓ Docker 和 Docker Compose 已安装${NC}"

# 检查环境变量文件
if [ ! -f .env ]; then
    echo -e "${YELLOW}警告: 未找到 .env 文件，从模板创建...${NC}"
    cp .env.example .env
    echo -e "${RED}请编辑 .env 文件配置数据库和 JWT 密钥后再运行此脚本${NC}"
    exit 1
fi

echo -e "${GREEN}✓ 环境变量文件已存在${NC}"

# 检查 SSL 证书
if [ ! -f nginx/ssl/cert.pem ] || [ ! -f nginx/ssl/key.pem ]; then
    echo -e "${YELLOW}警告: 未找到 SSL 证书${NC}"
    echo "请选择证书获取方式:"
    echo "1) 使用 Let's Encrypt 自动获取（需要域名）"
    echo "2) 使用自签名证书（测试用）"
    echo "3) 手动上传证书"
    read -p "请选择 [1-3]: " choice

    case $choice in
        1)
            echo "使用 Let's Encrypt 获取证书..."
            ./scripts/init-ssl.sh
            ;;
        2)
            echo "生成自签名证书..."
            ./scripts/generate-self-signed-cert.sh
            ;;
        3)
            echo -e "${RED}请将证书文件放在 nginx/ssl/ 目录:${NC}"
            echo "  - nginx/ssl/cert.pem (证书)"
            echo "  - nginx/ssl/key.pem (私钥)"
            exit 1
            ;;
        *)
            echo -e "${RED}无效选择${NC}"
            exit 1
            ;;
    esac
fi

echo -e "${GREEN}✓ SSL 证书已就绪${NC}"

# 拉取最新代码（可选）
if [ -d .git ]; then
    read -p "是否拉取最新代码? [y/N]: " pull
    if [[ $pull =~ ^[Yy]$ ]]; then
        git pull
    fi
fi

# 构建并启动
echo ""
echo "开始构建和启动服务..."
docker-compose down
docker-compose build --no-cache
docker-compose up -d

# 等待服务启动
echo ""
echo "等待服务启动..."
sleep 5

# 检查服务状态
echo ""
echo "检查服务状态..."
docker-compose ps

# 健康检查
echo ""
echo "执行健康检查..."
if curl -sf http://localhost:8080/health > /dev/null; then
    echo -e "${GREEN}✓ 后端服务运行正常${NC}"
else
    echo -e "${RED}✗ 后端服务可能未正常启动，请检查日志${NC}"
    docker-compose logs plms-app
fi

echo ""
echo "========================================="
echo -e "${GREEN}部署完成!${NC}"
echo "========================================="
echo ""
echo "访问地址:"
echo "  HTTP:  http://$(curl -s ifconfig.me 2>/dev/null || echo 'your-server-ip')"
echo "  HTTPS: https://$(curl -s ifconfig.me 2>/dev/null || echo 'your-server-ip')"
echo ""
echo "查看日志:"
echo "  docker-compose logs -f plms-app"
echo "  docker-compose logs -f nginx"
echo ""
echo "常用命令:"
echo "  停止服务: docker-compose down"
echo "  重启服务: docker-compose restart"
echo "  更新部署: ./deploy.sh"
echo ""
