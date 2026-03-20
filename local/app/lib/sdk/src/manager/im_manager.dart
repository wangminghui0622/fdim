import 'dart:async';
import 'dart:convert';
import 'package:flutter/foundation.dart';

import '../../flutter_openim_sdk.dart';
import '../../../core/http_client.dart';
import '../../../core/config.dart';
import '../../../core/ws_client.dart';

class IMManager {
  late ConversationManager conversationManager;
  late FriendshipManager friendshipManager;
  late MessageManager messageManager;
  late GroupManager groupManager;
  late UserManager userManager;

  late OnConnectListener _connectListener;

  late String userID;
  late UserInfo userInfo;
  bool isLogined = false;
  String? token;

  StreamSubscription? _wsStatusSub;
  StreamSubscription? _wsMessageSub;

  IMManager() {
    conversationManager = ConversationManager();
    friendshipManager = FriendshipManager();
    messageManager = MessageManager();
    groupManager = GroupManager();
    userManager = UserManager();
  }

  /// Initialize the SDK
  Future<dynamic> initSDK({
    required int platformID,
    required String apiAddr,
    required String wsAddr,
    required String dataDir,
    required OnConnectListener listener,
    int logLevel = 6,
    bool isNeedEncryption = false,
    bool isCompression = false,
    bool isLogStandardOutput = true,
    String? logFilePath,
    String? operationID,
  }) async {
    _connectListener = listener;

    // Update config with provided addresses
    await Config.setServer(apiAddr, wsAddr);
    HttpClient.updateBaseUrl(apiAddr);

    // Listen to WS status changes
    _wsStatusSub?.cancel();
    _wsStatusSub = WsClient.statusStream.listen((status) {
      switch (status) {
        case WsStatus.connecting:
          _connectListener.connecting();
          break;
        case WsStatus.connected:
          _connectListener.connectSuccess();
          break;
        case WsStatus.disconnected:
          if (isLogined) {
            _connectListener.connectFailed(0, 'Connection lost');
          }
          break;
      }
    });

    // Listen to WS messages and dispatch to managers
    _wsMessageSub?.cancel();
    _wsMessageSub = WsClient.messageStream.listen(_handleWsMessage);

    return true;
  }

  /// Login
  Future<UserInfo> login({
    required String userID,
    required String token,
    String? operationID,
    Future<UserInfo> Function()? defaultValue,
    bool checkLoginStatus = true,
  }) async {
    this.isLogined = true;
    this.userID = userID;
    this.token = token;

    // Store in config
    Config.userID = userID;
    Config.token = token;

    // Connect WebSocket
    WsClient.connect();

    try {
      return this.userInfo = await userManager.getSelfUserInfo();
    } catch (error, stackTrace) {
      debugPrint('login e: $error  s: $stackTrace');
      if (null != defaultValue) {
        return this.userInfo = await (defaultValue.call());
      }
      // Return basic info
      return this.userInfo = UserInfo(userID: userID);
    }
  }

  /// Logout
  Future<dynamic> logout({String? operationID}) async {
    WsClient.disconnect();
    isLogined = false;
    token = null;
    await Config.clearLogin();
  }

  /// Get login status
  Future<int?> getLoginStatus({String? operationID}) async {
    if (!isLogined) return LoginStatus.logout;
    if (token == null || token!.isEmpty) return LoginStatus.logout;
    return LoginStatus.logged;
  }

  Future<String> getLoginUserID() async => userID;
  Future<UserInfo> getLoginUserInfo() async => userInfo;

  void dispose() {
    _wsStatusSub?.cancel();
    _wsMessageSub?.cancel();
  }

  /// Handle incoming WebSocket messages (official OpenIM WsResp protocol)
  void _handleWsMessage(WsResp resp) {
    try {
      switch (resp.reqIdentifier) {
        case WsProtocol.wsPushMsg:
          _handlePushMsg(resp.data);
          break;
        case WsProtocol.wsKickOnlineMsg:
          _connectListener.kickedOffline();
          break;
        case WsProtocol.wsLogoutMsg:
          _connectListener.userTokenExpired();
          break;
        default:
          debugPrint('[SDK] Unhandled reqIdentifier: ${resp.reqIdentifier}');
      }
    } catch (e) {
      debugPrint('[SDK] Error dispatching WsResp reqIdentifier=${resp.reqIdentifier}: $e');
    }
  }

  /// Handle WSPushMsg (2001) - dispatch by MsgData.contentType
  Future<void> _handlePushMsg(dynamic data) async {
    if (data == null) return;

    // data is MsgData JSON (may be raw bytes or parsed map)
    Map<String, dynamic> msgData;
    if (data is Map<String, dynamic>) {
      msgData = data;
    } else if (data is String) {
      msgData = jsonDecode(data) as Map<String, dynamic>;
    } else if (data is List<int>) {
      msgData = jsonDecode(utf8.decode(data)) as Map<String, dynamic>;
    } else {
      debugPrint('[SDK] Unknown push data type: ${data.runtimeType}');
      return;
    }

    final contentType = msgData['contentType'] as int? ?? 0;
    final content = msgData['content'];
    
    // 调试日志：打印所有收到的消息类型
    debugPrint('[SDK] Received push message: contentType=$contentType, isNotification=${MessageType.isNotificationType(contentType)}');

    // Decode notification detail from content field
    // Backend wraps tips in NotificationElem: {"detail": "<tips_json_string>"}
    Map<String, dynamic>? detail;
    if (content is String && content.isNotEmpty) {
      try {
        final parsed = jsonDecode(content);
        if (parsed is Map<String, dynamic>) {
          // Unwrap NotificationElem: extract inner "detail" JSON string
          final innerDetail = parsed['detail'];
          if (innerDetail is String && innerDetail.isNotEmpty) {
            detail = jsonDecode(innerDetail) as Map<String, dynamic>?;
          } else {
            detail = parsed;
          }
        }
      } catch (_) {
        // content is not JSON (e.g. text message)
      }
    }

    if (!MessageType.isNotificationType(contentType)) {
      // Regular message
      final msg = Message.fromJson(msgData);
      messageManager.msgListener.recvNewMessage(msg);
      // Trigger conversation list refresh to update unread counts
      conversationManager.listener.conversationChanged([]);
      // Update total unread count
      _updateTotalUnreadCount();
      return;
    }

    // Notification messages - dispatch by contentType
    switch (contentType) {
      // ===== Friend notifications =====
      case MessageType.friendApplicationApprovedNotification:
        if (detail != null) {
          final info = FriendApplicationInfo.fromJson(_flattenFriendApplicationTips(detail));
          friendshipManager.listener.friendApplicationAccepted(info);
        }
        break;
      case MessageType.friendApplicationRejectedNotification:
        if (detail != null) {
          final info = FriendApplicationInfo.fromJson(_flattenFriendApplicationTips(detail));
          friendshipManager.listener.friendApplicationRejected(info);
        }
        break;
      case MessageType.friendApplicationNotification:
        if (detail != null) {
          final info = FriendApplicationInfo.fromJson(_flattenFriendApplicationTips(detail));
          friendshipManager.listener.friendApplicationAdded(info);
        }
        break;
      case MessageType.friendAddedNotification:
        if (detail != null) {
          final info = FriendInfo.fromJson(detail['friend'] ?? detail);
          friendshipManager.listener.friendAdded(info);
        }
        break;
      case MessageType.friendDeletedNotification:
        if (detail != null) {
          final info = FriendInfo.fromJson(detail['friend'] ?? detail);
          friendshipManager.listener.friendDeleted(info);
        }
        break;
      case MessageType.friendRemarkSetNotification:
        if (detail != null) {
          final info = FriendInfo.fromJson(detail['friend'] ?? detail);
          friendshipManager.listener.friendInfoChanged(info);
        }
        break;
      case MessageType.blackAddedNotification:
        if (detail != null) {
          final info = BlacklistInfo.fromJson(detail['black'] ?? detail);
          friendshipManager.listener.blackAdded(info);
        }
        break;
      case MessageType.blackDeletedNotification:
        if (detail != null) {
          final info = BlacklistInfo.fromJson(detail['black'] ?? detail);
          friendshipManager.listener.blackDeleted(info);
        }
        break;
      // FriendInfoUpdated (1209/1210) - 好友信息更新通知
      case 1209:
      case 1210:
        if (detail != null) {
          final info = FriendInfo.fromJson(detail['friend'] ?? detail);
          friendshipManager.listener.friendInfoChanged(info);
        }
        break;

      // ===== Conversation notifications =====
      case MessageType.conversationChangeNotification:
      case MessageType.burnAfterReadingNotification: // 1701 - 阅后即焚/私聊设置变更
      case 1702: // ConversationUnreadNotification - 【修复问题2】会话未读数变更通知
      case 1703: // ClearConversationNotification
      case 1704: // ConversationDeleteNotification
        // 【修复问题2】触发会话列表刷新，确保未读数实时更新
        debugPrint('[SDK] 【DEBUG】Conversation notification received, contentType=$contentType');
        if (detail != null) {
          debugPrint('[SDK] 【DEBUG】Conversation notification detail: $detail');
        }
        debugPrint('[SDK] 【DEBUG】Triggering conversationChanged');
        conversationManager.listener.conversationChanged([]);
        // Update total unread count
        debugPrint('[SDK] 【DEBUG】Updating total unread count');
        _updateTotalUnreadCount();
        break;

      // ===== User notifications =====
      case MessageType.userInfoUpdatedNotification:
        if (detail != null) {
          final uid = detail['userID'] as String?;
          if (uid == userID) {
            // Self info updated - refresh
            userManager.getSelfUserInfo().then((info) {
              userInfo = info;
              userManager.listener.selfInfoUpdated(info);
            }).catchError((_) {});
          }
        }
        break;

      // ===== Group notifications =====
      case MessageType.groupCreatedNotification:
        if (detail != null) {
          final info = GroupInfo.fromJson(detail['group'] ?? detail);
          groupManager.listener.joinedGroupAdded(info);
        }
        break;
      case MessageType.groupInfoSetNotification:
      case MessageType.groupInfoSetAnnouncementNotification:
      case MessageType.groupInfoSetNameNotification:
        if (detail != null) {
          final info = GroupInfo.fromJson(detail['group'] ?? detail);
          groupManager.listener.groupInfoChanged(info);
        }
        break;
      case MessageType.joinGroupApplicationNotification:
        if (detail != null) {
          final info = GroupApplicationInfo.fromJson(detail);
          groupManager.listener.groupApplicationAdded(info);
        }
        break;
      case MessageType.groupApplicationAcceptedNotification:
        if (detail != null) {
          final info = GroupApplicationInfo.fromJson(detail);
          groupManager.listener.groupApplicationAccepted(info);
        }
        break;
      case MessageType.groupApplicationRejectedNotification:
        if (detail != null) {
          final info = GroupApplicationInfo.fromJson(detail);
          groupManager.listener.groupApplicationRejected(info);
        }
        break;
      case MessageType.memberQuitNotification:
        if (detail != null) {
          final info = GroupMembersInfo.fromJson(detail['quitUser'] ?? detail);
          groupManager.listener.groupMemberDeleted(info);
        }
        break;
      case MessageType.memberKickedNotification:
        if (detail != null) {
          final kicked = detail['kickedUserList'];
          if (kicked is List && kicked.isNotEmpty) {
            final info = GroupMembersInfo.fromJson(kicked.first);
            groupManager.listener.groupMemberDeleted(info);
          }
        }
        break;
      case MessageType.memberInvitedNotification:
      case MessageType.memberEnterNotification:
        if (detail != null) {
          final invited = detail['invitedUserList'] ?? detail['entrantUser'];
          if (invited is List && invited.isNotEmpty) {
            final info = GroupMembersInfo.fromJson(invited.first);
            groupManager.listener.groupMemberAdded(info);
          } else if (invited is Map<String, dynamic>) {
            final info = GroupMembersInfo.fromJson(invited);
            groupManager.listener.groupMemberAdded(info);
          }
        }
        break;
      case MessageType.groupOwnerTransferredNotification:
      case MessageType.groupMemberInfoChangedNotification:
      case MessageType.groupMemberSetToAdminNotification:
      case MessageType.groupMemberSetToOrdinaryUserNotification:
        if (detail != null) {
          final changed = detail['changedUser'] ?? detail['newGroupOwner'] ?? detail;
          if (changed is Map<String, dynamic>) {
            final info = GroupMembersInfo.fromJson(changed);
            groupManager.listener.groupMemberInfoChanged(info);
          }
        }
        break;
      case MessageType.groupMemberMutedNotification:
      case MessageType.groupMemberCancelMutedNotification:
        if (detail != null) {
          final mutedUser = detail['mutedUser'] ?? detail;
          if (mutedUser is Map<String, dynamic>) {
            final info = GroupMembersInfo.fromJson(mutedUser);
            groupManager.listener.groupMemberInfoChanged(info);
          }
        }
        break;
      case MessageType.groupMutedNotification:
      case MessageType.groupCancelMutedNotification:
        if (detail != null) {
          final info = GroupInfo.fromJson(detail['group'] ?? detail);
          groupManager.listener.groupInfoChanged(info);
        }
        break;
      case MessageType.dismissGroupNotification:
        if (detail != null) {
          final info = GroupInfo.fromJson(detail['group'] ?? detail);
          groupManager.listener.groupDismissed(info);
        }
        break;

      // ===== Message notifications =====
      case MessageType.revokeMessageNotification:
        if (detail != null) {
          final info = RevokedInfo.fromJson(detail);
          messageManager.msgListener.newRecvMessageRevoked(info);
        }
        break;
      case MessageType.hasReadReceipt:
        debugPrint('[SDK] Received hasReadReceipt notification, detail=$detail');
        if (detail != null) {
          final info = ReadReceiptInfo.fromJson(detail);
          debugPrint('[SDK] ReadReceiptInfo parsed: conversationID=${info.conversationID}, seqs=${info.seqs}, userID=${info.userID}, msgIDList=${info.msgIDList}');

          if ((info.msgIDList == null || info.msgIDList!.isEmpty) &&
              info.seqs != null &&
              info.seqs!.isNotEmpty &&
              info.conversationID != null &&
              info.conversationID!.isNotEmpty) {
            await _convertSeqsToMsgIDList(info);
          }

          debugPrint('[SDK] After conversion: msgIDList=${info.msgIDList}');

          final conversationID = info.conversationID ?? '';
          if (conversationID.startsWith('sg_')) {
            messageManager.msgListener.recvGroupReadReceipt([info]);
          } else {
            messageManager.msgListener.recvC2CReadReceipt([info]);
          }
        } else {
          debugPrint('[SDK] WARNING: hasReadReceipt detail is null!');
        }
        break;

      default:
        // Unknown notification - still dispatch as message
        final msg = Message.fromJson(msgData);
        messageManager.msgListener.recvNewMessage(msg);
        break;
    }
  }

  /// Update total unread count by fetching from server
  void _updateTotalUnreadCount() {
    conversationManager.getTotalUnreadMsgCount().then((count) {
      conversationManager.listener.totalUnreadMessageCountChanged(count);
    }).catchError((e) {
      debugPrint('[SDK] Failed to get total unread count: $e');
    });
  }


  Future<void> _convertSeqsToMsgIDList(ReadReceiptInfo info) async {
    try {
      if (info.seqs == null || info.seqs!.isEmpty || info.conversationID == null) {
        return;
      }

      final data = await HttpClient.post('/msg/pull_msg_by_seq', data: {
        'userID': Config.userID,
        'conversationID': info.conversationID,
        'seqs': info.seqs,
      });

      final msgIDList = <String>[];
      if (data is Map) {
        final msgsList = data['msgs'];
        if (msgsList is List) {
          for (final item in msgsList) {
            if (item is Map) {
              final msgList = item['Msgs'];
              if (msgList is List) {
                for (final msgJson in msgList) {
                  if (msgJson is Map) {
                    final clientMsgID = msgJson['clientMsgID'] as String?;
                    if (clientMsgID != null && clientMsgID.isNotEmpty) {
                      msgIDList.add(clientMsgID);
                    }
                  }
                }
              }
            }
          }
        }
      }

      info.msgIDList = msgIDList;
      debugPrint('[SDK] _convertSeqsToMsgIDList: converted ${info.seqs!.length} seqs to ${msgIDList.length} msgIDs');
    } catch (e) {
      debugPrint('[SDK] _convertSeqsToMsgIDList error: $e');
      info.msgIDList = [];
    }
  }

  /// Flatten backend FriendApplicationTips structure to flat FriendApplicationInfo fields.
  /// Backend sends: {"fromToUserID": {"fromUserID": "A", "toUserID": "B"}, "handleMsg": "..."}
  /// FriendApplicationInfo expects flat: {"fromUserID": "A", "toUserID": "B", "handleMsg": "..."}
  Map<String, dynamic> _flattenFriendApplicationTips(Map<String, dynamic> tips) {
    final result = Map<String, dynamic>.from(tips);
    final fromTo = tips['fromToUserID'];
    if (fromTo is Map<String, dynamic>) {
      result['fromUserID'] ??= fromTo['fromUserID'];
      result['toUserID'] ??= fromTo['toUserID'];
    }
    return result;
  }
}
