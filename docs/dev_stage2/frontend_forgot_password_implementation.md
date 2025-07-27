# 前端忘记密码功能实现

## 📋 功能概述

为前端登录页面添加了完整的忘记密码功能，支持三步重置密码流程：

1. **输入邮箱** - 用户输入邮箱地址
2. **验证验证码** - 输入邮箱收到的6位验证码
3. **设置新密码** - 输入新密码完成重置

## 🔧 实现细节

### API接口更新

在 `src/services/api.ts` 中新增了重置密码相关的API调用：

```typescript
// 发送验证码
sendVerificationCode: (email: string) =>
  api.post<ApiResponse>('/auth/send-verification-code', {
    email,
  }),

// 验证验证码
verifyCode: (email: string, code: string) =>
  api.post<ApiResponse<{ reset_token: string }>>('/auth/verify-code', {
    email,
    code,
  }),

// 重置密码
resetPassword: (resetToken: string, newPassword: string) =>
  api.post<ApiResponse>('/auth/reset-password', {
    reset_token: resetToken,
    new_password: newPassword,
  }),
```

### 前端组件更新

在 `src/pages/LoginPage.tsx` 中添加了：

#### 1. 新增状态管理
```typescript
const [resetLoading, setResetLoading] = useState(false)
const [resetStep, setResetStep] = useState(1) // 1: 输入邮箱, 2: 输入验证码, 3: 设置新密码
const [resetEmail, setResetEmail] = useState('')
const [resetToken, setResetToken] = useState('')
```

#### 2. 处理函数
- `handleSendResetCode` - 发送验证码到邮箱
- `handleVerifyCode` - 验证验证码并获取重置令牌
- `handleResetPassword` - 使用重置令牌设置新密码

#### 3. UI组件
- 新增"重置密码"Tab页
- 三步式表单界面，根据 `resetStep` 状态显示不同步骤
- 在登录表单下方添加"忘记密码？"链接

## 🎯 用户流程

### 步骤1：输入邮箱
- 用户点击"忘记密码？"链接
- 切换到"重置密码"Tab页
- 输入邮箱地址，点击"发送验证码"

### 步骤2：验证验证码
- 系统显示验证码已发送的提示
- 用户输入收到的6位验证码
- 点击"验证"按钮进行验证

### 步骤3：设置新密码
- 验证成功后，显示密码设置界面
- 用户输入新密码和确认密码
- 点击"重置密码"完成重置

## 🔗 后端API对应

该功能完全对应后端的重置密码API：

- `POST /api/v1/auth/send-verification-code` - 发送验证码
- `POST /api/v1/auth/verify-code` - 验证验证码并获取重置令牌
- `POST /api/v1/auth/reset-password` - 使用令牌重置密码

详细API文档参见：[两步密码重置API文档](./two_step_password_reset_api.md)

## ✅ 功能特点

- **安全性** - 使用一次性验证码和临时重置令牌
- **用户友好** - 清晰的三步流程指引
- **错误处理** - 完善的错误提示和异常处理
- **状态管理** - 合理的步骤状态控制
- **响应式** - 适配移动端和桌面端

## 🎨 界面设计

- 与现有登录/注册界面风格保持一致
- 使用Ant Design组件库
- 清晰的步骤指示和用户引导
- 友好的错误提示和成功反馈

## 📅 实现时间

**实现日期**: 2025-07-27
**相关分支**: `xucheng/feature/v0.2/game-engine`
**状态**: ✅ 已完成
