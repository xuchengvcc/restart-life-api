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
- Go 1.21+
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

**Windows用户：**
```cmd
# 构建镜像
scripts\build.bat

# 启动所有服务
scripts\start.bat

# 停止服务
scripts\stop.bat
```

**Linux/macOS用户：**
```bash
# 给脚本添加执行权限
chmod +x scripts/*.sh

# 构建镜像
./scripts/build.sh

# 启动所有服务
./scripts/start.sh

# 停止服务
./scripts/stop.sh
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
