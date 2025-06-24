#!/bin/bash

# 测试认证API的脚本

BASE_URL="http://localhost:8080/api/v1"

echo "=== 测试用户认证系统 ==="

# 1. 测试用户注册
echo "1. 测试用户注册..."
REGISTER_RESPONSE=$(curl -s -X POST \
  "${BASE_URL}/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }')

echo "注册响应: $REGISTER_RESPONSE"

# 提取访问Token（需要jq工具）
if command -v jq &> /dev/null; then
    ACCESS_TOKEN=$(echo $REGISTER_RESPONSE | jq -r '.data.access_token')
    echo "访问Token: $ACCESS_TOKEN"
else
    echo "请安装jq工具来解析JSON响应"
    ACCESS_TOKEN="your_token_here"
fi

echo ""

# 2. 测试用户登录
echo "2. 测试用户登录..."
LOGIN_RESPONSE=$(curl -s -X POST \
  "${BASE_URL}/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }')

echo "登录响应: $LOGIN_RESPONSE"

# 如果有jq，更新token
if command -v jq &> /dev/null; then
    ACCESS_TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.data.access_token')
fi

echo ""

# 3. 测试获取用户信息
echo "3. 测试获取用户信息..."
PROFILE_RESPONSE=$(curl -s -X GET \
  "${BASE_URL}/auth/profile" \
  -H "Authorization: Bearer $ACCESS_TOKEN")

echo "用户信息响应: $PROFILE_RESPONSE"

echo ""

# 4. 测试更新用户信息
echo "4. 测试更新用户信息..."
UPDATE_RESPONSE=$(curl -s -X PUT \
  "${BASE_URL}/auth/profile" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "bio": "这是我的个人简介",
    "country": "中国"
  }')

echo "更新用户信息响应: $UPDATE_RESPONSE"

echo ""

# 5. 测试修改密码
echo "5. 测试修改密码..."
CHANGE_PASSWORD_RESPONSE=$(curl -s -X POST \
  "${BASE_URL}/auth/change-password" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "old_password": "password123",
    "new_password": "newpassword123"
  }')

echo "修改密码响应: $CHANGE_PASSWORD_RESPONSE"

echo ""

# 6. 测试用新密码登录
echo "6. 测试用新密码登录..."
NEW_LOGIN_RESPONSE=$(curl -s -X POST \
  "${BASE_URL}/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "newpassword123"
  }')

echo "新密码登录响应: $NEW_LOGIN_RESPONSE"

echo ""

# 7. 测试登出
echo "7. 测试登出..."
LOGOUT_RESPONSE=$(curl -s -X POST \
  "${BASE_URL}/auth/logout" \
  -H "Authorization: Bearer $ACCESS_TOKEN")

echo "登出响应: $LOGOUT_RESPONSE"

echo ""
echo "=== 测试完成 ==="
