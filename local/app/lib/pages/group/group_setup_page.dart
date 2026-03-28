import 'package:flutter/material.dart';
import 'package:flutter_easyloading/flutter_easyloading.dart';
import 'package:get/get.dart';

import '../../core/apis/group_api.dart';
import '../../core/models/group_info.dart';

class GroupSetupPage extends StatefulWidget {
  const GroupSetupPage({super.key});

  @override
  State<GroupSetupPage> createState() => _GroupSetupPageState();
}

class _GroupSetupPageState extends State<GroupSetupPage> {
  final _nameCtrl = TextEditingController();
  final _noticeCtrl = TextEditingController();
  final _introCtrl = TextEditingController();
  GroupInfo? _group;
  bool _loading = true;

  @override
  void initState() {
    super.initState();
    final groupID = Get.arguments?['groupID'] ?? '';
    if (groupID.isNotEmpty) _load(groupID);
  }

  @override
  void dispose() {
    _nameCtrl.dispose();
    _noticeCtrl.dispose();
    _introCtrl.dispose();
    super.dispose();
  }

  Future<void> _load(String groupID) async {
    try {
      final info = await GroupApi.getGroupsInfo(groupID);
      if (mounted && info != null) {
        setState(() {
          _group = info;
          _nameCtrl.text = info.groupName;
          _noticeCtrl.text = info.notification;
          _introCtrl.text = info.introduction;
          _loading = false;
        });
      } else {
        if (mounted) setState(() => _loading = false);
      }
    } catch (_) {
      if (mounted) setState(() => _loading = false);
    }
  }

  Future<void> _save() async {
    if (_group == null) return;
    EasyLoading.show();
    try {
      await GroupApi.setGroupInfo(
        groupID: _group!.groupID,
        groupName: _nameCtrl.text.trim(),
        notification: _noticeCtrl.text.trim(),
        introduction: _introCtrl.text.trim(),
      );
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
        title: const Text('群设置'),
        actions: [
          TextButton(
            onPressed: _save,
            child: const Text('保存', style: TextStyle(color: Color(0xFF0089FF), fontWeight: FontWeight.w600)),
          ),
        ],
      ),
      body: _loading
          ? const Center(child: CircularProgressIndicator())
          : _group == null
              ? const Center(child: Text('群组不存在'))
              : ListView(
                  padding: const EdgeInsets.all(16),
                  children: [
                    TextField(
                      controller: _nameCtrl,
                      decoration: const InputDecoration(
                        labelText: '群名称',
                        border: OutlineInputBorder(),
                      ),
                    ),
                    const SizedBox(height: 16),
                    TextField(
                      controller: _noticeCtrl,
                      maxLines: 3,
                      decoration: const InputDecoration(
                        labelText: '群公告',
                        border: OutlineInputBorder(),
                      ),
                    ),
                    const SizedBox(height: 16),
                    TextField(
                      controller: _introCtrl,
                      maxLines: 2,
                      decoration: const InputDecoration(
                        labelText: '群简介',
                        border: OutlineInputBorder(),
                      ),
                    ),
                  ],
                ),
    );
  }
}
