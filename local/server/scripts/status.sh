#!/bin/bash
#
# fdim 查看所有服务状态
# 用法: ./scripts/status.sh
#

PID_DIR="/tmp"

RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

SERVICES=(
    "user" "auth" "admin" "chat" "conversation"
    "msg" "third" "bot"
    "msgtransfer" "push" "msggateway" "cron"
    "api"
)

printf "%-16s %-8s %-6s\n" "SERVICE" "STATUS" "PID"
printf "%-16s %-8s %-6s\n" "-------" "------" "---"

for name in "${SERVICES[@]}"; do
    pid_file="$PID_DIR/fdim_${name}.pid"
    if [ -f "$pid_file" ]; then
        pid=$(cat "$pid_file")
        if kill -0 "$pid" 2>/dev/null; then
            printf "%-16s ${GREEN}%-8s${NC} %-6s\n" "$name" "运行中" "$pid"
        else
            printf "%-16s ${RED}%-8s${NC} %-6s\n" "$name" "已退出" "$pid"
        fi
    else
        printf "%-16s ${RED}%-8s${NC} %-6s\n" "$name" "未启动" "-"
    fi
done
