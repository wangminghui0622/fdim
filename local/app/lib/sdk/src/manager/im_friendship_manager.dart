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
    return _parseFriendApplicationList(data);
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
    final list = data is Map ? (data['friendRequests'] ?? data) : data;
    if (list is List) {
      return list.map((e) => FriendApplicationInfo.fromJson(e)).toList();
    }
    return [];
  }
}
