import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../core/config.dart';
import '../../core/controllers/auth_controller.dart';
import '../../core/controllers/im_controller.dart';
import '../../routes/app_routes.dart';

class SplashPage extends StatefulWidget {
  const SplashPage({super.key});

  @override
  State<SplashPage> createState() => _SplashPageState();
}

class _SplashPageState extends State<SplashPage> {
  @override
  void initState() {
    super.initState();
    _checkLoginAndInit();
  }

  Future<void> _checkLoginAndInit() async {
    await Future.delayed(const Duration(milliseconds: 1500));
    
    if (Config.isLoggedIn) {
      // Token exists, need to reconnect WebSocket and init SDK
      try {
        final imCtrl = Get.find<IMController>();
        final authCtrl = Get.find<AuthController>();
        
        // Initialize SDK and login via SDK to reconnect WebSocket
        await imCtrl.initOpenIM();
        final info = await imCtrl.loginSDK();
        authCtrl.userInfo.value = info;
        
        Get.offAllNamed(AppRoutes.home);
      } catch (e) {
        debugPrint('[Splash] SDK init/login error: $e');
        // If SDK login fails, clear tokens and go to login page
        await Config.clearLogin();
        Get.offAllNamed(AppRoutes.login);
      }
    } else {
      Get.offAllNamed(AppRoutes.login);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFF1B72EC),
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Container(
              width: 80,
              height: 80,
              decoration: BoxDecoration(
                color: Colors.white,
                borderRadius: BorderRadius.circular(20),
              ),
              child: const Icon(
                Icons.chat_bubble_rounded,
                size: 48,
                color: Color(0xFF1B72EC),
              ),
            ),
            const SizedBox(height: 24),
            const Text(
              'FDIM',
              style: TextStyle(
                fontSize: 32,
                fontWeight: FontWeight.bold,
                color: Colors.white,
              ),
            ),
            const SizedBox(height: 8),
            const Text(
              '即时通讯',
              style: TextStyle(fontSize: 16, color: Colors.white70),
            ),
          ],
        ),
      ),
    );
  }
}
