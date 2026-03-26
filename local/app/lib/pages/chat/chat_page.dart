import 'dart:async';
import 'dart:convert';
import 'dart:io';

import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_easyloading/flutter_easyloading.dart';
import 'package:get/get.dart';
import 'package:image_picker/image_picker.dart';
import 'package:focus_detector_v2/focus_detector_v2.dart';
import 'package:path_provider/path_provider.dart';
import 'package:record/record.dart';

import '../../core/config.dart';
import '../../core/controllers/im_controller.dart';
import '../../core/upload_service.dart';
import '../../models/signaling_info.dart';
import '../../core/local_store.dart';
import '../../sdk/flutter_openim_sdk.dart';
import '../../routes/app_routes.dart';
import '../../widgets/media_preview_page.dart';
import '../../widgets/emoji_picker_widget.dart';
import '../../utils/audio_player_manager.dart';
import '../../utils/file_utils.dart';

class ChatPage extends StatefulWidget {
  const ChatPage({super.key});

  @override
  State<ChatPage> createState() => _ChatPageState();
}

class _ChatPageState extends State<ChatPage> {
  final _inputController = TextEditingController();
  final _scrollController = ScrollController();
  final _focusNode = FocusNode();
  final _messages = <Message>[];
  final _pendingNewMessages = <Message>[];
  bool _loading = true;
  bool _loadingMore = false;
  bool _isOlderEnd = false;   // 已加载到最早的消息
  bool _isNewerEnd = false;   // 已加载到最新的消息
  bool _isOnline = false;
  bool _showEmojiPicker = false;
  bool _showMorePanel = false;
  bool _sendButtonVisible = false;
  bool _showScrollDown = false;
  int _unreadBelow = 0;
  bool _isFirstLoad = true;
  bool _isLoadingMore = false;
  static const int _pageSize = 40;
  Timer? _onlineTimer;
  StreamSubscription? _newMsgSub;
  StreamSubscription? _readReceiptSub;
  StreamSubscription? _convChangedSub;
  StreamSubscription? _msgRevokedSub;
  final _audioPlayer = AudioPlayerManager();
  final _audioRecorder = AudioRecorder();
  final _playedVoiceIds = <String>{};
  bool _voiceMode = false;
  bool _isRecording = false;
  String? _recordFilePath;
  int? _recordStartAt;
  double _cachedKeyboardHeight = 290;
  bool _hideEmojiWhenKeyboardShows = false;
  bool _pendingScrollAfterKeyboardShows = false;

  late String conversationID;
  late String userID;
  late String groupID;
  late String showName;
  late String faceURL;
  late int sessionType;

  @override
  void initState() {
    super.initState();
    final args = Get.arguments as Map<String, dynamic>;
    conversationID = args['conversationID'] ?? '';
    userID = args['userID'] ?? '';
    groupID = args['groupID'] ?? '';
    showName = args['showName'] ?? '';
    faceURL = args['faceURL'] ?? '';
    sessionType = args['sessionType'] ?? 1;

    Get.find<IMController>().currentChatConversationID.value = conversationID;
    // 与官方 Go SDK 一致：告知 SDK 当前活跃会话，新消息不递增 unreadCount
    OpenIM.iMManager.setActiveConversation(conversationID);
    _playedVoiceIds.addAll(LocalStore.getPlayedVoiceIds(conversationID));
    _loadMessages();
    _clearUnreadCount();

    // Listen for new messages via SDK
    final imCtrl = Get.find<IMController>();
    _newMsgSub = imCtrl.newMessageSubject.listen(_onNewMessage);

    // Listen for read receipts (C2C only)
    _readReceiptSub = imCtrl.c2cReadReceiptSubject.listen(_onReadReceipt);
    _msgRevokedSub = imCtrl.msgRevokedSubject.listen(_onMessageRevoked);

    // 【修复问题2】监听会话变更通知，确保未读数实时更新
    _convChangedSub = imCtrl.conversationChangedSubject.listen((_) {
      // 会话变更时不需要特殊处理，会话列表页会自动刷新
    });

    if (sessionType == 1 && userID.isNotEmpty) {
      _checkOnlineStatus();
      _onlineTimer = Timer.periodic(
        const Duration(seconds: 15),
        (_) => _checkOnlineStatus(),
      );
    }

    _scrollController.addListener(() {
      // 上滑加载更早消息
      if (_scrollController.position.pixels <=
          _scrollController.position.minScrollExtent + 50) {
        _loadOlderMessages();
      }
      // 下滑加载更新消息
      final distToBottom = _scrollController.position.maxScrollExtent -
          _scrollController.position.pixels;
      if (distToBottom < 80) {
        _loadNewerMessages();
      }
      // 显示/隐藏向下箭头
      final shouldShow = distToBottom > 80;
      if (shouldShow != _showScrollDown) {
        setState(() => _showScrollDown = shouldShow);
      }
    });

    _focusNode.addListener(() {
      if (!mounted) return;
      setState(() {});
      if (_focusNode.hasFocus) {
        Future.delayed(const Duration(milliseconds: 80), _scrollToBottom);
        Future.delayed(const Duration(milliseconds: 220), _scrollToBottom);
      }
    });

    // Listen to input controller changes to toggle send button visibility
    _inputController.addListener(() {
      setState(() {
        _sendButtonVisible = _inputController.text.trim().isNotEmpty;
      });
    });
  }

  @override
  void dispose() {
    // 与官方 ChatLogic.onClose 一致：离开时再次标记已读，确保最后收到的消息也被标记
    _markAsRead();
    // 清除活跃会话标记
    OpenIM.iMManager.setActiveConversation(null);
    final imCtrl = Get.find<IMController>();
    if (imCtrl.currentChatConversationID.value == conversationID) {
      imCtrl.currentChatConversationID.value = '';
    }
    _inputController.dispose();
    _scrollController.dispose();
    _focusNode.dispose();
    _audioRecorder.dispose();
    _onlineTimer?.cancel();
    _newMsgSub?.cancel();
    _readReceiptSub?.cancel();
    _convChangedSub?.cancel();
    _msgRevokedSub?.cancel();
    _audioPlayer.dispose();
    super.dispose();
  }

  void _onNewMessage(Message msg) {
    final match =
        (msg.isSingleChat && (msg.sendID == userID || msg.recvID == userID)) ||
        (msg.isGroupChat && msg.groupID == groupID);
    if (!match || !mounted) return;

    final existingMessageIndex = _messages.indexWhere(
      (e) => e.clientMsgID != null && e.clientMsgID == msg.clientMsgID,
    );
    if (existingMessageIndex >= 0) {
      final existing = _messages[existingMessageIndex];
      final shouldMergeSelfEcho =
          msg.sendID == Config.userID &&
          ((msg.seq ?? 0) > 0 || (msg.serverMsgID ?? '').isNotEmpty);
      if (shouldMergeSelfEcho) {
        setState(() {
          existing.serverMsgID = msg.serverMsgID ?? existing.serverMsgID;
          existing.sendTime = msg.sendTime ?? existing.sendTime;
          existing.createTime = msg.createTime ?? existing.createTime;
          existing.seq = msg.seq ?? existing.seq;
          existing.status = MessageStatus.succeeded;
        });
        if ((existing.seq ?? 0) > 0) {
          LocalStore.setMaxSeq(conversationID, existing.seq!);
          LocalStore.putMessage(conversationID, existing);
          LocalStore.updateLatestMsg(conversationID, existing);
        }
      }
      return;
    }

    final existingPendingIndex = _pendingNewMessages.indexWhere(
      (e) => e.clientMsgID != null && e.clientMsgID == msg.clientMsgID,
    );
    if (existingPendingIndex >= 0) {
      final existing = _pendingNewMessages[existingPendingIndex];
      final shouldMergeSelfEcho =
          msg.sendID == Config.userID &&
          ((msg.seq ?? 0) > 0 || (msg.serverMsgID ?? '').isNotEmpty);
      if (shouldMergeSelfEcho) {
        setState(() {
          existing.serverMsgID = msg.serverMsgID ?? existing.serverMsgID;
          existing.sendTime = msg.sendTime ?? existing.sendTime;
          existing.createTime = msg.createTime ?? existing.createTime;
          existing.seq = msg.seq ?? existing.seq;
          existing.status = MessageStatus.succeeded;
        });
      }
      return;
    }

    final isIncoming = msg.sendID != Config.userID;
    final wasAtBottom = _isNearBottom();

    if (!isIncoming || wasAtBottom) {
      setState(() => _messages.add(msg));
      if ((msg.seq ?? 0) > 0) {
        LocalStore.setMaxSeq(conversationID, msg.seq!);
      }
      if (wasAtBottom) {
        _scrollToBottom();
      }
      if (isIncoming) {
        _markMessageAsRead(msg, true);
      }
      return;
    }

    setState(() {
      _pendingNewMessages.add(msg);
      _unreadBelow = _pendingNewMessages.where((e) => e.sendID != Config.userID).length;
      _showScrollDown = true;
    });
  }

  void _onMessageRevoked(RevokedInfo info) {
    if (!mounted) return;
    debugPrint(
      '[Revoke][Client] notification received: '
      'clientMsgID=${info.clientMsgID} revokerID=${info.revokerID} '
      'revokeTime=${info.revokeTime} sessionType=${info.sessionType}',
    );
    final clientMsgID = info.clientMsgID;
    if (clientMsgID == null || clientMsgID.isEmpty) return;

    final revokeTime = info.revokeTime ?? DateTime.now().millisecondsSinceEpoch;
    final isSelfRevoke = info.revokerID == Config.userID;
    final revokeText = isSelfRevoke ? '你撤回了一条消息' : '对方撤回了一条消息';

    final index = _messages.indexWhere((e) => e.clientMsgID == clientMsgID);
    if (index < 0) return;

    final old = _messages[index];
    // 已经是撤回状态则跳过（发送方本地已先更新，服务端推送到达时不重复处理）
    if (old.contentType == MessageType.revokeMessageNotification) return;
    final revokedMsg = Message(
      clientMsgID: old.clientMsgID,
      serverMsgID: old.serverMsgID,
      createTime: old.createTime,
      sendTime: revokeTime,
      sendID: info.revokerID ?? old.sendID,
      senderNickname: info.revokerNickname ?? old.senderNickname,
      senderFaceUrl: old.senderFaceUrl,
      senderPlatformID: old.senderPlatformID,
      sessionType: old.sessionType,
      contentType: MessageType.revokeMessageNotification,
      status: MessageStatus.succeeded,
      textElem: TextElem(content: revokeText),
    )
      ..seq = old.seq
      ..isRead = true
      ..recvID = old.recvID
      ..groupID = old.groupID;

    setState(() {
      _messages[index] = revokedMsg;
    });

    LocalStore.putMessage(conversationID, revokedMsg);
    LocalStore.updateLatestMsg(conversationID, revokedMsg);
  }

  void _onReadReceipt(List<ReadReceiptInfo> list) {
    if (!mounted) return;
    try {
      for (var readInfo in list) {
        if (readInfo.msgIDList != null && readInfo.msgIDList!.isNotEmpty) {
          for (var e in _messages) {
            if (readInfo.msgIDList!.contains(e.clientMsgID)) {
              // 与官方 chat_logic.dart onRecvC2CReadReceipt 一致
              e.isRead = true;
              e.hasReadTime = DateTime.now().millisecondsSinceEpoch;
              // 已读状态已由 IMManager._handlePushMsg 写入本地 DB
            }
          }
        }
      }
      if (mounted) setState(() {});
    } catch (e) {
      debugPrint('[Chat] _onReadReceipt error: $e');
    }
  }

 
  /// 过滤掉信令消息（customType 200-204），它们不应出现在聊天记录中
  List<Message> _filterSignalingMessages(List<Message> msgs) {
    return msgs.where((msg) {
      if (msg.contentType != MessageType.custom) return true;
      try {
        final raw = msg.customElem?.data;
        if (raw == null) return true;
        final customType = jsonDecode(raw)['customType'] as int?;
        if (customType == null) return true;
        // 信令消息 200-204：不显示
        return customType < CustomMessageType.callingInvite ||
               customType > CustomMessageType.callingHungup;
      } catch (_) {
        return true;
      }
    }).toList();
  }

  /// 按照官方逻辑加载历史消息
  Future<void> _loadMessages() async {
    try {
      debugPrint('[Chat] Loading messages for conversationID: $conversationID, isFirstLoad: $_isFirstLoad');

      final result = await OpenIM.iMManager.messageManager.getAdvancedHistoryMessageList(
        conversationID: conversationID,
        count: _isFirstLoad ? _pageSize : _messages.length,
        startMsg: _isFirstLoad ? null : _messages.firstOrNull,
      );

      if (result.messageList == null || result.messageList!.isEmpty) {
        if (mounted) {
          setState(() {
            _loading = false;
            _isOlderEnd = true;
            _isNewerEnd = true;
          });
        }
        return;
      }

      final msgs = _filterSignalingMessages(result.messageList!);

      if (_isFirstLoad) {
        _isFirstLoad = false;
        // 首次加载：清空并替换所有消息
        if (mounted) {
          setState(() {
            _messages.clear();
            _messages.addAll(msgs);
            _loading = false;
            _isOlderEnd = result.isEnd == true;
            _isNewerEnd = true;
          });
          // 官方：首次加载后直接滚到底部
          _scrollToBottom();
        }
      } else {
        // 加载更多：插入到顶部
        if (mounted) {
          setState(() {
            _messages.insertAll(0, msgs);
            _loading = false;
            _isOlderEnd = result.isEnd == true;
          });
        }
      }
    } catch (e) {
      debugPrint('[Chat] loadMessages error: $e');
      if (mounted) setState(() => _loading = false);
    }
  }

  /// 上滑加载更早消息（按照官方逻辑）
  Future<void> _loadOlderMessages() async {
    if (_isLoadingMore || _isOlderEnd) return;
    setState(() => _isLoadingMore = true);
    try {
      final result = await OpenIM.iMManager.messageManager.getAdvancedHistoryMessageList(
        conversationID: conversationID,
        count: _pageSize,
        startMsg: _messages.firstOrNull,
      );

      if (result.messageList != null && result.messageList!.isNotEmpty) {
        // 保持滚动位置
        final offset = _scrollController.offset;
        final oldHeight = _scrollController.position.maxScrollExtent;
        
        final older = _filterSignalingMessages(result.messageList!);
        setState(() {
          _messages.insertAll(0, older);
          _isOlderEnd = result.isEnd == true;
        });
        
        WidgetsBinding.instance.addPostFrameCallback((_) {
          final newHeight = _scrollController.position.maxScrollExtent;
          _scrollController.jumpTo(offset + (newHeight - oldHeight));
        });
      }
    } catch (e) {
      debugPrint('[Chat] _loadOlderMessages error: $e');
    } finally {
      setState(() => _isLoadingMore = false);
    }
  }

  /// 下滑加载更新消息
  /// 下滑加载更新消息（按照官方逻辑，实际上不需要，因为 WebSocket 会推送）
  Future<void> _loadNewerMessages() async {
    if (!_isNearBottom() || _pendingNewMessages.isEmpty) return;
    _appendPendingMessages(scrollToBottom: false);
  }

  /// 按照官方逻辑：进入时清零未读数
  void _clearUnreadCount() {
    final conv = LocalStore.getConversation(conversationID);
    if (conv != null && conv.unreadCount > 0) {
      OpenIM.iMManager.conversationManager.markConversationMessageAsRead(
        conversationID: conversationID,
      );
    }
  }

  /// 移除不再需要的 seq 边界更新逻辑
  void _updateSeqBounds(List<Message> msgs) {
    // 官方实现中没有这个逻辑，保留空方法以避免编译错误
  }

  Future<void> _checkOnlineStatus() async {
    if (userID.isEmpty) return;
    try {
      final result = await OpenIM.iMManager.userManager.getUserStatus([userID]);
      if (result.isNotEmpty && mounted) {
        final old = _isOnline;
        _isOnline = result.first.isOnline;
        if (old != _isOnline) setState(() {});
      }
    } catch (_) {}
  }

  bool _isNearBottom() {
    if (!_scrollController.hasClients) return true;
    return (_scrollController.position.maxScrollExtent -
            _scrollController.position.pixels) <=
        80;
  }

  void _appendPendingMessages({bool scrollToBottom = true}) {
    if (_pendingNewMessages.isEmpty || !mounted) return;

    final incoming = _pendingNewMessages.where((e) => e.sendID != Config.userID).toList();
    setState(() {
      _messages.addAll(_pendingNewMessages);
      _pendingNewMessages.clear();
      _unreadBelow = 0;
    });

    for (final msg in incoming) {
      _markMessageAsRead(msg, true);
    }

    if (scrollToBottom) {
      _scrollToBottom();
    }
  }

  void _scrollToBottom() {
    WidgetsBinding.instance.addPostFrameCallback((_) async {
      if (!mounted || !_scrollController.hasClients) return;

      double? lastTarget;
      for (var i = 0; i < 3; i++) {
        if (!mounted || !_scrollController.hasClients) return;
        final target = _scrollController.position.maxScrollExtent;
        if (lastTarget != null && (target - lastTarget).abs() < 1) {
          break;
        }
        _scrollController.jumpTo(target);
        lastTarget = target;
        await WidgetsBinding.instance.endOfFrame;
      }

      if (_pendingNewMessages.isNotEmpty) {
        _appendPendingMessages(scrollToBottom: false);
        WidgetsBinding.instance.addPostFrameCallback((_) async {
          if (!mounted || !_scrollController.hasClients) return;
          double? lastPendingTarget;
          for (var i = 0; i < 3; i++) {
            if (!mounted || !_scrollController.hasClients) return;
            final target = _scrollController.position.maxScrollExtent;
            if (lastPendingTarget != null &&
                (target - lastPendingTarget).abs() < 1) {
              break;
            }
            _scrollController.jumpTo(target);
            lastPendingTarget = target;
            await WidgetsBinding.instance.endOfFrame;
          }
        });
      }
    });
  }

  Future<void> _markAsRead() async {
    if (!mounted) return;
    try {
      await OpenIM.iMManager.conversationManager.markConversationMessageAsRead(
        conversationID: conversationID,
      );
      final now = DateTime.now().millisecondsSinceEpoch;
      for (final msg in _messages) {
        if (msg.sendID != Config.userID && msg.isRead != true) {
          msg.isRead = true;
          msg.hasReadTime = now;
          if (msg.clientMsgID != null) {
            await LocalStore.markMessageRead(
              conversationID,
              msg.clientMsgID!,
              now,
            );
          }
        }
      }
      if (mounted) setState(() {});
    } catch (e) {
      debugPrint('[Chat] _markAsRead error: $e');
    }
  }

  Future<void> _markMessageAsRead(Message msg, bool visible) async {
    if (!visible) return;
    if ((msg.contentType ?? 0) >= 1000 || msg.contentType == MessageType.voice) {
      return;
    }
    if (msg.sendID == Config.userID || msg.isRead == true) return;

    final now = DateTime.now().millisecondsSinceEpoch;
    try {
      await OpenIM.iMManager.conversationManager.markConversationMessageAsRead(
        conversationID: conversationID,
      );
    } catch (e) {
      debugPrint('[Chat] _markMessageAsRead request error: $e');
    } finally {
      msg.isRead = true;
      msg.hasReadTime = now;
      if (msg.clientMsgID != null) {
        await LocalStore.markMessageRead(conversationID, msg.clientMsgID!, now);
      }
      if (mounted) setState(() {});
    }
  }

  /// 与官方一致：滚动到第一条未读消息位置
  /// [firstUnreadIdx] 是未读消息在 _messages 列表中的索引
  /// 当前实现保留接口，必要时可扩展为搜索消息定位。
  void _scrollToFirstUnread(int firstUnreadIdx) {
    if (!_scrollController.hasClients || firstUnreadIdx < 0 || firstUnreadIdx >= _messages.length) {
      return;
    }
    final ratio = _messages.length <= 1 ? 1.0 : firstUnreadIdx / (_messages.length - 1);
    final target = _scrollController.position.maxScrollExtent * ratio;
    _scrollController.jumpTo(target.clamp(
      _scrollController.position.minScrollExtent,
      _scrollController.position.maxScrollExtent,
    ));
  }

  Future<void> _pickAndSendImage() async {
    final picker = ImagePicker();
    final file = await picker.pickImage(
      source: ImageSource.gallery,
      imageQuality: 70,
    );
    if (file == null) return;

    final tempMsg = await OpenIM.iMManager.messageManager
        .createImageMessageFromFullPath(imagePath: file.path);
    tempMsg.sessionType = sessionType;
    tempMsg.recvID = sessionType == ConversationType.single ? userID : '';
    tempMsg.groupID = sessionType != ConversationType.single ? groupID : '';
    tempMsg.isRead = false;

    setState(() => _messages.add(tempMsg));
    _scrollToBottom();

    try {
      // 尝试上传图片到服务器；失败时降级用本地路径
      try {
        final uploadedUrl = await UploadService.uploadFile(file.path, group: 'image');
        tempMsg.pictureElem?.sourcePicture?.url = uploadedUrl;
        tempMsg.pictureElem?.bigPicture?.url = uploadedUrl;
        tempMsg.pictureElem?.snapshotPicture?.url = uploadedUrl;
        debugPrint('[Chat] image uploaded: $uploadedUrl');
      } catch (e) {
        debugPrint('[Chat] image upload failed, sending with local path: $e');
      }

      await OpenIM.iMManager.messageManager.sendMessage(
        message: tempMsg,
        offlinePushInfo: OfflinePushInfo(),
        userID: sessionType == ConversationType.single ? userID : null,
        groupID: sessionType != ConversationType.single ? groupID : null,
      );
      if (mounted) {
        setState(() {
          tempMsg.status = MessageStatus.succeeded;
        });
        _scrollToBottom();
      }
    } catch (e) {
      if (mounted) setState(() => tempMsg.status = MessageStatus.failed);
      EasyLoading.showToast('发送失败');
      debugPrint('[Chat] image send error: $e');
    }
  }

  Future<void> _sendTextMessage() async {
    final text = _inputController.text.trim();
    if (text.isEmpty) return;
    _inputController.clear();

    final tempMsg = await OpenIM.iMManager.messageManager.createTextMessage(
      text: text,
    );
    tempMsg.sessionType = sessionType;
    tempMsg.recvID = sessionType == ConversationType.single ? userID : '';
    tempMsg.groupID = sessionType != ConversationType.single ? groupID : '';
    tempMsg.isRead = false;

    setState(() => _messages.add(tempMsg));
    _scrollToBottom();

    try {
      await OpenIM.iMManager.messageManager.sendMessage(
        message: tempMsg,
        offlinePushInfo: OfflinePushInfo(),
        userID: sessionType == ConversationType.single ? userID : null,
        groupID: sessionType != ConversationType.single ? groupID : null,
      );
      if (mounted) {
        setState(() {
          tempMsg.status = MessageStatus.succeeded;
        });
        _scrollToBottom();
      }
    } catch (e) {
      if (mounted) setState(() => tempMsg.status = MessageStatus.failed);
      EasyLoading.showToast('发送失败');
    }
  }

  Future<void> _toggleVoiceMode() async {
    if (_isRecording) {
      await _cancelVoiceRecord();
    }

    final nextVoiceMode = !_voiceMode;
    if (nextVoiceMode) {
      // 切换到语音模式
      setState(() {
        _voiceMode = true;
        _showEmojiPicker = false;
        _showMorePanel = false;
      });
      _focusNode.unfocus();
    } else {
      // 切换到文字模式
      setState(() {
        _voiceMode = false;
        _showEmojiPicker = false;
        _showMorePanel = false;
        _pendingScrollAfterKeyboardShows = true;
      });
      WidgetsBinding.instance.addPostFrameCallback((_) {
        if (!mounted) return;
        FocusScope.of(context).requestFocus(_focusNode);
        _scrollToBottom();
        Future.delayed(const Duration(milliseconds: 80), () {
          if (mounted) _scrollToBottom();
        });
        Future.delayed(const Duration(milliseconds: 220), () {
          if (mounted) _scrollToBottom();
        });
        Future.delayed(const Duration(milliseconds: 360), () {
          if (mounted) _scrollToBottom();
        });
      });
    }
  }

  Future<void> _startVoiceRecord() async {
    try {
      final hasPermission = await _audioRecorder.hasPermission();
      if (!hasPermission) {
        EasyLoading.showToast('请先开启麦克风权限');
        return;
      }

      final dir = await getTemporaryDirectory();
      final filePath = '${dir.path}/voice_${DateTime.now().millisecondsSinceEpoch}.m4a';
      await _audioRecorder.start(
        const RecordConfig(
          encoder: AudioEncoder.aacLc,
          bitRate: 128000,
          sampleRate: 44100,
        ),
        path: filePath,
      );

      if (!mounted) return;
      setState(() {
        _isRecording = true;
        _recordFilePath = filePath;
        _recordStartAt = DateTime.now().millisecondsSinceEpoch;
      });
    } catch (e) {
      EasyLoading.showToast('开始录音失败');
      debugPrint('[Chat] _startVoiceRecord error: $e');
    }
  }

  Future<void> _stopVoiceRecordAndSend() async {
    if (!_isRecording) return;
    try {
      final path = await _audioRecorder.stop();
      final startAt = _recordStartAt;
      final durationMs = startAt == null
          ? 0
          : DateTime.now().millisecondsSinceEpoch - startAt;
      final duration = (durationMs / 1000).ceil();

      if (!mounted) return;
      setState(() {
        _isRecording = false;
        _recordStartAt = null;
      });

      if (path == null || path.isEmpty || duration <= 0) {
        EasyLoading.showToast('录音失败');
        return;
      }
      if (duration < 1) {
        try {
          final file = File(path);
          if (await file.exists()) {
            await file.delete();
          }
        } catch (_) {}
        EasyLoading.showToast('说话时间太短');
        return;
      }

      final tempMsg = await OpenIM.iMManager.messageManager
          .createSoundMessageFromFullPath(
        soundPath: path,
        duration: duration,
      );
      tempMsg.sessionType = sessionType;
      tempMsg.recvID = sessionType == ConversationType.single ? userID : '';
      tempMsg.groupID = sessionType != ConversationType.single ? groupID : '';
      tempMsg.isRead = false;
      tempMsg.soundElem?.soundPath = path;
      tempMsg.soundElem?.sourceUrl = path;

      setState(() => _messages.add(tempMsg));
      _scrollToBottom();

      try {
        // 尝试上传音频文件到服务器，获取永久 URL；失败时降级用本地路径
        try {
          final uploadedUrl = await UploadService.uploadFile(path, group: 'voice');
          tempMsg.soundElem?.sourceUrl = uploadedUrl;
          debugPrint('[Chat] voice uploaded: $uploadedUrl');
        } catch (e) {
          debugPrint('[Chat] voice upload failed, sending with local path: $e');
        }

        await OpenIM.iMManager.messageManager.sendMessage(
          message: tempMsg,
          offlinePushInfo: OfflinePushInfo(),
          userID: sessionType == ConversationType.single ? userID : null,
          groupID: sessionType != ConversationType.single ? groupID : null,
        );
        if (mounted) {
          setState(() {
            tempMsg.status = MessageStatus.succeeded;
          });
          _scrollToBottom();
        }
      } catch (e) {
        if (mounted) setState(() => tempMsg.status = MessageStatus.failed);
        EasyLoading.showToast('语音发送失败');
        debugPrint('[Chat] voice send error: $e');
      }
    } catch (e) {
      if (mounted) {
        setState(() {
          _isRecording = false;
          _recordStartAt = null;
        });
      }
      EasyLoading.showToast('结束录音失败');
      debugPrint('[Chat] _stopVoiceRecordAndSend error: $e');
    }
  }

  Future<void> _cancelVoiceRecord() async {
    try {
      if (_isRecording) {
        final path = await _audioRecorder.stop();
        if (path != null && path.isNotEmpty) {
          final file = File(path);
          if (await file.exists()) {
            await file.delete();
          }
        }
      }
    } catch (_) {}
    if (!mounted) return;
    setState(() {
      _isRecording = false;
      _recordFilePath = null;
      _recordStartAt = null;
    });
  }

  @override
  Widget build(BuildContext context) {
    final keyboardInset = MediaQuery.of(context).viewInsets.bottom;
    final emojiPanelVisible = _showEmojiPicker || _showMorePanel || _hideEmojiWhenKeyboardShows;
    final emojiPanelAnimationDuration = _hideEmojiWhenKeyboardShows
        ? Duration.zero
        : const Duration(milliseconds: 180);
    final keyboardHostHeight = keyboardInset > 0
        ? (keyboardInset > _cachedKeyboardHeight
            ? keyboardInset
            : _cachedKeyboardHeight)
        : 0.0;
    final bottomPanelHeight = emojiPanelVisible
        ? _cachedKeyboardHeight
        : keyboardHostHeight;
    final bottomHostHeight = bottomPanelHeight;
    if (keyboardInset > 180 && keyboardInset > _cachedKeyboardHeight) {
      WidgetsBinding.instance.addPostFrameCallback((_) {
        if (!mounted) return;
        setState(() => _cachedKeyboardHeight = keyboardInset);
      });
    }
    if (_hideEmojiWhenKeyboardShows && keyboardInset > 0) {
      WidgetsBinding.instance.addPostFrameCallback((_) {
        if (!mounted) return;
        setState(() {
          _showEmojiPicker = false;
          _showMorePanel = false;
          _hideEmojiWhenKeyboardShows = false;
        });
      });
    }
    if (_pendingScrollAfterKeyboardShows && keyboardInset > 0) {
      WidgetsBinding.instance.addPostFrameCallback((_) {
        if (!mounted) return;
        _scrollToBottom();
        Future.delayed(const Duration(milliseconds: 80), () {
          if (mounted) _scrollToBottom();
        });
        setState(() => _pendingScrollAfterKeyboardShows = false);
      });
    }

    return Scaffold(
      resizeToAvoidBottomInset: false,
      appBar: AppBar(
        centerTitle: true,
        title: Column(
          crossAxisAlignment: CrossAxisAlignment.center,
          children: [
            Text(showName, style: const TextStyle(fontSize: 17)),
            if (sessionType == 1)
              Text(
                _isOnline ? '在线' : '离线',
                style: TextStyle(
                  fontSize: 12,
                  color: _isOnline ? Colors.green : Colors.grey,
                ),
              ),
          ],
        ),
        actions: [
          if (sessionType == ConversationType.single)
            IconButton(
              icon: const Icon(Icons.phone_outlined),
              onPressed: _startCall,
            ),
          IconButton(
            icon: const Icon(Icons.more_horiz),
            onPressed: () => Get.toNamed(
              AppRoutes.chatSettings,
              arguments: {
                'conversationID': conversationID,
                'userID': userID,
                'groupID': groupID,
                'sessionType': sessionType,
              },
            ),
          ),
        ],
      ),
      body: Stack(
        children: [
          Positioned.fill(
            child: _loading
                ? const Center(child: CircularProgressIndicator())
                : Stack(
                    children: [
                      GestureDetector(
                        onTap: () {
                          FocusScope.of(context).unfocus();
                          if (_showEmojiPicker || _showMorePanel) {
                            setState(() {
                              _showEmojiPicker = false;
                              _showMorePanel = false;
                            });
                          }
                        },
                        child: ListView.builder(
                          controller: _scrollController,
                          padding: EdgeInsets.only(
                            left: 12,
                            right: 12,
                            top: 8,
                            bottom: 88 + bottomHostHeight,
                          ),
                          itemCount: _messages.length,
                          itemBuilder: (_, i) {
                            final msg = _messages[i];
                            final showTime = _shouldShowTime(i);
                            return Column(
                              mainAxisSize: MainAxisSize.min,
                              children: [
                                if (showTime) _buildTimeSeparator(msg),
                                _buildMessageItem(msg),
                              ],
                            );
                          },
                        ),
                      ),
                    ],
                  ),
          ),
          Positioned(
            left: 0,
            right: 0,
            bottom: 0,
            child: Container(
              color: const Color(0xFFF0F2F6),
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  _buildInputBar(),
                  SizedBox(
                    height: bottomHostHeight,
                    child: Align(
                      alignment: Alignment.topCenter,
                      child: IgnorePointer(
                        ignoring: keyboardInset > 0,
                        child: AnimatedContainer(
                          duration: emojiPanelAnimationDuration,
                          curve: Curves.easeOutCubic,
                          height: emojiPanelVisible ? _cachedKeyboardHeight : 0,
                          child: ClipRect(
                            child: SizedBox(
                              height: _cachedKeyboardHeight,
                              child: _showMorePanel
                                  ? _buildMoreActionsPanel()
                                  : IgnorePointer(
                                      ignoring: !_showEmojiPicker,
                                      child: EmojiPickerSheet(
                                        onEmojiSelected: (emoji) {
                                          _inputController.text += emoji;
                                        },
                                        onBackspacePressed: () {
                                          final text = _inputController.text;
                                          if (text.isNotEmpty) {
                                            _inputController.text = text.substring(0, text.length - 1);
                                          }
                                        },
                                      ),
                                    ),
                            ),
                          ),
                        ),
                      ),
                    ),
                  ),
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildMessageItem(Message msg) {
    final isMe = msg.sendID == Config.userID;

    // 通知类型消息（contentType > 1000）：居中显示，无头像/已读标签（与官方一致）
    if ((msg.contentType ?? 0) > 1000) {
      return Padding(
        padding: const EdgeInsets.symmetric(vertical: 4),
        child: _buildBubble(msg, isMe),
      );
    }

    // 【修复问题3】使用 FocusDetector 检测消息可见性（与官方一致）
    return FocusDetector(
      onVisibilityGained: () {
        _markMessageAsRead(msg, true);
      },
      onVisibilityLost: () {
        _markMessageAsRead(msg, false);
      },
      child: Padding(
        padding: const EdgeInsets.symmetric(vertical: 4),
        child: Row(
          mainAxisAlignment: isMe
              ? MainAxisAlignment.end
              : MainAxisAlignment.start,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            if (!isMe)
              _buildAvatar(msg.senderFaceUrl ?? '', msg.senderNickname ?? ''),
            if (!isMe) const SizedBox(width: 8),
            Flexible(
              child: Column(
                crossAxisAlignment: isMe
                    ? CrossAxisAlignment.end
                    : CrossAxisAlignment.start,
                children: [
                  _buildBubble(msg, isMe),
                  // Show read status tag below message bubble (only for sender in single chat)
                  if (isMe &&
                      msg.status == MessageStatus.succeeded &&
                      sessionType == ConversationType.single)
                    _buildReadTag(msg),
                ],
              ),
            ),
            if (isMe) const SizedBox(width: 8),
            if (isMe) _buildAvatar(Config.faceURL, Config.nickname),
          ],
        ),
      ),
    );
  }

  /// 获取消息时间戳（毫秒），带归一化和回退
  static int _getMsgTimestamp(Message msg) {
    int ts = msg.sendTime ?? msg.createTime ?? 0;
    if (ts <= 0) return 0;
    // 秒级时间戳 < 10^12，毫秒级 >= 10^12
    return ts < 1000000000000 ? ts * 1000 : ts;
  }

  /// 是否在该消息前显示时间分隔：第一条始终显示；与上一条间隔 >= 5 分钟时显示
  bool _shouldShowTime(int index) {
    if (index == 0) return true;
    final cur = _getMsgTimestamp(_messages[index]);
    final prev = _getMsgTimestamp(_messages[index - 1]);
    // 如果任一时间戳无效，显示时间分隔
    if (cur <= 0 || prev <= 0) return true;
    final diff = (cur - prev).abs();
    final show = diff >= 5 * 60 * 1000; // 5 分钟
    debugPrint('[Chat] _shouldShowTime[$index]: cur=$cur prev=$prev diff=${diff}ms show=$show');
    return show;
  }

  Widget _buildTimeSeparator(Message msg) {
    int ms = _getMsgTimestamp(msg);
    // 如果无有效时间戳，用当前时间作为回退
    if (ms <= 0) ms = DateTime.now().millisecondsSinceEpoch;
    // 显式标记为 UTC，再转为手机本地时区
    final dt = DateTime.fromMillisecondsSinceEpoch(ms, isUtc: true).toLocal();
    final now = DateTime.now();
    final today = DateTime(now.year, now.month, now.day);
    final msgDay = DateTime(dt.year, dt.month, dt.day);
    final diff = today.difference(msgDay).inDays;
    String suffix = '';
    if (diff == 0) {
      suffix = '(今天)';
    } else if (diff == 1) {
      suffix = '(昨天)';
    }
    final text = '${dt.year}年${_pad(dt.month)}月${_pad(dt.day)}日 ${_pad(dt.hour)}:${_pad(dt.minute)}$suffix';
    debugPrint('[Chat] _buildTimeSeparator: ms=$ms localHour=${dt.hour} text=$text');
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8),
      child: Center(
        child: Text(
          text,
          style: const TextStyle(fontSize: 12, color: Color(0xFF999999)),
        ),
      ),
    );
  }

  static String _pad(int n) => n.toString().padLeft(2, '0');

  Widget _buildReadTag(Message msg) {
    final isRead = msg.isRead == true;
    return Padding(
      padding: const EdgeInsets.only(top: 2, bottom: 2),
      child: Text(
        isRead ? '已读' : '未读',
        style: const TextStyle(fontSize: 11, color: Color(0xFF999999)),
      ),
    );
  }

  Widget _buildAvatar(String url, String name) {
    return CircleAvatar(
      radius: 18,
      backgroundColor: Colors.grey[300],
      backgroundImage: url.isNotEmpty ? NetworkImage(url) : null,
      child: url.isEmpty
          ? Text(
              name.isNotEmpty ? name[0].toUpperCase() : '?',
              style: const TextStyle(color: Colors.white, fontSize: 14),
            )
          : null,
    );
  }

  Future<void> _showMessageMenu(
    Message msg,
    bool isMe,
    LongPressStartDetails details,
  ) async {
    final canCopy = msg.contentType == MessageType.text && msg.textContent.isNotEmpty;
    final msgSendTime = msg.sendTime ?? msg.createTime ?? 0;
    final msgTimeMs = msgSendTime < 1000000000000 ? msgSendTime * 1000 : msgSendTime;
    final canRevoke = isMe;
    if (!canCopy && !canRevoke) return;

    final overlay = Overlay.of(context).context.findRenderObject() as RenderBox;
    final selected = await showMenu<String>(
      context: context,
      position: RelativeRect.fromRect(
        Rect.fromCenter(
          center: details.globalPosition,
          width: 1,
          height: 1,
        ),
        Offset.zero & overlay.size,
      ),
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
      items: [
        if (canCopy)
          const PopupMenuItem<String>(
            value: 'copy',
            child: Text('复制'),
          ),
        if (canRevoke)
          const PopupMenuItem<String>(
            value: 'revoke',
            child: Text('撤回'),
          ),
      ],
    );

    if (!mounted || selected == null) return;
    if (selected == 'copy') {
      Clipboard.setData(ClipboardData(text: msg.textContent));
      EasyLoading.showToast('已复制');
      return;
    }
    if (selected == 'revoke') {
      final canRevokeNow =
          (DateTime.now().millisecondsSinceEpoch - msgTimeMs) <=
          const Duration(minutes: 1).inMilliseconds;
      if (!canRevokeNow) {
        debugPrint(
          '[Revoke][Client] blocked by time window: '
          'clientMsgID=${msg.clientMsgID} seq=${msg.seq} '
          'sendTime=$msgTimeMs now=${DateTime.now().millisecondsSinceEpoch}',
        );
        EasyLoading.showToast('仅支持1分钟内撤回');
        return;
      }
      if ((msg.seq ?? 0) <= 0) {
        debugPrint(
          '[Revoke][Client] blocked by missing seq: '
          'clientMsgID=${msg.clientMsgID} serverMsgID=${msg.serverMsgID} '
          'seq=${msg.seq} sendTime=${msg.sendTime}',
        );
        EasyLoading.showToast('消息尚未同步完成，暂时不能撤回');
        return;
      }
      try {
        debugPrint(
          '[Revoke][Client] request start: '
          'conversationID=$conversationID clientMsgID=${msg.clientMsgID} '
          'serverMsgID=${msg.serverMsgID} seq=${msg.seq}',
        );
        await OpenIM.iMManager.messageManager.revokeMessage(
          conversationID: conversationID,
          seq: msg.seq!,
        );
        debugPrint(
          '[Revoke][Client] request success: '
          'conversationID=$conversationID clientMsgID=${msg.clientMsgID} seq=${msg.seq}',
        );
        // 立即本地更新（不等服务端推送回来），与官方 SDK 行为一致
        _onMessageRevoked(RevokedInfo(
          revokerID: Config.userID,
          clientMsgID: msg.clientMsgID,
          revokeTime: DateTime.now().millisecondsSinceEpoch,
          sessionType: msg.sessionType,
          sourceMessageSendID: msg.sendID,
          sourceMessageSenderNickname: msg.senderNickname,
          sourceMessageSendTime: msg.sendTime,
        ));
      } catch (e) {
        debugPrint(
          '[Revoke][Client] request failed: '
          'conversationID=$conversationID clientMsgID=${msg.clientMsgID} '
          'seq=${msg.seq} error=$e',
        );
        EasyLoading.showToast('撤回失败');
      }
    }
  }

  /// 底部更多操作面板（微信风格：相册/拍摄/位置/语音通话/视频通话/红包/转账）
  Widget _buildMoreActionsPanel() {
    final isSingle = sessionType == ConversationType.single;
    final actions = <_MoreAction>[
      _MoreAction(icon: Icons.photo_library_outlined, label: '相册', onTap: _pickAndSendImage),
      _MoreAction(icon: Icons.camera_alt_outlined, label: '拍摄', onTap: _takePhotoAndSend),
      _MoreAction(icon: Icons.location_on_outlined, label: '位置', onTap: () {
        EasyLoading.showToast('位置功能开发中');
      }),
      if (isSingle)
        _MoreAction(icon: Icons.phone_outlined, label: '语音通话', onTap: () => _directCall(CallType.audio)),
      if (isSingle)
        _MoreAction(icon: Icons.videocam_outlined, label: '视频通话', onTap: () => _directCall(CallType.video)),
      _MoreAction(icon: Icons.redeem_outlined, label: '红包', onTap: () {
        EasyLoading.showToast('红包功能开发中');
      }),
      _MoreAction(icon: Icons.swap_horiz, label: '转账', onTap: () {
        EasyLoading.showToast('转账功能开发中');
      }),
    ];

    return Container(
      color: const Color(0xFFF0F2F6),
      padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 16),
      child: GridView.builder(
        physics: const NeverScrollableScrollPhysics(),
        gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
          crossAxisCount: 4,
          mainAxisSpacing: 12,
          crossAxisSpacing: 16,
          childAspectRatio: 0.85,
        ),
        itemCount: actions.length,
        itemBuilder: (_, i) {
          final action = actions[i];
          return GestureDetector(
            onTap: () {
              setState(() => _showMorePanel = false);
              action.onTap();
            },
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                Container(
                  width: 56,
                  height: 56,
                  decoration: BoxDecoration(
                    color: Colors.white,
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: Icon(action.icon, size: 28, color: const Color(0xFF4A5568)),
                ),
                const SizedBox(height: 6),
                Text(
                  action.label,
                  style: const TextStyle(fontSize: 11, color: Color(0xFF8E9AB0)),
                  overflow: TextOverflow.ellipsis,
                ),
              ],
            ),
          );
        },
      ),
    );
  }

  /// 拍摄照片并发送
  Future<void> _takePhotoAndSend() async {
    final picker = ImagePicker();
    final photo = await picker.pickImage(source: ImageSource.camera, imageQuality: 70);
    if (photo == null) return;

    final tempMsg = await OpenIM.iMManager.messageManager
        .createImageMessageFromFullPath(imagePath: photo.path);
    tempMsg.sessionType = sessionType;
    tempMsg.recvID = sessionType == ConversationType.single ? userID : '';
    tempMsg.groupID = sessionType != ConversationType.single ? groupID : '';
    tempMsg.isRead = false;

    setState(() => _messages.add(tempMsg));
    _scrollToBottom();

    try {
      // 尝试上传图片到服务器；失败时降级用本地路径
      try {
        final uploadedUrl = await UploadService.uploadFile(photo.path, group: 'image');
        tempMsg.pictureElem?.sourcePicture?.url = uploadedUrl;
        tempMsg.pictureElem?.bigPicture?.url = uploadedUrl;
        tempMsg.pictureElem?.snapshotPicture?.url = uploadedUrl;
        debugPrint('[Chat] camera image uploaded: $uploadedUrl');
      } catch (e) {
        debugPrint('[Chat] camera image upload failed, sending with local path: $e');
      }

      await OpenIM.iMManager.messageManager.sendMessage(
        message: tempMsg,
        offlinePushInfo: OfflinePushInfo(),
        userID: sessionType == ConversationType.single ? userID : null,
        groupID: sessionType != ConversationType.single ? groupID : null,
      );
      if (mounted) {
        setState(() => tempMsg.status = MessageStatus.succeeded);
        _scrollToBottom();
      }
    } catch (e) {
      if (mounted) setState(() => tempMsg.status = MessageStatus.failed);
      EasyLoading.showToast('发送失败');
      debugPrint('[Chat] camera send error: $e');
    }
  }

  /// 直接发起指定类型的通话（从更多面板调用）
  void _directCall(CallType callType) {
    debugPrint('[ChatPage] _directCall: callType=$callType, userID=$userID, sessionType=$sessionType');
    if (sessionType != ConversationType.single || userID.isEmpty) {
      EasyLoading.showToast('仅支持单聊通话');
      return;
    }
    final imCtrl = Get.find<IMController>();
    if (imCtrl.isRtcBusy) {
      EasyLoading.showToast('当前正在通话中');
      return;
    }
    imCtrl.call(callType: callType, inviteeUserIDList: [userID]);
  }

  /// 发起通话（弹框选择视频/语音，AppBar 调用）
  void _startCall() {
    if (sessionType != ConversationType.single || userID.isEmpty) {
      EasyLoading.showToast('仅支持单聊通话');
      return;
    }
    final imCtrl = Get.find<IMController>();
    if (imCtrl.isRtcBusy) {
      EasyLoading.showToast('当前正在通话中');
      return;
    }
    showModalBottomSheet(
      context: context,
      backgroundColor: Colors.white,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(16)),
      ),
      builder: (_) => SafeArea(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            ListTile(
              title: Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: const [
                  Icon(Icons.videocam_outlined, color: Colors.blue, size: 22),
                  SizedBox(width: 8),
                  Text('视频通话', style: TextStyle(fontSize: 16)),
                ],
              ),
              onTap: () {
                Navigator.pop(context);
                imCtrl.call(callType: CallType.video, inviteeUserIDList: [userID]);
              },
            ),
            const Divider(height: 1),
            ListTile(
              title: Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: const [
                  Icon(Icons.phone_outlined, color: Colors.blue, size: 22),
                  SizedBox(width: 8),
                  Text('语音通话', style: TextStyle(fontSize: 16)),
                ],
              ),
              onTap: () {
                Navigator.pop(context);
                imCtrl.call(callType: CallType.audio, inviteeUserIDList: [userID]);
              },
            ),
            const Divider(height: 1),
            ListTile(
              title: const Text('取消', textAlign: TextAlign.center, style: TextStyle(fontSize: 16)),
              onTap: () => Navigator.pop(context),
            ),
          ],
        ),
      ),
    );
  }

  /// 从历史拉取的撤回消息中解析显示文本
  String _parseRevokeDisplayText(Message msg) {
    try {
      final detail = msg.notificationElem?.detail;
      if (detail != null && detail.isNotEmpty) {
        final map = jsonDecode(detail);
        if (map is Map<String, dynamic>) {
          final revokerID = map['revokerID'] as String?;
          if (revokerID == Config.userID) {
            return '你撤回了一条消息';
          }
          return '对方撤回了一条消息';
        }
      }
    } catch (_) {}
    // 无法解析时根据 sendID 判断
    if (msg.sendID == Config.userID) {
      return '你撤回了一条消息';
    }
    return '对方撤回了一条消息';
  }

  Widget _buildBubble(Message msg, bool isMe) {
    if ((msg.contentType ?? 0) > 1000) {
      String displayText = msg.textContent;

      // 撤回消息：从 notificationElem.detail 解析出 RevokedInfo，显示友好文本
      if (msg.contentType == MessageType.revokeMessageNotification) {
        if (msg.textElem?.content != null && msg.textElem!.content!.isNotEmpty) {
          displayText = msg.textElem!.content!;
        } else {
          displayText = _parseRevokeDisplayText(msg);
        }
      }

      return Center(
        child: Container(
          margin: const EdgeInsets.symmetric(vertical: 8),
          padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 4),
          decoration: BoxDecoration(
            color: Colors.grey[200],
            borderRadius: BorderRadius.circular(4),
          ),
          child: Text(
            displayText,
            style: const TextStyle(fontSize: 12, color: Colors.grey),
          ),
        ),
      );
    }

    // 通话结果消息（customType == 210）
    if (msg.contentType == MessageType.custom) {
      try {
        final customData = msg.customElem?.data;
        if (customData != null) {
          final map = jsonDecode(customData);
          if (map['customType'] == CustomMessageType.callResult) {
            final info = CallResultInfo.fromJson(map['data']);
            final text = info.getDisplayText(isMe);
            return GestureDetector(
              onTap: () => _directCall(info.isVideo ? CallType.video : CallType.audio),
              onLongPressStart: (details) => _showMessageMenu(msg, isMe, details),
              child: Container(
                padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
                decoration: BoxDecoration(
                  color: isMe ? const Color(0xFF0089FF) : Colors.white,
                  borderRadius: BorderRadius.circular(8),
                ),
                child: Row(
                  mainAxisSize: MainAxisSize.min,
                  children: isMe
                      ? [
                          Text(text, style: TextStyle(fontSize: 15, color: isMe ? Colors.white : Colors.black87)),
                          const SizedBox(width: 6),
                          Icon(
                            info.isVideo ? Icons.videocam : Icons.phone,
                            size: 20,
                            color: isMe ? Colors.white : const Color(0xFF0089FF),
                          ),
                        ]
                      : [
                          Icon(
                            info.isVideo ? Icons.videocam : Icons.phone,
                            size: 20,
                            color: const Color(0xFF0089FF),
                          ),
                          const SizedBox(width: 6),
                          Text(text, style: const TextStyle(fontSize: 15, color: Colors.black87)),
                        ],
                ),
              ),
            );
          }
        }
      } catch (_) {}
    }

    // Image message
    if (msg.contentType == MessageType.picture) {
      final url =
          msg.pictureElem?.sourcePicture?.url ??
          msg.pictureElem?.snapshotPicture?.url ??
          '';
      final width = msg.pictureElem?.sourcePicture?.width?.toDouble() ?? 200;
      final height = msg.pictureElem?.sourcePicture?.height?.toDouble() ?? 200;

      // Calculate display size maintaining aspect ratio
      final maxWidth = MediaQuery.of(context).size.width * 0.6;
      final maxHeight = 300.0;
      double displayWidth = width;
      double displayHeight = height;

      if (displayWidth > maxWidth) {
        displayHeight = displayHeight * (maxWidth / displayWidth);
        displayWidth = maxWidth;
      }
      if (displayHeight > maxHeight) {
        displayWidth = displayWidth * (maxHeight / displayHeight);
        displayHeight = maxHeight;
      }

      return GestureDetector(
        onTap: () {
          Navigator.push(
            context,
            MaterialPageRoute(
              builder: (context) => MediaPreviewPage(
                url: url,
                isVideo: false,
                heroTag: msg.clientMsgID,
              ),
            ),
          );
        },
        onLongPressStart: (details) => _showMessageMenu(msg, isMe, details),
        child: Container(
          width: displayWidth,
          height: displayHeight,
          constraints: BoxConstraints(
            minWidth: 100,
            minHeight: 100,
            maxWidth: maxWidth,
            maxHeight: maxHeight,
          ),
          decoration: BoxDecoration(
            borderRadius: BorderRadius.circular(8),
            boxShadow: [
              BoxShadow(
                color: Colors.black.withValues(alpha: 0.1),
                blurRadius: 4,
                offset: const Offset(0, 2),
              ),
            ],
          ),
          child: ClipRRect(
            borderRadius: BorderRadius.circular(8),
            child: url.isNotEmpty
                ? Image.network(
                    url,
                    fit: BoxFit.cover,
                    loadingBuilder: (context, child, loadingProgress) {
                      if (loadingProgress == null) return child;
                      return Center(
                        child: CircularProgressIndicator(
                          value: loadingProgress.expectedTotalBytes != null
                              ? loadingProgress.cumulativeBytesLoaded /
                                    loadingProgress.expectedTotalBytes!
                              : null,
                        ),
                      );
                    },
                    errorBuilder: (context, error, stackTrace) => Container(
                      color: Colors.grey[200],
                      child: const Center(
                        child: Icon(
                          Icons.broken_image,
                          size: 48,
                          color: Colors.grey,
                        ),
                      ),
                    ),
                  )
                : Container(
                    color: Colors.grey[200],
                    child: const Center(
                      child: Icon(Icons.image, size: 48, color: Colors.grey),
                    ),
                  ),
          ),
        ),
      );
    }

    // Video message
    if (msg.contentType == MessageType.video) {
      final snapshotUrl = msg.videoElem?.snapshotUrl ?? '';
      final duration = msg.videoElem?.duration ?? 0;
      final width = msg.videoElem?.snapshotWidth?.toDouble() ?? 200;
      final height = msg.videoElem?.snapshotHeight?.toDouble() ?? 200;

      // Calculate display size maintaining aspect ratio
      final maxWidth = MediaQuery.of(context).size.width * 0.6;
      final maxHeight = 300.0;
      double displayWidth = width;
      double displayHeight = height;

      if (displayWidth > maxWidth) {
        displayHeight = displayHeight * (maxWidth / displayWidth);
        displayWidth = maxWidth;
      }
      if (displayHeight > maxHeight) {
        displayWidth = displayWidth * (maxHeight / displayHeight);
        displayHeight = maxHeight;
      }

      return GestureDetector(
        onTap: () {
          final videoUrl = msg.videoElem?.videoUrl ?? '';
          if (videoUrl.isNotEmpty) {
            Navigator.push(
              context,
              MaterialPageRoute(
                builder: (context) => MediaPreviewPage(
                  url: videoUrl,
                  isVideo: true,
                  heroTag: msg.clientMsgID,
                ),
              ),
            );
          } else {
            EasyLoading.showToast('视频地址无效');
          }
        },
        onLongPressStart: (details) => _showMessageMenu(msg, isMe, details),
        child: Container(
          width: displayWidth,
          height: displayHeight,
          constraints: BoxConstraints(
            minWidth: 100,
            minHeight: 100,
            maxWidth: maxWidth,
            maxHeight: maxHeight,
          ),
          decoration: BoxDecoration(
            borderRadius: BorderRadius.circular(8),
            boxShadow: [
              BoxShadow(
                color: Colors.black.withValues(alpha: 0.1),
                blurRadius: 4,
                offset: const Offset(0, 2),
              ),
            ],
          ),
          child: ClipRRect(
            borderRadius: BorderRadius.circular(8),
            child: Stack(
              fit: StackFit.expand,
              children: [
                snapshotUrl.isNotEmpty
                    ? Image.network(
                        snapshotUrl,
                        fit: BoxFit.cover,
                        errorBuilder: (context, error, stackTrace) => Container(
                          color: Colors.grey[200],
                          child: const Center(
                            child: Icon(
                              Icons.videocam,
                              size: 48,
                              color: Colors.grey,
                            ),
                          ),
                        ),
                      )
                    : Container(
                        color: Colors.grey[200],
                        child: const Center(
                          child: Icon(
                            Icons.videocam,
                            size: 48,
                            color: Colors.grey,
                          ),
                        ),
                      ),
                // Play button overlay
                Center(
                  child: Container(
                    width: 50,
                    height: 50,
                    decoration: BoxDecoration(
                      color: Colors.black.withValues(alpha: 0.5),
                      shape: BoxShape.circle,
                    ),
                    child: const Icon(
                      Icons.play_arrow,
                      color: Colors.white,
                      size: 32,
                    ),
                  ),
                ),
                // Duration label
                if (duration > 0)
                  Positioned(
                    right: 8,
                    bottom: 8,
                    child: Container(
                      padding: const EdgeInsets.symmetric(
                        horizontal: 6,
                        vertical: 2,
                      ),
                      decoration: BoxDecoration(
                        color: Colors.black.withValues(alpha: 0.6),
                        borderRadius: BorderRadius.circular(4),
                      ),
                      child: Text(
                        _formatDuration(duration),
                        style: const TextStyle(
                          color: Colors.white,
                          fontSize: 11,
                        ),
                      ),
                    ),
                  ),
              ],
            ),
          ),
        ),
      );
    }

    // Voice/Sound message
    if (msg.contentType == MessageType.voice) {
      final duration = msg.soundElem?.duration ?? 0;
      final msgId = msg.clientMsgID ?? '';
      final isPlayed = isMe || _playedVoiceIds.contains(msgId);
      // 优先用 sourceUrl，若为空或本地路径不存在则回退到 soundPath
      final sourceUrl = msg.soundElem?.sourceUrl ?? '';
      final soundPath = msg.soundElem?.soundPath ?? '';
      final playKey = sourceUrl.isNotEmpty ? sourceUrl : soundPath;
      final iconColor = isMe ? Colors.white : const Color(0xFF0089FF);

      return GestureDetector(
        onTap: () async {
          if (playKey.isNotEmpty) {
            final ok = await _audioPlayer.playOrPause(playKey);
            if (ok) {
              if (!isMe && msgId.isNotEmpty && !_playedVoiceIds.contains(msgId)) {
                setState(() => _playedVoiceIds.add(msgId));
                LocalStore.addPlayedVoiceId(conversationID, msgId);
              }
            } else {
              EasyLoading.showToast('语音文件不可用');
            }
          } else {
            EasyLoading.showToast('语音地址无效');
          }
        },
        onLongPressStart: (details) => _showMessageMenu(msg, isMe, details),
        child: Row(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.center,
          children: [
            Container(
              constraints: BoxConstraints(
                minWidth: 100,
                maxWidth: MediaQuery.of(context).size.width * 0.5,
              ),
              padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
              decoration: BoxDecoration(
                color: isMe ? const Color(0xFF0089FF) : Colors.white,
                borderRadius: BorderRadius.circular(8),
                boxShadow: [
                  BoxShadow(
                    color: Colors.black.withValues(alpha: 0.05),
                    blurRadius: 4,
                    offset: const Offset(0, 2),
                  ),
                ],
              ),
              child: ValueListenableBuilder<String?>(
                valueListenable: _audioPlayer.playingStateNotifier,
                builder: (_, playingUrl, __) {
                  final isPlayingThis = playingUrl == playKey && playKey.isNotEmpty;
                  return Row(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      if (isPlayingThis)
                        _VoiceWaveWidget(color: iconColor)
                      else
                        Icon(Icons.mic, size: 20, color: iconColor),
                      const SizedBox(width: 8),
                      Text(
                        '${duration}s',
                        style: TextStyle(
                          fontSize: 14,
                          color: isMe ? Colors.white : Colors.black87,
                        ),
                      ),
                    ],
                  );
                },
              ),
            ),
            if (!isPlayed) ...[
              const SizedBox(width: 6),
              Container(
                width: 8,
                height: 8,
                decoration: const BoxDecoration(
                  color: Colors.red,
                  shape: BoxShape.circle,
                ),
              ),
            ],
          ],
        ),
      );
    }

    // File message
    if (msg.contentType == MessageType.file) {
      final fileName = msg.fileElem?.fileName ?? '文件';
      final fileSize = msg.fileElem?.fileSize ?? 0;
      return GestureDetector(
        onTap: () async {
          final fileUrl = msg.fileElem?.sourceUrl ?? '';
          if (fileUrl.isNotEmpty) {
            await FileUtils.downloadAndOpen(
              context: context,
              url: fileUrl,
              fileName: fileName,
            );
          } else {
            EasyLoading.showToast('文件地址无效');
          }
        },
        onLongPressStart: (details) => _showMessageMenu(msg, isMe, details),
        child: Container(
          constraints: BoxConstraints(
            maxWidth: MediaQuery.of(context).size.width * 0.65,
          ),
          padding: const EdgeInsets.all(12),
          decoration: BoxDecoration(
            color: isMe ? const Color(0xFF0089FF) : Colors.white,
            borderRadius: BorderRadius.circular(8),
            boxShadow: [
              BoxShadow(
                color: Colors.black.withValues(alpha: 0.05),
                blurRadius: 4,
                offset: const Offset(0, 2),
              ),
            ],
          ),
          child: Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              Icon(
                Icons.insert_drive_file,
                size: 32,
                color: isMe ? Colors.white : const Color(0xFF0089FF),
              ),
              const SizedBox(width: 12),
              Flexible(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      fileName,
                      style: TextStyle(
                        fontSize: 14,
                        color: isMe ? Colors.white : Colors.black87,
                      ),
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                    const SizedBox(height: 2),
                    Text(
                      _formatFileSize(fileSize),
                      style: TextStyle(
                        fontSize: 12,
                        color: isMe ? Colors.white70 : Colors.grey,
                      ),
                    ),
                  ],
                ),
              ),
            ],
          ),
        ),
      );
    }

    // Text message (default)
    return GestureDetector(
      onLongPressStart: (details) => _showMessageMenu(msg, isMe, details),
      child: Container(
        constraints: BoxConstraints(
          maxWidth: MediaQuery.of(context).size.width * 0.65,
        ),
        padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
        decoration: BoxDecoration(
          color: isMe ? const Color(0xFF0089FF) : Colors.white,
          borderRadius: BorderRadius.only(
            topLeft: const Radius.circular(12),
            topRight: const Radius.circular(12),
            bottomLeft: Radius.circular(isMe ? 12 : 2),
            bottomRight: Radius.circular(isMe ? 2 : 12),
          ),
          boxShadow: [
            BoxShadow(
              color: Colors.black.withValues(alpha: 0.05),
              blurRadius: 4,
              offset: const Offset(0, 2),
            ),
          ],
        ),
        child: Row(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.end,
          children: [
            Flexible(
              child: Text(
                msg.textContent,
                style: TextStyle(
                  fontSize: 15,
                  color: isMe ? Colors.white : Colors.black87,
                ),
              ),
            ),
            if (msg.status == MessageStatus.sending) ...[
              const SizedBox(width: 4),
              const SizedBox(
                width: 12,
                height: 12,
                child: CircularProgressIndicator(strokeWidth: 1.5),
              ),
            ],
            if (msg.status == MessageStatus.failed) ...[
              const SizedBox(width: 4),
              const Icon(Icons.error, size: 14, color: Colors.red),
            ],
          ],
        ),
      ),
    );
  }

  Widget _buildInputBar() {
    final attachedToBottomPanel = _showEmojiPicker || _showMorePanel || _hideEmojiWhenKeyboardShows;
    final safeBottom = MediaQuery.of(context).padding.bottom;
    final bottomPadding = attachedToBottomPanel ? 8.0 : safeBottom + 8;
    return Container(
      padding: EdgeInsets.only(
        left: 12,
        right: 8,
        top: 8,
        bottom: bottomPadding,
      ),
      decoration: const BoxDecoration(
        color: Colors.transparent,
      ),
      child: Row(
        children: [
          GestureDetector(
            onTap: _toggleVoiceMode,
            child: Container(
              width: 32,
              height: 32,
              decoration: const BoxDecoration(color: Colors.transparent),
              child: Icon(
                _voiceMode ? Icons.keyboard_alt_outlined : Icons.mic_none,
                color: const Color(0xFF8E9AB0),
                size: 24,
              ),
            ),
          ),
          const SizedBox(width: 8),
          Expanded(
            child: _voiceMode
                ? GestureDetector(
                    onLongPressStart: (_) => _startVoiceRecord(),
                    onLongPressEnd: (_) => _stopVoiceRecordAndSend(),
                    onLongPressCancel: _cancelVoiceRecord,
                    child: Container(
                      height: 40,
                      alignment: Alignment.center,
                      decoration: BoxDecoration(
                        color: _isRecording ? const Color(0xFFE8F3FF) : Colors.white,
                        borderRadius: BorderRadius.circular(20),
                        border: Border.all(
                          color: _isRecording
                              ? const Color(0xFF0089FF)
                              : const Color(0xFFD9DEE8),
                        ),
                      ),
                      child: Text(
                        _isRecording ? '松开 发送语音' : '按住 说话',
                        style: TextStyle(
                          fontSize: 15,
                          color: _isRecording
                              ? const Color(0xFF0089FF)
                              : const Color(0xFF0C1C33),
                          fontWeight: FontWeight.w500,
                        ),
                      ),
                    ),
                  )
                : Container(
                    decoration: BoxDecoration(
                      color: Colors.white,
                      borderRadius: BorderRadius.circular(20),
                    ),
                    child: TextField(
                      controller: _inputController,
                      focusNode: _focusNode,
                      maxLines: 4,
                      minLines: 1,
                      decoration: const InputDecoration(
                        hintText: '输入消息...',
                        hintStyle: TextStyle(color: Color(0xFF8E9AB0)),
                        border: InputBorder.none,
                        contentPadding: EdgeInsets.symmetric(
                          horizontal: 16,
                          vertical: 8,
                        ),
                      ),
                      style: const TextStyle(fontSize: 17, color: Color(0xFF0C1C33)),
                      textInputAction: TextInputAction.send,
                      onSubmitted: (_) => _sendTextMessage(),
                      onTap: () {
                        if (_showEmojiPicker || _showMorePanel) {
                          setState(() {
                            _voiceMode = false;
                            _hideEmojiWhenKeyboardShows = true;
                            _pendingScrollAfterKeyboardShows = true;
                          });
                          FocusScope.of(context).requestFocus(_focusNode);
                          return;
                        }
                        if (!_focusNode.hasFocus) {
                          setState(() => _pendingScrollAfterKeyboardShows = true);
                          WidgetsBinding.instance.addPostFrameCallback((_) {
                            if (mounted) {
                              FocusScope.of(context).requestFocus(_focusNode);
                            }
                          });
                        }
                      },
                    ),
                  ),
          ),
          const SizedBox(width: 8),
          GestureDetector(
            onTap: () {
              if (_showEmojiPicker) {
                setState(() {
                  _voiceMode = false;
                  _hideEmojiWhenKeyboardShows = true;
                });
                FocusScope.of(context).requestFocus(_focusNode);
                return;
              }

              _focusNode.unfocus();
              setState(() {
                _voiceMode = false;
                _showMorePanel = false;
                _hideEmojiWhenKeyboardShows = false;
                _showEmojiPicker = true;
              });
              _scrollToBottom();
            },
            child: Container(
              width: 32,
              height: 32,
              decoration: const BoxDecoration(color: Colors.transparent),
              child: Icon(
                _showEmojiPicker
                    ? Icons.keyboard
                    : Icons.emoji_emotions_outlined,
                color: const Color(0xFF8E9AB0),
                size: 24,
              ),
            ),
          ),
          const SizedBox(width: 8),
          GestureDetector(
            onTap: _sendButtonVisible
                ? _sendTextMessage
                : () {
                    if (_showMorePanel) {
                      // 已打开 more 面板，切回键盘
                      setState(() {
                        _showMorePanel = false;
                        _hideEmojiWhenKeyboardShows = true;
                      });
                      FocusScope.of(context).requestFocus(_focusNode);
                      return;
                    }
                    _focusNode.unfocus();
                    setState(() {
                      _voiceMode = false;
                      _showEmojiPicker = false;
                      _hideEmojiWhenKeyboardShows = false;
                      _showMorePanel = true;
                    });
                    _scrollToBottom();
                  },
            child: Container(
              width: 32,
              height: 32,
              decoration: const BoxDecoration(color: Colors.transparent),
              child: Icon(
                _sendButtonVisible
                    ? Icons.send
                    : Icons.add_circle_outline,
                color: _sendButtonVisible
                    ? const Color(0xFF0089FF)
                    : const Color(0xFF8E9AB0),
                size: 24,
              ),
            ),
          ),
          const SizedBox(width: 8),
        ],
      ),
    );
  }

  String _formatDuration(int seconds) {
    final minutes = seconds ~/ 60;
    final secs = seconds % 60;
    return '${minutes.toString().padLeft(2, '0')}:${secs.toString().padLeft(2, '0')}';
  }

  String _formatFileSize(int bytes) {
    if (bytes < 1024) {
      return '${bytes}B';
    } else if (bytes < 1024 * 1024) {
      return '${(bytes / 1024).toStringAsFixed(1)}KB';
    } else if (bytes < 1024 * 1024 * 1024) {
      return '${(bytes / (1024 * 1024)).toStringAsFixed(1)}MB';
    } else {
      return '${(bytes / (1024 * 1024 * 1024)).toStringAsFixed(1)}GB';
    }
  }
}

class _MoreAction {
  final IconData icon;
  final String label;
  final VoidCallback onTap;
  _MoreAction({required this.icon, required this.label, required this.onTap});
}

/// 语音播放时的声波动画（3根竖条交替波动）
class _VoiceWaveWidget extends StatefulWidget {
  final Color color;
  const _VoiceWaveWidget({required this.color});

  @override
  State<_VoiceWaveWidget> createState() => _VoiceWaveWidgetState();
}

class _VoiceWaveWidgetState extends State<_VoiceWaveWidget>
    with SingleTickerProviderStateMixin {
  late AnimationController _ctrl;

  @override
  void initState() {
    super.initState();
    _ctrl = AnimationController(
      vsync: this,
      duration: const Duration(milliseconds: 600),
    )..repeat(reverse: true);
  }

  @override
  void dispose() {
    _ctrl.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return AnimatedBuilder(
      animation: _ctrl,
      builder: (_, __) {
        return Row(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.center,
          children: List.generate(3, (i) {
            // 每根bar有不同的相位偏移，产生波浪效果
            final phase = (i * 0.3).clamp(0.0, 1.0);
            final t = (_ctrl.value + phase) % 1.0;
            final height = 6.0 + 10.0 * (0.5 + 0.5 * _sin(t));
            return Container(
              width: 3,
              height: height,
              margin: EdgeInsets.only(left: i == 0 ? 0 : 2),
              decoration: BoxDecoration(
                color: widget.color,
                borderRadius: BorderRadius.circular(1.5),
              ),
            );
          }),
        );
      },
    );
  }

  double _sin(double t) {
    // 简单正弦近似：t 在 0~1 之间
    final x = t * 3.14159 * 2;
    return (x - x * x * x / 6 + x * x * x * x * x / 120).clamp(-1.0, 1.0);
  }
}
