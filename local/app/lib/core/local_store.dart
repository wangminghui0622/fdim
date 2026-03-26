import 'dart:convert';
import 'package:flutter/foundation.dart';
import 'package:hive/hive.dart';

import '../sdk/flutter_openim_sdk.dart';
import 'config.dart';

/// LocalStore - 模拟官方 openim-sdk-core 中的 SQLite 本地数据库。
///
/// 官方 Go SDK 在本地维护三类持久化数据：
/// 1. conversations   — 会话列表（含 unreadCount、latestMsg 等）
/// 2. messages        — 消息（含 isRead / hasReadTime）
/// 3. hasReadSeq      — 每个会话已读到的 seq（标记已读的 watermark）
///
/// 本类用 Hive Box 实现上述功能，key 设计：
///   - conversations box: key = conversationID, value = ConversationInfo JSON
///   - messages box:      key = "convID:clientMsgID", value = Message JSON
///   - meta box:          key = "hasReadSeq:convID", value = int
///                        key = "maxSeq:convID",     value = int
class LocalStore {
  LocalStore._();

  static Box? _convBox;
  static Box? _msgBox;
  static Box? _metaBox;

  static String get _userPrefix => Config.userID;

  /// 记录用户刚主动标记已读的会话 ID，防止服务端同步时回填旧未读数
  static final Set<String> _justReadConvIDs = {};

  /// 初始化（登录后调用一次）
  static Future<void> init() async {
    final uid = _userPrefix;
    if (uid.isEmpty) return;
    _convBox = await Hive.openBox('conv_$uid');
    _msgBox = await Hive.openBox('msg_$uid');
    _metaBox = await Hive.openBox('meta_$uid');
    debugPrint('[LocalStore] initialized for user $uid');
  }

  /// 登出时关闭
  static Future<void> close() async {
    await _convBox?.close();
    await _msgBox?.close();
    await _metaBox?.close();
    _convBox = null;
    _msgBox = null;
    _metaBox = null;
  }

  // ======================= Conversation =======================

  /// 保存 / 更新一条会话
  static Future<void> putConversation(ConversationInfo conv) async {
    final box = _convBox;
    if (box == null) return;
    await box.put(conv.conversationID, jsonEncode(conv.toJson()));
  }

  /// 批量保存会话（合并式）
  /// unreadCount 合并策略：
  /// 1. 若本地已建立 maxSeq/hasReadSeq，则优先使用本地 seq 水位差计算未读，避免服务端旧 unreadCount 回灌。
  /// 2. 若本地刚主动标记过已读，则保持 0，直到收到新的消息递增未读。
  /// 3. 若本地还没有有效 seq 水位，则回退到本地/服务端 unreadCount 的较大值。
  static Future<void> putConversations(List<ConversationInfo> list) async {
    final box = _convBox;
    if (box == null) return;
    debugPrint('[LocalStore] putConversations: ${list.length} conversations from server');
    for (final conv in list) {
      final local = getConversation(conv.conversationID);
      final serverUnread = conv.unreadCount;
      if (local != null) {
        final localUnread = local.unreadCount;
        final localMaxSeq = getMaxSeq(conv.conversationID);
        final localHasReadSeq = getHasReadSeq(conv.conversationID);
        if (localMaxSeq > 0) {
          final derivedUnread = localMaxSeq > localHasReadSeq
              ? localMaxSeq - localHasReadSeq
              : 0;
          conv.unreadCount = derivedUnread;
          debugPrint('[LocalStore] ${conv.conversationID}: derive by seq (server=$serverUnread, local=$localUnread, maxSeq=$localMaxSeq, hasReadSeq=$localHasReadSeq) -> $derivedUnread');
          if (derivedUnread == 0) {
            _justReadConvIDs.remove(conv.conversationID);
          }
        } else if (_justReadConvIDs.contains(conv.conversationID)) {
          conv.unreadCount = 0;
          debugPrint('[LocalStore] ${conv.conversationID}: justRead, keep 0 (server=$serverUnread, local=$localUnread)');
          if (serverUnread == 0) {
            _justReadConvIDs.remove(conv.conversationID);
            debugPrint('[LocalStore] ${conv.conversationID}: server synced, remove from justRead set');
          }
        } else {
          conv.unreadCount = localUnread > serverUnread ? localUnread : serverUnread;
          debugPrint('[LocalStore] ${conv.conversationID}: merge max (server=$serverUnread, local=$localUnread) -> ${conv.unreadCount}');
        }
      } else {
        debugPrint('[LocalStore] ${conv.conversationID}: new conversation, use server unread=$serverUnread');
      }
      await box.put(conv.conversationID, jsonEncode(conv.toJson()));
    }
  }

  /// 获取所有本地会话
  static List<ConversationInfo> getAllConversations() {
    final box = _convBox;
    if (box == null) return [];
    final result = <ConversationInfo>[];
    for (final raw in box.values) {
      try {
        final map = jsonDecode(raw as String) as Map<String, dynamic>;
        result.add(ConversationInfo.fromJson(map));
      } catch (_) {}
    }
    return result;
  }

  /// 获取单条会话
  static ConversationInfo? getConversation(String conversationID) {
    final box = _convBox;
    if (box == null) return null;
    final raw = box.get(conversationID);
    if (raw == null) return null;
    try {
      return ConversationInfo.fromJson(jsonDecode(raw as String));
    } catch (_) {
      return null;
    }
  }

  /// 删除会话
  static Future<void> deleteConversation(String conversationID) async {
    await _convBox?.delete(conversationID);
  }

  /// 设置会话 unreadCount（本地维护，与官方 SDK 一致）
  static Future<void> setUnreadCount(String conversationID, int count) async {
    final conv = getConversation(conversationID);
    if (conv != null) {
      conv.unreadCount = count;
      await putConversation(conv);
    }
  }

  /// 本地递增 unreadCount（收到新消息时调用）
  static Future<int> incrementUnread(String conversationID) async {
    final conv = getConversation(conversationID);
    if (conv != null) {
      _justReadConvIDs.remove(conversationID);
      conv.unreadCount += 1;
      await putConversation(conv);
      return conv.unreadCount;
    }
    return 0;
  }

  /// 本地清零 unreadCount（标记已读时调用）
  static Future<void> clearUnread(String conversationID) async {
    _justReadConvIDs.add(conversationID);
    await setUnreadCount(conversationID, 0);
  }

  /// 计算本地总未读数
  static int getTotalUnreadCount() {
    int total = 0;
    for (final conv in getAllConversations()) {
      total += conv.unreadCount;
    }
    return total;
  }

  /// 更新会话的 latestMsg 和 latestMsgSendTime
  static Future<void> updateLatestMsg(
      String conversationID, Message msg) async {
    final conv = getConversation(conversationID);
    if (conv != null) {
      conv.latestMsg = msg;
      conv.latestMsgSendTime = msg.sendTime;
      await putConversation(conv);
    }
  }

  // ======================= Message =======================

  static String _msgKey(String conversationID, String clientMsgID) =>
      '$conversationID:$clientMsgID';

  /// 保存一条消息到本地
  static Future<void> putMessage(
      String conversationID, Message msg) async {
    final box = _msgBox;
    if (box == null || msg.clientMsgID == null) return;
    await box.put(
        _msgKey(conversationID, msg.clientMsgID!), jsonEncode(msg.toJson()));
  }

  /// 批量保存消息
  static Future<void> putMessages(
      String conversationID, List<Message> msgs) async {
    final box = _msgBox;
    if (box == null) return;
    for (final msg in msgs) {
      if (msg.clientMsgID == null) continue;
      await box.put(
          _msgKey(conversationID, msg.clientMsgID!), jsonEncode(msg.toJson()));
    }
  }

  /// 获取单条消息
  static Message? getMessage(String conversationID, String clientMsgID) {
    final box = _msgBox;
    if (box == null) return null;
    final raw = box.get(_msgKey(conversationID, clientMsgID));
    if (raw == null) return null;
    try {
      return Message.fromJson(jsonDecode(raw as String));
    } catch (_) {
      return null;
    }
  }

  /// 获取会话中所有本地消息（按 seq 排序）
  static List<Message> getMessages(String conversationID) {
    final box = _msgBox;
    if (box == null) return [];
    final prefix = '$conversationID:';
    final result = <Message>[];
    for (final key in box.keys) {
      if ((key as String).startsWith(prefix)) {
        try {
          final map = jsonDecode(box.get(key) as String) as Map<String, dynamic>;
          result.add(Message.fromJson(map));
        } catch (_) {}
      }
    }
    result.sort((a, b) => (a.seq ?? 0).compareTo(b.seq ?? 0));
    return result;
  }

  /// 删除一条消息
  static Future<void> deleteMessage(
      String conversationID, String clientMsgID) async {
    await _msgBox?.delete(_msgKey(conversationID, clientMsgID));
  }

  // ======================= Read Status =======================

  /// 标记消息已读
  /// 双写策略：
  /// 1) 写入 message DB（conversationID:clientMsgID → msg JSON with isRead=true）
  /// 2) 写入 meta DB（readMsg:clientMsgID → readTime）作为 conversationID 无关的备份
  /// 这样即使 read receipt 的 conversationID 与本地存储不一致，也能恢复已读状态
  static Future<void> markMessageRead(
      String conversationID, String clientMsgID, int readTime) async {
    // 始终写入 conversationID-independent 的备份
    await _metaBox?.put('readMsg:$clientMsgID', readTime);
    // 尝试更新 message DB（可能因 conversationID 不匹配而找不到）
    final msg = getMessage(conversationID, clientMsgID);
    if (msg != null) {
      msg.isRead = true;
      msg.hasReadTime = readTime;
      await putMessage(conversationID, msg);
    }
  }

  /// 批量标记已读（通过 msgIDList）
  static Future<void> markMessagesRead(
      String conversationID, List<String> msgIDList, int readTime) async {
    for (final id in msgIDList) {
      await markMessageRead(conversationID, id, readTime);
    }
  }

  /// 检查消息是否已被标记为已读（conversation-independent 查询）
  static int? getMessageReadTime(String clientMsgID) {
    return _metaBox?.get('readMsg:$clientMsgID') as int?;
  }

  /// 获取 hasReadSeq（已读水位线）
  static int getHasReadSeq(String conversationID) {
    return _metaBox?.get('hasReadSeq:$conversationID') ?? 0;
  }

  /// 设置 hasReadSeq
  static Future<void> setHasReadSeq(String conversationID, int seq) async {
    final current = getHasReadSeq(conversationID);
    if (seq > current) {
      await _metaBox?.put('hasReadSeq:$conversationID', seq);
    }
  }

  /// 获取 peerReadSeq（对方已读水位线：对方已读到的最大 seq）
  /// 用于判断自己发的消息是否已被对方阅读
  static int getPeerReadSeq(String conversationID) {
    return _metaBox?.get('peerReadSeq:$conversationID') ?? 0;
  }

  /// 设置 peerReadSeq
  static Future<void> setPeerReadSeq(String conversationID, int seq) async {
    final current = getPeerReadSeq(conversationID);
    if (seq > current) {
      await _metaBox?.put('peerReadSeq:$conversationID', seq);
    }
  }

  /// 根据 hasReadSeq 恢复消息已读状态
  /// 官方逻辑: 自己发出的消息 isRead 由对端已读回执决定；
  /// 对方发的消息 seq <= hasReadSeq 就是已读
  static void applyHasReadSeq(
      String conversationID, List<Message> msgs, String myUserID) {
    final readSeq = getHasReadSeq(conversationID);
    for (final msg in msgs) {
      // 对方发的消息：如果 seq <= hasReadSeq，标记已读
      if (msg.sendID != myUserID && (msg.seq ?? 0) > 0 && (msg.seq ?? 0) <= readSeq) {
        msg.isRead = true;
        msg.hasReadTime ??= DateTime.now().millisecondsSinceEpoch;
      }
      // 自己发的消息：多策略恢复已读状态
      if (msg.sendID == myUserID) {
        // 策略 A: peerReadSeq 水位线（最可靠，不依赖单条消息匹配）
        final peerSeq = getPeerReadSeq(conversationID);
        if (peerSeq > 0 && (msg.seq ?? 0) > 0 && (msg.seq ?? 0) <= peerSeq) {
          msg.isRead = true;
          msg.hasReadTime ??= DateTime.now().millisecondsSinceEpoch;
          continue;
        }
        // 策略 B: 从 message DB 缓存恢复
        if (msg.clientMsgID != null) {
          final cached = getMessage(conversationID, msg.clientMsgID!);
          if (cached != null && cached.isRead == true) {
            msg.isRead = true;
            msg.hasReadTime = cached.hasReadTime;
            continue;
          }
          // 策略 C: 从 meta DB 的 conversationID-independent 记录恢复
          final metaReadTime = getMessageReadTime(msg.clientMsgID!);
          if (metaReadTime != null) {
            msg.isRead = true;
            msg.hasReadTime = metaReadTime;
          }
        }
      }
    }
  }

  // ======================= Voice Played =======================

  /// 获取已播放语音消息ID集合
  static Set<String> getPlayedVoiceIds(String conversationID) {
    final raw = _metaBox?.get('playedVoice:$conversationID');
    if (raw is String && raw.isNotEmpty) {
      return Set<String>.from(jsonDecode(raw));
    }
    return {};
  }

  /// 标记语音消息为已播放
  static Future<void> addPlayedVoiceId(String conversationID, String msgId) async {
    final ids = getPlayedVoiceIds(conversationID);
    ids.add(msgId);
    await _metaBox?.put('playedVoice:$conversationID', jsonEncode(ids.toList()));
  }

  // ======================= Meta =======================

  /// 获取 maxSeq
  static int getMaxSeq(String conversationID) {
    return _metaBox?.get('maxSeq:$conversationID') ?? 0;
  }

  /// 设置 maxSeq
  static Future<void> setMaxSeq(String conversationID, int seq) async {
    final current = getMaxSeq(conversationID);
    if (seq > current) {
      await _metaBox?.put('maxSeq:$conversationID', seq);
    }
  }
}
