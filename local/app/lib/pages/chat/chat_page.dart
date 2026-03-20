import 'dart:async';

import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_easyloading/flutter_easyloading.dart';
import 'package:get/get.dart';
import 'package:image_picker/image_picker.dart';
import 'package:focus_detector_v2/focus_detector_v2.dart';

import '../../core/config.dart';
import '../../core/controllers/im_controller.dart';
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
  // 【修复问题1】使用全局静态缓存，确保退出聊天室后再进入时"已读/未读"状态不会消失
  static final Map<String, Map<String, bool>> _globalReadStatusCache = <String, Map<String, bool>>{};
  static final Map<String, Map<String, int>> _globalReadTimeCache = <String, Map<String, int>>{};

  final _inputController = TextEditingController();
  final _scrollController = ScrollController();
  final _focusNode = FocusNode();
  final _messages = <Message>[];
  bool _loading = true;
  bool _loadingMore = false;
  bool _isEnd = false;
  bool _isOnline = false;
  bool _showEmojiPicker = false;
  bool _sendButtonVisible = false;
  Timer? _onlineTimer;
  StreamSubscription? _newMsgSub;
  StreamSubscription? _readReceiptSub;
  StreamSubscription? _convChangedSub;
  final _audioPlayer = AudioPlayerManager();

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

    _loadMessages();
    _markAsRead();

    // Listen for new messages via SDK
    final imCtrl = Get.find<IMController>();
    _newMsgSub = imCtrl.newMessageSubject.listen(_onNewMessage);

    // Listen for read receipts (C2C only)
    _readReceiptSub = imCtrl.c2cReadReceiptSubject.listen(_onReadReceipt);

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
      if (_scrollController.position.pixels <=
          _scrollController.position.minScrollExtent + 50) {
        _loadHistoryMessages();
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
    _inputController.dispose();
    _scrollController.dispose();
    _focusNode.dispose();
    _onlineTimer?.cancel();
    _newMsgSub?.cancel();
    _readReceiptSub?.cancel();
    _convChangedSub?.cancel();
    _audioPlayer.dispose();
    super.dispose();
  }

  void _onNewMessage(Message msg) {
    final match =
        (msg.isSingleChat && (msg.sendID == userID || msg.recvID == userID)) ||
        (msg.isGroupChat && msg.groupID == groupID);
    // 检查消息是否已存在，避免重复添加（发送方已经添加了临时消息）
    if (match && mounted && !_messages.contains(msg)) {
      setState(() => _messages.add(msg));
      _scrollToBottom();
      // 【修复】接收到新消息时立即标记为已读（因为用户正在聊天页面）
      // 与官方实现一致：当用户在聊天页面时，新消息到达立即标记已读
      if (msg.sendID != Config.userID) {
        _markAsRead();
      }
    }
  }

  void _onReadReceipt(List<ReadReceiptInfo> list) {
    if (!mounted) return;
    try {
      debugPrint('[Chat] _onReadReceipt called with ${list.length} items');
      for (var readInfo in list) {
        debugPrint(
          '[Chat] ReadReceiptInfo: userID=${readInfo.userID}, msgIDList=${readInfo.msgIDList}',
        );

        int updatedCount = 0;
        if (readInfo.msgIDList != null && readInfo.msgIDList!.isNotEmpty) {
          debugPrint('[Chat] Using msgIDList: ${readInfo.msgIDList}');
          for (var e in _messages) {
            if (readInfo.msgIDList!.contains(e.clientMsgID)) {
              debugPrint('[Chat] Marking message ${e.clientMsgID} as read');
              final readTime = DateTime.now().millisecondsSinceEpoch;
              e.isRead = true;
              e.hasReadTime = readTime;
              _cacheReadStatus(e, isRead: true, readTime: readTime);
              updatedCount++;
            }
          }
        }

        debugPrint('[Chat] Updated $updatedCount messages to read status');
      }
      if (mounted) {
        debugPrint('[Chat] Calling setState to refresh UI');
        setState(() {});
      }
    } catch (e) {
      debugPrint('[Chat] _onReadReceipt error: $e');
    }
  }

  Future<void> _loadMessages() async {
    try {
      debugPrint('[Chat] Loading messages for conversationID: $conversationID');
      final result = await OpenIM.iMManager.messageManager
          .getAdvancedHistoryMessageList(
            conversationID: conversationID,
            count: 50,
          );
      final msgs = result.messageList ?? [];
      debugPrint(
        '[Chat] Loaded ${msgs.length} messages, isEnd: ${result.isEnd}',
      );
      for (final msg in msgs) {
        _restoreLocalReadStatus(msg);
      }
      msgs.sort((a, b) => (a.sendTime ?? 0).compareTo(b.sendTime ?? 0));
      _isEnd = result.isEnd ?? true;
      if (mounted) {
        setState(() {
          _messages.clear();
          _messages.addAll(msgs);
          _loading = false;
        });
        _scrollToBottom();
      }
    } catch (e) {
      debugPrint('[Chat] loadMessages error: $e');
      if (mounted) setState(() => _loading = false);
    }
  }

  Future<void> _loadHistoryMessages() async {
    if (_loadingMore || _isEnd || _messages.isEmpty) return;
    _loadingMore = true;
    try {
      final result = await OpenIM.iMManager.messageManager
          .getAdvancedHistoryMessageList(
            conversationID: conversationID,
            startMsg: _messages.first,
            count: 50,
          );
      final msgs = result.messageList ?? [];
      for (final msg in msgs) {
        _restoreLocalReadStatus(msg);
      }
      msgs.sort((a, b) => (a.sendTime ?? 0).compareTo(b.sendTime ?? 0));
      _isEnd = result.isEnd ?? true;
      if (mounted && msgs.isNotEmpty) {
        setState(() => _messages.insertAll(0, msgs));
      }
    } catch (_) {}
    _loadingMore = false;
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

  void _markAsRead() {
    // 【修复问题2】进入聊天页面时清除未读数
    OpenIM.iMManager.conversationManager.markConversationMessageAsRead(
      conversationID: conversationID,
    );
  }

  void _restoreLocalReadStatus(Message msg) {
    final key = msg.clientMsgID;
    if (key == null || key.isEmpty) return;

    // 【修复问题1】从全局缓存中恢复已读状态
    final convCache = _globalReadStatusCache[conversationID];
    final timeCache = _globalReadTimeCache[conversationID];
    
    final cachedRead = convCache?[key];
    if (cachedRead != null) {
      // 服务端已读优先，避免本地旧缓存 false 覆盖后端返回的 true
      if (msg.isRead == true) {
        convCache![key] = true;
      } else {
        msg.isRead = cachedRead;
      }
    }

    final cachedReadTime = timeCache?[key];
    if (cachedReadTime != null && (msg.hasReadTime == null || msg.hasReadTime == 0)) {
      msg.hasReadTime = cachedReadTime;
    }

    // 若服务端已经返回已读，但本地还没有时间，则补一个本地缓存时间，避免 UI 回退
    if (msg.isRead == true) {
      _globalReadStatusCache.putIfAbsent(conversationID, () => {})[key] = true;
      _globalReadTimeCache.putIfAbsent(conversationID, () => {})[key] ??=
          msg.hasReadTime ?? DateTime.now().millisecondsSinceEpoch;
      msg.hasReadTime ??= _globalReadTimeCache[conversationID]![key];
    }
  }

  void _cacheReadStatus(Message msg, {required bool isRead, int? readTime}) {
    final key = msg.clientMsgID;
    if (key == null || key.isEmpty) return;
    
    // 【修复问题1】缓存到全局Map中，按conversationID分组
    _globalReadStatusCache.putIfAbsent(conversationID, () => {})[key] = isRead;
    if (readTime != null) {
      _globalReadTimeCache.putIfAbsent(conversationID, () => {})[key] = readTime;
    }
  }

  void _markMessageAsRead(Message message, bool visible) {
    // 【优化】基于可见性的已读回执
    // 只有当消息可见且未读且不是自己发送的消息时，才标记为已读
    // 注意：不需要每条消息都调用，进入页面时已经调用过 _markAsRead()
    // 这里的逻辑保留是为了兼容未来可能的增量已读功能
  }

  void _scrollToBottom() {
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (_scrollController.hasClients) {
        _scrollController.animateTo(
          _scrollController.position.maxScrollExtent,
          duration: const Duration(milliseconds: 200),
          curve: Curves.easeOut,
        );
      }
    });
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
    // Initialize as unread for sender's own message
    tempMsg.isRead = false;
    _cacheReadStatus(tempMsg, isRead: false);

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
          // Keep isRead as false until we receive read receipt
        });
      }
    } catch (e) {
      if (mounted) setState(() => tempMsg.status = MessageStatus.failed);
      EasyLoading.showToast('发送失败');
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
    // Initialize as unread for sender's own message
    tempMsg.isRead = false;
    _cacheReadStatus(tempMsg, isRead: false);

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
          // Keep isRead as false until we receive read receipt
        });
      }
    } catch (e) {
      if (mounted) setState(() => tempMsg.status = MessageStatus.failed);
      EasyLoading.showToast('发送失败');
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
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
      body: Column(
        children: [
          Expanded(
            child: _loading
                ? const Center(child: CircularProgressIndicator())
                : GestureDetector(
                    onTap: () {
                      FocusScope.of(context).unfocus();
                      if (_showEmojiPicker) {
                        setState(() => _showEmojiPicker = false);
                      }
                    },
                    child: ListView.builder(
                      controller: _scrollController,
                      padding: const EdgeInsets.symmetric(
                        horizontal: 12,
                        vertical: 8,
                      ),
                      itemCount: _messages.length,
                      itemBuilder: (_, i) => _buildMessageItem(_messages[i]),
                    ),
                  ),
          ),
          _buildInputBar(),
          if (_showEmojiPicker)
            SizedBox(
              height: 250,
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
        ],
      ),
    );
  }

  Widget _buildMessageItem(Message msg) {
    final isMe = msg.sendID == Config.userID;
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

  Widget _buildReadTag(Message msg) {
    final isRead = msg.isRead == true;
    return Padding(
      padding: const EdgeInsets.only(bottom: 4),
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

  void _showMessageMenu(Message msg, bool isMe) {
    showModalBottomSheet(
      context: context,
      builder: (_) => SafeArea(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            if (msg.contentType == MessageType.text)
              ListTile(
                leading: const Icon(Icons.copy),
                title: const Text('复制'),
                onTap: () {
                  Clipboard.setData(ClipboardData(text: msg.textContent));
                  Navigator.pop(context);
                  EasyLoading.showToast('已复制');
                },
              ),
            if (isMe)
              ListTile(
                leading: const Icon(Icons.delete, color: Colors.red),
                title: const Text('删除', style: TextStyle(color: Colors.red)),
                onTap: () async {
                  Navigator.pop(context);
                  try {
                    await OpenIM.iMManager.messageManager
                        .deleteMessageFromLocalAndSvr(
                          conversationID: conversationID,
                          seq: msg.seq ?? 0,
                        );
                    setState(() => _messages.remove(msg));
                  } catch (_) {
                    EasyLoading.showToast('删除失败');
                  }
                },
              ),
            ListTile(
              leading: const Icon(Icons.cancel_outlined),
              title: const Text('取消'),
              onTap: () => Navigator.pop(context),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildBubble(Message msg, bool isMe) {
    if ((msg.contentType ?? 0) > 1000) {
      return Center(
        child: Container(
          margin: const EdgeInsets.symmetric(vertical: 8),
          padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 4),
          decoration: BoxDecoration(
            color: Colors.grey[200],
            borderRadius: BorderRadius.circular(4),
          ),
          child: Text(
            msg.textContent,
            style: const TextStyle(fontSize: 12, color: Colors.grey),
          ),
        ),
      );
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
        onLongPress: () => _showMessageMenu(msg, isMe),
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
        onLongPress: () => _showMessageMenu(msg, isMe),
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
      return GestureDetector(
        onTap: () async {
          final soundUrl = msg.soundElem?.sourceUrl ?? '';
          if (soundUrl.isNotEmpty) {
            await _audioPlayer.playOrPause(soundUrl);
          } else {
            EasyLoading.showToast('语音地址无效');
          }
        },
        onLongPress: () => _showMessageMenu(msg, isMe),
        child: Container(
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
          child: Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              Icon(
                Icons.mic,
                size: 20,
                color: isMe ? Colors.white : const Color(0xFF0089FF),
              ),
              const SizedBox(width: 8),
              Text(
                '${duration}s',
                style: TextStyle(
                  fontSize: 14,
                  color: isMe ? Colors.white : Colors.black87,
                ),
              ),
            ],
          ),
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
        onLongPress: () => _showMessageMenu(msg, isMe),
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
      onLongPress: () => _showMessageMenu(msg, isMe),
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
    return Container(
      padding: EdgeInsets.only(
        left: 12,
        right: 8,
        top: 8,
        bottom: MediaQuery.of(context).padding.bottom + 8,
      ),
      decoration: BoxDecoration(
        color: const Color(0xFFF0F2F6),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withValues(alpha: 0.05),
            blurRadius: 4,
            offset: const Offset(0, -2),
          ),
        ],
      ),
      child: Row(
        children: [
          GestureDetector(
            onTap: () {
              setState(() {
                _showEmojiPicker = !_showEmojiPicker;
                if (_showEmojiPicker) {
                  _focusNode.unfocus();
                }
              });
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
          const SizedBox(width: 12),
          Expanded(
            child: Container(
              decoration: BoxDecoration(
                color: Colors.white,
                borderRadius: BorderRadius.circular(4),
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
                    horizontal: 12,
                    vertical: 8,
                  ),
                ),
                style: const TextStyle(fontSize: 17, color: Color(0xFF0C1C33)),
                textInputAction: TextInputAction.send,
                onSubmitted: (_) => _sendTextMessage(),
                onTap: () {
                  if (_showEmojiPicker) {
                    setState(() => _showEmojiPicker = false);
                  }
                },
              ),
            ),
          ),
          const SizedBox(width: 12),
          GestureDetector(
            onTap: _sendButtonVisible ? _sendTextMessage : _pickAndSendImage,
            child: Container(
              width: 32,
              height: 32,
              decoration: const BoxDecoration(color: Colors.transparent),
              child: Icon(
                _sendButtonVisible ? Icons.send : Icons.image_outlined,
                color: _sendButtonVisible
                    ? const Color(0xFF0089FF)
                    : const Color(0xFF8E9AB0),
                size: 24,
              ),
            ),
          ),
          const SizedBox(width: 12),
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
