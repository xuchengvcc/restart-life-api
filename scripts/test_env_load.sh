#!/bin/bash

# 测试环境变量加载脚本
echo "=== 测试 .env.test 环境变量加载 ==="

# 方法1：直接source加载
echo "方法1: 使用 source 加载"
set -a
source .env.test
set +a

echo "DATABASE_MYSQL_HOST: $DATABASE_MYSQL_HOST"
echo "DATABASE_MYSQL_DATABASE: $DATABASE_MYSQL_DATABASE"
echo "REDIS_HOST: $REDIS_HOST"
echo "LOG_LEVEL: $LOG_LEVEL"

echo ""
echo "方法2: 导出所有变量供子进程使用"
export $(cat .env.test | grep -v '^#' | grep -v '^$' | xargs)

echo "已导出的测试环境变量："
env | grep -E "(DATABASE|REDIS|SERVER_MODE|CONFIG_ENV|LOG_LEVEL)" | sort

echo ""
echo "=== 测试完成 ==="
