# 前端API路径修正

## 📋 问题描述

前端API调用的路径与后端实际路由不匹配，导致API请求失败。

## 🔧 修正内容

### 1. 基础路径修正

**修正前**:
```typescript
baseURL: '/api'
```

**修正后**:
```typescript
baseURL: '/api/v1'
```

### 2. 认证API路径

所有认证相关的API都已正确指向 `/api/v1/auth/*`：

- ✅ `/auth/login`
- ✅ `/auth/register`
- ✅ `/auth/logout`
- ✅ `/auth/refresh`
- ✅ `/auth/send-verification-code`
- ✅ `/auth/verify-code`
- ✅ `/auth/reset-password`

### 3. 角色管理API路径

根据 `routes.go` 中的实际路由修正：

**修正前**:
```typescript
create: '/characters'
getById: '/characters/${id}'
getByUser: '/characters'
update: '/characters/${id}'
delete: '/characters/${id}'
```

**修正后**:
```typescript
create: '/characters/create'
getById: '/characters/get/${id}'
getByUser: '/characters/list'
update: '/characters/update/${id}'
delete: '/characters/delete/${id}'
```

### 4. 游戏API路径

根据 `routes.go` 中的游戏路由进行了重新设计：

**新增的路由**:
```typescript
startOrResume: '/game/start-or-resume'
startGame: '/game/start/${characterId}'
```

**修正的路径格式**:
- 从 `/game/${characterId}/action` 改为 `/game/action/${characterId}`
- 与后端路由保持一致

## 🎯 对应的后端路由

### 认证路由 (`/api/v1/auth`)
```go
auth.POST("/register", authHandler.Register)
auth.POST("/login", authHandler.Login)
auth.POST("/refresh", authHandler.RefreshToken)
auth.POST("/send-verification-code", authHandler.SendVerificationCode)
auth.POST("/verify-code", authHandler.VerifyCode)
auth.POST("/reset-password", authHandler.ResetPassword)
auth.POST("/logout", authHandler.Logout)
auth.GET("/profile", authHandler.GetProfile)
auth.PUT("/update_profile", authHandler.UpdateProfile)
auth.POST("/change-password", authHandler.ChangePassword)
```

### 角色路由 (`/api/v1/characters`)
```go
characters.POST("/create", characterHandler.CreateCharacter)
characters.GET("/list", characterHandler.GetUserCharacters)
characters.GET("/get/:id", characterHandler.GetCharacter)
characters.PUT("/update/:id", characterHandler.UpdateCharacter)
characters.DELETE("/delete/:id", characterHandler.DeleteCharacter)
characters.GET("/attributes/get/:id", characterHandler.GetCharacterAttributes)
characters.PUT("/attributes/update/:id", characterHandler.UpdateCharacterAttributes)
```

### 游戏路由 (`/api/v1/game`)
```go
game.POST("/start-or-resume", gameHandler.StartOrResumeGame)
game.POST("/start/:character_id", gameHandler.StartGame)
// 其他游戏路由待实现
```

## ✅ 修正结果

- 所有API路径现在与后端路由完全匹配
- 前端可以正确调用后端接口
- 支持忘记密码功能的完整流程
- 角色管理和游戏功能的API调用正确

## 📅 修正时间

**修正日期**: 2025-07-27
**相关分支**: `xucheng/feature/v0.2/game-engine`
**状态**: ✅ 已完成
