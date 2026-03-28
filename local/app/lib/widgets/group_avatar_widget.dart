import 'dart:math';
import 'package:flutter/material.dart';

/// 群组合头像组件
/// - 2人：左右各半圆
/// - 3人：三等分扇形
/// - 4+人：四象限（群主+随机3人）
class GroupAvatarWidget extends StatelessWidget {
  final List<String> faceURLs;
  final List<String> names;
  final double size;

  const GroupAvatarWidget({
    super.key,
    required this.faceURLs,
    required this.names,
    this.size = 48,
  });

  @override
  Widget build(BuildContext context) {
    final count = faceURLs.length;
    if (count == 0) {
      return _fallbackAvatar();
    }
    if (count == 1) {
      return _singleAvatar(faceURLs[0], names.isNotEmpty ? names[0] : '');
    }

    // 取前2/3/4个头像
    final displayCount = count >= 4 ? 4 : count;
    final urls = faceURLs.take(displayCount).toList();
    final nms = names.take(displayCount).toList();

    return SizedBox(
      width: size,
      height: size,
      child: ClipOval(
        child: Stack(
          children: List.generate(displayCount, (i) {
            return _buildSector(i, displayCount, urls[i], nms.length > i ? nms[i] : '');
          }),
        ),
      ),
    );
  }

  Widget _buildSector(int index, int total, String url, String name) {
    if (total == 2) {
      return _buildHalf(index, url, name);
    } else if (total == 3) {
      return _buildThird(index, url, name);
    } else {
      return _buildQuarter(index, url, name);
    }
  }

  /// 2人：左半圆 / 右半圆
  Widget _buildHalf(int index, String url, String name) {
    final isLeft = index == 0;
    return Positioned(
      left: isLeft ? 0 : size / 2,
      top: 0,
      child: ClipRect(
        child: SizedBox(
          width: size / 2,
          height: size,
          child: _avatarImage(url, name, width: size, height: size,
            alignment: isLeft ? Alignment.centerLeft : Alignment.centerRight),
        ),
      ),
    );
  }

  /// 3人：三等分
  Widget _buildThird(int index, String url, String name) {
    // 使用网格布局：上面1个居中，下面2个
    if (index == 0) {
      // 上方居中
      return Positioned(
        left: size / 4,
        top: 0,
        child: ClipOval(
          child: SizedBox(
            width: size / 2,
            height: size / 2,
            child: _tileImage(url, name, size / 2),
          ),
        ),
      );
    } else {
      // 下方左/右
      final isLeft = index == 1;
      return Positioned(
        left: isLeft ? 0 : size / 2,
        top: size / 2,
        child: ClipOval(
          child: SizedBox(
            width: size / 2,
            height: size / 2,
            child: _tileImage(url, name, size / 2),
          ),
        ),
      );
    }
  }

  /// 4人：四象限
  Widget _buildQuarter(int index, String url, String name) {
    final row = index ~/ 2;
    final col = index % 2;
    final half = size / 2;
    return Positioned(
      left: col * half,
      top: row * half,
      child: SizedBox(
        width: half,
        height: half,
        child: Padding(
          padding: const EdgeInsets.all(0.5),
          child: ClipOval(
            child: _tileImage(url, name, half - 1),
          ),
        ),
      ),
    );
  }

  Widget _tileImage(String url, String name, double tileSize) {
    if (url.isNotEmpty) {
      return Image.network(
        url,
        width: tileSize,
        height: tileSize,
        fit: BoxFit.cover,
        errorBuilder: (_, __, ___) => _letterTile(name, tileSize),
      );
    }
    return _letterTile(name, tileSize);
  }

  Widget _letterTile(String name, double tileSize) {
    final letter = name.isNotEmpty ? name[0].toUpperCase() : '?';
    final colors = [
      const Color(0xFF4CAF50),
      const Color(0xFF2196F3),
      const Color(0xFFFF9800),
      const Color(0xFF9C27B0),
      const Color(0xFFE91E63),
      const Color(0xFF00BCD4),
    ];
    final color = colors[name.hashCode.abs() % colors.length];
    return Container(
      width: tileSize,
      height: tileSize,
      color: color,
      alignment: Alignment.center,
      child: Text(
        letter,
        style: TextStyle(
          color: Colors.white,
          fontSize: tileSize * 0.45,
          fontWeight: FontWeight.w600,
        ),
      ),
    );
  }

  Widget _avatarImage(String url, String name, {
    required double width,
    required double height,
    Alignment alignment = Alignment.center,
  }) {
    if (url.isNotEmpty) {
      return Image.network(
        url,
        width: width,
        height: height,
        fit: BoxFit.cover,
        alignment: alignment,
        errorBuilder: (_, __, ___) => _letterTile(name, min(width, height)),
      );
    }
    return _letterTile(name, min(width, height));
  }

  Widget _singleAvatar(String url, String name) {
    return SizedBox(
      width: size,
      height: size,
      child: ClipOval(
        child: _tileImage(url, name, size),
      ),
    );
  }

  Widget _fallbackAvatar() {
    return Container(
      width: size,
      height: size,
      decoration: const BoxDecoration(
        color: Color(0xFF4CAF50),
        shape: BoxShape.circle,
      ),
      child: Icon(Icons.group, color: Colors.white, size: size * 0.5),
    );
  }
}
