# Achievement & Stats API (T-007)

## Scope
This document describes the four newly implemented endpoints that replaced placeholder handlers:

1. `GET /api/v1/achievements/categories`
2. `GET /api/v1/achievements/:character_id`
3. `GET /api/v1/stats/:character_id`
4. `GET /api/v1/stats/:character_id/timeline`

All endpoints require `Authorization: Bearer <token>`.

## Common Response Envelope
```json
{
  "success": true,
  "data": {}
}
```

Failure:
```json
{
  "success": false,
  "error": {
    "code": 1000,
    "message": "error message",
    "details": "optional"
  }
}
```

## 1) Achievement Categories
`GET /api/v1/achievements/categories`

Returns user-level category summary aggregated from owned characters.

### Response `data`
```json
[
  {
    "id": "milestone",
    "name": "里程碑",
    "description": "年龄与人生阶段关键节点",
    "total_count": 4,
    "unlocked_count": 2
  }
]
```

## 2) Character Achievements
`GET /api/v1/achievements/:character_id`

Returns achievement list for one character. Ownership is enforced.

### Path Params
- `character_id` (string, required)

### Response `data`
```json
{
  "character_id": "char_1",
  "character_name": "tester",
  "items": [
    {
      "id": "age_18",
      "title": "成年礼",
      "description": "角色达到 18 岁",
      "category": "milestone",
      "unlocked": true,
      "progress": 18,
      "max_progress": 18,
      "unlocked_at": 1744953600000
    }
  ],
  "unlocked_count": 1,
  "total_count": 9
}
```

## 3) Character Stats
`GET /api/v1/stats/:character_id`

Returns current game/statistics snapshot for one character. Ownership is enforced.

### Path Params
- `character_id` (string, required)

### Response `data`
```json
{
  "character_id": "char_1",
  "character_name": "tester",
  "current_age": 22,
  "life_stage": "young_adult",
  "is_game_active": true,
  "total_playtime": 120,
  "money": 15000,
  "event_count": 6,
  "pending_decision": false,
  "attributes": {
    "intelligence": 60
  },
  "last_event": {
    "age": 22,
    "description": "event",
    "impact": "impact"
  }
}
```

## 4) Character Timeline
`GET /api/v1/stats/:character_id/timeline`

Returns timeline projection of event history for one character.

### Path Params
- `character_id` (string, required)

### Response `data`
```json
{
  "character_id": "char_1",
  "timeline": [
    {
      "age": 18,
      "description": "进入大学",
      "impact": "learning",
      "created_at": 1744953600000
    }
  ],
  "total": 1
}
```

## Compatibility Notes
- Legacy game endpoints used by `.frontend` are aliased in routes:
  - `GET /api/v1/game/history/:character_id` -> event history
  - `GET /api/v1/game/load/:character_id` -> load game
  - `POST /api/v1/game/next-turn/:character_id` -> advance game
  - `POST /api/v1/game/decision/:character_id` -> legacy decision adapter
- Legacy decision adapter accepts either `option_type` or `decision_id` and maps values to:
  - `conservative`
  - `moderate`
  - `aggressive`
