#!/bin/bash

# VS Code Go 语言服务器重置脚本
echo "=== 重置 VS Code Go 语言服务器 ==="

# 1. 检查 Go 环境
echo "1. 检查 Go 环境..."
go version
go env GOMOD GOPATH GOROOT

# 2. 清理模块缓存
echo ""
echo "2. 清理并重新下载依赖..."
go clean -modcache 2>/dev/null || true
go mod download
go mod tidy

# 3. 重新生成 go.sum
echo ""
echo "3. 验证模块..."
go mod verify

# 4. 安装/更新 Go 工具
echo ""
echo "4. 安装 Go 语言工具..."
go install -a golang.org/x/tools/gopls@latest
go install -a golang.org/x/tools/cmd/goimports@latest
go install -a github.com/go-delve/delve/cmd/dlv@latest
go install -a honnef.co/go/tools/cmd/staticcheck@latest

# 5. 构建项目以检查错误
echo ""
echo "5. 构建项目..."
go build ./...

# 6. 运行测试以确保一切正常
echo ""
echo "6. 运行测试..."
go test ./... -v

echo ""
echo "=== 重置完成 ==="
echo ""
echo "现在请："
echo "1. 重启 VS Code"
echo "2. 打开 restart-life-api.code-workspace 工作区文件"
echo "3. 等待 Go 语言服务器重新索引（可能需要几分钟）"
echo "4. 测试 Ctrl+Click 跳转功能"
