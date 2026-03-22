import 'package:flutter/foundation.dart';

import '../../flutter_openim_sdk.dart';
import '../../../core/http_client.dart';
import '../../../core/config.dart';

class FriendshipManager {
  late OnFriendshipListener listener;

  Future setFriendshipListener(OnFriendshipListener listener) async {
    this.listener = listener;
  }

  Future<List<FriendInfo>> getFriendsInfo({
    required List<String> userIDList,
    bool filterBlack = false,
    String? operationID,
  }) async {
    final data = await HttpClient.post('/friend/get_designated_friends', data: {
      'ownerUserID': Config.userID,
      'friendUserIDs': userIDList,
    });
    return _parseFriendList(data);
  }

  Future<dynamic> addFriend({
    required String userID,
    String? reason,
    String? operationID,
  }) async {
    await HttpClient.post('/friend/add_friend', data: {
      'fromUserID': Config.userID,
      'toUserID': userID,
      'reqMsg': reason ?? '',
    });
  }

  Future<List<FriendApplicationInfo>> getFriendApplicationListAsRecipient({
    int offset = 0,
    int count = 100,
    String? operationID,
  }) async {
    final data = await HttpClient.post('/friend/get_friend_apply_list', data: {
      'userID': Config.userID,
      'pagination': {
        'pageNumber': offset ~/ count + 1,
        'showNumber': count,
      },
    });
    final list = _parseFriendApplicationList(data);
    // 服务端不返回 fromNickname/fromFaceURL，需要额外查询用户信息补全
    await _enrichApplicationsWithUserInfo(list);
    return list;
  }

  Future<List<FriendApplicationInfo>> getFriendApplicationListAsApplicant({
    int offset = 0,
    int count = 100,
    String? operationID,
  }) async {
    final data = await HttpClient.post('/friend/get_self_friend_apply_list', data: {
      'userID': Config.userID,
      'pagination': {
        'pageNumber': offset ~/ count + 1,
        'showNumber': count,
      },
    });
    return _parseFriendApplicationList(data);
  }

  Future<List<FriendInfo>> getFriendList({
    String? operationID,
    bool filterBlack = false,
  }) async {
    final data = await HttpClient.post('/friend/get_friend_list', data: {
      'userID': Config.userID,
      'pagination': {'pageNumber': 1, 'showNumber': 10000},
    });
    return _parseFriendList(data);
  }

  Future<List<FriendInfo>> getFriendListPage({
    bool filterBlack = false,
    int offset = 0,
    int count = 40,
    String? operationID,
  }) async {
    final data = await HttpClient.post('/friend/get_friend_list', data: {
      'userID': Config.userID,
      'pagination': {
        'pageNumber': offset ~/ count + 1,
        'showNumber': count,
      },
    });
    return _parseFriendList(data);
  }

  Future<dynamic> acceptFriendApplication({
    required String userID,
    String? handleMsg,
    String? operationID,
  }) async {
    await HttpClient.post('/friend/add_friend_response', data: {
      'toUserID': Config.userID,
      'fromUserID': userID,
      'handleResult': 1,
      'handleMsg': handleMsg ?? '',
    });
  }

  Future<dynamic> refuseFriendApplication({
    required String userID,
    String? handleMsg,
    String? operationID,
  }) async {
    await HttpClient.post('/friend/add_friend_response', data: {
      'toUserID': Config.userID,
      'fromUserID': userID,
      'handleResult': -1,
      'handleMsg': handleMsg ?? '',
    });
  }

  Future<dynamic> deleteFriend({
    required String userID,
    String? operationID,
  }) async {
    await HttpClient.post('/friend/delete_friend', data: {
      'ownerUserID': Config.userID,
      'friendUserID': userID,
    });
  }

  Future<dynamic> addBlacklist({
    required String userID,
    String? ex,
    String? operationID,
  }) async {
    await HttpClient.post('/friend/add_black', data: {
      'ownerUserID': Config.userID,
      'blackUserID': userID,
    });
  }

  Future<List<BlacklistInfo>> getBlacklist({String? operationID}) async {
    final data = await HttpClient.post('/friend/get_black_list', data: {
      'userID': Config.userID,
      'pagination': {'pageNumber': 1, 'showNumber': 10000},
    });
    if (data == null) return [];
    final list = data['blacks'] ?? data;
    if (list is List) {
      return list.map((e) => BlacklistInfo.fromJson(e)).toList();
    }
    return [];
  }

  Future<dynamic> removeBlacklist({
    required String userID,
    String? operationID,
  }) async {
    await HttpClient.post('/friend/remove_black', data: {
      'ownerUserID': Config.userID,
      'blackUserID': userID,
    });
  }

  Future<List<FriendshipInfo>> checkFriend({
    required List<String> userIDList,
    String? operationID,
  }) async {
    final data = await HttpClient.post('/friend/is_friend', data: {
      'userID1': Config.userID,
      'userID2': userIDList.first,
    });
    if (data == null) return [];
    // Build result
    return userIDList.map((uid) => FriendshipInfo(
      userID: uid,
      result: (data['inUser1Friends'] == true || data['inUser2Friends'] == true) ? 1 : 0,
    )).toList();
  }

  Future<List<SearchFriendsInfo>> searchFriends({
    List<String> keywordList = const [],
    bool isSearchUserID = false,
    bool isSearchNickname = false,
    bool isSearchRemark = false,
    String? operationID,
  }) async {
    final data = await HttpClient.post('/friend/search', data: {
      'keyword': keywordList.isNotEmpty ? keywordList.first : '',
      'pagination': {'pageNumber': 1, 'showNumber': 20},
    });
    if (data == null) return [];
    final list = data['users'] ?? data['friends'] ?? data;
    if (list is List) {
      return list.map((e) => SearchFriendsInfo.fromJson(e)).toList();
    }
    return [];
  }

  // Helpers
  List<FriendInfo> _parseFriendList(dynamic data) {
    if (data == null) return [];
    final list = data is Map ? (data['friendsInfo'] ?? data) : data;
    if (list is List) {
      return list.map((e) => FriendInfo.fromJson(e)).toList();
    }
    return [];
  }

  List<FriendApplicationInfo> _parseFriendApplicationList(dynamic data) {
    if (data == null) return [];
    debugPrint('[FriendshipManager] raw application data: $data');
    final list = data is Map ? (data['friendRequests'] ?? data) : data;
    if (list is List) {
      final result = list.map((e) {
        debugPrint('[FriendshipManager] raw item: $e');
        return FriendApplicationInfo.fromJson(e);
      }).toList();
      for (final r in result) {
        debugPrint('[FriendshipManager] parsed: fromUserID=${r.fromUserID}, fromNickname=${r.fromNickname}, reqMsg=${r.reqMsg}');
      }
      return result;
    }
    return [];
  }

  /// 服务端好友申请列表不含 nickname/faceURL，需要额外查询补全
  Future<void> _enrichApplicationsWithUserInfo(List<FriendApplicationInfo> list) async {
    // 收集所有需要补全的 userID（fromUserID 和 toUserID）
    final needIDs = <String>{};
    for (final apply in list) {
      if ((apply.fromNickname ?? '').isEmpty && (apply.fromUserID ?? '').isNotEmpty) {
        needIDs.add(apply.fromUserID!);
      }
      if ((apply.toNickname ?? '').isEmpty && (apply.toUserID ?? '').isNotEmpty) {
        needIDs.add(apply.toUserID!);
      }
    }
    debugPrint('[FriendshipManager] enrich: needIDs=$needIDs');
    if (needIDs.isEmpty) return;

    try {
      // 直接调 HTTP 接口，避免 getUsersInfo 响应格式不匹配
      final data = await HttpClient.post('/user/get_users_info', data: {
        'userIDs': needIDs.toList(),
      });
      debugPrint('[FriendshipManager] enrich raw response: $data (type: ${data.runtimeType})');

      // 从响应中提取用户列表（兼容多种服务端格式）
      final userMap = <String, Map<String, dynamic>>{};
      List? userList;
      if (data is List) {
        userList = data;
      } else if (data is Map) {
        userList = data['users'] ?? data['usersData'] ?? data['usersInfo'] ?? data['userInfo'];
      }
      debugPrint('[FriendshipManager] enrich userList: $userList');

      if (userList is List) {
        for (final item in userList) {
          if (item is Map) {
            final m = Map<String, dynamic>.from(item);
            // 用户信息可能直接在 item 中，也可能嵌套在 userInfo 字段中
            final info = m.containsKey('nickname') ? m : (m['userInfo'] ?? m);
            if (info is Map) {
              final uid = info['userID']?.toString();
              if (uid != null && uid.isNotEmpty) {
                userMap[uid] = Map<String, dynamic>.from(info);
                debugPrint('[FriendshipManager] enrich found user: $uid => nickname=${info['nickname']}');
              }
            }
          }
        }
      }

      debugPrint('[FriendshipManager] enrich: resolved ${userMap.length} users');

      // 回填昵称和头像
      for (final apply in list) {
        if ((apply.fromNickname ?? '').isEmpty) {
          final info = userMap[apply.fromUserID];
          if (info != null) {
            apply.fromNickname = info['nickname']?.toString();
            apply.fromFaceURL = info['faceURL']?.toString();
          }
        }
        if ((apply.toNickname ?? '').isEmpty) {
          final info = userMap[apply.toUserID];
          if (info != null) {
            apply.toNickname = info['nickname']?.toString();
            apply.toFaceURL = info['faceURL']?.toString();
          }
        }
      }
      debugPrint('[FriendshipManager] enrich done, first item fromNickname=${list.isNotEmpty ? list.first.fromNickname : "N/A"}');
    } catch (e, st) {
      debugPrint('[FriendshipManager] enrichApplicationsWithUserInfo error: $e\n$st');
    }
  }
}
