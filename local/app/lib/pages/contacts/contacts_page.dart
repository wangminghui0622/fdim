import 'dart:async';

import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../core/config.dart';
import '../../core/controllers/im_controller.dart';
import '../../sdk/flutter_openim_sdk.dart';
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
        const Duration(seconds: 15), (_) => _refreshOnlineStatus());

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
            final idx = _friends
                .indexWhere((f) => f.friendUserID == item.userID);
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
    
    Get.toNamed(AppRoutes.chat, arguments: {
      'conversationID': conversationID,
      'userID': friend.friendUserID ?? '',
      'groupID': '',
      'showName': friend.showName,
      'faceURL': friend.faceURL ?? '',
      'sessionType': ConversationType.single,
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('联系人'),
        actions: [
          IconButton(
            icon: const Icon(Icons.person_add_outlined),
            onPressed: () => Get.toNamed(AppRoutes.addFriend),
          ),
        ],
      ),
      body: RefreshIndicator(
        onRefresh: () async {
          await _loadFriends();
          await _loadApplyCount();
        },
        child: ListView(
          children: [
            // 搜索栏
            Padding(
              padding: const EdgeInsets.fromLTRB(12, 8, 12, 4),
              child: TextField(
                controller: _searchCtrl,
                decoration: InputDecoration(
                  hintText: '搜索好友',
                  prefixIcon: const Icon(Icons.search, size: 20),
                  border: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(24),
                      borderSide: BorderSide.none),
                  filled: true,
                  fillColor: const Color(0xFFF5F5F5),
                  contentPadding: const EdgeInsets.symmetric(horizontal: 16),
                  isDense: true,
                ),
                onChanged: (v) => setState(() => _searchKeyword = v),
              ),
            ),
            // 好友申请入口
            ListTile(
              leading: CircleAvatar(
                backgroundColor: const Color(0xFFFF9500),
                child: Badge(
                  isLabelVisible: _applyCount > 0,
                  label: Text('$_applyCount'),
                  child: const Icon(Icons.person_add, color: Colors.white),
                ),
              ),
              title: const Text('好友申请'),
              trailing: const Icon(Icons.chevron_right),
              onTap: () async {
                final imCtrl = Get.find<IMController>();
                imCtrl.clearFriendApplyBadge();
                await Get.toNamed(AppRoutes.friendApply);
                _loadApplyCount();
              },
            ),
            const Divider(height: 1),
            // 群组入口
            ListTile(
              leading: const CircleAvatar(
                backgroundColor: Color(0xFF4CAF50),
                child: Icon(Icons.group, color: Colors.white),
              ),
              title: const Text('我的群组'),
              trailing: const Icon(Icons.chevron_right),
              onTap: () => Get.toNamed(AppRoutes.groupList),
            ),
            const Divider(height: 1),
            // 入群申请入口
            ListTile(
              leading: const CircleAvatar(
                backgroundColor: Color(0xFF2196F3),
                child: Icon(Icons.group_add, color: Colors.white),
              ),
              title: const Text('入群申请'),
              trailing: const Icon(Icons.chevron_right),
              onTap: () => Get.toNamed(AppRoutes.groupRequests),
            ),
            const Divider(height: 1),
            // 好友列表
            if (_loading)
              const Padding(
                padding: EdgeInsets.all(32),
                child: Center(child: CircularProgressIndicator()),
              )
            else if (_friends.isEmpty)
              const Padding(
                padding: EdgeInsets.all(48),
                child: Center(
                    child: Text('暂无好友', style: TextStyle(color: Colors.grey))),
              )
            else
              ..._filteredFriends.map((f) => _buildFriendItem(f)),
          ],
        ),
      ),
    );
  }

  Widget _buildFriendItem(FriendInfo friend) {
    return ListTile(
      leading: Stack(
        children: [
          CircleAvatar(
            radius: 22,
            backgroundColor: Colors.grey[300],
            backgroundImage: (friend.faceURL ?? '').isNotEmpty
                ? NetworkImage(friend.faceURL!)
                : null,
            child: (friend.faceURL ?? '').isEmpty
                ? Text(
                    friend.showName.isNotEmpty
                        ? friend.showName[0].toUpperCase()
                        : '?',
                    style: const TextStyle(color: Colors.white, fontSize: 18))
                : null,
          ),
          Positioned(
            right: 0,
            bottom: 0,
            child: Container(
              width: 12,
              height: 12,
              decoration: BoxDecoration(
                color: friend.isOnline ? Colors.green : Colors.grey,
                shape: BoxShape.circle,
                border: Border.all(color: Colors.white, width: 2),
              ),
            ),
          ),
        ],
      ),
      title: Text(friend.showName),
      onTap: () => _openChat(friend),
    );
  }
}
