import '../../flutter_openim_sdk.dart';
import '../../../core/http_client.dart';
import '../../../core/config.dart';

class UserManager {
  late OnUserListener listener;

  Future setUserListener(OnUserListener listener) async {
    this.listener = listener;
  }

  Future<List<PublicUserInfo>> getUsersInfo({
    required List<String> userIDList,
    String? operationID,
  }) async {
    final data = await HttpClient.post('/user/get_users_info', data: {
      'userIDs': userIDList,
    });
    if (data is List) {
      return data.map((e) => PublicUserInfo.fromJson(e)).toList();
    }
    if (data is Map && data['users'] is List) {
      return (data['users'] as List)
          .map((e) => PublicUserInfo.fromJson(e))
          .toList();
    }
    return [];
  }

  Future<UserInfo> getSelfUserInfo({String? operationID}) async {
    final list = await getUsersInfo(userIDList: [Config.userID]);
    if (list.isNotEmpty) {
      final pub = list.first;
      return UserInfo(
        userID: pub.userID,
        nickname: pub.nickname,
        faceURL: pub.faceURL,
        ex: pub.ex,
      );
    }
    return UserInfo(userID: Config.userID, nickname: Config.nickname);
  }

  Future<String?> setSelfInfo({
    String? nickname,
    String? faceURL,
    int? gender,
    String? phoneNumber,
    String? email,
    int? birth,
    int? globalRecvMsgOpt,
    String? ex,
    String? operationID,
  }) async {
    final info = <String, dynamic>{
      'userID': Config.userID,
    };
    if (nickname != null) info['nickname'] = nickname;
    if (faceURL != null) info['faceURL'] = faceURL;
    if (gender != null) info['gender'] = gender;
    if (phoneNumber != null) info['phoneNumber'] = phoneNumber;
    if (email != null) info['email'] = email;
    if (birth != null) info['birth'] = birth;
    if (ex != null) info['ex'] = ex;

    await HttpClient.post('/user/update_user_info', data: info);

    // Update local config
    if (nickname != null) Config.nickname = nickname;
    if (faceURL != null) Config.faceURL = faceURL;

    return null;
  }

  Future<List<UserStatusInfo>> getUserStatus(
    List<String> userIDs, {
    String? operationID,
  }) async {
    final data = await HttpClient.post(
      '/user/get_users_online_status',
      data: {'userIDs': userIDs},
      showErrorToast: false,
    );
    if (data == null) return [];
    final result = data['successResult'];
    if (result is List) {
      return result.map((e) => UserStatusInfo.fromJson(e)).toList();
    }
    return [];
  }

  Future<List<UserStatusInfo>> subscribeUsersStatus(
    List<String> userIDs, {
    String? operationID,
  }) async {
    return getUserStatus(userIDs);
  }

  Future unsubscribeUsersStatus(
    List<String> userIDs, {
    String? operationID,
  }) async {}

  Future<List<UserStatusInfo>> getSubscribeUsersStatus({
    String? operationID,
  }) async {
    return [];
  }
}
