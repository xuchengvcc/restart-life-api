-- 创建角色决策表
CREATE TABLE IF NOT EXISTS game_decision (
    character_id VARCHAR(36) PRIMARY KEY,
    options JSON NOT NULL,
    previous_options JSON,
    previous_decision TINYINT COMMENT '1: conservative, 2: moderate, 3: aggressive',
    created_at BIGINT UNSIGNED NOT NULL COMMENT '创建时间（毫秒时间戳）',
    updated_at BIGINT UNSIGNED NOT NULL COMMENT '更新时间（毫秒时间戳）',
    FOREIGN KEY (character_id) REFERENCES character_tab(character_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
