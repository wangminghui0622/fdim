#!/bin/bash
#
# fdim 一键编译+启动所有服务
# 用法: ./scripts/start_all.sh
# 停止: ./scripts/stop_all.sh
#

set -e

# 项目根目录（脚本所在目录的上一级）
PROJECT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
CONFIG_DIR="$PROJECT_DIR/etc"
LOG_DIR="$PROJECT_DIR/logs"
PID_DIR="/tmp"

mkdir -p "$LOG_DIR"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

info()  { echo -e "${GREEN}[INFO]${NC}  $1"; }
warn()  { echo -e "${YELLOW}[WARN]${NC}  $1"; }
error() { echo -e "${RED}[ERROR]${NC} $1"; }

# 等待端口就绪
# 参数: $1=服务名 $2=端口 $3=超时秒数(默认30)
wait_for_port() {
    local name=$1
    local port=$2
    local timeout=${3:-30}
    local elapsed=0

    while [ $elapsed -lt $timeout ]; do
        if ss -tlnp | grep -q ":${port} " 2>/dev/null; then
            info "$name (:$port) 就绪 ✓"
            return 0
        fi
        sleep 1
        elapsed=$((elapsed + 1))
    done
    error "$name (:$port) 启动超时（${timeout}秒），查看日志: $LOG_DIR/${name}.log"
    return 1
}

# 编译单个服务
# 参数: $1=服务名 $2=go源文件路径
build_service() {
    local name=$1
    local src=$2
    info "编译 $name..."
    go build -o "$BIN_DIR/$name" "$src" || { error "编译 $name 失败"; exit 1; }
}

# 启动单个服务（编译后的二进制）
# 参数: $1=服务名 $2=配置文件 $3=端口(可选，用于等待就绪)
start_service() {
    local name=$1
    local config=$2
    local port=$3
    local pid_file="$PID_DIR/fdim_${name}.pid"
    local log_file="$LOG_DIR/${name}.log"
    local bin="$BIN_DIR/$name"

    # 检查是否已经在运行
    if [ -f "$pid_file" ] && kill -0 "$(cat "$pid_file")" 2>/dev/null; then
        warn "$name 已在运行 (PID: $(cat "$pid_file"))"
        return
    fi

    # 启动二进制
    nohup "$bin" -f "$config" > "$log_file" 2>&1 &
    local pid=$!
    echo $pid > "$pid_file"

    # 如果指定了端口，等待端口就绪；否则等 2 秒后检查进程
    if [ -n "$port" ]; then
        wait_for_port "$name" "$port" 30 || return 1
    else
        sleep 2
        if kill -0 $pid 2>/dev/null; then
            info "$name 启动成功 (PID: $pid)"
        else
            error "$name 启动失败，查看日志: $log_file"
            return 1
        fi
    fi
}

cd "$PROJECT_DIR"
BIN_DIR="$PROJECT_DIR/bin"
mkdir -p "$BIN_DIR"

# ========== 编译所有服务 ==========
info "========== 编译所有服务 =========="

build_service "user"         "./rpc/user/"
build_service "auth"         "./rpc/auth/"
build_service "msg"          "./rpc/msg/"
build_service "third"        "./rpc/third/"
build_service "admin"        "./rpc/admin/"
build_service "conversation" "./rpc/conversation/"
build_service "chat"         "./rpc/chat/"
build_service "bot"          "./rpc/bot/"
build_service "msgtransfer"  "./Infrastructure_service/msgtransfer/"
build_service "push"         "./Infrastructure_service/push/"
build_service "msggateway"   "./Infrastructure_service/msggateway/"
build_service "cron"         "./Infrastructure_service/cron/"
build_service "api"          "./api/"

info "========== 编译完成 =========="

# ========== 检查基础设施 ==========
info "========== 检查基础设施 =========="

check_port() {
    local name=$1
    local port=$2
    if ss -tlnp | grep -q ":${port} " 2>/dev/null; then
        info "$name (:$port) ✓"
    else
        error "$name (:$port) 未启动！"
        return 1
    fi
}

INFRA_OK=true
check_port "Redis"   6379  || INFRA_OK=false
check_port "MongoDB" 27017 || INFRA_OK=false
check_port "NATS"    4222  || INFRA_OK=false
check_port "Etcd"    2379  || INFRA_OK=false

if [ "$INFRA_OK" = false ]; then
    error "基础设施未就绪，请先启动 Redis、MongoDB、NATS、Etcd"
    exit 1
fi

info "========== 基础设施就绪 =========="

# ========== 1. 无 RPC 依赖的基础服务 ==========
info "========== [Layer 1] 启动基础 RPC 服务（无 RPC 依赖）=========="

start_service "auth"   "$CONFIG_DIR/auth.yaml"          8081  # 无依赖
start_service "third"  "$CONFIG_DIR/third.yaml"         8085  # 无依赖
start_service "admin"  "$CONFIG_DIR/admin.yaml"         8087  # 无依赖

# ========== 2. 启动 user 和 conversation（msg 依赖 conversation）==========
info "========== [Layer 2] 启动 user 和 conversation 服务 =========="

start_service "user"         "$CONFIG_DIR/user.yaml"         8080  # conversation 需要它
start_service "conversation" "$CONFIG_DIR/conversation.yaml" 8083  # msg 需要它

# 等待 conversation 完全注册到 etcd（重要！）
info "等待 conversation 注册到 etcd..."
sleep 3

# ========== 3. 启动 msg（依赖 conversation）==========
info "========== [Layer 3] 启动 msg 服务 =========="

start_service "msg"          "$CONFIG_DIR/msg.yaml"          8082  # 依赖 conversation

# ========== 4. 启动其他 RPC 服务 ==========
info "========== [Layer 4] 启动其他 RPC 服务 =========="

start_service "chat"         "$CONFIG_DIR/chat.yaml"         8084  # 依赖 admin.rpc
start_service "bot"          "$CONFIG_DIR/bot.yaml"          8086  # 依赖 msg.rpc

# ========== 5. Infrastructure 服务 ==========
info "========== [Layer 5] 启动 Infrastructure 服务 =========="

start_service "msgtransfer" "$CONFIG_DIR/msgtransfer.yaml"           # 依赖 NATS/MongoDB，无端口
start_service "msggateway"  "$CONFIG_DIR/msggateway.yaml"    8088    # 依赖 user/msg
start_service "push"        "$CONFIG_DIR/push.yaml"          8089    # 依赖 msggateway
start_service "cron"        "$CONFIG_DIR/cron.yaml"                  # 依赖 msg/conversation

# ========== 6. API 网关（依赖所有 RPC）==========
info "========== [Layer 6] 启动 API 网关 =========="

start_service "api" "$CONFIG_DIR/api.yaml" 10002

# ========== 启动完成 ==========
echo ""
info "========== 全部服务启动完成 =========="
echo ""
echo "  API            http://0.0.0.0:10002"
echo "  WebSocket      ws://0.0.0.0:10001"
echo ""
echo "  日志目录:  $LOG_DIR/"
echo "  查看状态:  ./scripts/status.sh"
echo "  停止所有:  ./scripts/stop_all.sh"
echo ""
