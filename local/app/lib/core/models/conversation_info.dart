class ConversationInfo {
  String conversationID;
  int conversationType; // 1:单聊 3:群聊
  String userID;
  String groupID;
  String showName;
  String faceURL;
  int recvMsgOpt;
  int unreadCount;
  int latestMsgSendTime;
  String latestMsg;
  bool isPinned;
  bool isPrivateChat;
  String draftText;
  int draftTextTime;
  String ex;

  ConversationInfo({
    this.conversationID = '',
    this.conversationType = 1,
    this.userID = '',
    this.groupID = '',
    this.showName = '',
    this.faceURL = '',
    this.recvMsgOpt = 0,
    this.unreadCount = 0,
    this.latestMsgSendTime = 0,
    this.latestMsg = '',
    this.isPinned = false,
    this.isPrivateChat = false,
    this.draftText = '',
    this.draftTextTime = 0,
    this.ex = '',
  });

  factory ConversationInfo.fromJson(Map<String, dynamic> json) =>
      ConversationInfo(
        conversationID: json['conversationID'] ?? '',
        conversationType: json['conversationType'] ?? 1,
        userID: json['userID'] ?? '',
        groupID: json['groupID'] ?? '',
        showName: json['showName'] ?? '',
        faceURL: json['faceURL'] ?? '',
        recvMsgOpt: json['recvMsgOpt'] ?? 0,
        unreadCount: json['unreadCount'] ?? 0,
        latestMsgSendTime: json['latestMsgSendTime'] ?? 0,
        latestMsg: json['latestMsg'] ?? '',
        isPinned: json['isPinned'] ?? false,
        isPrivateChat: json['isPrivateChat'] ?? false,
        draftText: json['draftText'] ?? '',
        draftTextTime: json['draftTextTime'] ?? 0,
        ex: json['ex'] ?? '',
      );

  Map<String, dynamic> toJson() => {
        'conversationID': conversationID,
        'conversationType': conversationType,
        'userID': userID,
        'groupID': groupID,
        'showName': showName,
        'faceURL': faceURL,
        'recvMsgOpt': recvMsgOpt,
        'unreadCount': unreadCount,
        'latestMsgSendTime': latestMsgSendTime,
        'latestMsg': latestMsg,
        'isPinned': isPinned,
        'isPrivateChat': isPrivateChat,
        'draftText': draftText,
        'draftTextTime': draftTextTime,
        'ex': ex,
      };
}
