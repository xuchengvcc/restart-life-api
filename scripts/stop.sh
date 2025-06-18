#!/bin/bash
# 停止脚本 - stop.sh

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== Restart Life API Docker Stop Script ===${NC}"

# 检查docker-compose是否存在
if ! command -v docker-compose &> /dev/null; then
    echo -e "${RED}Error: docker-compose is not installed or not in PATH${NC}"
    exit 1
fi

echo -e "${YELLOW}Stopping all services...${NC}"
docker-compose down

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ All services stopped successfully!${NC}"
    
    # 询问是否清理数据
    echo -e "${YELLOW}Do you want to remove volumes (this will delete all data)? (y/N): ${NC}"
    read -r choice
    
    if [[ "$choice" =~ ^[Yy]$ ]]; then
        echo -e "${YELLOW}Removing volumes...${NC}"
        docker-compose down -v
        echo -e "${GREEN}✅ Volumes removed successfully!${NC}"
    else
        echo -e "${BLUE}Data volumes preserved.${NC}"
    fi
    
    echo -e "${GREEN}=== Other Useful Commands ===${NC}"
    echo -e "🧹 Clean up everything:    ${BLUE}docker-compose down -v --remove-orphans${NC}"
    echo -e "📊 View stopped containers: ${BLUE}docker-compose ps -a${NC}"
    echo -e "🔍 Remove unused images:   ${BLUE}docker image prune -f${NC}"
    
else
    echo -e "${RED}❌ Failed to stop services${NC}"
    exit 1
fi
