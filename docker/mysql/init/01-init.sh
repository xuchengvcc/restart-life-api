#!/bin/bash
# MySQL初始化脚本

set -e

echo "Starting MySQL initialization..."

# 等待MySQL服务启动
while ! mysqladmin ping -h"localhost" --silent; do
    echo "Waiting for MySQL to start..."
    sleep 2
done

echo "MySQL is ready!"

# 创建额外的数据库和用户（如果需要）
mysql -u root -p"$MYSQL_ROOT_PASSWORD" <<-EOSQL
    -- 创建测试数据库
    CREATE DATABASE IF NOT EXISTS restart_life_test;
    
    -- 为应用用户授权
    GRANT ALL PRIVILEGES ON restart_life_dev.* TO 'restart_user'@'%';
    GRANT ALL PRIVILEGES ON restart_life_test.* TO 'restart_user'@'%';
    
    -- 刷新权限
    FLUSH PRIVILEGES;
    
    -- 显示数据库列表
    SHOW DATABASES;
EOSQL

echo "MySQL initialization completed!"
