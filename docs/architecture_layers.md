# 分层架构重构说明

## 🏗️ 新的分层架构

重构后的项目采用了标准的四层架构模式：

```
┌─────────────────┐
│   Handler Layer │  ← HTTP处理层 (Controllers)
├─────────────────┤
│   Service Layer │  ← 业务逻辑层
├─────────────────┤
│Repository Layer │  ← 仓储抽象层
├─────────────────┤
│     DAO Layer   │  ← 数据访问对象层
├─────────────────┤
│  Database Layer │  ← 数据库连接层
└─────────────────┘
```

## 📁 目录结构

```
internal/
├── api/handlers/        # Handler Layer - HTTP请求处理
├── services/           # Service Layer - 业务逻辑
├── repository/         # Repository Layer - 仓储抽象
├── dao/               # DAO Layer - 数据访问对象
├── database/          # Database Layer - 数据库连接管理
├── models/            # 数据模型
└── utils/             # 工具类
```

## 🔄 各层职责

### 1. Handler Layer (API层)
- **位置**: `internal/api/handlers/`
- **职责**:
  - 处理HTTP请求和响应
  - 参数验证和数据绑定
  - 调用Service层处理业务逻辑
  - 错误处理和状态码设置

### 2. Service Layer (服务层)
- **位置**: `internal/services/`
- **职责**:
  - 核心业务逻辑处理
  - 数据验证和业务规则检查
  - 事务管理
  - 调用Repository层进行数据操作

### 3. Repository Layer (仓储层)
- **位置**: `internal/repository/`
- **职责**:
  - 提供数据访问的抽象接口
  - 业务相关的数据组合和处理
  - 调用DAO层进行具体数据操作
  - 数据模型的业务逻辑包装

### 4. DAO Layer (数据访问对象层)
- **位置**: `internal/dao/`
- **职责**:
  - 直接的数据库操作 (CRUD)
  - SQL查询的封装
  - 数据库事务处理
  - 原始数据的读写

### 5. Database Layer (数据库层)
- **位置**: `internal/database/`
- **职责**:
  - 数据库连接管理
  - 连接池配置
  - 数据库初始化

## 🌊 数据流向

### 请求流向 (下行)
```
HTTP Request → Handler → Service → Repository → DAO → Database
```

### 响应流向 (上行)
```
Database → DAO → Repository → Service → Handler → HTTP Response
```

## 🔧 依赖注入

容器 (`cmd/server/container.go`) 中的初始化顺序：

```go
1. initUtils()         // 工具类初始化
2. initDAOs()         // DAO层初始化
3. initRepositories() // Repository层初始化
4. initServices()     // Service层初始化
5. initMiddleware()   // 中间件初始化
6. initHandlers()     // Handler层初始化
```

## ✅ 重构优势

1. **职责分离**: 每层有明确的职责边界
2. **可测试性**: 各层可以独立进行单元测试
3. **可维护性**: 代码结构清晰，易于维护和扩展
4. **可扩展性**: 新增功能时可以按层次添加
5. **复用性**: DAO层的方法可以被多个Repository复用
6. **解耦**: 上层不直接依赖数据库，通过接口抽象

## 🎯 Context支持

所有层都支持 `context.Context`，实现：
- 请求超时控制
- 请求取消机制
- 分布式追踪
- 资源清理

## 📝 示例

### Handler调用Service
```go
user, err := h.authService.GetProfile(c.Request.Context(), userID)
```

### Service调用Repository
```go
user, err := s.userRepo.GetByID(ctx, userID)
```

### Repository调用DAO
```go
user, err := r.userDAO.SelectByID(ctx, userID)
```

这种分层架构符合企业级应用开发的最佳实践，提供了良好的可维护性和扩展性。
