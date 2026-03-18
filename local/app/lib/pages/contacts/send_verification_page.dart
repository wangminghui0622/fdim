import 'package:flutter/material.dart';
import 'package:flutter_easyloading/flutter_easyloading.dart';
import 'package:get/get.dart';

import '../../core/apis/friend_api.dart';

class SendVerificationPage extends StatefulWidget {
  const SendVerificationPage({super.key});

  @override
  State<SendVerificationPage> createState() => _SendVerificationPageState();
}

class _SendVerificationPageState extends State<SendVerificationPage> {
  final _msgCtrl = TextEditingController(text: '请求添加你为好友');
  late final String _userID;
  late final String _nickname;
  late final String _faceURL;

  @override
  void initState() {
    super.initState();
    _userID = Get.arguments?['userID'] ?? '';
    _nickname = Get.arguments?['nickname'] ?? '';
    _faceURL = Get.arguments?['faceURL'] ?? '';
  }

  @override
  void dispose() {
    _msgCtrl.dispose();
    super.dispose();
  }

  Future<void> _send() async {
    EasyLoading.show();
    try {
      await FriendApi.addFriend(_userID, reqMsg: _msgCtrl.text.trim());
      EasyLoading.showToast('已发送申请');
      Get.back(result: true);
    } catch (_) {
      EasyLoading.showToast('发送失败');
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('好友验证')),
      body: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          const SizedBox(height: 16),
          Center(
            child: CircleAvatar(
              radius: 36,
              backgroundColor: Colors.grey[300],
              backgroundImage:
                  _faceURL.isNotEmpty ? NetworkImage(_faceURL) : null,
              child: _faceURL.isEmpty
                  ? Text(
                      _nickname.isNotEmpty ? _nickname[0].toUpperCase() : '?',
                      style: const TextStyle(fontSize: 24, color: Colors.white))
                  : null,
            ),
          ),
          const SizedBox(height: 8),
          Center(
            child: Text(_nickname.isNotEmpty ? _nickname : _userID,
                style: const TextStyle(fontSize: 18, fontWeight: FontWeight.bold)),
          ),
          const SizedBox(height: 24),
          TextField(
            controller: _msgCtrl,
            maxLines: 3,
            decoration: const InputDecoration(
              labelText: '验证信息',
              border: OutlineInputBorder(),
            ),
          ),
          const SizedBox(height: 24),
          SizedBox(
            width: double.infinity,
            height: 48,
            child: ElevatedButton(
              onPressed: _send,
              style: ElevatedButton.styleFrom(
                backgroundColor: const Color(0xFF1B72EC),
                foregroundColor: Colors.white,
                shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(12)),
              ),
              child: const Text('发送', style: TextStyle(fontSize: 16)),
            ),
          ),
        ],
      ),
    );
  }
}
