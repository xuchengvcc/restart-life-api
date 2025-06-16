-- 创建游戏角色表
CREATE TABLE IF NOT EXISTS characters (
    character_id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    user_id CHAR(36) NOT NULL,
    character_name VARCHAR(100) NOT NULL,
    birth_country VARCHAR(100) NOT NULL,
    birth_year INTEGER NOT NULL,
    current_age INTEGER NOT NULL DEFAULT 0,
    gender VARCHAR(20) NOT NULL,
    race VARCHAR(50) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    -- 角色属性
    intelligence INTEGER DEFAULT 50 CHECK (intelligence >= 0 AND intelligence <= 100),
    emotional_intelligence INTEGER DEFAULT 50 CHECK (emotional_intelligence >= 0 AND emotional_intelligence <= 100),
    memory INTEGER DEFAULT 50 CHECK (memory >= 0 AND memory <= 100),
    imagination INTEGER DEFAULT 50 CHECK (imagination >= 0 AND imagination <= 100),
    physical_fitness INTEGER DEFAULT 50 CHECK (physical_fitness >= 0 AND physical_fitness <= 100),
    appearance INTEGER DEFAULT 50 CHECK (appearance >= 0 AND appearance <= 100),
    
    -- 游戏状态
    life_stage VARCHAR(50) DEFAULT 'birth',
    current_status VARCHAR(100) DEFAULT 'healthy',
    happiness_level INTEGER DEFAULT 50 CHECK (happiness_level >= 0 AND happiness_level <= 100),
    health_level INTEGER DEFAULT 100 CHECK (health_level >= 0 AND health_level <= 100),
    money INTEGER DEFAULT 0 CHECK (money >= 0),
    
    -- 当前位置和活动
    current_location VARCHAR(200),
    current_activity VARCHAR(200),
    
    -- 游戏进度
    total_playtime INTEGER DEFAULT 0 COMMENT '总游戏时间（分钟）',
    game_completed BOOLEAN DEFAULT FALSE,
    final_age INTEGER,
    death_cause VARCHAR(200),
    
    -- 外键约束
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建索引
CREATE INDEX idx_characters_user_id ON characters(user_id);
CREATE INDEX idx_characters_is_active ON characters(is_active);
CREATE INDEX idx_characters_created_at ON characters(created_at);
CREATE INDEX idx_characters_life_stage ON characters(life_stage);
CREATE INDEX idx_characters_game_completed ON characters(game_completed); 