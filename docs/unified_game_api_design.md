# 游戏API统一设计方案

## 问题背景

之前的设计中，游戏操作需要两个不同的API：
- `POST /api/v1/game/advance/{character_id}` - 推进游戏年龄
- `POST /api/v1/game/decision/{character_id}` - 做出决策

这种设计存在以下问题：
1. **用户体验差**: 需要知道当前状态来决定调用哪个API
2. **前端逻辑复杂**: 需要判断游戏状态来选择合适的API
3. **API职责重叠**: 两个API实际上都是推进游戏进程
4. **开发维护成本高**: 需要维护两套类似的逻辑

## 优化方案

### 统一为智能推进API

将游戏推进统一为一个智能API：`POST /api/v1/game/advance/{character_id}`

### 智能判断逻辑

```
用户请求 → 检查游戏状态 → 智能处理
                ↓
        有待处理决策？
            ↓          ↓
          是            否
            ↓          ↓
    需要option_type   自动推进年龄
      处理决策        生成新事件
```

### API设计

#### 请求格式
```json
{
  "option_type": "conservative|moderate|aggressive"  // 可选，仅在有待处理决策时需要
}
```

#### 响应格式
```json
{
  "success": true,
  "data": {
    "character_id": "uuid",
    "current_age": 25,
    "is_game_active": true,
    "pending_decision": {  // 如果有新的决策
      "options": {
        "conservative": { "option_text": "稳健选择" },
        "moderate": { "option_text": "中庸选择" },
        "aggressive": { "option_text": "激进选择" }
      }
    },
    "key_events": [  // 新发生的事件
      {
        "description": "你遇到了一个重要的人生选择...",
        "impact": "这将影响你的未来发展"
      }
    ]
  }
}
```

## 使用场景

### 场景1：正常推进游戏
```http
POST /api/v1/game/advance/123
Content-Type: application/json

{}
```

**系统处理**：
1. 检查无待处理决策
2. 年龄+1
3. 调用AI生成新事件（可能有决策）
4. 返回新状态

### 场景2：有待处理决策时推进
```http
POST /api/v1/game/advance/123
Content-Type: application/json

{
  "option_type": "aggressive"
}
```

**系统处理**：
1. 检查有待处理决策
2. 处理用户选择
3. 调用AI生成决策结果
4. 可能生成新事件和决策
5. 返回新状态

### 场景3：错误处理
```http
POST /api/v1/game/advance/123
Content-Type: application/json

{}
```

**系统响应**（有待处理决策时）：
```json
{
  "success": false,
  "error": {
    "code": 1001,
    "message": "有待处理的决策，请提供option_type参数"
  }
}
```

## 技术实现

### GameHandler层
```go
func (h *GameHandler) AdvanceGame(c *gin.Context) {
    var req models.GameProgressRequest
    c.ShouldBindJSON(&req)  // 可选参数

    gameState, err := h.gameService.AdvanceGameSmart(
        c.Request.Context(),
        characterID,
        req.OptionType
    )
    // 处理响应...
}
```

### GameService层
```go
func (s *gameService) AdvanceGameSmart(ctx context.Context, characterID string, optionType string) (*models.GameState, error) {
    gameState, err := s.GetGameState(ctx, characterID)

    if gameState.PendingDecision != nil {
        // 有待处理决策，必须提供option_type
        if optionType == "" {
            return nil, fmt.Errorf("有待处理的决策，请提供option_type参数")
        }
        return s.MakeDecision(ctx, characterID, optionType)
    }

    // 没有待处理决策，正常推进
    return s.AdvanceGame(ctx, characterID)
}
```

## 向后兼容

为了保持向后兼容性，保留旧的`MakeDecision` API但标记为废弃：

```go
// @deprecated
func (h *GameHandler) MakeDecision(c *gin.Context) {
    // 内部调用新的AdvanceGameSmart方法
    gameState, err := h.gameService.AdvanceGameSmart(ctx, characterID, req.OptionType)
}
```

## 优势总结

### 用户体验提升
- **操作简化**: 只需要一个"推进游戏"按钮
- **智能处理**: 系统自动判断当前应该执行什么操作
- **错误提示清晰**: 明确告知用户需要什么操作

### 前端开发简化
- **单一API**: 前端只需要调用一个API
- **状态无关**: 不需要判断游戏状态来选择API
- **错误处理统一**: 统一的错误码和处理逻辑

### 后端架构优化
- **逻辑集中**: 游戏推进逻辑集中管理
- **代码复用**: 减少重复代码
- **维护成本低**: 只需要维护一套推进逻辑

### AI调用优化
- **统一入口**: 所有游戏推进都通过统一的AI调用
- **上下文一致**: 保证AI处理的连贯性
- **性能提升**: 避免多次AI调用

## 迁移指南

### 前端迁移
```javascript
// 旧方式 - 需要判断状态
if (gameState.pendingDecision) {
    await api.makeDecision(characterId, { option_type: "conservative" });
} else {
    await api.advanceGame(characterId);
}

// 新方式 - 统一调用
await api.advanceGame(characterId, { option_type: optionType });
```

### API调用示例
```javascript
// 推进游戏（无决策）
const response = await fetch(`/api/v1/game/advance/${characterId}`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({})
});

// 推进游戏（有决策）
const response = await fetch(`/api/v1/game/advance/${characterId}`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ option_type: "aggressive" })
});
```

这个设计极大地简化了游戏交互逻辑，提供了更好的用户体验和更清晰的API设计。
