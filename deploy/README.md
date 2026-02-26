# PLMS Docker 分步部署指南

## 概述

分两步部署：
1. **本地编译打包**（开发机执行）
2. **服务器部署**（云服务器执行）

## 目录结构

```
deploy/
├── build-local.sh      # 本地编译打包脚本
├── deploy-server.sh    # 服务器部署脚本
├── Dockerfile          # Docker 镜像构建文件
├── nginx-docker.conf   # Nginx 配置模板
└── README.md           # 本文件
```

---

## 步骤1: 本地编译打包

在开发机（Mac/Windows）上执行：

```bash
# 进入项目目录
cd /Users/wangyao/GolandProjects/PLMS

# 执行打包脚本
chmod +x deploy/build-local.sh
./deploy/build-local.sh
```

输出：`deploy/plms-deploy.tar.gz`

### 打包内容

- `plms-server` - 编译好的 Linux 二进制文件
- `Dockerfile` - Docker 构建文件
- `deploy-server.sh` - 服务器部署脚本
- `nginx/` - Nginx 配置和 SSL 证书目录
- `.env.example` - 环境变量模板

---

## 步骤2: 上传到服务器

```bash
# 上传到云服务器（修改为你的服务器IP）
scp deploy/plms-deploy.tar.gz root@your-server-ip:/opt/

# 登录服务器
ssh root@your-server-ip
```

---

## 步骤3: 服务器部署

在云服务器上执行：

```bash
cd /opt

# 解压
tar xzvf plms-deploy.tar.gz -C plms/
cd plms

# 配置环境变量
vi .env
```

编辑 `.env` 文件，填入数据库信息：

```env
# 数据库配置（必填）
DB_HOST=your-db-host
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your-password
DB_NAME=plms

# JWT密钥（必填，至少32位随机字符）
JWT_SECRET=your-super-secret-key-here
```

然后执行部署：

```bash
chmod +x deploy-server.sh
./deploy-server.sh
```

---

## 手动分步（如果不使用脚本）

### 1. 本地编译

```bash
cd /Users/wangyao/GolandProjects/PLMS

# 下载依赖
export GOPROXY=https://goproxy.cn,direct
go mod download

# 编译 Linux 版本
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64
go build -ldflags="-w -s" -o plms-server ./cmd/server/main.go
```

### 2. 上传文件

```bash
# 创建目录
ssh root@server "mkdir -p /opt/plms/nginx/ssl"

# 上传文件
scp plms-server root@server:/opt/plms/
scp deploy/Dockerfile root@server:/opt/plms/
scp nginx/nginx.conf root@server:/opt/plms/nginx/
scp nginx/conf.d/plms-docker.conf root@server:/opt/plms/nginx/conf.d/
```

### 3. 服务器启动

```bash
ssh root@server
cd /opt/plms

# 创建 .env 文件
cat > .env << EOF
APP_NAME=PLMS
APP_PORT=8080
DB_HOST=your-db-host
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your-password
DB_NAME=plms
JWT_SECRET=your-secret
EOF

# 生成证书
mkdir -p nginx/ssl
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
    -keyout nginx/ssl/key.pem \
    -out nginx/ssl/cert.pem \
    -subj "/C=CN/O=PLMS/CN=localhost"

# 构建镜像
docker build -f Dockerfile -t plms:latest .

# 启动应用（监听 127.0.0.1:8081）
docker run -d \
    --name plms-app \
    --restart unless-stopped \
    -p 127.0.0.1:8081:8080 \
    --env-file .env \
    plms:latest

# 启动 Nginx（监听 0.0.0.0:8080）
docker run -d \
    --name plms-nginx \
    --restart unless-stopped \
    -p 8080:443 \
    -v /opt/plms/nginx/nginx.conf:/etc/nginx/nginx.conf:ro \
    -v /opt/plms/nginx/conf.d/plms-docker.conf:/etc/nginx/conf.d/default.conf:ro \
    -v /opt/plms/nginx/ssl:/etc/nginx/ssl:ro \
    nginx:alpine
```

---

## 端口说明

| 服务 | 容器 | 宿主机 | 说明 |
|------|------|--------|------|
| 应用 | 8080 | 127.0.0.1:8081 | 仅本机访问 |
| Nginx | 443 | 0.0.0.0:8080 | HTTPS 对外 |

访问地址：`https://服务器IP:8080`

---

## 常用命令

```bash
# 查看容器
docker ps

# 查看日志
docker logs -f plms-app
docker logs -f plms-nginx

# 重启
docker restart plms-app
docker restart plms-nginx

# 停止
docker stop plms-app plms-nginx

# 删除容器
docker rm -f plms-app plms-nginx

# 更新部署（新版本）
# 1. 本地重新打包上传
# 2. 解压到 /opt/plms
# 3. docker build -t plms:latest .
# 4. docker restart plms-app
```

---

## 防火墙配置

```bash
# Ubuntu
ufw allow 8080/tcp

# CentOS
firewall-cmd --permanent --add-port=8080/tcp
firewall-cmd --reload
```
