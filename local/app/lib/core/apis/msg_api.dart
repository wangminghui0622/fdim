import 'dart:convert';
import 'package:uuid/uuid.dart';

import '../http_client.dart';
import '../config.dart';
import '../models/message_info.dart';

class MsgApi {
  static const _uuid = Uuid();

  static Future<int> getNewestSeq(String conversationID) async {
    final data = await HttpClient.post('/msg/newest_seq', data: {
      'userID': Config.userID,
      'conversationID': conversationID,
    });
    if (data == null) return 0;
    return (data['maxSeq'] ?? 0) as int;
  }

  static Future<List<MessageInfo>> pullMsgBySeq({
    required String conversationID,
    required List<int> seqs,
  }) async {
    final data = await HttpClient.post('/msg/pull_msg_by_seq', data: {
      'userID': Config.userID,
      'conversationID': conversationID,
      'seqs': seqs,
    });
    if (data == null) return [];
    final msgs = data['msgs'] ?? data;
    if (msgs is List) {
      return msgs
          .map((e) => MessageInfo.fromJson(e as Map<String, dynamic>))
          .toList();
    }
    return [];
  }

  static Future<MessageInfo?> sendMsg({
    required String recvID,
    required String groupID,
    required int sessionType,
    required int contentType,
    required String content,
  }) async {
    final clientMsgID = _uuid.v4();
    final data = await HttpClient.post('/msg/send_msg', data: {
      'recvID': recvID,
      'sendMsg': {
        'sendID': Config.userID,
        'recvID': recvID,
        'groupID': groupID,
        'senderNickname': Config.nickname,
        'senderFaceURL': Config.faceURL,
        'senderPlatformID': Config.platformID,
        'content': jsonEncode({'content': content}),
        'contentType': contentType,
        'sessionType': sessionType,
        'clientMsgID': clientMsgID,
        'sendTime': DateTime.now().millisecondsSinceEpoch,
      },
    });
    if (data != null && data is Map<String, dynamic>) {
      return MessageInfo.fromJson(data);
    }
    return null;
  }

  static Future<void> revokeMsg({
    required String conversationID,
    required int seq,
  }) async {
    await HttpClient.post('/msg/revoke_msg', data: {
      'userID': Config.userID,
      'conversationID': conversationID,
      'seq': seq,
    });
  }

  static Future<void> markMsgsAsRead({
    required String conversationID,
    required List<int> seqs,
  }) async {
    await HttpClient.post('/msg/mark_msgs_as_read', data: {
      'userID': Config.userID,
      'conversationID': conversationID,
      'seqs': seqs,
    }, showErrorToast: false);
  }

  static Future<void> clearConversationMsg(String conversationID) async {
    await HttpClient.post('/msg/clear_conversation_msg', data: {
      'userID': Config.userID,
      'conversationIDs': [conversationID],
    });
  }

  /// 获取群消息已读/未读成员列表
  static Future<Map<String, List<GroupMsgReadUser>>> getGroupMsgReadUsers({
    required String groupID,
    required String conversationID,
    required int seq,
  }) async {
    final data = await HttpClient.post('/msg/get_group_msg_read_users', data: {
      'groupID': groupID,
      'conversationID': conversationID,
      'seq': seq,
    });
    if (data == null) {
      return {'readUsers': [], 'unreadUsers': []};
    }
    final readUsers = (data['readUsers'] as List? ?? [])
        .map((e) => GroupMsgReadUser.fromJson(e as Map<String, dynamic>))
        .toList();
    final unreadUsers = (data['unreadUsers'] as List? ?? [])
        .map((e) => GroupMsgReadUser.fromJson(e as Map<String, dynamic>))
        .toList();
    return {'readUsers': readUsers, 'unreadUsers': unreadUsers};
  }
}

/// 群消息已读用户信息
class GroupMsgReadUser {
  final String userID;
  final String nickname;
  final String faceURL;

  GroupMsgReadUser({
    required this.userID,
    required this.nickname,
    required this.faceURL,
  });

  factory GroupMsgReadUser.fromJson(Map<String, dynamic> json) {
    return GroupMsgReadUser(
      userID: json['userID'] ?? '',
      nickname: json['nickname'] ?? '',
      faceURL: json['faceURL'] ?? '',
    );
  }
}
