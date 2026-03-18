import 'package:flutter/material.dart';
import 'package:flutter_easyloading/flutter_easyloading.dart';
import 'package:get/get.dart';

import '../../core/controllers/auth_controller.dart';

class ChangePasswordPage extends StatefulWidget {
  const ChangePasswordPage({super.key});

  @override
  State<ChangePasswordPage> createState() => _ChangePasswordPageState();
}

class _ChangePasswordPageState extends State<ChangePasswordPage> {
  final _currentController = TextEditingController();
  final _newController = TextEditingController();
  final _confirmController = TextEditingController();
  bool _obscure = true;

  @override
  void dispose() {
    _currentController.dispose();
    _newController.dispose();
    _confirmController.dispose();
    super.dispose();
  }

  Future<void> _submit() async {
    final current = _currentController.text.trim();
    final newPwd = _newController.text.trim();
    final confirm = _confirmController.text.trim();

    if (current.isEmpty || newPwd.isEmpty || confirm.isEmpty) {
      EasyLoading.showToast('请填写所有字段');
      return;
    }
    if (newPwd != confirm) {
      EasyLoading.showToast('两次密码不一致');
      return;
    }
    if (newPwd.length < 6) {
      EasyLoading.showToast('新密码至少6位');
      return;
    }

    EasyLoading.show(status: '修改中...');
    try {
      await Get.find<AuthController>().changePassword(
        currentPassword: current,
        newPassword: newPwd,
      );
      EasyLoading.showToast('修改成功');
      Get.back();
    } catch (e) {
      EasyLoading.dismiss();
      EasyLoading.showToast('修改失败: $e');
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('修改密码')),
      body: Padding(
        padding: const EdgeInsets.all(24),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            TextField(
              controller: _currentController,
              obscureText: _obscure,
              decoration: InputDecoration(
                labelText: '当前密码',
                prefixIcon: const Icon(Icons.lock_outline),
                suffixIcon: IconButton(
                  icon: Icon(
                      _obscure ? Icons.visibility_off : Icons.visibility),
                  onPressed: () => setState(() => _obscure = !_obscure),
                ),
                border: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(12)),
              ),
            ),
            const SizedBox(height: 16),
            TextField(
              controller: _newController,
              obscureText: _obscure,
              decoration: InputDecoration(
                labelText: '新密码',
                prefixIcon: const Icon(Icons.lock),
                border: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(12)),
              ),
            ),
            const SizedBox(height: 16),
            TextField(
              controller: _confirmController,
              obscureText: _obscure,
              decoration: InputDecoration(
                labelText: '确认新密码',
                prefixIcon: const Icon(Icons.lock),
                border: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(12)),
              ),
            ),
            const SizedBox(height: 32),
            SizedBox(
              height: 48,
              child: ElevatedButton(
                onPressed: _submit,
                style: ElevatedButton.styleFrom(
                  backgroundColor: const Color(0xFF1B72EC),
                  foregroundColor: Colors.white,
                  shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(12)),
                ),
                child: const Text('确认修改', style: TextStyle(fontSize: 16)),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
