import 'package:flutter/material.dart';
import 'package:flutter_easyloading/flutter_easyloading.dart';

import '../../sdk/fdim_sdk.dart';

class BlacklistPage extends StatefulWidget {
  const BlacklistPage({super.key});

  @override
  State<BlacklistPage> createState() => _BlacklistPageState();
}

class _BlacklistPageState extends State<BlacklistPage> {
  List<BlacklistInfo> _list = [];
  bool _loading = true;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    try {
      final list = await FDIM.iMManager.friendshipManager.getBlacklist();
      if (mounted) setState(() { _list = list; _loading = false; });
    } catch (_) {
      if (mounted) setState(() => _loading = false);
    }
  }

  Future<void> _remove(BlacklistInfo item) async {
    final confirm = await showDialog<bool>(
      context: context,
      builder: (_) => AlertDialog(
        title: const Text('移除黑名单'),
        content: Text('确定将 ${item.nickname ?? item.userID} 移除黑名单？'),
        actions: [
          TextButton(onPressed: () => Navigator.pop(context, false), child: const Text('取消')),
          TextButton(onPressed: () => Navigator.pop(context, true), child: const Text('确定')),
        ],
      ),
    );
    if (confirm != true) return;
    EasyLoading.show();
    try {
      await FDIM.iMManager.friendshipManager.removeBlacklist(userID: item.userID ?? '');
      EasyLoading.showToast('已移除');
      _load();
    } catch (_) {
      EasyLoading.showToast('操作失败');
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('黑名单')),
      body: RefreshIndicator(
        onRefresh: _load,
        child: _loading
            ? const Center(child: CircularProgressIndicator())
            : _list.isEmpty
                ? ListView(children: const [
                    SizedBox(height: 120),
                    Center(child: Text('暂无黑名单', style: TextStyle(color: Colors.grey))),
                  ])
                : ListView.separated(
                    itemCount: _list.length,
                    separatorBuilder: (_, __) => const Divider(height: 1),
                    itemBuilder: (_, i) {
                      final item = _list[i];
                      return ListTile(
                        leading: CircleAvatar(
                          backgroundColor: Colors.grey[300],
                          backgroundImage: (item.faceURL ?? '').isNotEmpty
                              ? NetworkImage(item.faceURL!)
                              : null,
                          child: (item.faceURL ?? '').isEmpty
                              ? Text((item.nickname ?? '?').isNotEmpty
                                  ? item.nickname![0].toUpperCase()
                                  : '?')
                              : null,
                        ),
                        title: Text(item.nickname ?? item.userID ?? ''),
                        trailing: TextButton(
                          onPressed: () => _remove(item),
                          child: const Text('移除'),
                        ),
                      );
                    },
                  ),
      ),
    );
  }
}
