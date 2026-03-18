#!/bin/bash
#
# fdim 一键停止所有服务
# 用法: ./scripts/stop_all.sh
#

PID_DIR="/tmp"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

info()  { echo -e "${GREEN}[INFO]${NC}  $1"; }
warn()  { echo -e "${YELLOW}[WARN]${NC}  $1"; }

# 所有服务（按启动的反序停止）
SERVICES=(
    "api"
    "cron"
    "msggateway"
    "push"
    "msgtransfer"
    "bot"
    "third"
    "msg"
    "conversation"
    "chat"
    "admin"
    "auth"
    "user"
)

info "========== 停止所有 fdim 服务 =========="

for name in "${SERVICES[@]}"; do
    pid_file="$PID_DIR/fdim_${name}.pid"
    if [ -f "$pid_file" ]; then
        pid=$(cat "$pid_file")
        # go run 会产生子进程，需要杀整个进程组
        if kill -0 "$pid" 2>/dev/null; then
            kill -- -"$pid" 2>/dev/null || kill "$pid" 2>/dev/null
            info "$name 已停止 (PID: $pid)"
        else
            warn "$name 父进程已退出 (PID: $pid)"
        fi
        rm -f "$pid_file"
    fi
done

# 兜底：杀掉所有残留的 fdim 相关进程
REMAINING=$(pgrep -f "fdim/(rpc|api|Infrastructure_service)" 2>/dev/null || true)
if [ -n "$REMAINING" ]; then
    warn "清理残留进程..."
    echo "$REMAINING" | xargs kill 2>/dev/null || true
    sleep 1
    # 如果还没死，强制杀
    STILL=$(pgrep -f "fdim/(rpc|api|Infrastructure_service)" 2>/dev/null || true)
    if [ -n "$STILL" ]; then
        echo "$STILL" | xargs kill -9 2>/dev/null || true
        warn "已强制终止残留进程"
    fi
fi

info "========== 全部服务已停止 =========="
