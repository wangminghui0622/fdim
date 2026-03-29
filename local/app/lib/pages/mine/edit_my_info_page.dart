import 'package:flutter/material.dart';
import 'package:flutter_easyloading/flutter_easyloading.dart';
import 'package:get/get.dart';
import 'package:image_picker/image_picker.dart';

import '../../core/config.dart';
import '../../core/apis/user_api.dart';
import '../../core/upload_service.dart';
import '../../sdk/fdim_sdk.dart';

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

  Future<void> _pickAndUploadAvatar() async {
    final source = await showModalBottomSheet<ImageSource>(
      context: context,
      builder: (ctx) => SafeArea(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            ListTile(
              leading: const Icon(Icons.camera_alt),
              title: const Text('拍照'),
              onTap: () => Navigator.pop(ctx, ImageSource.camera),
            ),
            ListTile(
              leading: const Icon(Icons.photo_library),
              title: const Text('从相册选择'),
              onTap: () => Navigator.pop(ctx, ImageSource.gallery),
            ),
            ListTile(
              leading: const Icon(Icons.close),
              title: const Text('取消'),
              onTap: () => Navigator.pop(ctx),
            ),
          ],
        ),
      ),
    );
    if (source == null) return;

    final picker = ImagePicker();
    final xFile = await picker.pickImage(source: source, imageQuality: 80, maxWidth: 512);
    if (xFile == null) return;

    EasyLoading.show(status: '上传中...');
    try {
      final url = await UploadService.uploadFile(xFile.path, group: 'avatar');
      if (mounted) {
        setState(() => _faceUrlCtrl.text = url);
      }
      EasyLoading.showToast('头像上传成功');
    } catch (e) {
      debugPrint('[Avatar] Upload error: $e');
      EasyLoading.showToast('上传失败');
    }
  }

  Future<void> _save() async {
    final nickname = _nicknameCtrl.text.trim();
    if (nickname.isEmpty) {
      EasyLoading.showToast('昵称不能为空');
      return;
    }
    EasyLoading.show();
    try {
      await FDIM.iMManager.userManager.setSelfInfo(
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
            child: const Text('保存', style: TextStyle(color: Color(0xFF0089FF), fontWeight: FontWeight.w600)),
          ),
        ],
      ),
      body: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          // 头像
          GestureDetector(
            behavior: HitTestBehavior.opaque,
            onTap: () {
              debugPrint('[EditMyInfo] Avatar tapped');
              _pickAndUploadAvatar();
            },
            child: Center(
              child: SizedBox(
                width: 100,
                height: 100,
                child: Stack(
                  clipBehavior: Clip.none,
                  children: [
                    CircleAvatar(
                      radius: 48,
                      backgroundColor: const Color(0xFFE0E0E0),
                      backgroundImage: _faceUrlCtrl.text.isNotEmpty
                          ? NetworkImage(_faceUrlCtrl.text)
                          : null,
                      child: _faceUrlCtrl.text.isEmpty
                          ? const Icon(Icons.person, size: 48, color: Colors.white)
                          : null,
                    ),
                    Positioned(
                      bottom: 0,
                      right: 0,
                      child: Container(
                        padding: const EdgeInsets.all(4),
                        decoration: const BoxDecoration(
                          color: Color(0xFF0089FF),
                          shape: BoxShape.circle,
                        ),
                        child: const Icon(Icons.camera_alt, size: 16, color: Colors.white),
                      ),
                    ),
                  ],
                ),
              ),
            ),
          ),
          const SizedBox(height: 8),
          Center(
            child: Text('点击更换头像', style: TextStyle(fontSize: 12, color: Colors.grey[500])),
          ),
          const SizedBox(height: 24),
          TextField(
            controller: _nicknameCtrl,
            decoration: const InputDecoration(
              labelText: '昵称',
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
