import 'package:flutter/foundation.dart';

import '../../flutter_openim_sdk.dart';
import '../../../core/http_client.dart';
import '../../../core/config.dart';

class ConversationManager {
  late OnConversationListener listener;

  Future setConversationListener(OnConversationListener listener) async {
    this.listener = listener;
  }

  Future<List<ConversationInfo>> getAllConversationList({
    String? operationID,
  }) async {
    final data = await HttpClient.post('/conversation/get_all_conversations', data: {
      'ownerUserID': Config.userID,
    });
    return _parseConversationList(data);
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
    return _parseConversationList(data);
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
    final list = await getAllConversationList();
    int total = 0;
    for (final c in list) {
      total += c.unreadCount;
    }
    return total;
  }

  Future markConversationMessageAsRead({
    required String conversationID,
    String? operationID,
  }) async {
    try {
      await HttpClient.post('/msg/mark_conversation_as_read', data: {
        'userID': Config.userID,
        'conversationID': conversationID,
        'hasReadSeq': 0,
        'seqs': [],
      }, showErrorToast: false);

      // 与官方体验保持一致：标记成功后立即刷新会话列表与总未读数，
      // 不完全依赖服务端通知链路，避免本地 WS/通知时序差异导致红点残留。
      listener.conversationChanged([]);
      try {
        final total = await getTotalUnreadMsgCount();
        listener.totalUnreadMessageCountChanged(total is int ? total : 0);
      } catch (e) {
        debugPrint('[SDK] refresh total unread after markConversationMessageAsRead error: $e');
      }
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

            final showName = (msgMap['senderNickname'] ??
                    msgMap['senderName'] ??
                    msgMap['showName'] ??
                    '')
                .toString();
            final faceURL = (msgMap['senderFaceUrl'] ??
                    msgMap['senderFaceURL'] ??
                    msgMap['faceURL'] ??
                    '')
                .toString();
            final latestMsgSendTime = msgMap['sendTime'] ??
                msgMap['latestMsgSendTime'] ??
                msgMap['latestMsgRecvTime'] ??
                msgMap['LatestMsgRecvTime'] ??
                0;

            final conv = ConversationInfo.fromJson({
              'conversationID': e['conversationID'],
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
