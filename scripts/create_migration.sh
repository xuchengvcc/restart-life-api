#!/bin/bash

# 创建新的数据库迁移文件
set -e

# 颜色输出
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

log_info() {
    echo -e "${YELLOW}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查参数
if [ $# -eq 0 ]; then
    log_error "请提供迁移名称"
    echo ""
    echo "用法: $0 <migration_name>"
    echo ""
    echo "示例:"
    echo "  $0 add_user_preferences_table"
    echo "  $0 add_character_skills_column"
    echo "  $0 create_game_events_table"
    exit 1
fi

MIGRATION_NAME=$1
MIGRATIONS_DIR="migrations"
TIMESTAMP=$(date +%s)

# 确保迁移目录存在
mkdir -p "$MIGRATIONS_DIR"

# 生成文件名（使用时间戳作为版本号）
UP_FILE="${MIGRATIONS_DIR}/${TIMESTAMP}_${MIGRATION_NAME}.up.sql"
DOWN_FILE="${MIGRATIONS_DIR}/${TIMESTAMP}_${MIGRATION_NAME}.down.sql"

# 创建up迁移文件
cat > "$UP_FILE" << EOF
-- Migration: $MIGRATION_NAME
-- Created at: $(date)
-- Description: TODO: 添加迁移描述

-- TODO: 在此添加数据库变更SQL语句
-- 示例:
-- CREATE TABLE example_table (
--     id INT AUTO_INCREMENT PRIMARY KEY,
--     name VARCHAR(255) NOT NULL,
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
-- );

EOF

# 创建down迁移文件
cat > "$DOWN_FILE" << EOF
-- Rollback for migration: $MIGRATION_NAME
-- Created at: $(date)

-- TODO: 在此添加回滚SQL语句
-- 示例:
-- DROP TABLE IF EXISTS example_table;

EOF

log_success "迁移文件已创建："
echo "  UP:   $UP_FILE"
echo "  DOWN: $DOWN_FILE"
echo ""
log_info "请编辑这些文件并添加相应的SQL语句"
echo ""
echo "下一步："
echo "1. 编辑 $UP_FILE 添加数据库变更"
echo "2. 编辑 $DOWN_FILE 添加回滚逻辑"
echo "3. 运行 scripts/init_db.sh 应用迁移"
