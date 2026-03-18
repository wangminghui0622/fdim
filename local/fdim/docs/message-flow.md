# fdim 消息流转详解

## 目录

- [架构总览](#架构总览)
- [核心组件职责](#核心组件职责)
- [单聊消息流转](#单聊消息流转)
- [群聊消息流转](#群聊消息流转)
- [各服务详细操作](#各服务详细操作)
- [NATS Topic 说明](#nats-topic-说明)
- [Redis 数据结构](#redis-数据结构)
- [MongoDB 数据结构](#mongodb-数据结构)

---

## 架构总览

```
┌──────────┐     HTTP      ┌──────────┐    gRPC     ┌──────────┐
│  客户端   │─────────────▶│ API 服务  │───────────▶│ msg RPC  │
│ (H5/App) │              │ :10002   │             │ :8085    │
└──────────┘              └──────────┘             └────┬─────┘
     ▲                                                  │
     │ WebSocket                              ┌─────────┴─────────┐
     │                                        │                   │
     │                                   Redis INCR          NATS Produce
     │                                   (分配 seq)         [toRedis] topic
     │                                        │                   │
     │                                        ▼                   ▼
     │                                  ┌──────────┐      ┌──────────────┐
     │                                  │  Redis   │      │ msgtransfer  │
     │                                  └──────────┘      └──────┬───────┘
     │                                        ▲            ┌─────┼─────┐
     │                                        │             │     │     │
     │                                  写消息缓存    NATS   │  NATS  │  NATS
     │                                        │      [toMongo]│ [toPush]│
     │                                        │             │     │     │
     │                                        │             ▼     │     ▼
     │                                        │      ┌────────┐  │ ┌────────┐
     │                                        └──────│MongoDB │  │ │ push   │
     │                                               │持久存储 │  │ │ 服务    │
     │                                               └────────┘  │ └───┬────┘
     │                                                           │     │
     │                                                           │  gRPC 调用
     │                                                           │     │
     │                                                           │     ▼
     │                                   ┌───────────────────────┘ ┌────────────┐
     │                                   │                         │ msggateway │
     │                                   │                         │ :10001     │
     │                                   │                         └─────┬──────┘
     │                                   │                               │
     └───────────────────────────────────┴───────────────── WebSocket ───┘
```

---

## 核心组件职责

| 组件 | 类型 | 职责 |
|------|------|------|
| **API** | HTTP 服务 (:10002) | 唯一的 HTTP 入口，转发请求到 RPC 服务 |
| **msg RPC** | gRPC 服务 (:8085) | 消息核心逻辑：校验、分配 seq、投递 NATS |
| **msgtransfer** | 后台服务 | 消息中转站：消费 NATS，写 Redis 缓存 + MongoDB + 推送 |
| **push** | gRPC + 后台服务 (:8089) | 在线推送 + 离线推送 |
| **msggateway** | gRPC + WebSocket (:8088/:10001) | 管理 WebSocket 连接，推送消息给在线用户 |
| **Redis** | 基础设施 (:6379) | 消息缓存（7天）、seq 计数器、在线状态 |
| **MongoDB** | 基础设施 (:27017) | 消息永久存储、历史记录查询 |
| **NATS** | 基础设施 (:4222) | 服务间异步消息队列 |

---

## 单聊消息完整流程

> 场景：用户 A (userID="alice") **第一次**发送 "你好" 给用户 B (userID="bob")

### 关键点

1. **会话创建**：msgtransfer 在处理消息时自动创建会话（为发送者和接收者各创建一份）
2. **seq 分配**：msg RPC 通过 Redis INCR 原子性分配序列号
3. **异步处理**：消息发送后立即返回，后续存储和推送全部异步
4. **已读未读**：通过 conversation 表和 Redis 的 has_read_seq 管理

### 流程图

```
用户A 发送 "你好"
    │
    │ ① HTTP POST /msg/send
    │    Body: { sendID:"alice", recvID:"bob", content:"你好", sessionType:1 }
    ▼
┌──────────┐
│ API 服务  │
│          │  校验 token → 从 header 提取 userID
│          │  转换为 gRPC 请求
└────┬─────┘
     │
     │ ② gRPC: msg.SendMsg(SendMsgReq)
     ▼
┌──────────┐
│ msg RPC  │
│          │  a. 校验参数（clientMsgID、sendID 不为空）
│          │  b. 生成会话ID:
│          │     sort(["alice","bob"]) → "si_alice_bob"
│          │  c. Redis INCR → conversation:max_seq:si_alice_bob → seq=1
│          │  d. 填充: msgData.Seq=1, SendTime=当前时间, ServerMsgID
│          │  e. 构造 PushMsgDataToMQ protobuf
│          │  f. 返回给客户端: { seq:1, serverMsgID, sendTime }
└────┬─────┘
     │
     │ ③ NATS Produce → topic: "toRedis"
     │    key: "si_alice_bob"
     │    value: PushMsgDataToMQ { msgData, conversationID }
     ▼
┌──────────────┐
│ msgtransfer  │  消费 [toRedis] (groupID: "redis")
│              │
│  步骤 1:     │  ④ **创建会话** (第一次发消息时)
│  CreateSingle│     gRPC: conversation.CreateSingleChatConversations
│  ChatConv    │     - 为 alice 创建会话: { ownerUserID:"alice", userID:"bob", conversationID:"si_alice_bob" }
│              │     - 为 bob 创建会话: { ownerUserID:"bob", userID:"alice", conversationID:"si_alice_bob" }
│              │     → MongoDB: db.conversation.insertMany([...])
│              │
│  步骤 2:     │  ⑤ 写 Redis 消息缓存
│  SetMessages │     SET msg:si_alice_bob:1 → {消息JSON} (TTL 7天)
│  ToCache     │
│              │
│  步骤 3:     │  ⑥ NATS Produce → topic: "toMongo"
│  转发toMongo │     value: MsgDataToMongoByMQ { conversationID, msgData[], lastSeq }
│              │
│  步骤 4:     │  ⑦ NATS Produce → topic: "toPush"
│  转发toPush  │     value: PushMsgDataToMQ { msgData, conversationID }
└──────┬───────┘
       │
       ├──── [toMongo] ────────────────────────────────┐
       │                                               ▼
       │                                    ┌──────────────────┐
       │                                    │ msgtransfer      │
       │                                    │ (toMongo handler) │
       │                                    │                  │
       │                                    │ ⑧ MongoDB Insert: │
       │                                    │ db.stream_msg.   │
       │                                    │ insertMany([{    │
       │                                    │   conversation_id│
       │                                    │   seq: 1         │
       │                                    │   send_id: alice │
       │                                    │   recv_id: bob   │
       │                                    │   content: 你好   │
       │                                    │   send_time: ... │
       │                                    │ }])              │
       │                                    └──────────────────┘
       │
       └──── [toPush] ─────────────────────────────────┐
                                                       ▼
                                            ┌──────────────────┐
                                            │ push 服务         │
                                            │                  │
                                            │ 解析消息:         │
                                            │  sessionType=1   │
                                            │  → 单聊           │
                                            │  → pushToUserID  │
                                            │    = "bob"       │
                                            └────────┬─────────┘
                                                     │
                                                     │ ⑨ gRPC: msggateway.OnlinePushMsg
                                                     │    { msgData, pushToUserID:"bob" }
                                                     ▼
                                            ┌──────────────────┐
                                            │ msggateway       │
                                            │                  │
                                            │ 查找 bob 的      │
                                            │ WebSocket 连接   │
                                            │                  │
                                            │ 如果在线:         │
                                            │ ⑩ ws.Send(msgData)│
                                            │                  │
                                            │ 如果不在线:       │
                                            │ 返回空列表        │
                                            └──────────────────┘
                                                     │
                                                     │ ⑩ WebSocket 推送
                                                     ▼
                                                 用户B 收到 "你好" ✓
                                                 
                                                 此时 bob 的 Tabbar 会显示:
                                                 - 会话列表出现 alice
                                                 - 未读小红点: 1
```

### 数据变化时间线

| 步骤 | 时间 | Redis | MongoDB | NATS |
|------|------|-------|---------|------|
| ③ msg RPC | 0ms | `conversation:max_seq:si_alice_bob = 1` | — | produce → [toRedis] |
| ④ msgtransfer | ~5ms | — | `conversation` 表新增 2 条（alice 和 bob 各一条） | — |
| ⑤ msgtransfer | ~5ms | `msg:si_alice_bob:1 = {消息JSON}` | — | — |
| ⑥ msgtransfer | ~5ms | — | — | produce → [toMongo] |
| ⑦ msgtransfer | ~5ms | — | — | produce → [toPush] |
| ⑧ toMongo handler | ~20ms | — | `stream_msg` 新增一条文档 | — |
| ⑨⑩ push+gateway | ~10ms | — | — | — |

**用户 A 只等了步骤 ③（~2ms），后续全部异步。**

---

## 会话管理详解

### 会话的作用

会话（Conversation）是 IM 系统的核心概念，用于：
1. **组织消息**：将用户之间的消息归类到一个会话中
2. **管理未读数**：记录每个用户在每个会话中的已读位置
3. **显示列表**：客户端的 Tabbar 消息列表就是会话列表
4. **存储设置**：置顶、免打扰、草稿等会话级别的设置

### 会话 ID 规则

| 类型 | 前缀 | 生成规则 | 示例 |
|------|------|----------|------|
| **单聊** | `si_` | `si_` + sort(userID1, userID2) | `si_alice_bob` |
| **群聊** | `sg_` | `sg_` + groupID | `sg_group001` |
| **通知** | `n_` | `n_` + sendID + `_` + recvID | `n_system_alice` |

**为什么单聊要 sort？**
- alice 发给 bob：conversationID = `si_alice_bob`
- bob 发给 alice：conversationID = `si_alice_bob`（相同）
- 保证双方使用同一个会话 ID，消息存储在一起

### 会话创建时机

**官方实现**：msgtransfer 在处理第一条消息时自动创建会话

```go
// E:/im/official/server/open-im-server/internal/msgtransfer/online_history_msg_handler.go:301-311
case constant.SingleChatType, constant.NotificationChatType:
    req := &pbconv.CreateSingleChatConversationsReq{
        RecvID:           msg.RecvID,
        SendID:           msg.SendID,
        ConversationID:   conversationID,
        ConversationType: msg.SessionType,
    }
    if err := och.conversationClient.CreateSingleChatConversations(ctx, req); err != nil {
        log.ZWarn(ctx, "single chat or notification first create conversation error", err)
    }
```

**创建内容**：
- 为**发送者**创建一条会话记录：`{ ownerUserID: alice, userID: bob, conversationID: si_alice_bob }`
- 为**接收者**创建一条会话记录：`{ ownerUserID: bob, userID: alice, conversationID: si_alice_bob }`

### 会话数据结构（MongoDB）

```json
// db.conversation 集合
{
  "_id": ObjectId("..."),
  "owner_user_id": "alice",           // 会话所有者
  "user_id": "bob",                   // 对方用户ID（单聊）
  "group_id": "",                     // 群ID（群聊）
  "conversation_id": "si_alice_bob",  // 会话ID
  "conversation_type": 1,             // 1=单聊, 3=群聊
  "is_pinned": false,                 // 是否置顶
  "recv_msg_opt": 0,                  // 接收消息选项（0=正常, 1=不接收, 2=接收但不提示）
  "is_private_chat": false,           // 是否私聊
  "burn_duration": 0,                 // 阅后即焚时长
  "group_at_type": 0,                 // 群@类型
  "is_not_in_group": false,           // 是否已不在群
  "update_unread_count_time": 0,      // 更新未读数时间
  "attached_info": "",                // 附加信息
  "ex": "",                           // 扩展字段
  "max_seq": 0,                       // 最大seq（用于增量拉取）
  "min_seq": 0,                       // 最小seq
  "create_time": ISODate("..."),      // 创建时间
  "draft_text_time": 0,               // 草稿时间
  "draft_text": ""                    // 草稿内容
}
```

### 会话查询

**客户端获取会话列表**：
```
GET /conversation/get_sorted_conversation_list
→ conversation RPC: GetSortedConversationList
→ 查询 MongoDB: db.conversation.find({ owner_user_id: "alice" })
→ 返回 alice 的所有会话（按最后消息时间排序）
```

**为什么需要两条会话记录？**
- alice 的会话列表：显示 bob（ownerUserID=alice）
- bob 的会话列表：显示 alice（ownerUserID=bob）
- 每个用户有自己的会话设置（置顶、免打扰等）

---

## 已读未读机制

### 核心概念

1. **maxSeq**：会话中的最大消息序列号（所有消息的最新位置）
2. **hasReadSeq**：用户在该会话中已读到的序列号
3. **未读数** = maxSeq - hasReadSeq

### 数据存储

| 数据 | 存储位置 | Key 格式 | 示例 |
|------|----------|----------|------|
| **maxSeq** | Redis | `conversation:max_seq:{conversationID}` | `conversation:max_seq:si_alice_bob = 5` |
| **hasReadSeq** | Redis | `conversation:has_read:{conversationID}:{userID}` | `conversation:has_read:si_alice_bob:bob = 3` |
| **会话信息** | MongoDB | `conversation` 集合 | `{ owner_user_id: "bob", max_seq: 5 }` |

### 标记已读流程

> 场景：bob 打开和 alice 的聊天页面，标记已读到 seq=5

```
客户端 bob
    │
    │ ① HTTP POST /msg/mark_conversation_as_read
    │    Body: { userID:"bob", conversationID:"si_alice_bob", hasReadSeq:5 }
    ▼
┌──────────┐
│ API 服务  │
└────┬─────┘
     │
     │ ② gRPC: msg.MarkConversationAsRead
     ▼
┌──────────┐
│ msg RPC  │
│          │  a. 获取会话信息（验证会话存在）
│          │     gRPC: conversation.GetConversation
│          │     → 从 MongoDB 查询会话
│          │
│          │  b. 获取当前已读位置
│          │     Redis GET: conversation:has_read:si_alice_bob:bob → 3
│          │
│          │  c. 标记消息为已读（单聊）
│          │     生成 seqs = [4, 5]（从 hasReadSeq+1 到 新hasReadSeq）
│          │     Redis: 为每条消息设置已读标记
│          │
│          │  d. 更新已读位置
│          │     Redis SET: conversation:has_read:si_alice_bob:bob = 5
│          │
│          │  e. 发送已读回执通知给 alice
│          │     NotificationSender.NotificationWithSessionType
│          │     → contentType = 2200 (HasReadReceipt)
│          │     → 通过 WebSocket 推送给 alice
│          │
│          │  f. 发送未读数变更通知给 bob 自己
│          │     ConversationUnreadChangeNotification
│          │     → unreadCount = maxSeq - hasReadSeq = 0
└──────────┘
```

### 已读回执

**alice 收到的已读回执**：
```json
{
  "contentType": 2200,  // HasReadReceipt
  "content": {
    "markAsReadUserID": "bob",
    "conversationID": "si_alice_bob",
    "seqs": [4, 5],
    "hasReadSeq": 5
  }
}
```

**客户端处理**：
- alice 的聊天界面显示：bob 已读（双勾变蓝）
- alice 发送的消息下方显示"已读"状态

### 未读小红点

**bob 的未读数计算**：
```
maxSeq = 5（会话最新消息）
hasReadSeq = 3（bob 已读到第3条）
未读数 = 5 - 3 = 2
```

**客户端显示**：
- Tabbar 消息列表：alice 的会话右侧显示红色数字 "2"
- 打开聊天页面后，调用 MarkConversationAsRead，未读数清零

### 获取未读数

```
GET /msg/get_conversations_has_read_and_max_seq
→ msg RPC: GetConversationsHasReadAndMaxSeq
→ 返回: {
    "seqs": {
      "si_alice_bob": {
        "hasReadSeq": 3,
        "maxSeq": 5
      }
    }
  }
→ 客户端计算: unreadCount = 5 - 3 = 2
```

### 关键点

1. **会话必须先创建**：MarkConversationAsRead 依赖 conversation.GetConversation，如果会话不存在会报错
2. **msgtransfer 负责创建会话**：在处理第一条消息时自动创建
3. **已读回执是可选的**：只有单聊才发送已读回执给对方
4. **未读数实时更新**：通过 WebSocket 推送 ConversationUnreadChangeNotification

---

## 群聊消息流转

> 场景：用户 A (userID="alice") 在群 "group001" 中发送 "大家好"
> 群成员：alice, bob, charlie

### 流程图

```
用户A 发送 "大家好"
    │
    │ ① HTTP POST /msg/send
    │    Body: { sendID:"alice", groupID:"group001", content:"大家好", sessionType:3 }
    ▼
┌──────────┐
│ API 服务  │
└────┬─────┘
     │
     │ ② gRPC: msg.SendMsg
     ▼
┌──────────┐
│ msg RPC  │
│          │  a. 生成会话ID: "sg_group001"
│          │  b. Redis INCR → conversation:max_seq:sg_group001 → seq=5
│          │     （假设群里已有4条消息）
│          │  c. 填充 msgData.Seq=5
│          │  d. 返回给 alice: { seq:5, ... }
└────┬─────┘
     │
     │ ③ NATS Produce → [toRedis]
     ▼
┌──────────────┐
│ msgtransfer  │
│              │  ④ Redis: SET msg:sg_group001:5 → {消息JSON}
│              │  ⑤ NATS → [toMongo]
│              │  ⑥ NATS → [toPush]
└──────┬───────┘
       │
       ├──── [toMongo] → ⑦ MongoDB 写入（同单聊）
       │
       └──── [toPush] ────────────────────────────────┐
                                                      ▼
                                            ┌─────────────────┐
                                            │ push 服务        │
                                            │                 │
                                            │ 解析消息:        │
                                            │  sessionType=3  │
                                            │  → 群聊          │
                                            │  → groupID =    │
                                            │    "group001"   │
                                            │                 │
                                            │ ⑧ gRPC 调用     │
                                            │ GroupClient.    │
                                            │ GetGroupMember  │
                                            │ UserIDs         │
                                            │ → ["alice",     │
                                            │    "bob",       │
                                            │    "charlie"]   │
                                            │                 │
                                            │ 排除发送者 alice  │
                                            │ → ["bob",       │
                                            │    "charlie"]   │
                                            └────────┬────────┘
                                                     │
                                                     │ ⑨ gRPC: msggateway.
                                                     │    OnlineBatchPushOneMsg
                                                     │    { msgData,
                                                     │      pushToUserIDs:
                                                     │      ["bob","charlie"] }
                                                     ▼
                                            ┌──────────────────┐
                                            │ msggateway       │
                                            │                  │
                                            │ 遍历每个 userID:  │
                                            │                  │
                                            │ bob 在线?        │
                                            │  → ⑩ ws.Send ✓  │
                                            │                  │
                                            │ charlie 在线?    │
                                            │  → ⑩ ws.Send ✓  │
                                            │                  │
                                            │ （不在线的忽略）   │
                                            └──────────────────┘
                                                     │
                                              ⑩ WebSocket 推送
                                                     │
                                              ┌──────┴──────┐
                                              ▼             ▼
                                          用户B 收到     用户C 收到
                                          "大家好" ✓    "大家好" ✓
```

### 群聊 vs 单聊的区别

| | 单聊 | 群聊 |
|---|---|---|
| **会话ID** | `si_` + sort(sendID, recvID) | `sg_` + groupID |
| **sessionType** | 1 | 3 (WriteGroupChatType) |
| **push 推送** | 直接推给 recvID | 先查群成员列表，排除发送者，批量推 |
| **gRPC 方法** | `OnlinePushMsg` | `OnlineBatchPushOneMsg` |
| **消息存储** | 一个会话一份 | 一个会话一份（所有群成员共享） |

---

## 各服务详细操作

### 1. msg RPC — SendMsg

```
输入: SendMsgReq { msgData: { sendID, recvID/groupID, content, sessionType, clientMsgID } }

处理:
  1. 校验: clientMsgID 和 sendID 不为空
  2. 生成 conversationID:
     - 单聊: si_ + sort(sendID, recvID)
     - 群聊: sg_ + groupID
  3. Redis INCR: conversation:max_seq:{conversationID} → 返回新 seq
  4. 填充字段: msgData.Seq, SendTime, ServerMsgID
  5. 构造 PushMsgDataToMQ { msgData, conversationID }
  6. proto.Marshal → NATS produce → [toRedis] topic

输出: SendMsgResp { serverMsgID, clientMsgID, sendTime, seq }

Redis 操作:
  - INCR conversation:max_seq:{conversationID}

NATS 操作:
  - Produce → [toRedis] topic, key=conversationID
```

### 2. msgtransfer — toRedis handler

```
消费: NATS [toRedis] topic (groupID: "redis")

处理:
  1. proto.Unmarshal → PushMsgDataToMQ
  2. 转换 protobuf MsgData → model.MsgDoc
  3. 写入 Redis 消息缓存（使用消息自带的 seq，不重新分配）
  4. 转发到 [toMongo] topic
  5. 转发到 [toPush] topic

Redis 操作:
  - SET msg:{conversationID}:{seq} → {消息JSON}  (TTL 7天)

NATS 操作:
  - Produce → [toMongo] topic (MsgDataToMongoByMQ)
  - Produce → [toPush] topic (PushMsgDataToMQ)
```

### 3. msgtransfer — toMongo handler

```
消费: NATS [toMongo] topic (groupID: "mongo")

处理:
  1. proto.Unmarshal → MsgDataToMongoByMQ
  2. 转换 protobuf MsgData[] → model.MsgDoc[]
  3. 批量写入 MongoDB

MongoDB 操作:
  - db.stream_msg.insertMany([...])
  - 索引: (conversation_id + seq) 唯一索引
```

### 4. push 服务 — toPush handler

```
消费: NATS [toPush] topic (groupID: "push")

处理:
  1. proto.Unmarshal → PushMsgDataToMQ
  2. 判断 sessionType:
     - 单聊 (1): pushToUserID = recvID
     - 群聊 (3): gRPC 调用 GroupClient.GetGroupMemberUserIDs → 排除发送者
  3. gRPC 调用 msggateway 推送

gRPC 调用:
  - 单聊: msggateway.OnlinePushMsg({ msgData, pushToUserID })
  - 群聊: msggateway.OnlineBatchPushOneMsg({ msgData, pushToUserIDs[] })
```

### 5. msggateway — WebSocket 推送

```
接收: gRPC OnlinePushMsg / OnlineBatchPushOneMsg

处理:
  1. 根据 userID 查找该用户的 WebSocket 连接
  2. 如果在线: 通过 WebSocket 发送消息
  3. 如果不在线: 返回空（push 服务可能触发离线推送）
```

---

## NATS Topic 说明

```
                    msg RPC
                       │
                       ▼
               ┌───────────────┐
               │  [toRedis]    │  所有新消息的入口
               │  groupID:redis│
               └───────┬───────┘
                       │
                  msgtransfer 消费
                       │
              ┌────────┴────────┐
              ▼                 ▼
     ┌───────────────┐  ┌───────────────┐
     │  [toMongo]    │  │  [toPush]     │
     │  groupID:mongo│  │  groupID:push │
     └───────┬───────┘  └───────┬───────┘
             │                  │
        msgtransfer 消费    push 服务消费
             │                  │
        MongoDB 写入       WebSocket 推送
```

| Topic | 生产者 | 消费者 | GroupID | 消息格式 |
|-------|--------|--------|---------|----------|
| `toRedis` | msg RPC | msgtransfer | redis | PushMsgDataToMQ |
| `toMongo` | msgtransfer | msgtransfer | mongo | MsgDataToMongoByMQ |
| `toPush` | msgtransfer | push 服务 | push | PushMsgDataToMQ |
| `toOfflinePush` | push 服务 | push 服务 | offlinePush | PushMsgDataToMQ |

---

## Redis 数据结构

| Key 格式 | 类型 | 用途 | TTL |
|----------|------|------|-----|
| `conversation:max_seq:{conversationID}` | String (int64) | 会话最大 seq 计数器 | 永不过期 |
| `conversation:min_seq:{conversationID}` | String (int64) | 会话最小保留 seq | 永不过期 |
| `msg:{conversationID}:{seq}` | String (JSON) | 消息内容缓存 | 7 天 |
| `conversation:has_read:{conversationID}:{userID}` | String (int64) | 用户已读 seq | 永不过期 |

### 示例

```
# 会话 si_alice_bob 已有 5 条消息
conversation:max_seq:si_alice_bob = 5

# 第 5 条消息的缓存
msg:si_alice_bob:5 = {
  "conversation_id": "si_alice_bob",
  "seq": 5,
  "send_id": "alice",
  "recv_id": "bob",
  "content": "你好",
  "content_type": 101,
  "send_time": 1709568000000,
  ...
}

# bob 已读到第 4 条
conversation:has_read:si_alice_bob:bob = 4
```

---

## MongoDB 数据结构

### 集合: `stream_msg`

```json
{
  "_id": ObjectId("..."),
  "conversation_id": "si_alice_bob",
  "seq": 5,
  "send_id": "alice",
  "recv_id": "bob",
  "group_id": "",
  "client_msg_id": "uuid-xxxx-xxxx",
  "server_msg_id": "uuid-xxxx-xxxx",
  "sender_platform_id": 5,
  "sender_nickname": "Alice",
  "sender_face_url": "https://...",
  "session_type": 1,
  "msg_from": 100,
  "content_type": 101,
  "content": "你好",
  "create_time": ISODate("2025-03-04T13:00:00Z"),
  "send_time": 1709568000000,
  "status": 0,
  "options": {},
  "at_user_ids": [],
  "attached_info": "",
  "ex": ""
}
```

### 索引

| 索引 | 字段 | 唯一 | 用途 |
|------|------|------|------|
| 1 | `(conversation_id, seq)` | ✅ | 按会话+seq 精确查找，防止重复 |
| 2 | `(conversation_id)` | ❌ | 按会话查询所有消息 |
| 3 | `(send_time)` | ❌ | 按时间范围查询、清理过期消息 |

---

## 客户端拉取消息

客户端不仅通过 WebSocket 实时接收推送，还需要**增量拉取**来保证不丢消息：

```
客户端本地记录: lastSeq = 3 （上次拉到 seq=3）

① 客户端请求: GET /msg/pull { conversationID:"si_alice_bob", startSeq:4 }

② API → msg RPC → Redis 查询:
   GET msg:si_alice_bob:4
   GET msg:si_alice_bob:5
   → 返回 2 条新消息

③ 如果 Redis 缓存过期（超过7天）:
   → 回退到 MongoDB 查询:
   db.stream_msg.find({ conversation_id:"si_alice_bob", seq: {$gte: 4} })

④ 客户端更新 lastSeq = 5
```

这就是为什么 **Redis 和 MongoDB 都要存**：
- **Redis**: 7天内的消息秒级返回
- **MongoDB**: 更早的历史消息兜底查询
