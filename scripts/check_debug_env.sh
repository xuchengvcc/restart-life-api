#!/bin/bash

echo "=== VS Code 调试环境验证 ==="

# 检查 VS Code 配置文件
echo "检查 VS Code 配置文件..."
if [ -f ".vscode/launch.json" ]; then
    echo "✅ launch.json - 调试配置文件存在"
else
    echo "❌ launch.json - 调试配置文件缺失"
fi

if [ -f ".vscode/tasks.json" ]; then
    echo "✅ tasks.json - 任务配置文件存在"
else
    echo "❌ tasks.json - 任务配置文件缺失"
fi

if [ -f ".vscode/settings.json" ]; then
    echo "✅ settings.json - VS Code 设置文件存在"
else
    echo "❌ settings.json - VS Code 设置文件缺失"
fi

if [ -f ".vscode/api-test.http" ]; then
    echo "✅ api-test.http - HTTP 测试文件存在"
else
    echo "❌ api-test.http - HTTP 测试文件缺失"
fi

echo ""

# 检查 Go 环境
echo "检查 Go 开发环境..."
if command -v go &> /dev/null; then
    echo "✅ Go 已安装: $(go version)"
else
    echo "❌ Go 未安装"
    exit 1
fi

# 检查项目构建
echo ""
echo "检查项目构建..."
if go build -o /tmp/test-build cmd/server/*.go; then
    echo "✅ 项目构建成功"
    rm -f /tmp/test-build
else
    echo "❌ 项目构建失败"
    exit 1
fi

# 检查测试
echo ""
echo "检查单元测试..."
if cd cmd/server && go test -v; then
    echo "✅ 单元测试通过"
    cd ../..
else
    echo "❌ 单元测试失败"
    cd ../..
    exit 1
fi

# 检查调试构建
echo ""
echo "检查调试构建..."
if go build -gcflags="all=-N -l" -o build/debug-test cmd/server/*.go; then
    echo "✅ 调试版本构建成功"
    rm -f build/debug-test
else
    echo "❌ 调试版本构建失败"
    exit 1
fi

echo ""
echo "🎉 VS Code 调试环境验证完成！"
echo ""
echo "使用方法："
echo "1. 在 VS Code 中打开项目"
echo "2. 按 F5 开始调试"
echo "3. 选择 '启动 Restart Life API' 配置"
echo "4. 在代码中设置断点"
echo "5. 使用 .vscode/api-test.http 测试 API"
echo ""
echo "详细说明请查看: .vscode/DEBUG_GUIDE.md"
