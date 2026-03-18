import '../models/user_info.dart';
import '../models/group_info.dart';
import '../models/message.dart';

/// Listener for service-level events
class OnListenerForService {
  Function(FriendApplicationInfo i)? onFriendApplicationAdded;
  Function(FriendApplicationInfo i)? onFriendApplicationAccepted;
  Function(GroupApplicationInfo info)? onGroupApplicationAccepted;
  Function(GroupApplicationInfo info)? onGroupApplicationAdded;
  Function(Message msg)? onRecvNewMessage;

  OnListenerForService({
    this.onFriendApplicationAdded,
    this.onFriendApplicationAccepted,
    this.onGroupApplicationAccepted,
    this.onGroupApplicationAdded,
    this.onRecvNewMessage,
  });

  void friendApplicationAccepted(FriendApplicationInfo u) {
    onFriendApplicationAccepted?.call(u);
  }

  void friendApplicationAdded(FriendApplicationInfo u) {
    onFriendApplicationAdded?.call(u);
  }

  void groupApplicationAccepted(GroupApplicationInfo info) {
    onGroupApplicationAccepted?.call(info);
  }

  void groupApplicationAdded(GroupApplicationInfo info) {
    onGroupApplicationAdded?.call(info);
  }

  void recvNewMessage(Message msg) {
    onRecvNewMessage?.call(msg);
  }
}
