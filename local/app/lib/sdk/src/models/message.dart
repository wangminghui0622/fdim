import 'dart:convert';
import '../enum/conversation_type.dart';

class Message {
  String? clientMsgID;
  String? serverMsgID;
  int? createTime;
  int? sendTime;
  int? sessionType;
  String? sendID;
  String? recvID;
  int? msgFrom;
  int? contentType;
  int? senderPlatformID;
  String? senderNickname;
  String? senderFaceUrl;
  String? groupID;
  String? localEx;
  int? seq;
  bool? isRead;
  int? hasReadTime;
  int? status;
  bool? isReact;
  bool? isExternalExtensions;
  OfflinePushInfo? offlinePush;
  String? attachedInfo;
  String? ex;
  Map<String, dynamic> exMap = {};
  PictureElem? pictureElem;
  SoundElem? soundElem;
  VideoElem? videoElem;
  FileElem? fileElem;
  AtTextElem? atTextElem;
  LocationElem? locationElem;
  CustomElem? customElem;
  QuoteElem? quoteElem;
  MergeElem? mergeElem;
  NotificationElem? notificationElem;
  FaceElem? faceElem;
  AttachedInfoElem? attachedInfoElem;
  TextElem? textElem;
  CardElem? cardElem;
  AdvancedTextElem? advancedTextElem;
  TypingElem? typingElem;

  Message({
    this.clientMsgID,
    this.serverMsgID,
    this.createTime,
    this.sendTime,
    this.sessionType,
    this.sendID,
    this.recvID,
    this.msgFrom,
    this.contentType,
    this.senderPlatformID,
    this.senderNickname,
    this.senderFaceUrl,
    this.groupID,
    this.localEx,
    this.seq,
    this.isRead,
    this.hasReadTime,
    this.status,
    this.offlinePush,
    this.attachedInfo,
    this.ex,
    this.exMap = const <String, dynamic>{},
    this.pictureElem,
    this.soundElem,
    this.videoElem,
    this.fileElem,
    this.atTextElem,
    this.locationElem,
    this.customElem,
    this.quoteElem,
    this.mergeElem,
    this.notificationElem,
    this.faceElem,
    this.attachedInfoElem,
    this.isExternalExtensions,
    this.isReact,
    this.textElem,
    this.cardElem,
    this.advancedTextElem,
    this.typingElem,
  });

  Message.fromJson(Map<String, dynamic> json) {
    clientMsgID = json['clientMsgID'];
    serverMsgID = json['serverMsgID'];
    createTime = json['createTime'];
    sendTime = json['sendTime'];
    sendID = json['sendID'];
    recvID = json['recvID'];
    msgFrom = json['msgFrom'];
    contentType = json['contentType'];
    senderPlatformID = json['senderPlatformID'];
    senderNickname = json['senderNickname'];
    senderFaceUrl = json['senderFaceUrl'];
    groupID = json['groupID'];
    localEx = json['localEx'];
    seq = json['seq'];
    isRead = json['isRead'];
    status = json['status'];
    offlinePush = json['offlinePush'] != null
        ? OfflinePushInfo.fromJson(json['offlinePush'])
        : null;
    attachedInfo = json['attachedInfo'];
    ex = json['ex'];
    exMap = json['exMap'] ?? {};
    sessionType = json['sessionType'];
    
    // Parse content field if it exists (server may return content as Base64 encoded JSON string)
    dynamic contentData = json['content'];
    if (contentData is String && contentData.isNotEmpty) {
      try {
        // First try to decode from Base64
        final decoded = utf8.decode(base64Decode(contentData));
        contentData = jsonDecode(decoded);
      } catch (_) {
        // If Base64 decode fails, try direct JSON decode
        try {
          contentData = jsonDecode(contentData);
        } catch (_) {}
      }
    }
    
    // Convert to Map<String, dynamic> if needed
    Map<String, dynamic>? contentMap;
    if (contentData is Map) {
      contentMap = Map<String, dynamic>.from(contentData);
    }
    
    // Parse specific elem fields from json first, fallback to content field
    pictureElem = json['pictureElem'] != null
        ? PictureElem.fromJson(json['pictureElem'])
        : (contentMap != null && contentType == 102 ? PictureElem.fromJson(contentMap) : null);
    soundElem = json['soundElem'] != null
        ? SoundElem.fromJson(json['soundElem'])
        : (contentMap != null && contentType == 103 ? SoundElem.fromJson(contentMap) : null);
    videoElem = json['videoElem'] != null
        ? VideoElem.fromJson(json['videoElem'])
        : (contentMap != null && contentType == 104 ? VideoElem.fromJson(contentMap) : null);
    fileElem = json['fileElem'] != null
        ? FileElem.fromJson(json['fileElem'])
        : (contentMap != null && contentType == 105 ? FileElem.fromJson(contentMap) : null);
    atTextElem = json['atTextElem'] != null
        ? AtTextElem.fromJson(json['atTextElem'])
        : (contentMap != null && contentType == 106 ? AtTextElem.fromJson(contentMap) : null);
    locationElem = json['locationElem'] != null
        ? LocationElem.fromJson(json['locationElem'])
        : (contentMap != null && contentType == 107 ? LocationElem.fromJson(contentMap) : null);
    customElem = json['customElem'] != null
        ? CustomElem.fromJson(json['customElem'])
        : (contentMap != null && contentType == 110 ? CustomElem.fromJson(contentMap) : null);
    quoteElem = json['quoteElem'] != null
        ? QuoteElem.fromJson(json['quoteElem'])
        : (contentMap != null && contentType == 114 ? QuoteElem.fromJson(contentMap) : null);
    mergeElem = json['mergeElem'] != null
        ? MergeElem.fromJson(json['mergeElem'])
        : null;
    notificationElem = json['notificationElem'] != null
        ? NotificationElem.fromJson(json['notificationElem'])
        : (contentMap != null && (contentType ?? 0) > 1000
            ? NotificationElem.fromJson(contentMap)
            : null);
    faceElem = json['faceElem'] != null
        ? FaceElem.fromJson(json['faceElem'])
        : (contentMap != null && contentType == 107 ? FaceElem.fromJson(contentMap) : null);
    attachedInfoElem = json['attachedInfoElem'] != null
        ? AttachedInfoElem.fromJson(json['attachedInfoElem'])
        : null;
    hasReadTime = json['hasReadTime'] ?? attachedInfoElem?.hasReadTime;
    isExternalExtensions = json['isExternalExtensions'];
    isReact = json['isReact'];
    
    // Parse textElem from json or from parsed content
    textElem = json['textElem'] != null
        ? TextElem.fromJson(json['textElem'])
        : (contentMap != null && contentType == 101 ? TextElem.fromJson(contentMap) : null);
    
    cardElem = json['cardElem'] != null
        ? CardElem.fromJson(json['cardElem'])
        : (contentMap != null && contentType == 108 ? CardElem.fromJson(contentMap) : null);
    advancedTextElem = json['advancedTextElem'] != null
        ? AdvancedTextElem.fromJson(json['advancedTextElem'])
        : null;
    typingElem = json['typingElem'] != null
        ? TypingElem.fromJson(json['typingElem'])
        : null;
  }

  Map<String, dynamic> toJson() => {
        'clientMsgID': clientMsgID,
        'serverMsgID': serverMsgID,
        'createTime': createTime,
        'sendTime': sendTime,
        'sendID': sendID,
        'recvID': recvID,
        'msgFrom': msgFrom,
        'contentType': contentType,
        'senderPlatformID': senderPlatformID,
        'senderNickname': senderNickname,
        'senderFaceUrl': senderFaceUrl,
        'groupID': groupID,
        'localEx': localEx,
        'seq': seq,
        'isRead': isRead,
        'hasReadTime': hasReadTime,
        'status': status,
        'offlinePush': offlinePush?.toJson(),
        'attachedInfo': attachedInfo,
        'ex': ex,
        'exMap': exMap,
        'sessionType': sessionType,
        'pictureElem': pictureElem?.toJson(),
        'soundElem': soundElem?.toJson(),
        'videoElem': videoElem?.toJson(),
        'fileElem': fileElem?.toJson(),
        'atTextElem': atTextElem?.toJson(),
        'locationElem': locationElem?.toJson(),
        'customElem': customElem?.toJson(),
        'quoteElem': quoteElem?.toJson(),
        'mergeElem': mergeElem?.toJson(),
        'notificationElem': notificationElem?.toJson(),
        'faceElem': faceElem?.toJson(),
        'attachedInfoElem': attachedInfoElem?.toJson(),
        'isExternalExtensions': isExternalExtensions,
        'isReact': isReact,
        'textElem': textElem?.toJson(),
        'cardElem': cardElem?.toJson(),
        'advancedTextElem': advancedTextElem?.toJson(),
        'typingElem': typingElem?.toJson(),
      };

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is Message &&
          runtimeType == other.runtimeType &&
          clientMsgID == other.clientMsgID;

  @override
  int get hashCode => clientMsgID.hashCode;

  void update(Message message) {
    if (this != message) return;
    serverMsgID = message.serverMsgID;
    createTime = message.createTime;
    sendTime = message.sendTime;
    sendID = message.sendID;
    recvID = message.recvID;
    msgFrom = message.msgFrom;
    contentType = message.contentType;
    senderPlatformID = message.senderPlatformID;
    senderNickname = message.senderNickname;
    senderFaceUrl = message.senderFaceUrl;
    groupID = message.groupID;
    seq = message.seq;
    isRead = message.isRead;
    hasReadTime = message.hasReadTime;
    status = message.status;
    offlinePush = message.offlinePush;
    attachedInfo = message.attachedInfo;
    ex = message.ex;
    exMap = message.exMap;
    sessionType = message.sessionType;
    pictureElem = message.pictureElem;
    soundElem = message.soundElem;
    videoElem = message.videoElem;
    fileElem = message.fileElem;
    atTextElem = message.atTextElem;
    locationElem = message.locationElem;
    customElem = message.customElem;
    quoteElem = message.quoteElem;
    mergeElem = message.mergeElem;
    notificationElem = message.notificationElem;
    faceElem = message.faceElem;
    attachedInfoElem = message.attachedInfoElem;
    textElem = message.textElem;
    cardElem = message.cardElem;
    advancedTextElem = message.advancedTextElem;
    typingElem = message.typingElem;
  }

  bool get isSingleChat => sessionType == ConversationType.single;
  bool get isGroupChat =>
      sessionType == ConversationType.group ||
      sessionType == ConversationType.superGroup;

  bool get isCustomType => contentType == 110;

  String get textContent {
    if (textElem?.content != null) return textElem!.content!;
    if (atTextElem?.text != null) return atTextElem!.text!;
    if (notificationElem?.defaultTips != null) return notificationElem!.defaultTips!;
    return '[${_contentTypeName}]';
  }

  String get _contentTypeName {
    switch (contentType) {
      case 101: return '文本';
      case 102: return '图片';
      case 103: return '语音';
      case 104: return '视频';
      case 105: return '文件';
      case 106: return '@消息';
      case 108: return '名片';
      case 109: return '位置';
      case 110: return '自定义';
      case 114: return '引用';
      default: return '消息';
    }
  }
}

class TextElem {
  String? content;
  TextElem({this.content});
  TextElem.fromJson(Map<String, dynamic> json) { content = json['content']; }
  Map<String, dynamic> toJson() => {'content': content};
}

class PictureElem {
  String? sourcePath;
  PictureInfo? sourcePicture;
  PictureInfo? bigPicture;
  PictureInfo? snapshotPicture;

  PictureElem({this.sourcePath, this.sourcePicture, this.bigPicture, this.snapshotPicture});
  PictureElem.fromJson(Map<String, dynamic> json) {
    sourcePath = json['sourcePath'];
    sourcePicture = json['sourcePicture'] != null ? PictureInfo.fromJson(json['sourcePicture']) : null;
    bigPicture = json['bigPicture'] != null ? PictureInfo.fromJson(json['bigPicture']) : null;
    snapshotPicture = json['snapshotPicture'] != null ? PictureInfo.fromJson(json['snapshotPicture']) : null;
  }
  Map<String, dynamic> toJson() => {
    'sourcePath': sourcePath,
    'sourcePicture': sourcePicture?.toJson(),
    'bigPicture': bigPicture?.toJson(),
    'snapshotPicture': snapshotPicture?.toJson(),
  };
}

class PictureInfo {
  String? uuid;
  String? type;
  int? size;
  int? width;
  int? height;
  String? url;

  PictureInfo({this.uuid, this.type, this.size, this.width, this.height, this.url});
  PictureInfo.fromJson(Map<String, dynamic> json) {
    uuid = json['uuid']; type = json['type']; size = json['size'];
    width = json['width']; height = json['height']; url = json['url'];
  }
  Map<String, dynamic> toJson() => {'uuid': uuid, 'type': type, 'size': size, 'width': width, 'height': height, 'url': url};
}

class SoundElem {
  String? uuid;
  String? soundPath;
  String? sourceUrl;
  int? dataSize;
  int? duration;

  SoundElem({this.uuid, this.soundPath, this.sourceUrl, this.dataSize, this.duration});
  SoundElem.fromJson(Map<String, dynamic> json) {
    uuid = json['uuid']; soundPath = json['soundPath']; sourceUrl = json['sourceUrl'];
    dataSize = json['dataSize']; duration = json['duration'];
  }
  Map<String, dynamic> toJson() => {'uuid': uuid, 'soundPath': soundPath, 'sourceUrl': sourceUrl, 'dataSize': dataSize, 'duration': duration};
}

class VideoElem {
  String? videoPath;
  String? videoUUID;
  String? videoUrl;
  String? videoType;
  int? videoSize;
  int? duration;
  String? snapshotPath;
  String? snapshotUUID;
  int? snapshotSize;
  String? snapshotUrl;
  int? snapshotWidth;
  int? snapshotHeight;

  VideoElem({this.videoPath, this.videoUUID, this.videoUrl, this.videoType, this.videoSize, this.duration, this.snapshotPath, this.snapshotUUID, this.snapshotSize, this.snapshotUrl, this.snapshotWidth, this.snapshotHeight});
  VideoElem.fromJson(Map<String, dynamic> json) {
    videoPath = json['videoPath']; videoUUID = json['videoUUID']; videoUrl = json['videoUrl'];
    videoType = json['videoType']; videoSize = json['videoSize']; duration = json['duration'];
    snapshotPath = json['snapshotPath']; snapshotUUID = json['snapshotUUID']; snapshotSize = json['snapshotSize'];
    snapshotUrl = json['snapshotUrl']; snapshotWidth = json['snapshotWidth']; snapshotHeight = json['snapshotHeight'];
  }
  Map<String, dynamic> toJson() => {
    'videoPath': videoPath, 'videoUUID': videoUUID, 'videoUrl': videoUrl, 'videoType': videoType,
    'videoSize': videoSize, 'duration': duration, 'snapshotPath': snapshotPath, 'snapshotUUID': snapshotUUID,
    'snapshotSize': snapshotSize, 'snapshotUrl': snapshotUrl, 'snapshotWidth': snapshotWidth, 'snapshotHeight': snapshotHeight,
  };
}

class FileElem {
  String? filePath;
  String? uuid;
  String? sourceUrl;
  String? fileName;
  int? fileSize;

  FileElem({this.filePath, this.uuid, this.sourceUrl, this.fileName, this.fileSize});
  FileElem.fromJson(Map<String, dynamic> json) {
    filePath = json['filePath']; uuid = json['uuid']; sourceUrl = json['sourceUrl'];
    fileName = json['fileName']; fileSize = json['fileSize'];
  }
  Map<String, dynamic> toJson() => {'filePath': filePath, 'uuid': uuid, 'sourceUrl': sourceUrl, 'fileName': fileName, 'fileSize': fileSize};
}

class AtTextElem {
  String? text;
  List<String>? atUserList;
  bool? isAtSelf;
  List<AtUserInfo>? atUsersInfo;
  Message? quoteMessage;

  AtTextElem({this.text, this.atUserList, this.isAtSelf, this.atUsersInfo, this.quoteMessage});
  AtTextElem.fromJson(Map<String, dynamic> json) {
    text = json['text'];
    if (json['atUserList'] is List) atUserList = (json['atUserList'] as List).map((e) => '$e').toList();
    isAtSelf = json['isAtSelf'];
    if (json['atUsersInfo'] is List) atUsersInfo = (json['atUsersInfo'] as List).map((e) => AtUserInfo.fromJson(e)).toList();
    quoteMessage = json['quoteMessage'] != null ? Message.fromJson(json['quoteMessage']) : null;
  }
  Map<String, dynamic> toJson() => {
    'text': text, 'atUserList': atUserList, 'isAtSelf': isAtSelf,
    'atUsersInfo': atUsersInfo?.map((e) => e.toJson()).toList(),
    'quoteMessage': quoteMessage?.toJson(),
  };
}

class LocationElem {
  String? description;
  double? longitude;
  double? latitude;

  LocationElem({this.description, this.longitude, this.latitude});
  LocationElem.fromJson(Map<String, dynamic> json) {
    description = json['description'];
    longitude = json['longitude'] is int ? (json['longitude'] as int).toDouble() : json['longitude'];
    latitude = json['latitude'] is int ? (json['latitude'] as int).toDouble() : json['latitude'];
  }
  Map<String, dynamic> toJson() => {'description': description, 'longitude': longitude, 'latitude': latitude};
}

class CustomElem {
  String? data;
  String? extension;
  String? description;

  CustomElem({this.data, this.extension, this.description});
  CustomElem.fromJson(Map<String, dynamic> json) {
    data = json['data']; extension = json['extension']; description = json['description'];
  }
  Map<String, dynamic> toJson() => {'data': data, 'extension': extension, 'description': description};
}

class QuoteElem {
  String? text;
  Message? quoteMessage;

  QuoteElem({this.text, this.quoteMessage});
  QuoteElem.fromJson(Map<String, dynamic> json) {
    text = json['text'];
    if (json['quoteMessage'] is Map) quoteMessage = Message.fromJson(json['quoteMessage']);
  }
  Map<String, dynamic> toJson() => {'text': text, 'quoteMessage': quoteMessage?.toJson()};
}

class MergeElem {
  String? title;
  List<String>? abstractList;
  List<Message>? multiMessage;

  MergeElem({this.title, this.abstractList, this.multiMessage});
  MergeElem.fromJson(Map<String, dynamic> json) {
    title = json['title'];
    if (json['abstractList'] is List) abstractList = json['abstractList'].cast<String>();
    if (json['multiMessage'] is List) multiMessage = (json['multiMessage'] as List).map((e) => Message.fromJson(e)).toList();
  }
  Map<String, dynamic> toJson() => {'title': title, 'abstractList': abstractList, 'multiMessage': multiMessage?.map((e) => e.toJson()).toList()};
}

class NotificationElem {
  String? detail;
  String? defaultTips;

  NotificationElem({this.detail, this.defaultTips});
  NotificationElem.fromJson(Map<String, dynamic> json) { detail = json['detail']; defaultTips = json['defaultTips']; }
  Map<String, dynamic> toJson() => {'detail': detail, 'defaultTips': defaultTips};
}

class FaceElem {
  int? index;
  String? data;

  FaceElem({this.index, this.data});
  FaceElem.fromJson(Map<String, dynamic> json) { index = json['index']; data = json['data']; }
  Map<String, dynamic> toJson() => {'index': index, 'data': data};
}

class AttachedInfoElem {
  GroupHasReadInfo? groupHasReadInfo;
  bool? isPrivateChat;
  int? hasReadTime;
  int? burnDuration;
  bool? notSenderNotificationPush;

  AttachedInfoElem({this.groupHasReadInfo, this.isPrivateChat, this.hasReadTime, this.burnDuration, this.notSenderNotificationPush});
  AttachedInfoElem.fromJson(Map<String, dynamic> json) {
    groupHasReadInfo = json['groupHasReadInfo'] != null ? GroupHasReadInfo.fromJson(json['groupHasReadInfo']) : null;
    isPrivateChat = json['isPrivateChat']; hasReadTime = json['hasReadTime'];
    burnDuration = json['burnDuration']; notSenderNotificationPush = json['notSenderNotificationPush'];
  }
  Map<String, dynamic> toJson() => {
    'groupHasReadInfo': groupHasReadInfo?.toJson(), 'isPrivateChat': isPrivateChat,
    'hasReadTime': hasReadTime, 'burnDuration': burnDuration, 'notSenderNotificationPush': notSenderNotificationPush,
  };
}

class GroupHasReadInfo {
  int? hasReadCount;
  int? unreadCount;

  GroupHasReadInfo.fromJson(Map<String, dynamic> json) { hasReadCount = json['hasReadCount'] ?? 0; unreadCount = json['unreadCount'] ?? 0; }
  Map<String, dynamic> toJson() => {'hasReadCount': hasReadCount, 'unreadCount': unreadCount};
}

class OfflinePushInfo {
  String? title;
  String? desc;
  String? ex;
  String? iOSPushSound;
  bool? iOSBadgeCount;

  OfflinePushInfo({this.title, this.desc, this.ex, this.iOSPushSound, this.iOSBadgeCount});
  OfflinePushInfo.fromJson(Map<String, dynamic> json) {
    title = json['title']; desc = json['desc']; ex = json['ex'];
    iOSPushSound = json['iOSPushSound']; iOSBadgeCount = json['iOSBadgeCount'];
  }
  Map<String, dynamic> toJson() => {'title': title, 'desc': desc, 'ex': ex, 'iOSPushSound': iOSPushSound, 'iOSBadgeCount': iOSBadgeCount};
}

class AtUserInfo {
  String? atUserID;
  String? groupNickname;

  AtUserInfo({this.atUserID, this.groupNickname});
  AtUserInfo.fromJson(Map<String, dynamic> json) { atUserID = json['atUserID']; groupNickname = json['groupNickname']; }
  Map<String, dynamic> toJson() => {'atUserID': atUserID, 'groupNickname': groupNickname};
}

class CardElem {
  String? userID;
  String? nickname;
  String? faceURL;
  String? ex;

  CardElem({this.userID, this.nickname, this.faceURL, this.ex});
  CardElem.fromJson(Map<String, dynamic> json) { userID = json['userID']; nickname = json['nickname']; faceURL = json['faceURL']; ex = json['ex']; }
  Map<String, dynamic> toJson() => {'userID': userID, 'nickname': nickname, 'faceURL': faceURL, 'ex': ex};
}

class TypingElem {
  String? msgTips;
  TypingElem({this.msgTips});
  TypingElem.fromJson(Map<String, dynamic> json) { msgTips = json['msgTips']; }
  Map<String, dynamic> toJson() => {'msgTips': msgTips};
}

class AdvancedTextElem {
  String? text;
  List<MessageEntity>? messageEntityList;
  AdvancedTextElem({this.text, this.messageEntityList});
  AdvancedTextElem.fromJson(Map<String, dynamic> json) {
    text = json['text'];
    messageEntityList = json['messageEntityList'] != null
        ? (json['messageEntityList'] as List).map((e) => MessageEntity.fromJson(e)).toList() : null;
  }
  Map<String, dynamic> toJson() => {'text': text, 'messageEntityList': messageEntityList?.map((e) => e.toJson()).toList()};
}

class MessageEntity {
  String? type;
  int? offset;
  int? length;
  String? url;
  String? ex;

  MessageEntity({this.type, this.offset, this.length, this.url, this.ex});
  MessageEntity.fromJson(Map<String, dynamic> json) { type = json['type']; offset = json['offset']; length = json['length']; url = json['url']; ex = json['ex']; }
  Map<String, dynamic> toJson() => {'type': type, 'offset': offset, 'length': length, 'url': url, 'ex': ex};
}

class AdvancedMessage {
  List<Message>? messageList;
  bool? isEnd;
  int? errCode;
  String? errMsg;
  int? lastMinSeq;

  AdvancedMessage({this.messageList, this.isEnd, this.errCode, this.errMsg, this.lastMinSeq});
  AdvancedMessage.fromJson(Map<String, dynamic> json) {
    if (json['messageList'] is List) {
      messageList = (json['messageList'] as List).map((e) => Message.fromJson(e)).toList();
    }
    isEnd = json['isEnd'];
    errCode = json['errCode'];
    errMsg = json['errMsg'];
    lastMinSeq = json['lastMinSeq'];
  }
}

class RichMessageInfo {
  String? text;
  String? type;
  String? url;
  int? offset;
  int? length;
  String? ex;

  RichMessageInfo({this.text, this.type, this.url, this.offset, this.length, this.ex});
  Map<String, dynamic> toJson() => {'text': text, 'type': type, 'url': url, 'offset': offset, 'length': length, 'ex': ex};
}
