-- 创建用户账户表
CREATE TABLE IF NOT EXISTS user_tab (
    user_id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '用户ID',
    username VARCHAR(64) UNIQUE NOT NULL COMMENT '用户名',
    email VARCHAR(128) UNIQUE NOT NULL COMMENT '邮箱',
    password_hash VARCHAR(255) NOT NULL COMMENT '密码哈希',
    created_at BIGINT UNSIGNED NOT NULL COMMENT '毫秒时间戳',
    updated_at BIGINT UNSIGNED NOT NULL COMMENT '毫秒时间戳',
    last_login BIGINT UNSIGNED NULL COMMENT '最后登录时间 毫秒时间戳',
    is_active BOOLEAN DEFAULT TRUE COMMENT '是否激活 1:正常 0:禁用/未激活',
    avatar_url VARCHAR(500) COMMENT '头像URL',
    bio TEXT COMMENT '个人简介',
    birth_date DATE COMMENT '出生日期',
    gender TINYINT COMMENT '性别 0:未知 1:男 2:女 3:其他',
    country VARCHAR(100) COMMENT '国家'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户账户表';

-- 注意：username 和 email 字段的 UNIQUE 约束会自动创建索引，无需手动创建