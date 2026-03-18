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
        const Duration(seconds: 10), (_) => _loadConversations(silent: true));

    // Listen to SDK events via IMController subjects
    final imCtrl = Get.find<IMController>();
    _convChangedSub = imCtrl.conversationChangedSubject.listen((_) {
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
    
    Get.toNamed(AppRoutes.chat, arguments: {
      'conversationID': conv.conversationID,
      'userID': userID,
      'groupID': groupID,
      'showName': conv.showName,
      'faceURL': conv.faceURL,
      'sessionType': conv.conversationType,
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('消息'),
        actions: [
          IconButton(
            icon: const Icon(Icons.search),
            onPressed: () => Get.toNamed(AppRoutes.globalSearch),
          ),
          PopupMenuButton<String>(
            icon: const Icon(Icons.add_circle_outline),
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
                          child: Text('暂无消息',
                              style: TextStyle(color: Colors.grey)),
                        ),
                      ],
                    )
                  : ListView.separated(
                      itemCount: _conversations.length,
                      separatorBuilder: (_, __) => const Divider(
                          height: 1, indent: 76, endIndent: 16),
                      itemBuilder: (_, i) =>
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
        color: Colors.red,
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
                  child: const Text('取消')),
              TextButton(
                  onPressed: () => Navigator.pop(ctx, true),
                  child:
                      const Text('删除', style: TextStyle(color: Colors.red))),
            ],
          ),
        );
      },
      onDismissed: (_) async {
        await OpenIM.iMManager.conversationManager
            .deleteConversationAndDeleteAllMsg(conversationID: conv.conversationID);
        _loadConversations();
      },
      child: ListTile(
        leading: CircleAvatar(
          radius: 24,
          backgroundColor: Colors.grey[300],
          backgroundImage:
              (conv.faceURL ?? '').isNotEmpty ? NetworkImage(conv.faceURL!) : null,
          child: (conv.faceURL ?? '').isEmpty
              ? Text(
                  (conv.showName ?? '').isNotEmpty ? conv.showName![0].toUpperCase() : '?',
                  style: const TextStyle(fontSize: 20, color: Colors.white))
              : null,
        ),
        title: Row(
          children: [
            Expanded(
              child: Text(conv.showName ?? '',
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                  style: const TextStyle(fontWeight: FontWeight.w500)),
            ),
            if ((conv.latestMsgSendTime ?? 0) > 0)
              Text(
                _formatTime(conv.latestMsgSendTime ?? 0),
                style: const TextStyle(fontSize: 12, color: Colors.grey),
              ),
          ],
        ),
        subtitle: Row(
          children: [
            if ((conv.draftText ?? '').isNotEmpty) ...[
              const Text('[草稿] ',
                  style: TextStyle(color: Colors.red, fontSize: 13)),
              Expanded(
                child: Text(conv.draftText ?? '',
                    maxLines: 1,
                    overflow: TextOverflow.ellipsis,
                    style: const TextStyle(fontSize: 13, color: Colors.grey)),
              ),
            ] else
              Expanded(
                child: Text(
                  _parseLatestMsg(conv.latestMsg),
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                  style: const TextStyle(fontSize: 13, color: Colors.grey),
                ),
              ),
            if (conv.unreadCount > 0)
              Container(
                padding:
                    const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
                decoration: BoxDecoration(
                  color: Colors.red,
                  borderRadius: BorderRadius.circular(10),
                ),
                child: Text(
                  conv.unreadCount > 99 ? '99+' : '${conv.unreadCount}',
                  style: const TextStyle(color: Colors.white, fontSize: 11),
                ),
              ),
          ],
        ),
        onTap: () => _openChat(conv),
        contentPadding:
            const EdgeInsets.symmetric(horizontal: 16, vertical: 4),
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
