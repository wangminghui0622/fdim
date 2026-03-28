import '../../flutter_openim_sdk.dart';
import '../../../core/http_client.dart';
import '../../../core/config.dart';

class GroupManager {
  late OnGroupListener listener;

  Future setGroupListener(OnGroupListener listener) async {
    this.listener = listener;
  }

  Future<List<GroupInfo>> getJoinedGroupList({String? operationID}) async {
    final data = await HttpClient.post('/group/get_joined_group_list', data: {
      'fromUserID': Config.userID,
      'pagination': {'pageNumber': 1, 'showNumber': 1000},
    });
    if (data == null) return [];
    final list = data['groups'] ?? data;
    if (list is List) {
      return list.map((e) => GroupInfo.fromJson(e)).toList();
    }
    return [];
  }

  Future<List<GroupInfo>> getJoinedGroupListPage({
    int offset = 0,
    int count = 40,
    String? operationID,
  }) async {
    final data = await HttpClient.post('/group/get_joined_group_list', data: {
      'fromUserID': Config.userID,
      'pagination': {
        'pageNumber': offset ~/ count + 1,
        'showNumber': count,
      },
    });
    if (data == null) return [];
    final list = data['groups'] ?? data;
    if (list is List) {
      return list.map((e) => GroupInfo.fromJson(e)).toList();
    }
    return [];
  }

  Future<List<GroupInfo>> getGroupsInfo({
    required List<String> groupIDList,
    String? operationID,
  }) async {
    final data = await HttpClient.post('/group/get_groups_info', data: {
      'groupIDs': groupIDList,
    });
    if (data == null) return [];
    final list = data['groupInfos'] ?? data['groups'] ?? data;
    if (list is List) {
      return list.map((e) => GroupInfo.fromJson(e)).toList();
    }
    return [];
  }

  Future createGroup({
    required GroupInfo groupInfo,
    List<String> memberUserIDs = const [],
    List<String> adminUserIDs = const [],
    String? ownerUserID,
    String? operationID,
  }) async {
    await HttpClient.post('/group/create_group', data: {
      'groupInfo': groupInfo.toJson(),
      'memberUserIDs': memberUserIDs,
      'adminUserIDs': adminUserIDs,
      'ownerUserID': ownerUserID ?? Config.userID,
    });
  }

  Future joinGroup({
    required String groupID,
    String? reason,
    int joinSource = 2,
    String? operationID,
  }) async {
    await HttpClient.post('/group/join_group', data: {
      'groupID': groupID,
      'reqMsg': reason ?? '',
      'joinSource': joinSource,
      'inviterUserID': '',
    });
  }

  Future quitGroup({
    required String groupID,
    String? operationID,
  }) async {
    await HttpClient.post('/group/quit_group', data: {
      'groupID': groupID,
    });
  }

  Future<List<GroupMembersInfo>> getGroupMemberList({
    required String groupID,
    int offset = 0,
    int count = 40,
    int filter = 0,
    String? operationID,
  }) async {
    final data = await HttpClient.post('/group/get_group_member_list', data: {
      'groupID': groupID,
      'pagination': {
        'pageNumber': offset ~/ count + 1,
        'showNumber': count,
      },
      'filter': filter,
    });
    if (data == null) return [];
    final list = data['members'] ?? data;
    if (list is List) {
      return list.map((e) => GroupMembersInfo.fromJson(e)).toList();
    }
    return [];
  }

  Future<List<GroupApplicationInfo>> getGroupApplicationListAsRecipient({
    String? operationID,
  }) async {
    final data = await HttpClient.post('/group/get_recv_group_applicationList', data: {
      'fromUserID': Config.userID,
      'pagination': {'pageNumber': 1, 'showNumber': 100},
    });
    if (data == null) return [];
    final list = data['groupRequests'] ?? data;
    if (list is List) {
      return list.map((e) => GroupApplicationInfo.fromJson(e)).toList();
    }
    return [];
  }

  Future acceptGroupApplication({
    required String groupID,
    required String fromUserID,
    String? handleMsg,
    String? operationID,
  }) async {
    await HttpClient.post('/group/group_application_response', data: {
      'groupID': groupID,
      'fromUserID': fromUserID,
      'handledMsg': handleMsg ?? '',
      'handleResult': 1,
    });
  }

  Future refuseGroupApplication({
    required String groupID,
    required String fromUserID,
    String? handleMsg,
    String? operationID,
  }) async {
    await HttpClient.post('/group/group_application_response', data: {
      'groupID': groupID,
      'fromUserID': fromUserID,
      'handledMsg': handleMsg ?? '',
      'handleResult': -1,
    });
  }

  Future dismissGroup({
    required String groupID,
    String? operationID,
  }) async {
    await HttpClient.post('/group/dismiss_group', data: {
      'groupID': groupID,
    });
  }

  Future setGroupInfo({
    required GroupInfo groupInfo,
    String? operationID,
  }) async {
    await HttpClient.post('/group/set_group_info', data: {
      'groupInfoForSet': groupInfo.toJson(),
    });
  }

  Future kickGroupMember({
    required String groupID,
    required List<String> userIDList,
    String? reason,
    String? operationID,
  }) async {
    await HttpClient.post('/group/kick_group', data: {
      'groupID': groupID,
      'kickedUserIDs': userIDList,
      'reason': reason ?? '',
    });
  }

  Future inviteUserToGroup({
    required String groupID,
    required List<String> userIDList,
    String? reason,
    String? operationID,
  }) async {
    await HttpClient.post('/group/invite_user_to_group', data: {
      'groupID': groupID,
      'invitedUserIDs': userIDList,
      'reason': reason ?? '',
    });
  }

  Future transferGroupOwner({
    required String groupID,
    required String newOwnerUserID,
    String? operationID,
  }) async {
    await HttpClient.post('/group/transfer_group', data: {
      'groupID': groupID,
      'newOwnerUserID': newOwnerUserID,
    });
  }
}
