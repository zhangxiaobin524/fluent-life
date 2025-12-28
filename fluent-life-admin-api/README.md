# 流畅人生管理后台 API

## 项目说明

这是流畅人生应用的管理后台后端API服务，提供用户、帖子、房间和训练记录的管理功能。

## 功能特性

- 用户管理：查看、搜索、删除用户
- 帖子管理：查看、搜索、删除帖子
- 房间管理：查看、搜索、删除、开启/关闭房间
- 训练统计：查看训练记录统计和详细记录

## 技术栈

- Go 1.24+
- Gin Web框架
- GORM (PostgreSQL)
- JWT认证（简化版）

## 配置

编辑 `configs/config.yaml` 文件配置数据库连接：

```yaml
DB_HOST: localhost
DB_PORT: "5432"
DB_USER: zhangxiaobin
DB_PASSWORD: ""
DB_NAME: fluent_life
```

## 启动服务

```bash
go run cmd/server/main.go
```

服务默认运行在 `http://localhost:8082`

## API接口

### 管理员登录
- POST `/api/v1/admin/login`
- Body: `{ "username": "admin", "password": "admin123" }`

### 用户管理
- GET `/api/v1/admin/users` - 获取用户列表
- GET `/api/v1/admin/users/:id` - 获取用户详情
- DELETE `/api/v1/admin/users/:id` - 删除用户

### 帖子管理
- GET `/api/v1/admin/posts` - 获取帖子列表
- GET `/api/v1/admin/posts/:id` - 获取帖子详情
- DELETE `/api/v1/admin/posts/:id` - 删除帖子

### 房间管理
- GET `/api/v1/admin/rooms` - 获取房间列表
- GET `/api/v1/admin/rooms/:id` - 获取房间详情
- DELETE `/api/v1/admin/rooms/:id` - 删除房间
- PATCH `/api/v1/admin/rooms/:id/toggle` - 开启/关闭房间

### 训练统计
- GET `/api/v1/admin/training/stats` - 获取训练统计
- GET `/api/v1/admin/training/records` - 获取训练记录列表

## 默认管理员账号

- 用户名: `admin`
- 密码: `admin123`

