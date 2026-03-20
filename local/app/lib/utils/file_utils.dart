import 'dart:io';
import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:flutter_easyloading/flutter_easyloading.dart';
import 'package:path_provider/path_provider.dart';
import 'package:url_launcher/url_launcher.dart';

/// File download and open utilities
class FileUtils {
  static final Dio _dio = Dio();

  /// Download file to local storage
  static Future<String?> downloadFile({
    required String url,
    required String fileName,
    Function(int, int)? onProgress,
  }) async {
    try {
      final dir = await getApplicationDocumentsDirectory();
      final savePath = '${dir.path}/$fileName';

      await _dio.download(
        url,
        savePath,
        onReceiveProgress: (received, total) {
          if (total != -1) {
            onProgress?.call(received, total);
          }
        },
      );

      return savePath;
    } catch (e) {
      debugPrint('[FileUtils] Download error: $e');
      return null;
    }
  }

  /// Open file with system default app
  static Future<bool> openFile(String filePath) async {
    try {
      final file = File(filePath);
      if (!await file.exists()) {
        return false;
      }

      final uri = Uri.file(filePath);
      if (await canLaunchUrl(uri)) {
        return await launchUrl(uri, mode: LaunchMode.externalApplication);
      }
      return false;
    } catch (e) {
      debugPrint('[FileUtils] Open file error: $e');
      return false;
    }
  }

  /// Download and open file with progress
  static Future<void> downloadAndOpen({
    required BuildContext context,
    required String url,
    required String fileName,
  }) async {
    EasyLoading.show(status: '下载中...', maskType: EasyLoadingMaskType.black);

    final filePath = await downloadFile(
      url: url,
      fileName: fileName,
      onProgress: (received, total) {
        final progress = (received / total * 100).toStringAsFixed(0);
        EasyLoading.showProgress(received / total, status: '下载中 $progress%');
      },
    );

    EasyLoading.dismiss();

    if (filePath != null) {
      final opened = await openFile(filePath);
      if (!opened) {
        EasyLoading.showToast('无法打开文件');
      }
    } else {
      EasyLoading.showToast('下载失败');
    }
  }

  /// Open URL in browser
  static Future<bool> openUrl(String url) async {
    try {
      final uri = Uri.parse(url);
      if (await canLaunchUrl(uri)) {
        return await launchUrl(uri, mode: LaunchMode.externalApplication);
      }
      return false;
    } catch (e) {
      debugPrint('[FileUtils] Open URL error: $e');
      return false;
    }
  }
}
