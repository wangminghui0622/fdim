import 'dart:async';
import 'dart:convert';

import 'package:flutter/foundation.dart';
import 'package:get/get.dart';
import 'package:rxdart/rxdart.dart' as rx;

import '../../models/signaling_info.dart';
import '../../sdk/flutter_openim_sdk.dart';
import '../config.dart';
import '../live_controller.dart';
import '../local_store.dart';
import '../ws_client.dart';

/// IMController - mirrors the official OpenIM Flutter demo's im_controller.dart
/// Uses OpenIM.iMManager (our pure-Dart SDK) for all IM operations.
class IMController extends GetxController with LiveController {
  final conversations = <ConversationInfo>[].obs;
  final totalUnreadCount = 0.obs;
  final friendApplyCount = 0.obs;
  final isLoggedIn = false.obs;
  final wsStatus = WsStatus.disconnected.obs;
  final currentChatConversationID = ''.obs;

  // RxDart subjects for cross-page event broadcasting (same as official IMCallback)
  final conversationAddedSubject = rx.PublishSubject<List<ConversationInfo>>();
  final conversationChangedSubject = rx.PublishSubject<List<ConversationInfo>>();
  final newMessageSubject = rx.PublishSubject<Message>();
  final friendApplicationSubject = rx.PublishSubject<FriendApplicationInfo>();
  final friendAddedSubject = rx.PublishSubject<FriendInfo>();
  final friendDeletedSubject = rx.PublishSubject<FriendInfo>();
  final friendInfoChangedSubject = rx.PublishSubject<FriendInfo>();
  final selfInfoUpdatedSubject = rx.PublishSubject<UserInfo>();
  final c2cReadReceiptSubject = rx.PublishSubject<List<ReadReceiptInfo>>();
  final msgRevokedSubject = rx.PublishSubject<RevokedInfo>();

  // Callback hooks for specific pages to override
  Function(Message msg)? onRecvNewMessage;
  Function(List<ReadReceiptInfo> list)? onRecvC2CReadReceipt;

  StreamSubscription? _wsStatusSub;

  @override
  void onInit() {
    super.onInit();
    _wsStatusSub = WsClient.statusStream.listen((status) {
      wsStatus.value = status;
    });
  }

  @override
  void onClose() {
    _wsStatusSub?.cancel();
    conversationAddedSubject.close();
    conversationChangedSubject.close();
    newMessageSubject.close();
    friendApplicationSubject.close();
    friendAddedSubject.close();
    friendDeletedSubject.close();
    friendInfoChangedSubject.close();
    selfInfoUpdatedSubject.close();
    c2cReadReceiptSubject.close();
    msgRevokedSubject.close();
    onCloseLive();
    OpenIM.iMManager.dispose();
    super.onClose();
  }

  /// Initialize SDK and set all listeners (call after login tokens are stored)
  Future<void> initOpenIM() async {
    await OpenIM.iMManager.initSDK(
      platformID: Config.platformID,
      apiAddr: Config.apiUrl,
      wsAddr: Config.wsUrl,
      dataDir: '',
      listener: OnConnectListener(
        onConnecting: () {
          debugPrint('[SDK] Connecting...');
        },
        onConnectSuccess: () {
          debugPrint('[SDK] Connected');
        },
        onConnectFailed: (code, msg) {
          debugPrint('[SDK] Connect failed: $code $msg');
        },
        onKickedOffline: () {
          debugPrint('[SDK] Kicked offline');
        },
        onUserTokenExpired: () {
          debugPrint('[SDK] Token expired');
        },
      ),
    );

    // Set conversation listener
    OpenIM.iMManager.conversationManager.setConversationListener(
      OnConversationListener(
        onConversationChanged: (list) {
          conversationChangedSubject.add(list);
        },
        onNewConversation: (list) {
          conversationAddedSubject.add(list);
        },
        onTotalUnreadMessageCountChanged: (count) {
          totalUnreadCount.value = count;
        },
      ),
    );

    // Set message listener
    OpenIM.iMManager.messageManager.setAdvancedMsgListener(
      OnAdvancedMsgListener(
        onRecvNewMessage: (msg) {
          // 自定义消息：区分信令（200-204）和通话结果（210）
          if (msg.contentType == MessageType.custom) {
            int? customType;
            try {
              final raw = msg.customElem?.data;
              if (raw != null) {
                customType = jsonDecode(raw)['customType'] as int?;
              }
            } catch (_) {}

            // 信令消息（200-204）：仅在线推送，不进聊天列表
            if (customType != null &&
                customType >= CustomMessageType.callingInvite &&
                customType <= CustomMessageType.callingHungup) {
              // 忽略自己发出的信令回显（push 推送给发送方的多设备同步）
              if (msg.sendID == Config.userID) return;
              handleSignalingMessage(msg);
              return;
            }
            // 通话结果（210）及其他自定义消息 → 走正常消息流
          }
          newMessageSubject.add(msg);
          onRecvNewMessage?.call(msg);
        },
        onRecvC2CReadReceipt: (list) {
          c2cReadReceiptSubject.add(list);
          onRecvC2CReadReceipt?.call(list);
        },
        onNewRecvMessageRevoked: (info) {
          msgRevokedSubject.add(info);
        },
      ),
    );

    // Set friendship listener
    // Official pattern: on ANY friend application event, re-fetch count from server
    // (never manual increment — matches HomeLogic.getUnhandledFriendApplicationCount)
    OpenIM.iMManager.friendshipManager.setFriendshipListener(
      OnFriendshipListener(
        onFriendApplicationAdded: (info) {
          loadFriendApplyCount();
          friendApplicationSubject.add(info);
          debugPrint('[SDK] New friend application from ${info.fromUserID}');
        },
        onFriendApplicationAccepted: (info) {
          loadFriendApplyCount();
          friendApplicationSubject.add(info);
          // Trigger friend list refresh when application is accepted
          // Convert FriendApplicationInfo to FriendInfo to trigger friendAddedSubject
          final friendInfo = FriendInfo(
            friendUserID: info.fromUserID,
            nickname: info.fromNickname,
            faceURL: info.fromFaceURL,
          );
          friendAddedSubject.add(friendInfo);

          // 确保双方都有会话：判断对方是谁
          _ensureConversationOnAccepted(info);
        },
        onFriendApplicationRejected: (info) {
          loadFriendApplyCount();
          friendApplicationSubject.add(info);
        },
        onFriendAdded: (info) {
          friendAddedSubject.add(info);
        },
        onFriendDeleted: (info) {
          friendDeletedSubject.add(info);
        },
        onFriendInfoChanged: (info) {
          friendInfoChangedSubject.add(info);
        },
      ),
    );

    // Set user listener
    OpenIM.iMManager.userManager.setUserListener(
      OnUserListener(
        onSelfInfoUpdated: (info) {
          selfInfoUpdatedSubject.add(info);
          Config.nickname = info.nickname ?? '';
          Config.faceURL = info.faceURL ?? '';
        },
      ),
    );

    // Set group listener
    OpenIM.iMManager.groupManager.setGroupListener(OnGroupListener());
  }

  /// Login via SDK
  Future<UserInfo> loginSDK() async {
    final info = await OpenIM.iMManager.login(
      userID: Config.userID,
      token: Config.token,
    );
    // 与官方一致：登录后初始化本地 DB
    await LocalStore.init();

    try {
      await OpenIM.iMManager.conversationManager.getConversationListSplit(
        offset: 0,
        count: 100,
      );
    } catch (e) {
      debugPrint('[IMController] initial conversation sync failed: $e');
    }

    totalUnreadCount.value = LocalStore.getTotalUnreadCount();
    isLoggedIn.value = true;
    onInitLive();
    return info;
  }

  /// Logout via SDK
  Future<void> logoutSDK() async {
    await OpenIM.iMManager.logout();
    await LocalStore.close();
    isLoggedIn.value = false;
  }

  /// Load pending friend application count from server.
  /// Called after login to initialize badge, and on WS reconnect.
  Future<void> loadFriendApplyCount() async {
    try {
      final list = await OpenIM.iMManager.friendshipManager
          .getFriendApplicationListAsRecipient();
      final pending = list.where((e) => e.handleResult == 0).length;
      friendApplyCount.value = pending;
    } catch (_) {}
  }

  void clearFriendApplyBadge() {
    friendApplyCount.value = 0;
  }

  /// 好友申请通过后，在 A（申请者）本地创建两条消息并创建会话：
  /// 1) A 的申请语句（来自 A 自己，右侧）
  /// 2) 默认问候（来自 B，左侧）
  Future<void> _ensureConversationOnAccepted(FriendApplicationInfo info) async {
    try {
      final myID = Config.userID;

      if (info.fromUserID != myID) {
        // 我是 B（通过者），会话已由 _createLocalFirstMessage 创建，跳过
        return;
      }

      // 我是 A（申请者），对方是 B（通过者）
      final otherID = info.toUserID ?? '';
      if (otherID.isEmpty) return;

      final ids = [myID, otherID]..sort();
      final convID = 'si_${ids[0]}_${ids[1]}';

      // 如果会话已存在，不重复创建
      if (LocalStore.getConversation(convID) != null) return;

      // 推送中不含 reqMsg，需要从服务端获取完整的好友申请信息
      String reqMsg = info.reqMsg ?? '';
      if (reqMsg.isEmpty) {
        try {
          final applies = await OpenIM.iMManager.friendshipManager
              .getFriendApplicationListAsApplicant();
          for (final a in applies) {
            if (a.toUserID == otherID && (a.reqMsg ?? '').isNotEmpty) {
              reqMsg = a.reqMsg!;
              break;
            }
          }
        } catch (_) {}
      }
      debugPrint('[IMController] reqMsg for $otherID: "$reqMsg"');

      // 查询 B 的昵称和头像（推送中 toNickname 可能为空）
      String otherNickname = otherID;
      String? otherFaceURL;
      if ((info.toNickname ?? '').isNotEmpty) {
        otherNickname = info.toNickname!;
        otherFaceURL = info.toFaceURL;
      } else {
        try {
          final users = await OpenIM.iMManager.userManager.getUsersInfo(
            userIDList: [otherID],
          );
          if (users.isNotEmpty) {
            otherNickname = users.first.nickname ?? otherID;
            otherFaceURL = users.first.faceURL;
          }
        } catch (_) {}
      }

      final now = DateTime.now().millisecondsSinceEpoch;
      int unread = 0;
      Message? latestMsg;

      // 消息1: A 的申请语句（如果有）— A 自己发的，A 聊天页右侧
      if (reqMsg.isNotEmpty) {
        final msg1 = Message(
          clientMsgID: '${myID}_${now}_req',
          sendID: myID,
          recvID: otherID,
          senderNickname: Config.nickname.isNotEmpty ? Config.nickname : myID,
          senderFaceUrl: Config.faceURL,
          sessionType: ConversationType.single,
          contentType: MessageType.text,
          createTime: now,
          sendTime: now,
          status: MessageStatus.succeeded,
          isRead: true,
          textElem: TextElem(content: reqMsg),
        );
        await LocalStore.putMessage(convID, msg1);
        latestMsg = msg1;
      }

      // 消息2: 默认问候 — 来自 B，A 聊天页左侧
      final msg2 = Message(
        clientMsgID: '${otherID}_${now}_greet',
        sendID: otherID,
        recvID: myID,
        senderNickname: otherNickname,
        senderFaceUrl: otherFaceURL,
        sessionType: ConversationType.single,
        contentType: MessageType.text,
        createTime: now + 1,
        sendTime: now + 1,
        status: MessageStatus.succeeded,
        isRead: false,
        textElem: TextElem(content: '我们已经是好友了，可以开始聊天了'),
      );
      await LocalStore.putMessage(convID, msg2);
      latestMsg = msg2;
      unread++;

      // 创建会话（用 B 的昵称）
      final conv = ConversationInfo(
        conversationID: convID,
        conversationType: ConversationType.single,
        userID: otherID,
        showName: otherNickname,
        faceURL: otherFaceURL,
        latestMsg: latestMsg,
        latestMsgSendTime: now + 1,
      );
      conv.unreadCount = unread;
      await LocalStore.putConversation(conv);

      // 通知 UI 刷新
      OpenIM.iMManager.conversationManager.listener
          .conversationChanged([conv]);
      OpenIM.iMManager.conversationManager.listener
          .totalUnreadMessageCountChanged(LocalStore.getTotalUnreadCount());

      debugPrint('[IMController] Created conversation $convID with 2 messages for friend $otherNickname');
    } catch (e) {
      debugPrint('[IMController] _ensureConversationOnAccepted error: $e');
    }
  }
}
