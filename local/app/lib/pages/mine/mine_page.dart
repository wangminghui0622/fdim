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
      backgroundColor: const Color(0xFFF8F9FA),
      body: Obx(() {
        final user = authCtrl.userInfo.value;
        return SingleChildScrollView(
          child: Column(
            children: [
              // 顶部蓝色背景区域
              Stack(
                children: [
                  Container(
                    height: 138,
                    width: double.infinity,
                    color: const Color(0xFF0089FF),
                  ),
                  // 用户信息卡片
                  Container(
                    height: 98,
                    margin: const EdgeInsets.only(left: 16, right: 16, top: 90),
                    padding: const EdgeInsets.symmetric(horizontal: 16),
                    decoration: BoxDecoration(
                      color: Colors.white,
                      borderRadius: BorderRadius.circular(6),
                    ),
                    child: Row(
                      children: [
                        CircleAvatar(
                          radius: 24,
                          backgroundColor: const Color(0xFFE0E0E0),
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
                                    fontSize: 14,
                                    color: Colors.white,
                                    fontWeight: FontWeight.w500,
                                  ),
                                )
                              : null,
                        ),
                        const SizedBox(width: 10),
                        Expanded(
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            mainAxisAlignment: MainAxisAlignment.center,
                            children: [
                              Text(
                                user?.nickname ?? Config.nickname,
                                style: const TextStyle(
                                  fontSize: 17,
                                  color: Color(0xFF0C1C33),
                                  fontWeight: FontWeight.w500,
                                ),
                              ),
                              const SizedBox(height: 4),
                              Text(
                                user?.userID ?? Config.userID,
                                style: const TextStyle(
                                  fontSize: 14,
                                  color: Color(0xFF8E9AB0),
                                ),
                              ),
                            ],
                          ),
                        ),
                      ],
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 10),
              // 菜单项
              _buildMenuItem(
                icon: Icons.person,
                label: '我的信息',
                onTap: () => Get.toNamed(AppRoutes.myInfo),
                isTopRadius: true,
              ),
              _buildMenuItem(
                icon: Icons.settings,
                label: '账号设置',
                onTap: () => Get.toNamed(AppRoutes.settings),
              ),
              _buildMenuItem(
                icon: Icons.info_outline,
                label: '关于我们',
                onTap: () => Get.toNamed(AppRoutes.aboutUs),
              ),
              _buildMenuItem(
                icon: Icons.logout,
                label: '退出登录',
                onTap: () async {
                  final confirm = await showDialog<bool>(
                    context: context,
                    builder: (ctx) => AlertDialog(
                      title: const Text('确认退出'),
                      content: const Text('确定要退出登录吗？'),
                      actions: [
                        TextButton(
                          onPressed: () => Navigator.pop(ctx, false),
                          child: const Text('取消'),
                        ),
                        TextButton(
                          onPressed: () => Navigator.pop(ctx, true),
                          child: const Text(
                            '退出',
                            style: TextStyle(color: Color(0xFFFF3B30)),
                          ),
                        ),
                      ],
                    ),
                  );
                  if (confirm == true) {
                    authCtrl.logout();
                  }
                },
                isBottomRadius: true,
              ),
            ],
          ),
        );
      }),
    );
  }

  Widget _buildMenuItem({
    required IconData icon,
    required String label,
    bool isTopRadius = false,
    bool isBottomRadius = false,
    VoidCallback? onTap,
  }) {
    return Container(
      margin: const EdgeInsets.symmetric(horizontal: 16),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.only(
          topRight: Radius.circular(isTopRadius ? 6 : 0),
          topLeft: Radius.circular(isTopRadius ? 6 : 0),
          bottomLeft: Radius.circular(isBottomRadius ? 6 : 0),
          bottomRight: Radius.circular(isBottomRadius ? 6 : 0),
        ),
      ),
      child: InkWell(
        onTap: onTap,
        child: Container(
          height: 56,
          padding: const EdgeInsets.only(left: 12, right: 16),
          child: Row(
            children: [
              Icon(icon, size: 24, color: const Color(0xFF0C1C33)),
              const SizedBox(width: 11),
              Text(
                label,
                style: const TextStyle(fontSize: 17, color: Color(0xFF0C1C33)),
              ),
              const Spacer(),
              const Icon(
                Icons.chevron_right,
                size: 24,
                color: Color(0xFF8E9AB0),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
