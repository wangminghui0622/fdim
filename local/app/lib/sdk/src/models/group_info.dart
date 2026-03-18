class GroupInfo {
  String? groupID;
  String? groupName;
  String? notification;
  String? introduction;
  String? faceURL;
  String? ownerUserID;
  int? createTime;
  int? memberCount;
  int? status;
  String? creatorUserID;
  int? groupType;
  int? needVerification;
  int? lookMemberInfo;
  int? applyMemberFriend;
  int? notificationUpdateTime;
  String? notificationUserID;
  String? ex;

  GroupInfo({
    this.groupID, this.groupName, this.notification, this.introduction,
    this.faceURL, this.ownerUserID, this.createTime, this.memberCount,
    this.status, this.creatorUserID, this.groupType, this.needVerification,
    this.lookMemberInfo, this.applyMemberFriend, this.notificationUpdateTime,
    this.notificationUserID, this.ex,
  });

  GroupInfo.fromJson(Map<String, dynamic> json) {
    groupID = json['groupID']; groupName = json['groupName'];
    notification = json['notification']; introduction = json['introduction'];
    faceURL = json['faceURL']; ownerUserID = json['ownerUserID'];
    createTime = json['createTime']; memberCount = json['memberCount'];
    status = json['status']; creatorUserID = json['creatorUserID'];
    groupType = json['groupType']; needVerification = json['needVerification'];
    lookMemberInfo = json['lookMemberInfo']; applyMemberFriend = json['applyMemberFriend'];
    notificationUpdateTime = json['notificationUpdateTime'];
    notificationUserID = json['notificationUserID']; ex = json['ex'];
  }

  Map<String, dynamic> toJson() => {
    'groupID': groupID, 'groupName': groupName, 'notification': notification,
    'introduction': introduction, 'faceURL': faceURL, 'ownerUserID': ownerUserID,
    'createTime': createTime, 'memberCount': memberCount, 'status': status,
    'creatorUserID': creatorUserID, 'groupType': groupType,
    'needVerification': needVerification, 'lookMemberInfo': lookMemberInfo,
    'applyMemberFriend': applyMemberFriend, 'notificationUpdateTime': notificationUpdateTime,
    'notificationUserID': notificationUserID, 'ex': ex,
  };
}

class GroupMembersInfo {
  String? groupID;
  String? userID;
  String? nickname;
  String? faceURL;
  int? roleLevel;
  int? joinTime;
  int? joinSource;
  String? operatorUserID;
  String? ex;
  int? muteEndTime;
  String? inviterUserID;

  GroupMembersInfo({
    this.groupID, this.userID, this.nickname, this.faceURL, this.roleLevel,
    this.joinTime, this.joinSource, this.operatorUserID, this.ex,
    this.muteEndTime, this.inviterUserID,
  });

  GroupMembersInfo.fromJson(Map<String, dynamic> json) {
    groupID = json['groupID']; userID = json['userID']; nickname = json['nickname'];
    faceURL = json['faceURL']; roleLevel = json['roleLevel']; joinTime = json['joinTime'];
    joinSource = json['joinSource']; operatorUserID = json['operatorUserID']; ex = json['ex'];
    muteEndTime = json['muteEndTime']; inviterUserID = json['inviterUserID'];
  }

  Map<String, dynamic> toJson() => {
    'groupID': groupID, 'userID': userID, 'nickname': nickname, 'faceURL': faceURL,
    'roleLevel': roleLevel, 'joinTime': joinTime, 'joinSource': joinSource,
    'operatorUserID': operatorUserID, 'ex': ex, 'muteEndTime': muteEndTime, 'inviterUserID': inviterUserID,
  };
}

class GroupApplicationInfo {
  String? groupID;
  String? groupName;
  String? groupFaceURL;
  int? creatorUserID;
  String? ownerUserID;
  int? memberCount;
  String? userID;
  String? nickname;
  String? userFaceURL;
  int? gender;
  int? handleResult;
  String? reqMsg;
  String? handledMsg;
  int? reqTime;
  String? handleUserID;
  int? handledTime;
  String? ex;
  int? joinSource;
  String? inviterUserID;

  GroupApplicationInfo.fromJson(Map<String, dynamic> json) {
    groupID = json['groupID']; groupName = json['groupName']; groupFaceURL = json['groupFaceURL'];
    userID = json['userID']; nickname = json['nickname']; userFaceURL = json['userFaceURL'];
    gender = json['gender']; handleResult = json['handleResult']; reqMsg = json['reqMsg'];
    handledMsg = json['handledMsg']; reqTime = json['reqTime']; handleUserID = json['handleUserID'];
    handledTime = json['handledTime']; ex = json['ex']; joinSource = json['joinSource'];
    inviterUserID = json['inviterUserID']; memberCount = json['memberCount'];
  }

  Map<String, dynamic> toJson() => {
    'groupID': groupID, 'groupName': groupName, 'groupFaceURL': groupFaceURL,
    'userID': userID, 'nickname': nickname, 'userFaceURL': userFaceURL,
    'gender': gender, 'handleResult': handleResult, 'reqMsg': reqMsg,
    'handledMsg': handledMsg, 'reqTime': reqTime, 'handleUserID': handleUserID,
    'handledTime': handledTime, 'ex': ex, 'joinSource': joinSource,
    'inviterUserID': inviterUserID, 'memberCount': memberCount,
  };
}
