import 'package:flutter/material.dart';

/// Emoji picker widget for chat input
class EmojiPickerWidget extends StatelessWidget {
  final Function(String) onEmojiSelected;
  final VoidCallback? onBackspacePressed;

  const EmojiPickerWidget({
    super.key,
    required this.onEmojiSelected,
    this.onBackspacePressed,
  });

  // Common emoji categories
  static const List<List<String>> _emojiCategories = [
    // Faces
    [
      '😀',
      '😃',
      '😄',
      '😁',
      '😅',
      '😂',
      '🤣',
      '😊',
      '😇',
      '🙂',
      '🙃',
      '😉',
      '😌',
      '😍',
      '🥰',
      '😘',
      '😗',
      '😙',
      '😚',
      '😋',
      '😛',
      '😜',
      '🤪',
      '😝',
      '🤑',
      '🤗',
      '🤭',
      '🤫',
      '🤔',
      '🤐',
      '🤨',
      '😐',
      '😑',
      '😶',
      '😏',
      '😒',
      '🙄',
      '😬',
      '🤥',
      '😌',
      '😔',
      '😪',
      '🤤',
      '😴',
      '😷',
      '🤒',
      '🤕',
      '🤢',
      '🤮',
      '🤧',
    ],
    // Gestures
    [
      '👋',
      '🤚',
      '🖐️',
      '✋',
      '🖖',
      '👌',
      '🤌',
      '🤏',
      '✌️',
      '🤞',
      '🤟',
      '🤘',
      '🤙',
      '👈',
      '👉',
      '👆',
      '🖕',
      '👇',
      '☝️',
      '👍',
      '👎',
      '✊',
      '👊',
      '🤛',
      '🤜',
      '👏',
      '🙌',
      '👐',
      '🤲',
      '🤝',
      '🙏',
      '✍️',
      '💪',
      '🦾',
      '🦿',
      '🦵',
      '🦶',
      '👂',
      '🦻',
      '👃',
    ],
    // Hearts
    [
      '❤️',
      '🧡',
      '💛',
      '💚',
      '💙',
      '💜',
      '🖤',
      '🤍',
      '🤎',
      '💔',
      '❣️',
      '💕',
      '💞',
      '💓',
      '💗',
      '💖',
      '💘',
      '💝',
      '💟',
      '♥️',
    ],
    // Objects
    [
      '🎉',
      '🎊',
      '🎈',
      '🎁',
      '🎀',
      '🏆',
      '🥇',
      '🥈',
      '🥉',
      '⚽',
      '🏀',
      '🏈',
      '⚾',
      '🥎',
      '🎾',
      '🏐',
      '🏉',
      '🥏',
      '🎱',
      '🪀',
      '🏓',
      '🏸',
      '🏒',
      '🏑',
      '🥍',
      '🏏',
      '🪃',
      '🥅',
      '⛳',
      '🪁',
    ],
  ];

  static const List<String> _categoryIcons = ['😀', '👋', '❤️', '🎉'];

  @override
  Widget build(BuildContext context) {
    return Container(
      height: 250,
      color: const Color(0xFFF0F2F6),
      child: Column(
        children: [
          _buildCategoryTabs(),
          Expanded(
            child: _buildEmojiGrid(0), // Default to first category
          ),
        ],
      ),
    );
  }

  Widget _buildCategoryTabs() {
    return Container(
      height: 44,
      color: Colors.white,
      child: Row(
        children: List.generate(
          _categoryIcons.length,
          (index) => Expanded(
            child: GestureDetector(
              onTap: () {},
              child: Center(
                child: Text(
                  _categoryIcons[index],
                  style: const TextStyle(fontSize: 20),
                ),
              ),
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildEmojiGrid(int categoryIndex) {
    final emojis = _emojiCategories[categoryIndex];
    return GridView.builder(
      padding: const EdgeInsets.all(8),
      gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
        crossAxisCount: 8,
        mainAxisSpacing: 4,
        crossAxisSpacing: 4,
      ),
      itemCount: emojis.length + 1, // +1 for backspace
      itemBuilder: (context, index) {
        if (index == emojis.length) {
          // Backspace button
          return GestureDetector(
            onTap: onBackspacePressed,
            child: const Center(
              child: Icon(Icons.backspace, size: 20, color: Colors.grey),
            ),
          );
        }
        return GestureDetector(
          onTap: () => onEmojiSelected(emojis[index]),
          child: Center(
            child: Text(emojis[index], style: const TextStyle(fontSize: 24)),
          ),
        );
      },
    );
  }
}

/// Emoji picker sheet for bottom display
class EmojiPickerSheet extends StatefulWidget {
  final Function(String) onEmojiSelected;
  final VoidCallback? onBackspacePressed;

  const EmojiPickerSheet({
    super.key,
    required this.onEmojiSelected,
    this.onBackspacePressed,
  });

  @override
  State<EmojiPickerSheet> createState() => _EmojiPickerSheetState();
}

class _EmojiPickerSheetState extends State<EmojiPickerSheet>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(
      length: EmojiPickerWidget._emojiCategories.length,
      vsync: this,
    );
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      height: 280,
      color: const Color(0xFFF0F2F6),
      child: Column(
        children: [
          _buildTabBar(),
          Expanded(
            child: TabBarView(
              controller: _tabController,
              children: EmojiPickerWidget._emojiCategories
                  .asMap()
                  .entries
                  .map((entry) => _buildEmojiGrid(entry.value))
                  .toList(),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildTabBar() {
    return Container(
      height: 44,
      color: Colors.white,
      child: TabBar(
        controller: _tabController,
        indicatorColor: const Color(0xFF0089FF),
        labelColor: const Color(0xFF0089FF),
        unselectedLabelColor: Colors.grey,
        tabs: EmojiPickerWidget._categoryIcons
            .map((icon) => Tab(text: icon))
            .toList(),
      ),
    );
  }

  Widget _buildEmojiGrid(List<String> emojis) {
    return GridView.builder(
      padding: const EdgeInsets.all(8),
      gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
        crossAxisCount: 8,
        mainAxisSpacing: 4,
        crossAxisSpacing: 4,
      ),
      itemCount: emojis.length + 1,
      itemBuilder: (context, index) {
        if (index == emojis.length) {
          return GestureDetector(
            onTap: widget.onBackspacePressed,
            child: const Center(
              child: Icon(Icons.backspace, size: 20, color: Colors.grey),
            ),
          );
        }
        return GestureDetector(
          onTap: () => widget.onEmojiSelected(emojis[index]),
          child: Center(
            child: Text(emojis[index], style: const TextStyle(fontSize: 24)),
          ),
        );
      },
    );
  }
}
