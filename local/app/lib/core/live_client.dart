import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:rxdart/rxdart.dart';
import 'package:wakelock_plus/wakelock_plus.dart';

import '../models/signaling_info.dart';
import '../pages/call/call_page.dart';

class LiveClient {
  LiveClient._();
  static final LiveClient instance = LiveClient._();

  OverlayEntry? _holder;
  bool isBusy = false;
  String? currentRoomID;
  Future Function(int duration, bool isPositive)? onTapHangup;

  void close() {
    _holder?.remove();
    _holder = null;
    isBusy = false;
    currentRoomID = null;
    WakelockPlus.disable();
  }

  void closeByRoomID(String roomID) {
    if (currentRoomID == roomID) close();
  }

  void start(
    BuildContext ctx, {
    required PublishSubject<CallEvent> callEventSubject,
    required CallType callType,
    required CallState initState,
    required String inviterUserID,
    required List<String> inviteeUserIDList,
    String? roomID,
    String? groupID,
    Future<SignalingCertificate> Function()? onDial,
    Future<SignalingCertificate> Function()? onTapPickup,
    Future Function()? onTapCancel,
    Future Function(int duration, bool isPositive)? onTapHangup,
    Future Function()? onTapReject,
    Future<Map<String, dynamic>?> Function(String userID)? onSyncUserInfo,
    Function()? onWaitingAccept,
    Function()? onStartCalling,
    Function(dynamic error, dynamic stack)? onError,
    Function()? onClose,
    Function()? onRoomDisconnected,
  }) {
    if (isBusy) {
      debugPrint('[LiveClient] start() aborted: already busy');
      return;
    }
    isBusy = true;
    currentRoomID = roomID;
    this.onTapHangup = onTapHangup;
    debugPrint('[LiveClient] start() initState=$initState, callType=$callType, roomID=$roomID');

    try {
      FocusScope.of(ctx).requestFocus(FocusNode());
    } catch (_) {}

    _holder = OverlayEntry(
      builder: (context) => CallPage(
        callType: callType,
        initState: initState,
        callEventSubject: callEventSubject,
        roomID: roomID,
        userID: initState == CallState.call
            ? inviteeUserIDList.first
            : inviterUserID,
        onDial: onDial,
        onTapCancel: onTapCancel,
        onTapHangup: onTapHangup,
        onTapReject: onTapReject,
        onTapPickup: onTapPickup,
        onSyncUserInfo: onSyncUserInfo,
        onWaitingAccept: onWaitingAccept,
        onStartCalling: onStartCalling,
        onError: onError,
        onRoomDisconnected: onRoomDisconnected,
        onClose: () {
          onClose?.call();
          close();
        },
      ),
    );

    try {
      Overlay.of(ctx).insert(_holder!);
      debugPrint('[LiveClient] overlay inserted successfully');
    } catch (e) {
      debugPrint('[LiveClient] Overlay.of(ctx) failed: $e');
      isBusy = false;
      _holder = null;
      return;
    }
    WakelockPlus.enable();
  }
}
