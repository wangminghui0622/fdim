import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../routes/app_routes.dart';

class DiscoverPage extends StatelessWidget {
  const DiscoverPage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFFEDEDED),
      body: ListView(
        children: [
          // 广场
          _buildCell(
            context,
            icon: Icons.grid_view_rounded,
            iconColor: const Color(0xFFE67E22),
            title: '广场',
            topPadding: 10,
            showDivider: false,
            onTap: () {},
          ),
          // 朋友圈
          _buildCell(
            context,
            icon: Icons.camera_alt,
            iconColor: const Color(0xFF27AE60),
            title: '朋友圈',
            topPadding: 10,
            showDivider: false,
            onTap: () {},
          ),
          // 全部群组
          _buildCell(
            context,
            icon: Icons.groups,
            iconColor: const Color(0xFF16A085),
            title: '全部群组',
            topPadding: 10,
            showDivider: false,
            onTap: () => Get.toNamed(AppRoutes.groupList),
          ),
          // 搜一搜
          _buildCell(
            context,
            icon: Icons.search,
            iconColor: const Color(0xFFE74C3C),
            title: '搜一搜',
            topPadding: 10,
            showDivider: false,
            onTap: () => Get.toNamed(AppRoutes.globalSearch),
          ),
          // 分销推广
          _buildCell(
            context,
            icon: Icons.share,
            iconColor: const Color(0xFF3498DB),
            title: '分销推广',
            topPadding: 10,
            showDivider: false,
            onTap: () {},
          ),
        ],
      ),
    );
  }

  Widget _buildCell(
    BuildContext context, {
    required IconData icon,
    required Color iconColor,
    required String title,
    required double topPadding,
    required bool showDivider,
    VoidCallback? onTap,
  }) {
    return Container(
      padding: EdgeInsets.only(top: topPadding),
      color: const Color(0xFFEDEDED),
      child: Column(
        children: [
          GestureDetector(
            onTap: onTap,
            child: Container(
              height: 54,
              color: Colors.white,
              child: Row(
                children: [
                  const SizedBox(width: 16),
                  Icon(icon, size: 22, color: iconColor),
                  const SizedBox(width: 16),
                  Expanded(
                    child: Text(title, style: const TextStyle(fontSize: 16)),
                  ),
                  const Icon(Icons.arrow_forward_ios, size: 14, color: Colors.grey),
                  const SizedBox(width: 12),
                ],
              ),
            ),
          ),
          if (showDivider)
            Row(
              children: [
                Container(width: 54, color: Colors.white, height: 0.5),
                Expanded(
                  child: Container(
                    height: 0.5,
                    color: const Color(0x7DD9D9D9),
                  ),
                ),
              ],
            ),
        ],
      ),
    );
  }
}
