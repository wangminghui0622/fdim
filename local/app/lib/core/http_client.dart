import 'package:dio/dio.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter_easyloading/flutter_easyloading.dart';
import 'config.dart';

/// 统一 API 响应结构（与官方OpenIM保持一致）
class ApiResp {
  final int errCode;
  final String errMsg;
  final String errDlt;  // Error detail - 与官方保持一致
  final dynamic data;

  ApiResp({
    required this.errCode,
    required this.errMsg,
    this.errDlt = '',
    this.data
  });

  factory ApiResp.fromJson(Map<String, dynamic> json) {
    return ApiResp(
      errCode: json['errCode'] ?? 0,
      errMsg: json['errMsg'] ?? '',
      errDlt: json['errDlt'] ?? '',
      data: json['data'],
    );
  }

  bool get isSuccess => errCode == 0;
}

class HttpClient {
  HttpClient._();

  static final Dio _dio = Dio();

  static void init() {
    _dio.options
      ..baseUrl = Config.apiUrl
      ..connectTimeout = const Duration(seconds: 30)
      ..receiveTimeout = const Duration(seconds: 30)
      ..headers = {'Content-Type': 'application/json'};

    _dio.interceptors.add(InterceptorsWrapper(
      onRequest: (options, handler) {
        final token = Config.token;
        if (token.isNotEmpty) {
          options.headers['token'] = token;
        }
        options.headers['operationID'] =
            DateTime.now().millisecondsSinceEpoch.toString();
        if (kDebugMode) {
          debugPrint('[HTTP] ${options.method} ${options.path}');
        }
        handler.next(options);
      },
      onResponse: (response, handler) {
        handler.next(response);
      },
      onError: (error, handler) {
        if (kDebugMode) {
          debugPrint('[HTTP] Error: ${error.message}');
        }
        handler.next(error);
      },
    ));
  }

  /// 通用 POST 请求，返回 data 字段
  static Future<dynamic> post(
    String path, {
    dynamic data,
    bool showErrorToast = true,
    Options? options,
    CancelToken? cancelToken,
  }) async {
    try {
      data ??= {};
      final result = await _dio.post<Map<String, dynamic>>(
        path,
        data: data,
        options: options,
        cancelToken: cancelToken,
      );
      final resp = ApiResp.fromJson(result.data!);
      if (resp.isSuccess) {
        return resp.data;
      } else {
        if (showErrorToast) {
          EasyLoading.showToast(resp.errMsg.isNotEmpty ? resp.errMsg : '请求失败');
        }
        return Future.error(ApiException(resp.errCode, resp.errMsg));
      }
    } on DioException catch (e) {
      String desc;
      switch (e.type) {
        case DioExceptionType.connectionTimeout:
          desc = '连接超时';
          break;
        case DioExceptionType.sendTimeout:
          desc = '发送超时';
          break;
        case DioExceptionType.receiveTimeout:
          desc = '接收超时';
          break;
        case DioExceptionType.connectionError:
          desc = '无法连接服务器';
          break;
        case DioExceptionType.cancel:
          desc = '请求已取消';
          break;
        default:
          desc = e.message ?? '请检查网络连接';
      }
      final msg = '网络错误: $desc';
      if (showErrorToast) EasyLoading.showToast(msg);
      return Future.error(msg);
    } catch (e) {
      if (showErrorToast) EasyLoading.showToast(e.toString());
      return Future.error(e);
    }
  }

  /// 更新 baseUrl（切换服务器时用）
  static void updateBaseUrl(String url) {
    _dio.options.baseUrl = url;
  }
}

class ApiException implements Exception {
  final int code;
  final String message;
  ApiException(this.code, this.message);

  @override
  String toString() => 'ApiException($code): $message';
}
