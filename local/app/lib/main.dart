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
  runZonedGuarded(
    () async {
      WidgetsFlutterBinding.ensureInitialized();
      await Hive.initFlutter();
      await Config.init();
      runApp(const FDIMApp());
    },
    (error, stackTrace) {
      debugPrint('Uncaught error: $error\n$stackTrace');
    },
  );
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
          platform: TargetPlatform.iOS,
          colorScheme: ColorScheme.fromSeed(seedColor: const Color(0xFF1DC063)),
          scaffoldBackgroundColor: const Color(0xFFEDEDED),
          primaryColor: const Color(0xFF1DC063),
          canvasColor: const Color(0xFFFFFFFF),
          appBarTheme: const AppBarTheme(
            backgroundColor: Color(0xFFF2F2F2),
            elevation: 0.5,
            iconTheme: IconThemeData(color: Color(0xFF191919)),
            titleTextStyle: TextStyle(
              color: Color(0xFF191919),
              fontSize: 17,
              fontWeight: FontWeight.w600,
            ),
          ),
          dividerColor: const Color(0xFFD9D9D9),
          textTheme: const TextTheme(
            bodyLarge: TextStyle(color: Color(0xFF191919)),
            bodyMedium: TextStyle(color: Color(0xFF191919)),
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
