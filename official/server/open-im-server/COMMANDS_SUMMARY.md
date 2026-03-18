# OpenIM Docker Compose 部署命令总结

## 📋 完整部署流程命令

### 1. 进入项目目录
```bash
cd open-im-server
```

### 2. 创建 .env 环境变量文件
```bash
# Windows PowerShell
@"
# OpenIM Docker Compose 环境变量配置
DATA_DIR=./data

# MongoDB 配置
MONGO_IMAGE=mongo:7.0
MONGO_ADDRESS=mongo:27017
MONGO_USERNAME=openIM
MONGO_PASSWORD=openIM123

# Redis 配置
REDIS_IMAGE=redis:7.2-alpine
REDIS_ADDRESS=redis:6379
REDIS_PASSWORD=openIM123

# etcd 配置
ETCD_IMAGE=quay.io/coreos/etcd:v3.5.13
ETCD_ADDRESS=etcd:2379
ETCD_ROOT_USER=
ETCD_ROOT_PASSWORD=
ETCD_USERNAME=
ETCD_PASSWORD=

# Kafka 配置
KAFKA_IMAGE=apache/kafka:3.9.0
KAFKA_ADDRESS=kafka:9092
KAFKA_USERNAME=
KAFKA_PASSWORD=

# MinIO 配置
MINIO_IMAGE=minio/minio:RELEASE.2024-01-16T16-07-38Z
MINIO_PORT=10005
MINIO_CONSOLE_PORT=10006
MINIO_INTERNAL_ADDRESS=http://minio:9000
MINIO_EXTERNAL_ADDRESS=http://localhost:10005
MINIO_ACCESS_KEY_ID=root
MINIO_SECRET_ACCESS_KEY=openIM123

# OpenIM Server 配置
OPENIM_MSG_GATEWAY_PORT=10001
OPENIM_API_PORT=10002
OPENIM_SECRET=openIM123
LOG_IS_STDOUT=true
LOG_LEVEL=6

# Chat Service 配置
CHAT_API_PORT=10008
ADMIN_API_PORT=10009
API_URL=http://localhost:10002

# Web 前端镜像配置
OPENIM_WEB_FRONT_IMAGE=registry.cn-hangzhou.aliyuncs.com/openimsdk/openim-web-front:release-v3.8.3
OPENIM_WEB_FRONT_PORT=11001
OPENIM_ADMIN_FRONT_IMAGE=registry.cn-hangzhou.aliyuncs.com/openimsdk/openim-admin-front:release-v1.8.4
OPENIM_ADMIN_FRONT_PORT=11002

# 监控组件配置（可选）
PROMETHEUS_IMAGE=prom/prometheus:latest
PROMETHEUS_PORT=19090
ALERTMANAGER_IMAGE=prom/alertmanager:latest
ALERTMANAGER_PORT=19093
GRAFANA_IMAGE=grafana/grafana:latest
GRAFANA_PORT=19091
GRAFANA_URL=http://localhost:19091
NODE_EXPORTER_IMAGE=prom/node-exporter:latest
NODE_EXPORTER_PORT=19100
"@ | Out-File -FilePath .env -Encoding utf8
```

### 3. 启动所有服务
```bash
docker-compose up -d
```

### 4. 查看服务状态
```bash
docker-compose ps
```

### 5. 查看服务日志
```bash
# 查看所有服务日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f openim-server
docker-compose logs -f openim-chat
docker-compose logs -f etcd
```

### 6. 修复 etcd 配置（如果遇到冲突错误）
```bash
# 停止并删除 etcd 容器
docker-compose stop etcd
docker-compose rm -f etcd

# 重新启动 etcd
docker-compose up -d etcd

# 检查 etcd 状态
docker-compose ps etcd
docker logs etcd --tail 10
```

### 7. 修复 openim-chat 连接问题（如果注册无反应）
```bash
# 停止并删除 openim-chat 容器
docker-compose stop openim-chat
docker-compose rm -f openim-chat

# 重新启动 openim-chat
docker-compose up -d openim-chat

# 检查 openim-chat 状态
docker-compose ps openim-chat
docker logs openim-chat --tail 30
```

### 8. 重启单个服务
```bash
docker-compose restart <service-name>
# 例如：
docker-compose restart openim-server
docker-compose restart openim-chat
```

### 9. 停止所有服务
```bash
# 停止服务（保留数据）
docker-compose stop

# 停止并删除容器（保留数据）
docker-compose down

# 停止并删除容器和数据卷（⚠️ 会删除所有数据）
docker-compose down -v
```

### 10. 重新构建并启动（修改代码后）
```bash
docker-compose up -d --build
```

### 11. 检查服务健康状态
```bash
# 检查 openim-server
curl http://localhost:10002/healthz

# 检查 openim-chat
curl http://localhost:10008/healthz
```

### 12. 查看容器内环境变量（调试用）
```bash
# 查看 openim-chat 环境变量
docker exec openim-chat env | grep CHATENV

# 查看 openim-server 环境变量
docker exec openim-server env | grep IMENV
```

### 13. 查看容器内配置文件（调试用）
```bash
# 查看 openim-chat 配置文件
docker exec openim-chat cat /openim-chat/config/share.yml

# 查看 openim-server 配置文件
docker exec openim-server cat /openim-server/config/share.yml
```

## 🔧 常用故障排查命令

### 检查端口占用
```bash
# Windows
netstat -ano | findstr :10002
netstat -ano | findstr :10008
netstat -ano | findstr :11001
```

### 检查 Docker 网络
```bash
docker network ls
docker network inspect open-im-server_openim
```

### 检查容器日志（最后50行）
```bash
docker logs openim-server --tail 50
docker logs openim-chat --tail 50
docker logs etcd --tail 50
docker logs mongo --tail 50
docker logs redis --tail 50
docker logs kafka --tail 50
```

### 进入容器调试
```bash
# 进入 openim-server 容器
docker exec -it openim-server sh

# 进入 openim-chat 容器
docker exec -it openim-chat sh

# 进入 mongo 容器
docker exec -it mongo mongosh -u openIM -p openIM123
```

## 📝 重要配置文件位置

1. **docker-compose.yml**: `open-im-server/docker-compose.yml`
2. **.env**: `open-im-server/.env`
3. **chat 配置**: `chat/config/share.yml`
4. **openim-server 配置**: `open-im-server/config/share.yml`

## 🌐 访问地址

- **Web 聊天界面**: http://localhost:11001
- **管理后台**: http://localhost:11002
- **OpenIM Server API**: http://localhost:10002
- **Chat Service API**: http://localhost:10008
- **MinIO 控制台**: http://localhost:10006 (用户名: root, 密码: openIM123)

## ⚠️ 注意事项

1. 首次启动可能需要较长时间，因为需要：
   - 拉取基础镜像
   - 构建 openim-server 和 openim-chat 镜像
   - 初始化数据库

2. 如果遇到镜像拉取失败，可以：
   - 配置 Docker 镜像加速器
   - 手动拉取镜像：`docker pull <image-name>`

3. 修改配置后需要重启服务：
   ```bash
   docker-compose restart <service-name>
   ```

4. 修改代码后需要重新构建：
   ```bash
   docker-compose up -d --build
   ```
