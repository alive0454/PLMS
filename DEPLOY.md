# PLMS Docker 部署指南

本文档介绍如何在云服务器上使用 Docker 部署 PLMS 项目，并启用 HTTPS 访问。

## 目录

- [准备工作](#准备工作)
- [快速开始](#快速开始)
- [详细部署步骤](#详细部署步骤)
- [SSL 证书配置](#ssl-证书配置)
- [常用命令](#常用命令)
- [故障排查](#故障排查)

---

## 准备工作

### 1. 服务器要求

- **操作系统**: Ubuntu 20.04+ / CentOS 7+ / Debian 10+
- **内存**: 建议 2GB 以上
- **磁盘**: 建议 20GB 以上
- **网络**: 开放 80 (HTTP) 和 443 (HTTPS) 端口

### 2. 安装 Docker

```bash
# Ubuntu/Debian
curl -fsSL https://get.docker.com | sh
sudo usermod -aG docker $USER
newgrp docker

# CentOS/RHEL
curl -fsSL https://get.docker.com | sh
sudo systemctl start docker
sudo systemctl enable docker
sudo usermod -aG docker $USER
```

### 3. 安装 Docker Compose

```bash
# 安装 Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# 验证安装
docker-compose --version
```

### 4. 准备数据库

确保你有一个可访问的 MySQL/MariaDB 数据库：

```sql
-- 创建数据库
CREATE DATABASE IF NOT EXISTS plms CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 创建用户（可选，也可以使用 root）
CREATE USER 'plms'@'%' IDENTIFIED BY 'your-password';
GRANT ALL PRIVILEGES ON plms.* TO 'plms'@'%';
FLUSH PRIVILEGES;
```

---

## 快速开始

如果你已有服务器和域名，执行以下命令快速部署：

```bash
# 1. 克隆代码
git clone <your-repo-url>
cd PLMS

# 2. 配置环境变量
cp .env.example .env
# 编辑 .env 文件，填入数据库配置和 JWT 密钥

# 3. 运行部署脚本
chmod +x deploy.sh
./deploy.sh
```

---

## 详细部署步骤

### 步骤 1: 上传代码到服务器

```bash
# 方式 1: 使用 git
git clone <your-repo-url>
cd PLMS

# 方式 2: 使用 scp 上传本地代码
tar czvf plms.tar.gz --exclude='.git' --exclude='node_modules' .
scp plms.tar.gz root@your-server-ip:/opt/
ssh root@your-server-ip
cd /opt && tar xzvf plms.tar.gz && cd PLMS
```

### 步骤 2: 配置环境变量

```bash
cp .env.example .env
nano .env  # 或 vim .env
```

编辑以下内容：

```env
# 数据库配置（使用云数据库或服务器上的数据库）
DB_HOST=your-db-host      # 数据库地址
DB_PORT=3306              # 数据库端口
DB_USER=plms              # 数据库用户名
DB_PASSWORD=your-password # 数据库密码
DB_NAME=plms              # 数据库名

# JWT 密钥（生产环境务必设置复杂的随机字符串）
JWT_SECRET=your-super-secret-key-at-least-32-characters-long
```

### 步骤 3: 配置 SSL 证书

#### 方式 A: 使用 Let's Encrypt（推荐，需要域名）

```bash
# 运行部署脚本，选择选项 1
./deploy.sh
# 然后选择 1) 使用 Let's Encrypt 自动获取
```

或手动申请：

```bash
# 安装 certbot
sudo apt-get install certbot

# 申请证书（替换为你的域名）
sudo certbot certonly --standalone -d your-domain.com

# 复制证书到项目目录
sudo cp /etc/letsencrypt/live/your-domain.com/fullchain.pem nginx/ssl/cert.pem
sudo cp /etc/letsencrypt/live/your-domain.com/privkey.pem nginx/ssl/key.pem
sudo chmod 644 nginx/ssl/cert.pem
sudo chmod 600 nginx/ssl/key.pem
```

#### 方式 B: 使用自签名证书（测试用，无域名）

```bash
./scripts/generate-self-signed-cert.sh
```

#### 方式 C: 使用现有证书

将你的证书文件上传到服务器：

```bash
# 上传到 nginx/ssl/ 目录
scp your-cert.pem root@your-server-ip:/opt/PLMS/nginx/ssl/cert.pem
scp your-key.pem root@your-server-ip:/opt/PLMS/nginx/ssl/key.pem
```

### 步骤 4: 启动服务

```bash
# 构建并启动
docker-compose up -d --build

# 查看日志
docker-compose logs -f

# 查看服务状态
docker-compose ps
```

### 步骤 5: 验证部署

```bash
# 测试后端 API
curl http://localhost:8080/health

# 测试 HTTPS（替换为你的域名或 IP）
curl -k https://your-server-ip/health
```

---

## SSL 证书配置

### Let's Encrypt 自动续期

证书有效期为 90 天，需要设置自动续期：

```bash
# 编辑 crontab
sudo crontab -e

# 添加以下行（每周日凌晨 2 点检查续期）
0 2 * * 0 /opt/PLMS/scripts/renew-ssl.sh >> /var/log/letsencrypt-renewal.log 2>&1
```

### 手动续期

```bash
./scripts/renew-ssl.sh
```

---

## 常用命令

```bash
# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f                    # 所有服务
docker-compose logs -f plms-app           # 仅后端
docker-compose logs -f nginx              # 仅 Nginx

# 重启服务
docker-compose restart
docker-compose restart plms-app           # 仅重启后端

# 停止服务
docker-compose down

# 停止并删除数据卷（谨慎使用）
docker-compose down -v

# 进入容器
docker exec -it plms-app sh
docker exec -it plms-nginx sh

# 更新部署（拉取最新代码后）
git pull
docker-compose down
docker-compose up -d --build
```

---

## 故障排查

### 1. 无法访问 HTTPS

```bash
# 检查防火墙
sudo ufw status
sudo ufw allow 80
sudo ufw allow 443

# 检查端口监听
sudo netstat -tlnp | grep -E '80|443'

# 查看 Nginx 日志
docker-compose logs nginx
```

### 2. 后端服务无法连接数据库

```bash
# 检查网络连接
docker exec -it plms-app sh
wget -qO- http://your-db-host:3306

# 查看后端日志
docker-compose logs plms-app
```

### 3. 证书过期

```bash
# 手动续期
./scripts/renew-ssl.sh

# 强制更新证书
docker exec plms-nginx nginx -s reload
```

### 4. 内存不足

```bash
# 清理未使用的镜像和容器
docker system prune -a

# 查看内存使用
docker stats
```

---

## 生产环境建议

1. **数据库**: 使用云数据库服务（如阿里云 RDS、腾讯云 CDB）或单独部署数据库
2. **备份**: 定期备份数据库和 SSL 证书
3. **监控**: 配置日志监控和告警
4. **安全**: 
   - 修改默认的 JWT_SECRET
   - 修改 sys_user 表中的默认密码
   - 配置防火墙，仅开放必要端口

---

## 目录结构

```
PLMS/
├── cmd/                    # 应用入口
├── internal/               # 内部代码
├── nginx/                  # Nginx 配置
│   ├── conf.d/
│   │   └── plms.conf       # 站点配置
│   ├── ssl/                # SSL 证书目录
│   ├── logs/               # Nginx 日志
│   └── nginx.conf          # Nginx 主配置
├── scripts/                # 部署脚本
│   ├── generate-self-signed-cert.sh
│   ├── init-ssl.sh
│   └── renew-ssl.sh
├── Dockerfile              # Docker 构建文件
├── docker-compose.yml      # Docker Compose 配置
├── deploy.sh               # 一键部署脚本
├── .env                    # 环境变量（需手动创建）
└── .env.example            # 环境变量模板
```
