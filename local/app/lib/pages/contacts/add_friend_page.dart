import 'package:flutter/material.dart';
import 'package:flutter_easyloading/flutter_easyloading.dart';

import '../../core/config.dart';
import '../../core/apis/user_api.dart';
import '../../sdk/fdim_sdk.dart';

class AddFriendPage extends StatefulWidget {
  const AddFriendPage({super.key});

  @override
  State<AddFriendPage> createState() => _AddFriendPageState();
}

class _AddFriendPageState extends State<AddFriendPage> {
  final _controller = TextEditingController();
  List<Map<String, dynamic>> _results = [];
  bool _searching = false;
  bool _sending = false;

  Future<void> _search() async {
    final keyword = _controller.text.trim();
    if (keyword.isEmpty) return;
    setState(() => _searching = true);
    try {
      final users = await UserApi.searchUserFullInfo(keyword);
      setState(() => _results = users
          .whereType<Map<String, dynamic>>()
          .where((u) => u['userID'] != Config.userID)
          .toList());
    } catch (_) {}
    if (mounted) setState(() => _searching = false);
  }

  Future<void> _showAddDialog(String toUserID, String nickname) async {
    final msgController = TextEditingController(text: '我是 ${Config.nickname}');
    try {
      final result = await showDialog<String>(
        context: context,
        builder: (ctx) => AlertDialog(
          title: Text('添加 $nickname'),
          content: TextField(
            controller: msgController,
            decoration: const InputDecoration(
              hintText: '输入验证消息',
              border: OutlineInputBorder(),
            ),
            maxLines: 3,
          ),
          actions: [
            TextButton(onPressed: () => Navigator.pop(ctx), child: const Text('取消')),
            ElevatedButton(
              onPressed: () => Navigator.pop(ctx, msgController.text.trim()),
              style: ElevatedButton.styleFrom(
                backgroundColor: const Color(0xFF1B72EC),
                foregroundColor: Colors.white,
              ),
              child: const Text('发送'),
            ),
          ],
        ),
      );
      if (result != null) {
        _addFriend(toUserID, reqMsg: result);
      }
    } finally {
      // Ensure controller is disposed after dialog is completely closed
      await Future.delayed(const Duration(milliseconds: 100));
      msgController.dispose();
    }
  }

  Future<void> _addFriend(String toUserID, {String reqMsg = ''}) async {
    if (_sending) return;
    _sending = true;
    EasyLoading.show();
    try {
      await OpenIM.iMManager.friendshipManager.addFriend(
        userID: toUserID,
        reason: reqMsg,
      );
      EasyLoading.showToast('已发送申请');
    } catch (e) {
      EasyLoading.showToast('发送失败');
    } finally {
      _sending = false;
    }
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('添加好友')),
      body: Column(
        children: [
          Padding(
            padding: const EdgeInsets.all(16),
            child: TextField(
              controller: _controller,
              decoration: InputDecoration(
                hintText: '输入手机号/用户ID搜索',
                prefixIcon: const Icon(Icons.search),
                suffixIcon: IconButton(
                  icon: const Icon(Icons.arrow_forward),
                  onPressed: _search,
                ),
                border: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(12)),
              ),
              onSubmitted: (_) => _search(),
            ),
          ),
          if (_searching)
            const Padding(
              padding: EdgeInsets.all(32),
              child: CircularProgressIndicator(),
            )
          else
            Expanded(
              child: _results.isEmpty
                  ? const Center(
                      child: Text('搜索用户', style: TextStyle(color: Colors.grey)))
                  : ListView.separated(
                      itemCount: _results.length,
                      separatorBuilder: (context, index) => const Divider(height: 1),
                      itemBuilder: (_, i) {
                        final user = _results[i];
                        final uid = user['userID'] ?? '';
                        final name = user['nickname'] ?? uid;
                        final face = user['faceURL'] ?? '';
                        return ListTile(
                          leading: CircleAvatar(
                            backgroundColor: Colors.grey[300],
                            backgroundImage:
                                face.isNotEmpty ? NetworkImage(face) : null,
                            child: face.isEmpty
                                ? Text(name.isNotEmpty ? name[0].toUpperCase() : '?',
                                    style: const TextStyle(color: Colors.white))
                                : null,
                          ),
                          title: Text(name),
                          subtitle: Text('ID: $uid',
                              style: const TextStyle(fontSize: 12, color: Colors.grey)),
                          trailing: ElevatedButton(
                            onPressed: () => _showAddDialog(uid, name),
                            style: ElevatedButton.styleFrom(
                              backgroundColor: const Color(0xFF1B72EC),
                              foregroundColor: Colors.white,
                              minimumSize: const Size(70, 32),
                            ),
                            child: const Text('添加'),
                          ),
                        );
                      },
                    ),
            ),
        ],
      ),
    );
  }
}
