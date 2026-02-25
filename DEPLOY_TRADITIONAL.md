# PLMS 传统部署指南（非 Docker）

## 一键部署

```bash
# 1. 上传代码到服务器
scp -r PLMS root@your-server:/opt/

# 2. 在服务器上执行部署脚本
cd /opt/PLMS
chmod +x deploy-traditional.sh
sudo ./deploy-traditional.sh
```

## 手动部署步骤

### 1. 安装 Go

```bash
# 下载并安装 Go
wget https://go.dev/dl/go1.24.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.24.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
go version
```

### 2. 编译应用

```bash
cd /opt/PLMS
export GOPROXY=https://goproxy.cn,direct
go mod download
go build -o plms ./cmd/server/main.go
```

### 3. 配置环境变量

```bash
cp .env.example .env
nano .env
```

编辑数据库配置：
```env
APP_PORT=8081              # 内部端口，Nginx 会代理过来
DB_HOST=your-db-host
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your-password
DB_NAME=plms
JWT_SECRET=your-secret-key
```

### 4. 配置 Nginx

```bash
# Ubuntu/Debian
sudo nano /etc/nginx/sites-available/plms

# CentOS/RHEL
sudo nano /etc/nginx/conf.d/plms.conf
```

写入配置：
```nginx
server {
    listen 8080 ssl http2;
    server_name _;

    ssl_certificate /opt/PLMS/nginx/ssl/cert.pem;
    ssl_certificate_key /opt/PLMS/nginx/ssl/key.pem;

    ssl_protocols TLSv1.2 TLSv1.3;

    location / {
        proxy_pass http://127.0.0.1:8081;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### 5. 生成 SSL 证书

```bash
mkdir -p /opt/PLMS/nginx/ssl
cd /opt/PLMS/nginx/ssl

# 自签名证书（测试用）
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
    -keyout key.pem \
    -out cert.pem \
    -subj "/C=CN/ST=Beijing/O=PLMS/CN=localhost"

# 或 Let's Encrypt
# certbot certonly --standalone -d your-domain.com
```

### 6. 创建 systemd 服务

```bash
sudo nano /etc/systemd/system/plms.service
```

写入：
```ini
[Unit]
Description=PLMS Application
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/opt/PLMS
EnvironmentFile=/opt/PLMS/.env
ExecStart=/opt/PLMS/plms
Restart=always

[Install]
WantedBy=multi-user.target
```

### 7. 启动服务

```bash
sudo systemctl daemon-reload
sudo systemctl enable plms
sudo systemctl start plms
sudo systemctl restart nginx

# 检查状态
sudo systemctl status plms
sudo nginx -t
```

## 端口说明

| 服务 | 端口 | 说明 |
|------|------|------|
| Nginx | 8080 | HTTPS 对外端口 |
| Go 应用 | 8081 | 内部端口，Nginx 代理 |

访问地址：`https://your-server-ip:8080`

## 常用命令

```bash
# 查看应用日志
journalctl -u plms -f

# 重启应用
sudo systemctl restart plms

# 停止应用
sudo systemctl stop plms

# 重启 Nginx
sudo systemctl restart nginx

# 查看 Nginx 日志
sudo tail -f /var/log/nginx/error.log
```

## 防火墙

```bash
# Ubuntu
sudo ufw allow 8080/tcp

# CentOS
sudo firewall-cmd --permanent --add-port=8080/tcp
sudo firewall-cmd --reload
```
