import '../http_client.dart';
import '../config.dart';
import '../models/user_info.dart';

class FriendApi {
  static Future<List<FriendInfo>> getFriendList({
    int pageNumber = 1,
    int showNumber = 100,
  }) async {
    final data = await HttpClient.post('/friend/get_friend_list', data: {
      'userID': Config.userID,
      'pagination': {'pageNumber': pageNumber, 'showNumber': showNumber},
    });
    if (data == null) return [];
    final list = data['friendsInfo'] ?? data;
    if (list is List) {
      return list.map((e) => FriendInfo.fromJson(e)).toList();
    }
    return [];
  }

  static Future<void> addFriend(String toUserID, {String reqMsg = ''}) async {
    await HttpClient.post('/friend/add_friend', data: {
      'fromUserID': Config.userID,
      'toUserID': toUserID,
      'reqMsg': reqMsg,
    });
  }

  static Future<void> respondFriendApply({
    required String fromUserID,
    required int handleResult, // 1:同意 -1:拒绝
    String handleMsg = '',
  }) async {
    await HttpClient.post('/friend/add_friend_response', data: {
      'toUserID': Config.userID,
      'fromUserID': fromUserID,
      'handleResult': handleResult,
      'handleMsg': handleMsg,
    });
  }

  static Future<List<FriendApplicationInfo>> getFriendApplyList({
    int pageNumber = 1,
    int showNumber = 100,
  }) async {
    final data = await HttpClient.post('/friend/get_friend_apply_list', data: {
      'userID': Config.userID,
      'pagination': {'pageNumber': pageNumber, 'showNumber': showNumber},
    });
    if (data == null) return [];
    final list = data['friendRequests'] ?? data;
    if (list is List) {
      return list.map((e) => FriendApplicationInfo.fromJson(e)).toList();
    }
    return [];
  }

  static Future<void> deleteFriend(String friendUserID) async {
    await HttpClient.post('/friend/delete_friend', data: {
      'ownerUserID': Config.userID,
      'friendUserID': friendUserID,
    });
  }

  static Future<void> setFriendRemark(
      String friendUserID, String remark) async {
    await HttpClient.post('/friend/set_friend_remark', data: {
      'ownerUserID': Config.userID,
      'friendUserID': friendUserID,
      'remark': remark,
    });
  }

  static Future<void> addBlack(String blackUserID) async {
    await HttpClient.post('/friend/add_black', data: {
      'ownerUserID': Config.userID,
      'blackUserID': blackUserID,
    });
  }

  static Future<void> removeBlack(String blackUserID) async {
    await HttpClient.post('/friend/remove_black', data: {
      'ownerUserID': Config.userID,
      'blackUserID': blackUserID,
    });
  }

  /// 搜索好友 (chat: /friend/search)
  static Future<List<dynamic>> searchFriend(
    String keyword, {
    int pageNumber = 1,
    int showNumber = 20,
  }) async {
    final data = await HttpClient.post('/friend/search', data: {
      'keyword': keyword,
      'pagination': {'pageNumber': pageNumber, 'showNumber': showNumber},
    });
    if (data == null) return [];
    final users = data['users'] ?? data['friends'] ?? data;
    if (users is List) return users;
    return [];
  }
}
