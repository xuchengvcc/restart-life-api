# Step 2: 核心业务模块 (v0.2.x)

## 📋 阶段概述
- **阶段目标**: 实现游戏核心业务逻辑，包括用户认证、角色管理和基础游戏功能
- **预计时间**: 2025-02-03 ~ 2025-02-16 (2周)
- **关键交付**: 完整的用户系统、角色系统和游戏核心逻辑

## 🎯 详细任务

### 任务4: 用户认证系统 ✅ **已完成**
- **分支**: `xucheng/feature/v0.2/auth-system`
- **负责人**: xucheng
- **预计时间**: 3-4天
- **完成时间**: 2025-07-07

#### 开发目标
- [x] JWT认证实现
- [x] 用户注册/登录接口
- [x] 密码加密和验证
- [x] 多平台认证支持

#### 详细任务清单
- [x] 实现JWT Token生成和验证
- [x] 实现用户注册接口（邮箱验证）
- [x] 实现用户登录接口（用户名/邮箱）
- [x] 实现密码哈希和验证（bcrypt）
- [x] 实现Token刷新机制
- [x] 实现退出登录功能
- [x] 实现用户信息查询和更新
- [x] 支持第三方登录预留接口
- [x] 实现两步重置密码功能

#### 交付物
- [x] 用户认证API (`internal/api/handlers/auth.go`)
- [x] JWT令牌管理 (`internal/utils/jwt.go`)
- [x] 密码安全模块 (`internal/utils/password.go`)
- [x] 认证中间件 (`internal/api/middleware/auth.go`)
- [x] 用户服务层 (`internal/services/auth_service.go`)
- [x] 用户数据层 (`internal/repository/user_repository.go`)
- [x] 邮箱验证码系统 (`internal/services/verification_code_service.go`)
- [x] 两步重置密码API (`/verify-code`, `/reset-password`)

#### 已实现的API接口
```http
POST /api/v1/auth/register              # 用户注册 ✅
POST /api/v1/auth/login                 # 用户登录 ✅
POST /api/v1/auth/logout                # 用户登出 ✅
POST /api/v1/auth/refresh               # 刷新Token ✅
GET  /api/v1/auth/profile               # 获取用户信息 ✅
PUT  /api/v1/auth/profile               # 更新用户信息 ✅
POST /api/v1/auth/change-password       # 修改密码 ✅
POST /api/v1/auth/send-verification-code # 发送验证码 ✅
POST /api/v1/auth/verify-code           # 验证验证码并获取重置令牌 ✅
POST /api/v1/auth/reset-password        # 使用令牌重置密码 ✅
```

#### 数据模型
```go
type User struct {
    UserID       string    `json:"user_id" db:"user_id"`
    Username     string    `json:"username" db:"username"`
    Email        string    `json:"email" db:"email"`
    PasswordHash string    `json:"-" db:"password_hash"`
    CreatedAt    time.Time `json:"created_at" db:"created_at"`
    UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
    LastLogin    *time.Time `json:"last_login" db:"last_login"`
    IsActive     bool      `json:"is_active" db:"is_active"`
}

type LoginRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
    Username string `json:"username" binding:"required,min=3,max=50"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}
```

#### 🎉 完成摘要

**用户认证系统已全面完成！**

✅ **核心功能**
- JWT认证（访问令牌 + 刷新令牌）
- 用户注册/登录（支持用户名和邮箱）
- 密码安全（bcrypt加密，强度可配置）
- 用户资料管理（查询、更新、修改密码）
- 邮箱验证码系统（Redis存储，自动过期）
- 两步重置密码（验证码→重置令牌→新密码）

✅ **安全特性**
- JWT令牌安全验证
- 密码强度校验
- 邮箱格式验证
- 重放攻击防护
- 一次性验证码/令牌
- 用户状态检查

✅ **技术架构**
- 分层架构（Handler→Service→Repository→DAO）
- 依赖注入和容器管理
- 中间件支持（认证、CORS、日志等）
- 错误统一处理
- 配置化管理

✅ **测试验证**
- 单元测试（密码管理、JWT管理）
- API集成测试脚本
- 两步重置密码流程测试
- 文档和示例完备

#### 📊 技术指标
- API响应时间：< 100ms
- 密码哈希强度：bcrypt cost 12
- JWT令牌有效期：访问令牌1小时，刷新令牌24小时
- 验证码有效期：5分钟
- 重置令牌有效期：10分钟

#### 📋 相关文档
- [认证系统文档](../auth_system.md)
- [两步重置密码API文档](../two_step_password_reset_api.md)
- [认证实现报告](../auth_implementation_report.md)
- [API测试脚本](../../scripts/test_two_step_password_reset.sh)

---

### 任务5: 角色管理系统 ✅ **已完成**
- **分支**: `xucheng/feature/v0.2/character-system`
- **负责人**: xucheng
- **预计时间**: 4-5天
- **完成时间**: 2025-07-12

#### 开发目标
- [x] 角色创建接口
- [x] 角色属性管理
- [x] 角色数据存储和查询
- [x] 角色状态更新

#### 详细任务清单
- [x] 实现角色创建接口（随机生成属性）
- [x] 实现角色列表查询
- [x] 实现角色详情查询
- [x] 实现角色删除功能
- [x] 实现角色属性更新
- [x] 实现角色关系管理
- [x] 实现角色状态追踪
- [x] 角色数据验证和校验

#### 交付物
- [x] 角色管理API (`internal/api/handlers/character.go`)
- [x] 角色数据模型 (`internal/models/character.go`)
- [x] 角色服务层 (`internal/services/character_service.go`)
- [x] 角色数据层 (`internal/repository/character_repository.go`)
- [x] 角色属性计算逻辑
- [x] 角色关系管理模块

#### 已实现的API接口
```http
POST /api/v1/characters/create              # 创建角色 ✅
GET  /api/v1/characters/list                # 获取角色列表 ✅
GET  /api/v1/characters/get/:id             # 获取角色详情 ✅
PUT  /api/v1/characters/update/:id          # 更新角色信息 ✅
DELETE /api/v1/characters/delete/:id        # 删除角色 ✅
GET  /api/v1/characters/attributes/get/:id  # 获取角色属性 ✅
PUT  /api/v1/characters/attributes/update/:id # 更新角色属性 ✅
```

#### 数据模型
```go
type Character struct {
    CharacterID   string              `json:"character_id" db:"character_id"`
    UserID        uint                `json:"user_id" db:"user_id"`
    CharacterName string              `json:"character_name" db:"character_name"`
    BirthCountry  string              `json:"birth_country" db:"birth_country"`
    BirthYear     int                 `json:"birth_year" db:"birth_year"`
    CurrentAge    int                 `json:"current_age" db:"current_age"`
    Gender        int                 `json:"gender" db:"gender"`
    Race          int                 `json:"race" db:"race"`
    Attributes    CharacterAttributes `json:"attributes"`
    IsActive      bool                `json:"is_active" db:"is_active"`
    CreatedAt     int64               `json:"created_at" db:"created_at"`
    UpdatedAt     int64               `json:"updated_at" db:"updated_at"`

    // 扩展属性
    LifeStage           string  `json:"life_stage" db:"life_stage"`
    CurrentStatus       string  `json:"current_status" db:"current_status"`
    HappinessLevel      int     `json:"happiness_level" db:"happiness_level"`
    HealthLevel         int     `json:"health_level" db:"health_level"`
    Money               int64   `json:"money" db:"money"`
    // ... 更多字段
}

type CharacterAttributes struct {
    Intelligence          int `json:"intelligence" db:"intelligence"`           // 智力 0-100
    EmotionalIntelligence int `json:"emotional_intelligence" db:"emotional_intelligence"` // 情商 0-100
    Memory                int `json:"memory" db:"memory"`                       // 记忆力 0-100
    Imagination           int `json:"imagination" db:"imagination"`             // 想象力 0-100
    PhysicalFitness       int `json:"physical_fitness" db:"physical_fitness"`   // 体质 0-100
    Appearance            int `json:"appearance" db:"appearance"`               // 外貌 0-100
}

type CreateCharacterRequest struct {
    CharacterName string `json:"character_name" binding:"required,max=100"`
    BirthCountry  string `json:"birth_country" binding:"required"`
    BirthYear     int    `json:"birth_year" binding:"required,min=1800,max=2050"`
    Gender        int    `json:"gender" binding:"required,min=0,max=3"`
    Race          int    `json:"race" binding:"required,min=0"`
}
```

#### 🎉 完成摘要

**角色管理系统已全面完成！**

✅ **核心功能**
- 角色创建（支持随机属性生成）
- 角色查询（详情、列表、活跃角色筛选）
- 角色更新（基本信息和属性分离更新）
- 角色删除（软删除机制）
- 权限验证（用户只能操作自己的角色）
- 属性管理（6维属性系统）

✅ **数据管理**
- 完整的分层架构（Handler→Service→Repository→DAO）
- MySQL数据持久化存储
- 数据验证和约束检查
- 事务处理和错误恢复
- 索引优化和查询性能

✅ **业务逻辑**
- 随机属性生成算法
- 角色生命周期管理
- 多角色支持
- 角色状态追踪
- 扩展属性支持

✅ **技术特性**
- RESTful API设计
- JSON数据格式
- 统一错误处理
- 结构化日志记录
- 参数验证和绑定

✅ **测试验证**
- 完整的API测试脚本
- 角色CRUD操作测试
- 权限验证测试
- 数据完整性测试
- 错误场景测试

#### 📊 技术指标
- API响应时间：< 200ms
- 角色属性范围：30-70（初始随机值）
- 支持并发角色操作
- 数据库事务一致性保证
- 软删除机制保证数据可恢复

#### 📋 相关文档
- [角色管理API测试脚本](../../scripts/test_character_api.sh)
- [角色数据模型文档](../../internal/models/character.go)
- [数据库迁移脚本](../../migrations/000002_create_characters_table.up.sql)

---

### 任务6: 游戏核心逻辑
- **分支**: `xucheng/feature/v0.2/game-engine`
- **负责人**: xucheng
- **预计时间**: 5-6天

#### 开发目标
- [ ] 人生推进系统
- [ ] 事件生成系统
- [ ] 决策选择系统
- [ ] 属性变化计算

#### 详细任务清单
- [ ] 实现年龄推进机制
- [ ] 实现随机事件生成算法
- [ ] 实现人生阶段识别和切换
- [ ] 实现决策选择逻辑
- [ ] 实现属性变化计算
- [ ] 实现事件结果处理
- [ ] 实现游戏状态管理
- [ ] 实现存档点管理

#### 交付物
- [ ] 游戏引擎核心 (`internal/services/game_service.go`)
- [ ] 事件系统 (`internal/services/event_service.go`)
- [ ] 决策处理逻辑 (`internal/models/decision.go`)
- [ ] 属性计算引擎 (`internal/utils/attribute_calculator.go`)
- [ ] 游戏状态管理 (`internal/models/game_state.go`)
- [ ] 游戏API接口 (`internal/api/handlers/game.go`)

#### API接口设计
```http
POST /api/v1/game/start/:character_id    # 开始游戏
POST /api/v1/game/advance/:character_id  # 推进一年
GET  /api/v1/game/state/:character_id    # 获取游戏状态
POST /api/v1/game/decision/:character_id # 做出决策
GET  /api/v1/game/events/:character_id   # 获取事件历史
POST /api/v1/game/save/:character_id     # 保存游戏
POST /api/v1/game/load/:character_id     # 加载游戏
```

#### 数据模型
```go
type GameState struct {
    CharacterID   string          `json:"character_id" db:"character_id"`
    CurrentAge    int             `json:"current_age" db:"current_age"`
    LifeStage     string          `json:"life_stage" db:"life_stage"`
    KeyEvents     []Event         `json:"key_events"`
    PendingDecisions *Decision   `json:"pending_decisions"`
    Attributes    CharacterAttributes `json:"attributes"`
    Relationships string  `json:"relationships"`
    LastSaveTime  time.Time      `json:"last_save_time" db:"last_save_time"`
}

type Event struct {
    EventID     string    `json:"event_id" db:"event_id"`
    CharacterID string    `json:"character_id" db:"character_id"`
    Description string    `json:"description" db:"description"`
    Age         int       `json:"age" db:"age"`
    Impact      string    `json:"impact"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type Decision struct {
	CharacterID      string          `json:"character_id" db:"character_id"`
	Options          *DecisionOption `json:"options"`
	PreviousOptions  *DecisionOption `json:"previous_options" db:"previous_options"`   // 上一个决策选项
	PreviousDecision *int8           `json:"previous_decision" db:"previous_decision"` // 上一个决策 1: conservative, 2: moderate, 3: aggressive
	CreatedAt        int64           `json:"created_at" db:"created_at"`
	UpdatedAt        int64           `json:"updated_at" db:"updated_at"`
}

// DecisionOption 决策选项
type DecisionOption struct {
	Conservative DecisionDetails `json:"conservative"` // 保守选项
	Moderate     DecisionDetails `json:"moderate"`     // 中庸选项
	Aggressive   DecisionDetails `json:"aggressive"`   // 激进选项
}
```

#### 游戏核心算法
```go
// 事件生成算法
func GenerateEvents(character *Character, gameState *GameState) []Event

// 属性变化计算
func CalculateAttributeChanges(event Event, decision Decision) AttributeChanges

// 人生阶段判断
func DetermineLifeStage(age int) LifeStage

// 决策成功率计算
func CalculateSuccessRate(option DecisionOption, attributes CharacterAttributes) float64
```

---

## 📊 阶段验收标准

### 功能验收
- [x] 用户能够成功注册和登录
- [x] JWT认证系统工作正常
- [x] 用户能够创建和管理角色
- [x] 角色属性系统正常运行
- [ ] 游戏推进逻辑正确
- [ ] 事件生成系统工作正常
- [ ] 决策系统能够正确处理选择

### 技术验收
- [x] 所有API接口返回标准格式
- [x] 数据验证和错误处理完善
- [ ] 单元测试覆盖率 > 80%
- [x] API响应时间 < 500ms
- [ ] 并发处理能力测试通过
- [x] 数据库事务处理正确

### 安全验收
- [x] 密码安全存储（bcrypt）
- [x] JWT Token安全验证
- [x] API访问权限控制正确
- [x] 敏感信息不在日志中暴露
- [x] SQL注入防护有效

---

## 🧪 测试规划

### 单元测试
```go
// 认证系统测试 ✅
- TestUserRegistration
- TestUserLogin
- TestJWTGeneration
- TestPasswordHashing

// 角色系统测试 ✅
- TestCharacterCreation
- TestAttributeGeneration
- TestCharacterQuery
- TestCharacterUpdate

// 游戏系统测试
- TestEventGeneration
- TestDecisionProcessing
- TestAttributeCalculation
- TestGameStateUpdate
```

### 集成测试
```go
// API集成测试 ✅
- TestAuthEndpoints
- TestCharacterEndpoints
- TestGameEndpoints

// 数据库集成测试 ✅
- TestUserRepository
- TestCharacterRepository
- TestGameRepository
```

### 性能测试
- 并发用户登录测试（100用户）
- 角色创建性能测试
- 游戏推进响应时间测试
- 数据库查询性能测试

---

## 🔧 技术规范

### 新增依赖
```go
// 认证相关
- github.com/golang-jwt/jwt/v4      // JWT处理
- golang.org/x/crypto/bcrypt        // 密码加密

// 验证相关
- github.com/go-playground/validator/v10 // 数据验证
- github.com/gin-gonic/gin/binding  // 参数绑定

// 工具库
- github.com/google/uuid           // UUID生成
- github.com/shopspring/decimal    // 精确计算
```

### 错误处理规范
```go
type APIError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
}

type APIResponse struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   *APIError   `json:"error,omitempty"`
}
```

### 日志规范
```go
// 使用结构化日志
log.WithFields(logrus.Fields{
    "user_id": userID,
    "action":  "login",
    "ip":      clientIP,
}).Info("User login successful")
```

---

## ⏭️ 下一步计划

完成Step 2后，将进入Step 3：高级功能模块开发 (v0.3.x)
- 关系网络系统
- 成就系统
- 存档与统计系统

---

*创建时间: 2025-01-26*
*最后更新: 2025-07-12*
