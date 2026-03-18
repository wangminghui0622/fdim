import 'dart:async';

import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:hive_flutter/hive_flutter.dart';
import 'package:flutter_easyloading/flutter_easyloading.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';

import 'core/config.dart';
import 'core/controllers/auth_controller.dart';
import 'core/controllers/im_controller.dart';
import 'routes/app_pages.dart';
import 'routes/app_routes.dart';

void main() {
  runZonedGuarded(() async {
    WidgetsFlutterBinding.ensureInitialized();
    await Hive.initFlutter();
    await Config.init();
    runApp(const FDIMApp());
  }, (error, stackTrace) {
    debugPrint('Uncaught error: $error\n$stackTrace');
  });
}

class FDIMApp extends StatelessWidget {
  const FDIMApp({super.key});

  @override
  Widget build(BuildContext context) {
    return ScreenUtilInit(
      designSize: const Size(375, 812),
      minTextAdapt: true,
      builder: (context, child) => GetMaterialApp(
        debugShowCheckedModeBanner: false,
        title: 'FDIM',
        theme: ThemeData(
          colorScheme: ColorScheme.fromSeed(seedColor: const Color(0xFF1B72EC)),
          scaffoldBackgroundColor: const Color(0xFFF6F6F6),
          appBarTheme: const AppBarTheme(
            backgroundColor: Colors.white,
            elevation: 0.5,
            iconTheme: IconThemeData(color: Colors.black),
            titleTextStyle: TextStyle(
              color: Colors.black,
              fontSize: 18,
              fontWeight: FontWeight.w600,
            ),
          ),
        ),
        initialBinding: InitBinding(),
        getPages: AppPages.routes,
        initialRoute: AppRoutes.splash,
        builder: EasyLoading.init(),
      ),
    );
  }
}

class InitBinding extends Bindings {
  @override
  void dependencies() {
    Get.put<AuthController>(AuthController(), permanent: true);
    Get.put<IMController>(IMController(), permanent: true);
  }
}
