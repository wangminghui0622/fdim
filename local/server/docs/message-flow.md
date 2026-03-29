# fdim 消息流转详解

## 目录

- [架构总览](#架构总览)
- [核心组件职责](#核心组件职责)
- [单聊消息流转](#单聊消息流转)
- [群聊消息流转](#群聊消息流转)
- [各服务详细操作](#各服务详细操�?
- [NATS Topic 说明](#nats-topic-说明)
- [Redis 数据结构](#redis-数据结构)
- [MongoDB 数据结构](#mongodb-数据结构)

---

## 架构总览

```
┌──────────�?    HTTP      ┌──────────�?   gRPC     ┌──────────�?�? 客户�?  │─────────────▶│ API 服务  │───────────▶│ msg RPC  �?�?(H5/App) �?             �?:10002   �?            �?:8085    �?└──────────�?             └──────────�?            └────┬─────�?     �?                                                 �?     �?WebSocket                              ┌─────────┴─────────�?     �?                                       �?                  �?     �?                                  Redis INCR          NATS Produce
     �?                                  (分配 seq)         [toRedis] topic
     �?                                       �?                  �?     �?                                       �?                  �?     �?                                 ┌──────────�?     ┌──────────────�?     �?                                 �? Redis   �?     �?msgtransfer  �?     �?                                 └──────────�?     └──────┬───────�?     �?                                       �?           ┌─────┼─────�?     �?                                       �?            �?    �?    �?     �?                                 写消息缓�?   NATS   �? NATS  �? NATS
     �?                                       �?     [toMongo]�?[toPush]�?     �?                                       �?            �?    �?    �?     �?                                       �?            �?    �?    �?     �?                                       �?     ┌────────�? �?┌────────�?     �?                                       └──────│MongoDB �? �?�?push   �?     �?                                              │持久存�?�? �?�?服务    �?     �?                                              └────────�? �?└───┬────�?     �?                                                          �?    �?     �?                                                          �? gRPC 调用
     �?                                                          �?    �?     �?                                                          �?    �?     �?                                  ┌───────────────────────�?┌────────────�?     �?                                  �?                        �?msggateway �?     �?                                  �?                        �?:10001     �?     �?                                  �?                        └─────┬──────�?     �?                                  �?                              �?     └───────────────────────────────────┴───────────────── WebSocket ───�?```

---

## 核心组件职责

| 组件 | 类型 | 职责 |
|------|------|------|
| **API** | HTTP 服务 (:10002) | 唯一�?HTTP 入口，转发请求到 RPC 服务 |
| **msg RPC** | gRPC 服务 (:8085) | 消息核心逻辑：校验、分�?seq、投�?NATS |
| **msgtransfer** | 后台服务 | 消息中转站：消费 NATS，写 Redis 缓存 + MongoDB + 推�?|
| **push** | gRPC + 后台服务 (:8089) | 在线推�?+ 离线推�?|
| **msggateway** | gRPC + WebSocket (:8088/:10001) | 管理 WebSocket 连接，推送消息给在线用户 |
| **Redis** | 基础设施 (:6379) | 消息缓存�?天）、seq 计数器、在线状�?|
| **MongoDB** | 基础设施 (:27017) | 消息永久存储、历史记录查�?|
| **NATS** | 基础设施 (:4222) | 服务间异步消息队�?|

---

## 单聊消息完整流程

> 场景：用�?A (userID="alice") **第一�?*发�?"你好" 给用�?B (userID="bob")

### 关键�?
1. **会话创建**：msgtransfer 在处理消息时自动创建会话（为发送者和接收者各创建一份）
2. **seq 分配**：msg RPC 通过 Redis INCR 原子性分配序列号
3. **异步处理**：消息发送后立即返回，后续存储和推送全部异�?4. **已读未读**：通过 conversation 表和 Redis �?has_read_seq 管理

### 流程�?
```
用户A 发�?"你好"
    �?    �?�?HTTP POST /msg/send
    �?   Body: { sendID:"alice", recvID:"bob", content:"你好", sessionType:1 }
    �?┌──────────�?�?API 服务  �?�?         �? 校验 token �?�?header 提取 userID
�?         �? 转换�?gRPC 请求
└────┬─────�?     �?     �?�?gRPC: msg.SendMsg(SendMsgReq)
     �?┌──────────�?�?msg RPC  �?�?         �? a. 校验参数（clientMsgID、sendID 不为空）
�?         �? b. 生成会话ID:
�?         �?    sort(["alice","bob"]) �?"si_alice_bob"
�?         �? c. Redis INCR �?conversation:max_seq:si_alice_bob �?seq=1
�?         �? d. 填充: msgData.Seq=1, SendTime=当前时间, ServerMsgID
�?         �? e. 构�?PushMsgDataToMQ protobuf
�?         �? f. 返回给客户端: { seq:1, serverMsgID, sendTime }
└────┬─────�?     �?     �?�?NATS Produce �?topic: "toRedis"
     �?   key: "si_alice_bob"
     �?   value: PushMsgDataToMQ { msgData, conversationID }
     �?┌──────────────�?�?msgtransfer  �? 消费 [toRedis] (groupID: "redis")
�?             �?�? 步骤 1:     �? �?**创建会话** (第一次发消息�?
�? CreateSingle�?    gRPC: conversation.CreateSingleChatConversations
�? ChatConv    �?    - �?alice 创建会话: { ownerUserID:"alice", userID:"bob", conversationID:"si_alice_bob" }
�?             �?    - �?bob 创建会话: { ownerUserID:"bob", userID:"alice", conversationID:"si_alice_bob" }
�?             �?    �?MongoDB: db.conversation.insertMany([...])
�?             �?�? 步骤 2:     �? �?�?Redis 消息缓存
�? SetMessages �?    SET msg:si_alice_bob:1 �?{消息JSON} (TTL 7�?
�? ToCache     �?�?             �?�? 步骤 3:     �? �?NATS Produce �?topic: "toMongo"
�? 转发toMongo �?    value: MsgDataToMongoByMQ { conversationID, msgData[], lastSeq }
�?             �?�? 步骤 4:     �? �?NATS Produce �?topic: "toPush"
�? 转发toPush  �?    value: PushMsgDataToMQ { msgData, conversationID }
└──────┬───────�?       �?       ├──── [toMongo] ────────────────────────────────�?       �?                                              �?       �?                                   ┌──────────────────�?       �?                                   �?msgtransfer      �?       �?                                   �?(toMongo handler) �?       �?                                   �?                 �?       �?                                   �?�?MongoDB Insert: �?       �?                                   �?db.stream_msg.   �?       �?                                   �?insertMany([{    �?       �?                                   �?  conversation_id�?       �?                                   �?  seq: 1         �?       �?                                   �?  send_id: alice �?       �?                                   �?  recv_id: bob   �?       �?                                   �?  content: 你好   �?       �?                                   �?  send_time: ... �?       �?                                   �?}])              �?       �?                                   └──────────────────�?       �?       └──── [toPush] ─────────────────────────────────�?                                                       �?                                            ┌──────────────────�?                                            �?push 服务         �?                                            �?                 �?                                            �?解析消息:         �?                                            �? sessionType=1   �?                                            �? �?单聊           �?                                            �? �?pushToUserID  �?                                            �?   = "bob"       �?                                            └────────┬─────────�?                                                     �?                                                     �?�?gRPC: msggateway.OnlinePushMsg
                                                     �?   { msgData, pushToUserID:"bob" }
                                                     �?                                            ┌──────────────────�?                                            �?msggateway       �?                                            �?                 �?                                            �?查找 bob �?     �?                                            �?WebSocket 连接   �?                                            �?                 �?                                            �?如果在线:         �?                                            �?�?ws.Send(msgData)�?                                            �?                 �?                                            �?如果不在�?       �?                                            �?返回空列�?       �?                                            └──────────────────�?                                                     �?                                                     �?�?WebSocket 推�?                                                     �?                                                 用户B 收到 "你好" �?                                                 
                                                 此时 bob �?Tabbar 会显�?
                                                 - 会话列表出现 alice
                                                 - 未读小红�? 1
```

### 数据变化时间�?
| 步骤 | 时间 | Redis | MongoDB | NATS |
|------|------|-------|---------|------|
| �?msg RPC | 0ms | `conversation:max_seq:si_alice_bob = 1` | �?| produce �?[toRedis] |
| �?msgtransfer | ~5ms | �?| `conversation` 表新�?2 条（alice �?bob 各一条） | �?|
| �?msgtransfer | ~5ms | `msg:si_alice_bob:1 = {消息JSON}` | �?| �?|
| �?msgtransfer | ~5ms | �?| �?| produce �?[toMongo] |
| �?msgtransfer | ~5ms | �?| �?| produce �?[toPush] |
| �?toMongo handler | ~20ms | �?| `stream_msg` 新增一条文�?| �?|
| ⑨⑩ push+gateway | ~10ms | �?| �?| �?|

**用户 A 只等了步�?③（~2ms），后续全部异步�?*

---

## 会话管理详解

### 会话的作�?
会话（Conversation）是 IM 系统的核心概念，用于�?1. **组织消息**：将用户之间的消息归类到一个会话中
2. **管理未读�?*：记录每个用户在每个会话中的已读位置
3. **显示列表**：客户端�?Tabbar 消息列表就是会话列表
4. **存储设置**：置顶、免打扰、草稿等会话级别的设�?
### 会话 ID 规则

| 类型 | 前缀 | 生成规则 | 示例 |
|------|------|----------|------|
| **单聊** | `si_` | `si_` + sort(userID1, userID2) | `si_alice_bob` |
| **群聊** | `sg_` | `sg_` + groupID | `sg_group001` |
| **通知** | `n_` | `n_` + sendID + `_` + recvID | `n_system_alice` |

**为什么单聊要 sort�?*
- alice 发给 bob：conversationID = `si_alice_bob`
- bob 发给 alice：conversationID = `si_alice_bob`（相同）
- 保证双方使用同一个会�?ID，消息存储在一�?
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

**创建内容**�?- �?*发送�?*创建一条会话记录：`{ ownerUserID: alice, userID: bob, conversationID: si_alice_bob }`
- �?*接收�?*创建一条会话记录：`{ ownerUserID: bob, userID: alice, conversationID: si_alice_bob }`

### 会话数据结构（MongoDB�?
```json
// db.conversation 集合
{
  "_id": ObjectId("..."),
  "owner_user_id": "alice",           // 会话所有�?  "user_id": "bob",                   // 对方用户ID（单聊）
  "group_id": "",                     // 群ID（群聊）
  "conversation_id": "si_alice_bob",  // 会话ID
  "conversation_type": 1,             // 1=单聊, 3=群聊
  "is_pinned": false,                 // 是否置顶
  "recv_msg_opt": 0,                  // 接收消息选项�?=正常, 1=不接�? 2=接收但不提示�?  "is_private_chat": false,           // 是否私聊
  "burn_duration": 0,                 // 阅后即焚时长
  "group_at_type": 0,                 // 群@类型
  "is_not_in_group": false,           // 是否已不在群
  "update_unread_count_time": 0,      // 更新未读数时�?  "attached_info": "",                // 附加信息
  "ex": "",                           // 扩展字段
  "max_seq": 0,                       // 最大seq（用于增量拉取）
  "min_seq": 0,                       // 最小seq
  "create_time": ISODate("..."),      // 创建时间
  "draft_text_time": 0,               // 草稿时间
  "draft_text": ""                    // 草稿内容
}
```

### 会话查询

**客户端获取会话列�?*�?```
GET /conversation/get_sorted_conversation_list
�?conversation RPC: GetSortedConversationList
�?查询 MongoDB: db.conversation.find({ owner_user_id: "alice" })
�?返回 alice 的所有会话（按最后消息时间排序）
```

**为什么需要两条会话记录？**
- alice 的会话列表：显示 bob（ownerUserID=alice�?- bob 的会话列表：显示 alice（ownerUserID=bob�?- 每个用户有自己的会话设置（置顶、免打扰等）

---

## 已读未读机制

### 核心概念

1. **maxSeq**：会话中的最大消息序列号（所有消息的最新位置）
2. **hasReadSeq**：用户在该会话中已读到的序列�?3. **未读�?* = maxSeq - hasReadSeq

### 数据存储

| 数据 | 存储位置 | Key 格式 | 示例 |
|------|----------|----------|------|
| **maxSeq** | Redis | `conversation:max_seq:{conversationID}` | `conversation:max_seq:si_alice_bob = 5` |
| **hasReadSeq** | Redis | `conversation:has_read:{conversationID}:{userID}` | `conversation:has_read:si_alice_bob:bob = 3` |
| **会话信息** | MongoDB | `conversation` 集合 | `{ owner_user_id: "bob", max_seq: 5 }` |

### 标记已读流程

> 场景：bob 打开�?alice 的聊天页面，标记已读�?seq=5

```
客户�?bob
    �?    �?�?HTTP POST /msg/mark_conversation_as_read
    �?   Body: { userID:"bob", conversationID:"si_alice_bob", hasReadSeq:5 }
    �?┌──────────�?�?API 服务  �?└────┬─────�?     �?     �?�?gRPC: msg.MarkConversationAsRead
     �?┌──────────�?�?msg RPC  �?�?         �? a. 获取会话信息（验证会话存在）
�?         �?    gRPC: conversation.GetConversation
�?         �?    �?�?MongoDB 查询会话
�?         �?�?         �? b. 获取当前已读位置
�?         �?    Redis GET: conversation:has_read:si_alice_bob:bob �?3
�?         �?�?         �? c. 标记消息为已读（单聊�?�?         �?    生成 seqs = [4, 5]（从 hasReadSeq+1 �?新hasReadSeq�?�?         �?    Redis: 为每条消息设置已读标�?�?         �?�?         �? d. 更新已读位置
�?         �?    Redis SET: conversation:has_read:si_alice_bob:bob = 5
�?         �?�?         �? e. 发送已读回执通知�?alice
�?         �?    NotificationSender.NotificationWithSessionType
�?         �?    �?contentType = 2200 (HasReadReceipt)
�?         �?    �?通过 WebSocket 推送给 alice
�?         �?�?         �? f. 发送未读数变更通知�?bob 自己
�?         �?    ConversationUnreadChangeNotification
�?         �?    �?unreadCount = maxSeq - hasReadSeq = 0
└──────────�?```

### 已读回执

**alice 收到的已读回�?*�?```json
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

**客户端处�?*�?- alice 的聊天界面显示：bob 已读（双勾变蓝）
- alice 发送的消息下方显示"已读"状�?
### 未读小红�?
**bob 的未读数计算**�?```
maxSeq = 5（会话最新消息）
hasReadSeq = 3（bob 已读到第3条）
未读�?= 5 - 3 = 2
```

**客户端显�?*�?- Tabbar 消息列表：alice 的会话右侧显示红色数�?"2"
- 打开聊天页面后，调用 MarkConversationAsRead，未读数清零

### 获取未读�?
```
GET /msg/get_conversations_has_read_and_max_seq
�?msg RPC: GetConversationsHasReadAndMaxSeq
�?返回: {
    "seqs": {
      "si_alice_bob": {
        "hasReadSeq": 3,
        "maxSeq": 5
      }
    }
  }
�?客户端计�? unreadCount = 5 - 3 = 2
```

### 关键�?
1. **会话必须先创�?*：MarkConversationAsRead 依赖 conversation.GetConversation，如果会话不存在会报�?2. **msgtransfer 负责创建会话**：在处理第一条消息时自动创建
3. **已读回执是可选的**：只有单聊才发送已读回执给对方
4. **未读数实时更�?*：通过 WebSocket 推�?ConversationUnreadChangeNotification

---

## 群聊消息流转

> 场景：用�?A (userID="alice") 在群 "group001" 中发�?"大家�?
> 群成员：alice, bob, charlie

### 流程�?
```
用户A 发�?"大家�?
    �?    �?�?HTTP POST /msg/send
    �?   Body: { sendID:"alice", groupID:"group001", content:"大家�?, sessionType:3 }
    �?┌──────────�?�?API 服务  �?└────┬─────�?     �?     �?�?gRPC: msg.SendMsg
     �?┌──────────�?�?msg RPC  �?�?         �? a. 生成会话ID: "sg_group001"
�?         �? b. Redis INCR �?conversation:max_seq:sg_group001 �?seq=5
�?         �?    （假设群里已�?条消息）
�?         �? c. 填充 msgData.Seq=5
�?         �? d. 返回�?alice: { seq:5, ... }
└────┬─────�?     �?     �?�?NATS Produce �?[toRedis]
     �?┌──────────────�?�?msgtransfer  �?�?             �? �?Redis: SET msg:sg_group001:5 �?{消息JSON}
�?             �? �?NATS �?[toMongo]
�?             �? �?NATS �?[toPush]
└──────┬───────�?       �?       ├──── [toMongo] �?�?MongoDB 写入（同单聊�?       �?       └──── [toPush] ────────────────────────────────�?                                                      �?                                            ┌─────────────────�?                                            �?push 服务        �?                                            �?                �?                                            �?解析消息:        �?                                            �? sessionType=3  �?                                            �? �?群聊          �?                                            �? �?groupID =    �?                                            �?   "group001"   �?                                            �?                �?                                            �?�?gRPC 调用     �?                                            �?GroupClient.    �?                                            �?GetGroupMember  �?                                            �?UserIDs         �?                                            �?�?["alice",     �?                                            �?   "bob",       �?                                            �?   "charlie"]   �?                                            �?                �?                                            �?排除发送�?alice  �?                                            �?�?["bob",       �?                                            �?   "charlie"]   �?                                            └────────┬────────�?                                                     �?                                                     �?�?gRPC: msggateway.
                                                     �?   OnlineBatchPushOneMsg
                                                     �?   { msgData,
                                                     �?     pushToUserIDs:
                                                     �?     ["bob","charlie"] }
                                                     �?                                            ┌──────────────────�?                                            �?msggateway       �?                                            �?                 �?                                            �?遍历每个 userID:  �?                                            �?                 �?                                            �?bob 在线?        �?                                            �? �?�?ws.Send �? �?                                            �?                 �?                                            �?charlie 在线?    �?                                            �? �?�?ws.Send �? �?                                            �?                 �?                                            �?（不在线的忽略）   �?                                            └──────────────────�?                                                     �?                                              �?WebSocket 推�?                                                     �?                                              ┌──────┴──────�?                                              �?            �?                                          用户B 收到     用户C 收到
                                          "大家�? �?   "大家�? �?```

### 群聊 vs 单聊的区�?
| | 单聊 | 群聊 |
|---|---|---|
| **会话ID** | `si_` + sort(sendID, recvID) | `sg_` + groupID |
| **sessionType** | 1 | 3 (WriteGroupChatType) |
| **push 推�?* | 直接推给 recvID | 先查群成员列表，排除发送者，批量�?|
| **gRPC 方法** | `OnlinePushMsg` | `OnlineBatchPushOneMsg` |
| **消息存储** | 一个会话一�?| 一个会话一份（所有群成员共享�?|

---

## 各服务详细操�?
### 1. msg RPC �?SendMsg

```
输入: SendMsgReq { msgData: { sendID, recvID/groupID, content, sessionType, clientMsgID } }

处理:
  1. 校验: clientMsgID �?sendID 不为�?  2. 生成 conversationID:
     - 单聊: si_ + sort(sendID, recvID)
     - 群聊: sg_ + groupID
  3. Redis INCR: conversation:max_seq:{conversationID} �?返回�?seq
  4. 填充字段: msgData.Seq, SendTime, ServerMsgID
  5. 构�?PushMsgDataToMQ { msgData, conversationID }
  6. proto.Marshal �?NATS produce �?[toRedis] topic

输出: SendMsgResp { serverMsgID, clientMsgID, sendTime, seq }

Redis 操作:
  - INCR conversation:max_seq:{conversationID}

NATS 操作:
  - Produce �?[toRedis] topic, key=conversationID
```

### 2. msgtransfer �?toRedis handler

```
消费: NATS [toRedis] topic (groupID: "redis")

处理:
  1. proto.Unmarshal �?PushMsgDataToMQ
  2. 转换 protobuf MsgData �?model.MsgDoc
  3. 写入 Redis 消息缓存（使用消息自带的 seq，不重新分配�?  4. 转发�?[toMongo] topic
  5. 转发�?[toPush] topic

Redis 操作:
  - SET msg:{conversationID}:{seq} �?{消息JSON}  (TTL 7�?

NATS 操作:
  - Produce �?[toMongo] topic (MsgDataToMongoByMQ)
  - Produce �?[toPush] topic (PushMsgDataToMQ)
```

### 3. msgtransfer �?toMongo handler

```
消费: NATS [toMongo] topic (groupID: "mongo")

处理:
  1. proto.Unmarshal �?MsgDataToMongoByMQ
  2. 转换 protobuf MsgData[] �?model.MsgDoc[]
  3. 批量写入 MongoDB

MongoDB 操作:
  - db.stream_msg.insertMany([...])
  - 索引: (conversation_id + seq) 唯一索引
```

### 4. push 服务 �?toPush handler

```
消费: NATS [toPush] topic (groupID: "push")

处理:
  1. proto.Unmarshal �?PushMsgDataToMQ
  2. 判断 sessionType:
     - 单聊 (1): pushToUserID = recvID
     - 群聊 (3): gRPC 调用 GroupClient.GetGroupMemberUserIDs �?排除发送�?  3. gRPC 调用 msggateway 推�?
gRPC 调用:
  - 单聊: msggateway.OnlinePushMsg({ msgData, pushToUserID })
  - 群聊: msggateway.OnlineBatchPushOneMsg({ msgData, pushToUserIDs[] })
```

### 5. msggateway �?WebSocket 推�?
```
接收: gRPC OnlinePushMsg / OnlineBatchPushOneMsg

处理:
  1. 根据 userID 查找该用户的 WebSocket 连接
  2. 如果在线: 通过 WebSocket 发送消�?  3. 如果不在�? 返回空（push 服务可能触发离线推送）
```

---

## NATS Topic 说明

```
                    msg RPC
                       �?                       �?               ┌───────────────�?               �? [toRedis]    �? 所有新消息的入�?               �? groupID:redis�?               └───────┬───────�?                       �?                  msgtransfer 消费
                       �?              ┌────────┴────────�?              �?                �?     ┌───────────────�? ┌───────────────�?     �? [toMongo]    �? �? [toPush]     �?     �? groupID:mongo�? �? groupID:push �?     └───────┬───────�? └───────┬───────�?             �?                 �?        msgtransfer 消费    push 服务消费
             �?                 �?        MongoDB 写入       WebSocket 推�?```

| Topic | 生产�?| 消费�?| GroupID | 消息格式 |
|-------|--------|--------|---------|----------|
| `toRedis` | msg RPC | msgtransfer | redis | PushMsgDataToMQ |
| `toMongo` | msgtransfer | msgtransfer | mongo | MsgDataToMongoByMQ |
| `toPush` | msgtransfer | push 服务 | push | PushMsgDataToMQ |
| `toOfflinePush` | push 服务 | push 服务 | offlinePush | PushMsgDataToMQ |

---

## Redis 数据结构

| Key 格式 | 类型 | 用�?| TTL |
|----------|------|------|-----|
| `conversation:max_seq:{conversationID}` | String (int64) | 会话最�?seq 计数�?| 永不过期 |
| `conversation:min_seq:{conversationID}` | String (int64) | 会话最小保�?seq | 永不过期 |
| `msg:{conversationID}:{seq}` | String (JSON) | 消息内容缓存 | 7 �?|
| `conversation:has_read:{conversationID}:{userID}` | String (int64) | 用户已读 seq | 永不过期 |

### 示例

```
# 会话 si_alice_bob 已有 5 条消�?conversation:max_seq:si_alice_bob = 5

# �?5 条消息的缓存
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

# bob 已读到第 4 �?conversation:has_read:si_alice_bob:bob = 4
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

| 索引 | 字段 | 唯一 | 用�?|
|------|------|------|------|
| 1 | `(conversation_id, seq)` | �?| 按会�?seq 精确查找，防止重�?|
| 2 | `(conversation_id)` | �?| 按会话查询所有消�?|
| 3 | `(send_time)` | �?| 按时间范围查询、清理过期消�?|

---

## 客户端拉取消�?
客户端不仅通过 WebSocket 实时接收推送，还需�?*增量拉取**来保证不丢消息：

```
客户端本地记�? lastSeq = 3 （上次拉�?seq=3�?
�?客户端请�? GET /msg/pull { conversationID:"si_alice_bob", startSeq:4 }

�?API �?msg RPC �?Redis 查询:
   GET msg:si_alice_bob:4
   GET msg:si_alice_bob:5
   �?返回 2 条新消息

�?如果 Redis 缓存过期（超�?天）:
   �?回退�?MongoDB 查询:
   db.stream_msg.find({ conversation_id:"si_alice_bob", seq: {$gte: 4} })

�?客户端更�?lastSeq = 5
```

这就是为什�?**Redis �?MongoDB 都要�?*�?- **Redis**: 7天内的消息秒级返�?- **MongoDB**: 更早的历史消息兜底查�?