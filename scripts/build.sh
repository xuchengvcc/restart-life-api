#!/bin/bash
# 构建脚本 - build.sh

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 项目配置
PROJECT_NAME="restart-life-api"
VERSION=${1:-latest}
DOCKER_IMAGE="$PROJECT_NAME:$VERSION"

echo -e "${BLUE}=== Restart Life API Docker Build Script ===${NC}"

# 检查是否使用腾讯云镜像或中国镜像
USE_TENCENT_MIRROR=false
USE_CHINA_MIRROR=false
if [ "$2" = "tencent" ] || [ "$1" = "tencent" ]; then
    USE_TENCENT_MIRROR=true
    echo -e "${YELLOW}Using Tencent Cloud mirror for faster build...${NC}"
elif [ "$2" = "china" ] || [ "$1" = "china" ]; then
    USE_CHINA_MIRROR=true
    echo -e "${YELLOW}Using China mirror for faster build...${NC}"
fi

echo -e "${YELLOW}Building Docker image: $DOCKER_IMAGE${NC}"

# 检查Docker是否运行
if ! docker info > /dev/null 2>&1; then
    echo -e "${RED}Error: Docker is not running. Please start Docker and try again.${NC}"
    exit 1
fi

# 选择Dockerfile
DOCKERFILE="docker/Dockerfile"
if [ "$USE_TENCENT_MIRROR" = true ]; then
    DOCKERFILE="docker/Dockerfile.tencent"
    echo -e "${YELLOW}Using optimized Dockerfile for Tencent Cloud...${NC}"
elif [ "$USE_CHINA_MIRROR" = true ]; then
    DOCKERFILE="docker/Dockerfile.china"
    echo -e "${YELLOW}Using optimized Dockerfile for China network...${NC}"
fi

# 构建Docker镜像
echo -e "${BLUE}Building Docker image...${NC}"
docker build -f "$DOCKERFILE" -t "$DOCKER_IMAGE" .

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Docker image built successfully: $DOCKER_IMAGE${NC}"
    
    # 显示镜像信息
    echo -e "${BLUE}Image details:${NC}"
    docker images "$PROJECT_NAME" --format "table {{.Repository}}\t{{.Tag}}\t{{.ID}}\t{{.CreatedAt}}\t{{.Size}}"
    
    echo -e "${GREEN}✅ Build completed successfully!${NC}"
    echo -e "${YELLOW}Next steps:${NC}"
    if [ "$USE_TENCENT_MIRROR" = true ]; then
        echo -e "  1. Run: ${BLUE}./scripts/start.sh tencent${NC} to start the application with Tencent mirrors"
    else
        echo -e "  1. Run: ${BLUE}./scripts/start.sh${NC} to start the application"
    fi
    echo -e "  2. Or run: ${BLUE}docker-compose up -d${NC} to start all services"
else
    echo -e "${RED}❌ Docker build failed${NC}"
    echo -e "${YELLOW}Tip: If you're in China, try: ${BLUE}./scripts/build.sh china${NC}"
    echo -e "${YELLOW}Or if using Tencent Cloud: ${BLUE}./scripts/build.sh tencent${NC}"
    exit 1
fi
