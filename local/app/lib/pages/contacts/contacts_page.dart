import 'dart:async';

import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../core/config.dart';
import '../../core/controllers/im_controller.dart';
import '../../sdk/fdim_sdk.dart';
import '../../routes/app_routes.dart';

class ContactsPage extends StatefulWidget {
  const ContactsPage({super.key});

  @override
  State<ContactsPage> createState() => _ContactsPageState();
}

class _ContactsPageState extends State<ContactsPage> {
  List<FriendInfo> _friends = [];
  bool _loading = true;
  int _applyCount = 0;
  Timer? _onlineTimer;
  String _searchKeyword = '';
  final _searchCtrl = TextEditingController();
  StreamSubscription? _applyCountSub;
  StreamSubscription? _friendAddedSub;
  StreamSubscription? _friendDeletedSub;

  @override
  void initState() {
    super.initState();
    _loadFriends();
    _loadApplyCount();
    _onlineTimer = Timer.periodic(
      const Duration(seconds: 15),
      (_) => _refreshOnlineStatus(),
    );

    // Official pattern: IMController.loadFriendApplyCount() fetches from server
    // on every friend application event. We just sync the value to local state.
    final imCtrl = Get.find<IMController>();
    _applyCountSub = imCtrl.friendApplyCount.listen((count) {
      if (mounted) setState(() => _applyCount = count);
    });
    _friendAddedSub = imCtrl.friendAddedSubject.listen((_) {
      _loadFriends();
    });
    _friendDeletedSub = imCtrl.friendDeletedSubject.listen((_) {
      _loadFriends();
    });
  }

  @override
  void dispose() {
    _applyCountSub?.cancel();
    _friendAddedSub?.cancel();
    _friendDeletedSub?.cancel();
    _onlineTimer?.cancel();
    _searchCtrl.dispose();
    super.dispose();
  }

  List<FriendInfo> get _filteredFriends {
    if (_searchKeyword.isEmpty) return _friends;
    final kw = _searchKeyword.toLowerCase();
    return _friends.where((f) {
      final name = f.showName.toLowerCase();
      final uid = (f.friendUserID ?? '').toLowerCase();
      return name.contains(kw) || uid.contains(kw);
    }).toList();
  }

  Future<void> _loadFriends() async {
    try {
      final list = await OpenIM.iMManager.friendshipManager.getFriendList();
      if (mounted) {
        setState(() {
          _friends = list;
          _loading = false;
        });
        _refreshOnlineStatus();
      }
    } catch (_) {
      if (mounted) setState(() => _loading = false);
    }
  }

  Future<void> _loadApplyCount() async {
    try {
      final list = await OpenIM.iMManager.friendshipManager
          .getFriendApplicationListAsRecipient();
      final pending = list.where((e) => e.handleResult == 0).length;
      if (mounted) setState(() => _applyCount = pending);
      final imCtrl = Get.find<IMController>();
      imCtrl.friendApplyCount.value = pending;
    } catch (_) {}
  }

  Future<void> _refreshOnlineStatus() async {
    if (_friends.isEmpty) return;
    try {
      final userIDs = _friends
          .map((f) => f.friendUserID ?? '')
          .where((id) => id.isNotEmpty)
          .toList();
      final result = await OpenIM.iMManager.userManager.getUserStatus(userIDs);
      if (mounted) {
        setState(() {
          for (final item in result) {
            final idx = _friends.indexWhere(
              (f) => f.friendUserID == item.userID,
            );
            if (idx > -1) {
              _friends[idx].isOnline = item.isOnline;
            }
          }
        });
      }
    } catch (_) {}
  }

  void _openChat(FriendInfo friend) async {
    // 【修复问题3】从联系人进入聊天前，先确保会话存在
    // Sort userIDs to match backend conversationID generation
    final ids = [Config.userID, friend.friendUserID]..sort();
    final conversationID = 'si_${ids[0]}_${ids[1]}';

    try {
      // 尝试获取会话，如果不存在会返回一个占位符
      await OpenIM.iMManager.conversationManager.getOneConversation(
        sourceID: friend.friendUserID ?? '',
        sessionType: ConversationType.single,
      );
    } catch (e) {
      debugPrint('[Contacts] Failed to get conversation: $e');
    }

    Get.toNamed(
      AppRoutes.chat,
      arguments: {
        'conversationID': conversationID,
        'userID': friend.friendUserID ?? '',
        'groupID': '',
        'showName': friend.showName,
        'faceURL': friend.faceURL ?? '',
        'sessionType': ConversationType.single,
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFFEDEDED),
      body: RefreshIndicator(
        onRefresh: () async {
          await _loadFriends();
          await _loadApplyCount();
        },
        child: ListView(
          children: [
            const SizedBox(height: 10),
            // 好友申请入口
            _buildMenuItem(
              icon: Icons.person_add,
              iconColor: const Color(0xFFFF9500),
              label: '新的好友',
              count: _applyCount,
              onTap: () async {
                final imCtrl = Get.find<IMController>();
                imCtrl.clearFriendApplyBadge();
                await Get.toNamed(AppRoutes.friendApply);
                _loadApplyCount();
              },
            ),
            // 群组入口
            _buildMenuItem(
              icon: Icons.group,
              iconColor: const Color(0xFF4CAF50),
              label: '我的群聊',
              onTap: () => Get.toNamed(AppRoutes.groupList),
            ),
            // 入群申请入口
            _buildMenuItem(
              icon: Icons.group_add,
              iconColor: const Color(0xFF2196F3),
              label: '入群申请',
              onTap: () => Get.toNamed(AppRoutes.groupRequests),
            ),
            const SizedBox(height: 10),
            // 好友列表
            Container(
              color: Colors.white,
              child: _loading
                  ? const Padding(
                      padding: EdgeInsets.all(32),
                      child: Center(child: CircularProgressIndicator()),
                    )
                  : _friends.isEmpty
                  ? ConstrainedBox(
                      constraints: BoxConstraints(
                        minHeight: MediaQuery.of(context).size.height
                            - MediaQuery.of(context).padding.top
                            - kToolbarHeight   // AppBar
                            - kBottomNavigationBarHeight // TabBar
                            - 200,             // 菜单项 + 间距
                      ),
                      child: const Center(
                        child: Text(
                          '暂无好友',
                          style: TextStyle(color: Colors.grey),
                        ),
                      ),
                    )
                  : Column(
                      children: _filteredFriends
                          .asMap()
                          .entries
                          .map(
                            (e) => Column(
                              children: [
                                _buildFriendItem(e.value),
                                if (e.key < _filteredFriends.length - 1)
                                  const Divider(
                                    height: 1,
                                    indent: 76,
                                    endIndent: 16,
                                  ),
                              ],
                            ),
                          )
                          .toList(),
                    ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildMenuItem({
    required IconData icon,
    required Color iconColor,
    required String label,
    int count = 0,
    VoidCallback? onTap,
  }) {
    return Container(
      color: Colors.white,
      child: InkWell(
        onTap: onTap,
        child: Container(
          height: 60,
          padding: const EdgeInsets.symmetric(horizontal: 16),
          child: Row(
            children: [
              Container(
                width: 42,
                height: 42,
                decoration: BoxDecoration(
                  color: iconColor,
                  borderRadius: BorderRadius.circular(6),
                ),
                child: Icon(icon, color: Colors.white, size: 24),
              ),
              const SizedBox(width: 12),
              Text(
                label,
                style: const TextStyle(fontSize: 17, color: Color(0xFF0C1C33)),
              ),
              const Spacer(),
              if (count > 0)
                Container(
                  constraints: const BoxConstraints(minWidth: 16),
                  height: 16,
                  padding: const EdgeInsets.symmetric(horizontal: 4),
                  decoration: BoxDecoration(
                    color: const Color(0xFFFF3B30),
                    borderRadius: BorderRadius.circular(8),
                  ),
                  alignment: Alignment.center,
                  child: Text(
                    count > 99 ? '99+' : '$count',
                    style: const TextStyle(
                      color: Colors.white,
                      fontSize: 10,
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                ),
              const SizedBox(width: 4),
              const Icon(Icons.chevron_right, color: Color(0xFF8E9AB0)),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildFriendItem(FriendInfo friend) {
    return InkWell(
      onTap: () => _openChat(friend),
      child: Container(
        height: 60,
        padding: const EdgeInsets.symmetric(horizontal: 16),
        child: Row(
          children: [
            Stack(
              children: [
                CircleAvatar(
                  radius: 21,
                  backgroundColor: const Color(0xFFE0E0E0),
                  backgroundImage: (friend.faceURL ?? '').isNotEmpty
                      ? NetworkImage(friend.faceURL!)
                      : null,
                  child: (friend.faceURL ?? '').isEmpty
                      ? Text(
                          friend.showName.isNotEmpty
                              ? friend.showName[0].toUpperCase()
                              : '?',
                          style: const TextStyle(
                            color: Colors.white,
                            fontSize: 14,
                            fontWeight: FontWeight.w500,
                          ),
                        )
                      : null,
                ),
                Positioned(
                  right: 0,
                  bottom: 0,
                  child: Container(
                    width: 10,
                    height: 10,
                    decoration: BoxDecoration(
                      color: friend.isOnline
                          ? const Color(0xFF10CC6A)
                          : const Color(0xFF8E9AB0),
                      shape: BoxShape.circle,
                      border: Border.all(color: Colors.white, width: 1.5),
                    ),
                  ),
                ),
              ],
            ),
            const SizedBox(width: 12),
            Expanded(
              child: Text(
                friend.showName,
                style: const TextStyle(fontSize: 17, color: Color(0xFF0C1C33)),
                maxLines: 1,
                overflow: TextOverflow.ellipsis,
              ),
            ),
          ],
        ),
      ),
    );
  }
}
