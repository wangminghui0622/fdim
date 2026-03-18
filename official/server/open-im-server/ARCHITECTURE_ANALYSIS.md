# OpenIM Server 架构分析

## 📋 目录
1. [整体架构概述](#整体架构概述)
2. [核心服务组件](#核心服务组件)
3. [基础设施组件](#基础设施组件)
4. [数据流架构](#数据流架构)
5. [部署模式](#部署模式)
6. [技术栈](#技术栈)
7. [服务间通信](#服务间通信)

---

## 整体架构概述

OpenIM Server 是一个**微服务架构**的即时通讯服务器，采用**分层设计**，支持**单体部署**和**微服务部署**两种模式。

### 架构层次

```
┌─────────────────────────────────────────────────────────────┐
│                     客户端层 (Client Layer)                    │
│  Flutter App / Web / iOS / Android / Desktop                 │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       │ WebSocket / HTTP
                       │
┌──────────────────────▼──────────────────────────────────────┐
│                   网关层 (Gateway Layer)                       │
│  ┌──────────────────┐  ┌──────────────────┐               │
│  │  Message Gateway  │  │   REST API        │               │
│  │  (WebSocket)      │  │   (HTTP)          │               │
│  │  Port: 10001      │  │   Port: 10002     │               │
│  └────────┬─────────┘  └────────┬─────────┘               │
└───────────┼──────────────────────┼────────────────────────┘
            │                        │
            │ gRPC                   │ gRPC
            │                        │
┌───────────▼──────────────────────▼────────────────────────┐
│                    RPC 服务层 (RPC Layer)                     │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐        │
│  │   Auth   │ │   User   │ │  Group  │ │   Msg   │        │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘        │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐                    │
│  │ Friend  │ │Conversat.│ │  Third  │                    │
│  └──────────┘ └──────────┘ └──────────┘                    │
└───────────┬──────────────────────┬────────────────────────┘
            │                        │
            │                        │
┌───────────▼──────────────────────▼────────────────────────┐
│              消息处理层 (Message Processing Layer)           │
│  ┌──────────────────┐  ┌──────────────────┐                │
│  │ Message Transfer │  │   Push Service   │                │
│  │  (Kafka Consumer)│  │  (Offline Push)  │                │
│  └──────────────────┘  └──────────────────┘                │
└───────────┬──────────────────────┬────────────────────────┘
            │                        │
            │                        │
┌───────────▼──────────────────────▼────────────────────────┐
│                   数据存储层 (Data Layer)                     │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐     │
│  │ MongoDB  │ │  Redis   │ │  Kafka  │ │  MinIO  │     │
│  │ (持久化)  │ │ (缓存)   │ │ (消息队列)│ │ (对象存储)│     │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘     │
│  ┌──────────┐                                              │
│  │   etcd   │                                              │
│  │ (服务发现)│                                              │
│  └──────────┘                                              │
└─────────────────────────────────────────────────────────────┘
```

---

## 核心服务组件

### 1. **Message Gateway (消息网关)**
- **端口**: 10001
- **协议**: WebSocket
- **职责**:
  - 维护客户端长连接
  - 接收客户端消息
  - 推送消息到客户端
  - 管理用户在线状态
  - 消息压缩/解压缩
- **位置**: `cmd/openim-msggateway/`, `internal/msggateway/`

### 2. **REST API (API 服务)**
- **端口**: 10002
- **协议**: HTTP/REST
- **职责**:
  - 提供 RESTful API 接口
  - 用户认证授权
  - 业务逻辑处理
  - 调用 RPC 服务
- **位置**: `cmd/openim-api/`, `internal/api/`

### 3. **RPC 服务组**

#### 3.1 Auth Service (认证服务)
- **职责**: 用户登录、Token 验证、权限管理
- **位置**: `cmd/openim-rpc/openim-rpc-auth/`, `internal/rpc/auth/`

#### 3.2 User Service (用户服务)
- **职责**: 用户信息管理、用户资料、用户搜索
- **位置**: `cmd/openim-rpc/openim-rpc-user/`, `internal/rpc/user/`

#### 3.3 Friend Service (好友服务)
- **职责**: 好友关系管理、好友申请、好友列表
- **位置**: `cmd/openim-rpc/openim-rpc-friend/`, `internal/rpc/relation/`

#### 3.4 Group Service (群组服务)
- **职责**: 群组创建、成员管理、群组信息
- **位置**: `cmd/openim-rpc/openim-rpc-group/`, `internal/rpc/group/`

#### 3.5 Conversation Service (会话服务)
- **职责**: 会话列表、会话设置、会话同步
- **位置**: `cmd/openim-rpc/openim-rpc-conversation/`, `internal/rpc/conversation/`

#### 3.6 Message Service (消息服务)
- **职责**: 消息存储、消息同步、消息历史
- **位置**: `cmd/openim-rpc/openim-rpc-msg/`, `internal/rpc/msg/`

#### 3.7 Third Service (第三方服务)
- **职责**: 文件上传、对象存储、第三方集成
- **位置**: `cmd/openim-rpc/openim-rpc-third/`, `internal/rpc/third/`

### 4. **Message Transfer (消息传输服务)**
- **核心功能**: **Kafka 生产端和消费端**
- **职责**:
  - **Kafka 消费者**:
    - 消费 `ToRedisTopic` 主题（从 RPC Message Service 接收消息）
    - 消费 `ToMongoTopic` 主题（批量写入 MongoDB）
  - **Kafka 生产者**:
    - 生产消息到 `ToMongoTopic` 主题（发送给 MongoDB Consumer）
    - 生产消息到 `ToPushTopic` 主题（发送给 Push Service）
  - **业务处理**:
    - 写入 Redis 缓存（批量处理）
    - 批量写入 MongoDB（异步持久化）
    - 管理消息序列号（Seq）
    - 处理用户已读状态
- **位置**: `cmd/openim-msgtransfer/`, `internal/msgtransfer/`
- **关键代码**:
  - `init.go:125-132` - 创建 Kafka 消费者
  - `init.go:89-96` - 创建 Kafka 生产者
  - `online_history_msg_handler.go:267` - 写入 Redis
  - `online_history_msg_handler.go:319` - 发送到 MongoDB Topic
  - `online_msg_to_mongo_handler.go:57` - 批量写入 MongoDB

### 5. **Push Service (推送服务)**
- **职责**:
  - 离线推送通知
  - 支持多种推送平台 (FCM, JPush, Getui 等)
  - 推送消息管理
- **位置**: `cmd/openim-push/`, `internal/push/`

### 6. **Cron Task (定时任务)**
- **职责**: 定时任务执行、数据清理、统计任务
- **位置**: `cmd/openim-crontask/`, `internal/tools/cron/`

---

## 基础设施组件

### 1. **MongoDB**
- **用途**: 主数据库，存储所有业务数据
- **数据**:
  - 用户信息
  - 消息记录
  - 群组信息
  - 好友关系
  - 会话信息

### 2. **Redis**
- **用途**: 缓存和会话存储
- **数据**:
  - 用户在线状态
  - 缓存数据
  - 分布式锁
  - 临时数据

### 3. **Kafka**
- **用途**: 消息队列，异步消息处理
- **功能**:
  - 消息发布/订阅
  - 消息持久化
  - 消息分发

### 4. **etcd**
- **用途**: 服务发现和配置管理
- **功能**:
  - RPC 服务注册与发现
  - 配置中心
  - 分布式协调

### 5. **MinIO**
- **用途**: 对象存储服务
- **功能**:
  - 文件上传/下载
  - 图片、视频、音频存储
  - 文件管理

---

## 数据流架构

### 消息发送流程

```
客户端
  │
  │ WebSocket
  ▼
Message Gateway (10001)
  │
  │ 验证、路由
  ▼
RPC Message Service
  │
  │ 【第一步】直接发送到 Kafka
  ▼
Kafka (ToRedisTopic) ⭐ 第一步：消息先到 Kafka
  │
  │
  ▼
Message Transfer (消费 ToRedisTopic)
  │
  │ 【第二步】写入 Redis 缓存
  ├─→ Redis (缓存) ⭐ 第二步：写入 Redis
  │
  │ 【第三步】发送到 MongoDB Topic
  ├─→ Kafka (ToMongoTopic) ⭐ 第三步：发送到 MongoDB Topic
  │     │
  │     │
  │     ▼
  │  Message Transfer (MongoDB Consumer)
  │     │
  │     │ 【第四步】批量写入 MongoDB
  │     └─→ MongoDB (持久化) ⭐ 第四步：最终写入 MongoDB
  │
  │ 【同时】发送到推送 Topic
  └─→ Kafka (ToPushTopic)
        │
        ├─→ 在线用户 → Message Gateway → 客户端
        │
        └─→ Push Service → FCM/JPush/Getui → 客户端
```

### ⚠️ 重要：数据写入顺序

**数据写入顺序：Kafka → Redis → Kafka → MongoDB**

1. **第一步：消息先到 Kafka**
   - RPC Message Service 收到消息后，**直接发送到 Kafka** (`ToRedisTopic`)
   - 代码位置：`internal/rpc/msg/send.go` → `MsgDatabase.MsgToMQ()`
   - **不先写 MongoDB**，而是先到 Kafka

2. **第二步：写入 Redis 缓存**
   - Message Transfer 服务从 Kafka (`ToRedisTopic`) 消费消息
   - **先写入 Redis 缓存** (`BatchInsertChat2Cache`)
   - 代码位置：`internal/msgtransfer/online_history_msg_handler.go`

3. **第三步：发送到 MongoDB Topic**
   - 从 Redis 缓存后，将消息发送到 Kafka 的 `ToMongoTopic`
   - 代码位置：`internal/msgtransfer/online_history_msg_handler.go` → `MsgToMongoMQ()`

4. **第四步：批量写入 MongoDB**
   - MongoDB Consumer 从 Kafka (`ToMongoTopic`) 消费消息
   - **批量写入 MongoDB** (`BatchInsertChat2DB`)
   - 代码位置：`internal/msgtransfer/online_msg_to_mongo_handler.go`

**设计原因：**
- ✅ **异步处理**：提高响应速度，不阻塞消息发送
- ✅ **解耦**：消息发送与持久化分离
- ✅ **可靠性**：Kafka 保证消息不丢失
- ✅ **性能**：批量写入 MongoDB，提高吞吐量
- ✅ **缓存优先**：先写 Redis，快速响应查询

### ⚠️ 重要：同步/异步执行说明

**Redis 写入和 Kafka 消息发送都是同步的，但整体流程是异步的：**

1. **Redis 写入 (`BatchInsertChat2Cache`)**
   - ✅ **同步执行**：直接调用 Redis 客户端，阻塞等待写入完成
   - 代码位置：`pkg/common/storage/controller/msg_transfer.go:201`
   - 必须等待 Redis 写入成功后才继续后续操作

2. **Kafka 消息发送 (`MsgToMongoMQ`, `MsgToPushMQ`)**
   - ✅ **同步执行**：使用 `sarama.SyncProducer`，阻塞等待消息发送完成
   - 代码位置：`pkg/common/storage/kafka/producer.go:29` - 使用 `sarama.SyncProducer`
   - 必须等待 Kafka 确认消息已发送后才返回

3. **整体流程是异步的**
   - ✅ **异步执行**：整个消息处理流程在 Kafka Consumer 的 goroutine 中执行
   - Message Transfer 服务从 Kafka 消费消息后，在独立的 goroutine 中处理
   - 代码位置：`internal/msgtransfer/init.go:154` - `go func() { ... }`

**执行顺序（同步阻塞）：**
```
Kafka Consumer (异步 goroutine)
  │
  ├─→ BatchInsertChat2Cache() ⏸️ 同步等待 Redis 写入完成
  │
  ├─→ MsgToMongoMQ() ⏸️ 同步等待 Kafka 发送完成
  │
  └─→ MsgToPushMQ() ⏸️ 同步等待 Kafka 发送完成
```

**为什么使用同步操作？**
- ✅ **数据一致性**：确保 Redis 写入成功后再发送到 MongoDB Topic
- ✅ **错误处理**：可以立即捕获和处理错误
- ✅ **顺序保证**：保证操作顺序（Redis → Kafka → MongoDB）
- ✅ **可靠性**：确保消息不丢失（Kafka SyncProducer 等待确认）

**性能优化：**
- 虽然单个操作是同步的，但整个流程在独立的 goroutine 中执行
- Kafka Consumer 可以并发处理多个消息批次
- 使用批量处理提高吞吐量（`batcher` 组件）

### 用户登录流程

```
客户端
  │
  │ HTTP POST /auth/user_register
  ▼
REST API (10002)
  │
  │ gRPC
  ▼
RPC Auth Service
  │
  ├─→ 验证用户信息
  │
  ├─→ 生成 Token
  │
  └─→ 返回 Token
        │
        └─→ 客户端
              │
              │ WebSocket 连接 (Token)
              ▼
        Message Gateway
              │
              │ 验证 Token
              ▼
        RPC Auth Service
              │
              └─→ 建立连接
```

### 消息同步流程

```
客户端
  │
  │ HTTP GET /msg/sync
  ▼
REST API
  │
  │ gRPC
  ▼
RPC Message Service
  │
  ├─→ 查询 MongoDB (消息历史)
  │
  └─→ 返回消息列表
        │
        └─→ 客户端
```

---

## 部署模式

### 1. **单体模式 (Monolith)**
- **入口**: `cmd/main.go`
- **特点**: 所有服务在同一个进程中运行
- **适用场景**: 开发、测试、小规模部署
- **启动方式**:
  ```bash
  ./openim-server -c config/
  ```

### 2. **微服务模式 (Microservices)**
- **特点**: 每个服务独立运行
- **服务列表**:
  - `openim-api`
  - `openim-msggateway`
  - `openim-msgtransfer`
  - `openim-push`
  - `openim-rpc-auth`
  - `openim-rpc-user`
  - `openim-rpc-friend`
  - `openim-rpc-group`
  - `openim-rpc-conversation`
  - `openim-rpc-msg`
  - `openim-rpc-third`
  - `openim-crontask`
- **适用场景**: 生产环境、大规模部署、高可用

### 3. **Docker Compose 部署**
- **配置文件**: `docker-compose.yml`
- **服务编排**: 包含所有基础设施和业务服务
- **网络**: 使用 `openim` 桥接网络

---

## 技术栈

### 后端技术
- **语言**: Go (Golang)
- **Web 框架**: Gin (REST API)
- **RPC 框架**: gRPC
- **WebSocket**: Gorilla WebSocket
- **ORM**: MongoDB Driver
- **消息队列**: Kafka
- **服务发现**: etcd
- **缓存**: Redis

### 基础设施
- **数据库**: MongoDB
- **缓存**: Redis
- **消息队列**: Apache Kafka
- **对象存储**: MinIO
- **服务发现**: etcd

### 监控与运维
- **监控**: Prometheus
- **可视化**: Grafana
- **告警**: Alertmanager
- **日志**: 结构化日志

---

## 服务间通信

### 1. **客户端 ↔ 服务端**
- **WebSocket**: 客户端 ↔ Message Gateway (实时消息)
- **HTTP**: 客户端 ↔ REST API (业务接口)

### 2. **服务间通信**
- **gRPC**: REST API ↔ RPC Services
- **Kafka**: Message Gateway → Kafka → Message Transfer
- **etcd**: 服务注册与发现

### 3. **数据访问**
- **MongoDB**: 所有 RPC 服务直接访问
- **Redis**: 缓存和会话管理
- **MinIO**: 文件存储访问

---

## 关键设计模式

### 1. **微服务架构**
- 服务解耦，独立部署
- 通过 gRPC 进行服务间通信
- 使用 etcd 进行服务发现

### 2. **事件驱动**
- 使用 Kafka 进行异步消息处理
- 消息发布/订阅模式

### 3. **CQRS (命令查询职责分离)**
- 写入操作通过 RPC 服务
- 查询操作通过 REST API

### 4. **缓存策略**
- Redis 缓存热点数据
- 减少数据库压力

### 5. **长连接管理**
- WebSocket 连接池
- 用户在线状态管理
- 连接心跳检测

---

## 扩展性设计

### 1. **水平扩展**
- 支持多实例部署
- 通过 etcd 进行负载均衡
- Kafka 支持分区扩展

### 2. **高可用**
- 服务无状态设计
- 数据库主从复制
- Redis 集群支持

### 3. **性能优化**
- 消息压缩
- 连接池管理
- 异步处理
- 缓存策略

---

## 配置文件结构

```
config/
├── discovery.yml      # 服务发现配置
├── mongodb.yml        # MongoDB 配置
├── redis.yml          # Redis 配置
├── kafka.yml          # Kafka 配置
├── minio.yml          # MinIO 配置
├── log.yml            # 日志配置
├── notification.yml  # 通知配置
├── webhooks.yml       # Webhook 配置
└── share.yml          # 共享配置
```

---

## 多机部署：如何找到接收方

### 问题场景

在多机（多实例）部署的情况下，当需要推送消息给用户时，如何找到接收方所在的服务器实例？

### 解决方案架构

```
┌─────────────────────────────────────────────────────────────┐
│                   消息推送流程（多机部署）                      │
└─────────────────────────────────────────────────────────────┘

Push Service (任意实例)
  │
  │ 1. 查询用户在线状态
  ├─→ OnlineCache (Redis)
  │     │
  │     └─→ 返回：用户是否在线、在线平台
  │
  │ 2. 如果用户在线，查找 Message Gateway 实例
  ├─→ Service Discovery (etcd)
  │     │
  │     └─→ 返回：所有 Message Gateway 实例列表
  │
  │ 3. 向所有 Message Gateway 实例广播推送请求
  ├─→ Message Gateway Instance 1 (gRPC)
  │     │
  │     ├─→ 检查本地连接 (UserMap)
  │     │     │
  │     │     ├─→ 有连接 → 推送消息 ✅
  │     │     └─→ 无连接 → 返回失败 ❌
  │     │
  ├─→ Message Gateway Instance 2 (gRPC)
  │     │
  │     ├─→ 检查本地连接 (UserMap)
  │     │     │
  │     │     ├─→ 有连接 → 推送消息 ✅
  │     │     └─→ 无连接 → 返回失败 ❌
  │     │
  └─→ Message Gateway Instance N (gRPC)
        │
        └─→ 检查本地连接...
```

### 核心机制

#### 1. **在线状态管理（Redis）**

**用户连接时：**
- Message Gateway 将用户在线状态写入 Redis
- 通过 Redis Pub/Sub 机制同步到所有实例
- 代码位置：`internal/msggateway/online.go` → `ChangeOnlineStatus()`

**存储结构：**
```
Redis Key: openim:user:online:{userID}
Value: [platformID1, platformID2, ...]
```

**在线状态同步：**
- Message Gateway 定期更新 Redis 中的在线状态
- 其他服务通过 Redis 订阅在线状态变化
- 代码位置：`pkg/rpccache/online.go` → `doSubscribe()`

#### 2. **查找在线用户**

**Push Service 查询流程：**
```go
// 1. 查询用户是否在线
onlineUserIDs, offlineUserIDs, err := onlineCache.GetUsersOnline(ctx, userIDs)

// 2. 如果在线，获取所有 Message Gateway 实例
conns, err := discovery.GetConns(ctx, "MessageGateway")

// 3. 向所有实例广播推送请求
for _, conn := range conns {
    msgClient.SuperGroupOnlineBatchPushOneMsg(ctx, req)
}
```

**代码位置：**
- `internal/push/push_handler.go:199` - `GetConnsAndOnlinePush()`
- `internal/push/onlinepusher.go:67` - `GetConnsAndOnlinePush()`

#### 3. **消息路由策略**

**策略一：广播模式（DefaultAllNode）**
- 向所有 Message Gateway 实例发送 gRPC 请求
- 每个实例检查本地是否有该用户的连接
- 有连接的实例推送消息，无连接的返回失败
- 代码位置：`internal/push/onlinepusher.go:67`

**策略二：一致性哈希（K8sStaticConsistentHash）**
- 使用一致性哈希算法确定用户所在的 Message Gateway 实例
- 只向特定实例发送请求（Kubernetes 环境）
- 代码位置：`internal/push/onlinepusher.go:150`

#### 4. **本地连接管理**

**每个 Message Gateway 实例：**
- 维护本地的用户连接映射（UserMap）
- 用户连接时注册到本地 Map
- 用户断开时从本地 Map 移除
- 代码位置：`internal/msggateway/user_map.go`

**连接查找：**
```go
// 检查本地是否有该用户的连接
clients, ok := wsServer.GetUserAllCons(userID)
if ok {
    // 推送消息到本地连接
    for _, client := range clients {
        client.PushMessage(ctx, msgData)
    }
}
```

**代码位置：**
- `internal/msggateway/hub_server.go:136` - `pushToUser()`
- `internal/msggateway/ws_server.go:129` - `GetUserAllCons()`

### 完整消息推送流程

```
1. Message Transfer 发送消息到 Kafka (ToPushTopic)
   │
   ▼
2. Push Service 消费消息
   │
   ├─→ 查询用户在线状态 (OnlineCache → Redis)
   │
   ├─→ 如果在线：
   │     │
   │     ├─→ 获取所有 Message Gateway 实例 (etcd)
   │     │
   │     ├─→ 向所有实例广播推送请求 (gRPC)
   │     │     │
   │     │     ├─→ Instance 1: 检查本地连接 → 推送/失败
   │     │     ├─→ Instance 2: 检查本地连接 → 推送/失败
   │     │     └─→ Instance N: 检查本地连接 → 推送/失败
   │     │
   │     └─→ 汇总推送结果
   │
   └─→ 如果离线：
         │
         └─→ 发送离线推送 (FCM/JPush/Getui)
```

### 关键设计点

1. **广播机制**
   - ✅ 简单可靠：不需要知道用户具体在哪台服务器
   - ✅ 容错性强：即使某台服务器故障，其他服务器仍可处理
   - ⚠️ 性能开销：需要向所有实例发送请求

2. **在线状态缓存**
   - ✅ 快速查询：Redis 缓存用户在线状态
   - ✅ 实时同步：通过 Pub/Sub 机制同步状态变化
   - ✅ 减少无效请求：只向在线用户推送

3. **本地连接管理**
   - ✅ 高效查找：本地 Map 快速查找连接
   - ✅ 内存管理：连接断开时及时清理
   - ✅ 并发安全：使用锁保护共享数据结构

4. **服务发现**
   - ✅ 动态发现：通过 etcd 自动发现所有实例
   - ✅ 负载均衡：支持多种路由策略
   - ✅ 高可用：实例故障时自动剔除

### 性能优化

1. **批量推送**
   - 支持批量推送多个用户
   - 减少 gRPC 调用次数
   - 代码位置：`internal/push/onlinepusher.go:83`

2. **并发控制**
   - 使用 `errgroup` 控制并发数
   - 避免过多并发请求
   - 代码位置：`internal/push/onlinepusher.go:91`

3. **缓存策略**
   - LRU 缓存用户在线状态
   - 减少 Redis 查询次数
   - 代码位置：`pkg/rpccache/online.go:66`

---

## 进程/服务数量统计

### 部署模式对比

#### 1. **单体模式 (Monolith)**
- **进程数**: **1 个进程**
- **说明**: 所有业务服务在同一个进程中运行
- **入口**: `cmd/main.go`
- **包含服务**:
  - 7 个 RPC 服务（Auth, User, Friend, Group, Conversation, Msg, Third）
  - 3 个业务服务（API, Message Gateway, Message Transfer, Push）
  - 1 个定时任务服务（Cron Task）

#### 2. **微服务模式 (Microservices)**
- **业务服务进程**: **12 个进程**
  1. `openim-api` - REST API 服务
  2. `openim-msggateway` - 消息网关（WebSocket）
  3. `openim-msgtransfer` - 消息传输服务
  4. `openim-push` - 推送服务
  5. `openim-crontask` - 定时任务服务
  6. `openim-rpc-auth` - 认证 RPC 服务
  7. `openim-rpc-user` - 用户 RPC 服务
  8. `openim-rpc-friend` - 好友 RPC 服务
  9. `openim-rpc-group` - 群组 RPC 服务
  10. `openim-rpc-conversation` - 会话 RPC 服务
  11. `openim-rpc-msg` - 消息 RPC 服务
  12. `openim-rpc-third` - 第三方 RPC 服务

- **基础设施服务**: **5 个进程**
  13. `mongo` - MongoDB 数据库
  14. `redis` - Redis 缓存
  15. `etcd` - 服务发现和配置中心
  16. `kafka` - 消息队列
  17. `minio` - 对象存储

- **前端服务**: **2 个进程**
  18. `openim-web-front` - Web 前端（用户聊天界面）
  19. `openim-admin-front` - 管理后台前端

- **Chat 服务**: **1 个进程**
  20. `openim-chat` - 聊天服务（管理后台、客服等功能）

- **监控服务（可选）**: **4 个进程**
  21. `prometheus` - 监控数据收集
  22. `grafana` - 监控数据可视化
  23. `alertmanager` - 告警管理
  24. `node-exporter` - 节点指标采集

### 总计

| 部署模式 | 业务服务 | 基础设施 | 前端 | Chat | 监控 | **总计** |
|---------|---------|---------|------|------|------|---------|
| **单体模式** | 1 | 5 | 2 | 1 | 0-4 | **9-13 个进程** |
| **微服务模式** | 12 | 5 | 2 | 1 | 0-4 | **20-24 个进程** |

### 进程数量说明

**最小部署（单体模式，不含监控）**:
- 1 个业务进程 + 5 个基础设施 + 2 个前端 + 1 个 Chat = **9 个进程**

**标准部署（微服务模式，不含监控）**:
- 12 个业务进程 + 5 个基础设施 + 2 个前端 + 1 个 Chat = **20 个进程**

**完整部署（微服务模式，含监控）**:
- 12 个业务进程 + 5 个基础设施 + 2 个前端 + 1 个 Chat + 4 个监控 = **24 个进程**

### 进程扩展说明

根据 `start-config.yml` 配置，某些服务可以启动多个实例：

- `openim-push`: 默认 8 个实例
- `openim-msgtransfer`: 默认 8 个实例
- `openim-crontask`: 默认 4 个实例
- 其他服务: 默认 1 个实例

**实际运行进程数 = 基础进程数 + 扩展实例数**

例如：微服务模式 + 扩展实例 = 12 + (8-1) + (8-1) + (4-1) = **29 个业务进程**

---

## 总结

OpenIM Server 采用**微服务架构**，通过**分层设计**实现了高可用、高性能的即时通讯服务。核心特点：

1. **分层清晰**: Gateway → API → RPC → Data
2. **服务解耦**: 各服务独立部署和扩展
3. **异步处理**: 使用 Kafka 进行消息异步处理
4. **高可用**: 支持集群部署和服务发现
5. **可扩展**: 支持水平扩展和性能优化

该架构设计能够支持**千万级用户**、**十万级超大群组**和**百亿级消息**的即时通讯场景。

### 进程数量总结

- **单体模式**: 1 个业务进程（所有服务在一个进程中）
- **微服务模式**: 12 个业务服务进程（每个服务独立运行）
- **完整部署**: 20-24 个进程（包含基础设施、前端、监控）
- **可扩展**: 支持多实例部署，实际进程数可达 30+ 个
