# 认证系统快速启动指南

本指南帮助您快速测试已完成的用户认证系统。

## 1. 环境准备

### 1.1 确保MySQL运行
```bash
# 使用Docker启动MySQL（推荐）
docker run -d \
  --name mysql-dev \
  -e MYSQL_ROOT_PASSWORD=password \
  -e MYSQL_DATABASE=restart_life_dev \
  -p 3306:3306 \
  mysql:8.0

# 或者确保您的本地MySQL正在运行
```

### 1.2 确保Redis运行
```bash
# 使用Docker启动Redis（推荐）
docker run -d \
  --name redis-dev \
  -p 6379:6379 \
  redis:7.0

# 或者确保您的本地Redis正在运行
```

## 2. 初始化数据库

```bash
# 运行数据库初始化脚本
./scripts/init_db.sh

# 如果需要自定义数据库配置
DB_HOST=localhost DB_PORT=3306 DB_USER=root DB_PASSWORD=password DB_NAME=restart_life_dev ./scripts/init_db.sh
```

## 3. 启动API服务

```bash
# 方式1: 直接运行（需要安装Go）
go run cmd/server/main.go

# 方式2: 使用已编译的二进制（如果存在）
./build/restart-life-api

# 方式3: 使用Docker Compose
docker-compose up api
```

服务启动后会在 http://localhost:8080 监听。

## 4. 测试认证功能

### 4.1 运行自动化测试
```bash
# 运行单元测试
make test-auth

# 运行API集成测试（需要服务正在运行）
make test-auth-api

# 或直接运行测试脚本
./scripts/test_auth_api.sh
```

### 4.2 手动测试关键接口

**注册用户：**
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'
```

**用户登录：**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

**获取用户信息（需要替换TOKEN）：**
```bash
curl -X GET http://localhost:8080/api/v1/auth/profile \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## 5. 验证健康状态

```bash
# 检查服务健康状态
curl http://localhost:8080/health

# 预期响应
{
  "status": "healthy",
  "timestamp": "2025-01-26T10:00:00Z",
  "version": "v0.1.0"
}
```

## 6. 查看日志

```bash
# 如果使用Docker Compose
docker-compose logs -f api

# 如果直接运行
# 日志会输出到控制台
```

## 常见问题

### Q: MySQL连接失败
**A:** 检查MySQL是否正在运行，配置文件中的连接信息是否正确：
```bash
mysql -h localhost -P 3306 -u root -p
```

### Q: Redis连接失败
**A:** 检查Redis是否正在运行：
```bash
redis-cli ping
# 应该返回 PONG
```

### Q: 端口冲突
**A:** 检查8080端口是否被占用：
```bash
lsof -i :8080
# 或者在配置文件中修改端口
```

### Q: JWT Token错误
**A:** 确保配置文件中的`jwt_secret`不为空，且在生产环境使用强密钥。

## 下一步

认证系统完成后，您可以：

1. 继续开发角色管理系统（Step 2 任务5）
2. 实现游戏核心逻辑（Step 2 任务6）
3. 查看完整的开发计划：[docs/dev_plan/step2.md](docs/dev_plan/step2.md)

## 相关文档

- [认证系统详细说明](docs/auth_system.md)
- [开发计划 Step 2](docs/dev_plan/step2.md)
- [API文档](docs/api_documentation.md)
