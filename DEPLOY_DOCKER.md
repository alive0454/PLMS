# PLMS Docker 部署指南（纯 Docker，无 compose）

## 快速开始

```bash
# 1. 上传代码到服务器
scp -r PLMS root@your-server:/opt/
ssh root@your-server
cd /opt/PLMS

# 2. 配置环境变量
cp .env.example .env
nano .env  # 编辑数据库配置

# 3. 执行部署
chmod +x deploy-docker.sh
./deploy-docker.sh
```

## 手动部署步骤

### 1. 构建镜像

```bash
docker build -t plms:latest .
```

### 2. 创建网络

```bash
docker network create plms-network
```

### 3. 启动应用容器

```bash
docker run -d \
    --name plms-app \
    --network plms-network \
    --restart unless-stopped \
    -p 8081:8080 \
    -e DB_HOST=your-db-host \
    -e DB_PORT=3306 \
    -e DB_USER=root \
    -e DB_PASSWORD=your-password \
    -e DB_NAME=plms \
    -e JWT_SECRET=your-secret-key \
    plms:latest
```

### 4. 启动 Nginx 容器

```bash
docker run -d \
    --name plms-nginx \
    --network plms-network \
    --restart unless-stopped \
    -p 8080:443 \
    -v $(pwd)/nginx/nginx.conf:/etc/nginx/nginx.conf:ro \
    -v $(pwd)/nginx/conf.d/plms-docker.conf:/etc/nginx/conf.d/default.conf:ro \
    -v $(pwd)/nginx/ssl:/etc/nginx/ssl:ro \
    -v $(pwd)/nginx/logs:/var/log/nginx \
    nginx:alpine
```

## 常用命令

```bash
# 查看运行中的容器
docker ps

# 查看日志
docker logs -f plms-app
docker logs -f plms-nginx

# 重启容器
docker restart plms-app plms-nginx

# 停止容器
docker stop plms-app plms-nginx

# 删除容器
docker rm -f plms-app plms-nginx

# 进入容器
docker exec -it plms-app sh
docker exec -it plms-nginx sh

# 查看容器资源使用
docker stats
```

## 更新部署

```bash
# 拉取最新代码
git pull

# 重新构建镜像
docker build -t plms:latest .

# 重启容器
docker restart plms-app

# 或者先删除再启动
docker rm -f plms-app
docker run -d ...
```

## 端口说明

| 服务 | 容器内端口 | 映射到宿主机 | 说明 |
|------|-----------|-------------|------|
| plms-app | 8080 | 8081 | Go 应用（内部） |
| plms-nginx | 443 | 8080 | HTTPS 对外端口 |

访问地址：`https://服务器IP:8080`

## 防火墙

```bash
# 开放 8080 端口
sudo ufw allow 8080/tcp
# 或
sudo firewall-cmd --permanent --add-port=8080/tcp
sudo firewall-cmd --reload
```

## 文件结构

```
PLMS/
├── Dockerfile              # 构建 Go 应用镜像
├── deploy-docker.sh        # 一键部署脚本
├── nginx/
│   ├── nginx.conf          # Nginx 主配置
│   ├── conf.d/
│   │   └── plms-docker.conf # Docker 部署专用站点配置
│   └── ssl/                # SSL 证书
└── .env                    # 环境变量
```
