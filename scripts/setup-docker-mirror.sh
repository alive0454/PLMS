#!/bin/bash

# Docker 国内镜像源配置脚本（阿里云服务器优化版）

echo "========================================="
echo "  Docker 国内镜像源配置"
echo "========================================="
echo ""

# 创建 Docker 配置目录
mkdir -p /etc/docker

# 备份旧配置
if [ -f /etc/docker/daemon.json ]; then
    cp /etc/docker/daemon.json /etc/docker/daemon.json.bak
    echo "已备份原配置到 /etc/docker/daemon.json.bak"
fi

# 写入国内镜像源配置（阿里云服务器优先使用阿里云镜像）
cat > /etc/docker/daemon.json << 'EOF'
{
  "registry-mirrors": [
    "https://mirror.baidubce.com",
    "https://hub-mirror.c.163.com",
    "https://docker.m.daocloud.io"
  ]
}
EOF

echo "已配置以下国内镜像源:"
echo "  - 百度镜像: https://mirror.baidubce.com"
echo "  - 网易镜像: https://hub-mirror.c.163.com"
echo "  - DaoCloud: https://docker.m.daocloud.io"
echo ""

# 重启 Docker
echo "重启 Docker 服务..."
if command -v systemctl &> /dev/null; then
    systemctl daemon-reload
    systemctl restart docker
else
    service docker restart
fi

echo ""
echo "========================================="
echo "配置完成!"
echo "========================================="
echo ""
echo "测试拉取: docker pull hello-world"
echo ""
