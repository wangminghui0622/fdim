import 'package:flutter/material.dart';
import 'package:flutter_easyloading/flutter_easyloading.dart';
import 'package:get/get.dart';

import '../../core/controllers/auth_controller.dart';

class MyInfoPage extends StatefulWidget {
  const MyInfoPage({super.key});

  @override
  State<MyInfoPage> createState() => _MyInfoPageState();
}

class _MyInfoPageState extends State<MyInfoPage> {
  final _nicknameController = TextEditingController();
  bool _editing = false;

  @override
  void initState() {
    super.initState();
    final auth = Get.find<AuthController>();
    _nicknameController.text = auth.userInfo.value?.nickname ?? '';
  }

  @override
  void dispose() {
    _nicknameController.dispose();
    super.dispose();
  }

  Future<void> _save() async {
    final nickname = _nicknameController.text.trim();
    if (nickname.isEmpty) {
      EasyLoading.showToast('昵称不能为空');
      return;
    }
    EasyLoading.show(status: '保存中...');
    try {
      await Get.find<AuthController>().updateUserInfo(nickname: nickname);
      EasyLoading.showToast('保存成功');
      setState(() => _editing = false);
    } catch (e) {
      EasyLoading.showToast('保存失败');
    }
  }

  @override
  Widget build(BuildContext context) {
    final auth = Get.find<AuthController>();
    return Scaffold(
      appBar: AppBar(
        title: const Text('个人信息'),
        actions: [
          if (_editing)
            TextButton(
              onPressed: _save,
              child: const Text('保存',
                  style: TextStyle(color: Color(0xFF1B72EC))),
            ),
        ],
      ),
      body: Obx(() {
        final user = auth.userInfo.value;
        return ListView(
          children: [
            const SizedBox(height: 16),
            // 头像
            Center(
              child: CircleAvatar(
                radius: 48,
                backgroundColor: Colors.grey[300],
                backgroundImage: (user?.faceURL ?? '').isNotEmpty
                    ? NetworkImage(user!.faceURL!)
                    : null,
                child: (user?.faceURL ?? '').isEmpty
                    ? const Icon(Icons.camera_alt, size: 32, color: Colors.white)
                    : null,
              ),
            ),
            const SizedBox(height: 24),
            // 昵称
            Container(
              color: Colors.white,
              child: ListTile(
                title: const Text('昵称'),
                trailing: _editing
                    ? SizedBox(
                        width: 180,
                        child: TextField(
                          controller: _nicknameController,
                          textAlign: TextAlign.right,
                          decoration: const InputDecoration(
                            border: InputBorder.none,
                            hintText: '输入昵称',
                          ),
                        ),
                      )
                    : Row(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          Text(user?.nickname ?? '',
                              style: const TextStyle(color: Colors.grey)),
                          const SizedBox(width: 4),
                          const Icon(Icons.chevron_right, color: Colors.grey),
                        ],
                      ),
                onTap: () {
                  if (!_editing) setState(() => _editing = true);
                },
              ),
            ),
            const Divider(height: 1, indent: 16),
            // ID
            Container(
              color: Colors.white,
              child: ListTile(
                title: const Text('ID'),
                trailing: Text(user?.userID ?? '',
                    style: const TextStyle(color: Colors.grey)),
              ),
            ),
            const Divider(height: 1, indent: 16),
            // 性别
            Container(
              color: Colors.white,
              child: ListTile(
                title: const Text('性别'),
                trailing: Text(
                  user?.gender == 1
                      ? '男'
                      : user?.gender == 2
                          ? '女'
                          : '未设置',
                  style: const TextStyle(color: Colors.grey),
                ),
              ),
            ),
            const Divider(height: 1, indent: 16),
            // 手机号
            Container(
              color: Colors.white,
              child: ListTile(
                title: const Text('手机号'),
                trailing: Text(user?.phoneNumber ?? '',
                    style: const TextStyle(color: Colors.grey)),
              ),
            ),
          ],
        );
      }),
    );
  }
}
