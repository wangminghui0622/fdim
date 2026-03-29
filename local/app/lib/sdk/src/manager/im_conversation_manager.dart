import 'package:flutter/foundation.dart';

import '../../fdim_sdk.dart';
import '../../../core/http_client.dart';
import '../../../core/config.dart';
import '../../../core/local_store.dart';

class ConversationManager {
  late OnConversationListener listener;

  Future setConversationListener(OnConversationListener listener) async {
    this.listener = listener;
  }

  Future<List<ConversationInfo>> getAllConversationList({
    String? operationID,
  }) async {
    // 与官方一致：优先返回本地缓存，同时后台同步服务端
    final local = LocalStore.getAllConversations();
    if (local.isNotEmpty) {
      // 后台同步最新数据
      _syncFromServerInBackground();
      return local;
    }
    // 本地为空时从服务端拉取并写入本地
    final data = await HttpClient.post('/conversation/get_all_conversations', data: {
      'ownerUserID': Config.userID,
    });
    final list = _parseConversationList(data);
    await _syncConversationSeqs(list);
    await LocalStore.putConversations(list);
    return list;
  }

  void _syncFromServerInBackground() async {
    try {
      final data = await HttpClient.post('/conversation/get_all_conversations', data: {
        'ownerUserID': Config.userID,
      }, showErrorToast: false);
      final list = _parseConversationList(data);
      await _syncConversationSeqs(list);
      await LocalStore.putConversations(list);
    } catch (_) {}
  }

  Future<List<ConversationInfo>> getConversationListSplit({
    int offset = 0,
    int count = 20,
    String? operationID,
  }) async {
    final data = await HttpClient.post('/conversation/get_sorted_conversation_list', data: {
      'userID': Config.userID,
      'conversationIDs': <String>[],
      'pagination': {'pageNumber': offset ~/ count + 1, 'showNumber': count},
    });
    final list = _parseConversationList(data);
    await _syncConversationSeqs(list);
    // 与官方一致：拉取结果写入本地 DB（putConversations 会保留本地 unreadCount）
    await LocalStore.putConversations(list);
    // 返回全部本地会话（包括本地创建但服务端尚无的会话，如好友通过后的会话）
    return LocalStore.getAllConversations();
  }

  Future<ConversationInfo> getOneConversation({
    required String sourceID,
    required int sessionType,
    String? operationID,
  }) async {
    final data = await HttpClient.post('/conversation/get_conversations', data: {
      'ownerUserID': Config.userID,
      'conversationIDs': [_buildConversationID(sourceID, sessionType)],
    });
    final list = _parseConversationList(data);
    if (list.isNotEmpty) return list.first;

    // Return a placeholder if not found
    return ConversationInfo(
      conversationID: _buildConversationID(sourceID, sessionType),
      conversationType: sessionType,
      userID: sessionType == ConversationType.single ? sourceID : null,
      groupID: sessionType != ConversationType.single ? sourceID : null,
    );
  }

  Future<List<ConversationInfo>> getMultipleConversation({
    required List<String> conversationIDList,
    String? operationID,
  }) async {
    final data = await HttpClient.post('/conversation/get_conversations', data: {
      'ownerUserID': Config.userID,
      'conversationIDs': conversationIDList,
    });
    return _parseConversationList(data);
  }

  Future setConversationDraft({
    required String conversationID,
    required String draftText,
    String? operationID,
  }) async {
    await HttpClient.post('/conversation/set_conversations', data: {
      'ownerUserID': Config.userID,
      'conversations': [{
        'ownerUserID': Config.userID,
        'conversationID': conversationID,
        'draftText': draftText,
        'draftTextTime': DateTime.now().millisecondsSinceEpoch ~/ 1000,
      }],
    });
  }

  Future pinConversation({
    required String conversationID,
    required bool isPinned,
    String? operationID,
  }) async {
    await setConversation(conversationID, ConversationReq(isPinned: isPinned));
  }

  Future hideConversation({
    required String conversationID,
    String? operationID,
  }) async {
    // Not directly supported in our API yet - mark as hidden locally
  }

  Future<dynamic> getConversationIDBySessionType({
    required String sourceID,
    required int sessionType,
    String? operationID,
  }) async {
    return _buildConversationID(sourceID, sessionType);
  }

  Future<dynamic> getTotalUnreadMsgCount({String? operationID}) async {
    // 与官方一致：从本地 DB 计算总未读数（不走网络）
    return LocalStore.getTotalUnreadCount();
  }

  Future<void> _syncConversationSeqs(List<ConversationInfo> list) async {
    final conversationIDs = list
        .map((e) => e.conversationID)
        .where((e) => e.isNotEmpty)
        .toList();
    if (conversationIDs.isEmpty) return;

    try {
      final data = await HttpClient.post(
        '/msg/get_conversations_has_read_and_max_seq',
        data: {
          'userID': Config.userID,
          'conversationIDs': conversationIDs,
        },
        showErrorToast: false,
      );
      if (data is! Map || data['seqs'] is! Map) return;
      final seqs = data['seqs'] as Map;

      for (final conv in list) {
        final seqInfo = seqs[conv.conversationID];
        if (seqInfo is! Map) continue;
        final maxSeq = (seqInfo['maxSeq'] ?? seqInfo['MaxSeq'] ?? 0) as int;
        final hasReadSeq = (seqInfo['hasReadSeq'] ?? seqInfo['HasReadSeq'] ?? 0) as int;
        if (maxSeq > 0) {
          await LocalStore.setMaxSeq(conv.conversationID, maxSeq);
        }
        if (hasReadSeq > 0) {
          await LocalStore.setHasReadSeq(conv.conversationID, hasReadSeq);
        }
      }
    } catch (e) {
      debugPrint('[ConversationManager] _syncConversationSeqs error: $e');
    }
  }

  Future markConversationMessageAsRead({
    required String conversationID,
    String? operationID,
  }) async {
    // 与官方一致：
    // 1) 先在本地 DB 立即清零 unreadCount（UI 即时响应）
    // 2) 再通知服务端
    // 3) 推送包含真实数据的 ConversationInfo
    await LocalStore.clearUnread(conversationID);
    final updatedConv = LocalStore.getConversation(conversationID);
    if (updatedConv != null) {
      listener.conversationChanged([updatedConv]);
    }
    listener.totalUnreadMessageCountChanged(LocalStore.getTotalUnreadCount());

    try {
      final seqInfo = await FDIM.iMManager.messageManager.getHasReadAndMaxSeq(
        conversationID,
      );
      final maxSeq = seqInfo['maxSeq'] ?? LocalStore.getMaxSeq(conversationID);
      if (maxSeq > 0) {
        await LocalStore.setHasReadSeq(conversationID, maxSeq);
      }
      await HttpClient.post('/msg/mark_conversation_as_read', data: {
        'userID': Config.userID,
        'conversationID': conversationID,
        'hasReadSeq': maxSeq,
        'seqs': [],
      }, showErrorToast: false);
    } catch (e) {
      debugPrint('[SDK] markConversationMessageAsRead error: $e');
    }
  }

  Future deleteConversationAndDeleteAllMsg({
    required String conversationID,
    String? operationID,
  }) async {
    await HttpClient.post('/conversation/delete_conversations', data: {
      'ownerUserID': Config.userID,
      'conversationIDs': [conversationID],
    });
  }

  Future clearConversationAndDeleteAllMsg({
    required String conversationID,
    String? operationID,
  }) async {
    await HttpClient.post('/msg/clear_conversation_msg', data: {
      'userID': Config.userID,
      'conversationIDs': [conversationID],
    });
  }

  Future setConversation(
    String conversationID,
    ConversationReq req, {
    String? operationID,
  }) async {
    await HttpClient.post('/conversation/set_conversations', data: {
      'ownerUserID': Config.userID,
      'conversations': [{
        'ownerUserID': Config.userID,
        'conversationID': conversationID,
        ...req.toJson(),
      }],
    });
  }

  List<ConversationInfo> simpleSort(List<ConversationInfo> list) => list
    ..sort((a, b) {
      if ((a.isPinned == true && b.isPinned == true) ||
          (a.isPinned != true && b.isPinned != true)) {
        int aCompare = (a.draftTextTime ?? 0) > (a.latestMsgSendTime ?? 0)
            ? (a.draftTextTime ?? 0)
            : (a.latestMsgSendTime ?? 0);
        int bCompare = (b.draftTextTime ?? 0) > (b.latestMsgSendTime ?? 0)
            ? (b.draftTextTime ?? 0)
            : (b.latestMsgSendTime ?? 0);
        if (aCompare > bCompare) return -1;
        if (aCompare < bCompare) return 1;
        return 0;
      } else if (a.isPinned == true && b.isPinned != true) {
        return -1;
      } else {
        return 1;
      }
    });

  String get atAllTag => 'AtAllTag';

  // Helpers
  static String _buildConversationID(String sourceID, int sessionType) {
    if (sessionType == ConversationType.single) {
      // 【修复问题3】单聊会话ID必须按字典序排序用户ID（与后端一致）
      final ids = [Config.userID, sourceID];
      ids.sort(); // 按字典序排序
      return 'si_${ids[0]}_${ids[1]}';
    } else {
      return 'sg_$sourceID';
    }
  }

  List<ConversationInfo> _parseConversationList(dynamic data) {
    if (data == null) return [];
    // Try various response formats
    List? list;
    if (data is List) {
      list = data;
    } else if (data is Map) {
      list = data['conversations'] ??
          data['conversationElems'] ??
          data['unreadCounts'];
      // get_sorted_conversation_list returns {conversationElems: [{conversationID, unreadCount, msgInfo: {...}}]}
      if (list != null && list.isNotEmpty && list.first is Map) {
        if (list.first.containsKey('conversation')) {
          // Old format: {conversation: {...}, unreadCount: ...}
          return list.map((e) {
            final conv = ConversationInfo.fromJson(e['conversation']);
            if (e['unreadCount'] != null) conv.unreadCount = e['unreadCount'];
            return conv;
          }).toList();
        } else if (list.first.containsKey('msgInfo')) {
          // New format: {conversationID, unreadCount, msgInfo: {...latest message fields...}}
          return list.map((e) {
            final msgInfo = e['msgInfo'];
            final msgMap = msgInfo is Map
                ? Map<String, dynamic>.from(msgInfo)
                : <String, dynamic>{};

            Message? latestMsg;
            try {
              latestMsg = msgMap.isNotEmpty ? Message.fromJson(msgMap) : null;
            } catch (err) {
              debugPrint('[SDK] parse latest msg from msgInfo error: $err');
            }

            // 从 conversationID 推导 conversationType、userID、groupID
            final convID = (e['conversationID'] ?? '').toString();
            int? conversationType;
            String? userID;
            String? groupID;
            String showName = '';
            String faceURL = '';

            if (convID.startsWith('sg_')) {
              // 群聊会话：使用 groupName 和 groupFaceURL（与官方一致）
              conversationType = ConversationType.superGroup;
              groupID = convID.substring(3);
              showName = (msgMap['groupName'] ??
                      msgMap['GroupName'] ??
                      msgMap['showName'] ??
                      groupID ??
                      '')
                  .toString();
              faceURL = (msgMap['groupFaceURL'] ??
                      msgMap['GroupFaceURL'] ??
                      msgMap['faceURL'] ??
                      '')
                  .toString();
            } else if (convID.startsWith('si_')) {
              // 单聊会话：使用 senderName 和 senderFaceURL
              conversationType = ConversationType.single;
              final parts = convID.split('_');
              if (parts.length == 3) {
                userID = (parts[1] == Config.userID) ? parts[2] : parts[1];
              }
              showName = (msgMap['senderNickname'] ??
                      msgMap['senderName'] ??
                      msgMap['showName'] ??
                      '')
                  .toString();
              faceURL = (msgMap['senderFaceUrl'] ??
                      msgMap['senderFaceURL'] ??
                      msgMap['faceURL'] ??
                      '')
                  .toString();
            } else {
              showName = (msgMap['senderNickname'] ??
                      msgMap['senderName'] ??
                      msgMap['showName'] ??
                      '')
                  .toString();
              faceURL = (msgMap['senderFaceUrl'] ??
                      msgMap['senderFaceURL'] ??
                      msgMap['faceURL'] ??
                      '')
                  .toString();
            }

            final latestMsgSendTime = msgMap['sendTime'] ??
                msgMap['latestMsgSendTime'] ??
                msgMap['latestMsgRecvTime'] ??
                msgMap['LatestMsgRecvTime'] ??
                0;

            final conv = ConversationInfo.fromJson({
              'conversationID': convID,
              'conversationType': conversationType,
              'userID': userID,
              'groupID': groupID,
              'recvMsgOpt': e['recvMsgOpt'],
              'unreadCount': e['unreadCount'],
              'isPinned': e['isPinned'],
              'showName': showName,
              'faceURL': faceURL,
              'latestMsg': latestMsg?.toJson(),
              'latestMsgSendTime': latestMsgSendTime,
            });
            return conv;
          }).toList();
        }
      }
    }
    if (list is List) {
      return list.map((e) => ConversationInfo.fromJson(e)).toList();
    }
    return [];
  }
}
