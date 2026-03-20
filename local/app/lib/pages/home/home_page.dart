import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:persistent_bottom_nav_bar_v2/persistent_bottom_nav_bar_v2.dart';

import '../../core/controllers/im_controller.dart';
import '../contacts/contacts_page.dart';
import '../conversation/conversation_page.dart';
import '../mine/mine_page.dart';

class HomePage extends StatelessWidget {
  const HomePage({super.key});

  List<PersistentTabConfig> _tabs(IMController imCtrl) => [
    PersistentTabConfig(
      screen: const ConversationPage(),
      item: ItemConfig(
        icon: _buildBadge(
          const Icon(Icons.chat_bubble, size: 24),
          imCtrl.totalUnreadCount.value,
        ),
        inactiveIcon: _buildBadge(
          const Icon(Icons.chat_bubble_outline, size: 24),
          imCtrl.totalUnreadCount.value,
        ),
        title: '消息',
        activeForegroundColor: const Color(0xFF0089FF),
        inactiveForegroundColor: const Color(0xFF8E9AB0),
      ),
    ),
    PersistentTabConfig(
      screen: const ContactsPage(),
      item: ItemConfig(
        icon: _buildBadge(
          const Icon(Icons.people, size: 24),
          imCtrl.friendApplyCount.value,
        ),
        inactiveIcon: _buildBadge(
          const Icon(Icons.people_outline, size: 24),
          imCtrl.friendApplyCount.value,
        ),
        title: '联系人',
        activeForegroundColor: const Color(0xFF0089FF),
        inactiveForegroundColor: const Color(0xFF8E9AB0),
      ),
    ),
    PersistentTabConfig(
      screen: const MinePage(),
      item: ItemConfig(
        icon: const Icon(Icons.person, size: 24),
        inactiveIcon: const Icon(Icons.person_outline, size: 24),
        title: '我的',
        activeForegroundColor: const Color(0xFF0089FF),
        inactiveForegroundColor: const Color(0xFF8E9AB0),
      ),
    ),
  ];

  @override
  Widget build(BuildContext context) {
    final imCtrl = Get.find<IMController>();
    return Scaffold(
      backgroundColor: Colors.white,
      body: Obx(
        () => PersistentTabView(
          tabs: _tabs(imCtrl),
          navBarBuilder: (navBarConfig) => Style1BottomNavBar(
            navBarConfig: navBarConfig,
            navBarDecoration: const NavBarDecoration(
              color: Colors.white,
              boxShadow: [
                BoxShadow(
                  color: Colors.black12,
                  blurRadius: 0.5,
                  spreadRadius: 0.5,
                ),
              ],
            ),
          ),
          navBarOverlap: const NavBarOverlap.none(),
          screenTransitionAnimation: const ScreenTransitionAnimation.none(),
        ),
      ),
    );
  }

  Widget _buildBadge(Widget icon, int count) {
    return Stack(
      clipBehavior: Clip.none,
      children: [
        icon,
        if (count > 0)
          Positioned(
            right: -8,
            top: -6,
            child: Container(
              constraints: const BoxConstraints(minWidth: 16),
              height: 16,
              padding: const EdgeInsets.symmetric(horizontal: 4),
              decoration: BoxDecoration(
                color: const Color(0xFFFF3B30),
                borderRadius: BorderRadius.circular(8),
              ),
              alignment: Alignment.center,
              child: Text(
                count > 99 ? '99+' : '$count',
                style: const TextStyle(
                  color: Colors.white,
                  fontSize: 9,
                  fontWeight: FontWeight.w600,
                ),
              ),
            ),
          ),
      ],
    );
  }
}
