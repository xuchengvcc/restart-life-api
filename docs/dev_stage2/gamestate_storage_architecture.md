# GameState存储架构设计

## 设计决策

**推荐使用：Redis + Database 混合存储模式**

## 架构分析

### 当前GameState组成
```go
type GameState struct {
    // 基本信息（来自character表）
    CharacterID, CharacterName, CurrentAge, LifeStage, Attributes

    // 实时状态（需要缓存）
    KeyEvents []Event
    PendingDecision *Decision

    // 计算字段（可重建）
    Education, Career, Location, MaritalStatus
    HealthStatus, WealthLevel, Relationships

    // 会话数据
    Money, LastYearDescription
}
```

## 存储策略

### 1. Database存储（持久化）

**character_tab**：
- 角色基本信息、属性、游戏进度
- 当前年龄、人生阶段、游戏完成状态
- 教育、职业、健康、财富等状态

**game_events**：
- 关键事件历史
- 按(character_id, age)索引优化查询

**不需要game_state表**：
- GameState是计算得出的复合对象
- 可以从character和events重建

### 2. Redis存储（性能缓存）

**Key设计**：
```
game_state:{character_id}
```

**存储内容**：
- 完整的GameState JSON对象
- 包含KeyEvents数组和PendingDecision

**TTL策略**：
- 24小时自动过期
- 用户活跃时续期
- 减少内存占用

### 3. 数据流设计

```
读取GameState：
1. 优先从Redis获取
2. Redis miss时从Database重建
3. 重建后写入Redis（TTL 24h）

更新GameState：
1. 重要操作更新Database
2. 同步更新Redis缓存
3. 确保数据一致性
```

## 实现优势

### 性能优势
- **毫秒级响应**：Redis提供极快的读取速度
- **减少计算**：避免重复的GameState重建
- **并发友好**：Redis天然支持高并发访问

### 数据安全
- **持久化保证**：Database确保数据永不丢失
- **灾难恢复**：Redis故障时可从Database重建
- **数据一致性**：关键操作同步更新两端

### 资源优化
- **内存管理**：TTL自动清理非活跃用户数据
- **存储分离**：热数据在Redis，冷数据在Database
- **扩展性好**：可独立扩展Redis和Database

## 关键方法实现

### GetGameState逻辑
```go
func (s *gameService) GetGameState(ctx context.Context, characterID string) (*models.GameState, error) {
    // 1. 尝试从Redis获取
    if gameState := s.getGameStateFromRedis(characterID); gameState != nil {
        return gameState, nil
    }

    // 2. 从Database重建
    gameState, err := s.buildGameStateFromDatabase(ctx, characterID)
    if err != nil {
        return nil, err
    }

    // 3. 写入Redis缓存
    s.saveGameStateToRedis(characterID, gameState)

    return gameState, nil
}
```

### SaveGame逻辑
```go
func (s *gameService) SaveGame(ctx context.Context, characterID string) error {
    gameState, err := s.GetGameState(ctx, characterID)
    if err != nil {
        return err
    }

    // 1. 更新Database（character + events）
    if err := s.updateCharacterInDB(ctx, gameState); err != nil {
        return err
    }

    // 2. 更新Redis缓存
    s.saveGameStateToRedis(characterID, gameState)

    return nil
}
```

## 数据一致性保证

### 写入顺序
1. **Database优先**：关键状态先写入数据库
2. **Redis同步**：数据库成功后更新缓存
3. **错误处理**：Database失败时不更新Redis

### 缓存失效
- 角色属性变化时清除Redis缓存
- 强制从Database重建确保一致性
- 防止脏数据问题

## 不需要额外表的原因

### GameState是计算产物
- **Education, Career等**：从character字段计算得出
- **Relationships, PersonalGrowth**：基于年龄计算
- **WealthLevel**：从Money字段转换

### 实时性要求
- PendingDecision：临时决策状态，不需要持久化
- KeyEvents：从game_events表加载
- 游戏进行中的临时状态适合Redis存储

## 总结

**推荐方案**：不创建game_state表，使用Redis + 现有Database表
- **简化设计**：减少表之间的复杂关系
- **提升性能**：Redis提供极速访问
- **降低成本**：减少数据库存储和维护成本
- **灵活扩展**：便于后续架构调整

这种设计在保证数据安全的同时，最大化了性能和用户体验。
