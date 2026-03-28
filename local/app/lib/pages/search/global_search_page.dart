import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../core/apis/user_api.dart';
import '../../core/apis/friend_api.dart';
import '../../routes/app_routes.dart';

class GlobalSearchPage extends StatefulWidget {
  const GlobalSearchPage({super.key});

  @override
  State<GlobalSearchPage> createState() => _GlobalSearchPageState();
}

class _GlobalSearchPageState extends State<GlobalSearchPage> {
  final _ctrl = TextEditingController();
  List<dynamic> _userResults = [];
  List<dynamic> _friendResults = [];
  bool _searching = false;

  @override
  void dispose() {
    _ctrl.dispose();
    super.dispose();
  }

  Future<void> _search(String keyword) async {
    if (keyword.trim().isEmpty) {
      setState(() { _userResults = []; _friendResults = []; });
      return;
    }
    setState(() => _searching = true);
    try {
      final users = await UserApi.searchUserPublic(keyword.trim());
      final friends = await FriendApi.searchFriend(keyword.trim());
      if (mounted) {
        setState(() {
          _userResults = users;
          _friendResults = friends;
          _searching = false;
        });
      }
    } catch (_) {
      if (mounted) setState(() => _searching = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: TextField(
          controller: _ctrl,
          autofocus: true,
          style: const TextStyle(color: Color(0xFF0C1C33)),
          decoration: const InputDecoration(
            hintText: '搜索用户ID / 昵称',
            hintStyle: TextStyle(color: Color(0xFF8E9AB0)),
            border: InputBorder.none,
          ),
          onSubmitted: _search,
        ),
        actions: [
          IconButton(
            icon: const Icon(Icons.search),
            onPressed: () => _search(_ctrl.text),
          ),
        ],
      ),
      body: _searching
          ? const Center(child: CircularProgressIndicator())
          : ListView(
              children: [
                if (_friendResults.isNotEmpty) ...[
                  const Padding(
                    padding: EdgeInsets.fromLTRB(16, 16, 16, 8),
                    child: Text('好友', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 14)),
                  ),
                  ..._friendResults.map((f) {
                    final map = f is Map<String, dynamic> ? f : <String, dynamic>{};
                    final uid = map['userID'] ?? map['friendUserID'] ?? '';
                    final name = map['nickname'] ?? map['remark'] ?? uid;
                    final face = map['faceURL'] ?? '';
                    return ListTile(
                      leading: CircleAvatar(
                        backgroundColor: Colors.grey[300],
                        backgroundImage: face.isNotEmpty ? NetworkImage(face) : null,
                        child: face.isEmpty
                            ? Text(name.isNotEmpty ? name[0].toUpperCase() : '?')
                            : null,
                      ),
                      title: Text(name),
                      subtitle: Text(uid, style: const TextStyle(fontSize: 12, color: Colors.grey)),
                      onTap: () => Get.toNamed(AppRoutes.userProfile, arguments: {'userID': uid}),
                    );
                  }),
                ],
                if (_userResults.isNotEmpty) ...[
                  const Padding(
                    padding: EdgeInsets.fromLTRB(16, 16, 16, 8),
                    child: Text('用户', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 14)),
                  ),
                  ..._userResults.map((u) {
                    final map = u is Map<String, dynamic> ? u : <String, dynamic>{};
                    final uid = map['userID'] ?? '';
                    final name = map['nickname'] ?? uid;
                    final face = map['faceURL'] ?? '';
                    return ListTile(
                      leading: CircleAvatar(
                        backgroundColor: Colors.grey[300],
                        backgroundImage: face.isNotEmpty ? NetworkImage(face) : null,
                        child: face.isEmpty
                            ? Text(name.isNotEmpty ? name[0].toUpperCase() : '?')
                            : null,
                      ),
                      title: Text(name),
                      subtitle: Text(uid, style: const TextStyle(fontSize: 12, color: Colors.grey)),
                      onTap: () => Get.toNamed(AppRoutes.userProfile, arguments: {'userID': uid}),
                    );
                  }),
                ],
                if (_userResults.isEmpty && _friendResults.isEmpty && _ctrl.text.isNotEmpty && !_searching)
                  const Padding(
                    padding: EdgeInsets.all(48),
                    child: Center(child: Text('无结果', style: TextStyle(color: Colors.grey))),
                  ),
              ],
            ),
    );
  }
}
