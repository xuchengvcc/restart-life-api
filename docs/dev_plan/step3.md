# Step 3: 高级功能模块 (v0.3.x)

## 📋 阶段概述
- **阶段目标**: 实现高级游戏功能，包括关系网络、成就系统和数据管理
- **预计时间**: 2025-02-17 ~ 2025-03-02 (2周)
- **关键交付**: 完整的关系系统、成就系统和存档统计功能

## 🎯 详细任务

### 任务7: 关系网络系统
- **分支**: `xucheng/feature/v0.3/relationship-system`
- **负责人**: xucheng
- **预计时间**: 4-5天

#### 开发目标
- [ ] 社会关系管理
- [ ] 关系质量评估
- [ ] 关系事件处理
- [ ] 关系影响计算

#### 详细任务清单
- [ ] 实现关系类型定义和管理
- [ ] 实现关系质量评分系统
- [ ] 实现关系事件生成和处理
- [ ] 实现关系对游戏的影响计算
- [ ] 实现关系网络图谱生成
- [ ] 实现关系互动选择系统
- [ ] 实现关系历史追踪
- [ ] 实现关系冲突和和解机制

#### 交付物
- [ ] 关系管理API (`internal/api/handlers/relationship.go`)
- [ ] 关系网络模型 (`internal/models/relationship.go`)
- [ ] 关系影响算法 (`pkg/utils/relationship_calculator.go`)
- [ ] 关系事件处理器 (`internal/services/relationship_service.go`)
- [ ] 关系数据层 (`internal/repository/postgres/relationship_repo.go`)
- [ ] 关系网络可视化数据接口

#### API接口设计
```http
GET  /api/v1/relationships/:character_id         # 获取角色关系网络
POST /api/v1/relationships/:character_id         # 创建新关系
PUT  /api/v1/relationships/:relationship_id      # 更新关系状态
GET  /api/v1/relationships/:relationship_id/history # 获取关系历史
POST /api/v1/relationships/:relationship_id/interact # 关系互动
GET  /api/v1/relationships/network/:character_id  # 获取关系网络图
GET  /api/v1/relationships/events/:character_id   # 获取关系事件
```

#### 数据模型
```go
type Relationship struct {
    RelationshipID   string    `json:"relationship_id" db:"relationship_id"`
    CharacterID      string    `json:"character_id" db:"character_id"`
    RelatedPersonID  string    `json:"related_person_id" db:"related_person_id"`
    RelationType     string    `json:"relation_type" db:"relation_type"`
    Intimacy         int       `json:"intimacy" db:"intimacy"`         // 亲密度 0-100
    Trust            int       `json:"trust" db:"trust"`               // 信任度 0-100
    Support          int       `json:"support" db:"support"`           // 支持度 0-100
    Status           string    `json:"status" db:"status"`             // 关系状态
    StartAge         int       `json:"start_age" db:"start_age"`
    LastInteraction  time.Time `json:"last_interaction" db:"last_interaction"`
    CreatedAt        time.Time `json:"created_at" db:"created_at"`
    UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

type RelatedPerson struct {
    PersonID     string `json:"person_id" db:"person_id"`
    Name         string `json:"name" db:"name"`
    Age          int    `json:"age" db:"age"`
    Gender       string `json:"gender" db:"gender"`
    Occupation   string `json:"occupation" db:"occupation"`
    Personality  string `json:"personality" db:"personality"`
    Background   string `json:"background" db:"background"`
}

type RelationshipEvent struct {
    EventID        string    `json:"event_id" db:"event_id"`
    RelationshipID string    `json:"relationship_id" db:"relationship_id"`
    EventType      string    `json:"event_type" db:"event_type"`
    Description    string    `json:"description" db:"description"`
    ImpactType     string    `json:"impact_type" db:"impact_type"`
    ImpactValue    int       `json:"impact_value" db:"impact_value"`
    OccurredAt     time.Time `json:"occurred_at" db:"occurred_at"`
}
```

#### 关系类型分类
```go
const (
    // 家庭关系
    RelationTypeParent    = "parent"
    RelationTypeChild     = "child"
    RelationTypeSibling   = "sibling"
    RelationTypeSpouse    = "spouse"
    
    // 社会关系
    RelationTypeFriend    = "friend"
    RelationTypeColleague = "colleague"
    RelationTypeNeighbor  = "neighbor"
    RelationTypeMentor    = "mentor"
    
    // 专业关系
    RelationTypeBoss      = "boss"
    RelationTypeEmployee  = "employee"
    RelationTypePartner   = "partner"
    RelationTypeClient    = "client"
)
```

---

### 任务8: 成就系统
- **分支**: `xucheng/feature/v0.3/achievement-system`
- **负责人**: xucheng
- **预计时间**: 3-4天

#### 开发目标
- [ ] 成就定义和管理
- [ ] 成就解锁逻辑
- [ ] 成就统计和展示
- [ ] 特殊成就处理

#### 详细任务清单
- [ ] 设计成就分类和等级系统
- [ ] 实现成就触发条件检查
- [ ] 实现成就解锁和通知系统
- [ ] 实现成就进度追踪
- [ ] 实现成就统计和排行
- [ ] 实现特殊成就和隐藏成就
- [ ] 实现成就奖励系统
- [ ] 成就数据导出功能

#### 交付物
- [ ] 成就系统API (`internal/api/handlers/achievement.go`)
- [ ] 成就解锁引擎 (`internal/services/achievement_service.go`)
- [ ] 成就数据模型 (`internal/models/achievement.go`)
- [ ] 成就统计服务 (`pkg/utils/achievement_calculator.go`)
- [ ] 成就数据层 (`internal/repository/postgres/achievement_repo.go`)
- [ ] 成就触发器系统

#### API接口设计
```http
GET  /api/v1/achievements/:character_id           # 获取角色成就列表
GET  /api/v1/achievements/:character_id/unlocked  # 获取已解锁成就
GET  /api/v1/achievements/:character_id/progress  # 获取成就进度
GET  /api/v1/achievements/categories              # 获取成就分类
GET  /api/v1/achievements/leaderboard            # 获取成就排行榜
POST /api/v1/achievements/check/:character_id     # 检查成就触发
GET  /api/v1/achievements/stats/:character_id     # 获取成就统计
```

#### 数据模型
```go
type Achievement struct {
    AchievementID   string    `json:"achievement_id" db:"achievement_id"`
    Title           string    `json:"title" db:"title"`
    Description     string    `json:"description" db:"description"`
    Category        string    `json:"category" db:"category"`
    Rarity          string    `json:"rarity" db:"rarity"`           // common, rare, epic, legendary
    Points          int       `json:"points" db:"points"`
    Icon            string    `json:"icon" db:"icon"`
    Requirements    string    `json:"requirements" db:"requirements"` // JSON格式的条件
    IsHidden        bool      `json:"is_hidden" db:"is_hidden"`
    CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

type CharacterAchievement struct {
    CharacterID     string    `json:"character_id" db:"character_id"`
    AchievementID   string    `json:"achievement_id" db:"achievement_id"`
    UnlockedAt      time.Time `json:"unlocked_at" db:"unlocked_at"`
    Progress        float64   `json:"progress" db:"progress"`        // 0.0-1.0
    CurrentValue    int       `json:"current_value" db:"current_value"`
    TargetValue     int       `json:"target_value" db:"target_value"`
    IsCompleted     bool      `json:"is_completed" db:"is_completed"`
}

type AchievementCategory struct {
    CategoryID   string `json:"category_id" db:"category_id"`
    Name         string `json:"name" db:"name"`
    Description  string `json:"description" db:"description"`
    Icon         string `json:"icon" db:"icon"`
    SortOrder    int    `json:"sort_order" db:"sort_order"`
}
```

#### 成就分类系统
```go
const (
    // 职业成就
    AchievementCategoryCareer = "career"
    
    // 家庭成就
    AchievementCategoryFamily = "family"
    
    // 社会成就
    AchievementCategorySocial = "social"
    
    // 个人成就
    AchievementCategoryPersonal = "personal"
    
    // 财富成就
    AchievementCategoryWealth = "wealth"
    
    // 特殊成就
    AchievementCategorySpecial = "special"
)
```

#### 成就触发示例
```json
{
  "achievement_id": "first_job",
  "requirements": {
    "type": "event",
    "event_type": "career_start",
    "conditions": {
      "age": {"min": 16, "max": 30}
    }
  }
}
```

---

### 任务9: 存档与统计系统
- **分支**: `xucheng/feature/v0.3/save-stats-system`
- **负责人**: xucheng
- **预计时间**: 4-5天

#### 开发目标
- [ ] 游戏存档管理
- [ ] 数据统计接口
- [ ] 历史数据查询
- [ ] 导入导出功能

#### 详细任务清单
- [ ] 实现多存档位管理
- [ ] 实现自动存档功能
- [ ] 实现存档数据压缩和优化
- [ ] 实现云存档同步功能
- [ ] 实现详细数据统计分析
- [ ] 实现历史数据查询和筛选
- [ ] 实现数据导出（JSON/CSV）
- [ ] 实现数据导入和验证

#### 交付物
- [ ] 存档管理API (`internal/api/handlers/save.go`)
- [ ] 统计分析模块 (`internal/services/stats_service.go`)
- [ ] 数据导入导出 (`pkg/utils/data_export.go`)
- [ ] 历史数据服务 (`internal/services/history_service.go`)
- [ ] 存档数据层 (`internal/repository/postgres/save_repo.go`)
- [ ] 云存档同步机制

#### API接口设计
```http
# 存档管理
GET  /api/v1/saves/:character_id              # 获取存档列表
POST /api/v1/saves/:character_id              # 创建存档
GET  /api/v1/saves/:save_id                   # 获取存档详情
POST /api/v1/saves/:save_id/load              # 加载存档
DELETE /api/v1/saves/:save_id                 # 删除存档

# 统计分析
GET  /api/v1/stats/:character_id/overview     # 获取统计概览
GET  /api/v1/stats/:character_id/attributes   # 获取属性统计
GET  /api/v1/stats/:character_id/events       # 获取事件统计
GET  /api/v1/stats/:character_id/relationships # 获取关系统计
GET  /api/v1/stats/:character_id/timeline     # 获取时间线数据

# 数据导出
GET  /api/v1/export/:character_id/json        # 导出JSON格式
GET  /api/v1/export/:character_id/csv         # 导出CSV格式
POST /api/v1/import/:character_id             # 数据导入
```

#### 数据模型
```go
type GameSave struct {
    SaveID        string    `json:"save_id" db:"save_id"`
    CharacterID   string    `json:"character_id" db:"character_id"`
    SaveName      string    `json:"save_name" db:"save_name"`
    SaveType      string    `json:"save_type" db:"save_type"`      // auto, manual, milestone
    GameState     string    `json:"game_state" db:"game_state"`    // JSON压缩存储
    CurrentAge    int       `json:"current_age" db:"current_age"`
    LifeStage     string    `json:"life_stage" db:"life_stage"`
    Screenshot    string    `json:"screenshot" db:"screenshot"`    // 可选的截图
    Description   string    `json:"description" db:"description"`
    CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

type CharacterStats struct {
    CharacterID     string                 `json:"character_id"`
    TotalPlayTime   int                   `json:"total_play_time"`    // 分钟
    CurrentAge      int                   `json:"current_age"`
    TotalEvents     int                   `json:"total_events"`
    TotalDecisions  int                   `json:"total_decisions"`
    Achievements    int                   `json:"achievements"`
    AttributeStats  AttributeStatistics   `json:"attribute_stats"`
    RelationshipStats RelationshipStatistics `json:"relationship_stats"`
    CareerStats     CareerStatistics      `json:"career_stats"`
    LifeStages      []LifeStageStats      `json:"life_stages"`
}

type AttributeStatistics struct {
    Intelligence struct {
        Current int     `json:"current"`
        Peak    int     `json:"peak"`
        Average float64 `json:"average"`
    } `json:"intelligence"`
    
    Physical struct {
        Current int     `json:"current"`
        Peak    int     `json:"peak"`
        Average float64 `json:"average"`
    } `json:"physical"`
    
    // 其他属性统计...
}

type TimelineEvent struct {
    Age         int       `json:"age"`
    EventType   string    `json:"event_type"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    Impact      string    `json:"impact"`
    Timestamp   time.Time `json:"timestamp"`
}
```

#### 数据导出格式
```go
type ExportData struct {
    Character      Character              `json:"character"`
    GameHistory    []TimelineEvent        `json:"game_history"`
    Relationships  []Relationship         `json:"relationships"`
    Achievements   []CharacterAchievement `json:"achievements"`
    Statistics     CharacterStats         `json:"statistics"`
    ExportedAt     time.Time             `json:"exported_at"`
    Version        string                `json:"version"`
}
```

---

## 📊 阶段验收标准

### 功能验收
- [ ] 关系网络系统正常工作
- [ ] 关系质量评估准确
- [ ] 成就系统能正确解锁
- [ ] 成就进度追踪正常
- [ ] 存档系统稳定可靠
- [ ] 统计数据准确完整
- [ ] 数据导入导出功能正常

### 技术验收
- [ ] 关系网络算法效率高
- [ ] 成就触发逻辑准确
- [ ] 存档数据完整性保证
- [ ] 统计查询性能良好
- [ ] 数据导出格式规范
- [ ] 单元测试覆盖率 > 80%

### 性能验收
- [ ] 关系网络查询 < 200ms
- [ ] 成就检查响应 < 100ms
- [ ] 存档保存速度 < 1s
- [ ] 统计数据生成 < 500ms
- [ ] 数据导出效率高

---

## 🧪 测试规划

### 单元测试
```go
// 关系系统测试
- TestRelationshipCreation
- TestRelationshipQualityCalculation
- TestRelationshipEventProcessing
- TestRelationshipImpactCalculation

// 成就系统测试
- TestAchievementTrigger
- TestAchievementProgress
- TestAchievementUnlock
- TestAchievementStatistics

// 存档系统测试
- TestGameSaveCreation
- TestGameSaveLoading
- TestDataExport
- TestDataImport
```

### 集成测试
```go
// 系统集成测试
- TestRelationshipGameIntegration
- TestAchievementGameIntegration
- TestSaveLoadGameState
- TestStatisticsGeneration
```

### 压力测试
- 大量关系数据处理测试
- 频繁成就检查性能测试
- 批量存档操作测试
- 复杂统计查询压力测试

---

## 🔧 技术规范

### 新增依赖
```go
// 数据处理
- github.com/shopspring/decimal     // 精确数值计算
- github.com/go-sql-driver/mysql   // MySQL驱动
- github.com/klauspost/compress     // 数据压缩

// 数据导出
- github.com/gocarina/gocsv         // CSV处理
- encoding/json                     // JSON处理

// 图论算法
- github.com/yourbasic/graph        // 关系网络图算法
```

### 性能优化策略
```go
// 关系网络优化
- 使用Redis缓存热点关系数据
- 关系查询SQL索引优化
- 关系网络图算法优化

// 成就系统优化
- 成就触发条件预计算
- 批量成就检查机制
- 成就数据缓存策略

// 存档系统优化
- 增量存档机制
- 存档数据压缩
- 异步存档处理
```

### 数据库优化
```sql
-- 关系表索引
CREATE INDEX idx_relationships_character_id ON relationships(character_id);
CREATE INDEX idx_relationships_type ON relationships(relation_type);

-- 成就表索引
CREATE INDEX idx_character_achievements_character_id ON character_achievements(character_id);
CREATE INDEX idx_character_achievements_completed ON character_achievements(is_completed);

-- 存档表索引
CREATE INDEX idx_game_saves_character_id ON game_saves(character_id);
CREATE INDEX idx_game_saves_created_at ON game_saves(created_at);
```

---

## ⏭️ 下一步计划

完成Step 3后，将进入Step 4：性能优化与部署 (v0.4.x)
- 性能优化
- 监控与日志
- 部署与文档

---

*创建时间: 2025-01-26*
*最后更新: 2025-01-26* 