#!/bin/bash

# 数据库初始化脚本
set -e

echo "=== 初始化Restart Life数据库 ==="

# 配置变量
DB_HOST=${DB_HOST:-"localhost"}
DB_PORT=${DB_PORT:-3306}
DB_USER=${DB_USER:-"root"}
DB_PASSWORD=${DB_PASSWORD:-"password"}
DB_NAME=${DB_NAME:-"restart_life_dev"}

# 检查MySQL是否可连接
echo "检查MySQL连接..."
if ! command -v mysql &> /dev/null; then
    echo "错误: 请先安装MySQL客户端"
    exit 1
fi

# 测试连接
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" -e "SELECT 1;" > /dev/null 2>&1
if [ $? -ne 0 ]; then
    echo "错误: 无法连接到MySQL服务器"
    echo "请检查连接配置: $DB_USER@$DB_HOST:$DB_PORT"
    exit 1
fi

echo "MySQL连接成功"

# 创建数据库（如果不存在）
echo "创建数据库: $DB_NAME"
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" -e "CREATE DATABASE IF NOT EXISTS \`$DB_NAME\` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 运行迁移脚本
echo "运行数据库迁移..."

# 创建用户表
echo "创建用户表..."
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" < migrations/000001_create_users_table.up.sql

# 创建角色表
echo "创建角色表..."
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" < migrations/000002_create_characters_table.up.sql

echo "=== 数据库初始化完成 ==="

# 验证表是否创建成功
echo "验证表结构..."
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" -e "SHOW TABLES;"

echo ""
echo "数据库配置："
echo "  主机: $DB_HOST:$DB_PORT"
echo "  数据库: $DB_NAME"
echo "  用户: $DB_USER"
echo ""
echo "可以开始运行应用程序了！"
