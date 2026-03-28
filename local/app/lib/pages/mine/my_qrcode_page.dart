import 'dart:io';
import 'dart:ui' as ui;

import 'package:flutter/material.dart';
import 'package:flutter/rendering.dart';
import 'package:flutter_easyloading/flutter_easyloading.dart';
import 'package:get/get.dart';
import 'package:path_provider/path_provider.dart';
import 'package:qr_flutter/qr_flutter.dart';

import '../../core/config.dart';
import '../../core/controllers/auth_controller.dart';

class MyQrcodePage extends StatefulWidget {
  const MyQrcodePage({super.key});

  @override
  State<MyQrcodePage> createState() => _MyQrcodePageState();
}

class _MyQrcodePageState extends State<MyQrcodePage> {
  final GlobalKey _qrKey = GlobalKey();

  Future<void> _saveQrCode() async {
    try {
      EasyLoading.show(status: '保存中...');
      
      final boundary = _qrKey.currentContext?.findRenderObject() as RenderRepaintBoundary?;
      if (boundary == null) {
        EasyLoading.showError('保存失败');
        return;
      }
      
      final image = await boundary.toImage(pixelRatio: 3.0);
      final byteData = await image.toByteData(format: ui.ImageByteFormat.png);
      if (byteData == null) {
        EasyLoading.showError('保存失败');
        return;
      }
      
      final bytes = byteData.buffer.asUint8List();
      
      // 获取应用文档目录
      final directory = await getApplicationDocumentsDirectory();
      final qrDir = Directory('${directory.path}/qrcode');
      if (!await qrDir.exists()) {
        await qrDir.create(recursive: true);
      }
      
      final fileName = 'my_qrcode_${DateTime.now().millisecondsSinceEpoch}.png';
      final file = File('${qrDir.path}/$fileName');
      await file.writeAsBytes(bytes);
      
      EasyLoading.showSuccess('已保存到 ${qrDir.path}');
    } catch (e) {
      debugPrint('[MyQrcodePage] Save error: $e');
      EasyLoading.showError('保存失败');
    }
  }

  @override
  Widget build(BuildContext context) {
    final authCtrl = Get.find<AuthController>();
    final user = authCtrl.userInfo.value;
    final nickname = user?.nickname ?? Config.nickname;
    final userID = user?.userID ?? Config.userID;
    final faceURL = user?.faceURL ?? '';

    return Scaffold(
      backgroundColor: Colors.white,
      appBar: AppBar(
        title: const Text('我的二维码'),
        backgroundColor: Colors.white,
        elevation: 0,
      ),
      body: Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            // 昵称
            Text(
              nickname,
              style: const TextStyle(
                fontSize: 20,
                fontWeight: FontWeight.w600,
              ),
            ),
            const SizedBox(height: 8),
            // ID
            Text(
              'ID: $userID',
              style: const TextStyle(
                fontSize: 14,
                color: Colors.grey,
              ),
            ),
            const SizedBox(height: 24),
            // 二维码（头像在中央）- 用 RepaintBoundary 包裹以便截图
            RepaintBoundary(
              key: _qrKey,
              child: Container(
                color: Colors.white,
                padding: const EdgeInsets.all(16),
                child: Stack(
                  alignment: Alignment.center,
                  children: [
                    QrImageView(
                      data: 'fdim://user/$userID',
                      version: QrVersions.auto,
                      size: 200,
                      backgroundColor: Colors.white,
                    ),
                    Container(
                      width: 50,
                      height: 50,
                      decoration: BoxDecoration(
                        color: Colors.white,
                        borderRadius: BorderRadius.circular(4),
                      ),
                      padding: const EdgeInsets.all(2),
                      child: ClipRRect(
                        borderRadius: BorderRadius.circular(2),
                        child: faceURL.isNotEmpty
                            ? Image.network(faceURL, fit: BoxFit.cover)
                            : Container(
                                color: Colors.grey[300],
                                child: const Icon(Icons.person, size: 30, color: Colors.white),
                              ),
                      ),
                    ),
                  ],
                ),
              ),
            ),
            const SizedBox(height: 16),
            // 提示文字
            const Text(
              '扫一扫上面的二维码，加我为好友',
              style: TextStyle(
                fontSize: 13,
                color: Colors.grey,
              ),
            ),
            const SizedBox(height: 24),
            // 保存图片按钮
            TextButton(
              onPressed: _saveQrCode,
              child: const Text(
                '保存图片',
                style: TextStyle(
                  fontSize: 14,
                  color: Color(0xFF1B72EC),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
