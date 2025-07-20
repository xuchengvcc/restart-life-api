#!/bin/bash

# 游戏核心逻辑API测试脚本
# 测试游戏推进、决策选择等功能

set -e

# API基础URL
BASE_URL="http://localhost:8080/api/v1"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 全局变量
ACCESS_TOKEN=""
CHARACTER_ID=""
DECISION_ID=""

# 打印带颜色的消息
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查JSON响应中的success字段
check_response() {
    local response="$1"
    local operation="$2"

    if echo "$response" | jq -e '.success == true' > /dev/null; then
        print_success "$operation"
        return 0
    else
        print_error "$operation failed"
        echo "$response" | jq '.'
        return 1
    fi
}

# 用户注册
register_user() {
    print_info "注册测试用户..."

    local username="gametest_$(date +%s)"
    local email="gametest_$(date +%s)@example.com"
    local password="GameTest123!"

    local response=$(curl -s -X POST "$BASE_URL/auth/register" \
        -H "Content-Type: application/json" \
        -d "{
            \"username\": \"$username\",
            \"email\": \"$email\",
            \"password\": \"$password\"
        }")

    check_response "$response" "用户注册"

    # 登录获取token
    print_info "用户登录..."
    local login_response=$(curl -s -X POST "$BASE_URL/auth/login" \
        -H "Content-Type: application/json" \
        -d "{
            \"username\": \"$username\",
            \"password\": \"$password\"
        }")

    if check_response "$login_response" "用户登录"; then
        ACCESS_TOKEN=$(echo "$login_response" | jq -r '.data.access_token')
        print_success "获取访问令牌: ${ACCESS_TOKEN:0:20}..."
    else
        return 1
    fi
}

# 创建角色
create_character() {
    print_info "创建测试角色..."

    local response=$(curl -s -X POST "$BASE_URL/characters/create" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -d '{
            "character_name": "游戏测试角色",
            "birth_country": "中国",
            "birth_year": 2000,
            "gender": 1,
            "race": 2
        }')

    if check_response "$response" "角色创建"; then
        CHARACTER_ID=$(echo "$response" | jq -r '.data.character_id')
        print_success "角色ID: $CHARACTER_ID"

        # 显示角色属性
        print_info "角色属性:"
        echo "$response" | jq '.data.attributes'
    else
        return 1
    fi
}

# 开始游戏
start_game() {
    print_info "开始游戏..."

    local response=$(curl -s -X POST "$BASE_URL/game/start/$CHARACTER_ID" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $ACCESS_TOKEN")

    if check_response "$response" "开始游戏"; then
        print_success "游戏开始成功"
        print_info "初始游戏状态:"
        echo "$response" | jq '.data | {current_age, life_stage, education, career, location}'

        print_info "出生事件:"
        echo "$response" | jq '.data.current_events[0] | {title, description}'
    else
        return 1
    fi
}

# 获取游戏状态
get_game_state() {
    print_info "获取游戏状态..."

    local response=$(curl -s -X GET "$BASE_URL/game/state/$CHARACTER_ID" \
        -H "Authorization: Bearer $ACCESS_TOKEN")

    if check_response "$response" "获取游戏状态"; then
        print_info "当前游戏状态:"
        echo "$response" | jq '.data | {current_age, life_stage, education, career, location, is_game_active}'
    else
        return 1
    fi
}

# 推进游戏
advance_game() {
    local target_age="$1"
    print_info "推进游戏到 ${target_age} 岁..."

    local response=$(curl -s -X POST "$BASE_URL/game/advance/$CHARACTER_ID" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $ACCESS_TOKEN")

    if check_response "$response" "推进游戏"; then
        local current_age=$(echo "$response" | jq -r '.data.game_state.current_age')
        local life_stage=$(echo "$response" | jq -r '.data.game_state.life_stage')

        print_success "推进到 ${current_age} 岁 (${life_stage})"

        # 显示事件
        local event_count=$(echo "$response" | jq '.data.game_state.current_events | length')
        if [ "$event_count" -gt 0 ]; then
            print_info "最新事件:"
            echo "$response" | jq '.data.game_state.current_events[-1] | {title, description}'
        fi

        # 检查是否有决策
        local has_decision=$(echo "$response" | jq '.data.decision != null')
        if [ "$has_decision" = "true" ]; then
            print_warning "有待处理的决策!"
            DECISION_ID=$(echo "$response" | jq -r '.data.decision.decision_id')

            print_info "决策问题:"
            echo "$response" | jq '.data.decision | {question, context}'

            print_info "决策选项:"
            echo "$response" | jq '.data.decision.options[]'

            return 2  # 表示有决策待处理
        fi
    else
        return 1
    fi
}

# 做出决策
make_decision() {
    local option_type="$1"
    print_info "做出决策: $option_type"

    local response=$(curl -s -X POST "$BASE_URL/game/decision/$CHARACTER_ID" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -d "{
            \"decision_id\": \"$DECISION_ID\",
            \"option_type\": \"$option_type\"
        }")

    if check_response "$response" "做出决策"; then
        print_success "决策处理完成"

        print_info "决策结果:"
        echo "$response" | jq '.data.decision_result | {result_description, life_impact}'

        # 显示属性变化
        local attr_changes=$(echo "$response" | jq '.data.decision_result.attribute_changes')
        if [ "$attr_changes" != "null" ] && [ "$attr_changes" != "{}" ]; then
            print_info "属性变化:"
            echo "$attr_changes" | jq '.'
        fi

        # 显示状态变化
        local status_changes=$(echo "$response" | jq '.data.decision_result.status_changes')
        if [ "$status_changes" != "null" ] && [ "$status_changes" != "{}" ]; then
            print_info "状态变化:"
            echo "$status_changes" | jq '.'
        fi

        # 清除决策ID
        DECISION_ID=""
    else
        return 1
    fi
}

# 保存游戏
save_game() {
    print_info "保存游戏..."

    local response=$(curl -s -X POST "$BASE_URL/game/save/$CHARACTER_ID" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $ACCESS_TOKEN")

    check_response "$response" "保存游戏"
}

# 加载游戏
load_game() {
    print_info "加载游戏..."

    local response=$(curl -s -X POST "$BASE_URL/game/load/$CHARACTER_ID" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $ACCESS_TOKEN")

    if check_response "$response" "加载游戏"; then
        print_info "加载的游戏状态:"
        echo "$response" | jq '.data | {current_age, life_stage, total_playtime}'
    else
        return 1
    fi
}

# 获取事件历史
get_event_history() {
    print_info "获取事件历史..."

    local response=$(curl -s -X GET "$BASE_URL/game/events/$CHARACTER_ID?page=1&limit=10" \
        -H "Authorization: Bearer $ACCESS_TOKEN")

    if check_response "$response" "获取事件历史"; then
        local event_count=$(echo "$response" | jq '.data | length')
        print_info "事件历史数量: $event_count"

        if [ "$event_count" -gt 0 ]; then
            print_info "最近的事件:"
            echo "$response" | jq '.data[] | {age, title, description}' | head -20
        fi
    else
        return 1
    fi
}

# 模拟完整的游戏流程
simulate_game_flow() {
    print_info "开始模拟完整的游戏流程..."

    # 推进到不同年龄并处理决策
    local ages=(5 10 16 18 22 25)

    for age in "${ages[@]}"; do
        echo ""
        print_info "=== 推进到 $age 岁 ==="

        local result
        advance_game "$age"
        result=$?

        if [ $result -eq 2 ]; then
            # 有决策，随机选择一个选项
            local options=("conservative" "moderate" "aggressive")
            local selected_option=${options[$((RANDOM % ${#options[@]}))]}

            print_warning "自动选择决策选项: $selected_option"
            make_decision "$selected_option"

            # 再次推进到目标年龄
            if [ $? -eq 0 ]; then
                advance_game "$age"
            fi
        elif [ $result -ne 0 ]; then
            print_error "推进游戏失败"
            return 1
        fi

        sleep 1  # 避免请求过快
    done
}

# 主函数
main() {
    print_info "开始游戏核心逻辑API测试"
    echo "================================"

    # 检查依赖
    if ! command -v jq &> /dev/null; then
        print_error "需要安装 jq 工具来解析JSON响应"
        exit 1
    fi

    # 检查服务器是否运行
    if ! curl -s "$BASE_URL/health" > /dev/null; then
        print_error "API服务器未运行或无法访问: $BASE_URL"
        exit 1
    fi

    # 执行测试流程
    register_user || exit 1
    echo ""

    create_character || exit 1
    echo ""

    start_game || exit 1
    echo ""

    get_game_state || exit 1
    echo ""

    simulate_game_flow || exit 1
    echo ""

    save_game || exit 1
    echo ""

    load_game || exit 1
    echo ""

    get_event_history || exit 1
    echo ""

    print_success "游戏核心逻辑API测试完成!"
    print_info "角色ID: $CHARACTER_ID"
    print_info "可以在游戏中继续使用这个角色进行测试"
}

# 运行主函数
main "$@"
