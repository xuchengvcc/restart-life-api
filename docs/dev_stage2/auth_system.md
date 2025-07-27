# 用户认证系统

这是Restart Life API的用户认证系统，实现了完整的用户注册、登录、Token管理和用户信息管理功能。

## 功能特性

- ✅ 用户注册/登录
- ✅ JWT Token认证
- ✅ 访问Token和刷新Token
- ✅ 密码安全加密（bcrypt）
- ✅ 用户信息管理
- ✅ 密码修改
- ✅ 多平台登录支持（用户名/邮箱）
- ✅ Token中间件保护
- ✅ 完整的错误处理

## API接口

### 1. 用户注册
```http
POST /api/v1/auth/register
Content-Type: application/json

{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
}
```

**响应：**
```json
{
    "success": true,
    "data": {
        "user": {
            "user_id": 1,
            "username": "testuser",
            "email": "test@example.com",
            "created_at": 1640995200000,
            "updated_at": 1640995200000,
            "last_login": 1640995200000,
            "is_active": true
        },
        "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "expires_at": 1641081600
    }
}
```

### 2. 用户登录
```http
POST /api/v1/auth/login
Content-Type: application/json

{
    "username": "testuser",  // 支持用户名或邮箱
    "password": "password123"
}
```

### 3. 刷新Token
```http
POST /api/v1/auth/refresh
Content-Type: application/json

{
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### 4. 获取用户信息
```http
GET /api/v1/auth/profile
Authorization: Bearer YOUR_ACCESS_TOKEN
```

### 5. 更新用户信息
```http
PUT /api/v1/auth/profile
Authorization: Bearer YOUR_ACCESS_TOKEN
Content-Type: application/json

{
    "bio": "这是我的个人简介",
    "birth_date": "1990-01-01",
    "gender": 1,
    "country": "中国",
    "avatar_url": "https://example.com/avatar.jpg"
}
```

### 6. 修改密码
```http
POST /api/v1/auth/change-password
Authorization: Bearer YOUR_ACCESS_TOKEN
Content-Type: application/json

{
    "old_password": "password123",
    "new_password": "newpassword123"
}
```

### 7. 登出
```http
POST /api/v1/auth/logout
Authorization: Bearer YOUR_ACCESS_TOKEN
```

## 认证流程

1. **注册/登录** → 获取访问Token和刷新Token
2. **API调用** → 在Header中携带 `Authorization: Bearer ACCESS_TOKEN`
3. **Token过期** → 使用刷新Token获取新的访问Token
4. **登出** → 客户端清除Token（服务端暂时为无状态）

## 安全特性

- **密码加密**: 使用bcrypt算法，成本因子为12
- **JWT签名**: 使用HMAC-SHA256算法
- **Token过期**: 访问Token 24小时，刷新Token 7天
- **输入验证**: 严格的参数验证和清理
- **错误处理**: 统一的错误响应格式

## 数据库表结构

```sql
CREATE TABLE IF NOT EXISTS user_tab (
    user_id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '用户ID',
    username VARCHAR(64) UNIQUE NOT NULL COMMENT '用户名',
    email VARCHAR(128) UNIQUE NOT NULL COMMENT '邮箱',
    password_hash VARCHAR(255) NOT NULL COMMENT '密码哈希',
    created_at BIGINT UNSIGNED NOT NULL COMMENT '毫秒时间戳',
    updated_at BIGINT UNSIGNED NOT NULL COMMENT '毫秒时间戳',
    last_login BIGINT UNSIGNED NULL COMMENT '最后登录时间 毫秒时间戳',
    is_active BOOLEAN DEFAULT TRUE COMMENT '是否激活 1:正常 0:禁用/未激活',
    avatar_url VARCHAR(500) COMMENT '头像URL',
    bio TEXT COMMENT '个人简介',
    birth_date DATE COMMENT '出生日期',
    gender TINYINT COMMENT '性别 0:未知 1:男 2:女 3:其他',
    country VARCHAR(100) COMMENT '国家'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

## 配置说明

在 `configs/test.yaml` 中配置认证相关参数：

```yaml
auth:
  jwt_secret: "your-jwt-secret-key"  # JWT签名密钥（生产环境请使用强密钥）
  jwt_expiry: 24h                    # 访问Token过期时间
  refresh_expiry: 168h               # 刷新Token过期时间（7天）
```

## 错误代码

| 错误代码 | 说明 |
|---------|------|
| 1001 | 用户名或密码错误 |
| 1002 | Token已过期 |
| 1003 | Token无效 |
| 1004 | 用户不存在 |
| 1005 | 用户已存在 |
| 1006 | 权限不足 |
| 2001 | 数据验证失败 |
| 2002 | 输入参数无效 |
| 5001 | 内部服务器错误 |
| 5002 | 数据库操作错误 |

## 测试

运行单元测试：
```bash
cd cmd/server
go test -v
```

运行API测试：
```bash
# 确保服务器正在运行
./scripts/test_auth_api.sh
```

## 使用示例

### JavaScript/TypeScript客户端

```javascript
// 注册用户
const registerResponse = await fetch('/api/v1/auth/register', {
    method: 'POST',
    headers: {
        'Content-Type': 'application/json',
    },
    body: JSON.stringify({
        username: 'testuser',
        email: 'test@example.com',
        password: 'password123'
    })
});

const { data } = await registerResponse.json();
const { access_token, refresh_token } = data;

// 存储Token
localStorage.setItem('access_token', access_token);
localStorage.setItem('refresh_token', refresh_token);

// 调用需要认证的API
const profileResponse = await fetch('/api/v1/auth/profile', {
    headers: {
        'Authorization': `Bearer ${access_token}`
    }
});
```

### Go客户端

```go
// 创建认证请求
type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

req := LoginRequest{
    Username: "testuser",
    Password: "password123",
}

// 发送登录请求
jsonData, _ := json.Marshal(req)
resp, err := http.Post("http://localhost:8080/api/v1/auth/login",
    "application/json", bytes.NewBuffer(jsonData))

// 解析响应获取Token
var authResponse struct {
    Success bool `json:"success"`
    Data struct {
        AccessToken string `json:"access_token"`
        User struct {
            UserID   uint   `json:"user_id"`
            Username string `json:"username"`
            Email    string `json:"email"`
        } `json:"user"`
    } `json:"data"`
}

json.NewDecoder(resp.Body).Decode(&authResponse)
accessToken := authResponse.Data.AccessToken
```

## 下一步开发

- [ ] 邮箱验证功能
- [ ] 忘记密码功能
- [ ] 第三方登录（OAuth2）
- [ ] 多设备登录管理
- [ ] 登录日志记录
- [ ] 账户锁定机制
- [ ] Token黑名单机制
