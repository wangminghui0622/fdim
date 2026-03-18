class GroupInfo {
  String groupID;
  String groupName;
  String notification;
  String introduction;
  String faceURL;
  String ownerUserID;
  int createTime;
  int memberCount;
  int status;
  String creatorUserID;
  int groupType;
  int needVerification;
  String ex;

  GroupInfo({
    this.groupID = '',
    this.groupName = '',
    this.notification = '',
    this.introduction = '',
    this.faceURL = '',
    this.ownerUserID = '',
    this.createTime = 0,
    this.memberCount = 0,
    this.status = 0,
    this.creatorUserID = '',
    this.groupType = 0,
    this.needVerification = 0,
    this.ex = '',
  });

  factory GroupInfo.fromJson(Map<String, dynamic> json) => GroupInfo(
        groupID: json['groupID'] ?? '',
        groupName: json['groupName'] ?? '',
        notification: json['notification'] ?? '',
        introduction: json['introduction'] ?? '',
        faceURL: json['faceURL'] ?? '',
        ownerUserID: json['ownerUserID'] ?? '',
        createTime: json['createTime'] ?? 0,
        memberCount: json['memberCount'] ?? 0,
        status: json['status'] ?? 0,
        creatorUserID: json['creatorUserID'] ?? '',
        groupType: json['groupType'] ?? 0,
        needVerification: json['needVerification'] ?? 0,
        ex: json['ex'] ?? '',
      );
}

class GroupMemberInfo {
  String groupID;
  String userID;
  String nickname;
  String faceURL;
  int roleLevel; // 1:普通 2:群主 3:管理员
  int joinTime;
  int joinSource;
  String operatorUserID;
  String ex;
  int muteEndTime;

  GroupMemberInfo({
    this.groupID = '',
    this.userID = '',
    this.nickname = '',
    this.faceURL = '',
    this.roleLevel = 1,
    this.joinTime = 0,
    this.joinSource = 0,
    this.operatorUserID = '',
    this.ex = '',
    this.muteEndTime = 0,
  });

  factory GroupMemberInfo.fromJson(Map<String, dynamic> json) =>
      GroupMemberInfo(
        groupID: json['groupID'] ?? '',
        userID: json['userID'] ?? '',
        nickname: json['nickname'] ?? '',
        faceURL: json['faceURL'] ?? '',
        roleLevel: json['roleLevel'] ?? 1,
        joinTime: json['joinTime'] ?? 0,
        joinSource: json['joinSource'] ?? 0,
        operatorUserID: json['operatorUserID'] ?? '',
        ex: json['ex'] ?? '',
        muteEndTime: json['muteEndTime'] ?? 0,
      );

  bool get isOwner => roleLevel == 100;
  bool get isAdmin => roleLevel == 60;
}

class GroupApplicationInfo {
  String groupID;
  String groupName;
  String groupFaceURL;
  String userID;
  String nickname;
  String userFaceURL;
  String reqMsg;
  int reqTime;
  int handleResult; // 0: pending, 1: agreed, -1: rejected
  String handledMsg;
  String handleUserID;
  int joinSource;

  GroupApplicationInfo({
    this.groupID = '',
    this.groupName = '',
    this.groupFaceURL = '',
    this.userID = '',
    this.nickname = '',
    this.userFaceURL = '',
    this.reqMsg = '',
    this.reqTime = 0,
    this.handleResult = 0,
    this.handledMsg = '',
    this.handleUserID = '',
    this.joinSource = 0,
  });

  factory GroupApplicationInfo.fromJson(Map<String, dynamic> json) =>
      GroupApplicationInfo(
        groupID: json['groupID'] ?? '',
        groupName: json['groupName'] ?? '',
        groupFaceURL: json['groupFaceURL'] ?? '',
        userID: json['userID'] ?? '',
        nickname: json['nickname'] ?? '',
        userFaceURL: json['userFaceURL'] ?? '',
        reqMsg: json['reqMsg'] ?? '',
        reqTime: json['reqTime'] ?? 0,
        handleResult: json['handleResult'] ?? 0,
        handledMsg: json['handledMsg'] ?? '',
        handleUserID: json['handleUserID'] ?? '',
        joinSource: json['joinSource'] ?? 0,
      );

  bool get isPending => handleResult == 0;
}
