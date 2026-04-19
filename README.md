# 《重启人生》API服务

基于Go和Gin框架的高性能RESTful API服务，为《重启人生》游戏提供后端支持。

## 🎮 项目概述

《重启人生》是一款文字模拟人生游戏的后端API服务，提供：
- 用户认证和角色管理
- 游戏逻辑和事件处理
- 数据持久化和缓存
- 多平台客户端支持

## 🚀 快速开始

### 环境要求
- Go 1.23.8+
- MySQL 8.0
- Redis 7.0
- Docker & Docker Compose

### 本地开发
```bash
# 安装依赖
go mod tidy

# 启动数据库服务
docker-compose up -d postgres redis

# 启动API服务
go run cmd/server/main.go
```

### Docker部署（推荐）

使用我们提供的便捷脚本：

**Linux/macOS用户：**
```bash
# 给脚本添加执行权限
chmod +x scripts/*.sh

# 构建镜像
ENV_FILE=../.env.test ./scripts/build.sh

# 启动所有服务
ENV_FILE=../.env.test ./scripts/start.sh

# 停止服务
ENV_FILE=../.env.test ./scripts/stop.sh

# 使用腾讯云镜像优化构建和启动
ENV_FILE=../.env.test ./scripts/build.sh tencent
ENV_FILE=../.env.test ./scripts/start.sh tencent

# 使用中国网络优化构建和启动
ENV_FILE=../.env.test ./scripts/build.sh china
ENV_FILE=../.env.test ./scripts/start.sh china
```

**测试环境build image**
```
ENV_FILE=../.env.test docker-compose -f docker/docker-compose.tencent.yml build
```
**生产环境build image**
```
ENV_FILE=../.env.live docker-compose -f docker/docker-compose.tencent.yml build
```

**测试环境部署+运行**
```
ENV_FILE=../.env.test docker-compose -f docker/docker-compose.tencent.yml up -d
```
**生产环境部署+运行**
```
ENV_FILE=../.env.live docker-compose -f docker/docker-compose.tencent.yml up -d
```

**测试环境关闭**
```
ENV_FILE=../.env.test docker-compose -f docker/docker-compose.tencent.yml down
```
**生产环境关闭**
```
ENV_FILE=../.env.live docker-compose -f docker/docker-compose.tencent.yml down
```

**Windows用户：**
```cmd
# 构建镜像（官方源）
scripts\build.bat
# 构建镜像（腾讯云镜像）
scripts\build.bat tencent
# 构建镜像（中国网络优化）
scripts\build.bat china

# 启动所有服务（官方源）
scripts\start.bat
# 启动所有服务（腾讯云镜像）
scripts\start.bat tencent
# 启动所有服务（中国网络优化）
scripts\start.bat china

# 停止服务
scripts\stop.bat
```

服务启动后访问：
- 🚀 API服务: http://localhost:8080
- 🗄️ 数据库管理: http://localhost:8081
- 🔴 Redis管理: http://localhost:8082

详细使用说明请参考 [Docker开发指南](DOCKER.md)。

## 📂 项目结构

详见项目目录结构和技术文档。

## 🔗 相关链接

- [产品需求文档](prdtd/PRD.md)
- [后端技术设计文档](prdtd/后端技术设计文档_Backend_TD.md)
- [游戏规则设计](regulations/regulation.md)
- [前端Unity项目仓库](https://github.com/your-org/restart-life-unity)

## 🤝 贡献指南

欢迎提交Issue和Pull Request！

## 📄 许可证

MIT License

## ✨ 已完成功能

### 🔐 用户认证系统 (v0.2.0)
- ✅ 用户注册/登录
- ✅ JWT Token认证
- ✅ 访问Token和刷新Token
- ✅ 密码安全加密（bcrypt）
- ✅ 用户信息管理
- ✅ 密码修改
- ✅ 多平台登录支持（用户名/邮箱）
- ✅ Token中间件保护

详细文档: [认证系统说明](docs/auth_system.md)

## Dev Environment Self-Check (A-001)
Run one command before development to verify local prerequisites:

```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\dev-check.ps1
```

Useful options:

```powershell
# Skip service probes when MySQL/Redis are not running yet
powershell -ExecutionPolicy Bypass -File .\scripts\dev-check.ps1 -SkipMySQL -SkipRedis
```

The script checks:
- Go version
- Node.js + npm availability
- MySQL TCP connectivity
- Redis TCP connectivity

Non-zero exit code means at least one required check failed.

## Workflow Gate Check (A-004)
Run workflow gate check locally:

```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\check-workflow-gates.ps1 -FeatureId restart-life-api-baseline -WorkflowDocsRoot ..\workflow-docs
```

CI is configured to run this gate check as a blocking job.

## Script-Based Validation (No Local Test EXE)
If your machine policy blocks generated `.exe` files during local testing, use script-based validation:

```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\script-test.ps1 -FeatureId restart-life-api-baseline -WorkflowDocsRoot ..\workflow-docs -SkipMySQL -SkipRedis
```

This command runs:
- `scripts/dev-check.ps1`
- `scripts/check-workflow-gates.ps1`
