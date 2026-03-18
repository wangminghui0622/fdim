import 'package:flutter/material.dart';
import 'package:flutter_easyloading/flutter_easyloading.dart';
import 'package:get/get.dart';

import '../../core/http_client.dart';

class ForgetPasswordPage extends StatefulWidget {
  const ForgetPasswordPage({super.key});

  @override
  State<ForgetPasswordPage> createState() => _ForgetPasswordPageState();
}

class _ForgetPasswordPageState extends State<ForgetPasswordPage> {
  final _phoneCtrl = TextEditingController();
  final _codeCtrl = TextEditingController();
  final _pwdCtrl = TextEditingController();
  final _confirmCtrl = TextEditingController();
  bool _obscure = true;
  int _countdown = 0;

  @override
  void dispose() {
    _phoneCtrl.dispose();
    _codeCtrl.dispose();
    _pwdCtrl.dispose();
    _confirmCtrl.dispose();
    super.dispose();
  }

  Future<void> _sendCode() async {
    final phone = _phoneCtrl.text.trim();
    if (phone.isEmpty) {
      EasyLoading.showToast('请输入手机号');
      return;
    }
    EasyLoading.show();
    try {
      await HttpClient.post('/account/code/send', data: {
        'phoneNumber': phone,
        'usedFor': 2, // 2 = reset password
      });
      EasyLoading.showToast('验证码已发送');
      _startCountdown();
    } catch (_) {
      EasyLoading.showToast('发送失败');
    }
  }

  void _startCountdown() {
    _countdown = 60;
    Future.doWhile(() async {
      await Future.delayed(const Duration(seconds: 1));
      if (!mounted) return false;
      setState(() => _countdown--);
      return _countdown > 0;
    });
  }

  Future<void> _resetPassword() async {
    final phone = _phoneCtrl.text.trim();
    final code = _codeCtrl.text.trim();
    final pwd = _pwdCtrl.text.trim();
    final confirm = _confirmCtrl.text.trim();

    if (phone.isEmpty || code.isEmpty || pwd.isEmpty) {
      EasyLoading.showToast('请填写所有字段');
      return;
    }
    if (pwd != confirm) {
      EasyLoading.showToast('两次密码不一致');
      return;
    }
    if (pwd.length < 6) {
      EasyLoading.showToast('密码至少6位');
      return;
    }

    EasyLoading.show(status: '重置中...');
    try {
      await HttpClient.post('/account/password/reset', data: {
        'phoneNumber': phone,
        'verifyCode': code,
        'password': pwd,
      });
      EasyLoading.showToast('密码已重置，请重新登录');
      Get.back();
    } catch (e) {
      EasyLoading.showToast('重置失败');
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('忘记密码')),
      body: SingleChildScrollView(
        padding: const EdgeInsets.symmetric(horizontal: 32, vertical: 24),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            const SizedBox(height: 16),
            TextField(
              controller: _phoneCtrl,
              keyboardType: TextInputType.phone,
              decoration: InputDecoration(
                labelText: '手机号',
                prefixIcon: const Icon(Icons.phone_android),
                border: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(12)),
              ),
            ),
            const SizedBox(height: 16),
            Row(
              children: [
                Expanded(
                  child: TextField(
                    controller: _codeCtrl,
                    keyboardType: TextInputType.number,
                    decoration: InputDecoration(
                      labelText: '验证码',
                      prefixIcon: const Icon(Icons.sms),
                      border: OutlineInputBorder(
                          borderRadius: BorderRadius.circular(12)),
                    ),
                  ),
                ),
                const SizedBox(width: 12),
                SizedBox(
                  height: 48,
                  child: ElevatedButton(
                    onPressed: _countdown > 0 ? null : _sendCode,
                    child: Text(_countdown > 0 ? '${_countdown}s' : '获取验证码'),
                  ),
                ),
              ],
            ),
            const SizedBox(height: 16),
            TextField(
              controller: _pwdCtrl,
              obscureText: _obscure,
              decoration: InputDecoration(
                labelText: '新密码',
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
              controller: _confirmCtrl,
              obscureText: _obscure,
              decoration: InputDecoration(
                labelText: '确认新密码',
                prefixIcon: const Icon(Icons.lock_outline),
                border: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(12)),
              ),
            ),
            const SizedBox(height: 32),
            SizedBox(
              height: 48,
              child: ElevatedButton(
                onPressed: _resetPassword,
                style: ElevatedButton.styleFrom(
                  backgroundColor: const Color(0xFF1B72EC),
                  foregroundColor: Colors.white,
                  shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(12)),
                ),
                child: const Text('重置密码', style: TextStyle(fontSize: 16)),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
