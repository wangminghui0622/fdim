# OpenIM 完整 Web 系统部署指南

## 📋 系统架构

本系统包含以下组件：

1. **基础设施服务**
   - MongoDB（数据库）
   - Redis（缓存）
   - etcd（服务发现）
   - Kafka（消息队列）
   - MinIO（对象存储）

2. **核心服务**
   - `openim-server`：IM 核心服务（消息、用户、群组等）
   - `openim-chat`：业务服务（注册、登录、管理后台等）

3. **前端服务**
   - `openim-web-front`：Web 聊天界面（用户端）
   - `openim-admin-front`：管理后台界面（管理员）

## 🚀 快速开始

### 1. 前置要求

- Docker 和 Docker Compose 已安装
- 至少 4GB 可用内存
- 端口未被占用：`10001`, `10002`, `10008`, `11001`, `11002` 等

### 2. 一键启动

```bash
cd /Users/xx/workspace/src/open-im-server

# 启动所有服务（首次启动会自动拉取镜像）
docker compose up -d

# 查看服务状态
docker compose ps

# 查看日志（如果遇到问题）
docker compose logs -f openim-server
docker compose logs -f openim-chat
docker compose logs -f openim-web-front
```

### 3. 访问系统

启动成功后，访问以下地址：

#### 🌐 Web 聊天界面（用户端）
- **地址**：http://localhost:11001
- **功能**：注册、登录、加好友、聊天、建群等

#### 🔧 管理后台（管理员）
- **地址**：http://localhost:11002
- **功能**：用户管理、群组管理、消息管理等

#### 📡 API 服务
- **OpenIM Server API**：http://localhost:10002
- **Chat Service API**：http://localhost:10008

### 4. 默认配置

所有配置都在 `.env` 文件中，主要端口：

```env
OPENIM_API_PORT=10002          # OpenIM Server API
CHAT_API_PORT=10008            # Chat Service API
OPENIM_WEB_FRONT_PORT=11001    # Web 前端
OPENIM_ADMIN_FRONT_PORT=11002  # 管理后台
```

## 📝 使用说明

### 首次使用

1. **访问 Web 前端**：http://localhost:11001
2. **注册账号**：点击注册，输入手机号/邮箱，设置密码
3. **登录**：使用注册的账号登录
4. **添加好友**：搜索用户 ID 或手机号，发送好友请求
5. **开始聊天**：在好友列表中点击好友，开始发送消息

### 管理后台使用

1. **访问管理后台**：http://localhost:11002
2. **登录**：使用管理员账号（需要先在系统中创建）
3. **管理功能**：
   - 查看所有用户
   - 管理群组
   - 查看消息记录
   - 系统配置

## 🔍 故障排查

### 服务启动失败

```bash
# 查看所有服务状态
docker compose ps

# 查看特定服务日志
docker compose logs -f <service-name>

# 重启服务
docker compose restart <service-name>

# 重新构建并启动（如果修改了代码）
docker compose up -d --build
```

### 端口冲突

如果端口被占用，修改 `.env` 文件中的端口配置：

```bash
# 编辑 .env 文件
vim .env

# 修改端口（例如改为 12001）
OPENIM_WEB_FRONT_PORT=12001
```

### 镜像拉取失败

如果无法拉取官方镜像，可以使用国内镜像源：

```bash
# 编辑 .env 文件，取消注释以下行：
#OPENIM_WEB_FRONT_IMAGE=registry.cn-hangzhou.aliyuncs.com/openimsdk/openim-web-front:release-v3.8.3
#OPENIM_ADMIN_FRONT_IMAGE=registry.cn-hangzhou.aliyuncs.com/openimsdk/openim-admin-front:release-v1.8.4
```

### 数据库连接失败

确保 MongoDB、Redis 等服务已正常启动：

```bash
# 检查服务状态
docker compose ps mongo redis etcd kafka

# 查看日志
docker compose logs mongo
docker compose logs redis
```

## 🛑 停止服务

```bash
# 停止所有服务（保留数据）
docker compose stop

# 停止并删除容器（保留数据）
docker compose down

# 停止并删除容器和数据卷（⚠️ 会删除所有数据）
docker compose down -v
```

## 📚 相关文档

- [OpenIM 官方文档](https://docs.openim.io/)
- [Docker Compose 部署指南](https://docs.openim.io/guides/gettingStarted/dockerCompose)
- [API 文档](https://docs.openim.io/guides/development/api)

## 🎯 项目结构

```
/Users/xx/workspace/src/
├── open-im-server/          # IM 核心服务（当前目录）
│   ├── docker-compose.yml   # Docker Compose 配置
│   ├── .env                 # 环境变量配置
│   └── config/              # 服务配置文件
├── chat/                    # Chat 业务服务
└── openim-sdk-js-wasm/      # Web SDK（可选，用于自定义开发）
```

## ✅ 验证部署

启动后，运行以下命令验证服务是否正常：

```bash
# 检查所有容器是否运行
docker compose ps

# 检查 OpenIM Server 健康状态
curl http://localhost:10002/healthz

# 检查 Chat Service 健康状态
curl http://localhost:10008/healthz
```

所有服务显示为 `Up` 状态，且健康检查通过，即表示部署成功！
