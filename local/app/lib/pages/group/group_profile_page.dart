import 'package:flutter/material.dart';
import 'package:flutter_easyloading/flutter_easyloading.dart';
import 'package:get/get.dart';

import '../../core/config.dart';
import '../../core/apis/group_api.dart';
import '../../core/models/group_info.dart';
import '../../routes/app_routes.dart';

class GroupProfilePage extends StatefulWidget {
  const GroupProfilePage({super.key});

  @override
  State<GroupProfilePage> createState() => _GroupProfilePageState();
}

class _GroupProfilePageState extends State<GroupProfilePage> {
  GroupInfo? _group;
  List<GroupMemberInfo> _members = [];
  bool _loading = true;
  bool _isOwner = false;
  bool _isAdmin = false;

  @override
  void initState() {
    super.initState();
    final groupID = Get.arguments?['groupID'] ?? '';
    if (groupID.isNotEmpty) _load(groupID);
  }

  Future<void> _load(String groupID) async {
    try {
      final info = await GroupApi.getGroupsInfo(groupID);
      final members = await GroupApi.getGroupMemberList(groupID);
      if (mounted) {
        setState(() {
          _group = info;
          _members = members;
          _loading = false;
          _isOwner = info?.ownerUserID == Config.userID;
          _isAdmin = members.any((m) =>
              m.userID == Config.userID && (m.isOwner || m.isAdmin));
        });
      }
    } catch (_) {
      if (mounted) setState(() => _loading = false);
    }
  }

  Future<void> _inviteMembers() async {
    final result = await Get.toNamed(AppRoutes.selectContacts);
    if (result is List && result.isNotEmpty) {
      EasyLoading.show();
      try {
        await GroupApi.inviteUserToGroup(
          _group!.groupID,
          result.cast<String>(),
        );
        EasyLoading.showToast('已邀请');
        _load(_group!.groupID);
      } catch (_) {
        EasyLoading.showToast('邀请失败');
      }
    }
  }

  void _openGroupChat() {
    if (_group == null) return;
    Get.toNamed(AppRoutes.chat, arguments: {
      'conversationID': 'sg_${_group!.groupID}',
      'userID': '',
      'groupID': _group!.groupID,
      'showName': _group!.groupName,
      'faceURL': _group!.faceURL,
      'sessionType': 3,
    });
  }

  Future<void> _quitGroup() async {
    final confirm = await Get.dialog<bool>(AlertDialog(
      title: const Text('退出群聊'),
      content: const Text('确定要退出该群聊吗？'),
      actions: [
        TextButton(onPressed: () => Get.back(result: false), child: const Text('取消')),
        TextButton(onPressed: () => Get.back(result: true), child: const Text('确定')),
      ],
    ));
    if (confirm != true) return;
    EasyLoading.show();
    try {
      await GroupApi.quitGroup(_group!.groupID);
      EasyLoading.showToast('已退出');
      Get.back(result: true);
    } catch (_) {
      EasyLoading.showToast('操作失败');
    }
  }

  Future<void> _dismissGroup() async {
    final confirm = await Get.dialog<bool>(AlertDialog(
      title: const Text('解散群聊'),
      content: const Text('确定要解散该群聊吗？此操作不可恢复。'),
      actions: [
        TextButton(onPressed: () => Get.back(result: false), child: const Text('取消')),
        TextButton(
          onPressed: () => Get.back(result: true),
          child: const Text('确定', style: TextStyle(color: Colors.red)),
        ),
      ],
    ));
    if (confirm != true) return;
    EasyLoading.show();
    try {
      await GroupApi.dismissGroup(_group!.groupID);
      EasyLoading.showToast('已解散');
      Get.back(result: true);
    } catch (_) {
      EasyLoading.showToast('操作失败');
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('群资料'),
        actions: [
          if (_isOwner || _isAdmin)
            IconButton(
              icon: const Icon(Icons.settings),
              onPressed: () async {
                final result = await Get.toNamed(AppRoutes.groupSetup,
                    arguments: {'groupID': _group!.groupID});
                if (result == true) _load(_group!.groupID);
              },
            ),
        ],
      ),
      body: _loading
          ? const Center(child: CircularProgressIndicator())
          : _group == null
              ? const Center(child: Text('群组不存在'))
              : ListView(
                  children: [
                    const SizedBox(height: 24),
                    Center(
                      child: CircleAvatar(
                        radius: 44,
                        backgroundColor: Colors.blue[100],
                        backgroundImage: _group!.faceURL.isNotEmpty
                            ? NetworkImage(_group!.faceURL)
                            : null,
                        child: _group!.faceURL.isEmpty
                            ? Text(_group!.groupName.isNotEmpty
                                ? _group!.groupName[0].toUpperCase()
                                : 'G',
                                style: const TextStyle(fontSize: 32))
                            : null,
                      ),
                    ),
                    const SizedBox(height: 12),
                    Center(
                      child: Text(_group!.groupName,
                          style: const TextStyle(fontSize: 20, fontWeight: FontWeight.bold)),
                    ),
                    Center(
                      child: Text('ID: ${_group!.groupID}',
                          style: const TextStyle(fontSize: 13, color: Colors.grey)),
                    ),
                    const SizedBox(height: 8),
                    Center(
                      child: Text('${_group!.memberCount} 名成员',
                          style: const TextStyle(color: Colors.grey)),
                    ),
                    if (_group!.notification.isNotEmpty) ...[
                      const SizedBox(height: 16),
                      Padding(
                        padding: const EdgeInsets.symmetric(horizontal: 16),
                        child: Card(
                          child: Padding(
                            padding: const EdgeInsets.all(12),
                            child: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                const Text('群公告',
                                    style: TextStyle(fontWeight: FontWeight.bold, fontSize: 14)),
                                const SizedBox(height: 4),
                                Text(_group!.notification,
                                    style: const TextStyle(fontSize: 13)),
                              ],
                            ),
                          ),
                        ),
                      ),
                    ],
                    const SizedBox(height: 16),
                    const Padding(
                      padding: EdgeInsets.symmetric(horizontal: 16),
                      child: Text('群成员', style: TextStyle(fontWeight: FontWeight.bold)),
                    ),
                    const SizedBox(height: 8),
                    ..._members.take(20).map((m) => ListTile(
                          dense: true,
                          leading: CircleAvatar(
                            radius: 18,
                            backgroundColor: Colors.grey[300],
                            backgroundImage: m.faceURL.isNotEmpty
                                ? NetworkImage(m.faceURL)
                                : null,
                            child: m.faceURL.isEmpty
                                ? Text(m.nickname.isNotEmpty
                                    ? m.nickname[0].toUpperCase()
                                    : '?',
                                    style: const TextStyle(fontSize: 12))
                                : null,
                          ),
                          title: Text(m.nickname.isNotEmpty ? m.nickname : m.userID),
                          trailing: m.isOwner
                              ? const Text('群主',
                                  style: TextStyle(color: Colors.orange, fontSize: 12))
                              : m.isAdmin
                                  ? const Text('管理员',
                                      style: TextStyle(color: Colors.blue, fontSize: 12))
                                  : null,
                          onTap: () {
                            if (m.userID != Config.userID) {
                              Get.toNamed(AppRoutes.userProfile,
                                  arguments: {'userID': m.userID});
                            }
                          },
                        )),
                    if (_members.length > 20)
                      Padding(
                        padding: const EdgeInsets.all(16),
                        child: Center(
                          child: Text('还有 ${_members.length - 20} 位成员',
                              style: const TextStyle(color: Colors.grey)),
                        ),
                      ),
                    const SizedBox(height: 24),
                    if (_isOwner || _isAdmin)
                      Padding(
                        padding: const EdgeInsets.symmetric(horizontal: 24),
                        child: SizedBox(
                          width: double.infinity,
                          height: 48,
                          child: OutlinedButton.icon(
                            onPressed: _inviteMembers,
                            icon: const Icon(Icons.person_add),
                            label: const Text('邀请成员', style: TextStyle(fontSize: 16)),
                            style: OutlinedButton.styleFrom(
                              shape: RoundedRectangleBorder(
                                  borderRadius: BorderRadius.circular(12)),
                            ),
                          ),
                        ),
                      ),
                    if (_isOwner || _isAdmin) const SizedBox(height: 12),
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 24),
                      child: SizedBox(
                        width: double.infinity,
                        height: 48,
                        child: ElevatedButton(
                          onPressed: _openGroupChat,
                          style: ElevatedButton.styleFrom(
                            backgroundColor: const Color(0xFF1B72EC),
                            foregroundColor: Colors.white,
                            shape: RoundedRectangleBorder(
                                borderRadius: BorderRadius.circular(12)),
                          ),
                          child: const Text('发送消息', style: TextStyle(fontSize: 16)),
                        ),
                      ),
                    ),
                    const SizedBox(height: 12),
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 24),
                      child: SizedBox(
                        width: double.infinity,
                        height: 48,
                        child: _isOwner
                            ? OutlinedButton(
                                onPressed: _dismissGroup,
                                style: OutlinedButton.styleFrom(
                                  foregroundColor: Colors.red,
                                  side: const BorderSide(color: Colors.red),
                                  shape: RoundedRectangleBorder(
                                      borderRadius: BorderRadius.circular(12)),
                                ),
                                child: const Text('解散群聊', style: TextStyle(fontSize: 16)),
                              )
                            : OutlinedButton(
                                onPressed: _quitGroup,
                                style: OutlinedButton.styleFrom(
                                  foregroundColor: Colors.red,
                                  side: const BorderSide(color: Colors.red),
                                  shape: RoundedRectangleBorder(
                                      borderRadius: BorderRadius.circular(12)),
                                ),
                                child: const Text('退出群聊', style: TextStyle(fontSize: 16)),
                              ),
                      ),
                    ),
                    const SizedBox(height: 32),
                  ],
                ),
    );
  }
}
