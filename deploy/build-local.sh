#!/bin/bash

# 步骤1: 本地编译并打包（在开发机/Mac/Windows上执行）

set -e

echo "========================================="
echo "  步骤1: 本地编译打包"
echo "========================================="

cd "$(dirname "$0")/.."

# 创建输出目录
mkdir -p deploy/output

# 清理旧文件
rm -rf deploy/output/*

echo ""
echo "[1/4] 下载 Go 依赖..."
export GOPROXY=https://goproxy.cn,direct
export GO111MODULE=on
go mod download

echo ""
echo "[2/4] 编译 Linux 二进制文件..."
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64
go build -ldflags="-w -s" -o deploy/output/plms-server ./cmd/server/main.go

echo ""
echo "[3/4] 复制部署文件..."

# 复制必要文件
cp deploy/Dockerfile deploy/output/
cp deploy/deploy-server.sh deploy/output/
cp -r nginx deploy/output/

# 复制环境变量（如果存在）
if [ -f .env ]; then
    echo "复制现有 .env 配置..."
    cp .env deploy/output/
else
    echo "创建 .env 模板..."
    cat > deploy/output/.env << 'EOF'
# 应用配置
APP_NAME=PLMS
APP_VERSION=1.0.0
APP_PORT=8080
APP_ENV=production

# 数据库配置（必填）
DB_HOST=your-db-host
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your-password
DB_NAME=plms

# JWT密钥（必填，生产环境请设置复杂随机字符串）
JWT_SECRET=your-secret-key-min-32-characters-long
EOF
fi

echo ""
echo "[4/4] 创建部署包..."

# 创建最终部署包（使用 COPYFILE_DISABLE 避免 Mac 扩展属性）
cd deploy/output
COPYFILE_DISABLE=1 tar czvf ../plms-deploy.tar.gz .

echo ""
echo "========================================="
echo "  打包完成!"
echo "========================================="
echo ""
echo "部署包: deploy/plms-deploy.tar.gz"
echo ""
echo "下一步: 上传到服务器"
echo "  scp deploy/plms-deploy.tar.gz root@your-server:/opt/"
echo ""
