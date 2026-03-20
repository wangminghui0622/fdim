import 'dart:async';

import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:intl/intl.dart';

import '../../sdk/flutter_openim_sdk.dart';
import '../../core/config.dart';
import '../../core/controllers/im_controller.dart';
import '../../routes/app_routes.dart';

class ConversationPage extends StatefulWidget {
  const ConversationPage({super.key});

  @override
  State<ConversationPage> createState() => _ConversationPageState();
}

class _ConversationPageState extends State<ConversationPage> {
  List<ConversationInfo> _conversations = [];
  bool _loading = true;
  Timer? _refreshTimer;
  StreamSubscription? _convChangedSub;
  StreamSubscription? _convAddedSub;
  StreamSubscription? _newMsgSub;

  @override
  void initState() {
    super.initState();
    _loadConversations();
    _refreshTimer = Timer.periodic(
      const Duration(seconds: 10),
      (_) => _loadConversations(silent: true),
    );

    // Listen to SDK events via IMController subjects
    final imCtrl = Get.find<IMController>();
    _convChangedSub = imCtrl.conversationChangedSubject.listen((_) {
      debugPrint(
        '[ConversationPage] 【DEBUG】conversationChangedSubject triggered, reloading conversations',
      );
      _loadConversations(silent: true);
    });
    _convAddedSub = imCtrl.conversationAddedSubject.listen((_) {
      _loadConversations(silent: true);
    });
    _newMsgSub = imCtrl.newMessageSubject.listen((_) {
      _loadConversations(silent: true);
    });
  }

  @override
  void dispose() {
    _refreshTimer?.cancel();
    _convChangedSub?.cancel();
    _convAddedSub?.cancel();
    _newMsgSub?.cancel();
    super.dispose();
  }

  Future<void> _loadConversations({bool silent = false}) async {
    if (!silent && !_loading) setState(() => _loading = true);
    try {
      final list = await OpenIM.iMManager.conversationManager
          .getConversationListSplit(offset: 0, count: 50);
      final sorted = OpenIM.iMManager.conversationManager.simpleSort(list);
      if (mounted) setState(() => _conversations = sorted);
    } catch (_) {}
    if (mounted && _loading) setState(() => _loading = false);
  }

  void _openChat(ConversationInfo conv) {
    // 【修复问题3】从 conversationID 中提取 userID/groupID（如果会话对象中没有）
    String userID = conv.userID ?? '';
    String groupID = conv.groupID ?? '';

    // 如果 userID 为空，尝试从 conversationID 中提取
    if (userID.isEmpty && conv.conversationID.startsWith('si_')) {
      // 单聊格式: si_userID1_userID2
      final parts = conv.conversationID.split('_');
      if (parts.length == 3) {
        // 找出不是当前用户的那个 userID
        final id1 = parts[1];
        final id2 = parts[2];
        userID = (id1 == Config.userID) ? id2 : id1;
      }
    } else if (groupID.isEmpty && conv.conversationID.startsWith('sg_')) {
      // 群聊格式: sg_groupID
      final parts = conv.conversationID.split('_');
      if (parts.length == 2) {
        groupID = parts[1];
      }
    }

    Get.toNamed(
      AppRoutes.chat,
      arguments: {
        'conversationID': conv.conversationID,
        'userID': userID,
        'groupID': groupID,
        'showName': conv.showName,
        'faceURL': conv.faceURL,
        'sessionType': conv.conversationType,
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFFF8F9FA),
      appBar: AppBar(
        backgroundColor: Colors.white,
        elevation: 0.5,
        title: const Text(
          '消息',
          style: TextStyle(
            color: Color(0xFF0C1C33),
            fontSize: 17,
            fontWeight: FontWeight.w600,
          ),
        ),
        actions: [
          IconButton(
            icon: const Icon(Icons.search, color: Color(0xFF0C1C33)),
            onPressed: () => Get.toNamed(AppRoutes.globalSearch),
          ),
          PopupMenuButton<String>(
            icon: const Icon(
              Icons.add_circle_outline,
              color: Color(0xFF0C1C33),
            ),
            onSelected: (v) {
              if (v == 'group') Get.toNamed(AppRoutes.createGroup);
              if (v == 'friend') Get.toNamed(AppRoutes.addFriend);
            },
            itemBuilder: (_) => [
              const PopupMenuItem(value: 'friend', child: Text('添加好友')),
              const PopupMenuItem(value: 'group', child: Text('创建群聊')),
            ],
          ),
        ],
      ),
      body: _loading && _conversations.isEmpty
          ? const Center(child: CircularProgressIndicator())
          : RefreshIndicator(
              onRefresh: () => _loadConversations(),
              child: _conversations.isEmpty
                  ? ListView(
                      children: const [
                        SizedBox(height: 200),
                        Center(
                          child: Text(
                            '暂无消息',
                            style: TextStyle(color: Colors.grey),
                          ),
                        ),
                      ],
                    )
                  : ListView.separated(
                      itemCount: _conversations.length,
                      separatorBuilder: (context, index) =>
                          const Divider(height: 1, indent: 76, endIndent: 16),
                      itemBuilder: (context, i) =>
                          _buildConversationItem(_conversations[i]),
                    ),
            ),
    );
  }

  Widget _buildConversationItem(ConversationInfo conv) {
    return Dismissible(
      key: Key(conv.conversationID),
      direction: DismissDirection.endToStart,
      background: Container(
        color: const Color(0xFFFF3B30),
        alignment: Alignment.centerRight,
        padding: const EdgeInsets.only(right: 20),
        child: const Icon(Icons.delete, color: Colors.white),
      ),
      confirmDismiss: (_) async {
        return await showDialog<bool>(
          context: context,
          builder: (ctx) => AlertDialog(
            title: const Text('确认删除'),
            content: const Text('确定删除该会话？'),
            actions: [
              TextButton(
                onPressed: () => Navigator.pop(ctx, false),
                child: const Text('取消'),
              ),
              TextButton(
                onPressed: () => Navigator.pop(ctx, true),
                child: const Text(
                  '删除',
                  style: TextStyle(color: Color(0xFFFF3B30)),
                ),
              ),
            ],
          ),
        );
      },
      onDismissed: (_) async {
        await OpenIM.iMManager.conversationManager
            .deleteConversationAndDeleteAllMsg(
              conversationID: conv.conversationID,
            );
        _loadConversations();
      },
      child: Container(
        color: Colors.white,
        child: InkWell(
          onTap: () => _openChat(conv),
          child: Container(
            height: 68,
            padding: const EdgeInsets.symmetric(horizontal: 16),
            child: Row(
              children: [
                // 左侧：头像
                CircleAvatar(
                  radius: 24,
                  backgroundColor: const Color(0xFFE0E0E0),
                  backgroundImage: (conv.faceURL ?? '').isNotEmpty
                      ? NetworkImage(conv.faceURL!)
                      : null,
                  child: (conv.faceURL ?? '').isEmpty
                      ? Text(
                          (conv.showName ?? '').isNotEmpty
                              ? conv.showName![0].toUpperCase()
                              : '?',
                          style: const TextStyle(
                            fontSize: 14,
                            color: Colors.white,
                            fontWeight: FontWeight.w500,
                          ),
                        )
                      : null,
                ),
                const SizedBox(width: 12),
                // 中间：用户名和最后一条消息
                Expanded(
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Row(
                        children: [
                          Expanded(
                            child: Text(
                              conv.showName ?? '',
                              maxLines: 1,
                              overflow: TextOverflow.ellipsis,
                              style: const TextStyle(
                                fontSize: 17,
                                fontWeight: FontWeight.normal,
                                color: Color(0xFF0C1C33),
                              ),
                            ),
                          ),
                          if ((conv.latestMsgSendTime ?? 0) > 0)
                            Text(
                              _formatTime(conv.latestMsgSendTime ?? 0),
                              style: const TextStyle(
                                fontSize: 12,
                                color: Color(0xFF8E9AB0),
                              ),
                            ),
                        ],
                      ),
                      const SizedBox(height: 3),
                      Row(
                        children: [
                          Expanded(
                            child: Text(
                              (conv.draftText ?? '').isNotEmpty
                                  ? '[草稿] ${conv.draftText}'
                                  : _parseLatestMsg(conv.latestMsg),
                              maxLines: 1,
                              overflow: TextOverflow.ellipsis,
                              style: TextStyle(
                                fontSize: 14,
                                color: (conv.draftText ?? '').isNotEmpty
                                    ? const Color(0xFF0089FF)
                                    : const Color(0xFF8E9AB0),
                              ),
                            ),
                          ),
                          const SizedBox(width: 4),
                          // 免打扰图标或未读数
                          if (conv.recvMsgOpt == 2)
                            Container(
                              width: 16,
                              height: 16,
                              decoration: BoxDecoration(
                                color: const Color(0xFF8E9AB0),
                                shape: BoxShape.circle,
                              ),
                              child: const Icon(
                                Icons.notifications_off,
                                size: 10,
                                color: Colors.white,
                              ),
                            )
                          else if (conv.unreadCount > 0)
                            Container(
                              constraints: const BoxConstraints(minWidth: 16),
                              height: 16,
                              padding: const EdgeInsets.symmetric(
                                horizontal: 4,
                              ),
                              decoration: BoxDecoration(
                                color: const Color(0xFFFF3B30),
                                borderRadius: BorderRadius.circular(8),
                              ),
                              alignment: Alignment.center,
                              child: Text(
                                conv.unreadCount > 99
                                    ? '99+'
                                    : '${conv.unreadCount}',
                                style: const TextStyle(
                                  color: Colors.white,
                                  fontSize: 10,
                                  fontWeight: FontWeight.w600,
                                ),
                              ),
                            ),
                        ],
                      ),
                    ],
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }

  String _parseLatestMsg(Message? latestMsg) {
    if (latestMsg == null) return '';
    return latestMsg.textContent;
  }

  String _formatTime(int timestamp) {
    if (timestamp <= 0) return '';
    final dt = DateTime.fromMillisecondsSinceEpoch(timestamp);
    final now = DateTime.now();
    final diff = now.difference(dt);

    if (diff.inDays == 0) {
      return DateFormat('HH:mm').format(dt);
    } else if (diff.inDays == 1) {
      return '昨天';
    } else if (diff.inDays < 7) {
      const weekDays = ['一', '二', '三', '四', '五', '六', '日'];
      return '周${weekDays[dt.weekday - 1]}';
    } else {
      return DateFormat('MM/dd').format(dt);
    }
  }
}
