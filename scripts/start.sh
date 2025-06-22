#!/bin/bash
# 启动脚本 - start.sh

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

PROJECT_NAME="restart-life-api"

echo -e "${BLUE}=== Restart Life API Docker Start Script ===${NC}"

# 检查是否使用腾讯云镜像
USE_TENCENT_MIRROR=false
COMPOSE_FILE="docker-compose.yml"
if [ "$1" = "tencent" ]; then
    USE_TENCENT_MIRROR=true
    COMPOSE_FILE="docker-compose.tencent.yml"
    echo -e "${YELLOW}Using Tencent Cloud mirror for faster startup...${NC}"
fi

# 检查Docker是否运行
if ! docker info > /dev/null 2>&1; then
    echo -e "${RED}Error: Docker is not running. Please start Docker and try again.${NC}"
    exit 1
fi

# 检查docker-compose是否存在
if ! command -v docker-compose &> /dev/null; then
    echo -e "${RED}Error: docker-compose is not installed or not in PATH${NC}"
    exit 1
fi

# 函数：显示服务状态
show_status() {
    echo -e "${BLUE}=== Service Status ===${NC}"
    docker-compose -f "$COMPOSE_FILE" ps
}

# 函数：显示日志
show_logs() {
    echo -e "${BLUE}=== Recent Logs ===${NC}"
    docker-compose -f "$COMPOSE_FILE" logs --tail=20
}

# 函数：等待服务健康
wait_for_services() {
    echo -e "${YELLOW}Waiting for services to be healthy...${NC}"
    
    local max_attempts=60
    local attempt=0
    
    while [ $attempt -lt $max_attempts ]; do
        if docker-compose -f "$COMPOSE_FILE" ps | grep -q "Up (healthy)"; then
            local healthy_count=$(docker-compose -f "$COMPOSE_FILE" ps | grep -c "Up (healthy)" || echo "0")
            local total_services=$(docker-compose -f "$COMPOSE_FILE" ps | grep -c "Up" || echo "0")
            
            echo -e "${YELLOW}Healthy services: $healthy_count${NC}"
            
            if [ "$healthy_count" -ge 3 ]; then  # app, mysql, redis
                echo -e "${GREEN}✅ All core services are healthy!${NC}"
                return 0
            fi
        fi
        
        echo -n "."
        sleep 5
        ((attempt++))
    done
    
    echo -e "\n${YELLOW}⚠️  Services may still be starting up. Check status manually.${NC}"
    return 1
}

# 主启动流程
echo -e "${YELLOW}Starting all services...${NC}"
if [ "$USE_CN_MIRROR" = true ]; then
    echo -e "${YELLOW}Using compose file: $COMPOSE_FILE${NC}"
fi

# 创建必要的目录
mkdir -p logs

# 启动服务
docker-compose -f "$COMPOSE_FILE" up -d

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Services started successfully!${NC}"
    
    # 等待服务健康
    wait_for_services
    
    # 显示服务状态
    show_status
    
    echo -e "${GREEN}=== Service URLs ===${NC}"
    echo -e "🚀 API Server:          ${BLUE}http://localhost:8080${NC}"
    echo -e "🔍 Health Check:        ${BLUE}http://localhost:8080/health${NC}"
    echo -e "🗄️  Database Admin:      ${BLUE}http://localhost:8081${NC}"
    echo -e "🔴 Redis Commander:     ${BLUE}http://localhost:8082${NC}"
    
    echo -e "\n${GREEN}=== Useful Commands ===${NC}"
    if [ "$USE_TENCENT_MIRROR" = true ]; then
        echo -e "📊 View logs:           ${BLUE}docker-compose -f $COMPOSE_FILE logs -f${NC}"
        echo -e "📊 View app logs:       ${BLUE}docker-compose -f $COMPOSE_FILE logs -f app${NC}"
        echo -e "🔄 Restart services:    ${BLUE}docker-compose -f $COMPOSE_FILE restart${NC}"
        echo -e "🛑 Stop services:       ${BLUE}docker-compose -f $COMPOSE_FILE down${NC}"
        echo -e "🧹 Clean up:            ${BLUE}docker-compose -f $COMPOSE_FILE down -v${NC}"
    else
        echo -e "📊 View logs:           ${BLUE}docker-compose logs -f${NC}"
        echo -e "📊 View app logs:       ${BLUE}docker-compose logs -f app${NC}"
        echo -e "🔄 Restart services:    ${BLUE}docker-compose restart${NC}"
        echo -e "🛑 Stop services:       ${BLUE}docker-compose down${NC}"
        echo -e "🧹 Clean up:            ${BLUE}docker-compose down -v${NC}"
    fi
    
    # 显示最近的日志
    show_logs
    
else
    echo -e "${RED}❌ Failed to start services${NC}"
    echo -e "${YELLOW}Checking logs for errors:${NC}"
    docker-compose -f "$COMPOSE_FILE" logs
    echo -e "${YELLOW}Tip: If you're in China, try: ${BLUE}./scripts/start.sh tencent${NC}"
    exit 1
fi
