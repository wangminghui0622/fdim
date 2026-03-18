import '../models/user_info.dart';

class OnFriendshipListener {
  Function(BlacklistInfo info)? onBlackAdded;
  Function(BlacklistInfo info)? onBlackDeleted;
  Function(FriendInfo info)? onFriendAdded;
  Function(FriendApplicationInfo info)? onFriendApplicationAccepted;
  Function(FriendApplicationInfo info)? onFriendApplicationAdded;
  Function(FriendApplicationInfo info)? onFriendApplicationDeleted;
  Function(FriendApplicationInfo info)? onFriendApplicationRejected;
  Function(FriendInfo info)? onFriendDeleted;
  Function(FriendInfo info)? onFriendInfoChanged;

  OnFriendshipListener({
    this.onBlackAdded,
    this.onBlackDeleted,
    this.onFriendAdded,
    this.onFriendApplicationAccepted,
    this.onFriendApplicationAdded,
    this.onFriendApplicationDeleted,
    this.onFriendApplicationRejected,
    this.onFriendDeleted,
    this.onFriendInfoChanged,
  });

  void blackAdded(BlacklistInfo info) => onBlackAdded?.call(info);
  void blackDeleted(BlacklistInfo info) => onBlackDeleted?.call(info);
  void friendAdded(FriendInfo info) => onFriendAdded?.call(info);
  void friendApplicationAccepted(FriendApplicationInfo info) =>
      onFriendApplicationAccepted?.call(info);
  void friendApplicationAdded(FriendApplicationInfo info) =>
      onFriendApplicationAdded?.call(info);
  void friendApplicationDeleted(FriendApplicationInfo info) =>
      onFriendApplicationDeleted?.call(info);
  void friendApplicationRejected(FriendApplicationInfo info) =>
      onFriendApplicationRejected?.call(info);
  void friendDeleted(FriendInfo info) => onFriendDeleted?.call(info);
  void friendInfoChanged(FriendInfo info) => onFriendInfoChanged?.call(info);
}
