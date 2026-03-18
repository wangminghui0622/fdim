import 'package:flutter/material.dart';
import 'package:flutter_easyloading/flutter_easyloading.dart';
import 'package:get/get.dart';

import '../../core/config.dart';
import '../../core/apis/user_api.dart';
import '../../sdk/flutter_openim_sdk.dart';

class EditMyInfoPage extends StatefulWidget {
  const EditMyInfoPage({super.key});

  @override
  State<EditMyInfoPage> createState() => _EditMyInfoPageState();
}

class _EditMyInfoPageState extends State<EditMyInfoPage> {
  final _nicknameCtrl = TextEditingController();
  final _faceUrlCtrl = TextEditingController();
  int _gender = 0;
  final _phoneCtrl = TextEditingController();
  final _emailCtrl = TextEditingController();

  @override
  void initState() {
    super.initState();
    _nicknameCtrl.text = Config.nickname;
    _faceUrlCtrl.text = Config.faceURL;
    _loadDetail();
  }

  @override
  void dispose() {
    _nicknameCtrl.dispose();
    _faceUrlCtrl.dispose();
    _phoneCtrl.dispose();
    _emailCtrl.dispose();
    super.dispose();
  }

  Future<void> _loadDetail() async {
    try {
      final users = await UserApi.getUsersInfo([Config.userID]);
      if (users.isNotEmpty && mounted) {
        final u = users.first;
        setState(() {
          _nicknameCtrl.text = u.nickname;
          _faceUrlCtrl.text = u.faceURL;
          _gender = u.gender;
          _phoneCtrl.text = u.phoneNumber;
          _emailCtrl.text = u.email;
        });
      }
    } catch (_) {}
  }

  Future<void> _save() async {
    final nickname = _nicknameCtrl.text.trim();
    if (nickname.isEmpty) {
      EasyLoading.showToast('昵称不能为空');
      return;
    }
    EasyLoading.show();
    try {
      await OpenIM.iMManager.userManager.setSelfInfo(
        nickname: nickname,
        faceURL: _faceUrlCtrl.text.trim(),
        gender: _gender,
        phoneNumber: _phoneCtrl.text.trim(),
        email: _emailCtrl.text.trim(),
      );
      Config.nickname = nickname;
      Config.faceURL = _faceUrlCtrl.text.trim();
      EasyLoading.showToast('保存成功');
      Get.back(result: true);
    } catch (_) {
      EasyLoading.showToast('保存失败');
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('编辑资料'),
        actions: [
          TextButton(
            onPressed: _save,
            child: const Text('保存', style: TextStyle(color: Colors.white)),
          ),
        ],
      ),
      body: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          TextField(
            controller: _nicknameCtrl,
            decoration: const InputDecoration(
              labelText: '昵称',
              border: OutlineInputBorder(),
            ),
          ),
          const SizedBox(height: 16),
          TextField(
            controller: _faceUrlCtrl,
            decoration: const InputDecoration(
              labelText: '头像URL',
              border: OutlineInputBorder(),
            ),
          ),
          const SizedBox(height: 16),
          DropdownButtonFormField<int>(
            value: _gender,
            decoration: const InputDecoration(
              labelText: '性别',
              border: OutlineInputBorder(),
            ),
            items: const [
              DropdownMenuItem(value: 0, child: Text('未设置')),
              DropdownMenuItem(value: 1, child: Text('男')),
              DropdownMenuItem(value: 2, child: Text('女')),
            ],
            onChanged: (v) => setState(() => _gender = v ?? 0),
          ),
          const SizedBox(height: 16),
          TextField(
            controller: _phoneCtrl,
            keyboardType: TextInputType.phone,
            decoration: const InputDecoration(
              labelText: '手机号',
              border: OutlineInputBorder(),
            ),
          ),
          const SizedBox(height: 16),
          TextField(
            controller: _emailCtrl,
            keyboardType: TextInputType.emailAddress,
            decoration: const InputDecoration(
              labelText: '邮箱',
              border: OutlineInputBorder(),
            ),
          ),
        ],
      ),
    );
  }
}
