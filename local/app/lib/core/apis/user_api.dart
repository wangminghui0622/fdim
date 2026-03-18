import '../http_client.dart';
import '../config.dart';
import '../models/user_info.dart';

class UserApi {
  /// IM 层获取用户公开信息 (open-im-server: /user/get_users_info)
  static Future<List<UserInfo>> getUsersInfo(List<String> userIDs) async {
    final data = await HttpClient.post('/user/get_users_info', data: {
      'userIDs': userIDs,
    });
    if (data is List) {
      return data.map((e) => UserInfo.fromJson(e)).toList();
    }
    if (data is Map && data['users'] is List) {
      return (data['users'] as List)
          .map((e) => UserInfo.fromJson(e))
          .toList();
    }
    return [];
  }

  /// Chat 层获取用户完整信息 (chat: /user/find/full)
  static Future<List<dynamic>> getUserFullInfo(List<String> userIDs) async {
    final data = await HttpClient.post('/user/find/full', data: {
      'userIDs': userIDs,
      'pagination': {'pageNumber': 0, 'showNumber': userIDs.length},
      'platform': Config.platformID,
    });
    if (data is Map && data['users'] is List) {
      return data['users'] as List;
    }
    return [];
  }

  /// Chat 层搜索用户 (chat: /user/search/full)
  static Future<List<dynamic>> searchUserFullInfo(
    String keyword, {
    int pageNumber = 1,
    int showNumber = 20,
  }) async {
    final data = await HttpClient.post('/user/search/full', data: {
      'keyword': keyword,
      'pagination': {'pageNumber': pageNumber, 'showNumber': showNumber},
    });
    if (data is Map && data['users'] is List) {
      return data['users'] as List;
    }
    return [];
  }

  /// IM 层更新用户信息 (open-im-server: /user/update_user_info)
  static Future<void> updateUserInfo(Map<String, dynamic> info) async {
    await HttpClient.post('/user/update_user_info', data: info);
  }

  /// Chat 层更新用户信息 (chat: /user/update)
  static Future<void> chatUpdateUserInfo(Map<String, dynamic> info) async {
    await HttpClient.post('/user/update', data: {
      ...info,
      'platform': Config.platformID,
    });
  }

  /// Chat 层搜索公开用户信息 (chat: /user/search/public)
  static Future<List<dynamic>> searchUserPublic(
    String keyword, {
    int pageNumber = 1,
    int showNumber = 20,
  }) async {
    final data = await HttpClient.post('/user/search/public', data: {
      'keyword': keyword,
      'pagination': {'pageNumber': pageNumber, 'showNumber': showNumber},
    });
    if (data is Map && data['users'] is List) {
      return data['users'] as List;
    }
    return [];
  }

  /// 在线状态 (open-im-server: /user/get_users_online_status)
  static Future<List<Map<String, dynamic>>> getUsersOnlineStatus(
      List<String> userIDs) async {
    final data = await HttpClient.post(
      '/user/get_users_online_status',
      data: {'userIDs': userIDs},
      showErrorToast: false,
    );
    if (data == null) return [];
    final successResult = data['successResult'];
    if (successResult is List) {
      return successResult.cast<Map<String, dynamic>>();
    }
    return [];
  }
}
