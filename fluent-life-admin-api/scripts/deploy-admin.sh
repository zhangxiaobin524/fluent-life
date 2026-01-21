#!/bin/bash

# 一鍵部署後台管理（admin-api + admin-frontend）
# 用法：
#   cd /opt/fluent-life/fluent-life-admin-api
#   chmod +x scripts/deploy-admin.sh
#   scripts/deploy-admin.sh
#
# 可選：指定前端打到哪個後端（給 Vite build 用）
#   export VITE_ADMIN_API_BASE_URL=http://你的域名或IP:8082

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
COMPOSE_FILE="${ROOT_DIR}/docker-compose.admin.yml"

if [ ! -f "$COMPOSE_FILE" ]; then
  echo -e "${RED}❌ 找不到 ${COMPOSE_FILE}${NC}"
  exit 1
fi

cd "$ROOT_DIR"

if command -v docker &>/dev/null && docker compose version &>/dev/null; then
  COMPOSE="docker compose"
elif command -v docker-compose &>/dev/null; then
  COMPOSE="docker-compose"
else
  echo -e "${RED}❌ 未檢測到 docker compose${NC}"
  exit 1
fi

echo -e "${YELLOW}📦 構建並啟動 admin-api(8082) + admin-frontend(5172)...${NC}"
$COMPOSE -f "$COMPOSE_FILE" up -d --build

echo -e "${YELLOW}⏳ 等待服務啟動...${NC}"
sleep 8

API_BASE="${VITE_ADMIN_API_BASE_URL:-http://localhost:8082}"
echo -e "${YELLOW}🏥 測試後端（可選）: ${API_BASE}/health${NC}"
if curl -fsS "${API_BASE}/health" >/dev/null 2>&1; then
  echo -e "${GREEN}✅ admin-api 健康檢查通過${NC}"
else
  echo -e "${YELLOW}⚠️  健康檢查未通過（若你後端沒做 /health 可忽略），建議看日誌${NC}"
fi

SERVER_IP="$(hostname -I 2>/dev/null | awk '{print $1}')"
echo -e "${GREEN}✅ 部署完成${NC}"
echo "後台前端: http://${SERVER_IP:-localhost}:5172"
echo "後台 API:  ${API_BASE}"
echo "查看日誌:  $COMPOSE -f $COMPOSE_FILE logs -f"

