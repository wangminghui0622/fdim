import 'dart:convert';
import 'message.dart';
import '../enum/conversation_type.dart';

class ConversationInfo {
  String conversationID;
  int? conversationType;
  String? userID;
  String? groupID;
  String? showName;
  String? faceURL;
  int? recvMsgOpt;
  int unreadCount = 0;
  Message? latestMsg;
  int? latestMsgSendTime;
  String? draftText;
  int? draftTextTime;
  bool? isPinned;
  bool? isPrivateChat;
  int? burnDuration;
  bool? isMsgDestruct;
  int? msgDestructTime;
  String? ex;
  bool? isNotInGroup;
  int? groupAtType;

  ConversationInfo({
    required this.conversationID,
    this.conversationType,
    this.userID,
    this.groupID,
    this.showName,
    this.faceURL,
    this.recvMsgOpt,
    this.unreadCount = 0,
    this.latestMsg,
    this.latestMsgSendTime,
    this.draftText,
    this.draftTextTime,
    this.isPrivateChat,
    this.burnDuration,
    this.isPinned,
    this.isNotInGroup,
    this.ex,
    this.groupAtType,
    this.isMsgDestruct,
    this.msgDestructTime,
  });

  ConversationInfo.fromJson(Map<String, dynamic> json)
      : conversationID = json['conversationID'] ?? '' {
    conversationType = json['conversationType'];
    userID = json['userID'];
    groupID = json['groupID'];
    showName = json['showName'];
    faceURL = json['faceURL'];
    recvMsgOpt = json['recvMsgOpt'];
    unreadCount = json['unreadCount'] ?? 0;
    try {
      if (json['latestMsg'] is String && (json['latestMsg'] as String).isNotEmpty) {
        latestMsg = Message.fromJson(jsonDecode(json['latestMsg']));
      } else if (json['latestMsg'] is Map) {
        latestMsg = Message.fromJson(json['latestMsg']);
      }
    } catch (_) {}
    latestMsgSendTime = json['latestMsgSendTime'];
    draftText = json['draftText'];
    draftTextTime = json['draftTextTime'];
    isPinned = json['isPinned'];
    isPrivateChat = json['isPrivateChat'];
    burnDuration = json['burnDuration'];
    isNotInGroup = json['isNotInGroup'];
    groupAtType = json['groupAtType'];
    ex = json['ex'];
    isMsgDestruct = json['isMsgDestruct'];
    msgDestructTime = json['msgDestructTime'];
  }

  Map<String, dynamic> toJson() => {
        'conversationID': conversationID,
        'conversationType': conversationType,
        'userID': userID,
        'groupID': groupID,
        'showName': showName,
        'faceURL': faceURL,
        'recvMsgOpt': recvMsgOpt,
        'unreadCount': unreadCount,
        'latestMsg': latestMsg?.toJson(),
        'latestMsgSendTime': latestMsgSendTime,
        'draftText': draftText,
        'draftTextTime': draftTextTime,
        'isPinned': isPinned,
        'isPrivateChat': isPrivateChat,
        'burnDuration': burnDuration,
        'isNotInGroup': isNotInGroup,
        'groupAtType': groupAtType,
        'ex': ex,
        'isMsgDestruct': isMsgDestruct,
        'msgDestructTime': msgDestructTime,
      };

  bool get isSingleChat => conversationType == ConversationType.single;
  bool get isGroupChat =>
      conversationType == ConversationType.group ||
      conversationType == ConversationType.superGroup;
  bool get isValid => isSingleChat || (isGroupChat && !(isNotInGroup ?? false));

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is ConversationInfo &&
          runtimeType == other.runtimeType &&
          conversationID == other.conversationID;

  @override
  int get hashCode => conversationID.hashCode;
}
