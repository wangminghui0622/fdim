import '../models/group_info.dart';

class OnGroupListener {
  Function(GroupApplicationInfo info)? onGroupApplicationAccepted;
  Function(GroupApplicationInfo info)? onGroupApplicationAdded;
  Function(GroupApplicationInfo info)? onGroupApplicationDeleted;
  Function(GroupApplicationInfo info)? onGroupApplicationRejected;
  Function(GroupInfo info)? onGroupDismissed;
  Function(GroupInfo info)? onGroupInfoChanged;
  Function(GroupMembersInfo info)? onGroupMemberAdded;
  Function(GroupMembersInfo info)? onGroupMemberDeleted;
  Function(GroupMembersInfo info)? onGroupMemberInfoChanged;
  Function(GroupInfo info)? onJoinedGroupAdded;
  Function(GroupInfo info)? onJoinedGroupDeleted;

  OnGroupListener({
    this.onGroupApplicationAccepted,
    this.onGroupApplicationAdded,
    this.onGroupApplicationDeleted,
    this.onGroupApplicationRejected,
    this.onGroupDismissed,
    this.onGroupInfoChanged,
    this.onGroupMemberAdded,
    this.onGroupMemberDeleted,
    this.onGroupMemberInfoChanged,
    this.onJoinedGroupAdded,
    this.onJoinedGroupDeleted,
  });

  void groupApplicationAccepted(GroupApplicationInfo i) =>
      onGroupApplicationAccepted?.call(i);
  void groupApplicationAdded(GroupApplicationInfo i) =>
      onGroupApplicationAdded?.call(i);
  void groupApplicationDeleted(GroupApplicationInfo i) =>
      onGroupApplicationDeleted?.call(i);
  void groupApplicationRejected(GroupApplicationInfo i) =>
      onGroupApplicationRejected?.call(i);
  void groupDismissed(GroupInfo i) => onGroupDismissed?.call(i);
  void groupInfoChanged(GroupInfo i) => onGroupInfoChanged?.call(i);
  void groupMemberAdded(GroupMembersInfo i) => onGroupMemberAdded?.call(i);
  void groupMemberDeleted(GroupMembersInfo i) => onGroupMemberDeleted?.call(i);
  void groupMemberInfoChanged(GroupMembersInfo i) =>
      onGroupMemberInfoChanged?.call(i);
  void joinedGroupAdded(GroupInfo i) => onJoinedGroupAdded?.call(i);
  void joinedGroupDeleted(GroupInfo i) => onJoinedGroupDeleted?.call(i);
}
