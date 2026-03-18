import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../core/controllers/im_controller.dart';
import '../conversation/conversation_page.dart';
import '../contacts/contacts_page.dart';
import '../mine/mine_page.dart';

class HomePage extends StatefulWidget {
  const HomePage({super.key});

  @override
  State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {
  int _currentIndex = 0;
  final _pages = const [
    ConversationPage(),
    ContactsPage(),
    MinePage(),
  ];

  @override
  Widget build(BuildContext context) {
    final imCtrl = Get.find<IMController>();
    return Scaffold(
      body: IndexedStack(index: _currentIndex, children: _pages),
      bottomNavigationBar: Obx(() => BottomNavigationBar(
            currentIndex: _currentIndex,
            onTap: (i) => setState(() => _currentIndex = i),
            type: BottomNavigationBarType.fixed,
            selectedItemColor: const Color(0xFF1B72EC),
            unselectedItemColor: Colors.grey,
            items: [
              BottomNavigationBarItem(
                icon: _buildBadge(
                    const Icon(Icons.chat_bubble_outline), imCtrl.totalUnreadCount.value),
                activeIcon: _buildBadge(
                    const Icon(Icons.chat_bubble), imCtrl.totalUnreadCount.value),
                label: '消息',
              ),
              BottomNavigationBarItem(
                icon: _buildBadge(
                    const Icon(Icons.people_outline), imCtrl.friendApplyCount.value),
                activeIcon: _buildBadge(
                    const Icon(Icons.people), imCtrl.friendApplyCount.value),
                label: '联系人',
              ),
              const BottomNavigationBarItem(
                icon: Icon(Icons.person_outline),
                activeIcon: Icon(Icons.person),
                label: '我的',
              ),
            ],
          )),
    );
  }

  Widget _buildBadge(Widget icon, int count) {
    if (count <= 0) return icon;
    return Badge(
      label: Text(count > 99 ? '99+' : '$count',
          style: const TextStyle(fontSize: 10)),
      child: icon,
    );
  }
}
