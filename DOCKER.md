# Docker 开发环境使用指南

## 快速开始

### Windows 用户

1. **构建镜像**
   ```cmd
   scripts\build.bat
   ```

2. **启动服务**
   ```cmd
   scripts\start.bat
   ```

3. **停止服务**
   ```cmd
   scripts\stop.bat
   ```

### Linux/macOS 用户

1. **给脚本添加执行权限**
   ```bash
   chmod +x scripts/*.sh
   ```

2. **构建镜像**
   ```bash
   ./scripts/build.sh
   ```

3. **启动服务**
   ```bash
   ./scripts/start.sh
   ```

4. **停止服务**
   ```bash
   ./scripts/stop.sh
   ```

### 🇨🇳 中国用户优化

如果你在中国，网络访问Docker Hub较慢，可以使用腾讯云镜像：

**Windows用户：**
```cmd
scripts\build.bat tencent     # 使用腾讯云镜像构建
scripts\start.bat tencent     # 使用腾讯云镜像启动服务
```

**Linux/macOS用户：**
```bash
./scripts/build.sh tencent    # 使用腾讯云镜像构建
./scripts/start.sh tencent    # 使用腾讯云镜像启动服务
```

优化内容：
- 使用腾讯云Alpine镜像源
- 使用goproxy.cn作为Go模块代理
- 使用腾讯云容器镜像服务的MySQL/Redis镜像

## 服务说明

启动后可以访问以下服务：

- 🚀 **API 服务**: http://localhost:8080
- 🔍 **健康检查**: http://localhost:8080/health
- 🗄️ **数据库管理** (Adminer): http://localhost:8081
- 🔴 **Redis 管理** (Redis Commander): http://localhost:8082

## 目录结构

```
docker/
├── mysql/
│   ├── conf/my.cnf          # MySQL配置文件
│   └── init/01-init.sh      # MySQL初始化脚本
└── redis/
    └── redis.conf           # Redis配置文件

configs/
├── development.yaml         # 本地开发配置
├── docker.yaml             # Docker环境配置
└── production.yaml          # 生产环境配置

# Docker相关文件
Dockerfile                   # 应用镜像构建文件
Dockerfile.tencent          # 腾讯云镜像优化版Dockerfile
docker-compose.yml          # 服务编排文件
docker-compose.tencent.yml  # 腾讯云镜像优化版compose文件
.dockerignore               # Docker忽略文件

# 脚本文件
scripts/build.bat           # Windows构建脚本
scripts/build.sh            # Linux/macOS构建脚本
scripts/start.bat           # Windows启动脚本
scripts/start.sh            # Linux/macOS启动脚本
scripts/stop.bat            # Windows停止脚本
scripts/stop.sh             # Linux/macOS停止脚本
```

## 常用命令

### Docker Compose 命令

```bash
# 启动所有服务（后台运行）
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f app

# 重启服务
docker-compose restart

# 停止服务
docker-compose down

# 停止服务并删除数据卷
docker-compose down -v
```

### Docker 命令

```bash
# 查看容器状态
docker ps

# 进入应用容器
docker exec -it restart-life-api sh

# 进入MySQL容器
docker exec -it restart-mysql mysql -u root -p

# 进入Redis容器
docker exec -it restart-redis redis-cli

# 查看镜像
docker images

# 清理未使用的镜像
docker image prune -f
```

## 开发工作流

1. **首次启动**
   ```bash
   ./scripts/build.sh    # 构建镜像
   ./scripts/start.sh    # 启动服务
   ```

2. **代码修改后**
   ```bash
   ./scripts/build.sh    # 重新构建镜像
   docker-compose restart app  # 重启应用服务
   ```

3. **查看日志**
   ```bash
   docker-compose logs -f app
   ```

4. **停止服务**
   ```bash
   ./scripts/stop.sh
   ```

## 数据持久化

- **MySQL数据**: 存储在 `mysql_data` 数据卷中
- **Redis数据**: 存储在 `redis_data` 数据卷中
- **应用日志**: 映射到本地 `logs/` 目录

## 环境配置

应用会根据 `CONFIG_ENV` 环境变量加载对应的配置文件：

- `development` → `configs/development.yaml`
- `docker` → `configs/docker.yaml`
- `production` → `configs/production.yaml`

Docker环境默认使用 `docker` 配置，其中数据库主机名为服务名称（`mysql`, `redis`）。

## 故障排除

### 服务无法启动

1. 检查Docker是否运行
2. 检查端口是否被占用
3. 查看日志：`docker-compose logs`

### 数据库连接失败

1. 确保MySQL服务健康：`docker-compose ps`
2. 检查数据库配置
3. 查看MySQL日志：`docker-compose logs mysql`

### 应用构建失败

1. 检查Go版本和依赖
2. 清理Docker缓存：`docker system prune -f`
3. 重新构建：`./build.sh`

## 性能优化

### 开发环境优化

- 使用多阶段构建减小镜像大小
- 利用Docker层缓存加速构建
- 数据卷映射避免重复复制

### 资源限制

可以在 `docker-compose.yml` 中添加资源限制：

```yaml
deploy:
  resources:
    limits:
      cpus: '0.5'
      memory: 512M
    reservations:
      cpus: '0.25'
      memory: 256M
```

## 🔧 网络问题解决方案

### 常见网络错误

如果遇到以下错误：
```
failed to solve: alpine:latest: failed to resolve source metadata
```

这通常是因为网络连接问题导致无法访问Docker Hub。

### 解决方案

1. **使用中国镜像版本**（推荐）
   ```bash
   # 使用优化的中国镜像
   ./scripts/build.sh cn
   ./scripts/start.sh cn
   ```

2. **配置Docker镜像加速器**
   
   在Docker Desktop中配置镜像加速器：
   ```json
   {
     "registry-mirrors": [
       "https://mirror.ccs.tencentyun.com",
       "https://docker.mirrors.ustc.edu.cn",
       "https://hub-mirror.c.163.com"
     ]
   }
   ```

3. **手动拉取镜像**
   ```bash
   # 预先拉取所需镜像
   docker pull golang:1.23.8-alpine
   docker pull alpine:latest
   docker pull mysql:8.0
   docker pull redis:7-alpine
   ```

### 中国镜像版本的优势

- ✅ 使用阿里云Container Registry镜像
- ✅ 使用goproxy.cn作为Go模块代理
- ✅ 使用阿里云Alpine软件源
- ✅ 构建速度更快，成功率更高

### 网络测试

你可以使用以下命令测试网络连接：
```bash
# 测试Docker Hub连接
docker pull hello-world

# 测试Go代理连接
curl -I https://goproxy.cn

# 测试阿里云镜像
docker pull registry.cn-hangzhou.aliyuncs.com/acs/alpine:latest
```
