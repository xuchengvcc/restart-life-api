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

### Docker部署
```bash
# 构建并启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps
```

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
