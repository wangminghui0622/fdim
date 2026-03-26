import 'dart:io';

import 'package:dio/dio.dart';
import 'package:flutter/foundation.dart';

import 'http_client.dart';

/// 文件上传服务 — 使用服务端 Object Storage (MinIO) 的 FormData 上传流程
/// 流程：initiate_form_data → 上传文件 → complete_form_data → 返回永久 URL
class UploadService {
  UploadService._();

  static final Dio _uploadDio = Dio(BaseOptions(
    connectTimeout: const Duration(seconds: 60),
    receiveTimeout: const Duration(seconds: 120),
    sendTimeout: const Duration(seconds: 120),
  ));

  /// 上传本地文件，返回可访问的永久 URL
  /// [filePath] 本地文件路径
  /// [group] 存储分组，如 "voice", "image", "file"
  static Future<String> uploadFile(String filePath, {String group = ''}) async {
    final file = File(filePath);
    if (!await file.exists()) {
      throw Exception('文件不存在: $filePath');
    }

    final fileSize = await file.length();
    final fileName = filePath.split(Platform.pathSeparator).last;
    final contentType = _guessContentType(fileName);

    debugPrint('[Upload] initiate: name=$fileName, size=$fileSize, contentType=$contentType');

    // 1. 申请上传（获取预签名 URL 和 form 字段）
    final initResp = await HttpClient.post('/object/initiate_form_data', data: {
      'name': fileName,
      'size': fileSize,
      'contentType': contentType,
      'group': group,
    });

    if (initResp is! Map) {
      throw Exception('initiate_form_data 返回格式错误');
    }

    final uploadId = initResp['id'] as String? ?? '';
    final uploadUrl = initResp['url'] as String? ?? '';
    final fileField = initResp['file'] as String? ?? 'file';
    final formDataMap = initResp['formData'];

    if (uploadUrl.isEmpty || uploadId.isEmpty) {
      throw Exception('initiate_form_data 缺少 url 或 id');
    }

    debugPrint('[Upload] uploading to: $uploadUrl');

    // 2. 通过 multipart form-data 上传文件到预签名 URL
    final formData = FormData();

    // 添加服务端要求的 form 字段
    if (formDataMap is Map) {
      for (final entry in formDataMap.entries) {
        formData.fields.add(MapEntry(entry.key.toString(), entry.value.toString()));
      }
    }

    // 添加文件字段
    formData.files.add(MapEntry(
      fileField,
      await MultipartFile.fromFile(filePath, filename: fileName),
    ));

    final uploadResp = await _uploadDio.post(
      uploadUrl,
      data: formData,
      options: Options(
        headers: {'Content-Type': 'multipart/form-data'},
        validateStatus: (status) => status != null && status < 500,
      ),
    );
    debugPrint('[Upload] upload response: ${uploadResp.statusCode}');

    // 3. 确认上传完成，获取永久访问 URL
    final completeResp = await HttpClient.post('/object/complete_form_data', data: {
      'id': uploadId,
    });

    if (completeResp is! Map || (completeResp['url'] ?? '').toString().isEmpty) {
      throw Exception('complete_form_data 缺少 url');
    }

    final accessUrl = completeResp['url'] as String;
    debugPrint('[Upload] complete: url=$accessUrl');
    return accessUrl;
  }

  static String _guessContentType(String fileName) {
    final ext = fileName.split('.').last.toLowerCase();
    switch (ext) {
      case 'm4a':
      case 'aac':
        return 'audio/aac';
      case 'mp3':
        return 'audio/mpeg';
      case 'wav':
        return 'audio/wav';
      case 'ogg':
        return 'audio/ogg';
      case 'jpg':
      case 'jpeg':
        return 'image/jpeg';
      case 'png':
        return 'image/png';
      case 'gif':
        return 'image/gif';
      case 'webp':
        return 'image/webp';
      case 'mp4':
        return 'video/mp4';
      case 'pdf':
        return 'application/pdf';
      default:
        return 'application/octet-stream';
    }
  }
}
