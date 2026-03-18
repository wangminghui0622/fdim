import 'dart:async';

import 'package:flutter/foundation.dart';
import 'package:get/get.dart';
import 'package:rxdart/rxdart.dart' as rx;

import '../../sdk/flutter_openim_sdk.dart';
import '../config.dart';
import '../ws_client.dart';

/// IMController - mirrors the official OpenIM Flutter demo's im_controller.dart
/// Uses OpenIM.iMManager (our pure-Dart SDK) for all IM operations.
class IMController extends GetxController {
  final conversations = <ConversationInfo>[].obs;
  final totalUnreadCount = 0.obs;
  final friendApplyCount = 0.obs;
  final isLoggedIn = false.obs;
  final wsStatus = WsStatus.disconnected.obs;

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
    isLoggedIn.value = true;
    return info;
  }

  /// Logout via SDK
  Future<void> logoutSDK() async {
    await OpenIM.iMManager.logout();
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
}
