import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../core/controllers/im_controller.dart';
import '../../routes/app_routes.dart';
import '../contacts/contacts_page.dart';
import '../conversation/conversation_page.dart';
import '../discover/discover_page.dart';
import '../mine/mine_page.dart';

class HomePage extends StatefulWidget {
  const HomePage({super.key});

  @override
  State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {
  int _currentIndex = 0;
  static const double _iconSize = 25.0;
  static const Color _activeColor = Color(0xFF1B72EC);
  static const Color _inactiveColor = Colors.grey;

  final List<Widget> _pages = const [
    ConversationPage(),
    ContactsPage(),
    DiscoverPage(),
    MinePage(),
  ];

  final List<String> _titles = const ['消息', '联系人', '发现', '我'];

  @override
  Widget build(BuildContext context) {
    final imCtrl = Get.find<IMController>();
    return Obx(() {
      return Scaffold(
        appBar: AppBar(
          backgroundColor: const Color(0xFFF2F2F2),
          elevation: 0,
          centerTitle: true,
          title: Text(_titles[_currentIndex]),
          actions: _currentIndex == 3
              ? null
              : [
                  IconButton(
                    icon: const Icon(Icons.search, color: Color(0xFF191919)),
                    onPressed: () => Get.toNamed(AppRoutes.globalSearch),
                  ),
                  PopupMenuButton<String>(
                    icon: const Icon(Icons.add_circle_outline, color: Color(0xFF191919)),
                    onSelected: (v) {
                      if (v == 'group') Get.toNamed(AppRoutes.createGroup);
                      if (v == 'friend') Get.toNamed(AppRoutes.addFriend);
                    },
                    itemBuilder: (_) => [
                      const PopupMenuItem(value: 'friend', child: Text('添加好友')),
                      const PopupMenuItem(value: 'group', child: Text('发起群聊')),
                    ],
                  ),
                ],
        ),
        body: IndexedStack(
          index: _currentIndex,
          children: _pages,
        ),
        bottomNavigationBar: BottomNavigationBar(
          currentIndex: _currentIndex,
          type: BottomNavigationBarType.fixed,
          selectedItemColor: _activeColor,
          unselectedItemColor: _inactiveColor,
          selectedFontSize: 11,
          unselectedFontSize: 11,
          onTap: (i) => setState(() => _currentIndex = i),
          items: [
            BottomNavigationBarItem(
              icon: _buildBadge(
                Image.asset('assets/tabbar_mainframe@3x.png', width: _iconSize, height: _iconSize),
                imCtrl.totalUnreadCount.value,
              ),
              activeIcon: _buildBadge(
                ColorFiltered(
                  colorFilter: const ColorFilter.mode(_activeColor, BlendMode.srcIn),
                  child: Image.asset('assets/tabbar_mainframe@3x.png', width: _iconSize, height: _iconSize),
                ),
                imCtrl.totalUnreadCount.value,
              ),
              label: '消息',
            ),
            BottomNavigationBarItem(
              icon: _buildBadge(
                Image.asset('assets/tabbar_contacts@3x.png', width: _iconSize, height: _iconSize),
                imCtrl.friendApplyCount.value,
              ),
              activeIcon: _buildBadge(
                ColorFiltered(
                  colorFilter: const ColorFilter.mode(_activeColor, BlendMode.srcIn),
                  child: Image.asset('assets/tabbar_contacts@3x.png', width: _iconSize, height: _iconSize),
                ),
                imCtrl.friendApplyCount.value,
              ),
              label: '联系人',
            ),
            BottomNavigationBarItem(
              icon: Image.asset('assets/tabbar_discover@3x.png', width: _iconSize, height: _iconSize),
              activeIcon: ColorFiltered(
                colorFilter: const ColorFilter.mode(_activeColor, BlendMode.srcIn),
                child: Image.asset('assets/tabbar_discover@3x.png', width: _iconSize, height: _iconSize),
              ),
              label: '发现',
            ),
            BottomNavigationBarItem(
              icon: Image.asset('assets/tabbar_me@3x.png', width: _iconSize, height: _iconSize),
              activeIcon: ColorFiltered(
                colorFilter: const ColorFilter.mode(_activeColor, BlendMode.srcIn),
                child: Image.asset('assets/tabbar_me@3x.png', width: _iconSize, height: _iconSize),
              ),
              label: '我',
            ),
          ],
        ),
      );
    });
  }

  Widget _buildBadge(Widget icon, int count) {
    return Stack(
      clipBehavior: Clip.none,
      children: [
        icon,
        if (count > 0)
          Positioned(
            right: -8,
            top: -4,
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
