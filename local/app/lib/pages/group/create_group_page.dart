import 'package:flutter/material.dart';
import 'package:flutter_easyloading/flutter_easyloading.dart';
import 'package:get/get.dart';

import '../../core/apis/group_api.dart';
import '../../core/apis/friend_api.dart';
import '../../core/config.dart';
import '../../core/local_store.dart';
import '../../core/models/user_info.dart';
import '../../routes/app_routes.dart';
import '../../sdk/flutter_openim_sdk.dart' hide FriendInfo;

class CreateGroupPage extends StatefulWidget {
  const CreateGroupPage({super.key});

  @override
  State<CreateGroupPage> createState() => _CreateGroupPageState();
}

class _CreateGroupPageState extends State<CreateGroupPage> {
  final _nameCtrl = TextEditingController();
  List<FriendInfo> _friends = [];
  final Set<String> _selected = {};
  bool _loading = true;

  @override
  void initState() {
    super.initState();
    _loadFriends();
  }

  @override
  void dispose() {
    _nameCtrl.dispose();
    super.dispose();
  }

  Future<void> _loadFriends() async {
    try {
      final list = await FriendApi.getFriendList();
      if (mounted) setState(() { _friends = list; _loading = false; });
    } catch (_) {
      if (mounted) setState(() => _loading = false);
    }
  }

  Future<void> _createGroup() async {
    final name = _nameCtrl.text.trim();
    if (name.isEmpty) {
      EasyLoading.showToast('请输入群名称');
      return;
    }
    if (_selected.isEmpty) {
      EasyLoading.showToast('请选择至少一位成员');
      return;
    }
    EasyLoading.show();
    try {
      final groupInfo = await GroupApi.createGroup(
        groupName: name,
        memberUserIDs: _selected.toList(),
      );
      EasyLoading.showToast('创建成功');
      if (groupInfo != null && groupInfo.groupID.isNotEmpty) {
        // 手动插入本地会话，确保会话列表立刻显示新群
        final convID = 'sg_${groupInfo.groupID}';
        final conv = ConversationInfo(
          conversationID: convID,
          conversationType: ConversationType.superGroup,
          groupID: groupInfo.groupID,
          showName: groupInfo.groupName,
          faceURL: groupInfo.faceURL,
          latestMsgSendTime: DateTime.now().millisecondsSinceEpoch,
        );
        await LocalStore.putConversation(conv);
        OpenIM.iMManager.conversationManager.listener.conversationChanged([conv]);

        // 创建成功后直接进入群聊（与微信/官方一致）
        Get.back(result: true);
        Get.toNamed(AppRoutes.chat, arguments: {
          'conversationID': 'sg_${groupInfo.groupID}',
          'userID': '',
          'groupID': groupInfo.groupID,
          'showName': groupInfo.groupName,
          'faceURL': groupInfo.faceURL,
          'sessionType': 3,
        });
      } else {
        Get.back(result: true);
      }
    } catch (e) {
      debugPrint('[CreateGroup] Error: $e');
      EasyLoading.dismiss();
    }
  }

  @override
  Widget build(BuildContext context) {
    final canCreate = _nameCtrl.text.trim().isNotEmpty && _selected.isNotEmpty;
    return Scaffold(
      appBar: AppBar(
        title: const Text('创建群聊'),
      ),
      body: Column(
        children: [
          Padding(
            padding: const EdgeInsets.fromLTRB(16, 16, 16, 8),
            child: TextField(
              controller: _nameCtrl,
              onChanged: (_) => setState(() {}),
              decoration: InputDecoration(
                hintText: '请输入群名称',
                prefixIcon: const Icon(Icons.group, color: Color(0xFF0089FF)),
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(10),
                ),
                focusedBorder: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(10),
                  borderSide: const BorderSide(color: Color(0xFF0089FF), width: 1.5),
                ),
              ),
            ),
          ),
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 4),
            child: Align(
              alignment: Alignment.centerLeft,
              child: Text('选择成员 (${_selected.length})',
                  style: const TextStyle(fontSize: 14, fontWeight: FontWeight.w500, color: Colors.grey)),
            ),
          ),
          Expanded(
            child: _loading
                ? const Center(child: CircularProgressIndicator())
                : _friends.isEmpty
                    ? const Center(child: Text('暂无好友', style: TextStyle(color: Colors.grey)))
                    : ListView.builder(
                        itemCount: _friends.length,
                        itemBuilder: (_, i) {
                          final f = _friends[i];
                          final uid = f.friendUserID;
                          final checked = _selected.contains(uid);
                          return CheckboxListTile(
                            value: checked,
                            activeColor: const Color(0xFF0089FF),
                            onChanged: (v) {
                              setState(() {
                                if (v == true) {
                                  _selected.add(uid);
                                } else {
                                  _selected.remove(uid);
                                }
                              });
                            },
                            secondary: CircleAvatar(
                              backgroundColor: Colors.grey[300],
                              backgroundImage: f.faceURL.isNotEmpty
                                  ? NetworkImage(f.faceURL)
                                  : null,
                              child: f.faceURL.isEmpty
                                  ? Text(f.showName.isNotEmpty
                                      ? f.showName[0].toUpperCase()
                                      : '?')
                                  : null,
                            ),
                            title: Text(f.showName),
                          );
                        },
                      ),
          ),
          // 底部确认按钮
          SafeArea(
            child: Padding(
              padding: const EdgeInsets.fromLTRB(24, 8, 24, 16),
              child: SizedBox(
                width: double.infinity,
                height: 48,
                child: ElevatedButton(
                  onPressed: canCreate ? _createGroup : null,
                  style: ElevatedButton.styleFrom(
                    backgroundColor: const Color(0xFF0089FF),
                    disabledBackgroundColor: const Color(0xFFB0D4FF),
                    foregroundColor: Colors.white,
                    disabledForegroundColor: Colors.white70,
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(12),
                    ),
                    elevation: 0,
                  ),
                  child: Text(
                    _selected.isEmpty ? '确认创建' : '确认创建 (${_selected.length}人)',
                    style: const TextStyle(fontSize: 16, fontWeight: FontWeight.w600),
                  ),
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }
}
