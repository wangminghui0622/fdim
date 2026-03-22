import 'dart:async';
import 'dart:convert';
import 'package:flutter/foundation.dart';

import '../../flutter_openim_sdk.dart';
import '../../../core/http_client.dart';
import '../../../core/config.dart';
import '../../../core/ws_client.dart';
import '../../../core/local_store.dart';

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

  /// 与官方 Go SDK 一致：当前正在查看的会话 ID，
  /// 该会话的新消息不递增 unreadCount
  String? _activeConversationID;

  /// ChatPage 进入/离开时调用
  void setActiveConversation(String? conversationID) {
    _activeConversationID = conversationID;
  }

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
      // Regular message — 与官方一致：写入本地 DB 后再通知 UI
      final msg = Message.fromJson(msgData);
      final convID = _resolveConversationID(msg);
      if (convID.isNotEmpty) {
        // 0) 确保会话存在于本地 DB（新会话首条消息时创建骨架）
        await _ensureConversationExists(convID, msg);
        // 1) 写入本地消息 DB
        await LocalStore.putMessage(convID, msg);
        // 2) 更新会话 latestMsg
        await LocalStore.updateLatestMsg(convID, msg);
        // 3) 与官方 Go SDK 一致：
        //    - 不是自己发的消息
        //    - 不是当前正在查看的会话
        //    才递增 unreadCount
        if (msg.sendID != userID && convID != _activeConversationID) {
          await LocalStore.incrementUnread(convID);
        }
      }
      messageManager.msgListener.recvNewMessage(msg);
      // 与官方一致：推送包含真实数据的 ConversationInfo 列表
      _notifyConversationChanged(convID);
      // 更新总未读数（从本地 DB 计算，无需网络请求）
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
      case MessageType.burnAfterReadingNotification: // 1701
      case 1702: // ConversationUnreadNotification
      case 1703: // ClearConversationNotification
      case 1704: // ConversationDeleteNotification
        // 与官方一致：从服务端拉取最新会话列表 → 写入本地 DB → 推送真实数据
        await _syncConversationsFromServer();
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
        if (detail != null) {
          final info = ReadReceiptInfo.fromJson(detail);
          final convID = info.conversationID ?? '';
          final readTime = DateTime.now().millisecondsSinceEpoch;

          debugPrint('[SDK] hasReadReceipt: convID=$convID, '
              'hasReadSeq=${info.hasReadSeq}, '
              'seqs=${info.seqs}, '
              'msgIDList=${info.msgIDList}');

          // 策略1: hasReadSeq 水位线（最常见的服务端格式）
          // 含义：对方已读到该 seq → 我发的 seq<=该值 的消息都已被对方阅读
          if (info.hasReadSeq != null && info.hasReadSeq! > 0 && convID.isNotEmpty) {
            await LocalStore.setPeerReadSeq(convID, info.hasReadSeq!);
            debugPrint('[SDK] hasReadReceipt: stored peerReadSeq=${info.hasReadSeq} for $convID');
          }

          // 策略2: seqs → 转换为 msgIDList
          if ((info.msgIDList == null || info.msgIDList!.isEmpty) &&
              info.seqs != null &&
              info.seqs!.isNotEmpty &&
              convID.isNotEmpty) {
            await _convertSeqsToMsgIDList(info);
          }

          // 策略3: 按 msgIDList 逐条标记已读（双写到 meta box）
          if (info.msgIDList != null && info.msgIDList!.isNotEmpty && convID.isNotEmpty) {
            await LocalStore.markMessagesRead(convID, info.msgIDList!, readTime);
            debugPrint('[SDK] hasReadReceipt: marked ${info.msgIDList!.length} msgs read');
          }

          if (convID.startsWith('sg_')) {
            messageManager.msgListener.recvGroupReadReceipt([info]);
          } else {
            messageManager.msgListener.recvC2CReadReceipt([info]);
          }
        }
        break;

      default:
        // Unknown notification - still dispatch as message
        final msg = Message.fromJson(msgData);
        messageManager.msgListener.recvNewMessage(msg);
        break;
    }
  }

  /// 与官方一致：从本地 DB 计算总未读数并推送
  void _updateTotalUnreadCount() {
    final count = LocalStore.getTotalUnreadCount();
    conversationManager.listener.totalUnreadMessageCountChanged(count);
  }

  /// 从消息中解析 conversationID（与官方 SDK 行为一致）
  String _resolveConversationID(Message msg) {
    final st = msg.sessionType;
    if (st == ConversationType.single) {
      final ids = [msg.sendID ?? '', msg.recvID ?? '']..sort();
      if (ids[0].isEmpty || ids[1].isEmpty) return '';
      return 'si_${ids[0]}_${ids[1]}';
    }
    if (st == ConversationType.group || st == ConversationType.superGroup) {
      final gid = msg.groupID ?? '';
      if (gid.isEmpty) return '';
      return 'sg_$gid';
    }
    return '';
  }

  /// 确保会话存在于本地 DB（新会话首条消息时创建骨架）
  Future<void> _ensureConversationExists(String convID, Message msg) async {
    final existing = LocalStore.getConversation(convID);
    if (existing != null) return;
    // 创建骨架 ConversationInfo
    final st = msg.sessionType;
    String? cUserID;
    String? cGroupID;
    String? showName;
    String? faceURL;
    if (st == ConversationType.single) {
      // 单聊：对方的信息
      cUserID = (msg.sendID == userID) ? msg.recvID : msg.sendID;
      if (msg.sendID != userID) {
        // 收到别人的消息：直接用发送者昵称和头像
        showName = msg.senderNickname ?? msg.sendID ?? '';
        faceURL = msg.senderFaceUrl;
      } else {
        // 自己发的消息：尝试从服务端查询对方昵称
        showName = msg.recvID ?? '';
        try {
          final users = await userManager.getUsersInfo(userIDList: [cUserID ?? '']);
          if (users.isNotEmpty) {
            showName = users.first.nickname ?? showName;
            faceURL = users.first.faceURL;
          }
        } catch (_) {}
      }
    } else {
      cGroupID = msg.groupID;
      showName = msg.groupID ?? '';
    }
    final conv = ConversationInfo(
      conversationID: convID,
      conversationType: st,
      userID: cUserID,
      groupID: cGroupID,
      showName: showName,
      faceURL: faceURL,
      latestMsg: msg,
      latestMsgSendTime: msg.sendTime,
    );
    await LocalStore.putConversation(conv);
  }

  /// 通知会话变更（推送真实 ConversationInfo，与官方一致）
  void _notifyConversationChanged(String conversationID) {
    if (conversationID.isEmpty) {
      conversationManager.listener.conversationChanged([]);
      return;
    }
    final conv = LocalStore.getConversation(conversationID);
    if (conv != null) {
      conversationManager.listener.conversationChanged([conv]);
    } else {
      conversationManager.listener.conversationChanged([]);
    }
  }

  /// 从服务端同步会话列表到本地 DB，然后推送
  Future<void> _syncConversationsFromServer() async {
    try {
      // getConversationListSplit 内部已调用 putConversations（会保留本地 unreadCount）
      final list = await conversationManager.getConversationListSplit(
        offset: 0, count: 100,
      );
      conversationManager.listener.conversationChanged(list);
    } catch (e) {
      debugPrint('[SDK] _syncConversationsFromServer error: $e');
      conversationManager.listener.conversationChanged([]);
    }
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
