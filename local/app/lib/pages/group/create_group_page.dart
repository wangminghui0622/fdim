import 'package:flutter/material.dart';
import 'package:flutter_easyloading/flutter_easyloading.dart';
import 'package:get/get.dart';

import '../../core/apis/group_api.dart';
import '../../core/apis/friend_api.dart';
import '../../core/models/user_info.dart';

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
      await GroupApi.createGroup(
        groupName: name,
        memberUserIDs: _selected.toList(),
      );
      EasyLoading.showToast('创建成功');
      Get.back(result: true);
    } catch (e) {
      EasyLoading.showToast('创建失败');
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('创建群聊'),
        actions: [
          TextButton(
            onPressed: _createGroup,
            child: const Text('创建', style: TextStyle(color: Colors.white)),
          ),
        ],
      ),
      body: Column(
        children: [
          Padding(
            padding: const EdgeInsets.all(16),
            child: TextField(
              controller: _nameCtrl,
              decoration: const InputDecoration(
                labelText: '群名称',
                border: OutlineInputBorder(),
              ),
            ),
          ),
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 16),
            child: Align(
              alignment: Alignment.centerLeft,
              child: Text('选择成员 (${_selected.length})',
                  style: const TextStyle(fontSize: 14, color: Colors.grey)),
            ),
          ),
          const SizedBox(height: 8),
          Expanded(
            child: _loading
                ? const Center(child: CircularProgressIndicator())
                : ListView.builder(
                    itemCount: _friends.length,
                    itemBuilder: (_, i) {
                      final f = _friends[i];
                      final uid = f.friendUserID;
                      final checked = _selected.contains(uid);
                      return CheckboxListTile(
                        value: checked,
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
        ],
      ),
    );
  }
}
