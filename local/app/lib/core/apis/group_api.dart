import '../http_client.dart';
import '../config.dart';
import '../models/group_info.dart';

class GroupApi {
  static Future<GroupInfo?> getGroupsInfo(String groupID) async {
    final data = await HttpClient.post('/group/get_groups_info', data: {
      'groupIDs': [groupID],
    });
    if (data == null) return null;
    final list = data['groupInfos'] ?? data;
    if (list is List && list.isNotEmpty) {
      return GroupInfo.fromJson(list[0]);
    }
    return null;
  }

  static Future<List<GroupInfo>> getJoinedGroupList({
    int pageNumber = 1,
    int showNumber = 100,
  }) async {
    final data = await HttpClient.post('/group/get_joined_group_list', data: {
      'fromUserID': Config.userID,
      'pagination': {'pageNumber': pageNumber, 'showNumber': showNumber},
    });
    if (data == null) return [];
    final list = data['groups'] ?? data;
    if (list is List) {
      return list.map((e) => GroupInfo.fromJson(e)).toList();
    }
    return [];
  }

  static Future<List<GroupMemberInfo>> getGroupMemberList(
    String groupID, {
    int pageNumber = 1,
    int showNumber = 100,
  }) async {
    final data = await HttpClient.post('/group/get_group_member_list', data: {
      'groupID': groupID,
      'pagination': {'pageNumber': pageNumber, 'showNumber': showNumber},
    });
    if (data == null) return [];
    final list = data['members'] ?? data;
    if (list is List) {
      return list.map((e) => GroupMemberInfo.fromJson(e)).toList();
    }
    return [];
  }

  static Future<GroupInfo?> createGroup({
    required String groupName,
    required List<String> memberUserIDs,
    String faceURL = '',
  }) async {
    final data = await HttpClient.post('/group/create_group', data: {
      'ownerUserID': Config.userID,
      'groupInfo': {
        'groupName': groupName,
        'faceURL': faceURL,
        'groupType': 2,
      },
      'memberUserIDs': memberUserIDs,
    });
    if (data == null) return null;
    final groupData = data['groupInfo'] ?? data;
    if (groupData is Map<String, dynamic>) {
      return GroupInfo.fromJson(groupData);
    }
    return null;
  }

  static Future<void> joinGroup(String groupID, {String reqMsg = ''}) async {
    await HttpClient.post('/group/join_group', data: {
      'groupID': groupID,
      'userID': Config.userID,
      'reqMessage': reqMsg,
      'joinSource': 2,
    });
  }

  static Future<void> quitGroup(String groupID) async {
    await HttpClient.post('/group/quit_group', data: {
      'groupID': groupID,
      'userID': Config.userID,
    });
  }

  static Future<void> inviteUserToGroup(
      String groupID, List<String> userIDs) async {
    await HttpClient.post('/group/invite_user_to_group', data: {
      'groupID': groupID,
      'invitedUserIDs': userIDs,
    });
  }

  static Future<List<GroupApplicationInfo>> getRecvGroupApplicationList({
    int pageNumber = 1,
    int showNumber = 100,
  }) async {
    final data = await HttpClient.post('/group/get_recv_group_applicationList', data: {
      'userID': Config.userID,
      'pagination': {'pageNumber': pageNumber, 'showNumber': showNumber},
    });
    if (data == null) return [];
    final list = data['groupRequests'] ?? data;
    if (list is List) {
      return list.map((e) => GroupApplicationInfo.fromJson(e)).toList();
    }
    return [];
  }

  static Future<void> respondGroupApplication({
    required String groupID,
    required String fromUserID,
    required int handleResult, // 1: agree, -1: reject
    String handleMsg = '',
  }) async {
    await HttpClient.post('/group/group_application_response', data: {
      'groupID': groupID,
      'fromUserID': fromUserID,
      'handledMsg': handleMsg,
      'handleResult': handleResult,
    });
  }

  static Future<void> setGroupInfo({
    required String groupID,
    String? groupName,
    String? notification,
    String? introduction,
    String? faceURL,
    int? needVerification,
  }) async {
    final groupInfo = <String, dynamic>{'groupID': groupID};
    if (groupName != null) groupInfo['groupName'] = groupName;
    if (notification != null) groupInfo['notification'] = notification;
    if (introduction != null) groupInfo['introduction'] = introduction;
    if (faceURL != null) groupInfo['faceURL'] = faceURL;
    if (needVerification != null) groupInfo['needVerification'] = needVerification;
    await HttpClient.post('/group/set_group_info', data: {
      'groupInfoForSet': groupInfo,
    });
  }

  static Future<void> kickGroupMember(
      String groupID, List<String> userIDs) async {
    await HttpClient.post('/group/kick_group', data: {
      'groupID': groupID,
      'kickedUserIDs': userIDs,
    });
  }

  static Future<void> transferGroupOwner(
      String groupID, String newOwnerUserID) async {
    await HttpClient.post('/group/transfer_group', data: {
      'groupID': groupID,
      'oldOwnerUserID': Config.userID,
      'newOwnerUserID': newOwnerUserID,
    });
  }

  static Future<void> dismissGroup(String groupID) async {
    await HttpClient.post('/group/dismiss_group', data: {
      'groupID': groupID,
    });
  }

  static Future<void> setGroupMemberInfo({
    required String groupID,
    required String userID,
    String? nickname,
    String? faceURL,
    int? roleLevel,
    String? ex,
  }) async {
    final info = <String, dynamic>{'groupID': groupID, 'userID': userID};
    if (nickname != null) info['nickname'] = nickname;
    if (faceURL != null) info['faceURL'] = faceURL;
    if (roleLevel != null) info['roleLevel'] = roleLevel;
    if (ex != null) info['ex'] = ex;
    await HttpClient.post('/group/set_group_member_info', data: info);
  }
}
