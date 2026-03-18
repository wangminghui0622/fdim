import 'package:flutter/material.dart';
import 'package:flutter_easyloading/flutter_easyloading.dart';

import '../../sdk/flutter_openim_sdk.dart';

class FriendApplyPage extends StatefulWidget {
  const FriendApplyPage({super.key});

  @override
  State<FriendApplyPage> createState() => _FriendApplyPageState();
}

class _FriendApplyPageState extends State<FriendApplyPage> {
  List<FriendApplicationInfo> _applies = [];
  bool _loading = true;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    try {
      final list = await OpenIM.iMManager.friendshipManager
          .getFriendApplicationListAsRecipient();
      if (mounted) setState(() { _applies = list; _loading = false; });
    } catch (_) {
      if (mounted) setState(() => _loading = false);
    }
  }

  Future<void> _accept(FriendApplicationInfo apply) async {
    EasyLoading.show();
    try {
      await OpenIM.iMManager.friendshipManager.acceptFriendApplication(
          userID: apply.fromUserID ?? '');
      EasyLoading.showToast('已同意');
      _load();
    } catch (e) {
      EasyLoading.showToast('操作失败');
    }
  }

  Future<void> _reject(FriendApplicationInfo apply) async {
    EasyLoading.show();
    try {
      await OpenIM.iMManager.friendshipManager.refuseFriendApplication(
          userID: apply.fromUserID ?? '');
      EasyLoading.showToast('已拒绝');
      _load();
    } catch (e) {
      EasyLoading.showToast('操作失败');
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('好友申请')),
      body: _loading
          ? const Center(child: CircularProgressIndicator())
          : _applies.isEmpty
              ? const Center(child: Text('暂无申请', style: TextStyle(color: Colors.grey)))
              : RefreshIndicator(
                  onRefresh: _load,
                  child: ListView.separated(
                    itemCount: _applies.length,
                    separatorBuilder: (_, __) => const Divider(height: 1),
                    itemBuilder: (_, i) => _buildItem(_applies[i]),
                  ),
                ),
    );
  }

  Widget _buildItem(FriendApplicationInfo apply) {
    return ListTile(
      leading: CircleAvatar(
        backgroundColor: Colors.grey[300],
        backgroundImage: (apply.fromFaceURL ?? '').isNotEmpty
            ? NetworkImage(apply.fromFaceURL!) : null,
        child: (apply.fromFaceURL ?? '').isEmpty
            ? Text((apply.fromNickname ?? '').isNotEmpty ? apply.fromNickname![0].toUpperCase() : '?',
                style: const TextStyle(color: Colors.white))
            : null,
      ),
      title: Text((apply.fromNickname ?? '').isNotEmpty ? apply.fromNickname! : (apply.fromUserID ?? '')),
      subtitle: Text((apply.reqMsg ?? '').isNotEmpty ? apply.reqMsg! : '请求添加你为好友',
          style: const TextStyle(fontSize: 13, color: Colors.grey)),
      trailing: apply.handleResult == 0
          ? Row(mainAxisSize: MainAxisSize.min, children: [
              TextButton(onPressed: () => _reject(apply), child: const Text('拒绝')),
              const SizedBox(width: 4),
              ElevatedButton(
                onPressed: () => _accept(apply),
                style: ElevatedButton.styleFrom(
                  backgroundColor: const Color(0xFF1B72EC),
                  foregroundColor: Colors.white,
                  minimumSize: const Size(60, 32),
                ),
                child: const Text('同意'),
              ),
            ])
          : Text(apply.handleResult == 1 ? '已同意' : '已拒绝',
              style: const TextStyle(color: Colors.grey, fontSize: 13)),
    );
  }
}
