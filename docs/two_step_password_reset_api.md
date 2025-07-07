# 两步重置密码 API 文档

本文档描述了新的两步重置密码流程，这种设计更加安全且符合前后端分离架构的最佳实践。

## 流程概述

两步重置密码流程分为以下步骤：

1. **发送验证码** - 用户提供邮箱，系统发送验证码
2. **验证验证码并获取重置令牌** - 验证验证码正确后，返回临时重置令牌
3. **使用重置令牌重置密码** - 用户提供重置令牌和新密码完成重置

## API 接口

### 1. 发送重置密码验证码

**请求**
```
POST /api/v1/auth/send-verification-code
Content-Type: application/json

{
    "email": "user@example.com"
}
```

**响应**
```json
{
    "success": true,
    "message": "验证码已发送到您的邮箱",
    "data": {
        "message": "验证码已发送到您的邮箱"
    }
}
```

**错误响应**
- `400` - 邮箱地址无效
- `429` - 发送过于频繁，请稍后再试
- `500` - 邮件发送失败

---

### 2. 验证验证码并获取重置令牌

**请求**
```
POST /api/v1/auth/verify-code
Content-Type: application/json

{
    "email": "user@example.com",
    "code": "123456"
}
```

**响应**
```json
{
    "success": true,
    "message": "验证码验证成功",
    "data": {
        "message": "验证码验证成功",
        "reset_token": "eyJhbGciOiJIUzI1NiI..."
    }
}
```

**错误响应**
- `400` - 验证码无效
- `404` - 用户不存在
- `410` - 验证码已过期
- `500` - 验证失败

**重要说明**
- 验证码验证成功后会立即被删除（一次性使用）
- 重置令牌有效期为10分钟
- 重置令牌只能用于重置密码，不是访问令牌

---

### 3. 使用重置令牌重置密码

**请求**
```
POST /api/v1/auth/reset-password
Content-Type: application/json

{
    "reset_token": "eyJhbGciOiJIUzI1NiI...",
    "new_password": "newpassword123"
}
```

**响应**
```json
{
    "success": true,
    "message": "密码重置成功",
    "data": {
        "message": "密码重置成功"
    }
}
```

**错误响应**
- `400` - 新密码格式不符合要求
- `401` - 重置令牌无效或已过期
- `404` - 用户不存在
- `500` - 重置失败

**重要说明**
- 重置令牌使用后会立即被删除（一次性使用）
- 新密码最少6位字符
- 重置成功后用户需要使用新密码重新登录

## 前端集成示例

### JavaScript/TypeScript 示例

```javascript
class PasswordResetService {
    constructor(baseURL) {
        this.baseURL = baseURL;
    }

    // 步骤1：发送验证码
    async sendVerificationCode(email) {
        const response = await fetch(`${this.baseURL}/auth/send-verification-code`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email })
        });
        return response.json();
    }

    // 步骤2：验证验证码并获取重置令牌
    async verifyCodeAndGetToken(email, code) {
        const response = await fetch(`${this.baseURL}/auth/verify-code`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, code })
        });
        const result = await response.json();
        return result.data?.reset_token;
    }

    // 步骤3：使用重置令牌重置密码
    async resetPasswordWithToken(resetToken, newPassword) {
        const response = await fetch(`${this.baseURL}/auth/reset-password`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                reset_token: resetToken,
                new_password: newPassword
            })
        });
        return response.json();
    }
}

// 使用示例
const resetService = new PasswordResetService('http://localhost:8080/api/v1');

async function resetPassword(email, code, newPassword) {
    try {
        // 步骤1：发送验证码（通常在前面的页面完成）
        await resetService.sendVerificationCode(email);

        // 步骤2：验证验证码并获取重置令牌
        const resetToken = await resetService.verifyCodeAndGetToken(email, code);

        // 步骤3：使用重置令牌重置密码
        const result = await resetService.resetPasswordWithToken(resetToken, newPassword);

        if (result.success) {
            console.log('密码重置成功');
            // 跳转到登录页面
        }
    } catch (error) {
        console.error('密码重置失败:', error);
    }
}
```

## 安全特性

1. **验证码一次性使用** - 验证码验证成功后立即删除，防止重放攻击
2. **重置令牌时效性** - 重置令牌有效期10分钟，减少攻击窗口
3. **令牌一次性使用** - 重置令牌使用后立即删除
4. **邮箱验证** - 只有邮箱所有者才能收到验证码
5. **分离式设计** - 验证和重置分开，增加安全性

## 测试

使用提供的测试脚本进行完整流程测试：

```bash
./scripts/test_two_step_password_reset.sh
```

测试脚本会验证完整的两步重置密码流程，包括：
- 发送验证码
- 验证验证码并获取重置令牌
- 使用重置令牌重置密码
- 使用新密码登录验证

## 故障排除

### 常见问题

1. **验证码未收到**
   - 检查邮箱地址是否正确
   - 检查垃圾邮件文件夹
   - 确认邮件服务配置正确

2. **验证码验证失败**
   - 确认验证码输入正确
   - 检查验证码是否已过期（有效期5分钟）
   - 验证码只能使用一次

3. **重置令牌无效**
   - 检查令牌是否已过期（有效期10分钟）
   - 令牌只能使用一次
   - 确认令牌格式正确

4. **密码重置失败**
   - 检查新密码是否符合要求（最少6位）
   - 确认用户账户存在且有效
   - 检查服务器日志获取详细错误信息
