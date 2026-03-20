class UserInfo {
  String? userID;
  String? nickname;
  String? faceURL;
  String? ex;
  int? createTime;
  String? remark;
  int? globalRecvMsgOpt;
  int? appMangerLevel;
  int? gender;
  String? phoneNumber;
  String? email;
  int? birth;

  UserInfo({
    this.userID,
    this.nickname,
    this.faceURL,
    this.appMangerLevel,
    this.ex,
    this.createTime,
    this.remark,
    this.globalRecvMsgOpt,
    this.gender,
    this.phoneNumber,
    this.email,
    this.birth,
  });

  UserInfo.fromJson(Map<String, dynamic> json) {
    userID = json['userID'] ?? userID;
    nickname = json['nickname'] ?? nickname;
    faceURL = json['faceURL'] ?? faceURL;
    remark = json['remark'] ?? remark;
    ex = json['ex'] ?? ex;
    createTime = json['createTime'];
    globalRecvMsgOpt = json['globalRecvMsgOpt'];
    appMangerLevel = json['appMangerLevel'];
    gender = json['gender'];
    phoneNumber = json['phoneNumber'];
    email = json['email'];
    birth = json['birth'];
  }

  Map<String, dynamic> toJson() => {
        'appMangerLevel': appMangerLevel,
        'userID': userID,
        'nickname': nickname,
        'faceURL': faceURL,
        'ex': ex,
        'createTime': createTime,
        'remark': remark,
        'globalRecvMsgOpt': globalRecvMsgOpt,
        'gender': gender,
        'phoneNumber': phoneNumber,
        'email': email,
        'birth': birth,
      };

  String getShowName() => _isNull(remark) ?? _isNull(nickname) ?? userID!;

  static String? _isNull(String? value) {
    if (value == null || value.trim().isEmpty) return null;
    return value;
  }

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is UserInfo &&
          runtimeType == other.runtimeType &&
          userID == other.userID;

  @override
  int get hashCode => userID.hashCode;
}

class PublicUserInfo {
  String? userID;
  String? nickname;
  String? faceURL;
  int? appManagerLevel;
  String? ex;

  PublicUserInfo({
    this.userID,
    this.nickname,
    this.faceURL,
    this.appManagerLevel,
    this.ex,
  });

  PublicUserInfo.fromJson(Map<String, dynamic> json) {
    userID = json['userID'];
    nickname = json['nickname'];
    faceURL = json['faceURL'];
    appManagerLevel = json['appManagerLevel'];
    ex = json['ex'];
  }

  Map<String, dynamic> toJson() => {
        'userID': userID,
        'nickname': nickname,
        'faceURL': faceURL,
        'appMangerLevel': appManagerLevel,
        'ex': ex,
      };
}

class FriendInfo {
  String? ownerUserID;
  String? userID;
  String? nickname;
  String? faceURL;
  String? friendUserID;
  String? remark;
  String? ex;
  int? createTime;
  int? addSource;
  String? operatorUserID;
  bool isOnline = false;

  FriendInfo({
    this.ownerUserID,
    this.userID,
    this.nickname,
    this.faceURL,
    this.friendUserID,
    this.remark,
    this.ex,
    this.createTime,
    this.addSource,
    this.operatorUserID,
    this.isOnline = false,
  });

  String get showName => _isNull(remark) ?? _isNull(nickname) ?? (friendUserID ?? userID ?? '');

  FriendInfo.fromJson(Map<String, dynamic> json) {
    ownerUserID = json['ownerUserID'];
    userID = json['userID'];
    remark = json['remark'];
    createTime = json['createTime'];
    addSource = json['addSource'];
    operatorUserID = json['operatorUserID'];
    nickname = json['nickname'];
    faceURL = json['faceURL'];
    friendUserID = json['friendUserID'];
    ex = json['ex'];
  }

  Map<String, dynamic> toJson() => {
        'ownerUserID': ownerUserID,
        'userID': userID,
        'remark': remark,
        'createTime': createTime,
        'addSource': addSource,
        'operatorUserID': operatorUserID,
        'nickname': nickname,
        'faceURL': faceURL,
        'friendUserID': friendUserID,
        'ex': ex,
      };

  String getShowName() => showName;

  static String? _isNull(String? value) {
    if (value == null || value.trim().isEmpty) return null;
    return value;
  }
}

class BlacklistInfo {
  String? userID;
  String? nickname;
  String? ownerUserID;
  String? blockUserID;
  String? faceURL;
  int? gender;
  int? createTime;
  int? addSource;
  String? operatorUserID;
  String? ex;

  BlacklistInfo({
    this.ownerUserID,
    this.blockUserID,
    this.userID,
    this.nickname,
    this.faceURL,
    this.gender,
    this.createTime,
    this.addSource,
    this.operatorUserID,
    this.ex,
  });

  BlacklistInfo.fromJson(Map<String, dynamic> json) {
    ownerUserID = json['ownerUserID'];
    blockUserID = json['blockUserID'];
    userID = json['userID'];
    nickname = json['nickname'];
    faceURL = json['faceURL'];
    gender = json['gender'];
    createTime = json['createTime'];
    addSource = json['addSource'];
    operatorUserID = json['operatorUserID'];
    ex = json['ex'];
  }

  Map<String, dynamic> toJson() => {
        'ownerUserID': ownerUserID,
        'blockUserID': blockUserID,
        'userID': userID,
        'nickname': nickname,
        'faceURL': faceURL,
        'gender': gender,
        'createTime': createTime,
        'addSource': addSource,
        'operatorUserID': operatorUserID,
        'ex': ex,
      };
}

class FriendshipInfo {
  String? userID;
  int? result;

  FriendshipInfo({this.userID, this.result});

  FriendshipInfo.fromJson(Map<String, dynamic> json) {
    userID = json['userID'];
    result = json['result'];
  }

  Map<String, dynamic> toJson() => {'userID': userID, 'result': result};
}

class FriendApplicationInfo {
  String? fromUserID;
  String? fromNickname;
  String? fromFaceURL;
  String? toUserID;
  String? toNickname;
  String? toFaceURL;
  int? handleResult;
  String? reqMsg;
  int? createTime;
  String? handlerUserID;
  String? handleMsg;
  int? handleTime;
  String? ex;

  FriendApplicationInfo({
    this.fromUserID,
    this.fromNickname,
    this.fromFaceURL,
    this.toUserID,
    this.toNickname,
    this.toFaceURL,
    this.handleResult,
    this.reqMsg,
    this.createTime,
    this.handlerUserID,
    this.handleMsg,
    this.handleTime,
    this.ex,
  });

  FriendApplicationInfo.fromJson(Map<String, dynamic> json) {
    fromUserID = json['fromUserID'];
    fromNickname = json['fromNickname'];
    fromFaceURL = json['fromFaceURL'];
    toUserID = json['toUserID'];
    toNickname = json['toNickname'];
    toFaceURL = json['toFaceURL'];
    handleResult = json['handleResult'];
    reqMsg = json['reqMsg'];
    createTime = json['createTime'];
    handlerUserID = json['handlerUserID'];
    handleMsg = json['handleMsg'];
    handleTime = json['handleTime'];
    ex = json['ex'];
  }

  Map<String, dynamic> toJson() => {
        'fromUserID': fromUserID,
        'fromNickname': fromNickname,
        'fromFaceURL': fromFaceURL,
        'toUserID': toUserID,
        'toNickname': toNickname,
        'toFaceURL': toFaceURL,
        'handleResult': handleResult,
        'reqMsg': reqMsg,
        'createTime': createTime,
        'handlerUserID': handlerUserID,
        'handleMsg': handleMsg,
        'handleTime': handleTime,
        'ex': ex,
      };

  bool get isWaitingHandle => handleResult == 0;
  bool get isAgreed => handleResult == 1;
  bool get isRejected => handleResult == -1;
}

class UserStatusInfo {
  String? userID;
  int? status;
  List<int>? platformIDs;

  UserStatusInfo({this.userID, this.status, this.platformIDs});

  bool get isOnline => status == 1;

  UserStatusInfo.fromJson(Map<String, dynamic> json) {
    userID = json['userID'];
    status = json['status'];
    platformIDs = json['platformIDs'] == null
        ? []
        : List<int>.from(json['platformIDs'].map((x) => x));
  }

  Map<String, dynamic> toJson() => {
        'userID': userID,
        'status': status,
        'platformIDs':
            platformIDs != null ? List<dynamic>.from(platformIDs!.map((x) => x)) : [],
      };
}

class ReadReceiptInfo {
  String? userID;
  String? groupID;
  List<String>? msgIDList;
  int? readTime;
  int? msgFrom;
  int? contentType;
  int? sessionType;

  // 兼容后端当前 MarkAsReadTips 结构
  List<int>? _seqs;
  String? _conversationID;
  int? _hasReadSeq;

  ReadReceiptInfo({
    this.userID,
    this.groupID,
    this.msgIDList,
    this.readTime,
    this.msgFrom,
    this.contentType,
    this.sessionType,
  });

  ReadReceiptInfo.fromJson(Map<String, dynamic> json) {
    userID = json['markAsReadUserID'] ?? json['uid'] ?? json['userID'];
    groupID = json['groupID'];
    _conversationID = json['conversationID'] as String?;
    _hasReadSeq = json['hasReadSeq'];

    if (json['seqs'] is List) {
      _seqs = (json['seqs'] as List)
          .map((e) => e is int ? e : int.tryParse('$e') ?? 0)
          .where((e) => e > 0)
          .toList();
    }
    if (json['msgIDList'] is List) {
      msgIDList = (json['msgIDList'] as List).map((e) => '$e').toList();
    }

    readTime = json['readTime'];
    msgFrom = json['msgFrom'];
    contentType = json['contentType'];
    sessionType = json['sessionType'];
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    data['userID'] = userID;
    data['msgIDList'] = msgIDList;
    data['readTime'] = readTime;
    data['msgFrom'] = msgFrom;
    data['contentType'] = contentType;
    data['sessionType'] = sessionType;
    return data;
  }

  List<int>? get seqs => _seqs;
  String? get conversationID => _conversationID;
  int? get hasReadSeq => _hasReadSeq;
}

class RevokedInfo {
  String? revokerID;
  int? revokerRole;
  String? revokerNickname;
  String? clientMsgID;
  int? revokeTime;
  int? sourceMessageSendTime;
  String? sourceMessageSendID;
  String? sourceMessageSenderNickname;
  int? sessionType;

  RevokedInfo({
    this.revokerID,
    this.revokerRole,
    this.revokerNickname,
    this.clientMsgID,
    this.revokeTime,
    this.sourceMessageSendTime,
    this.sourceMessageSendID,
    this.sourceMessageSenderNickname,
    this.sessionType,
  });

  RevokedInfo.fromJson(Map<String, dynamic> json) {
    revokerID = json['revokerID'];
    revokerRole = json['revokerRole'];
    revokerNickname = json['revokerNickname'];
    clientMsgID = json['clientMsgID'];
    revokeTime = json['revokeTime'];
    sourceMessageSendTime = json['sourceMessageSendTime'];
    sourceMessageSendID = json['sourceMessageSendID'];
    sourceMessageSenderNickname = json['sourceMessageSenderNickname'];
    sessionType = json['sessionType'];
  }

  Map<String, dynamic> toJson() => {
        'revokerID': revokerID,
        'revokerRole': revokerRole,
        'revokerNickname': revokerNickname,
        'clientMsgID': clientMsgID,
        'revokeTime': revokeTime,
        'sourceMessageSendTime': sourceMessageSendTime,
        'sourceMessageSendID': sourceMessageSendID,
        'sourceMessageSenderNickname': sourceMessageSenderNickname,
        'sessionType': sessionType,
      };
}
