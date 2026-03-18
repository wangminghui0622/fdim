import 'package:flutter/material.dart';
import 'package:flutter_easyloading/flutter_easyloading.dart';
import 'package:get/get.dart';

import '../../core/controllers/auth_controller.dart';
import '../../core/http_client.dart';
import '../../routes/app_routes.dart';

class LoginPage extends StatefulWidget {
  const LoginPage({super.key});

  @override
  State<LoginPage> createState() => _LoginPageState();
}

class _LoginPageState extends State<LoginPage> {
  final _phoneController = TextEditingController();
  final _passwordController = TextEditingController();
  bool _obscurePassword = true;

  @override
  void dispose() {
    _phoneController.dispose();
    _passwordController.dispose();
    super.dispose();
  }

  Future<void> _login() async {
    final phone = _phoneController.text.trim();
    final password = _passwordController.text.trim();
    if (phone.isEmpty || password.isEmpty) {
      EasyLoading.showToast('请输入手机号和密码');
      return;
    }
    EasyLoading.show(status: '登录中...');
    try {
      await Get.find<AuthController>().login(
        phoneNumber: phone,
        password: password,
      );
      EasyLoading.dismiss();
      Get.offAllNamed(AppRoutes.home);
    } catch (e) {
      EasyLoading.dismiss();
      // Extract error message from ApiException
      String errorMsg = '登录失败';
      if (e is ApiException) {
        errorMsg = e.message.isNotEmpty ? e.message : '登录失败';
      } else {
        errorMsg = e.toString();
      }
      EasyLoading.showToast(errorMsg);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.white,
      body: SafeArea(
        child: SingleChildScrollView(
          padding: const EdgeInsets.symmetric(horizontal: 32),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              const SizedBox(height: 80),
              // Logo
              Center(
                child: Container(
                  width: 72,
                  height: 72,
                  decoration: BoxDecoration(
                    color: const Color(0xFF1B72EC),
                    borderRadius: BorderRadius.circular(18),
                  ),
                  child: const Icon(Icons.chat_bubble_rounded,
                      size: 40, color: Colors.white),
                ),
              ),
              const SizedBox(height: 24),
              const Center(
                child: Text('FDIM',
                    style:
                        TextStyle(fontSize: 28, fontWeight: FontWeight.bold)),
              ),
              const SizedBox(height: 48),
              // 手机号
              TextField(
                controller: _phoneController,
                keyboardType: TextInputType.phone,
                decoration: InputDecoration(
                  labelText: '手机号',
                  prefixIcon: const Icon(Icons.phone_android),
                  border: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(12)),
                ),
              ),
              const SizedBox(height: 16),
              // 密码
              TextField(
                controller: _passwordController,
                obscureText: _obscurePassword,
                decoration: InputDecoration(
                  labelText: '密码',
                  prefixIcon: const Icon(Icons.lock_outline),
                  suffixIcon: IconButton(
                    icon: Icon(_obscurePassword
                        ? Icons.visibility_off
                        : Icons.visibility),
                    onPressed: () =>
                        setState(() => _obscurePassword = !_obscurePassword),
                  ),
                  border: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(12)),
                ),
              ),
              const SizedBox(height: 32),
              // 登录按钮
              SizedBox(
                height: 48,
                child: ElevatedButton(
                  onPressed: _login,
                  style: ElevatedButton.styleFrom(
                    backgroundColor: const Color(0xFF1B72EC),
                    foregroundColor: Colors.white,
                    shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(12)),
                  ),
                  child:
                      const Text('登录', style: TextStyle(fontSize: 16)),
                ),
              ),
              const SizedBox(height: 16),
              // 注册 & 配置
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  TextButton(
                    onPressed: () => Get.toNamed(AppRoutes.register),
                    child: const Text('注册账号'),
                  ),
                  TextButton(
                    onPressed: () => Get.toNamed(AppRoutes.forgetPassword),
                    child: const Text('忘记密码'),
                  ),
                ],
              ),
              Align(
                alignment: Alignment.centerRight,
                child: TextButton(
                  onPressed: () => Get.toNamed(AppRoutes.serverConfig),
                  child: const Text('服务器配置',
                      style: TextStyle(color: Colors.grey)),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
