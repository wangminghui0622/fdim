import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../core/config.dart';
import '../../core/controllers/auth_controller.dart';
import '../../routes/app_routes.dart';

class MinePage extends StatelessWidget {
  const MinePage({super.key});

  @override
  Widget build(BuildContext context) {
    final authCtrl = Get.find<AuthController>();
    return Scaffold(
      appBar: AppBar(title: const Text('我的')),
      body: Obx(() {
        final user = authCtrl.userInfo.value;
        return ListView(
          children: [
            // 用户信息卡片
            Container(
              color: Colors.white,
              padding: const EdgeInsets.all(20),
              child: InkWell(
                onTap: () => Get.toNamed(AppRoutes.myInfo),
                child: Row(
                  children: [
                    CircleAvatar(
                      radius: 32,
                      backgroundColor: Colors.grey[300],
                      backgroundImage: (user?.faceURL ?? '').isNotEmpty
                          ? NetworkImage(user!.faceURL!)
                          : null,
                      child: (user?.faceURL ?? '').isEmpty
                          ? Text(
                              (user?.nickname ?? Config.nickname).isNotEmpty
                                  ? (user?.nickname ?? Config.nickname)[0]
                                      .toUpperCase()
                                  : '?',
                              style: const TextStyle(
                                  fontSize: 24, color: Colors.white))
                          : null,
                    ),
                    const SizedBox(width: 16),
                    Expanded(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(
                            user?.nickname ?? Config.nickname,
                            style: const TextStyle(
                                fontSize: 18, fontWeight: FontWeight.bold),
                          ),
                          const SizedBox(height: 4),
                          Text('ID: ${user?.userID ?? Config.userID}',
                              style: const TextStyle(
                                  fontSize: 13, color: Colors.grey)),
                        ],
                      ),
                    ),
                    const Icon(Icons.chevron_right, color: Colors.grey),
                  ],
                ),
              ),
            ),
            const SizedBox(height: 12),
            // 菜单项
            _buildMenuItem(Icons.edit, '编辑资料', () => Get.toNamed(AppRoutes.editMyInfo)),
            _buildMenuItem(Icons.settings, '设置', () => Get.toNamed(AppRoutes.settings)),
            _buildMenuItem(Icons.block, '黑名单', () => Get.toNamed(AppRoutes.blacklist)),
            _buildMenuItem(Icons.language, '语言', () => Get.toNamed(AppRoutes.languageSetup)),
            _buildMenuItem(Icons.info_outline, '关于', () => Get.toNamed(AppRoutes.aboutUs)),
            _buildMenuItem(Icons.dns_outlined, '服务器配置',
                () => Get.toNamed(AppRoutes.serverConfig)),
            const SizedBox(height: 24),
            // 退出登录
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 16),
              child: SizedBox(
                height: 48,
                child: OutlinedButton(
                  onPressed: () async {
                    final confirm = await showDialog<bool>(
                      context: context,
                      builder: (ctx) => AlertDialog(
                        title: const Text('确认退出'),
                        content: const Text('确定要退出登录吗？'),
                        actions: [
                          TextButton(
                              onPressed: () => Navigator.pop(ctx, false),
                              child: const Text('取消')),
                          TextButton(
                              onPressed: () => Navigator.pop(ctx, true),
                              child: const Text('退出',
                                  style: TextStyle(color: Colors.red))),
                        ],
                      ),
                    );
                    if (confirm == true) {
                      authCtrl.logout();
                    }
                  },
                  style: OutlinedButton.styleFrom(
                    foregroundColor: Colors.red,
                    side: const BorderSide(color: Colors.red),
                    shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(12)),
                  ),
                  child:
                      const Text('退出登录', style: TextStyle(fontSize: 16)),
                ),
              ),
            ),
          ],
        );
      }),
    );
  }

  Widget _buildMenuItem(IconData icon, String title, VoidCallback onTap) {
    return Container(
      color: Colors.white,
      child: ListTile(
        leading: Icon(icon, color: Colors.grey[700]),
        title: Text(title),
        trailing: const Icon(Icons.chevron_right, color: Colors.grey),
        onTap: onTap,
      ),
    );
  }
}
