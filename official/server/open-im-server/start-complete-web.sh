#!/bin/bash

# OpenIM 完整 Web 系统启动脚本
# 功能：一键启动所有服务，包括 Web 前端和管理后台

set -e

echo "=========================================="
echo "OpenIM 完整 Web 系统启动脚本"
echo "=========================================="
echo ""

# 检查 Docker 是否运行
if ! docker info > /dev/null 2>&1; then
    echo "❌ 错误：Docker 未运行，请先启动 Docker"
    exit 1
fi

# 检查 docker-compose.yml 是否存在
if [ ! -f "docker-compose.yml" ]; then
    echo "❌ 错误：找不到 docker-compose.yml 文件"
    exit 1
fi

# 检查 .env 文件是否存在
if [ ! -f ".env" ]; then
    echo "⚠️  警告：找不到 .env 文件，将使用默认配置"
fi

echo "📦 正在启动所有服务..."
echo ""

# 启动服务
docker compose up -d

echo ""
echo "⏳ 等待服务启动（30秒）..."
sleep 30

echo ""
echo "📊 服务状态："
docker compose ps

echo ""
echo "=========================================="
echo "✅ 启动完成！"
echo "=========================================="
echo ""
echo "🌐 Web 聊天界面（用户端）："
echo "   http://localhost:11001"
echo ""
echo "🔧 管理后台（管理员）："
echo "   http://localhost:11002"
echo ""
echo "📡 API 服务："
echo "   OpenIM Server API: http://localhost:10002"
echo "   Chat Service API:  http://localhost:10008"
echo ""
echo "📝 查看日志："
echo "   docker compose logs -f <service-name>"
echo ""
echo "🛑 停止服务："
echo "   docker compose stop"
echo ""
