import 'dart:async';
import 'dart:convert';
import 'dart:io';

import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:flutter_easyloading/flutter_easyloading.dart';
import 'package:get/get.dart';
import 'package:just_audio/just_audio.dart';
import 'package:rxdart/rxdart.dart';
import 'package:uuid/uuid.dart';

import '../models/signaling_info.dart';
import '../sdk/fdim_sdk.dart';
import 'config.dart';
import 'live_client.dart';
import 'rtc_api.dart';

mixin LiveController {
  final signalingSubject = PublishSubject<CallEvent>();

  final _audioPlayer = AudioPlayer(handleInterruptions: false);
  Timer? _callTimeoutTimer;

  bool get isRtcBusy => LiveClient.instance.isBusy;

  void onInitLive() {
    _signalingListener();
  }

  void onCloseLive() {
    _callTimeoutTimer?.cancel();
    signalingSubject.close();
    _stopSound();
  }

  void _signalingListener() {
    signalingSubject.stream.listen((event) async {
      if (event.state == CallState.beCalled) {
        debugPrint('[RTC] beCalled event received from ${event.data.invitation?.inviterUserID}');

        // 正在通话中 → 自动回复"忙线"
        if (LiveClient.instance.isBusy) {
          debugPrint('[RTC] already busy, auto-reply busy to ${event.data.invitation?.inviterUserID}');
          _autoReplyBusy(event.data);
          return;
        }

        final ctx = Get.overlayContext ?? Get.context;
        if (ctx == null) {
          debugPrint('[RTC] beCalled aborted: no overlay/context');
          return;
        }
        _playSound();
        final mediaType = event.data.invitation!.mediaType;
        final callType =
            mediaType == 'audio' ? CallType.audio : CallType.video;
        // 被叫方也启动超时保护（防止取消信令丢失）
        final timeout = event.data.invitation?.timeout ?? 30;
        _callTimeoutTimer?.cancel();
        _callTimeoutTimer = Timer(Duration(seconds: timeout + 5), () {
          if (LiveClient.instance.isBusy) {
            // CallPage._listenEvents 收到 timeout 后会调用 _onEnded → onClose → close()
            signalingSubject.add(CallEvent(CallState.timeout, event.data));
          }
        });

        LiveClient.instance.start(
          ctx,
          callEventSubject: signalingSubject,
          roomID: event.data.invitation!.roomID!,
          inviteeUserIDList: event.data.invitation!.inviteeUserIDList!,
          inviterUserID: event.data.invitation!.inviterUserID!,
          callType: callType,
          initState: CallState.beCalled,
          onSyncUserInfo: _onSyncUserInfo,
          onTapPickup: () => _onTapPickup(
            event.data..userID = Config.userID,
          ),
          onTapReject: () => _onTapReject(
            event.data..userID = Config.userID,
          ),
          onTapHangup: (duration, isPositive) => _onTapHangup(
            event.data..userID = Config.userID,
            duration,
            isPositive,
          ),
          onError: _onError,
          onRoomDisconnected: () {},
        );
      } else if (event.state == CallState.beBusy) {
        _callTimeoutTimer?.cancel();
        _stopSound();
        // 发持久化"对方正忙"消息到聊天记录
        _sendCallResult(
          signaling: event.data,
          result: CallResultType.busy,
        );
      } else if (event.state == CallState.beRejected ||
          event.state == CallState.beCanceled ||
          event.state == CallState.beAccepted ||
          event.state == CallState.timeout ||
          event.state == CallState.beHangup) {
        _callTimeoutTimer?.cancel();
        _stopSound();
      }
    });
  }

  /// 正在通话中时自动回复忙线
  Future<void> _autoReplyBusy(SignalingInfo signaling) async {
    try {
      final data = {
        'customType': CustomMessageType.callingBusy,
        'data': signaling.invitation!.toJson(),
      };
      final message = await OpenIM.iMManager.messageManager
          .createCustomMessage(
              data: jsonEncode(data), extension: '', description: '');
      await OpenIM.iMManager.messageManager.sendMessage(
        message: message,
        offlinePushInfo: OfflinePushInfo(),
        userID: signaling.invitation!.inviterUserID,
        isOnlineOnly: true,
      );
      debugPrint('[RTC] busy reply sent to ${signaling.invitation!.inviterUserID}');
    } catch (e) {
      debugPrint('[RTC] _autoReplyBusy error: $e');
    }
  }

  /// 发起通话
  void call({
    required CallType callType,
    required List<String> inviteeUserIDList,
    String? groupID,
  }) {
    debugPrint('[RTC] call() invoked: callType=$callType, invitees=$inviteeUserIDList');
    if (LiveClient.instance.isBusy) {
      debugPrint('[RTC] call() aborted: already busy');
      return;
    }
    final ctx = Get.overlayContext ?? Get.context;
    if (ctx == null) {
      debugPrint('[RTC] call() aborted: no overlay/context available');
      return;
    }

    final mediaType = callType == CallType.audio ? 'audio' : 'video';
    final inviterUserID = Config.userID;

    final signal = SignalingInfo(
      userID: inviterUserID,
      invitation: InvitationInfo(
        inviterUserID: inviterUserID,
        inviteeUserIDList: inviteeUserIDList,
        roomID: const Uuid().v4(),
        timeout: 30,
        mediaType: mediaType,
        sessionType: groupID != null ? 3 : 1,
        platformID: Platform.isAndroid ? 2 : 1,
        groupID: groupID,
      ),
    );

    LiveClient.instance.start(
      ctx,
      callEventSubject: signalingSubject,
      roomID: signal.invitation!.roomID,
      inviterUserID: inviterUserID,
      inviteeUserIDList: inviteeUserIDList,
      groupID: groupID,
      callType: callType,
      initState: CallState.call,
      onDial: () => _onDialSingle(signal),
      onTapCancel: () => _onTapCancel(signal),
      onTapHangup: (duration, isPositive) =>
          _onTapHangup(signal, duration, isPositive),
      onSyncUserInfo: _onSyncUserInfo,
      onWaitingAccept: () {
        _playSound();
        // 与官方一致：超时自动取消
        final timeout = signal.invitation?.timeout ?? 30;
        _callTimeoutTimer?.cancel();
        _callTimeoutTimer = Timer(Duration(seconds: timeout), () {
          if (LiveClient.instance.isBusy) {
            _onTimeout(signal);
            signalingSubject.add(CallEvent(CallState.timeout, signal));
          }
        });
      },
      onStartCalling: () {
        _callTimeoutTimer?.cancel();
        _stopSound();
      },
      onError: _onError,
      onRoomDisconnected: () {},
      onClose: () {
        _callTimeoutTimer?.cancel();
        _stopSound();
      },
    );
  }

  /// 拨出 → 发信令 + 获取 Token
  Future<SignalingCertificate> _onDialSingle(SignalingInfo signaling) async {
    debugPrint('[RTC] _onDialSingle: roomID=${signaling.invitation!.roomID}, invitee=${signaling.invitation!.inviteeUserIDList!.first}');
    final data = {
      'customType': CustomMessageType.callingInvite,
      'data': signaling.invitation!.toJson(),
    };
    final message = await OpenIM.iMManager.messageManager
        .createCustomMessage(
            data: jsonEncode(data), extension: '', description: '');
    debugPrint('[RTC] _onDialSingle: sending signaling message...');
    await OpenIM.iMManager.messageManager.sendMessage(
      message: message,
      offlinePushInfo: OfflinePushInfo(),
      userID: signaling.invitation!.inviteeUserIDList!.first,
      isOnlineOnly: true,
    );
    debugPrint('[RTC] _onDialSingle: signaling sent OK, getting RTC token...');
    final certificate = await RtcApi.getTokenForRTC(
      signaling.invitation!.roomID!,
      Config.userID,
    );
    debugPrint('[RTC] _onDialSingle: got certificate liveURL=${certificate.liveURL}');
    return certificate;
  }

  /// 接听
  Future<SignalingCertificate> _onTapPickup(SignalingInfo signaling) async {
    _stopSound();
    final data = {
      'customType': CustomMessageType.callingAccept,
      'data': signaling.invitation!.toJson(),
    };
    final message = await OpenIM.iMManager.messageManager
        .createCustomMessage(
            data: jsonEncode(data), extension: '', description: '');
    await OpenIM.iMManager.messageManager.sendMessage(
      message: message,
      offlinePushInfo: OfflinePushInfo(),
      userID: signaling.invitation!.inviterUserID,
      isOnlineOnly: true,
    );
    final certificate = await RtcApi.getTokenForRTC(
      signaling.invitation!.roomID!,
      Config.userID,
    );
    return certificate;
  }

  /// 拒绝
  Future<void> _onTapReject(SignalingInfo signaling) async {
    _stopSound();
    final data = {
      'customType': CustomMessageType.callingReject,
      'data': signaling.invitation!.toJson(),
    };
    final message = await OpenIM.iMManager.messageManager
        .createCustomMessage(
            data: jsonEncode(data), extension: '', description: '');
    final recvUserID =
        signaling.invitation!.inviterUserID == Config.userID
            ? signaling.invitation!.inviteeUserIDList!.first
            : signaling.invitation!.inviterUserID;
    await OpenIM.iMManager.messageManager.sendMessage(
      message: message,
      offlinePushInfo: OfflinePushInfo(),
      userID: recvUserID,
      isOnlineOnly: true,
    );
    // 拒绝方发持久化消息（对方看到"对方已拒绝"，自己看到"已拒绝"）
    _sendCallResult(
      signaling: signaling,
      result: CallResultType.rejected,
    );
  }

  /// 取消
  Future<void> _onTapCancel(SignalingInfo signaling) async {
    _stopSound();
    final data = {
      'customType': CustomMessageType.callingCancel,
      'data': signaling.invitation!.toJson(),
    };
    final message = await OpenIM.iMManager.messageManager
        .createCustomMessage(
            data: jsonEncode(data), extension: '', description: '');
    final recvUserID =
        signaling.invitation!.inviterUserID == Config.userID
            ? signaling.invitation!.inviteeUserIDList!.first
            : signaling.invitation!.inviterUserID;
    await OpenIM.iMManager.messageManager.sendMessage(
      message: message,
      offlinePushInfo: OfflinePushInfo(),
      userID: recvUserID,
      isOnlineOnly: true,
    );
    // 取消方发持久化消息（对方看到"对方已取消"，自己看到"已取消"）
    _sendCallResult(
      signaling: signaling,
      result: CallResultType.cancelled,
    );
  }

  /// 超时未接听 → 发取消信令 + 记录超时结果
  Future<void> _onTimeout(SignalingInfo signaling) async {
    _stopSound();
    final data = {
      'customType': CustomMessageType.callingCancel,
      'data': signaling.invitation!.toJson(),
    };
    final message = await OpenIM.iMManager.messageManager
        .createCustomMessage(
            data: jsonEncode(data), extension: '', description: '');
    final recvUserID =
        signaling.invitation!.inviterUserID == Config.userID
            ? signaling.invitation!.inviteeUserIDList!.first
            : signaling.invitation!.inviterUserID;
    try {
      await OpenIM.iMManager.messageManager.sendMessage(
        message: message,
        offlinePushInfo: OfflinePushInfo(),
        userID: recvUserID,
        isOnlineOnly: true,
      );
    } catch (_) {}
    _sendCallResult(
      signaling: signaling,
      result: CallResultType.timeout,
    );
  }

  /// 挂断
  /// isPositive=true: 用户主动挂断，发信令+completed结果
  /// isPositive=false: 网络中断，不发信令，发networkError结果
  Future<void> _onTapHangup(
    SignalingInfo signaling,
    int duration,
    bool isPositive,
  ) async {
    if (isPositive) {
      final data = {
        'customType': CustomMessageType.callingHungup,
        'data': signaling.invitation!.toJson(),
      };
      final message = await OpenIM.iMManager.messageManager
          .createCustomMessage(
              data: jsonEncode(data), extension: '', description: '');
      final recvUserID =
          signaling.invitation!.inviterUserID == Config.userID
              ? signaling.invitation!.inviteeUserIDList!.first
              : signaling.invitation!.inviterUserID;
      await OpenIM.iMManager.messageManager.sendMessage(
        message: message,
        offlinePushInfo: OfflinePushInfo(),
        userID: recvUserID,
        isOnlineOnly: true,
      );
    }
    // 挂断 → 发持久化通话结果消息
    final resultType = isPositive
        ? CallResultType.completed
        : CallResultType.networkError;
    _sendCallResult(
      signaling: signaling,
      result: resultType,
      duration: duration,
    );
    _stopSound();
  }

  Future<Map<String, dynamic>?> _onSyncUserInfo(String userID) async {
    try {
      final list =
          await OpenIM.iMManager.userManager.getUsersInfo(userIDList: [userID]);
      if (list.isNotEmpty) {
        return {
          'nickname': list.first.nickname,
          'faceURL': list.first.faceURL,
        };
      }
    } catch (e) {
      debugPrint('[RTC] syncUserInfo error: $e');
    }
    return null;
  }

  void _onError(dynamic error, dynamic stack) {
    debugPrint('[RTC] error: $error\n$stack');
    _callTimeoutTimer?.cancel();
    LiveClient.instance.close();
    _stopSound();
    try {
      EasyLoading.showToast('通话失败: $error');
    } catch (_) {}
  }

  /// 发送持久化通话结果消息到聊天记录
  Future<void> _sendCallResult({
    required SignalingInfo signaling,
    required String result,
    int duration = 0,
  }) async {
    try {
      final callResult = CallResultInfo(
        callType: signaling.invitation?.mediaType ?? 'audio',
        result: result,
        duration: duration,
      );
      final data = {
        'customType': CustomMessageType.callResult,
        'data': callResult.toJson(),
      };
      final message = await OpenIM.iMManager.messageManager
          .createCustomMessage(
              data: jsonEncode(data), extension: '', description: '');
      final recvUserID =
          signaling.invitation!.inviterUserID == Config.userID
              ? signaling.invitation!.inviteeUserIDList!.first
              : signaling.invitation!.inviterUserID;
      await OpenIM.iMManager.messageManager.sendMessage(
        message: message,
        offlinePushInfo: OfflinePushInfo(),
        userID: recvUserID,
        isOnlineOnly: false,
      );
      // 立即通知聊天页显示此消息（不等 push 回推）
      try {
        (this as dynamic).newMessageSubject.add(message);
      } catch (_) {}
    } catch (e) {
      debugPrint('[RTC] _sendCallResult error: $e');
    }
  }

  void _playSound() async {
    if (!_audioPlayer.playerState.playing) {
      // Use a simple ringtone tone; in production, load a bundled asset
      try {
        await _audioPlayer.setUrl(
            'https://actions.google.com/sounds/v1/alarms/beep_short.ogg');
        _audioPlayer.setLoopMode(LoopMode.one);
        _audioPlayer.setVolume(1.0);
        _audioPlayer.play();
      } catch (_) {}
    }
  }

  void _stopSound() async {
    if (_audioPlayer.playerState.playing) {
      _audioPlayer.stop();
    }
  }

  /// 处理收到的自定义消息（信令）— 由 IM Controller 调用
  void handleSignalingMessage(Message message) {
    try {
      if (message.contentType != MessageType.custom) return;
      final customData = message.customElem?.data;
      if (customData == null) return;
      final map = jsonDecode(customData);
      final customType = map['customType'] as int?;
      final data = map['data'];
      if (customType == null || data == null) return;

      final invitation = InvitationInfo.fromJson(data);
      final signalingInfo = SignalingInfo(
        userID: message.sendID,
        invitation: invitation,
      );

      switch (customType) {
        case CustomMessageType.callingInvite:
          signalingSubject.add(CallEvent(CallState.beCalled, signalingInfo));
          break;
        case CustomMessageType.callingAccept:
          signalingSubject.add(CallEvent(CallState.beAccepted, signalingInfo));
          break;
        case CustomMessageType.callingReject:
          signalingSubject.add(CallEvent(CallState.beRejected, signalingInfo));
          break;
        case CustomMessageType.callingCancel:
          signalingSubject.add(CallEvent(CallState.beCanceled, signalingInfo));
          break;
        case CustomMessageType.callingHungup:
          signalingSubject.add(CallEvent(CallState.beHangup, signalingInfo));
          break;
        case CustomMessageType.callingBusy:
          signalingSubject.add(CallEvent(CallState.beBusy, signalingInfo));
          break;
      }
    } catch (e) {
      debugPrint('[RTC] handleSignalingMessage error: $e');
    }
  }
}
