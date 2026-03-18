import 'dart:async';

import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_easyloading/flutter_easyloading.dart';
import 'package:get/get.dart';
import 'package:image_picker/image_picker.dart';

import '../../core/config.dart';
import '../../core/controllers/im_controller.dart';
import '../../sdk/flutter_openim_sdk.dart';
import '../../routes/app_routes.dart';

class ChatPage extends StatefulWidget {
  const ChatPage({super.key});

  @override
  State<ChatPage> createState() => _ChatPageState();
}

class _ChatPageState extends State<ChatPage> {
  final _inputController = TextEditingController();
  final _scrollController = ScrollController();
  final _messages = <Message>[];
  bool _loading = true;
  bool _loadingMore = false;
  bool _isEnd = false;
  bool _isOnline = false;
  Timer? _onlineTimer;
  StreamSubscription? _newMsgSub;
  StreamSubscription? _readReceiptSub;
  StreamSubscription? _convChangedSub;

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
          const Duration(seconds: 15), (_) => _checkOnlineStatus());
    }

    _scrollController.addListener(() {
      if (_scrollController.position.pixels <=
          _scrollController.position.minScrollExtent + 50) {
        _loadHistoryMessages();
      }
    });
  }

  @override
  void dispose() {
    _inputController.dispose();
    _scrollController.dispose();
    _onlineTimer?.cancel();
    _newMsgSub?.cancel();
    _readReceiptSub?.cancel();
    _convChangedSub?.cancel();
    super.dispose();
  }

  void _onNewMessage(Message msg) {
    final match = (msg.isSingleChat &&
            (msg.sendID == userID || msg.recvID == userID)) ||
        (msg.isGroupChat && msg.groupID == groupID);
    // 检查消息是否已存在，避免重复添加（发送方已经添加了临时消息）
    if (match && mounted && !_messages.contains(msg)) {
      setState(() => _messages.add(msg));
      _scrollToBottom();
      _markAsRead();
    }
  }

  void _onReadReceipt(List<ReadReceiptInfo> list) {
    if (!mounted) return;
    try {
      for (var readInfo in list) {
        // 【修复问题1】检查是否是当前会话的已读回执
        if (readInfo.userID == userID) {
          // 后端发送的是 hasReadSeq，表示对方已读到哪个seq
          final hasReadSeq = readInfo.hasReadSeq ?? readInfo.readTime ?? 0;
          debugPrint('[Chat] Received read receipt: userID=${readInfo.userID}, hasReadSeq=$hasReadSeq');
          
          // 更新所有 seq <= hasReadSeq 的消息为已读
          for (var e in _messages) {
            if ((e.seq ?? 0) > 0 && (e.seq ?? 0) <= hasReadSeq) {
              e.isRead = true;
              e.hasReadTime = DateTime.now().millisecondsSinceEpoch;
            }
          }
        }
      }
      if (mounted) setState(() {});
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
      debugPrint('[Chat] Loaded ${msgs.length} messages, isEnd: ${result.isEnd}');
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
      final result =
          await OpenIM.iMManager.userManager.getUserStatus([userID]);
      if (result.isNotEmpty && mounted) {
        final old = _isOnline;
        _isOnline = result.first.isOnline;
        if (old != _isOnline) setState(() {});
      }
    } catch (_) {}
  }

  void _markAsRead() {
    OpenIM.iMManager.conversationManager
        .markConversationMessageAsRead(conversationID: conversationID);
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
    final file = await picker.pickImage(source: ImageSource.gallery, imageQuality: 70);
    if (file == null) return;

    final tempMsg = await OpenIM.iMManager.messageManager
        .createImageMessageFromFullPath(imagePath: file.path);
    tempMsg.sessionType = sessionType;
    tempMsg.recvID = sessionType == ConversationType.single ? userID : '';
    tempMsg.groupID = sessionType != ConversationType.single ? groupID : '';
    // Initialize as unread for sender's own message
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

    final tempMsg = await OpenIM.iMManager.messageManager
        .createTextMessage(text: text);
    tempMsg.sessionType = sessionType;
    tempMsg.recvID = sessionType == ConversationType.single ? userID : '';
    tempMsg.groupID = sessionType != ConversationType.single ? groupID : '';
    // Initialize as unread for sender's own message
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
                    color: _isOnline ? Colors.green : Colors.grey),
              ),
          ],
        ),
        actions: [
          IconButton(
            icon: const Icon(Icons.more_horiz),
            onPressed: () => Get.toNamed(AppRoutes.chatSettings, arguments: {
              'conversationID': conversationID,
              'userID': userID,
              'groupID': groupID,
              'sessionType': sessionType,
            }),
          ),
        ],
      ),
      body: Column(
        children: [
          Expanded(
            child: _loading
                ? const Center(child: CircularProgressIndicator())
                : GestureDetector(
                    onTap: () => FocusScope.of(context).unfocus(),
                    child: ListView.builder(
                      controller: _scrollController,
                      padding: const EdgeInsets.symmetric(
                          horizontal: 12, vertical: 8),
                      itemCount: _messages.length,
                      itemBuilder: (_, i) =>
                          _buildMessageItem(_messages[i]),
                    ),
                  ),
          ),
          _buildInputBar(),
        ],
      ),
    );
  }

  Widget _buildMessageItem(Message msg) {
    final isMe = msg.sendID == Config.userID;
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4),
      child: Row(
        mainAxisAlignment:
            isMe ? MainAxisAlignment.end : MainAxisAlignment.start,
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          if (!isMe) _buildAvatar(msg.senderFaceUrl ?? '', msg.senderNickname ?? ''),
          if (!isMe) const SizedBox(width: 8),
          Flexible(
            child: Column(
              crossAxisAlignment: isMe ? CrossAxisAlignment.end : CrossAxisAlignment.start,
              children: [
                Row(
                  mainAxisSize: MainAxisSize.min,
                  crossAxisAlignment: CrossAxisAlignment.end,
                  children: [
                    _buildBubble(msg, isMe),
                    // Show read status tag next to message bubble (only for sender in single chat)
                    if (isMe && msg.status == MessageStatus.succeeded && sessionType == ConversationType.single) ...[
                      const SizedBox(width: 8),
                      _buildReadTag(msg),
                    ],
                  ],
                ),
              ],
            ),
          ),
          if (isMe) const SizedBox(width: 8),
          if (isMe) _buildAvatar(Config.faceURL, Config.nickname),
        ],
      ),
    );
  }

  Widget _buildReadTag(Message msg) {
    final isRead = msg.isRead == true;
    return Padding(
      padding: const EdgeInsets.only(bottom: 4),
      child: Text(
        isRead ? '已读' : '未读',
        style: const TextStyle(
          fontSize: 11,
          color: Color(0xFF999999),
        ),
      ),
    );
  }

  Widget _buildAvatar(String url, String name) {
    return CircleAvatar(
      radius: 18,
      backgroundColor: Colors.grey[300],
      backgroundImage: url.isNotEmpty ? NetworkImage(url) : null,
      child: url.isEmpty
          ? Text(name.isNotEmpty ? name[0].toUpperCase() : '?',
              style: const TextStyle(color: Colors.white, fontSize: 14))
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
                    await OpenIM.iMManager.messageManager.deleteMessageFromLocalAndSvr(
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
          child: Text(msg.textContent,
              style: const TextStyle(fontSize: 12, color: Colors.grey)),
        ),
      );
    }

    // Image message
    if (msg.contentType == MessageType.picture) {
      final url = msg.pictureElem?.sourcePicture?.url ??
          msg.pictureElem?.snapshotPicture?.url ?? '';
      return GestureDetector(
        onLongPress: () => _showMessageMenu(msg, isMe),
        child: Container(
          constraints: BoxConstraints(
              maxWidth: MediaQuery.of(context).size.width * 0.55,
              maxHeight: 200),
          decoration: BoxDecoration(
            borderRadius: BorderRadius.circular(8),
            boxShadow: [
              BoxShadow(
                  color: Colors.black.withValues(alpha: 0.05),
                  blurRadius: 4,
                  offset: const Offset(0, 2))
            ],
          ),
          child: ClipRRect(
            borderRadius: BorderRadius.circular(8),
            child: url.isNotEmpty
                ? Image.network(url, fit: BoxFit.cover,
                    errorBuilder: (context, error, stackTrace) =>
                        const Icon(Icons.broken_image, size: 48))
                : const Icon(Icons.image, size: 48, color: Colors.grey),
          ),
        ),
      );
    }

    return GestureDetector(
      onLongPress: () => _showMessageMenu(msg, isMe),
      child: Container(
        constraints:
            BoxConstraints(maxWidth: MediaQuery.of(context).size.width * 0.65),
        padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
        decoration: BoxDecoration(
          color: isMe ? const Color(0xFF1B72EC) : Colors.white,
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
                offset: const Offset(0, 2))
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
                    fontSize: 15, color: isMe ? Colors.white : Colors.black87),
              ),
            ),
            if (msg.status == MessageStatus.sending) ...[
              const SizedBox(width: 4),
              const SizedBox(
                  width: 12,
                  height: 12,
                  child: CircularProgressIndicator(strokeWidth: 1.5)),
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
        color: Colors.white,
        boxShadow: [
          BoxShadow(
              color: Colors.black.withValues(alpha: 0.05),
              blurRadius: 4,
              offset: const Offset(0, -2))
        ],
      ),
      child: Row(
        children: [
          Expanded(
            child: TextField(
              controller: _inputController,
              maxLines: 4,
              minLines: 1,
              decoration: InputDecoration(
                hintText: '输入消息...',
                border: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(20),
                    borderSide: BorderSide.none),
                filled: true,
                fillColor: const Color(0xFFF5F5F5),
                contentPadding:
                    const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
              ),
              textInputAction: TextInputAction.send,
              onSubmitted: (_) => _sendTextMessage(),
            ),
          ),
          IconButton(
            onPressed: _pickAndSendImage,
            icon: Icon(Icons.image_outlined, color: Colors.grey[600]),
          ),
          IconButton(
            onPressed: _sendTextMessage,
            icon: const Icon(Icons.send, color: Color(0xFF1B72EC)),
          ),
        ],
      ),
    );
  }
}
