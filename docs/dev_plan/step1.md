# Step 1: 基础架构搭建 (v0.1.x)

## 📋 阶段概述
- **阶段目标**: 搭建项目基础架构，建立开发环境和基础服务框架
- **预计时间**: 2025-01-26 ~ 2025-02-02 (1周)
- **关键交付**: 可运行的基础Web服务、数据库连接、项目架构

## 🎯 详细任务

### 任务1: 项目初始化与基础配置
- **分支**: `xucheng/feature/v0.1/init-project-structure`
- **负责人**: xucheng
- **预计时间**: 1-2天

#### 开发目标
- [x] 创建开发分支
- [x] 搭建Go项目基础架构
- [x] 配置Go modules和依赖管理
- [x] 创建标准的项目目录结构
- [x] 配置基础的开发环境

#### 详细任务清单
- [x] 初始化go.mod文件，配置项目依赖
- [x] 创建完整的目录结构（按照TD文档规范）
- [x] 编写基础的main.go入口文件
- [x] 配置开发和生产环境配置文件
- [x] 创建基础的Makefile或构建脚本
- [x] 配置.gitignore文件

#### 交付物
- [x] 开发分支创建
- [x] 完整的项目目录结构
- [x] go.mod配置文件
- [x] 基础的main.go入口文件
- [x] 开发环境配置文件
- [x] 构建和开发脚本

#### 项目目录结构
```
restart-life-api/
├── cmd/
│   └── server/
│       └── main.go              # 应用程序入口点
├── internal/
│   ├── api/
│   │   ├── handlers/            # HTTP处理器
│   │   ├── middleware/          # 中间件
│   │   └── routes/              # 路由定义
│   ├── models/                  # 数据模型
│   ├── services/                # 业务逻辑服务
│   ├── repository/              # 数据访问层
│   │   ├── postgres/
│   │   └── redis/
│   └── config/                  # 配置管理
├── pkg/
│   ├── utils/                   # 工具函数
│   └── database/                # 数据库连接
├── migrations/                  # 数据库迁移文件
├── docs/                        # API文档
├── docker/
│   └── Dockerfile
├── configs/                     # 配置文件
├── scripts/                     # 构建脚本
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

---

### 任务2: 数据库设计与连接
- **分支**: `xucheng/feature/v0.1/database-setup`
- **负责人**: xucheng
- **预计时间**: 2-3天

#### 开发目标
- [ ] PostgreSQL数据库连接配置
- [ ] Redis缓存连接配置
- [ ] 数据库迁移系统搭建
- [ ] 核心数据表结构设计

#### 详细任务清单
- [ ] 实现PostgreSQL连接池管理
- [ ] 实现Redis连接和基础操作
- [ ] 创建数据库迁移系统（支持up/down）
- [ ] 设计并实现核心数据表结构
- [ ] 编写数据库连接健康检查
- [ ] 配置数据库连接参数和环境变量

#### 交付物
- [ ] 数据库连接模块 (`pkg/database/`)
- [ ] 数据库迁移脚本 (`migrations/`)
- [ ] 核心表结构SQL文件
- [ ] 数据库配置文档
- [ ] 数据库健康检查接口

#### 核心数据表
```sql
-- 用户账户表
CREATE TABLE users (
    user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP,
    is_active BOOLEAN DEFAULT true
);

-- 游戏角色表
CREATE TABLE characters (
    character_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    character_name VARCHAR(100) NOT NULL,
    birth_country VARCHAR(100) NOT NULL,
    birth_year INTEGER NOT NULL,
    current_age INTEGER NOT NULL DEFAULT 0,
    gender VARCHAR(20) NOT NULL,
    race VARCHAR(50) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 更多表结构见完整迁移文件...
```

---

### 任务3: 基础服务框架
- **分支**: `xucheng/feature/v0.1/basic-framework`
- **负责人**: xucheng
- **预计时间**: 2-3天

#### 开发目标
- [ ] Gin Web框架集成
- [ ] 基础中间件实现（CORS、日志、错误处理）
- [ ] 配置管理系统（Viper）
- [ ] 健康检查接口

#### 详细任务清单
- [ ] 集成Gin框架，配置基础路由
- [ ] 实现CORS中间件，支持跨域请求
- [ ] 实现请求日志中间件
- [ ] 实现统一错误处理中间件
- [ ] 实现请求恢复中间件（panic recovery）
- [ ] 集成Viper配置管理
- [ ] 创建配置文件模板（dev/prod环境）
- [ ] 实现健康检查和服务状态接口

#### 交付物
- [ ] 完整的Web服务框架
- [ ] 中间件系统 (`internal/api/middleware/`)
- [ ] 配置管理模块 (`internal/config/`)
- [ ] 基础API接口 (`/health`, `/ping`)
- [ ] 服务启动和关闭逻辑

#### 基础中间件
```go
// 中间件列表
- CORS中间件 - 跨域支持
- Logger中间件 - 请求日志
- Recovery中间件 - 异常恢复
- RateLimit中间件 - 限流控制
- RequestID中间件 - 请求追踪
```

#### 健康检查接口
```
GET /health      - 服务健康状态
GET /ping        - 基础连通性检查
GET /ready       - 服务就绪状态
GET /metrics     - 基础指标信息
```

---

## 📊 阶段验收标准

### 功能验收
- [ ] 服务能正常启动和关闭
- [ ] 数据库连接正常，能执行基础CRUD
- [ ] Redis缓存连接正常，能读写数据
- [ ] 健康检查接口返回正确状态
- [ ] 基础中间件功能正常

### 技术验收
- [ ] 代码遵循Go规范，通过golint检查
- [ ] 单元测试覆盖率 > 70%
- [ ] 所有配置通过环境变量管理
- [ ] 日志格式统一，包含必要信息
- [ ] 错误处理统一，返回标准格式

### 文档验收
- [ ] 项目README完整，包含启动说明
- [ ] API文档基础框架搭建完成
- [ ] 数据库迁移文档完成
- [ ] 开发环境搭建文档完成

---

## 🔧 技术规范

### 依赖管理
```go
// 主要依赖包
- github.com/gin-gonic/gin          // Web框架
- github.com/lib/pq                 // PostgreSQL驱动
- github.com/go-redis/redis/v8      // Redis客户端
- github.com/spf13/viper            // 配置管理
- github.com/sirupsen/logrus        // 日志库
- github.com/golang-migrate/migrate // 数据库迁移
- github.com/google/uuid            // UUID生成
```

### 代码规范
- 使用gofmt格式化代码
- 遵循Go官方命名规范
- 每个公开函数需要注释
- 错误处理不能忽略
- 单元测试文件命名为 `*_test.go`

### 配置规范
```yaml
# config/development.yaml
server:
  port: 8080
  mode: debug
  
database:
  postgres:
    host: localhost
    port: 5432
    database: restart_life_dev
    username: postgres
    password: password
    sslmode: disable
    
  redis:
    host: localhost
    port: 6379
    password: ""
    database: 0
```

---

## ⏭️ 下一步计划

完成Step 1后，将进入Step 2：核心业务模块开发 (v0.2.x)
- 用户认证系统
- 角色管理系统  
- 游戏核心逻辑

---

*创建时间: 2025-01-26*
*最后更新: 2025-01-26*
