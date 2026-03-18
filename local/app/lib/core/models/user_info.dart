class UserInfo {
  String userID;
  String nickname;
  String faceURL;
  int gender;
  String phoneNumber;
  String email;
  String birth;
  int createTime;
  String ex;
  int globalRecvMsgOpt;

  UserInfo({
    this.userID = '',
    this.nickname = '',
    this.faceURL = '',
    this.gender = 0,
    this.phoneNumber = '',
    this.email = '',
    this.birth = '',
    this.createTime = 0,
    this.ex = '',
    this.globalRecvMsgOpt = 0,
  });

  factory UserInfo.fromJson(Map<String, dynamic> json) => UserInfo(
        userID: json['userID'] ?? '',
        nickname: json['nickname'] ?? '',
        faceURL: json['faceURL'] ?? '',
        gender: json['gender'] ?? 0,
        phoneNumber: json['phoneNumber'] ?? '',
        email: json['email'] ?? '',
        birth: json['birth'] ?? '',
        createTime: json['createTime'] ?? 0,
        ex: json['ex'] ?? '',
        globalRecvMsgOpt: json['globalRecvMsgOpt'] ?? 0,
      );

  Map<String, dynamic> toJson() => {
        'userID': userID,
        'nickname': nickname,
        'faceURL': faceURL,
        'gender': gender,
        'phoneNumber': phoneNumber,
        'email': email,
        'birth': birth,
        'createTime': createTime,
        'ex': ex,
        'globalRecvMsgOpt': globalRecvMsgOpt,
      };
}

class FriendInfo {
  String ownerUserID;
  String friendUserID;
  String nickname;
  String faceURL;
  String remark;
  int createTime;
  String ex;
  bool isOnline;

  FriendInfo({
    this.ownerUserID = '',
    this.friendUserID = '',
    this.nickname = '',
    this.faceURL = '',
    this.remark = '',
    this.createTime = 0,
    this.ex = '',
    this.isOnline = false,
  });

  factory FriendInfo.fromJson(Map<String, dynamic> json) => FriendInfo(
        ownerUserID: json['ownerUserID'] ?? '',
        friendUserID: json['friendUserID'] ?? '',
        nickname: json['nickname'] ?? '',
        faceURL: json['faceURL'] ?? '',
        remark: json['remark'] ?? '',
        createTime: json['createTime'] ?? 0,
        ex: json['ex'] ?? '',
      );

  String get showName => remark.isNotEmpty ? remark : nickname;
}

class FriendApplicationInfo {
  String fromUserID;
  String fromNickname;
  String fromFaceURL;
  String toUserID;
  String toNickname;
  String toFaceURL;
  int handleResult; // 0:未处理 1:同意 -1:拒绝
  String reqMsg;
  String handlerUserID;
  String handleMsg;
  int handleTime;
  int createTime;

  FriendApplicationInfo({
    this.fromUserID = '',
    this.fromNickname = '',
    this.fromFaceURL = '',
    this.toUserID = '',
    this.toNickname = '',
    this.toFaceURL = '',
    this.handleResult = 0,
    this.reqMsg = '',
    this.handlerUserID = '',
    this.handleMsg = '',
    this.handleTime = 0,
    this.createTime = 0,
  });

  factory FriendApplicationInfo.fromJson(Map<String, dynamic> json) =>
      FriendApplicationInfo(
        fromUserID: json['fromUserID'] ?? '',
        fromNickname: json['fromNickname'] ?? '',
        fromFaceURL: json['fromFaceURL'] ?? '',
        toUserID: json['toUserID'] ?? '',
        toNickname: json['toNickname'] ?? '',
        toFaceURL: json['toFaceURL'] ?? '',
        handleResult: json['handleResult'] ?? 0,
        reqMsg: json['reqMsg'] ?? '',
        handlerUserID: json['handlerUserID'] ?? '',
        handleMsg: json['handleMsg'] ?? '',
        handleTime: json['handleTime'] ?? 0,
        createTime: json['createTime'] ?? 0,
      );
}
