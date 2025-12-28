# GitHub Actions CI/CD 配置指南

## 📋 概述

本项目已配置 GitHub Actions 用于自动化构建、测试和部署。

## 🔧 工作流程

### 1. CI - 持续集成 (`ci.yml`)

**触发条件：**
- 推送到 `main` 或 `develop` 分支
- 创建 Pull Request 到 `main` 或 `develop` 分支

**执行步骤：**
- ✅ 后端 Go 代码构建和测试
- ✅ 前端 React 代码构建
- ✅ Docker 镜像构建测试

### 2. CD - 自动部署 (`deploy.yml`)

**触发条件：**
- 推送到 `main` 分支
- 创建版本标签（如 `v1.0.0`）

**执行步骤：**
- 📦 同步代码到服务器
- 🐳 构建 Docker 镜像
- 🚀 部署服务
- 🏥 健康检查

### 3. 手动部署 (`deploy-manual.yml`)

**触发条件：**
- 手动触发（GitHub Actions UI）

**功能：**
- 可选择部署环境
- 可选择跳过测试

### 4. Docker 镜像构建 (`docker-build.yml`)

**触发条件：**
- 推送到 `main` 分支
- 创建版本标签
- 手动触发

**功能：**
- 构建并推送 Docker 镜像到 GitHub Container Registry

## 🔐 配置 GitHub Secrets

在 GitHub 仓库设置中添加以下 Secrets：

### 必需的 Secrets

1. **SSH_PRIVATE_KEY**
   - 用于 SSH 连接到服务器的私钥
   - 生成方式：
     ```bash
     ssh-keygen -t ed25519 -C "github-actions" -f ~/.ssh/github_actions
     ```
   - 将 `~/.ssh/github_actions` 的内容复制到 GitHub Secrets

2. **SERVER_HOST**
   - 服务器 IP 地址或域名
   - 例如：`120.55.250.184`

3. **SERVER_USER**
   - SSH 登录用户名
   - 例如：`root` 或 `ubuntu`

### 可选的 Secrets

4. **VITE_API_BASE_URL**
   - 前端 API 基础 URL（用于 Docker 构建）
   - 例如：`http://120.55.250.184:8081/api/v1`

## 📝 配置步骤

### 步骤 1：生成 SSH 密钥对

在本地机器上：

```bash
# 生成 SSH 密钥对
ssh-keygen -t ed25519 -C "github-actions" -f ~/.ssh/github_actions

# 查看公钥（需要添加到服务器）
cat ~/.ssh/github_actions.pub

# 查看私钥（需要添加到 GitHub Secrets）
cat ~/.ssh/github_actions
```

### 步骤 2：配置服务器 SSH 访问

在服务器上：

```bash
# 将公钥添加到 authorized_keys
mkdir -p ~/.ssh
chmod 700 ~/.ssh
echo "你的公钥内容" >> ~/.ssh/authorized_keys
chmod 600 ~/.ssh/authorized_keys

# 测试 SSH 连接（从本地）
ssh -i ~/.ssh/github_actions user@your-server-ip
```

### 步骤 3：在 GitHub 仓库中配置 Secrets

1. 进入 GitHub 仓库
2. 点击 **Settings** → **Secrets and variables** → **Actions**
3. 点击 **New repository secret**
4. 添加以下 Secrets：

| Secret 名称 | 值 | 说明 |
|------------|-----|------|
| `SSH_PRIVATE_KEY` | `~/.ssh/github_actions` 文件内容 | SSH 私钥 |
| `SERVER_HOST` | `120.55.250.184` | 服务器 IP |
| `SERVER_USER` | `root` | SSH 用户名 |
| `VITE_API_BASE_URL` | `http://120.55.250.184:8081/api/v1` | 前端 API URL（可选） |

### 步骤 4：准备服务器环境

确保服务器上已安装：

```bash
# Docker
docker --version

# Docker Compose
docker compose version

# 或
docker-compose --version
```

### 步骤 5：创建部署目录

在服务器上：

```bash
# 创建部署目录
sudo mkdir -p /opt/fluent-life/fluent-life-api
sudo mkdir -p /opt/fluent-life/fluent-life-frontend

# 设置权限
sudo chown -R $USER:$USER /opt/fluent-life
```

### 步骤 6：创建 .env 文件

在服务器上创建 `/opt/fluent-life/fluent-life-api/.env`：

```bash
cd /opt/fluent-life/fluent-life-api
nano .env
```

内容示例：

```env
# 数据库配置
DB_PASSWORD=your_password_here

# JWT 配置
JWT_SECRET=your_jwt_secret_here

# 前端 API URL
VITE_API_BASE_URL=http://120.55.250.184:8081/api/v1
```

## 🚀 使用方法

### 自动部署

1. **推送到 main 分支**：自动触发 CI 和 CD
   ```bash
   git push origin main
   ```

2. **创建版本标签**：自动部署
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

### 手动部署

1. 进入 GitHub 仓库
2. 点击 **Actions** 标签
3. 选择 **CD - Manual Deploy**
4. 点击 **Run workflow**
5. 选择环境和其他选项
6. 点击 **Run workflow**

### 查看部署状态

1. 进入 GitHub 仓库
2. 点击 **Actions** 标签
3. 查看工作流程运行状态

## 🔍 故障排查

### 问题 1：SSH 连接失败

**错误信息：**
```
Permission denied (publickey)
```

**解决方案：**
1. 检查 `SSH_PRIVATE_KEY` Secret 是否正确
2. 检查服务器上的 `~/.ssh/authorized_keys` 是否包含公钥
3. 检查服务器 SSH 配置是否允许密钥认证

### 问题 2：部署失败

**错误信息：**
```
Docker Compose 未安装
```

**解决方案：**
1. 在服务器上安装 Docker Compose
2. 检查 `docker compose version` 或 `docker-compose --version`

### 问题 3：.env 文件不存在

**错误信息：**
```
.env 文件不存在
```

**解决方案：**
1. 在服务器上创建 `/opt/fluent-life/fluent-life-api/.env`
2. 参考 `env.example` 文件

### 问题 4：服务健康检查失败

**解决方案：**
1. 检查服务器日志：`docker compose logs`
2. 检查端口是否被占用
3. 检查防火墙配置

## 📊 工作流程状态徽章

在 README.md 中添加状态徽章：

```markdown
![CI](https://github.com/your-username/your-repo/workflows/CI%20-%20Build%20and%20Test/badge.svg)
![CD](https://github.com/your-username/your-repo/workflows/CD%20-%20Deploy%20to%20Server/badge.svg)
```

## 🔄 自定义配置

### 修改部署路径

编辑 `.github/workflows/deploy.yml`：

```yaml
- name: Copy files to server
  run: |
    rsync -avz ... \
      ./fluent-life-api/ ${{ secrets.SERVER_USER }}@${{ secrets.SERVER_HOST }}:/your/custom/path/
```

### 添加部署前/后脚本

在 `deploy.yml` 中添加步骤：

```yaml
- name: Run pre-deploy script
  run: |
    ssh ${{ secrets.SERVER_USER }}@${{ secrets.SERVER_HOST }} << 'EOF'
      # 你的脚本
    EOF
```

### 多环境部署

创建不同的 workflow 文件：
- `.github/workflows/deploy-staging.yml`
- `.github/workflows/deploy-production.yml`

## 📚 相关文档

- [GitHub Actions 文档](https://docs.github.com/en/actions)
- [Docker Compose 文档](https://docs.docker.com/compose/)
- [SSH 密钥管理](https://docs.github.com/en/authentication/connecting-to-github-with-ssh)

## ✅ 检查清单

- [ ] SSH 密钥对已生成
- [ ] 服务器 SSH 访问已配置
- [ ] GitHub Secrets 已配置
- [ ] 服务器 Docker 和 Docker Compose 已安装
- [ ] 部署目录已创建
- [ ] `.env` 文件已创建
- [ ] 首次部署测试成功

## 🎉 完成！

配置完成后，每次推送到 `main` 分支都会自动触发部署。你可以在 GitHub Actions 页面查看部署状态和日志。

