import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../core/apis/group_api.dart';
import '../../core/models/group_info.dart';
import '../../routes/app_routes.dart';

class GroupListPage extends StatefulWidget {
  const GroupListPage({super.key});

  @override
  State<GroupListPage> createState() => _GroupListPageState();
}

class _GroupListPageState extends State<GroupListPage> {
  List<GroupInfo> _groups = [];
  bool _loading = true;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    try {
      final list = await GroupApi.getJoinedGroupList();
      if (mounted) setState(() { _groups = list; _loading = false; });
    } catch (_) {
      if (mounted) setState(() => _loading = false);
    }
  }

  void _openGroupChat(GroupInfo g) {
    Get.toNamed(AppRoutes.chat, arguments: {
      'conversationID': 'sg_${g.groupID}',
      'userID': '',
      'groupID': g.groupID,
      'showName': g.groupName,
      'faceURL': g.faceURL,
      'sessionType': 3,
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('我的群组'),
        actions: [
          IconButton(
            icon: const Icon(Icons.add),
            onPressed: () async {
              final result = await Get.toNamed(AppRoutes.createGroup);
              if (result == true) _load();
            },
          ),
        ],
      ),
      body: RefreshIndicator(
        onRefresh: _load,
        child: _loading
            ? const Center(child: CircularProgressIndicator())
            : _groups.isEmpty
                ? ListView(
                    children: const [
                      SizedBox(height: 120),
                      Center(child: Text('暂无群组', style: TextStyle(color: Colors.grey))),
                    ],
                  )
                : ListView.separated(
                    itemCount: _groups.length,
                    separatorBuilder: (_, __) => const Divider(height: 1),
                    itemBuilder: (_, i) {
                      final g = _groups[i];
                      return ListTile(
                        leading: CircleAvatar(
                          radius: 22,
                          backgroundColor: Colors.blue[100],
                          backgroundImage: g.faceURL.isNotEmpty
                              ? NetworkImage(g.faceURL)
                              : null,
                          child: g.faceURL.isEmpty
                              ? Text(g.groupName.isNotEmpty
                                  ? g.groupName[0].toUpperCase()
                                  : 'G')
                              : null,
                        ),
                        title: Text(g.groupName),
                        subtitle: Text('${g.memberCount} 人',
                            style: const TextStyle(fontSize: 12, color: Colors.grey)),
                        trailing: const Icon(Icons.chevron_right, size: 18),
                        onTap: () => _openGroupChat(g),
                        onLongPress: () {
                          Get.toNamed(AppRoutes.groupProfile, arguments: {'groupID': g.groupID});
                        },
                      );
                    },
                  ),
      ),
    );
  }
}
