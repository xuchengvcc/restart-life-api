#!/bin/bash
# 网络问题解决脚本 - fix-network.sh

set -e

echo "=== Docker网络问题修复工具 ==="

# 检查Docker是否运行
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker未运行，请先启动Docker"
    exit 1
fi

echo "🔍 正在诊断网络问题..."

# 测试Docker Hub连接
echo "测试Docker Hub连接..."
if docker pull hello-world > /dev/null 2>&1; then
    echo "✅ Docker Hub连接正常"
    DOCKER_HUB_OK=true
else
    echo "❌ Docker Hub连接失败"
    DOCKER_HUB_OK=false
fi

# 测试阿里云镜像
echo "测试阿里云镜像连接..."
if docker pull registry.cn-hangzhou.aliyuncs.com/library/hello-world > /dev/null 2>&1; then
    echo "✅ 阿里云镜像连接正常"
    ALIYUN_OK=true
else
    echo "❌ 阿里云镜像连接失败"
    ALIYUN_OK=false
fi

# 测试腾讯云镜像
echo "测试腾讯云镜像连接..."
if docker pull ccr.ccs.tencentyun.com/library/hello-world > /dev/null 2>&1; then
    echo "✅ 腾讯云镜像连接正常"
    TENCENT_OK=true
else
    echo "❌ 腾讯云镜像连接失败"
    TENCENT_OK=false
fi

echo ""
echo "=== 诊断结果 ==="

if [ "$DOCKER_HUB_OK" = true ]; then
    echo "推荐使用：普通版本"
    echo "命令：./scripts/build.sh"
elif [ "$ALIYUN_OK" = true ]; then
    echo "推荐使用：阿里云版本"
    echo "命令：./scripts/build.sh aliyun"
elif [ "$TENCENT_OK" = true ]; then
    echo "推荐使用：腾讯云版本"
    echo "命令：./scripts/build.sh tencent"
else
    echo "⚠️  所有镜像源都无法连接"
    echo "建议："
    echo "1. 检查网络连接"
    echo "2. 配置Docker镜像加速器"
    echo "3. 使用VPN或代理"
    echo "4. 使用离线构建方式"
fi

echo ""
echo "=== 镜像加速器配置建议 ==="
echo "在Docker Desktop的设置中添加以下镜像加速器："
echo ""
echo "{"
echo '  "registry-mirrors": ['
echo '    "https://mirror.ccs.tencentyun.com",'
echo '    "https://docker.mirrors.ustc.edu.cn",'
echo '    "https://hub-mirror.c.163.com",'
echo '    "https://registry.cn-hangzhou.aliyuncs.com"'
echo '  ]'
echo "}"
