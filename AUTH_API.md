# 登录认证系统 API 文档

## 概述

本系统实现了基于 JWT 的登录认证功能，包含以下特性：

- JWT Token 认证
- 默认密码强制修改
- 密码 BCrypt 加密存储
- 跨域支持

## 接口列表

### 1. 用户登录

**接口地址**: `POST /api/v1/auth/login`

**无需认证**

**请求参数**:
```json
{
  "username": "admin",
  "password": "admin123"
}
```

**成功响应**:
```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "admin",
      "name": "系统管理员",
      "role": "admin",
      "isDefaultPassword": true,
      "avatar": "https://..."
    }
  }
}
```

**错误响应**:
```json
{
  "code": 401,
  "message": "用户名或密码错误",
  "data": null
}
```

### 2. 获取当前用户信息

**接口地址**: `GET /api/v1/auth/current-user`

**需要认证**: `Authorization: Bearer {token}`

**成功响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "username": "admin",
    "name": "系统管理员",
    "role": "admin",
    "isDefaultPassword": false,
    "avatar": "https://..."
  }
}
```

### 3. 修改密码

**接口地址**: `POST /api/v1/auth/change-password`

**需要认证**: `Authorization: Bearer {token}`

**请求参数**:
```json
{
  "oldPassword": "admin123",
  "newPassword": "newpassword123"
}
```

**成功响应**:
```json
{
  "code": 200,
  "message": "密码修改成功",
  "data": null
}
```

**错误响应**:
```json
{
  "code": 400,
  "message": "旧密码错误",
  "data": null
}
```

### 4. 登出

**接口地址**: `POST /api/v1/auth/logout`

**需要认证**: `Authorization: Bearer {token}`

**成功响应**:
```json
{
  "code": 200,
  "message": "登出成功",
  "data": null
}
```

### 5. 检查是否使用默认密码

**接口地址**: `GET /api/v1/auth/check-default-password`

**需要认证**: `Authorization: Bearer {token}`

**成功响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "isDefaultPassword": false
  }
}
```

## 权限控制说明

### 公开接口（无需认证）

- `POST /api/v1/auth/login` - 登录

### 需要认证但不受默认密码限制的接口

- `GET /api/v1/auth/current-user` - 获取当前用户信息
- `POST /api/v1/auth/change-password` - 修改密码
- `POST /api/v1/auth/logout` - 登出
- `GET /api/v1/auth/check-default-password` - 检查默认密码状态

### 需要认证且受默认密码限制的接口

所有其他接口在使用默认密码登录时会被拦截，返回：

```json
{
  "code": 403,
  "message": "请先修改默认密码",
  "data": {
    "requireChangePassword": true
  }
}
```

## 环境变量配置

```bash
# JWT 密钥（生产环境请修改）
JWT_SECRET=your-secret-key-change-in-production
```

## 默认用户

| 用户名 | 密码 | 角色 |
|--------|------|------|
| admin | admin123 | admin |
| user | user123 | user |

## 初始化数据库

运行 SQL 脚本初始化默认用户：

```bash
mysql -u root -p plms < scripts/init_users.sql
```

或者手动执行：

```sql
-- 密码: admin123
-- 哈希: $2a$10$Lhhzfkmp3h7r5uHrzgCorOHjhi6JFI.AffFhJVvNX5TlNiakk8Chq

-- 密码: user123
-- 哈希: $2a$10$q0ToiSZFi8ahfIzfoHBt8eZ7YNv5pdht4RrAVi14n3qUKRAgZFTI.

INSERT INTO sys_user (username, password, name, role, is_default_password, status) VALUES 
('admin', '$2a$10$Lhhzfkmp3h7r5uHrzgCorOHjhi6JFI.AffFhJVvNX5TlNiakk8Chq', '系统管理员', 'admin', 1, 1),
('user', '$2a$10$q0ToiSZFi8ahfIzfoHBt8eZ7YNv5pdht4RrAVi14n3qUKRAgZFTI.', '普通用户', 'user', 1, 1);
```

## 前端集成建议

### 请求拦截器

```javascript
// 自动添加 Authorization 请求头
axios.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});
```

### 响应拦截器

```javascript
// 处理 401 状态码
axios.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      window.location.href = '/user/login';
    }
    // 处理 403 默认密码强制修改
    if (error.response?.status === 403 && 
        error.response?.data?.data?.requireChangePassword) {
      window.location.href = '/user/change-password';
    }
    return Promise.reject(error);
  }
);
```

## 项目文件结构

```
internal/
├── middleware/
│   └── auth.go          # JWT 认证中间件、跨域中间件
├── handlers/
│   └── auth.go          # 认证相关接口处理器
├── services/
│   └── auth.go          # 认证服务（登录、密码修改、Token生成）
├── models/
│   └── sys_user.go      # 系统用户模型
└── config/
    └── config.go        # 配置（添加 GetEnv 函数）

cmd/server/
└── main.go              # 更新路由配置

scripts/
├── generate_password.go # 密码哈希生成工具
└── init_users.sql       # 初始化用户 SQL
```
