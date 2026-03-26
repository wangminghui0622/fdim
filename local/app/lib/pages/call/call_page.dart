import 'dart:async';

import 'package:flutter/material.dart';
import 'package:livekit_client/livekit_client.dart';
import 'package:permission_handler/permission_handler.dart';
import 'package:rxdart/rxdart.dart';

import '../../models/signaling_info.dart';

class CallPage extends StatefulWidget {
  final CallType callType;
  final CallState initState;
  final PublishSubject<CallEvent> callEventSubject;
  final String? roomID;
  final String? userID;
  final Future<SignalingCertificate> Function()? onDial;
  final Future<SignalingCertificate> Function()? onTapPickup;
  final Future Function()? onTapCancel;
  final Future Function(int duration, bool isPositive)? onTapHangup;
  final Future Function()? onTapReject;
  final Future<Map<String, dynamic>?> Function(String userID)? onSyncUserInfo;
  final Function()? onWaitingAccept;
  final Function()? onStartCalling;
  final Function(dynamic error, dynamic stack)? onError;
  final Function()? onClose;
  final Function()? onRoomDisconnected;

  const CallPage({
    super.key,
    required this.callType,
    required this.initState,
    required this.callEventSubject,
    this.roomID,
    this.userID,
    this.onDial,
    this.onTapPickup,
    this.onTapCancel,
    this.onTapHangup,
    this.onTapReject,
    this.onSyncUserInfo,
    this.onWaitingAccept,
    this.onStartCalling,
    this.onError,
    this.onClose,
    this.onRoomDisconnected,
  });

  @override
  State<CallPage> createState() => _CallPageState();
}

class _CallPageState extends State<CallPage> {
  late CallState _callState;
  late CallType _callType;
  String _nickname = '';
  String _faceURL = '';
  bool _micEnabled = true;
  bool _speakerOn = true;
  bool _cameraEnabled = true;

  Room? _room;
  EventsListener<RoomEvent>? _listener;
  LocalParticipant? _localParticipant;
  RemoteParticipant? _remoteParticipant;

  Timer? _durationTimer;
  int _callDuration = 0;
  StreamSubscription? _eventSub;
  SignalingCertificate? _certificate;

  @override
  void initState() {
    super.initState();
    _callState = widget.initState;
    _callType = widget.callType;
    _loadUserInfo();
    _listenEvents();

    if (_callState == CallState.call) {
      _dial();
    }
  }

  @override
  void dispose() {
    _durationTimer?.cancel();
    _eventSub?.cancel();
    (() async {
      _room?.removeListener(_onRoomUpdate);
      await _listener?.dispose();
      await _room?.disconnect();
      await _room?.dispose();
    })();
    super.dispose();
  }

  void _loadUserInfo() async {
    if (widget.userID == null) return;
    final info = await widget.onSyncUserInfo?.call(widget.userID!);
    if (info != null && mounted) {
      setState(() {
        _nickname = info['nickname'] ?? '';
        _faceURL = info['faceURL'] ?? '';
      });
    }
  }

  void _listenEvents() {
    _eventSub = widget.callEventSubject.stream.listen((event) {
      if (!mounted) return;
      switch (event.state) {
        case CallState.beAccepted:
          _onAccepted();
          break;
        case CallState.beRejected:
          _onEnded('对方已拒绝');
          break;
        case CallState.beCanceled:
          _onEnded('对方已取消');
          break;
        case CallState.beHangup:
          _onEnded('通话结束');
          break;
        case CallState.timeout:
          _onEnded('无人接听');
          break;
        default:
          break;
      }
    });
  }

  /// 发起通话
  void _dial() async {
    debugPrint('[CallPage] _dial() started');
    try {
      _certificate = await widget.onDial?.call();
      debugPrint('[CallPage] _dial() certificate: liveURL=${_certificate?.liveURL}, token=${_certificate?.token != null ? "present" : "null"}');
      if (_certificate == null) {
        debugPrint('[CallPage] _dial() failed: certificate is null');
        _onEnded('获取通话凭证失败');
        return;
      }
      widget.onWaitingAccept?.call();
      if (mounted) setState(() => _callState = CallState.call);
      debugPrint('[CallPage] _dial() success, waiting for accept');
    } catch (e, s) {
      debugPrint('[CallPage] _dial() error: $e');
      widget.onError?.call(e, s);
    }
  }

  /// 对方接听后，双方连接 LiveKit
  void _onAccepted() async {
    widget.onStartCalling?.call();
    setState(() => _callState = CallState.calling);
    await _connectRoom();
  }

  /// 被叫方接听
  void _pickup() async {
    try {
      _certificate = await widget.onTapPickup?.call();
      if (_certificate == null) {
        _onEnded('获取通话凭证失败');
        return;
      }
      setState(() => _callState = CallState.calling);
      await _connectRoom();
    } catch (e, s) {
      widget.onError?.call(e, s);
    }
  }

  Future<void> _connectRoom() async {
    if (_certificate == null) return;
    try {
      // 请求摄像头和麦克风权限
      final permissions = await [
        Permission.camera,
        Permission.microphone,
      ].request();
      if (permissions[Permission.microphone]?.isDenied == true) {
        debugPrint('[RTC] microphone permission denied');
      }
      if (permissions[Permission.camera]?.isDenied == true) {
        debugPrint('[RTC] camera permission denied');
      }

      _room = Room();
      _listener = _room!.createListener();
      // ignore: deprecated_member_use
      await _room!.connect(
        _certificate!.liveURL!,
        _certificate!.token!,
        roomOptions: RoomOptions(
          dynacast: true,
          adaptiveStream: false,
          defaultCameraCaptureOptions: const CameraCaptureOptions(
            params: VideoParametersPresets.h1080_169,
          ),
          defaultVideoPublishOptions: const VideoPublishOptions(
            videoEncoding: VideoEncoding(
              maxBitrate: 3000 * 1000, // 3 Mbps
              maxFramerate: 30,
            ),
          ),
        ),
      );
      if (!mounted) return;
      _room!.addListener(_onRoomUpdate);
      _setupListeners();
      _localParticipant = _room!.localParticipant;

      // 默认开启扬声器
      await Hardware.instance.setSpeakerphoneOn(true);

      // 发布音视频
      await _localParticipant?.setMicrophoneEnabled(true);
      if (_callType == CallType.video) {
        await _localParticipant?.setCameraEnabled(true);
      }

      _startDurationTimer();
      setState(() {});
    } catch (e, s) {
      widget.onError?.call(e, s);
    }
  }

  void _setupListeners() {
    _listener!
      ..on<RoomDisconnectedEvent>((event) {
        widget.onRoomDisconnected?.call();
        widget.onClose?.call();
      })
      ..on<ParticipantConnectedEvent>((_) {
        _remoteParticipant = _room?.remoteParticipants.values.firstOrNull;
        if (mounted) setState(() {});
      })
      ..on<ParticipantDisconnectedEvent>((_) {
        _remoteParticipant = null;
        if (mounted) setState(() {});
        // 对方离开，结束通话
        _onEnded('通话结束');
      })
      ..on<TrackSubscribedEvent>((_) {
        if (mounted) setState(() {});
      })
      ..on<TrackUnsubscribedEvent>((_) {
        if (mounted) setState(() {});
      });
  }

  void _onRoomUpdate() {
    _remoteParticipant = _room?.remoteParticipants.values.firstOrNull;
    if (mounted) setState(() {});
  }

  void _startDurationTimer() {
    _durationTimer = Timer.periodic(const Duration(seconds: 1), (_) {
      if (mounted) setState(() => _callDuration++);
    });
  }

  void _onEnded(String reason) {
    _durationTimer?.cancel();
    if (!mounted) return;
    setState(() => _callState = CallState.hangup);
    Future.delayed(const Duration(milliseconds: 500), () {
      widget.onClose?.call();
    });
  }

  void _hangup() async {
    await widget.onTapHangup?.call(_callDuration, true);
    _onEnded('通话结束');
  }

  void _cancel() async {
    await widget.onTapCancel?.call();
    _onEnded('已取消');
  }

  void _reject() async {
    await widget.onTapReject?.call();
    _onEnded('已拒绝');
  }

  void _toggleMic() async {
    _micEnabled = !_micEnabled;
    await _localParticipant?.setMicrophoneEnabled(_micEnabled);
    if (mounted) setState(() {});
  }

  void _toggleSpeaker() async {
    _speakerOn = !_speakerOn;
    await Hardware.instance.setSpeakerphoneOn(_speakerOn);
    if (mounted) setState(() {});
  }

  void _toggleCamera() async {
    _cameraEnabled = !_cameraEnabled;
    await _localParticipant?.setCameraEnabled(_cameraEnabled);
    if (mounted) setState(() {});
  }

  void _switchCamera() async {
    final publications = _localParticipant?.videoTrackPublications ?? [];
    for (final pub in publications) {
      if (pub.track is LocalVideoTrack) {
        try {
          final track = pub.track as LocalVideoTrack;
          final current = (track.currentOptions as CameraCaptureOptions).cameraPosition;
          final next = current == CameraPosition.front
              ? CameraPosition.back
              : CameraPosition.front;
          await track.setCameraPosition(next);
        } catch (_) {}
      }
    }
  }

  String _formatDuration(int seconds) {
    final m = (seconds ~/ 60).toString().padLeft(2, '0');
    final s = (seconds % 60).toString().padLeft(2, '0');
    return '$m:$s';
  }

  @override
  Widget build(BuildContext context) {
    return Material(
      color: const Color(0xFF0078D4),
      child: Container(
        decoration: const BoxDecoration(
          gradient: LinearGradient(
            begin: Alignment.topCenter,
            end: Alignment.bottomCenter,
            colors: [
              Color(0xFF0089FF),
              Color(0xFF0068C9),
            ],
          ),
        ),
        child: SafeArea(
          child: Stack(
            children: [
              // 远端视频（全屏）
              if (_callType == CallType.video &&
                  _callState == CallState.calling &&
                  _remoteParticipant != null)
                Positioned.fill(
                  child: _buildRemoteVideo(),
                ),

              // 本地视频（小窗，右上角）
              if (_callType == CallType.video &&
                  _callState == CallState.calling &&
                  _localParticipant != null)
                Positioned(
                  top: 60,
                  right: 16,
                  width: 120,
                  height: 160,
                  child: ClipRRect(
                    borderRadius: BorderRadius.circular(12),
                    child: _buildLocalVideo(),
                  ),
                ),

              // 顶部信息
              Positioned(
                top: 40,
                left: 0,
                right: 0,
                child: _buildTopInfo(),
              ),

              // 底部操作按钮
              Positioned(
                bottom: 40,
                left: 0,
                right: 0,
                child: _buildActions(),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildRemoteVideo() {
    final tracks = _remoteParticipant?.videoTrackPublications ?? [];
    for (final pub in tracks) {
      if (pub.track != null && !pub.isScreenShare) {
        return VideoTrackRenderer(pub.track as VideoTrack);
      }
    }
    return Container(color: const Color(0xFF005AA0));
  }

  Widget _buildLocalVideo() {
    final tracks = _localParticipant?.videoTrackPublications ?? [];
    for (final pub in tracks) {
      if (pub.track != null && !pub.isScreenShare) {
        return VideoTrackRenderer(pub.track as VideoTrack);
      }
    }
    return Container(color: const Color(0xFF005AA0));
  }

  Widget _buildTopInfo() {
    String statusText;
    switch (_callState) {
      case CallState.call:
        statusText = '等待对方接听...';
        break;
      case CallState.beCalled:
        statusText = _callType == CallType.video ? '视频通话邀请' : '语音通话邀请';
        break;
      case CallState.calling:
        statusText = _formatDuration(_callDuration);
        break;
      default:
        statusText = '';
    }

    return Column(
      children: [
        if (_faceURL.isNotEmpty)
          CircleAvatar(
            radius: 40,
            backgroundImage: NetworkImage(_faceURL),
          )
        else
          CircleAvatar(
            radius: 40,
            backgroundColor: Colors.white24,
            child: Text(
              _nickname.isNotEmpty ? _nickname[0] : '?',
              style: const TextStyle(fontSize: 32, color: Colors.white),
            ),
          ),
        const SizedBox(height: 12),
        Text(
          _nickname.isNotEmpty ? _nickname : (widget.userID ?? ''),
          style: const TextStyle(
            fontSize: 20,
            color: Colors.white,
            fontWeight: FontWeight.w600,
          ),
        ),
        const SizedBox(height: 6),
        Text(
          statusText,
          style: const TextStyle(fontSize: 14, color: Colors.white70),
        ),
      ],
    );
  }

  Widget _buildActions() {
    switch (_callState) {
      case CallState.call:
        // 发起方：显示取消按钮
        return _buildSingleAction(
          icon: Icons.call_end,
          color: Colors.red,
          label: '取消',
          onTap: _cancel,
        );

      case CallState.beCalled:
        // 被叫方：显示接听和拒绝
        return Row(
          mainAxisAlignment: MainAxisAlignment.spaceEvenly,
          children: [
            _buildActionButton(
              icon: Icons.call_end,
              color: Colors.red,
              label: '拒绝',
              onTap: _reject,
            ),
            _buildActionButton(
              icon: Icons.call,
              color: Colors.green,
              label: '接听',
              onTap: _pickup,
            ),
          ],
        );

      case CallState.calling:
        // 通话中：麦克风、挂断、扬声器 / 翻转摄像头
        return Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceEvenly,
              children: [
                _buildActionButton(
                  icon: _micEnabled ? Icons.mic : Icons.mic_off,
                  color: Colors.white24,
                  label: _micEnabled ? '静音' : '取消静音',
                  onTap: _toggleMic,
                ),
                _buildActionButton(
                  icon: Icons.call_end,
                  color: Colors.red,
                  label: '挂断',
                  onTap: _hangup,
                  size: 64,
                ),
                if (_callType == CallType.video)
                  _buildActionButton(
                    icon: Icons.cameraswitch,
                    color: Colors.white24,
                    label: '翻转',
                    onTap: _switchCamera,
                  )
                else
                  _buildActionButton(
                    icon:
                        _speakerOn ? Icons.volume_up : Icons.volume_off,
                    color: Colors.white24,
                    label: _speakerOn ? '扬声器' : '听筒',
                    onTap: _toggleSpeaker,
                  ),
              ],
            ),
            if (_callType == CallType.video) ...[
              const SizedBox(height: 16),
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                children: [
                  _buildActionButton(
                    icon: _cameraEnabled
                        ? Icons.videocam
                        : Icons.videocam_off,
                    color: Colors.white24,
                    label: _cameraEnabled ? '关摄像头' : '开摄像头',
                    onTap: _toggleCamera,
                  ),
                  _buildActionButton(
                    icon:
                        _speakerOn ? Icons.volume_up : Icons.volume_off,
                    color: Colors.white24,
                    label: _speakerOn ? '扬声器' : '听筒',
                    onTap: _toggleSpeaker,
                  ),
                ],
              ),
            ],
          ],
        );

      default:
        return const SizedBox.shrink();
    }
  }

  Widget _buildSingleAction({
    required IconData icon,
    required Color color,
    required String label,
    required VoidCallback onTap,
  }) {
    return Center(
      child: _buildActionButton(
          icon: icon, color: color, label: label, onTap: onTap, size: 64),
    );
  }

  Widget _buildActionButton({
    required IconData icon,
    required Color color,
    required String label,
    required VoidCallback onTap,
    double size = 52,
  }) {
    return GestureDetector(
      onTap: onTap,
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Container(
            width: size,
            height: size,
            decoration: BoxDecoration(
              color: color,
              shape: BoxShape.circle,
            ),
            child: Icon(icon, color: Colors.white, size: size * 0.5),
          ),
          const SizedBox(height: 6),
          Text(
            label,
            style: const TextStyle(fontSize: 12, color: Colors.white70),
          ),
        ],
      ),
    );
  }
}
