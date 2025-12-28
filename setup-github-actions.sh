#!/bin/bash

# GitHub Actions 快速设置脚本

set -e

echo "🚀 GitHub Actions CI/CD 设置脚本"
echo ""

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 检查是否在项目根目录
if [ ! -d ".github/workflows" ]; then
    echo -e "${RED}❌ 请在项目根目录运行此脚本${NC}"
    exit 1
fi

echo -e "${BLUE}📋 步骤 1: 生成 SSH 密钥对${NC}"
echo ""

# 检查是否已存在密钥
if [ -f ~/.ssh/github_actions ] || [ -f ~/.ssh/github_actions.pub ]; then
    echo -e "${YELLOW}⚠️  检测到已存在的密钥文件${NC}"
    read -p "是否重新生成？(y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        rm -f ~/.ssh/github_actions ~/.ssh/github_actions.pub
    else
        echo -e "${GREEN}✅ 使用现有密钥${NC}"
        SKIP_KEYGEN=true
    fi
fi

if [ "$SKIP_KEYGEN" != "true" ]; then
    read -p "请输入服务器 IP 或域名（用于密钥注释）: " SERVER_HOST
    if [ -z "$SERVER_HOST" ]; then
        SERVER_HOST="github-actions"
    fi
    
    ssh-keygen -t ed25519 -C "github-actions@${SERVER_HOST}" -f ~/.ssh/github_actions -N ""
    echo -e "${GREEN}✅ SSH 密钥对已生成${NC}"
fi

echo ""
echo -e "${BLUE}📋 步骤 2: 显示公钥（需要添加到服务器）${NC}"
echo ""
echo -e "${YELLOW}════════════════════════════════════════════════════════${NC}"
cat ~/.ssh/github_actions.pub
echo -e "${YELLOW}════════════════════════════════════════════════════════${NC}"
echo ""
echo -e "${YELLOW}请将上面的公钥添加到服务器的 ~/.ssh/authorized_keys${NC}"
echo ""

read -p "是否现在添加到服务器？(y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    read -p "请输入服务器用户名（如 root）: " SERVER_USER
    read -p "请输入服务器 IP 或域名: " SERVER_HOST
    
    if [ -n "$SERVER_USER" ] && [ -n "$SERVER_HOST" ]; then
        echo "正在添加公钥到服务器..."
        ssh-copy-id -i ~/.ssh/github_actions.pub ${SERVER_USER}@${SERVER_HOST} || {
            echo -e "${YELLOW}⚠️  自动添加失败，请手动添加${NC}"
            echo "运行以下命令："
            echo "  cat ~/.ssh/github_actions.pub | ssh ${SERVER_USER}@${SERVER_HOST} 'mkdir -p ~/.ssh && cat >> ~/.ssh/authorized_keys'"
        }
    fi
fi

echo ""
echo -e "${BLUE}📋 步骤 3: 显示私钥（需要添加到 GitHub Secrets）${NC}"
echo ""
echo -e "${YELLOW}════════════════════════════════════════════════════════${NC}"
cat ~/.ssh/github_actions
echo -e "${YELLOW}════════════════════════════════════════════════════════${NC}"
echo ""
echo -e "${YELLOW}请将上面的私钥添加到 GitHub Secrets: SSH_PRIVATE_KEY${NC}"
echo ""

echo -e "${BLUE}📋 步骤 4: GitHub Secrets 配置清单${NC}"
echo ""
echo "请在 GitHub 仓库中添加以下 Secrets："
echo ""
echo -e "${GREEN}1. SSH_PRIVATE_KEY${NC}"
echo "   值: 上面的私钥内容（整个文件内容）"
echo ""
echo -e "${GREEN}2. SERVER_HOST${NC}"
read -p "   请输入服务器 IP 或域名: " SERVER_HOST_INPUT
if [ -n "$SERVER_HOST_INPUT" ]; then
    echo "   值: ${SERVER_HOST_INPUT}"
    SERVER_HOST_FINAL=$SERVER_HOST_INPUT
else
    SERVER_HOST_FINAL="your-server-ip"
fi
echo ""
echo -e "${GREEN}3. SERVER_USER${NC}"
read -p "   请输入服务器用户名（如 root）: " SERVER_USER_INPUT
if [ -n "$SERVER_USER_INPUT" ]; then
    echo "   值: ${SERVER_USER_INPUT}"
    SERVER_USER_FINAL=$SERVER_USER_INPUT
else
    SERVER_USER_FINAL="root"
fi
echo ""
echo -e "${GREEN}4. VITE_API_BASE_URL (可选)${NC}"
if [ -n "$SERVER_HOST_FINAL" ]; then
    VITE_API_BASE_URL="http://${SERVER_HOST_FINAL}:8081/api/v1"
    echo "   值: ${VITE_API_BASE_URL}"
else
    echo "   值: http://your-server-ip:8081/api/v1"
fi
echo ""

echo -e "${BLUE}📋 步骤 5: 服务器准备${NC}"
echo ""
echo "请在服务器上运行以下命令："
echo ""
echo -e "${GREEN}# 1. 安装 Docker 和 Docker Compose${NC}"
echo "   # Ubuntu/Debian:"
echo "   curl -fsSL https://get.docker.com -o get-docker.sh"
echo "   sh get-docker.sh"
echo "   sudo usermod -aG docker \$USER"
echo ""
echo -e "${GREEN}# 2. 创建部署目录${NC}"
echo "   sudo mkdir -p /opt/fluent-life/fluent-life-api"
echo "   sudo mkdir -p /opt/fluent-life/fluent-life-frontend"
echo "   sudo chown -R \$USER:\$USER /opt/fluent-life"
echo ""
echo -e "${GREEN}# 3. 创建 .env 文件${NC}"
echo "   cd /opt/fluent-life/fluent-life-api"
echo "   # 复制 env.example 并编辑"
echo "   # 设置 DB_PASSWORD, JWT_SECRET, VITE_API_BASE_URL 等"
echo ""

echo -e "${GREEN}✅ 设置脚本完成！${NC}"
echo ""
echo "下一步："
echo "1. 将公钥添加到服务器"
echo "2. 在 GitHub 仓库中添加 Secrets"
echo "3. 在服务器上准备环境"
echo "4. 推送到 main 分支测试自动部署"
echo ""

