import '../models/user_info.dart';

class OnUserListener {
  Function(UserInfo info)? onSelfInfoUpdated;
  Function(UserStatusInfo info)? onUserStatusChanged;

  OnUserListener({this.onSelfInfoUpdated, this.onUserStatusChanged});

  void selfInfoUpdated(UserInfo info) => onSelfInfoUpdated?.call(info);
  void userStatusChanged(UserStatusInfo info) =>
      onUserStatusChanged?.call(info);
}
