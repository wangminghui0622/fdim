import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../routes/app_routes.dart';

class SettingsPage extends StatelessWidget {
  const SettingsPage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('设置')),
      body: ListView(
        children: [
          const SizedBox(height: 12),
          Container(
            color: Colors.white,
            child: ListTile(
              title: const Text('修改密码'),
              trailing: const Icon(Icons.chevron_right, color: Colors.grey),
              onTap: () => Get.toNamed(AppRoutes.changePassword),
            ),
          ),
          const Divider(height: 1, indent: 16),
          Container(
            color: Colors.white,
            child: const ListTile(
              title: Text('消息通知'),
              trailing: Icon(Icons.chevron_right, color: Colors.grey),
            ),
          ),
          const Divider(height: 1, indent: 16),
          Container(
            color: Colors.white,
            child: ListTile(
              title: const Text('黑名单'),
              trailing: const Icon(Icons.chevron_right, color: Colors.grey),
              onTap: () => Get.toNamed(AppRoutes.blacklist),
            ),
          ),
          const Divider(height: 1, indent: 16),
          Container(
            color: Colors.white,
            child: ListTile(
              title: const Text('语言'),
              trailing: const Icon(Icons.chevron_right, color: Colors.grey),
              onTap: () => Get.toNamed(AppRoutes.languageSetup),
            ),
          ),
          const Divider(height: 1, indent: 16),
          Container(
            color: Colors.white,
            child: ListTile(
              title: const Text('关于'),
              trailing: const Text('v1.0.0', style: TextStyle(color: Colors.grey)),
              onTap: () => Get.toNamed(AppRoutes.aboutUs),
            ),
          ),
        ],
      ),
    );
  }
}
