-- 创建游戏角色表
CREATE TABLE IF NOT EXISTS character_tab (
    character_id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    user_id INT UNSIGNED NOT NULL,
    character_name VARCHAR(100) NOT NULL,
    birth_country VARCHAR(100) NOT NULL,
    birth_year INT NOT NULL,
    current_age INT NOT NULL DEFAULT 0,
    gender TINYINT COMMENT '性别 0:未知 1:男 2:女 3:其他',
    race TINYINT COMMENT '种族类型 如 0:未知 1:白人 2:黄种人 3:黑人 ...',
    is_active BOOLEAN DEFAULT TRUE,
    created_at BIGINT UNSIGNED NOT NULL COMMENT '创建时间（毫秒时间戳）',
    updated_at BIGINT UNSIGNED NOT NULL COMMENT '更新时间（毫秒时间戳）',
    
    -- 角色属性
    intelligence INT DEFAULT 50 CHECK (intelligence >= 0 AND intelligence <= 100) COMMENT '智力，学习、推理、解决问题能力',
    emotional_intelligence INT DEFAULT 50 CHECK (emotional_intelligence >= 0 AND emotional_intelligence <= 100) COMMENT '情商，理解和管理情绪的能力',
    memory INT DEFAULT 50 CHECK (memory >= 0 AND memory <= 100) COMMENT '记忆力，学习和回忆信息的能力',
    imagination INT DEFAULT 50 CHECK (imagination >= 0 AND imagination <= 100) COMMENT '想象力，创新、创造、艺术表现',
    physical_fitness INT DEFAULT 50 CHECK (physical_fitness >= 0 AND physical_fitness <= 100) COMMENT '体质，健康、耐力、力量等身体素质',
    appearance INT DEFAULT 50 CHECK (appearance >= 0 AND appearance <= 100) COMMENT '外貌，吸引力、社交影响',
    
    -- 角色经历描述
    career_desc TEXT COMMENT '职业履历描述',
    education_desc TEXT COMMENT '教育经历描述',
    
    -- 角色扩展信息
    education_level VARCHAR(100) COMMENT '最高学历',
    marital_status VARCHAR(50) COMMENT '婚姻状况',
    current_country VARCHAR(100) COMMENT '现居国家',
    current_location VARCHAR(150) COMMENT '地理位置或场所。例如 “北京·清华大学',
    current_activity VARCHAR(200) COMMENT '描述角色此刻的主要行为状态，便于事件触发、属性变化、剧情推进等',
    personality VARCHAR(100) COMMENT '性格特征描述',
    career VARCHAR(100) COMMENT '当前职业描述',
    skill_tendency VARCHAR(100) COMMENT '技能倾向描述',
    family_background TEXT COMMENT '家庭背景描述',
    social_relationships TEXT COMMENT '社会关系描述',
    
    -- 游戏状态
    life_stage VARCHAR(50) DEFAULT 'birth' COMMENT '当前人生阶段，如 birth, childhood, adolescence, adulthood, old_age',
    current_status VARCHAR(100) DEFAULT 'healthy' COMMENT '当前状态，如 healthy, sick, injured, tired, deceased, excited',
    happiness_level INT DEFAULT 50 CHECK (happiness_level >= 0 AND happiness_level <= 100),
    health_level INT DEFAULT 100 CHECK (health_level >= 0 AND health_level <= 100),
    money BIGINT DEFAULT 0 COMMENT '资产（可为负，表示负债）',
    
    -- 游戏进度
    total_playtime INT DEFAULT 0 COMMENT '总游戏时间（分钟）',
    game_completed BOOLEAN DEFAULT FALSE,
    final_age INT,
    death_cause VARCHAR(200),
    
    -- 外键约束和索引
    FOREIGN KEY (user_id) REFERENCES user_tab(user_id) ON DELETE CASCADE,
    INDEX idx_character_user_id (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;