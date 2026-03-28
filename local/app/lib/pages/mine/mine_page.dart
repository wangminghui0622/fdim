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
      backgroundColor: const Color(0xFFEDEDED),
      body: Obx(() {
        final user = authCtrl.userInfo.value;
        return ListView(
          children: [
            // 头像+昵称+ID 卡片（微信风格）
            Container(
              padding: const EdgeInsets.only(top: 10),
              color: const Color(0xFFEDEDED),
              child: Container(
                color: Colors.white,
                padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 16),
                child: GestureDetector(
                  behavior: HitTestBehavior.opaque,
                  onTap: () => Get.toNamed(AppRoutes.myInfo),
                  child: Row(
                    children: [
                      CircleAvatar(
                        radius: 30,
                        backgroundColor: Colors.grey[300],
                        backgroundImage: (user?.faceURL ?? '').isNotEmpty
                            ? NetworkImage(user!.faceURL!)
                            : null,
                        child: (user?.faceURL ?? '').isEmpty
                            ? const Icon(Icons.person, size: 36, color: Colors.white)
                            : null,
                      ),
                      const SizedBox(width: 16),
                      Expanded(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(
                              user?.nickname ?? Config.nickname,
                              style: const TextStyle(fontSize: 17, fontWeight: FontWeight.w500),
                            ),
                            const SizedBox(height: 12),
                            Text(
                              'ID: ${user?.userID ?? Config.userID}',
                              style: const TextStyle(fontSize: 13, color: Colors.grey),
                            ),
                          ],
                        ),
                      ),
                      Column(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          GestureDetector(
                            onTap: () => Get.toNamed(AppRoutes.myQrcode),
                            child: const Icon(Icons.qr_code, color: Colors.grey, size: 20),
                          ),
                          const SizedBox(height: 16),
                          GestureDetector(
                            onTap: () => Get.toNamed(AppRoutes.myInfo),
                            child: const Icon(Icons.arrow_forward_ios, size: 14, color: Colors.grey),
                          ),
                        ],
                      ),
                    ],
                  ),
                ),
              ),
            ),
            // 钱包
            _buildSection(children: [
              _buildCell(Icons.account_balance_wallet, '钱包', iconColor: const Color(0xFFE67E22)),
            ]),
            // 收藏 / 相册 / 卡包 / 表情
            _buildSection(children: [
              _buildCell(Icons.bookmark, '收藏', iconColor: const Color(0xFFE74C3C)),
              _buildCell(Icons.photo, '相册', iconColor: const Color(0xFF3498DB)),
              _buildCell(Icons.credit_card, '卡包', iconColor: const Color(0xFF9B59B6)),
              _buildCell(Icons.tag_faces, '表情', iconColor: const Color(0xFFF1C40F)),
            ]),
            // 设置
            _buildSection(children: [
              _buildCell(Icons.settings, '设置', iconColor: const Color(0xFF1B72EC),
                  onTap: () => Get.toNamed(AppRoutes.settings)),
            ]),
            // 退出登录
            _buildSection(children: [
              _buildCell(Icons.logout, '退出登录', textColor: Colors.red,
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
                        child: const Text('退出',
                            style: TextStyle(color: Colors.red)),
                      ),
                    ],
                  ),
                );
                if (confirm == true) authCtrl.logout();
              }),
            ]),
          ],
        );
      }),
    );
  }

  Widget _buildSection({required List<Widget> children}) {
    return Container(
      padding: const EdgeInsets.only(top: 10),
      color: const Color(0xFFEDEDED),
      child: Column(children: children),
    );
  }

  Widget _buildCell(IconData icon, String title,
      {VoidCallback? onTap, Color? iconColor, Color textColor = const Color(0xFF191919)}) {
    return Container(
      color: Colors.white,
      height: 50,
      child: ListTile(
        leading: Icon(icon, size: 22, color: iconColor ?? textColor),
        title: Text(title, style: TextStyle(fontSize: 16, color: textColor)),
        trailing: const Icon(Icons.arrow_forward_ios, size: 14, color: Colors.grey),
        onTap: onTap,
      ),
    );
  }
}
