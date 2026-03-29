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
    // 等待第一帧渲染完成后再执行初始化
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _checkLoginAndInit();
    });
  }

  Future<void> _checkLoginAndInit() async {
    if (Config.isLoggedIn) {
      // Token exists, need to reconnect WebSocket and init SDK
      try {
        final imCtrl = Get.find<IMController>();
        final authCtrl = Get.find<AuthController>();
        
        // Initialize SDK and login via SDK to reconnect WebSocket
        await imCtrl.initFDIM();
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
    // 黑色背景，与原生启动屏无缝衔接
    return const Scaffold(
      backgroundColor: Colors.black,
    );
  }
}
