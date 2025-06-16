# Step 2: 核心业务模块 (v0.2.x)

## 📋 阶段概述
- **阶段目标**: 实现游戏核心业务逻辑，包括用户认证、角色管理和基础游戏功能
- **预计时间**: 2025-02-03 ~ 2025-02-16 (2周)
- **关键交付**: 完整的用户系统、角色系统和游戏核心逻辑

## 🎯 详细任务

### 任务4: 用户认证系统
- **分支**: `xucheng/feature/v0.2/auth-system`
- **负责人**: xucheng
- **预计时间**: 3-4天

#### 开发目标
- [ ] JWT认证实现
- [ ] 用户注册/登录接口
- [ ] 密码加密和验证
- [ ] 多平台认证支持

#### 详细任务清单
- [ ] 实现JWT Token生成和验证
- [ ] 实现用户注册接口（邮箱验证）
- [ ] 实现用户登录接口（用户名/邮箱）
- [ ] 实现密码哈希和验证（bcrypt）
- [ ] 实现Token刷新机制
- [ ] 实现退出登录功能
- [ ] 实现用户信息查询和更新
- [ ] 支持第三方登录预留接口

#### 交付物
- [ ] 用户认证API (`internal/api/handlers/auth.go`)
- [ ] JWT令牌管理 (`pkg/utils/jwt.go`)
- [ ] 密码安全模块 (`pkg/utils/password.go`)
- [ ] 认证中间件 (`internal/api/middleware/auth.go`)
- [ ] 用户服务层 (`internal/services/auth_service.go`)
- [ ] 用户数据层 (`internal/repository/postgres/user_repo.go`)

#### API接口设计
```http
POST /api/v1/auth/register    # 用户注册
POST /api/v1/auth/login       # 用户登录
POST /api/v1/auth/logout      # 用户登出
POST /api/v1/auth/refresh     # 刷新Token
GET  /api/v1/auth/profile     # 获取用户信息
PUT  /api/v1/auth/profile     # 更新用户信息
POST /api/v1/auth/change-password # 修改密码
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

---

### 任务5: 角色管理系统
- **分支**: `xucheng/feature/v0.2/character-system`
- **负责人**: xucheng
- **预计时间**: 4-5天

#### 开发目标
- [ ] 角色创建接口
- [ ] 角色属性管理
- [ ] 角色数据存储和查询
- [ ] 角色状态更新

#### 详细任务清单
- [ ] 实现角色创建接口（随机生成属性）
- [ ] 实现角色列表查询
- [ ] 实现角色详情查询
- [ ] 实现角色删除功能
- [ ] 实现角色属性更新
- [ ] 实现角色关系管理
- [ ] 实现角色状态追踪
- [ ] 角色数据验证和校验

#### 交付物
- [ ] 角色管理API (`internal/api/handlers/character.go`)
- [ ] 角色数据模型 (`internal/models/character.go`)
- [ ] 角色服务层 (`internal/services/character_service.go`)
- [ ] 角色数据层 (`internal/repository/postgres/character_repo.go`)
- [ ] 角色属性计算逻辑
- [ ] 角色关系管理模块

#### API接口设计
```http
POST /api/v1/characters          # 创建角色
GET  /api/v1/characters          # 获取角色列表
GET  /api/v1/characters/:id      # 获取角色详情
PUT  /api/v1/characters/:id      # 更新角色信息
DELETE /api/v1/characters/:id    # 删除角色
GET  /api/v1/characters/:id/attributes # 获取角色属性
PUT  /api/v1/characters/:id/attributes # 更新角色属性
GET  /api/v1/characters/:id/relationships # 获取角色关系
```

#### 数据模型
```go
type Character struct {
    CharacterID   string            `json:"character_id" db:"character_id"`
    UserID        string            `json:"user_id" db:"user_id"`
    CharacterName string            `json:"character_name" db:"character_name"`
    BirthCountry  string            `json:"birth_country" db:"birth_country"`
    BirthYear     int               `json:"birth_year" db:"birth_year"`
    CurrentAge    int               `json:"current_age" db:"current_age"`
    Gender        string            `json:"gender" db:"gender"`
    Race          string            `json:"race" db:"race"`
    Attributes    CharacterAttributes `json:"attributes"`
    IsActive      bool              `json:"is_active" db:"is_active"`
    CreatedAt     time.Time         `json:"created_at" db:"created_at"`
    UpdatedAt     time.Time         `json:"updated_at" db:"updated_at"`
}

type CharacterAttributes struct {
    Intelligence int `json:"intelligence" db:"intelligence"` // 智力 0-100
    Physical     int `json:"physical" db:"physical"`         // 体质 0-100
    Charm        int `json:"charm" db:"charm"`               // 魅力 0-100
    Willpower    int `json:"willpower" db:"willpower"`       // 意志力 0-100
    Creativity   int `json:"creativity" db:"creativity"`     // 创造力 0-100
}

type CreateCharacterRequest struct {
    CharacterName string `json:"character_name" binding:"required,max=100"`
    BirthCountry  string `json:"birth_country" binding:"required"`
    BirthYear     int    `json:"birth_year" binding:"required,min=1800,max=2050"`
    Gender        string `json:"gender" binding:"required,oneof=male female"`
    Race          string `json:"race" binding:"required"`
}
```

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
- [ ] 属性计算引擎 (`pkg/utils/attribute_calculator.go`)
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
    CurrentEvents []Event         `json:"current_events"`
    PendingDecisions []Decision   `json:"pending_decisions"`
    Attributes    CharacterAttributes `json:"attributes"`
    Relationships []Relationship  `json:"relationships"`
    LastSaveTime  time.Time      `json:"last_save_time" db:"last_save_time"`
}

type Event struct {
    EventID     string    `json:"event_id" db:"event_id"`
    CharacterID string    `json:"character_id" db:"character_id"`
    EventType   string    `json:"event_type" db:"event_type"`
    Title       string    `json:"title" db:"title"`
    Description string    `json:"description" db:"description"`
    Age         int       `json:"age" db:"age"`
    Impact      EventImpact `json:"impact"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type Decision struct {
    DecisionID  string    `json:"decision_id" db:"decision_id"`
    EventID     string    `json:"event_id" db:"event_id"`
    Question    string    `json:"question" db:"question"`
    Options     []DecisionOption `json:"options"`
    Deadline    *time.Time `json:"deadline" db:"deadline"`
}

type DecisionOption struct {
    OptionID    string    `json:"option_id"`
    Text        string    `json:"text"`
    Consequences string  `json:"consequences"`
    Requirements map[string]int `json:"requirements"`
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
- [ ] 用户能够成功注册和登录
- [ ] JWT认证系统工作正常
- [ ] 用户能够创建和管理角色
- [ ] 角色属性系统正常运行
- [ ] 游戏推进逻辑正确
- [ ] 事件生成系统工作正常
- [ ] 决策系统能够正确处理选择

### 技术验收
- [ ] 所有API接口返回标准格式
- [ ] 数据验证和错误处理完善
- [ ] 单元测试覆盖率 > 80%
- [ ] API响应时间 < 500ms
- [ ] 并发处理能力测试通过
- [ ] 数据库事务处理正确

### 安全验收
- [ ] 密码安全存储（bcrypt）
- [ ] JWT Token安全验证
- [ ] API访问权限控制正确
- [ ] 敏感信息不在日志中暴露
- [ ] SQL注入防护有效

---

## 🧪 测试规划

### 单元测试
```go
// 认证系统测试
- TestUserRegistration
- TestUserLogin
- TestJWTGeneration
- TestPasswordHashing

// 角色系统测试
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
// API集成测试
- TestAuthEndpoints
- TestCharacterEndpoints
- TestGameEndpoints

// 数据库集成测试
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
*最后更新: 2025-01-26* 