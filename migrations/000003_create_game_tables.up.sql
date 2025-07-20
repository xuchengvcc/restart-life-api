-- 创建游戏事件表
CREATE TABLE IF NOT EXISTS game_events (
    event_id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    character_id VARCHAR(36) NOT NULL,
    age INT NOT NULL,
    description TEXT NOT NULL,
    impact VARCHAR(255) NOT NULL COMMENT '事件影响简要描述',
    created_at BIGINT UNSIGNED NOT NULL COMMENT '创建时间（毫秒时间戳）',
    INDEX idx_event_character_id (character_id),
    INDEX idx_event_character_age (character_id, age)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
