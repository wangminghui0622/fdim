import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../core/config.dart';
import '../../core/apis/user_api.dart';
import '../../core/apis/friend_api.dart';
import '../../core/models/user_info.dart';
import '../../routes/app_routes.dart';

class UserProfilePage extends StatefulWidget {
  const UserProfilePage({super.key});

  @override
  State<UserProfilePage> createState() => _UserProfilePageState();
}

class _UserProfilePageState extends State<UserProfilePage> {
  UserInfo? _user;
  bool _loading = true;
  bool _isFriend = false;

  @override
  void initState() {
    super.initState();
    final userID = Get.arguments?['userID'] ?? '';
    if (userID.isNotEmpty) _loadUser(userID);
  }

  Future<void> _loadUser(String userID) async {
    try {
      final users = await UserApi.getUsersInfo([userID]);
      if (users.isNotEmpty && mounted) {
        setState(() { _user = users.first; _loading = false; });
      }
      // 检查是否已经是好友
      final friends = await FriendApi.getFriendList();
      _isFriend = friends.any((f) => f.friendUserID == userID);
      if (mounted) setState(() {});
    } catch (_) {
      if (mounted) setState(() => _loading = false);
    }
  }

  void _addFriend() {
    Get.toNamed(AppRoutes.sendVerification, arguments: {
      'userID': _user!.userID,
      'nickname': _user!.nickname,
      'faceURL': _user!.faceURL,
    });
  }

  void _openChat() {
    final uid = _user!.userID;
    final myID = Config.userID;
    // Sort userIDs to match backend conversationID generation
    final ids = [myID, uid]..sort();
    Get.toNamed(AppRoutes.chat, arguments: {
      'conversationID': 'si_${ids[0]}_${ids[1]}',
      'userID': uid,
      'groupID': '',
      'showName': _user!.nickname,
      'faceURL': _user!.faceURL,
      'sessionType': 1,
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('用户资料')),
      body: _loading
          ? const Center(child: CircularProgressIndicator())
          : _user == null
              ? const Center(child: Text('用户不存在'))
              : Column(
                  children: [
                    const SizedBox(height: 32),
                    CircleAvatar(
                      radius: 44,
                      backgroundColor: Colors.grey[300],
                      backgroundImage: _user!.faceURL.isNotEmpty
                          ? NetworkImage(_user!.faceURL)
                          : null,
                      child: _user!.faceURL.isEmpty
                          ? Text(
                              _user!.nickname.isNotEmpty
                                  ? _user!.nickname[0].toUpperCase()
                                  : '?',
                              style: const TextStyle(
                                  fontSize: 32, color: Colors.white))
                          : null,
                    ),
                    const SizedBox(height: 16),
                    Text(_user!.nickname,
                        style: const TextStyle(
                            fontSize: 20, fontWeight: FontWeight.bold)),
                    const SizedBox(height: 4),
                    Text('ID: ${_user!.userID}',
                        style:
                            const TextStyle(fontSize: 14, color: Colors.grey)),
                    const SizedBox(height: 8),
                    if (_user!.gender > 0)
                      Text(_user!.gender == 1 ? '男' : '女',
                          style: const TextStyle(color: Colors.grey)),
                    const Spacer(),
                    Padding(
                      padding: const EdgeInsets.all(24),
                      child: _user!.userID == Config.userID
                          ? const SizedBox()
                          : _isFriend
                              ? SizedBox(
                                  width: double.infinity,
                                  height: 48,
                                  child: ElevatedButton(
                                    onPressed: _openChat,
                                    style: ElevatedButton.styleFrom(
                                      backgroundColor:
                                          const Color(0xFF1B72EC),
                                      foregroundColor: Colors.white,
                                      shape: RoundedRectangleBorder(
                                          borderRadius:
                                              BorderRadius.circular(12)),
                                    ),
                                    child: const Text('发消息',
                                        style: TextStyle(fontSize: 16)),
                                  ),
                                )
                              : SizedBox(
                                  width: double.infinity,
                                  height: 48,
                                  child: ElevatedButton(
                                    onPressed: _addFriend,
                                    style: ElevatedButton.styleFrom(
                                      backgroundColor:
                                          const Color(0xFF1B72EC),
                                      foregroundColor: Colors.white,
                                      shape: RoundedRectangleBorder(
                                          borderRadius:
                                              BorderRadius.circular(12)),
                                    ),
                                    child: const Text('添加好友',
                                        style: TextStyle(fontSize: 16)),
                                  ),
                                ),
                    ),
                  ],
                ),
    );
  }
}
