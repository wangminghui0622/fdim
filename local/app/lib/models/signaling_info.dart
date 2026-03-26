class SignalingInfo {
  String? userID;
  InvitationInfo? invitation;

  SignalingInfo({this.userID, this.invitation});

  SignalingInfo.fromJson(Map<String, dynamic> json) {
    invitation = json['invitation'] == null
        ? null
        : InvitationInfo.fromJson(json['invitation']);
    userID = json['userID'] ?? invitation?.inviterUserID;
  }

  Map<String, dynamic> toJson() {
    return {
      'userID': userID,
      'invitation': invitation?.toJson(),
    };
  }
}

class InvitationInfo {
  String? inviterUserID;
  List<String>? inviteeUserIDList;
  String? groupID;
  String? roomID;
  int? timeout;
  String? mediaType;
  int? sessionType;
  int? platformID;

  InvitationInfo({
    this.inviterUserID,
    this.inviteeUserIDList,
    this.groupID,
    this.roomID,
    this.timeout,
    this.mediaType,
    this.sessionType,
    this.platformID,
  });

  InvitationInfo.fromJson(Map<String, dynamic> json) {
    inviterUserID = json['inviterUserID'];
    inviteeUserIDList = json['inviteeUserIDList']?.cast<String>();
    groupID = json['groupID'];
    roomID = json['roomID'];
    timeout = json['timeout'];
    mediaType = json['mediaType'];
    sessionType = json['sessionType'];
    platformID = json['platformID'];
  }

  Map<String, dynamic> toJson() {
    return {
      'inviterUserID': inviterUserID,
      'inviteeUserIDList': inviteeUserIDList,
      'groupID': groupID,
      'roomID': roomID,
      'timeout': timeout,
      'mediaType': mediaType,
      'sessionType': sessionType,
      'platformID': platformID,
    };
  }
}

class SignalingCertificate {
  String? token;
  String? roomID;
  String? liveURL;
  List<String>? busyLineUserIDList;

  SignalingCertificate({
    this.token,
    this.roomID,
    this.liveURL,
    this.busyLineUserIDList,
  });

  SignalingCertificate.fromJson(Map<String, dynamic> json) {
    token = json['token'];
    roomID = json['roomID'];
    liveURL = json['serverUrl'];
    busyLineUserIDList = json['busyLineUserIDList']?.cast<String>();
  }

  Map<String, dynamic> toJson() {
    return {
      'token': token,
      'roomID': roomID,
      'serverUrl': liveURL,
      'busyLineUserIDList': busyLineUserIDList,
    };
  }
}

enum CallType { audio, video }

enum CallState {
  call,
  beCalled,
  reject,
  beRejected,
  calling,
  beAccepted,
  hangup,
  beHangup,
  connecting,
  cancel,
  beCanceled,
  timeout,
}

class CallEvent {
  CallState state;
  SignalingInfo data;
  dynamic fields;

  CallEvent(this.state, this.data, {this.fields});
}

class CustomMessageType {
  static const int callingInvite = 200;
  static const int callingAccept = 201;
  static const int callingReject = 202;
  static const int callingCancel = 203;
  static const int callingHungup = 204;
  static const int call = 205;
  /// 通话结果（持久化到聊天记录）
  static const int callResult = 210;
}

/// 通话结果类型（只记录动作，显示文本根据 isMe 在渲染时决定）
class CallResultType {
  static const String rejected = 'rejected';     // 拒绝（发送方=拒绝者）
  static const String cancelled = 'cancelled';   // 取消（发送方=取消者）
  static const String timeout = 'timeout';       // 超时未接听
  static const String completed = 'completed';   // 通话完成
}

/// 通话结果消息体
class CallResultInfo {
  String callType;  // 'audio' | 'video'
  String result;    // CallResultType
  int duration;     // 通话时长（秒），仅 completed 时有效

  CallResultInfo({
    required this.callType,
    required this.result,
    this.duration = 0,
  });

  CallResultInfo.fromJson(Map<String, dynamic> json)
      : callType = json['callType'] ?? 'audio',
        result = json['result'] ?? '',
        duration = json['duration'] ?? 0;

  Map<String, dynamic> toJson() => {
        'callType': callType,
        'result': result,
        'duration': duration,
      };

  /// 获取显示文本，isMe 表示消息是否是自己发的
  String getDisplayText(bool isMe) {
    switch (result) {
      case CallResultType.rejected:
        return isMe ? '已拒绝' : '对方已拒绝';
      case CallResultType.cancelled:
        return isMe ? '已取消' : '对方已取消';
      case CallResultType.timeout:
        return '未接听';
      case CallResultType.completed:
        final min = (duration ~/ 60).toString().padLeft(2, '0');
        final sec = (duration % 60).toString().padLeft(2, '0');
        return '通话时长 $min:$sec';
      default:
        return '通话结束';
    }
  }

  bool get isVideo => callType == 'video';
}
