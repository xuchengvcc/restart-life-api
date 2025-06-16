# 《重启人生》后端技术设计文档 (Backend TDD)

## 文档信息
- **版本**: v1.0
- **创建日期**: 2025-01-26
- **最后更新**: 2025-01-26
- **作者**: 后端技术团队

## 1. 多平台后端架构概览

### 1.1 多平台支持技术栈
- **编程语言**: Go 1.21+
- **Web框架**: Gin框架 (高性能HTTP服务)
- **数据库**: PostgreSQL 15+ (ACID事务支持)
- **缓存**: Redis 7+ (会话存储和高速缓存)
- **API网关**: Nginx/Kong (跨域、限流、负载均衡)
- **配置管理**: Viper (多环境配置)
- **日志**: Logrus/Zap (结构化日志)
- **监控**: Prometheus + Grafana (性能监控)
- **认证**: JWT + OAuth2.0 (多平台统一认证)

### 1.2 多平台架构设计原则
- **平台无关性**: API设计不依赖特定客户端平台
- **标准化协议**: 严格遵循RESTful API标准
- **统一认证**: 支持多种认证方式的统一认证体系
- **跨域支持**: 完善的CORS配置支持Web客户端
- **内容协商**: 支持多种数据格式(JSON/XML)
- **版本控制**: API版本管理支持客户端兼容性
- **错误标准化**: 统一的错误码和错误信息格式

### 1.3 多平台系统架构图
```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Unity Client  │    │                  │    │                 │
├─────────────────┤    │   API Gateway    │    │  Go API Server  │
│   Web Browser   │◄──►│  (Nginx/Kong)    │◄──►│  (Gin + HTTP)   │
├─────────────────┤    │  - CORS Support  │    │  - RESTful API  │
│  WeChat Mini    │    │  - Rate Limiting │    │  - Multi-Auth   │
│    Program      │    │  - Load Balance  │    │  - JSON API     │
├─────────────────┤    │  - SSL/TLS      │    │                 │
│   Mobile App    │    │                  │    │                 │
├─────────────────┤    └──────────────────┘    └─────────────────┘
│   Desktop App   │                                      │
└─────────────────┘                                      │
                                               ┌─────────────────┐
┌─────────────────┐                           │   PostgreSQL    │
│  Admin Panel    │◄──────────────────────────┤   Database      │
│  (Web Dashboard)│                           │   - JSONB       │
└─────────────────┘                           │   - Full Text   │
                                               └─────────────────┘
                                                         │
                                               ┌─────────────────┐
                                               │   Redis Cache   │
                                               │   - Sessions    │
                                               │   - API Cache   │
                                               │   - Rate Limit  │
                                               └─────────────────┘
```

### 1.4 支持的客户端平台

#### 1.4.1 游戏客户端
- **Unity客户端**: PC、Mac、Linux桌面游戏
- **Unity移动端**: iOS、Android原生应用
- **Web游戏**: WebGL构建的浏览器游戏

#### 1.4.2 Web应用
- **现代浏览器**: Chrome、Firefox、Safari、Edge
- **响应式设计**: 支持桌面和移动浏览器
- **PWA支持**: 渐进式Web应用特性

#### 1.4.3 小程序生态
- **微信小程序**: 微信生态内的轻应用
- **支付宝小程序**: 支付宝生态适配
- **抖音小程序**: 字节跳动生态扩展

#### 1.4.4 原生移动应用
- **React Native**: 跨平台移动应用开发
- **Flutter**: Google跨平台解决方案
- **原生开发**: iOS Swift、Android Kotlin

### 1.5 多平台特性支持

#### 1.5.1 认证系统适配
- **统一Token**: JWT标准令牌所有平台通用
- **OAuth2.0**: 支持第三方登录(微信、QQ、微博等)
- **小程序登录**: 微信小程序专用登录流程
- **游客模式**: 免注册游戏体验

#### 1.5.2 数据同步机制
- **云存档**: 跨设备数据同步
- **离线支持**: 客户端本地数据存储
- **冲突解决**: 数据冲突的自动和手动解决
- **增量同步**: 优化网络传输效率

#### 1.5.3 性能优化
- **CDN分发**: 静态资源全球加速
- **API缓存**: 高频接口Redis缓存
- **数据库优化**: 读写分离、连接池优化
- **压缩传输**: Gzip压缩减少传输量

### 1.6 项目结构
```
restart-life-api/
├── cmd/
│   └── server/
│       └── main.go              # 应用程序入口点
├── internal/
│   ├── api/
│   │   ├── handlers/            # HTTP处理器
│   │   │   ├── auth.go
│   │   │   ├── character.go
│   │   │   ├── game.go
│   │   │   └── achievement.go
│   │   ├── middleware/          # 中间件
│   │   │   ├── auth.go
│   │   │   ├── cors.go
│   │   │   └── ratelimit.go
│   │   └── routes/              # 路由定义
│   │       └── routes.go
│   ├── models/                  # 数据模型
│   │   ├── user.go
│   │   ├── character.go
│   │   ├── event.go
│   │   └── achievement.go
│   ├── services/                # 业务逻辑服务
│   │   ├── character_service.go
│   │   ├── game_service.go
│   │   └── event_service.go
│   ├── repository/              # 数据访问层
│   │   ├── postgres/
│   │   │   ├── character_repo.go
│   │   │   ├── event_repo.go
│   │   │   └── user_repo.go
│   │   └── redis/
│   │       └── cache_repo.go
│   └── config/                  # 配置管理
│       └── config.go
├── pkg/
│   ├── utils/                   # 工具函数
│   │   ├── jwt.go
│   │   ├── password.go
│   │   └── validator.go
│   └── database/                # 数据库连接
│       ├── postgres.go
│       └── redis.go
├── migrations/                  # 数据库迁移文件
├── docs/                        # API文档
├── docker/
│   └── Dockerfile
├── go.mod
├── go.sum
└── README.md
```

## 2. 数据库设计

### 2.1 数据库架构设计原则

基于游戏数据库设计最佳实践，我们采用以下设计原则：

1. **数据规范化**: 遵循第三范式，减少数据冗余
2. **性能优化**: 合理使用索引和分区
3. **扩展性**: 支持水平扩展和数据分片
4. **数据完整性**: 使用外键约束和触发器

### 2.2 核心数据表设计

#### 2.2.1 用户和角色表

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
    birth_year INTEGER NOT NULL CHECK (birth_year BETWEEN 1800 AND 2050),
    current_age INTEGER NOT NULL DEFAULT 0,
    gender VARCHAR(20) NOT NULL,
    race VARCHAR(50) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 角色属性表
CREATE TABLE character_attributes (
    attribute_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    character_id UUID NOT NULL REFERENCES characters(character_id) ON DELETE CASCADE,
    intelligence INTEGER NOT NULL DEFAULT 50 CHECK (intelligence BETWEEN 0 AND 100),
    constitution INTEGER NOT NULL DEFAULT 50 CHECK (constitution BETWEEN 0 AND 100),
    charisma INTEGER NOT NULL DEFAULT 50 CHECK (charisma BETWEEN 0 AND 100),
    willpower INTEGER NOT NULL DEFAULT 50 CHECK (willpower BETWEEN 0 AND 100),
    creativity INTEGER NOT NULL DEFAULT 50 CHECK (creativity BETWEEN 0 AND 100),
    academic_skill INTEGER NOT NULL DEFAULT 0 CHECK (academic_skill BETWEEN 0 AND 100),
    social_skill INTEGER NOT NULL DEFAULT 0 CHECK (social_skill BETWEEN 0 AND 100),
    athletic_skill INTEGER NOT NULL DEFAULT 0 CHECK (athletic_skill BETWEEN 0 AND 100),
    artistic_skill INTEGER NOT NULL DEFAULT 0 CHECK (artistic_skill BETWEEN 0 AND 100),
    business_skill INTEGER NOT NULL DEFAULT 0 CHECK (business_skill BETWEEN 0 AND 100),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

#### 2.2.2 事件系统表

```sql
-- 事件模板表
CREATE TABLE event_templates (
    template_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_name VARCHAR(200) NOT NULL,
    event_type VARCHAR(50) NOT NULL, -- 'random', 'choice', 'development', 'relationship', 'era'
    description TEXT NOT NULL,
    min_age INTEGER,
    max_age INTEGER,
    required_attributes JSONB, -- 触发条件
    probability_weight DECIMAL(5,3) DEFAULT 1.0,
    era_start INTEGER, -- 时代限制
    era_end INTEGER,
    countries TEXT[], -- 国家限制
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 事件选择选项表
CREATE TABLE event_choices (
    choice_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    template_id UUID NOT NULL REFERENCES event_templates(template_id) ON DELETE CASCADE,
    choice_text TEXT NOT NULL,
    choice_order INTEGER NOT NULL,
    requirements JSONB, -- 选择条件
    effects JSONB NOT NULL, -- 选择结果
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 角色事件历史表
CREATE TABLE character_events (
    event_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    character_id UUID NOT NULL REFERENCES characters(character_id) ON DELETE CASCADE,
    template_id UUID NOT NULL REFERENCES event_templates(template_id),
    event_age INTEGER NOT NULL,
    chosen_option_id UUID REFERENCES event_choices(choice_id),
    event_result JSONB NOT NULL, -- 事件结果详情
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### 2.3 Go数据模型定义

```go
// internal/models/user.go
package models

import (
    "time"
    "github.com/google/uuid"
)

// User 用户模型
type User struct {
    ID           uuid.UUID  `json:"id" db:"user_id"`
    Username     string     `json:"username" db:"username"`
    Email        string     `json:"email" db:"email"`
    PasswordHash string     `json:"-" db:"password_hash"`
    CreatedAt    time.Time  `json:"created_at" db:"created_at"`
    UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
    LastLogin    *time.Time `json:"last_login" db:"last_login"`
    IsActive     bool       `json:"is_active" db:"is_active"`
}

// Character 角色模型
type Character struct {
    ID            uuid.UUID              `json:"id" db:"character_id"`
    UserID        uuid.UUID              `json:"user_id" db:"user_id"`
    Name          string                 `json:"character_name" db:"character_name"`
    BirthCountry  string                 `json:"birth_country" db:"birth_country"`
    BirthYear     int                    `json:"birth_year" db:"birth_year"`
    CurrentAge    int                    `json:"current_age" db:"current_age"`
    Gender        string                 `json:"gender" db:"gender"`
    Race          string                 `json:"race" db:"race"`
    IsActive      bool                   `json:"is_active" db:"is_active"`
    CreatedAt     time.Time              `json:"created_at" db:"created_at"`
    UpdatedAt     time.Time              `json:"updated_at" db:"updated_at"`
    
    // 关联数据
    Attributes    *CharacterAttributes   `json:"attributes,omitempty"`
    Health        *CharacterHealth       `json:"health,omitempty"`
    Economics     *CharacterEconomics    `json:"economics,omitempty"`
    Relationships []CharacterRelationship `json:"relationships,omitempty"`
}

// CharacterAttributes 角色属性
type CharacterAttributes struct {
    AttributeID    uuid.UUID `json:"attribute_id" db:"attribute_id"`
    CharacterID    uuid.UUID `json:"character_id" db:"character_id"`
    Intelligence   int       `json:"intelligence" db:"intelligence"`
    Constitution   int       `json:"constitution" db:"constitution"`
    Charisma       int       `json:"charisma" db:"charisma"`
    Willpower      int       `json:"willpower" db:"willpower"`
    Creativity     int       `json:"creativity" db:"creativity"`
    AcademicSkill  int       `json:"academic_skill" db:"academic_skill"`
    SocialSkill    int       `json:"social_skill" db:"social_skill"`
    AthleticSkill  int       `json:"athletic_skill" db:"athletic_skill"`
    ArtisticSkill  int       `json:"artistic_skill" db:"artistic_skill"`
    BusinessSkill  int       `json:"business_skill" db:"business_skill"`
    UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}
```

## 3. 多平台API设计

### 3.1 RESTful API标准

#### 3.1.1 API设计原则
- **资源导向**: 使用名词而非动词定义端点
- **HTTP方法**: 正确使用GET、POST、PUT、DELETE、PATCH
- **状态码**: 使用标准HTTP状态码
- **幂等性**: 确保相同操作的可重复性
- **无状态**: 每个请求包含完整的处理信息

#### 3.1.2 URL结构规范
```
基础格式: https://api.restartlife.com/v1/{resource}
示例:
- GET    /api/v1/users                 # 获取用户列表
- POST   /api/v1/users                 # 创建用户
- GET    /api/v1/users/{id}            # 获取特定用户
- PUT    /api/v1/users/{id}            # 更新用户
- DELETE /api/v1/users/{id}            # 删除用户
- GET    /api/v1/characters/{id}/events # 获取角色事件
```

#### 3.1.3 HTTP状态码使用
- **200 OK**: 请求成功
- **201 Created**: 资源创建成功
- **204 No Content**: 成功但无返回内容
- **400 Bad Request**: 请求参数错误
- **401 Unauthorized**: 未认证
- **403 Forbidden**: 无权限
- **404 Not Found**: 资源不存在
- **409 Conflict**: 资源冲突
- **422 Unprocessable Entity**: 数据验证失败
- **429 Too Many Requests**: 请求频率超限
- **500 Internal Server Error**: 服务器内部错误

### 3.2 多平台认证系统

#### 3.2.1 JWT统一认证

```go
// JWT配置结构
type JWTConfig struct {
    SecretKey      string        `json:"secret_key"`
    AccessTTL      time.Duration `json:"access_ttl"`
    RefreshTTL     time.Duration `json:"refresh_ttl"`
    Issuer         string        `json:"issuer"`
    Audience       []string      `json:"audience"`
}

// JWT Claims结构
type Claims struct {
    UserID       string   `json:"user_id"`
    Username     string   `json:"username"`
    Platform     string   `json:"platform"`
    Permissions  []string `json:"permissions"`
    jwt.RegisteredClaims
}

// 生成多平台适配的JWT Token
func GenerateToken(userID, username, platform string) (string, error) {
    claims := Claims{
        UserID:   userID,
        Username: username,
        Platform: platform,
        Permissions: getPermissionsByPlatform(platform),
        RegisteredClaims: jwt.RegisteredClaims{
            Issuer:    "restart-life-api",
            Subject:   userID,
            Audience:  []string{"unity", "web", "wechat", "mobile"},
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            NotBefore: jwt.NewNumericDate(time.Now()),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(jwtSecret))
}
```

#### 3.2.2 多平台登录支持

```go
// 统一登录请求结构
type UniversalLoginRequest struct {
    Platform     string `json:"platform" validate:"required,oneof=unity web wechat mobile"`
    LoginType    string `json:"login_type" validate:"required"`
    Credentials  map[string]interface{} `json:"credentials"`
    DeviceInfo   DeviceInfo `json:"device_info,omitempty"`
}

// 设备信息
type DeviceInfo struct {
    DeviceID   string `json:"device_id"`
    Platform   string `json:"platform"`
    Version    string `json:"version"`
    UserAgent  string `json:"user_agent,omitempty"`
}

// 统一登录响应
type UniversalLoginResponse struct {
    AccessToken  string    `json:"access_token"`
    RefreshToken string    `json:"refresh_token"`
    TokenType    string    `json:"token_type"`
    ExpiresIn    int       `json:"expires_in"`
    User         UserInfo  `json:"user"`
    Permissions  []string  `json:"permissions"`
}

// 多平台登录处理器
func (h *AuthHandler) UniversalLogin(c *gin.Context) {
    var req UniversalLoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse{
            Code:    "INVALID_REQUEST",
            Message: "请求数据格式错误",
            Details: err.Error(),
        })
        return
    }

    var user *models.User
    var err error

    switch req.LoginType {
    case "username_password":
        user, err = h.authService.LoginWithPassword(
            req.Credentials["username"].(string),
            req.Credentials["password"].(string),
        )
    case "wechat_code":
        user, err = h.authService.LoginWithWeChatCode(
            req.Credentials["code"].(string),
        )
    case "guest":
        user, err = h.authService.CreateGuestUser(req.DeviceInfo)
    case "oauth":
        user, err = h.authService.LoginWithOAuth(
            req.Credentials["provider"].(string),
            req.Credentials["token"].(string),
        )
    default:
        c.JSON(http.StatusBadRequest, ErrorResponse{
            Code:    "UNSUPPORTED_LOGIN_TYPE",
            Message: "不支持的登录方式",
        })
        return
    }

    if err != nil {
        c.JSON(http.StatusUnauthorized, ErrorResponse{
            Code:    "LOGIN_FAILED",
            Message: "登录失败",
            Details: err.Error(),
        })
        return
    }

    // 生成Token
    accessToken, refreshToken, err := h.authService.GenerateTokens(
        user.ID.String(), user.Username, req.Platform,
    )
    if err != nil {
        c.JSON(http.StatusInternalServerError, ErrorResponse{
            Code:    "TOKEN_GENERATION_FAILED",
            Message: "令牌生成失败",
        })
        return
    }

    response := UniversalLoginResponse{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        TokenType:    "Bearer",
        ExpiresIn:    3600,
        User: UserInfo{
            ID:       user.ID.String(),
            Username: user.Username,
            Email:    user.Email,
        },
        Permissions: getPermissionsByPlatform(req.Platform),
    }

    c.JSON(http.StatusOK, response)
}
```

### 3.3 跨域支持(CORS)

#### 3.3.1 CORS中间件配置

```go
// CORS配置结构
type CORSConfig struct {
    AllowOrigins     []string `json:"allow_origins"`
    AllowMethods     []string `json:"allow_methods"`
    AllowHeaders     []string `json:"allow_headers"`
    ExposeHeaders    []string `json:"expose_headers"`
    AllowCredentials bool     `json:"allow_credentials"`
    MaxAge           int      `json:"max_age"`
}

// CORS中间件实现
func CORSMiddleware(config CORSConfig) gin.HandlerFunc {
    return func(c *gin.Context) {
        origin := c.Request.Header.Get("Origin")
        
        // 检查Origin是否被允许
        if isOriginAllowed(origin, config.AllowOrigins) {
            c.Header("Access-Control-Allow-Origin", origin)
        }
        
        c.Header("Access-Control-Allow-Methods", strings.Join(config.AllowMethods, ", "))
        c.Header("Access-Control-Allow-Headers", strings.Join(config.AllowHeaders, ", "))
        c.Header("Access-Control-Expose-Headers", strings.Join(config.ExposeHeaders, ", "))
        
        if config.AllowCredentials {
            c.Header("Access-Control-Allow-Credentials", "true")
        }
        
        if config.MaxAge > 0 {
            c.Header("Access-Control-Max-Age", strconv.Itoa(config.MaxAge))
        }
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(http.StatusNoContent)
            return
        }
        
        c.Next()
    }
}

// 默认CORS配置
func DefaultCORSConfig() CORSConfig {
    return CORSConfig{
        AllowOrigins: []string{
            "http://localhost:3000",    // React开发服务器
            "http://localhost:8080",    // Vue开发服务器
            "https://restartlife.com",  // 生产环境
            "https://*.restartlife.com", // 子域名
        },
        AllowMethods: []string{
            "GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH",
        },
        AllowHeaders: []string{
            "Origin", "Content-Type", "Accept", "Authorization",
            "X-Requested-With", "X-Platform", "X-Version",
        },
        ExposeHeaders: []string{
            "X-Total-Count", "X-Page-Count", "X-Rate-Limit",
        },
        AllowCredentials: true,
        MaxAge:           86400, // 24小时
    }
}
```

### 3.4 平台特定适配

#### 3.4.1 微信小程序适配

```go
// 微信小程序登录
type WeChatMiniProgramLoginRequest struct {
    Code          string `json:"code" validate:"required"`
    EncryptedData string `json:"encrypted_data,omitempty"`
    IV            string `json:"iv,omitempty"`
}

// 微信小程序登录处理
func (h *AuthHandler) WeChatMiniProgramLogin(c *gin.Context) {
    var req WeChatMiniProgramLoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse{
            Code:    "INVALID_REQUEST",
            Message: "请求数据格式错误",
        })
        return
    }
    
    // 调用微信API获取session_key和openid
    sessionData, err := h.wechatService.Code2Session(req.Code)
    if err != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse{
            Code:    "WECHAT_AUTH_FAILED",
            Message: "微信授权失败",
            Details: err.Error(),
        })
        return
    }
    
    // 查找或创建用户
    user, err := h.authService.FindOrCreateWeChatUser(sessionData.OpenID, sessionData.UnionID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, ErrorResponse{
            Code:    "USER_CREATE_FAILED",
            Message: "用户创建失败",
        })
        return
    }
    
    // 生成JWT令牌
    token, err := h.authService.GenerateToken(user.ID.String(), user.Username, "wechat")
    if err != nil {
        c.JSON(http.StatusInternalServerError, ErrorResponse{
            Code:    "TOKEN_GENERATION_FAILED",
            Message: "令牌生成失败",
        })
        return
    }
    
    response := UniversalLoginResponse{
        AccessToken: token,
        TokenType:   "Bearer",
        ExpiresIn:   3600,
        User: UserInfo{
            ID:       user.ID.String(),
            Username: user.Username,
            OpenID:   sessionData.OpenID,
        },
    }
    
    c.JSON(http.StatusOK, response)
}
```

#### 3.4.2 Unity客户端优化

```go
// Unity专用数据格式
type UnityGameData struct {
    Character      *Character               `json:"character"`
    GameState      *GameState               `json:"game_state"`
    Events         []Event                  `json:"events"`
    Achievements   []Achievement            `json:"achievements"`
    Settings       *ClientSettings          `json:"settings"`
    
    // Unity特定字段
    AssetVersion   string                   `json:"asset_version"`
    ClientVersion  string                   `json:"client_version"`
    Platform       string                   `json:"platform"`
}

// Unity批量数据同步
func (h *GameHandler) UnityDataSync(c *gin.Context) {
    userID := c.GetString("user_id")
    characterID := c.Param("character_id")
    
    // 获取完整游戏数据
    gameData, err := h.gameService.GetCompleteGameData(userID, characterID)
    if err != nil {
        c.JSON(http.StatusNotFound, ErrorResponse{
            Code:    "GAME_DATA_NOT_FOUND",
            Message: "游戏数据不存在",
        })
        return
    }
    
    // Unity专用格式化
    unityData := UnityGameData{
        Character:      gameData.Character,
        GameState:      gameData.GameState,
        Events:         gameData.RecentEvents,
        Achievements:   gameData.Achievements,
        Settings:       gameData.Settings,
        AssetVersion:   getCurrentAssetVersion(),
        ClientVersion:  getMinimumClientVersion(),
        Platform:       "unity",
    }
    
    c.JSON(http.StatusOK, unityData)
}
```

### 3.5 统一错误处理

```go
// 标准错误响应结构
type ErrorResponse struct {
    Code      string                 `json:"code"`
    Message   string                 `json:"message"`
    Details   interface{}            `json:"details,omitempty"`
    TraceID   string                 `json:"trace_id"`
    Timestamp time.Time              `json:"timestamp"`
    Path      string                 `json:"path"`
    Method    string                 `json:"method"`
    Platform  string                 `json:"platform,omitempty"`
}

// 错误码定义
const (
    ErrCodeInvalidRequest      = "INVALID_REQUEST"
    ErrCodeUnauthorized        = "UNAUTHORIZED"
    ErrCodeForbidden          = "FORBIDDEN"
    ErrCodeNotFound           = "NOT_FOUND"
    ErrCodeConflict           = "CONFLICT"
    ErrCodeValidationFailed   = "VALIDATION_FAILED"
    ErrCodeRateLimited        = "RATE_LIMITED"
    ErrCodeInternalError      = "INTERNAL_ERROR"
    ErrCodeServiceUnavailable = "SERVICE_UNAVAILABLE"
)

// 全局错误处理中间件
func ErrorHandlingMiddleware() gin.HandlerFunc {
    return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
        err := recovered.(error)
        
        traceID := c.GetString("trace_id")
        platform := c.GetHeader("X-Platform")
        
        errorResponse := ErrorResponse{
            Code:      ErrCodeInternalError,
            Message:   "服务器内部错误",
            TraceID:   traceID,
            Timestamp: time.Now(),
            Path:      c.Request.URL.Path,
            Method:    c.Request.Method,
            Platform:  platform,
        }
        
        // 根据错误类型设置不同的状态码和消息
        switch e := err.(type) {
        case *ValidationError:
            errorResponse.Code = ErrCodeValidationFailed
            errorResponse.Message = "数据验证失败"
            errorResponse.Details = e.Fields
            c.JSON(http.StatusUnprocessableEntity, errorResponse)
        case *AuthError:
            errorResponse.Code = ErrCodeUnauthorized
            errorResponse.Message = e.Message
            c.JSON(http.StatusUnauthorized, errorResponse)
        case *NotFoundError:
            errorResponse.Code = ErrCodeNotFound
            errorResponse.Message = e.Message
            c.JSON(http.StatusNotFound, errorResponse)
        default:
            logger.Error("Unhandled error", "error", err, "trace_id", traceID)
            c.JSON(http.StatusInternalServerError, errorResponse)
        }
    })
}
```

### 3.6 API版本控制

```go
// 版本控制中间件
func APIVersionMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        version := c.GetHeader("API-Version")
        if version == "" {
            version = c.Query("version")
        }
        if version == "" {
            version = "v1" // 默认版本
        }
        
        c.Set("api_version", version)
        c.Header("API-Version", version)
        c.Next()
    }
}

// 版本路由注册
func RegisterVersionedRoutes(r *gin.Engine) {
    v1 := r.Group("/api/v1")
    {
        v1.Use(APIVersionMiddleware())
        registerV1Routes(v1)
    }
    
    v2 := r.Group("/api/v2")
    {
        v2.Use(APIVersionMiddleware())
        registerV2Routes(v2)
    }
}
```

### 3.7 API限流与安全

```go
// 限流配置
type RateLimitConfig struct {
    RequestsPerMinute int    `json:"requests_per_minute"`
    BurstSize         int    `json:"burst_size"`
    KeyPrefix         string `json:"key_prefix"`
}

// 基于Redis的限流中间件
func RateLimitMiddleware(config RateLimitConfig, redis *redis.Client) gin.HandlerFunc {
    return func(c *gin.Context) {
        key := generateRateLimitKey(c, config.KeyPrefix)
        
        current, err := redis.Incr(context.Background(), key).Result()
        if err != nil {
            logger.Error("Rate limit check failed", "error", err)
            c.Next()
            return
        }
        
        if current == 1 {
            redis.Expire(context.Background(), key, time.Minute)
        }
        
        if current > int64(config.RequestsPerMinute) {
            c.JSON(http.StatusTooManyRequests, ErrorResponse{
                Code:    "RATE_LIMITED",
                Message: "请求频率过高，请稍后再试",
            })
            c.Abort()
            return
        }
        
        // 设置响应头
        c.Header("X-RateLimit-Limit", strconv.Itoa(config.RequestsPerMinute))
        c.Header("X-RateLimit-Remaining", strconv.FormatInt(int64(config.RequestsPerMinute)-current, 10))
        c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(time.Minute).Unix(), 10))
        
        c.Next()
    }
}

// 生成限流键
func generateRateLimitKey(c *gin.Context, prefix string) string {
    userID := c.GetString("user_id")
    if userID != "" {
        return fmt.Sprintf("%s:user:%s", prefix, userID)
    }
    
    clientIP := c.ClientIP()
    return fmt.Sprintf("%s:ip:%s", prefix, clientIP)
}
```

### 3.8 中间件组合配置

```go
// 多平台中间件设置
func SetupMultiPlatformMiddleware(r *gin.Engine, config *Config) {
    // 基础中间件
    r.Use(gin.Logger())
    r.Use(gin.Recovery())
    
    // 安全中间件
    r.Use(SecurityHeadersMiddleware())
    
    // CORS支持
    r.Use(CORSMiddleware(DefaultCORSConfig()))
    
    // 请求ID和追踪
    r.Use(TraceIDMiddleware())
    
    // 错误处理
    r.Use(ErrorHandlingMiddleware())
    
    // 限流
    r.Use(RateLimitMiddleware(config.RateLimit, config.Redis))
    
    // 内容验证
    r.Use(ValidationMiddleware())
    
    // API版本控制
    r.Use(APIVersionMiddleware())
    
    // 平台检测
    r.Use(PlatformDetectionMiddleware())
}

// 平台检测中间件
func PlatformDetectionMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        platform := c.GetHeader("X-Platform")
        if platform == "" {
            userAgent := c.GetHeader("User-Agent")
            platform = detectPlatformFromUserAgent(userAgent)
        }
        
        c.Set("platform", platform)
        c.Next()
    }
}

func detectPlatformFromUserAgent(userAgent string) string {
    switch {
    case strings.Contains(userAgent, "Unity"):
        return "unity"
    case strings.Contains(userAgent, "MicroMessenger"):
        return "wechat"
    case strings.Contains(userAgent, "Mozilla"):
        return "web"
    default:
        return "unknown"
    }
}
            req.Credentials["username"].(string),
            req.Credentials["password"].(string),
        )
    case "wechat_code":
        user, err = h.authService.LoginWithWeChatCode(
            req.Credentials["code"].(string),
        )
    case "guest":
        user, err = h.authService.CreateGuestUser(req.DeviceInfo)
    case "oauth":
        user, err = h.authService.LoginWithOAuth(
            req.Credentials["provider"].(string),
            req.Credentials["token"].(string),
        )
    default:
        c.JSON(http.StatusBadRequest, ErrorResponse{
            Code:    "UNSUPPORTED_LOGIN_TYPE",
            Message: "不支持的登录方式",
        })
        return
    }

    if err != nil {
        c.JSON(http.StatusUnauthorized, ErrorResponse{
            Code:    "LOGIN_FAILED",
            Message: "登录失败",
            Details: err.Error(),
        })
        return
    }

    // 生成Token
    accessToken, refreshToken, err := h.authService.GenerateTokens(
        user.ID.String(), user.Username, req.Platform,
    )
    if err != nil {
        c.JSON(http.StatusInternalServerError, ErrorResponse{
            Code:    "TOKEN_GENERATION_FAILED",
            Message: "令牌生成失败",
        })
        return
    }

    response := UniversalLoginResponse{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        TokenType:    "Bearer",
        ExpiresIn:    3600,
        User: UserInfo{
            ID:       user.ID.String(),
            Username: user.Username,
            Email:    user.Email,
        },
        Permissions: getPermissionsByPlatform(req.Platform),
    }

    c.JSON(http.StatusOK, response)
}
```

### 3.3 跨域支持(CORS)

#### 3.3.1 CORS中间件配置

```go
// CORS配置结构
type CORSConfig struct {
    AllowOrigins     []string `json:"allow_origins"`
    AllowMethods     []string `json:"allow_methods"`
    AllowHeaders     []string `json:"allow_headers"`
    ExposeHeaders    []string `json:"expose_headers"`
    AllowCredentials bool     `json:"allow_credentials"`
    MaxAge           int      `json:"max_age"`
}

// CORS中间件实现
func CORSMiddleware(config CORSConfig) gin.HandlerFunc {
    return func(c *gin.Context) {
        origin := c.Request.Header.Get("Origin")
        
        // 检查Origin是否被允许
        if isOriginAllowed(origin, config.AllowOrigins) {
            c.Header("Access-Control-Allow-Origin", origin)
        }
        
        c.Header("Access-Control-Allow-Methods", strings.Join(config.AllowMethods, ", "))
        c.Header("Access-Control-Allow-Headers", strings.Join(config.AllowHeaders, ", "))
        c.Header("Access-Control-Expose-Headers", strings.Join(config.ExposeHeaders, ", "))
        
        if config.AllowCredentials {
            c.Header("Access-Control-Allow-Credentials", "true")
        }
        
        if config.MaxAge > 0 {
            c.Header("Access-Control-Max-Age", strconv.Itoa(config.MaxAge))
        }
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(http.StatusNoContent)
            return
        }
        
        c.Next()
    }
}

// 默认CORS配置
func DefaultCORSConfig() CORSConfig {
    return CORSConfig{
        AllowOrigins: []string{
            "http://localhost:3000",    // React开发服务器
            "http://localhost:8080",    // Vue开发服务器
            "https://restartlife.com",  // 生产环境
            "https://*.restartlife.com", // 子域名
        },
        AllowMethods: []string{
            "GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH",
        },
        AllowHeaders: []string{
            "Origin", "Content-Type", "Accept", "Authorization",
            "X-Requested-With", "X-Platform", "X-Version",
        },
        ExposeHeaders: []string{
            "X-Total-Count", "X-Page-Count", "X-Rate-Limit",
        },
        AllowCredentials: true,
        MaxAge:           86400, // 24小时
    }
}
```

### 3.4 内容协商与版本控制

#### 3.4.1 API版本控制

```go
// 版本控制中间件
func APIVersionMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        version := c.GetHeader("API-Version")
        if version == "" {
            version = c.Query("version")
        }
        if version == "" {
            version = "v1" // 默认版本
        }
        
        c.Set("api_version", version)
        c.Header("API-Version", version)
        c.Next()
    }
}

// 版本路由注册
func RegisterVersionedRoutes(r *gin.Engine) {
    v1 := r.Group("/api/v1")
    {
        v1.Use(APIVersionMiddleware())
        registerV1Routes(v1)
    }
    
    v2 := r.Group("/api/v2")
    {
        v2.Use(APIVersionMiddleware())
        registerV2Routes(v2)
    }
}
```

#### 3.4.2 内容协商

```go
// 内容协商中间件
func ContentNegotiationMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        accept := c.GetHeader("Accept")
        contentType := c.GetHeader("Content-Type")
        
        // 支持的内容类型
        supportedTypes := []string{
            "application/json",
            "application/xml",
            "text/plain",
        }
        
        // 确定响应格式
        responseFormat := negotiateContent(accept, supportedTypes)
        c.Set("response_format", responseFormat)
        
        // 验证请求格式
        if !isContentTypeSupported(contentType, supportedTypes) {
            c.JSON(http.StatusUnsupportedMediaType, ErrorResponse{
                Code:    "UNSUPPORTED_MEDIA_TYPE",
                Message: "不支持的媒体类型",
            })
            c.Abort()
            return
        }
        
        c.Next()
    }
}

// 响应格式化函数
func FormatResponse(c *gin.Context, data interface{}) {
    format := c.GetString("response_format")
    
    switch format {
    case "application/xml":
        c.XML(http.StatusOK, data)
    case "text/plain":
        c.String(http.StatusOK, fmt.Sprintf("%v", data))
    default:
        c.JSON(http.StatusOK, data)
    }
}
```

### 3.5 统一错误处理

#### 3.5.1 标准错误响应格式

```go
// 标准错误响应结构
type ErrorResponse struct {
    Code      string                 `json:"code"`
    Message   string                 `json:"message"`
    Details   interface{}            `json:"details,omitempty"`
    TraceID   string                 `json:"trace_id"`
    Timestamp time.Time              `json:"timestamp"`
    Path      string                 `json:"path"`
    Method    string                 `json:"method"`
    Platform  string                 `json:"platform,omitempty"`
}

// 错误码定义
const (
    ErrCodeInvalidRequest      = "INVALID_REQUEST"
    ErrCodeUnauthorized        = "UNAUTHORIZED"
    ErrCodeForbidden          = "FORBIDDEN"
    ErrCodeNotFound           = "NOT_FOUND"
    ErrCodeConflict           = "CONFLICT"
    ErrCodeValidationFailed   = "VALIDATION_FAILED"
    ErrCodeRateLimited        = "RATE_LIMITED"
    ErrCodeInternalError      = "INTERNAL_ERROR"
    ErrCodeServiceUnavailable = "SERVICE_UNAVAILABLE"
)

// 全局错误处理中间件
func ErrorHandlingMiddleware() gin.HandlerFunc {
    return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
        err := recovered.(error)
        
        traceID := c.GetString("trace_id")
        platform := c.GetHeader("X-Platform")
        
        errorResponse := ErrorResponse{
            Code:      ErrCodeInternalError,
            Message:   "服务器内部错误",
            TraceID:   traceID,
            Timestamp: time.Now(),
            Path:      c.Request.URL.Path,
            Method:    c.Request.Method,
            Platform:  platform,
        }
        
        // 根据错误类型设置不同的状态码和消息
        switch e := err.(type) {
        case *ValidationError:
            errorResponse.Code = ErrCodeValidationFailed
            errorResponse.Message = "数据验证失败"
            errorResponse.Details = e.Fields
            c.JSON(http.StatusUnprocessableEntity, errorResponse)
        case *AuthError:
            errorResponse.Code = ErrCodeUnauthorized
            errorResponse.Message = e.Message
            c.JSON(http.StatusUnauthorized, errorResponse)
        case *NotFoundError:
            errorResponse.Code = ErrCodeNotFound
            errorResponse.Message = e.Message
            c.JSON(http.StatusNotFound, errorResponse)
        default:
            logger.Error("Unhandled error", "error", err, "trace_id", traceID)
            c.JSON(http.StatusInternalServerError, errorResponse)
        }
    })
}
```

### 3.6 平台特定适配

#### 3.6.1 微信小程序适配

```go
// 微信小程序登录
type WeChatMiniProgramLoginRequest struct {
    Code          string `json:"code" validate:"required"`
    EncryptedData string `json:"encrypted_data,omitempty"`
    IV            string `json:"iv,omitempty"`
}

// 微信小程序登录处理
func (h *AuthHandler) WeChatMiniProgramLogin(c *gin.Context) {
    var req WeChatMiniProgramLoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse{
            Code:    "INVALID_REQUEST",
            Message: "请求数据格式错误",
        })
        return
    }
    
    // 调用微信API获取session_key和openid
    sessionData, err := h.wechatService.Code2Session(req.Code)
    if err != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse{
            Code:    "WECHAT_AUTH_FAILED",
            Message: "微信授权失败",
            Details: err.Error(),
        })
        return
    }
    
    // 查找或创建用户
    user, err := h.authService.FindOrCreateWeChatUser(sessionData.OpenID, sessionData.UnionID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, ErrorResponse{
            Code:    "USER_CREATE_FAILED",
            Message: "用户创建失败",
        })
        return
    }
    
    // 生成JWT令牌
    token, err := h.authService.GenerateToken(user.ID.String(), user.Username, "wechat")
    if err != nil {
        c.JSON(http.StatusInternalServerError, ErrorResponse{
            Code:    "TOKEN_GENERATION_FAILED",
            Message: "令牌生成失败",
        })
        return
    }
    
    response := UniversalLoginResponse{
        AccessToken: token,
        TokenType:   "Bearer",
        ExpiresIn:   3600,
        User: UserInfo{
            ID:       user.ID.String(),
            Username: user.Username,
            OpenID:   sessionData.OpenID,
        },
    }
    
    c.JSON(http.StatusOK, response)
}
```

#### 3.6.2 Unity客户端优化

```go
// Unity专用数据格式
type UnityGameData struct {
    Character      *Character               `json:"character"`
    GameState      *GameState               `json:"game_state"`
    Events         []Event                  `json:"events"`
    Achievements   []Achievement            `json:"achievements"`
    Settings       *ClientSettings          `json:"settings"`
    
    // Unity特定字段
    AssetVersion   string                   `json:"asset_version"`
    ClientVersion  string                   `json:"client_version"`
    Platform       string                   `json:"platform"`
}

// Unity批量数据同步
func (h *GameHandler) UnityDataSync(c *gin.Context) {
    userID := c.GetString("user_id")
    characterID := c.Param("character_id")
    
    // 获取完整游戏数据
    gameData, err := h.gameService.GetCompleteGameData(userID, characterID)
    if err != nil {
        c.JSON(http.StatusNotFound, ErrorResponse{
            Code:    "GAME_DATA_NOT_FOUND",
            Message: "游戏数据不存在",
        })
        return
    }
    
    // Unity专用格式化
    unityData := UnityGameData{
        Character:      gameData.Character,
        GameState:      gameData.GameState,
        Events:         gameData.RecentEvents,
        Achievements:   gameData.Achievements,
        Settings:       gameData.Settings,
        AssetVersion:   getCurrentAssetVersion(),
        ClientVersion:  getMinimumClientVersion(),
        Platform:       "unity",
    }
    
    c.JSON(http.StatusOK, unityData)
}
```

### 3.7 API限流与安全

#### 3.7.1 请求限流

```go
// 限流配置
type RateLimitConfig struct {
    RequestsPerMinute int    `json:"requests_per_minute"`
    BurstSize         int    `json:"burst_size"`
    KeyPrefix         string `json:"key_prefix"`
}

// 基于Redis的限流中间件
func RateLimitMiddleware(config RateLimitConfig, redis *redis.Client) gin.HandlerFunc {
    return func(c *gin.Context) {
        key := generateRateLimitKey(c, config.KeyPrefix)
        
        current, err := redis.Incr(context.Background(), key).Result()
        if err != nil {
            logger.Error("Rate limit check failed", "error", err)
            c.Next()
            return
        }
        
        if current == 1 {
            redis.Expire(context.Background(), key, time.Minute)
        }
        
        if current > int64(config.RequestsPerMinute) {
            c.JSON(http.StatusTooManyRequests, ErrorResponse{
                Code:    "RATE_LIMITED",
                Message: "请求频率过高，请稍后再试",
            })
            c.Abort()
            return
        }
        
        // 设置响应头
        c.Header("X-RateLimit-Limit", strconv.Itoa(config.RequestsPerMinute))
        c.Header("X-RateLimit-Remaining", strconv.FormatInt(int64(config.RequestsPerMinute)-current, 10))
        c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(time.Minute).Unix(), 10))
        
        c.Next()
    }
}

// 生成限流键
func generateRateLimitKey(c *gin.Context, prefix string) string {
    userID := c.GetString("user_id")
    if userID != "" {
        return fmt.Sprintf("%s:user:%s", prefix, userID)
    }
    
    clientIP := c.ClientIP()
    return fmt.Sprintf("%s:ip:%s", prefix, clientIP)
}
```

#### 3.7.2 输入验证和安全

```go
// 输入验证中间件
func ValidationMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 检查Content-Length
        if c.Request.ContentLength > MaxRequestSize {
            c.JSON(http.StatusRequestEntityTooLarge, ErrorResponse{
                Code:    "REQUEST_TOO_LARGE",
                Message: "请求体过大",
            })
            c.Abort()
            return
        }
        
        // 验证User-Agent
        userAgent := c.GetHeader("User-Agent")
        if !isValidUserAgent(userAgent) {
            c.JSON(http.StatusBadRequest, ErrorResponse{
                Code:    "INVALID_USER_AGENT",
                Message: "无效的User-Agent",
            })
            c.Abort()
            return
        }
        
        c.Next()
    }
}

// 安全头中间件
func SecurityHeadersMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("X-Content-Type-Options", "nosniff")
        c.Header("X-Frame-Options", "DENY")
        c.Header("X-XSS-Protection", "1; mode=block")
        c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        c.Header("Content-Security-Policy", "default-src 'self'")
        c.Next()
    }
}
```

### 3.8 API文档化

#### 3.8.1 Swagger/OpenAPI集成

```go
// Swagger配置
func setupSwagger(r *gin.Engine) {
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// API文档注释示例
// @title 重启人生 API
// @version 1.0
// @description 多平台人生模拟游戏后端API
// @termsOfService https://restartlife.com/terms
// @contact.name API Support
// @contact.url https://restartlife.com/support
// @contact.email support@restartlife.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host api.restartlife.com
// @BasePath /api/v1
// @schemes https http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT Authorization header using the Bearer scheme.
```

## 4. 配置管理和部署

### 4.1 配置管理

```go
// internal/config/config.go
package config

import (
    "fmt"
    "restart-life-api/pkg/database"
    "time"
    
    "github.com/spf13/viper"
)

type Config struct {
    Server   ServerConfig            `mapstructure:"server"`
    Database database.PostgresConfig `mapstructure:"database"`
    Redis    database.RedisConfig    `mapstructure:"redis"`
    JWT      JWTConfig              `mapstructure:"jwt"`
    App      AppConfig              `mapstructure:"app"`
}

type ServerConfig struct {
    Port         int           `mapstructure:"port"`
    Host         string        `mapstructure:"host"`
    ReadTimeout  time.Duration `mapstructure:"read_timeout"`
    WriteTimeout time.Duration `mapstructure:"write_timeout"`
    IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

type JWTConfig struct {
    SecretKey   string        `mapstructure:"secret_key"`
    ExpiresIn   time.Duration `mapstructure:"expires_in"`
    RefreshIn   time.Duration `mapstructure:"refresh_in"`
}

type AppConfig struct {
    Name        string `mapstructure:"name"`
    Version     string `mapstructure:"version"`
    Environment string `mapstructure:"environment"`
    LogLevel    string `mapstructure:"log_level"`
}

func LoadConfig(configPath string) (*Config, error) {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(configPath)
    viper.AddConfigPath(".")
    
    // 设置环境变量前缀
    viper.SetEnvPrefix("RESTART_LIFE")
    viper.AutomaticEnv()
    
    // 设置默认值
    setDefaults()
    
    if err := viper.ReadInConfig(); err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }
    
    var config Config
    if err := viper.Unmarshal(&config); err != nil {
        return nil, fmt.Errorf("failed to unmarshal config: %w", err)
    }
    
    return &config, nil
}

func setDefaults() {
    // Server defaults
    viper.SetDefault("server.port", 8080)
    viper.SetDefault("server.host", "localhost")
    viper.SetDefault("server.read_timeout", "10s")
    viper.SetDefault("server.write_timeout", "10s")
    viper.SetDefault("server.idle_timeout", "60s")
    
    // Database defaults
    viper.SetDefault("database.host", "localhost")
    viper.SetDefault("database.port", 5432)
    viper.SetDefault("database.user", "postgres")
    viper.SetDefault("database.dbname", "restart_life_db")
    viper.SetDefault("database.sslmode", "disable")
    viper.SetDefault("database.max_open_conns", 25)
    viper.SetDefault("database.max_idle_conns", 5)
    viper.SetDefault("database.max_lifetime", "15m")
    
    // Redis defaults
    viper.SetDefault("redis.host", "localhost")
    viper.SetDefault("redis.port", 6379)
    viper.SetDefault("redis.db", 0)
    viper.SetDefault("redis.pool_size", 10)
    viper.SetDefault("redis.min_idle_conns", 2)
    
    // JWT defaults
    viper.SetDefault("jwt.expires_in", "1h")
    viper.SetDefault("jwt.refresh_in", "168h") // 7 days
    
    // App defaults
    viper.SetDefault("app.name", "Restart Life API")
    viper.SetDefault("app.version", "1.0.0")
    viper.SetDefault("app.environment", "development")
    viper.SetDefault("app.log_level", "info")
}
```

### 4.2 Docker部署配置

```dockerfile
# docker/Dockerfile
FROM golang:1.21-alpine AS builder

# 安装必要的工具
RUN apk add --no-cache git ca-certificates tzdata

# 设置工作目录
WORKDIR /app

# 复制go mod文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/server/main.go

# 最终镜像
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .
COPY --from=builder /app/config.yaml .

# 创建非root用户
RUN adduser -D -s /bin/sh appuser
USER appuser

EXPOSE 8080

CMD ["./main"]
```

```yaml
# docker-compose.yml
version: '3.8'

services:
  api:
    build:
      context: .
      dockerfile: docker/Dockerfile
    ports:
      - "8080:8080"
    environment:
      - RESTART_LIFE_DATABASE_HOST=postgres
      - RESTART_LIFE_DATABASE_PORT=5432
      - RESTART_LIFE_DATABASE_USER=postgres
      - RESTART_LIFE_DATABASE_PASSWORD=password
      - RESTART_LIFE_DATABASE_DBNAME=restart_life_db
      - RESTART_LIFE_REDIS_HOST=redis
      - RESTART_LIFE_REDIS_PORT=6379
      - RESTART_LIFE_JWT_SECRET_KEY=your-secret-key-here
    depends_on:
      - postgres
      - redis
    networks:
      - restart-life-network

  postgres:
    image: postgres:15-alpine
    environment:
      - POSTGRES_DB=restart_life_db
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=password
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./migrations:/docker-entrypoint-initdb.d
    networks:
      - restart-life-network

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    networks:
      - restart-life-network

volumes:
  postgres_data:
  redis_data:

networks:
  restart-life-network:
    driver: bridge
```

## 5. 性能优化和监控

### 5.1 性能优化

```go
// pkg/monitoring/metrics.go
package monitoring

import (
    "time"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    // HTTP请求指标
    HTTPRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "http_request_duration_seconds",
            Help: "HTTP请求处理时间",
        },
        []string{"method", "endpoint", "status_code"},
    )

    // 数据库查询指标
    DatabaseQueryDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "database_query_duration_seconds",
            Help: "数据库查询时间",
        },
        []string{"query_type", "table"},
    )

    // 缓存命中率
    CacheHitRate = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "cache_hits_total",
            Help: "缓存命中次数",
        },
        []string{"cache_type", "hit_miss"},
    )

    // 活跃用户数
    ActiveUsers = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "active_users_current",
            Help: "当前活跃用户数",
        },
    )
)

// RecordHTTPRequest 记录HTTP请求指标
func RecordHTTPRequest(method, endpoint, statusCode string, duration time.Duration) {
    HTTPRequestDuration.WithLabelValues(method, endpoint, statusCode).Observe(duration.Seconds())
}

// RecordDatabaseQuery 记录数据库查询指标
func RecordDatabaseQuery(queryType, table string, duration time.Duration) {
    DatabaseQueryDuration.WithLabelValues(queryType, table).Observe(duration.Seconds())
}
```

### 5.2 缓存策略

```go
// internal/repository/redis/cache_repo.go
package redis

import (
    "context"
    "encoding/json"
    "fmt"
    "time"
    
    "github.com/go-redis/redis/v8"
)

// 缓存键命名规范
const (
    CharacterDetailsKey    = "char:details:%s"
    CharacterAttributesKey = "char:attrs:%s"
    UserCharactersKey      = "user:chars:%s"
    EventTemplatesKey      = "events:templates:%s:%d"
    AchievementsKey        = "achievements:all"
    UserSessionKey         = "session:user:%s"
)

// 缓存TTL设置
const (
    CharacterDetailsTTL    = 5 * time.Minute
    CharacterAttributesTTL = 10 * time.Minute
    EventTemplatesTTL      = 1 * time.Hour
    AchievementsTTL        = 24 * time.Hour
    UserSessionTTL         = 24 * time.Hour
)

type CacheRepository struct {
    client *redis.Client
}

func NewCacheRepository(client *redis.Client) *CacheRepository {
    return &CacheRepository{
        client: client,
    }
}

// Set 设置缓存
func (r *CacheRepository) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
    data, err := json.Marshal(value)
    if err != nil {
        return fmt.Errorf("failed to marshal value: %w", err)
    }
    
    return r.client.Set(ctx, key, data, ttl).Err()
}

// Get 获取缓存
func (r *CacheRepository) Get(ctx context.Context, key string, dest interface{}) error {
    data, err := r.client.Get(ctx, key).Result()
    if err != nil {
        if err == redis.Nil {
            return fmt.Errorf("cache miss for key: %s", key)
        }
        return fmt.Errorf("failed to get cache: %w", err)
    }
    
    return json.Unmarshal([]byte(data), dest)
}

// SetCharacterDetails 缓存角色详情
func (r *CacheRepository) SetCharacterDetails(ctx context.Context, characterID string, character interface{}) error {
    key := fmt.Sprintf(CharacterDetailsKey, characterID)
    return r.Set(ctx, key, character, CharacterDetailsTTL)
}

// GetCharacterDetails 获取角色详情缓存
func (r *CacheRepository) GetCharacterDetails(ctx context.Context, characterID string, dest interface{}) error {
    key := fmt.Sprintf(CharacterDetailsKey, characterID)
    return r.Get(ctx, key, dest)
}
```

---

## 附录

### A. API接口文档

#### A.1 认证接口
```
POST /api/v1/auth/register
POST /api/v1/auth/login
POST /api/v1/auth/refresh
```

#### A.2 角色管理接口
```
POST /api/v1/characters
GET /api/v1/characters
GET /api/v1/characters/{id}
DELETE /api/v1/characters/{id}
```

#### A.3 游戏进程接口
```
POST /api/v1/characters/{id}/advance
POST /api/v1/characters/{character_id}/events/{event_id}/choose
GET /api/v1/characters/{id}/events
```

### B. 数据库迁移

```sql
-- migrations/001_initial_schema.up.sql
-- 创建初始数据库结构
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 用户表
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

-- 角色表
CREATE TABLE characters (
    character_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    character_name VARCHAR(100) NOT NULL,
    birth_country VARCHAR(100) NOT NULL,
    birth_year INTEGER NOT NULL CHECK (birth_year BETWEEN 1800 AND 2050),
    current_age INTEGER NOT NULL DEFAULT 0,
    gender VARCHAR(20) NOT NULL,
    race VARCHAR(50) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_characters_user_id ON characters(user_id);
```

### C. Go依赖管理

```go
// go.mod
module restart-life-api

go 1.21

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/lib/pq v1.10.9
    github.com/google/uuid v1.3.0
    github.com/golang-jwt/jwt/v5 v5.0.0
    github.com/go-playground/validator/v10 v10.15.5
    github.com/go-redis/redis/v8 v8.11.5
    github.com/spf13/viper v1.16.0
    golang.org/x/crypto v0.14.0
    github.com/joho/godotenv v1.4.0
    github.com/prometheus/client_golang v1.16.0
    github.com/sirupsen/logrus v1.9.3
)
```

## 6. 多平台部署配置

### 6.1 Docker多阶段构建

```dockerfile
# docker/Dockerfile.multiplatform
FROM golang:1.21-alpine AS builder

# 安装必要的工具
RUN apk add --no-cache git ca-certificates tzdata

# 设置工作目录
WORKDIR /app

# 复制go mod文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用(支持多平台)
ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -a -installsuffix cgo -ldflags="-w -s" \
    -o main cmd/server/main.go

# 最终镜像
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .
COPY --from=builder /app/config.yaml .

# 创建非root用户
RUN adduser -D -s /bin/sh appuser
USER appuser

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

EXPOSE 8080

CMD ["./main"]
```

### 6.2 Nginx多平台网关配置

```nginx
# nginx/nginx.conf
events {
    worker_connections 1024;
}

http {
    # 日志格式
    log_format main '$remote_addr - $remote_user [$time_local] "$request" '
                   '$status $body_bytes_sent "$http_referer" '
                   '"$http_user_agent" "$http_x_forwarded_for" '
                   '"$http_x_platform" "$request_time"';

    access_log /var/log/nginx/access.log main;

    # 上游服务器
    upstream api_servers {
        server api:8080 max_fails=3 fail_timeout=30s;
        keepalive 32;
    }

    # 基于平台的路由
    map $http_x_platform $backend_type {
        default         "standard";
        "unity"         "game";
        "wechat"        "miniprogram";
        "web"           "webapp";
    }

    # 限流配置
    limit_req_zone $binary_remote_addr zone=general:10m rate=10r/s;
    limit_req_zone $http_x_platform zone=platform:10m rate=20r/s;
    
    server {
        listen 80;
        server_name api.restartlife.com;
        
        # 重定向到HTTPS
        return 301 https://$server_name$request_uri;
    }

    server {
        listen 443 ssl http2;
        server_name api.restartlife.com;

        # SSL配置
        ssl_certificate /etc/nginx/ssl/cert.pem;
        ssl_certificate_key /etc/nginx/ssl/key.pem;
        ssl_protocols TLSv1.2 TLSv1.3;
        ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512;

        # 平台特定配置
        location ~ ^/api/v[0-9]+/platforms/(unity|web|wechat|mobile)/ {
            # 平台特定的限流
            limit_req zone=platform burst=5 nodelay;
            
            # 平台检测
            set $detected_platform $1;
            proxy_set_header X-Platform $detected_platform;
            
            proxy_pass http://api_servers;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }

        # 通用API接口
        location /api/ {
            # 通用限流
            limit_req zone=general burst=10 nodelay;
            
            # CORS预检请求
            if ($request_method = 'OPTIONS') {
                add_header Access-Control-Allow-Origin *;
                add_header Access-Control-Allow-Methods 'GET, POST, PUT, DELETE, OPTIONS';
                add_header Access-Control-Allow-Headers 'DNT,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range,Authorization,X-Platform,X-Version';
                add_header Access-Control-Max-Age 1728000;
                add_header Content-Type 'text/plain; charset=utf-8';
                add_header Content-Length 0;
                return 204;
            }

            # 代理到后端
            proxy_pass http://api_servers;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
            
            # 设置超时
            proxy_connect_timeout 5s;
            proxy_send_timeout 60s;
            proxy_read_timeout 60s;
            
            # 缓存设置
            proxy_cache_bypass $http_pragma;
            proxy_cache_revalidate on;
        }

        # Unity WebGL特殊处理
        location ~* \.(unity3d|unityweb|data|symbols\.json)$ {
            add_header Access-Control-Allow-Origin *;
            add_header Access-Control-Allow-Headers Range;
            
            # 支持断点续传
            proxy_set_header Range $http_range;
            proxy_set_header If-Range $http_if_range;
            proxy_no_cache $http_range $http_if_range;
            
            proxy_pass http://api_servers;
        }

        # 微信小程序域名校验
        location = /MP_verify_xxxxxxxxxxxx.txt {
            return 200 'xxxxxxxxxxxx';
            add_header Content-Type text/plain;
        }

        # 健康检查
        location /health {
            proxy_pass http://api_servers/health;
            access_log off;
        }

        # 静态文件服务
        location /static/ {
            alias /usr/share/nginx/html/static/;
            expires 30d;
            add_header Cache-Control "public, immutable";
        }
    }
}
```

### 6.3 Docker Compose多环境配置

```yaml
# docker-compose.multiplatform.yml
version: '3.8'

services:
  api:
    build:
      context: .
      dockerfile: docker/Dockerfile.multiplatform
      platforms:
        - linux/amd64
        - linux/arm64
    ports:
      - "8080:8080"
    environment:
      - RESTART_LIFE_DATABASE_HOST=postgres
      - RESTART_LIFE_DATABASE_PORT=5432
      - RESTART_LIFE_DATABASE_USER=postgres
      - RESTART_LIFE_DATABASE_PASSWORD=password
      - RESTART_LIFE_DATABASE_DBNAME=restart_life_db
      - RESTART_LIFE_REDIS_HOST=redis
      - RESTART_LIFE_REDIS_PORT=6379
      - RESTART_LIFE_JWT_SECRET_KEY=${JWT_SECRET_KEY}
      - RESTART_LIFE_CORS_ALLOWED_ORIGINS=${CORS_ALLOWED_ORIGINS}
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    networks:
      - restart-life-network
    deploy:
      replicas: 2
      resources:
        limits:
          memory: 512M
        reservations:
          memory: 256M

  postgres:
    image: postgres:15-alpine
    environment:
      - POSTGRES_DB=restart_life_db
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=password
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./migrations:/docker-entrypoint-initdb.d
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 30s
      timeout: 10s
      retries: 3
    networks:
      - restart-life-network

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 30s
      timeout: 10s
      retries: 3
    networks:
      - restart-life-network

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf
      - ./nginx/ssl:/etc/nginx/ssl
      - static_files:/usr/share/nginx/html/static
    depends_on:
      - api
    networks:
      - restart-life-network

  # Prometheus监控
  prometheus:
    image: prom/prometheus:latest
    ports:
      - "9090:9090"
    volumes:
      - ./monitoring/prometheus.yml:/etc/prometheus/prometheus.yml
    networks:
      - restart-life-network

  # Grafana仪表板
  grafana:
    image: grafana/grafana:latest
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
    volumes:
      - grafana_data:/var/lib/grafana
    networks:
      - restart-life-network

volumes:
  postgres_data:
  redis_data:
  grafana_data:
  static_files:

networks:
  restart-life-network:
    driver: bridge
```

### 6.4 Kubernetes多平台部署

```yaml
# k8s/namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: restart-life

---
# k8s/configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: api-config
  namespace: restart-life
data:
  config.yaml: |
    server:
      port: 8080
      host: "0.0.0.0"
    database:
      host: postgres-service
      port: 5432
      user: postgres
      dbname: restart_life_db
    redis:
      host: redis-service
      port: 6379
    cors:
      allowed_origins:
        - "https://restartlife.com"
        - "https://*.restartlife.com"

---
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: api-deployment
  namespace: restart-life
spec:
  replicas: 3
  selector:
    matchLabels:
      app: restart-life-api
  template:
    metadata:
      labels:
        app: restart-life-api
    spec:
      containers:
      - name: api
        image: restart-life/api:latest
        ports:
        - containerPort: 8080
        env:
        - name: RESTART_LIFE_JWT_SECRET_KEY
          valueFrom:
            secretKeyRef:
              name: api-secrets
              key: jwt-secret
        volumeMounts:
        - name: config-volume
          mountPath: /app/config.yaml
          subPath: config.yaml
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
      volumes:
      - name: config-volume
        configMap:
          name: api-config

---
# k8s/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: api-service
  namespace: restart-life
spec:
  selector:
    app: restart-life-api
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
  type: ClusterIP

---
# k8s/ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: api-ingress
  namespace: restart-life
  annotations:
    kubernetes.io/ingress.class: nginx
    cert-manager.io/cluster-issuer: letsencrypt-prod
    nginx.ingress.kubernetes.io/cors-allow-origin: "*"
    nginx.ingress.kubernetes.io/cors-allow-methods: "GET, POST, PUT, DELETE, OPTIONS"
    nginx.ingress.kubernetes.io/cors-allow-headers: "DNT,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range,Authorization,X-Platform,X-Version"
spec:
  tls:
  - hosts:
    - api.restartlife.com
    secretName: api-tls
  rules:
  - host: api.restartlife.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: api-service
            port:
              number: 80
```

### 6.5 CI/CD多平台构建

```yaml
# .github/workflows/multiplatform-deploy.yml
name: Multi-Platform Build and Deploy

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

env:
  REGISTRY: ghcr.io
  IMAGE_NAME: restart-life/api

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: 1.21
    
    - name: Run tests
      run: |
        go mod download
        go test -v ./...
        go test -race -coverprofile=coverage.out ./...
    
    - name: Upload coverage to Codecov
      uses: codecov/codecov-action@v3

  build:
    needs: test
    runs-on: ubuntu-latest
    steps:
    - name: Checkout
      uses: actions/checkout@v3

    - name: Set up Docker Buildx
      uses: docker/setup-buildx-action@v2

    - name: Log in to Container Registry
      uses: docker/login-action@v2
      with:
        registry: ${{ env.REGISTRY }}
        username: ${{ github.actor }}
        password: ${{ secrets.GITHUB_TOKEN }}

    - name: Extract metadata
      id: meta
      uses: docker/metadata-action@v4
      with:
        images: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}
        tags: |
          type=ref,event=branch
          type=ref,event=pr
          type=sha,prefix={{branch}}-
          type=raw,value=latest,enable={{is_default_branch}}

    - name: Build and push Docker image
      uses: docker/build-push-action@v4
      with:
        context: .
        file: docker/Dockerfile.multiplatform
        platforms: linux/amd64,linux/arm64
        push: true
        tags: ${{ steps.meta.outputs.tags }}
        labels: ${{ steps.meta.outputs.labels }}
        cache-from: type=gha
        cache-to: type=gha,mode=max

  deploy:
    needs: build
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    steps:
    - name: Deploy to Kubernetes
      run: |
        # 这里添加部署到K8s集群的逻辑
        echo "Deploying to production cluster"
```

### 6.6 环境配置管理

```yaml
# config/production.yaml
server:
  port: 8080
  host: "0.0.0.0"
  read_timeout: "30s"
  write_timeout: "30s"

database:
  host: ${DB_HOST}
  port: ${DB_PORT}
  user: ${DB_USER}
  password: ${DB_PASSWORD}
  dbname: ${DB_NAME}
  max_open_conns: 50
  max_idle_conns: 10

redis:
  host: ${REDIS_HOST}
  port: ${REDIS_PORT}
  password: ${REDIS_PASSWORD}
  pool_size: 20

cors:
  allowed_origins:
    - "https://restartlife.com"
    - "https://www.restartlife.com"
    - "https://app.restartlife.com"
    - "https://miniprogram.restartlife.com"
  allowed_methods:
    - "GET"
    - "POST" 
    - "PUT"
    - "DELETE"
    - "OPTIONS"
    - "PATCH"

rate_limit:
  requests_per_minute: 60
  burst_size: 10

jwt:
  secret_key: ${JWT_SECRET_KEY}
  expires_in: "24h"
  refresh_in: "168h"

logging:
  level: "info"
  format: "json"

monitoring:
  enable_metrics: true
  metrics_port: 9090
  health_check_path: "/health"
```

### D. 版本历史

- **v1.0** (2025-01-26)：初版后端技术设计文档完成，包含多平台支持

---

*本文档为《重启人生》游戏后端技术设计文档，专注于Go服务端的多平台技术架构和实现方案。* 