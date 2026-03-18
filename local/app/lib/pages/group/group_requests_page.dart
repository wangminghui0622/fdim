import 'package:flutter/material.dart';
import 'package:flutter_easyloading/flutter_easyloading.dart';

import '../../core/apis/group_api.dart';
import '../../core/models/group_info.dart';

class GroupRequestsPage extends StatefulWidget {
  const GroupRequestsPage({super.key});

  @override
  State<GroupRequestsPage> createState() => _GroupRequestsPageState();
}

class _GroupRequestsPageState extends State<GroupRequestsPage> {
  List<GroupApplicationInfo> _list = [];
  bool _loading = true;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    try {
      final list = await GroupApi.getRecvGroupApplicationList();
      if (mounted) setState(() { _list = list; _loading = false; });
    } catch (_) {
      if (mounted) setState(() => _loading = false);
    }
  }

  Future<void> _handle(GroupApplicationInfo item, int result) async {
    EasyLoading.show();
    try {
      await GroupApi.respondGroupApplication(
        groupID: item.groupID,
        fromUserID: item.userID,
        handleResult: result,
      );
      EasyLoading.showToast(result == 1 ? '已同意' : '已拒绝');
      _load();
    } catch (_) {
      EasyLoading.showToast('操作失败');
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('入群申请')),
      body: RefreshIndicator(
        onRefresh: _load,
        child: _loading
            ? const Center(child: CircularProgressIndicator())
            : _list.isEmpty
                ? ListView(
                    children: const [
                      SizedBox(height: 120),
                      Center(child: Text('暂无申请', style: TextStyle(color: Colors.grey))),
                    ],
                  )
                : ListView.separated(
                    itemCount: _list.length,
                    separatorBuilder: (_, __) => const Divider(height: 1),
                    itemBuilder: (_, i) {
                      final item = _list[i];
                      return ListTile(
                        leading: CircleAvatar(
                          backgroundColor: Colors.grey[300],
                          backgroundImage: item.userFaceURL.isNotEmpty
                              ? NetworkImage(item.userFaceURL)
                              : null,
                          child: item.userFaceURL.isEmpty
                              ? Text(item.nickname.isNotEmpty
                                  ? item.nickname[0].toUpperCase()
                                  : '?')
                              : null,
                        ),
                        title: Text(item.nickname.isNotEmpty ? item.nickname : item.userID),
                        subtitle: Text(
                          '申请加入 ${item.groupName}${item.reqMsg.isNotEmpty ? '\n${item.reqMsg}' : ''}',
                          maxLines: 2,
                        ),
                        trailing: item.isPending
                            ? Row(
                                mainAxisSize: MainAxisSize.min,
                                children: [
                                  TextButton(
                                    onPressed: () => _handle(item, -1),
                                    child: const Text('拒绝'),
                                  ),
                                  const SizedBox(width: 4),
                                  FilledButton(
                                    onPressed: () => _handle(item, 1),
                                    child: const Text('同意'),
                                  ),
                                ],
                              )
                            : Text(
                                item.handleResult == 1 ? '已同意' : '已拒绝',
                                style: const TextStyle(color: Colors.grey, fontSize: 13),
                              ),
                      );
                    },
                  ),
      ),
    );
  }
}
