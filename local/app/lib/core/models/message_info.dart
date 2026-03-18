import 'dart:convert';

class MessageInfo {
  String clientMsgID;
  String serverMsgID;
  int createTime;
  int sendTime;
  int sessionType; // 1:单聊 3:群聊
  String sendID;
  String recvID;
  String groupID;
  int contentType;
  int platformID;
  String senderNickname;
  String senderFaceUrl;
  String content;
  int seq;
  bool isRead;
  int status; // 1:发送中 2:发送成功 3:发送失败
  bool isRevoked;
  String ex;

  // 解析后的内容
  TextElem? textElem;
  PictureElem? pictureElem;
  SoundElem? soundElem;
  VideoElem? videoElem;
  FileElem? fileElem;

  MessageInfo({
    this.clientMsgID = '',
    this.serverMsgID = '',
    this.createTime = 0,
    this.sendTime = 0,
    this.sessionType = 1,
    this.sendID = '',
    this.recvID = '',
    this.groupID = '',
    this.contentType = 101,
    this.platformID = 0,
    this.senderNickname = '',
    this.senderFaceUrl = '',
    this.content = '',
    this.seq = 0,
    this.isRead = false,
    this.status = 2,
    this.isRevoked = false,
    this.ex = '',
    this.textElem,
    this.pictureElem,
    this.soundElem,
    this.videoElem,
    this.fileElem,
  });

  factory MessageInfo.fromJson(Map<String, dynamic> json) {
    final msg = MessageInfo(
      clientMsgID: json['clientMsgID'] ?? '',
      serverMsgID: json['serverMsgID'] ?? '',
      createTime: json['createTime'] ?? 0,
      sendTime: json['sendTime'] ?? 0,
      sessionType: json['sessionType'] ?? 1,
      sendID: json['sendID'] ?? '',
      recvID: json['recvID'] ?? '',
      groupID: json['groupID'] ?? '',
      contentType: json['contentType'] ?? 101,
      platformID: json['platformID'] ?? 0,
      senderNickname: json['senderNickname'] ?? '',
      senderFaceUrl: json['senderFaceUrl'] ?? '',
      content: json['content'] ?? '',
      seq: json['seq'] ?? 0,
      isRead: json['isRead'] ?? false,
      status: json['status'] ?? 2,
      isRevoked: json['isRevoked'] ?? false,
      ex: json['ex'] ?? '',
    );
    msg._parseContent();
    return msg;
  }

  void _parseContent() {
    if (content.isEmpty) return;
    try {
      final map = jsonDecode(content) as Map<String, dynamic>?;
      if (map == null) return;
      switch (contentType) {
        case ContentType.text:
          textElem = TextElem.fromJson(map);
          break;
        case ContentType.picture:
          pictureElem = PictureElem.fromJson(map);
          break;
        case ContentType.sound:
          soundElem = SoundElem.fromJson(map);
          break;
        case ContentType.video:
          videoElem = VideoElem.fromJson(map);
          break;
        case ContentType.file:
          fileElem = FileElem.fromJson(map);
          break;
      }
    } catch (_) {}
  }

  Map<String, dynamic> toJson() => {
        'clientMsgID': clientMsgID,
        'serverMsgID': serverMsgID,
        'createTime': createTime,
        'sendTime': sendTime,
        'sessionType': sessionType,
        'sendID': sendID,
        'recvID': recvID,
        'groupID': groupID,
        'contentType': contentType,
        'platformID': platformID,
        'senderNickname': senderNickname,
        'senderFaceUrl': senderFaceUrl,
        'content': content,
        'seq': seq,
        'isRead': isRead,
        'status': status,
        'isRevoked': isRevoked,
        'ex': ex,
      };

  bool get isSingleChat => sessionType == 1;
  bool get isGroupChat => sessionType == 3;

  String get textContent {
    if (textElem != null) return textElem!.content;
    if (contentType == ContentType.text) {
      try {
        final map = jsonDecode(content);
        return map['content'] ?? content;
      } catch (_) {
        return content;
      }
    }
    return '[${ContentType.name(contentType)}]';
  }
}

class ContentType {
  static const int text = 101;
  static const int picture = 102;
  static const int sound = 103;
  static const int video = 104;
  static const int file = 105;
  static const int atText = 106;
  static const int merger = 115;
  static const int card = 116;
  static const int location = 117;
  static const int custom = 110;
  static const int typing = 113;
  static const int quote = 114;
  // 通知类
  static const int revokeNotification = 111;
  static const int friendNotificationBegin = 1200;
  static const int groupNotificationBegin = 1500;

  static String name(int type) {
    switch (type) {
      case text:
        return '文本';
      case picture:
        return '图片';
      case sound:
        return '语音';
      case video:
        return '视频';
      case file:
        return '文件';
      case atText:
        return '@消息';
      case location:
        return '位置';
      case card:
        return '名片';
      case merger:
        return '合并消息';
      case quote:
        return '引用';
      case custom:
        return '自定义';
      case typing:
        return '正在输入';
      default:
        if (type > 1000) return '通知';
        return '消息';
    }
  }
}

class TextElem {
  String content;
  TextElem({this.content = ''});
  factory TextElem.fromJson(Map<String, dynamic> json) =>
      TextElem(content: json['content'] ?? '');
}

class PictureElem {
  String sourcePath;
  PictureInfo? sourcePicture;
  PictureInfo? bigPicture;
  PictureInfo? snapshotPicture;

  PictureElem({
    this.sourcePath = '',
    this.sourcePicture,
    this.bigPicture,
    this.snapshotPicture,
  });

  factory PictureElem.fromJson(Map<String, dynamic> json) => PictureElem(
        sourcePath: json['sourcePath'] ?? '',
        sourcePicture: json['sourcePicture'] != null
            ? PictureInfo.fromJson(json['sourcePicture'])
            : null,
        bigPicture: json['bigPicture'] != null
            ? PictureInfo.fromJson(json['bigPicture'])
            : null,
        snapshotPicture: json['snapshotPicture'] != null
            ? PictureInfo.fromJson(json['snapshotPicture'])
            : null,
      );

  String get url =>
      bigPicture?.url ?? sourcePicture?.url ?? snapshotPicture?.url ?? '';
}

class PictureInfo {
  String uuid;
  String type;
  int size;
  int width;
  int height;
  String url;

  PictureInfo({
    this.uuid = '',
    this.type = '',
    this.size = 0,
    this.width = 0,
    this.height = 0,
    this.url = '',
  });

  factory PictureInfo.fromJson(Map<String, dynamic> json) => PictureInfo(
        uuid: json['uuid'] ?? '',
        type: json['type'] ?? '',
        size: json['size'] ?? 0,
        width: json['width'] ?? 0,
        height: json['height'] ?? 0,
        url: json['url'] ?? '',
      );
}

class SoundElem {
  String uuid;
  String soundPath;
  String sourceUrl;
  int dataSize;
  int duration;

  SoundElem({
    this.uuid = '',
    this.soundPath = '',
    this.sourceUrl = '',
    this.dataSize = 0,
    this.duration = 0,
  });

  factory SoundElem.fromJson(Map<String, dynamic> json) => SoundElem(
        uuid: json['uuid'] ?? '',
        soundPath: json['soundPath'] ?? '',
        sourceUrl: json['sourceUrl'] ?? '',
        dataSize: json['dataSize'] ?? 0,
        duration: json['duration'] ?? 0,
      );
}

class VideoElem {
  String videoPath;
  String videoUUID;
  String videoUrl;
  String videoType;
  int videoSize;
  int duration;
  String snapshotPath;
  String snapshotUUID;
  int snapshotSize;
  String snapshotUrl;
  int snapshotWidth;
  int snapshotHeight;

  VideoElem({
    this.videoPath = '',
    this.videoUUID = '',
    this.videoUrl = '',
    this.videoType = '',
    this.videoSize = 0,
    this.duration = 0,
    this.snapshotPath = '',
    this.snapshotUUID = '',
    this.snapshotSize = 0,
    this.snapshotUrl = '',
    this.snapshotWidth = 0,
    this.snapshotHeight = 0,
  });

  factory VideoElem.fromJson(Map<String, dynamic> json) => VideoElem(
        videoPath: json['videoPath'] ?? '',
        videoUUID: json['videoUUID'] ?? '',
        videoUrl: json['videoUrl'] ?? '',
        videoType: json['videoType'] ?? '',
        videoSize: json['videoSize'] ?? 0,
        duration: json['duration'] ?? 0,
        snapshotPath: json['snapshotPath'] ?? '',
        snapshotUUID: json['snapshotUUID'] ?? '',
        snapshotSize: json['snapshotSize'] ?? 0,
        snapshotUrl: json['snapshotUrl'] ?? '',
        snapshotWidth: json['snapshotWidth'] ?? 0,
        snapshotHeight: json['snapshotHeight'] ?? 0,
      );
}

class FileElem {
  String filePath;
  String uuid;
  String sourceUrl;
  String fileName;
  int fileSize;

  FileElem({
    this.filePath = '',
    this.uuid = '',
    this.sourceUrl = '',
    this.fileName = '',
    this.fileSize = 0,
  });

  factory FileElem.fromJson(Map<String, dynamic> json) => FileElem(
        filePath: json['filePath'] ?? '',
        uuid: json['uuid'] ?? '',
        sourceUrl: json['sourceUrl'] ?? '',
        fileName: json['fileName'] ?? '',
        fileSize: json['fileSize'] ?? 0,
      );
}
