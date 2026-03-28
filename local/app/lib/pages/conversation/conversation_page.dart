import 'dart:async';
import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:focus_detector_v2/focus_detector_v2.dart';
import 'package:get/get.dart';

import '../../core/apis/group_api.dart';
import '../../models/signaling_info.dart';
import '../../sdk/fdim_sdk.dart';
import '../../core/config.dart';
import '../../core/controllers/im_controller.dart';
import '../../routes/app_routes.dart';
import '../../widgets/group_avatar_widget.dart';

class ConversationPage extends StatefulWidget {
  const ConversationPage({super.key});

  @override
  State<ConversationPage> createState() => _ConversationPageState();
}

class _ConversationPageState extends State<ConversationPage> {
  List<ConversationInfo> _conversations = [];
  bool _loading = true;
  StreamSubscription? _convChangedSub;
  StreamSubscription? _convAddedSub;
  // 群成员头像缓存：groupID -> [{url, name}, ...]
  final Map<String, List<Map<String, String>>> _groupMemberAvatars = {};
  final Set<String> _loadingGroupAvatars = {};

  @override
  void initState() {
    super.initState();
    _loadConversations();

    // 与官方 conversation_logic.dart 一致：事件驱动，不轮询
    final imCtrl = Get.find<IMController>();
    _convChangedSub = imCtrl.conversationChangedSubject.listen((newList) {
      // 与官方 onChanged 一致：用回调中的真实 ConversationInfo 替换列表中的旧对象
      _onConversationChanged(newList);
    });
    _convAddedSub = imCtrl.conversationAddedSubject.listen((newList) {
      _onConversationChanged(newList);
    });
  }

  @override
  void dispose() {
    _convChangedSub?.cancel();
    _convAddedSub?.cancel();
    super.dispose();
  }

  /// 与官方 ConversationLogic.onChanged 一致：
  /// 用新的 ConversationInfo 替换列表中的旧对象，然后重新排序
  void _onConversationChanged(List<ConversationInfo> newList) {
    if (!mounted) return;
    if (newList.isEmpty) {
      // 空列表意味着需要全量刷新
      _loadConversations(silent: true);
      return;
    }
    // 替换已有的，添加新的
    for (var newConv in newList) {
      _conversations.removeWhere((e) => e.conversationID == newConv.conversationID);
    }
    _conversations.insertAll(0, newList);
    final sorted = OpenIM.iMManager.conversationManager.simpleSort(_conversations);
    setState(() => _conversations = sorted);
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
    // 与官方一致：markConversationMessageAsRead 会立即清零本地 unreadCount 并通知
    OpenIM.iMManager.conversationManager.markConversationMessageAsRead(
      conversationID: conv.conversationID,
    );

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

    // 确保 sessionType 正确：从 conversationType 获取，如果为 null 则从 conversationID 推导
    int sessionType = conv.conversationType ?? ConversationType.single;
    if (sessionType == ConversationType.single && groupID.isNotEmpty) {
      sessionType = ConversationType.superGroup;
    }

    Get.toNamed(
      AppRoutes.chat,
      arguments: {
        'conversationID': conv.conversationID,
        'userID': userID,
        'groupID': groupID,
        'showName': conv.showName,
        'faceURL': conv.faceURL,
        'sessionType': sessionType,
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFFEDEDED),
      body: FocusDetector(
        onFocusGained: () => _loadConversations(silent: true),
        child: _loading && _conversations.isEmpty
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
                _buildConversationAvatar(conv),
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
                            SizedBox(
                              width: 168,
                              child: Text(
                                _formatTime(conv.latestMsgSendTime ?? 0),
                                textAlign: TextAlign.left,
                                style: const TextStyle(
                                  fontSize: 12,
                                  color: Color(0xFF8E9AB0),
                                ),
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
                                  : _parseLatestMsg(conv.latestMsg, conv: conv),
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
                          else if (_displayUnreadCount(conv) > 0)
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
                                _displayUnreadCount(conv) > 99
                                    ? '99+'
                                    : '${_displayUnreadCount(conv)}',
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

  Widget _buildConversationAvatar(ConversationInfo conv) {
    final groupID = conv.groupID ?? '';
    final isGroup = conv.isGroupChat || conv.conversationID.startsWith('sg_');

    if (isGroup && groupID.isNotEmpty) {
      // 懒加载群成员头像
      _ensureGroupAvatarsLoaded(groupID);
      final members = _groupMemberAvatars[groupID];
      if (members != null && members.isNotEmpty) {
        return GroupAvatarWidget(
          size: 48,
          faceURLs: members.map((m) => m['url'] ?? '').toList(),
          names: members.map((m) => m['name'] ?? '').toList(),
        );
      }
      // 未加载完成时显示默认群图标
      return Container(
        width: 48, height: 48,
        decoration: const BoxDecoration(
          color: Color(0xFF4CAF50),
          shape: BoxShape.circle,
        ),
        child: const Icon(Icons.group, color: Colors.white, size: 22),
      );
    }

    // 单聊头像
    return CircleAvatar(
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
                fontSize: 14, color: Colors.white, fontWeight: FontWeight.w500,
              ),
            )
          : null,
    );
  }

  void _ensureGroupAvatarsLoaded(String groupID) {
    if (_groupMemberAvatars.containsKey(groupID)) return;
    if (_loadingGroupAvatars.contains(groupID)) return;
    _loadingGroupAvatars.add(groupID);
    GroupApi.getGroupMemberList(groupID, showNumber: 20).then((members) {
      if (!mounted) return;
      // 群主排第一，其余按加入顺序
      members.sort((a, b) {
        if (a.roleLevel == 100) return -1;
        if (b.roleLevel == 100) return 1;
        return 0;
      });
      // 取群主 + 最多3个其他成员
      final display = <Map<String, String>>[];
      for (final m in members) {
        display.add({'url': m.faceURL, 'name': m.nickname});
        if (display.length >= 4) break;
      }
      _groupMemberAvatars[groupID] = display;
      _loadingGroupAvatars.remove(groupID);
      if (mounted) setState(() {});
    }).catchError((_) {
      _loadingGroupAvatars.remove(groupID);
    });
  }

  String _parseLatestMsg(Message? latestMsg, {ConversationInfo? conv}) {
    if (latestMsg == null) return '';
    // 撤回消息：优先使用 textElem，否则从 notificationElem.detail 解析
    if (latestMsg.contentType == MessageType.revokeMessageNotification) {
      if (latestMsg.textElem?.content != null && latestMsg.textElem!.content!.isNotEmpty) {
        return latestMsg.textElem!.content!;
      }
      try {
        final detail = latestMsg.notificationElem?.detail;
        if (detail != null && detail.isNotEmpty) {
          final map = _tryJsonDecode(detail);
          if (map is Map<String, dynamic>) {
            final revokerID = map['revokerID'] as String?;
            return revokerID == Config.userID ? '你撤回了一条消息' : '对方撤回了一条消息';
          }
        }
      } catch (_) {}
      return latestMsg.sendID == Config.userID ? '你撤回了一条消息' : '对方撤回了一条消息';
    }
    // 通话结果自定义消息：显示友好文本（多种解析方式兜底）
    if (latestMsg.contentType == MessageType.custom) {
      final info = _tryParseCallResult(latestMsg);
      if (info != null) {
        final isMe = latestMsg.sendID == Config.userID;
        final icon = info.isVideo ? '📹' : '📞';
        return '[通话] ${info.getDisplayText(isMe)} $icon';
      }
    }
    // 群聊消息：显示发送者昵称前缀（与官方微信/OpenIM一致）
    String text = latestMsg.textContent;
    final isGroup = conv?.isGroupChat == true ||
        (conv?.conversationID.startsWith('sg_') ?? false);
    if (isGroup && latestMsg.sendID != Config.userID) {
      final senderName = latestMsg.senderNickname ?? '';
      if (senderName.isNotEmpty) {
        text = '$senderName: $text';
      }
    }
    return text;
  }

  CallResultInfo? _tryParseCallResult(Message msg) {
    // 方式1：直接从 customElem.data 解析
    try {
      final customData = msg.customElem?.data;
      if (customData != null && customData.isNotEmpty) {
        final map = jsonDecode(customData);
        if (map is Map && map['customType'] == CustomMessageType.callResult) {
          return CallResultInfo.fromJson(Map<String, dynamic>.from(map['data']));
        }
      }
    } catch (_) {}
    // 方式2：从 toJson 的 customElem 重新解析
    try {
      final json = msg.toJson();
      final elemMap = json['customElem'];
      if (elemMap is Map) {
        final innerData = elemMap['data'];
        if (innerData is String && innerData.isNotEmpty) {
          final map = jsonDecode(innerData);
          if (map is Map && map['customType'] == CustomMessageType.callResult) {
            return CallResultInfo.fromJson(Map<String, dynamic>.from(map['data']));
          }
        }
      }
    } catch (_) {}
    // 方式3：嵌套解一层
    try {
      final customData = msg.customElem?.data;
      if (customData != null && customData.isNotEmpty) {
        final outer = jsonDecode(customData);
        if (outer is Map && outer['data'] is String) {
          final inner = jsonDecode(outer['data']);
          if (inner is Map && inner['customType'] == CustomMessageType.callResult) {
            return CallResultInfo.fromJson(Map<String, dynamic>.from(inner['data']));
          }
        }
      }
    } catch (_) {}
    return null;
  }

  dynamic _tryJsonDecode(String s) {
    try { return const JsonDecoder().convert(s); } catch (_) { return null; }
  }

  String _formatTime(int timestamp) {
    if (timestamp <= 0) return '';
    final normalized = timestamp < 1000000000000 ? timestamp * 1000 : timestamp;
    // 显式标记为 UTC，再转为手机本地时区
    final dt = DateTime.fromMillisecondsSinceEpoch(normalized, isUtc: true).toLocal();
    final now = DateTime.now();
    final today = DateTime(now.year, now.month, now.day);
    final msgDay = DateTime(dt.year, dt.month, dt.day);
    final diff = today.difference(msgDay).inDays;
    final base = '${dt.year}年${dt.month.toString().padLeft(2, '0')}月${dt.day.toString().padLeft(2, '0')}日 ${dt.hour.toString().padLeft(2, '0')}:${dt.minute.toString().padLeft(2, '0')}';
    if (diff == 0) return '$base(今天)';
    if (diff == 1) return '$base(昨天)';
    return base;
  }

  /// 与官方一致：直接使用 conv.unreadCount（由 LocalStore 实时维护）
  int _displayUnreadCount(ConversationInfo conv) {
    return conv.unreadCount;
  }
}
