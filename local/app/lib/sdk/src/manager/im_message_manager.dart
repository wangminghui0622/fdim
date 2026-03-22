import 'dart:convert';
import 'package:flutter/foundation.dart';
import 'package:uuid/uuid.dart';
import '../../flutter_openim_sdk.dart';
import '../../../core/http_client.dart';
import '../../../core/config.dart';
import '../../../core/local_store.dart';

class MessageManager {
  OnMsgSendProgressListener? msgSendProgressListener;
  late OnAdvancedMsgListener msgListener;
  OnCustomBusinessListener? customBusinessListener;

  Future setAdvancedMsgListener(OnAdvancedMsgListener listener) async {
    msgListener = listener;
  }

  void setMsgSendProgressListener(OnMsgSendProgressListener listener) {
    msgSendProgressListener = listener;
  }

  Future setCustomBusinessListener(OnCustomBusinessListener listener) async {
    customBusinessListener = listener;
  }

  /// Send a message
  Future<Message> sendMessage({
    required Message message,
    required OfflinePushInfo offlinePushInfo,
    String? userID,
    String? groupID,
    bool isOnlineOnly = false,
    String? operationID,
  }) async {
    message.sendID = Config.userID;
    message.senderNickname = Config.nickname;
    message.senderFaceUrl = Config.faceURL;
    message.status = MessageStatus.sending;

    final sendMsgData = <String, dynamic>{
      'clientMsgID': message.clientMsgID,
      'serverMsgID': message.serverMsgID,
      'sendID': message.sendID,
      'recvID': userID ?? '',
      'groupID': groupID ?? '',
      'senderNickname': message.senderNickname,
      'senderFaceURL': message.senderFaceUrl,
      'senderPlatformID': Config.platformID,
      'contentType': message.contentType,
      'sessionType': message.sessionType,
      'createTime': message.createTime,
    };

    // Add content based on type (JSON-encode to string for server)
    dynamic contentObj;
    if (message.textElem != null) {
      contentObj = message.textElem!.toJson();
    } else if (message.pictureElem != null) {
      contentObj = message.pictureElem!.toJson();
    } else if (message.soundElem != null) {
      contentObj = message.soundElem!.toJson();
    } else if (message.videoElem != null) {
      contentObj = message.videoElem!.toJson();
    } else if (message.fileElem != null) {
      contentObj = message.fileElem!.toJson();
    } else if (message.atTextElem != null) {
      contentObj = message.atTextElem!.toJson();
    } else if (message.locationElem != null) {
      contentObj = message.locationElem!.toJson();
    } else if (message.customElem != null) {
      contentObj = message.customElem!.toJson();
    } else if (message.quoteElem != null) {
      contentObj = message.quoteElem!.toJson();
    } else if (message.cardElem != null) {
      contentObj = message.cardElem!.toJson();
    } else if (message.mergeElem != null) {
      contentObj = message.mergeElem!.toJson();
    } else if (message.faceElem != null) {
      contentObj = message.faceElem!.toJson();
    }
    if (contentObj != null) {
      sendMsgData['content'] = contentObj is String ? contentObj : jsonEncode(contentObj);
    }

    try {
      final resp = await HttpClient.post('/msg/send_msg', data: {
        'recvID': userID ?? '',
        'sendMsg': sendMsgData,
      });
      message.status = MessageStatus.succeeded;
      if (resp is Map) {
        message.serverMsgID = resp['serverMsgID'] ?? message.serverMsgID;
        message.sendTime = resp['sendTime'] ?? message.sendTime;
      }
      // 与官方一致：发送成功后写入本地 DB
      final convID = _resolveConversationID(message, userID, groupID);
      if (convID.isNotEmpty) {
        await LocalStore.putMessage(convID, message);
        await LocalStore.updateLatestMsg(convID, message);
      }
      return message;
    } catch (e) {
      message.status = MessageStatus.failed;
      rethrow;
    }
  }

  /// Create a text message
  Future<Message> createTextMessage({
    required String text,
    String? operationID,
  }) async {
    final now = DateTime.now().millisecondsSinceEpoch;
    return Message(
      clientMsgID: const Uuid().v4(),
      createTime: now,
      sendTime: now,
      sendID: Config.userID,
      senderNickname: Config.nickname,
      senderFaceUrl: Config.faceURL,
      senderPlatformID: Config.platformID,
      contentType: MessageType.text,
      sessionType: ConversationType.single,
      status: MessageStatus.sending,
      textElem: TextElem(content: text),
    );
  }

  /// Create an image message
  Future<Message> createImageMessage({
    required String imagePath,
    String? operationID,
  }) async {
    final now = DateTime.now().millisecondsSinceEpoch;
    return Message(
      clientMsgID: const Uuid().v4(),
      createTime: now,
      sendTime: now,
      sendID: Config.userID,
      contentType: MessageType.picture,
      sessionType: ConversationType.single,
      status: MessageStatus.sending,
      pictureElem: PictureElem(sourcePath: imagePath),
    );
  }

  Future<Message> createImageMessageFromFullPath({
    required String imagePath,
    String? operationID,
  }) async {
    return createImageMessage(imagePath: imagePath);
  }

  /// Create a sound message
  Future<Message> createSoundMessage({
    required String soundPath,
    required int duration,
    String? operationID,
  }) async {
    final now = DateTime.now().millisecondsSinceEpoch;
    return Message(
      clientMsgID: const Uuid().v4(),
      createTime: now,
      sendTime: now,
      sendID: Config.userID,
      contentType: MessageType.voice,
      sessionType: ConversationType.single,
      status: MessageStatus.sending,
      soundElem: SoundElem(soundPath: soundPath, duration: duration),
    );
  }

  Future<Message> createSoundMessageFromFullPath({
    required String soundPath,
    required int duration,
    String? operationID,
  }) async {
    return createSoundMessage(soundPath: soundPath, duration: duration);
  }

  /// Create a video message
  Future<Message> createVideoMessage({
    required String videoPath,
    required String videoType,
    required int duration,
    required String snapshotPath,
    String? operationID,
  }) async {
    final now = DateTime.now().millisecondsSinceEpoch;
    return Message(
      clientMsgID: const Uuid().v4(),
      createTime: now,
      sendTime: now,
      sendID: Config.userID,
      contentType: MessageType.video,
      sessionType: ConversationType.single,
      status: MessageStatus.sending,
      videoElem: VideoElem(
        videoPath: videoPath,
        videoType: videoType,
        duration: duration,
        snapshotPath: snapshotPath,
      ),
    );
  }

  Future<Message> createVideoMessageFromFullPath({
    required String videoPath,
    required String videoType,
    required int duration,
    required String snapshotPath,
    String? operationID,
  }) async {
    return createVideoMessage(
      videoPath: videoPath,
      videoType: videoType,
      duration: duration,
      snapshotPath: snapshotPath,
    );
  }

  /// Create a file message
  Future<Message> createFileMessage({
    required String filePath,
    required String fileName,
    String? operationID,
  }) async {
    final now = DateTime.now().millisecondsSinceEpoch;
    return Message(
      clientMsgID: const Uuid().v4(),
      createTime: now,
      sendTime: now,
      sendID: Config.userID,
      contentType: MessageType.file,
      sessionType: ConversationType.single,
      status: MessageStatus.sending,
      fileElem: FileElem(filePath: filePath, fileName: fileName),
    );
  }

  Future<Message> createFileMessageFromFullPath({
    required String filePath,
    required String fileName,
    String? operationID,
  }) async {
    return createFileMessage(filePath: filePath, fileName: fileName);
  }

  /// Create a location message
  Future<Message> createLocationMessage({
    required double latitude,
    required double longitude,
    required String description,
    String? operationID,
  }) async {
    final now = DateTime.now().millisecondsSinceEpoch;
    return Message(
      clientMsgID: const Uuid().v4(),
      createTime: now,
      sendTime: now,
      sendID: Config.userID,
      contentType: MessageType.location,
      sessionType: ConversationType.single,
      status: MessageStatus.sending,
      locationElem: LocationElem(
        latitude: latitude,
        longitude: longitude,
        description: description,
      ),
    );
  }

  /// Create a custom message
  Future<Message> createCustomMessage({
    required String data,
    required String extension,
    required String description,
    String? operationID,
  }) async {
    final now = DateTime.now().millisecondsSinceEpoch;
    return Message(
      clientMsgID: const Uuid().v4(),
      createTime: now,
      sendTime: now,
      sendID: Config.userID,
      contentType: MessageType.custom,
      sessionType: ConversationType.single,
      status: MessageStatus.sending,
      customElem: CustomElem(
        data: data,
        extension: extension,
        description: description,
      ),
    );
  }

  /// Create a quote message
  Future<Message> createQuoteMessage({
    required String text,
    required Message quoteMsg,
    String? operationID,
  }) async {
    final now = DateTime.now().millisecondsSinceEpoch;
    return Message(
      clientMsgID: const Uuid().v4(),
      createTime: now,
      sendTime: now,
      sendID: Config.userID,
      contentType: MessageType.quote,
      sessionType: ConversationType.single,
      status: MessageStatus.sending,
      quoteElem: QuoteElem(text: text, quoteMessage: quoteMsg),
    );
  }

  /// Create a card message
  Future<Message> createCardMessage({
    required String userID,
    required String nickname,
    String? faceURL,
    String? ex,
    String? operationID,
  }) async {
    final now = DateTime.now().millisecondsSinceEpoch;
    return Message(
      clientMsgID: const Uuid().v4(),
      createTime: now,
      sendTime: now,
      sendID: Config.userID,
      contentType: MessageType.card,
      sessionType: ConversationType.single,
      status: MessageStatus.sending,
      cardElem: CardElem(
        userID: userID,
        nickname: nickname,
        faceURL: faceURL,
        ex: ex,
      ),
    );
  }

  /// Create a face (emoji) message
  Future<Message> createFaceMessage({
    int index = -1,
    String? data,
    String? operationID,
  }) async {
    final now = DateTime.now().millisecondsSinceEpoch;
    return Message(
      clientMsgID: const Uuid().v4(),
      createTime: now,
      sendTime: now,
      sendID: Config.userID,
      contentType: MessageType.customFace,
      sessionType: ConversationType.single,
      status: MessageStatus.sending,
      faceElem: FaceElem(index: index, data: data),
    );
  }

  /// Create a forward message
  Future<Message> createForwardMessage({
    required Message message,
    String? operationID,
  }) async {
    final now = DateTime.now().millisecondsSinceEpoch;
    return Message(
      clientMsgID: const Uuid().v4(),
      createTime: now,
      sendTime: now,
      sendID: Config.userID,
      contentType: message.contentType,
      sessionType: message.sessionType,
      status: MessageStatus.sending,
      textElem: message.textElem,
      pictureElem: message.pictureElem,
      soundElem: message.soundElem,
      videoElem: message.videoElem,
      fileElem: message.fileElem,
      atTextElem: message.atTextElem,
      locationElem: message.locationElem,
      customElem: message.customElem,
      quoteElem: message.quoteElem,
      mergeElem: message.mergeElem,
      faceElem: message.faceElem,
      cardElem: message.cardElem,
    );
  }

  /// Create a merged message
  Future<Message> createMergerMessage({
    required List<Message> messageList,
    required String title,
    required List<String> summaryList,
    String? operationID,
  }) async {
    final now = DateTime.now().millisecondsSinceEpoch;
    return Message(
      clientMsgID: const Uuid().v4(),
      createTime: now,
      sendTime: now,
      sendID: Config.userID,
      contentType: MessageType.merger,
      sessionType: ConversationType.single,
      status: MessageStatus.sending,
      mergeElem: MergeElem(
        title: title,
        abstractList: summaryList,
        multiMessage: messageList,
      ),
    );
  }

  /// Create an @ text message
  Future<Message> createTextAtMessage({
    required String text,
    required List<String> atUserIDList,
    List<AtUserInfo> atUserInfoList = const [],
    Message? quoteMessage,
    String? operationID,
  }) async {
    final now = DateTime.now().millisecondsSinceEpoch;
    return Message(
      clientMsgID: const Uuid().v4(),
      createTime: now,
      sendTime: now,
      sendID: Config.userID,
      contentType: MessageType.atText,
      sessionType: ConversationType.superGroup,
      status: MessageStatus.sending,
      atTextElem: AtTextElem(
        text: text,
        atUserList: atUserIDList,
        atUsersInfo: atUserInfoList,
        quoteMessage: quoteMessage,
      ),
    );
  }

  /// Get advanced history message list (seq-based pagination)
  Future<AdvancedMessage> getAdvancedHistoryMessageList({
    String? conversationID,
    Message? startMsg,
    int? count,
    String? operationID,
  }) async {
    final n = count ?? 40;
    debugPrint('[SDK] getAdvancedHistoryMessageList: conversationID=$conversationID, count=$n');

    // 1. Determine the upper bound seq
    int endSeq;
    if (startMsg != null && (startMsg.seq ?? 0) > 0) {
      endSeq = startMsg.seq! - 1;
      debugPrint('[SDK] Using startMsg.seq: $endSeq');
    } else {
      // Get newest seq for this conversation
      debugPrint('[SDK] Fetching newest seq for conversationID: $conversationID');
      final seqData = await HttpClient.post('/msg/newest_seq', data: {
        'userID': Config.userID,
        'conversationID': conversationID ?? '',
      });
      debugPrint('[SDK] newest_seq response: $seqData');
      if (seqData == null) {
        debugPrint('[SDK] seqData is null');
        endSeq = 0;
      } else {
        endSeq = (seqData is Map ? (seqData['maxSeq'] ?? 0) : 0) as int;
      }
      debugPrint('[SDK] Got maxSeq: $endSeq');
    }

    List<Message> messages = [];
    int startSeq = 1;

    if (endSeq <= 0) {
      debugPrint('[SDK] endSeq <= 0, skipping server fetch');
    } else {
    // 2. Build seq list: [max(1, endSeq-n+1) .. endSeq]
    startSeq = endSeq - n + 1 > 0 ? endSeq - n + 1 : 1;
    final seqs = List<int>.generate(endSeq - startSeq + 1, (i) => startSeq + i);
    debugPrint('[SDK] Pulling messages from seq $startSeq to $endSeq (${seqs.length} seqs)');

    // 3. Pull messages by seq
    final data = await HttpClient.post('/msg/pull_msg_by_seq', data: {
      'userID': Config.userID,
      'conversationID': conversationID ?? '',
      'seqs': seqs,
    });
    debugPrint('[SDK] pull_msg_by_seq response type: ${data.runtimeType}');
    debugPrint('[SDK] pull_msg_by_seq response data: $data');

    if (data != null && data is Map) {
      final msgsList = data['msgs'];
      debugPrint('[SDK] msgsList type: ${msgsList.runtimeType}');

      if (msgsList is List && msgsList.isNotEmpty) {
        final firstItem = msgsList[0];
        debugPrint('[SDK] firstItem type: ${firstItem.runtimeType}');

        if (firstItem is Map && firstItem.containsKey('Msgs')) {
          final msgList = firstItem['Msgs'];
          debugPrint('[SDK] msgList type: ${msgList.runtimeType}, length: ${msgList is List ? msgList.length : 'N/A'}');
          if (msgList is List) {
            messages = msgList
                .whereType<Map>()
                .map((e) => Message.fromJson(Map<String, dynamic>.from(e)))
                .toList();
          }
        } else {
          debugPrint('[SDK] Detected flat msg list format, length: ${msgsList.length}');
          messages = msgsList
              .whereType<Map>()
              .map((e) => Message.fromJson(Map<String, dynamic>.from(e)))
              .toList();
        }
      }
    }

    // Sort by seq ascending
    messages.sort((a, b) => (a.seq ?? 0).compareTo(b.seq ?? 0));

    // 服务端不返回 status 字段（这是客户端概念），但能从服务端拉到的消息一定是发送成功的
    for (final msg in messages) {
      msg.status ??= MessageStatus.succeeded;
    }

    } // end if (endSeq > 0)

    // 与官方一致：先从本地 DB 恢复已读状态，再写入本地 DB
    // 顺序很重要：必须先 apply 再 put，否则 put 会用服务端的 isRead=false 覆盖本地已有的 isRead=true
    if (conversationID != null && conversationID.isNotEmpty) {
      LocalStore.applyHasReadSeq(conversationID, messages, Config.userID);
      await LocalStore.putMessages(conversationID, messages);
    }

    // 合并本地创建的消息（如好友通过后的申请语句和问候消息，没有 serverMsgID）
    if (conversationID != null && conversationID.isNotEmpty) {
      final localMsgs = LocalStore.getMessages(conversationID);
      final existingIDs = messages.map((m) => m.clientMsgID).toSet();
      for (final lm in localMsgs) {
        if (lm.clientMsgID != null && !existingIDs.contains(lm.clientMsgID)) {
          messages.add(lm);
        }
      }
      // 按 sendTime 排序（本地消息没有 seq）
      messages.sort((a, b) => (a.sendTime ?? 0).compareTo(b.sendTime ?? 0));
    }

    debugPrint('[SDK] Returning ${messages.length} messages, isEnd: ${startSeq <= 1}');

    return AdvancedMessage(
      messageList: messages,
      isEnd: startSeq <= 1,
    );
  }

  /// 获取会话的 hasReadSeq 和 maxSeq（与官方一致）
  Future<Map<String, int>> getHasReadAndMaxSeq(String conversationID) async {
    final data = await HttpClient.post(
      '/msg/get_conversations_has_read_and_max_seq',
      data: {
        'userID': Config.userID,
        'conversationIDs': [conversationID],
      },
    );
    int maxSeq = 0;
    int hasReadSeq = 0;
    debugPrint('[SDK] getHasReadAndMaxSeq raw response: $data');
    if (data != null && data is Map && data['seqs'] is Map) {
      final seqInfo = data['seqs'][conversationID];
      debugPrint('[SDK] getHasReadAndMaxSeq seqInfo for $conversationID: $seqInfo');
      if (seqInfo is Map) {
        maxSeq = (seqInfo['maxSeq'] ?? seqInfo['MaxSeq'] ?? 0) as int;
        hasReadSeq = (seqInfo['hasReadSeq'] ?? seqInfo['HasReadSeq'] ?? 0) as int;
      }
    }
    debugPrint('[SDK] getHasReadAndMaxSeq: conv=$conversationID maxSeq=$maxSeq hasReadSeq=$hasReadSeq');
    return {'maxSeq': maxSeq, 'hasReadSeq': hasReadSeq};
  }

  /// 按 seq 范围拉取消息 [fromSeq..toSeq]，返回按 seq 升序
  Future<List<Message>> getMessagesBySeqRange({
    required String conversationID,
    required int fromSeq,
    required int toSeq,
  }) async {
    if (fromSeq > toSeq || toSeq <= 0) return [];
    final seqs = List<int>.generate(toSeq - fromSeq + 1, (i) => fromSeq + i);
    debugPrint('[SDK] getMessagesBySeqRange: conv=$conversationID from=$fromSeq to=$toSeq (${seqs.length} seqs)');

    final data = await HttpClient.post('/msg/pull_msg_by_seq', data: {
      'userID': Config.userID,
      'conversationID': conversationID,
      'seqs': seqs,
    });

    List<Message> messages = [];
    if (data != null && data is Map) {
      final msgsList = data['msgs'];
      if (msgsList is List && msgsList.isNotEmpty) {
        final firstItem = msgsList[0];
        if (firstItem is Map && firstItem.containsKey('Msgs')) {
          final msgList = firstItem['Msgs'];
          if (msgList is List) {
            messages = msgList
                .whereType<Map>()
                .map((e) => Message.fromJson(Map<String, dynamic>.from(e)))
                .toList();
          }
        } else {
          messages = msgsList
              .whereType<Map>()
              .map((e) => Message.fromJson(Map<String, dynamic>.from(e)))
              .toList();
        }
      }
    }

    messages.sort((a, b) => (a.seq ?? 0).compareTo(b.seq ?? 0));
    for (final msg in messages) {
      msg.status ??= MessageStatus.succeeded;
    }

    // 恢复本地已读状态并缓存
    if (conversationID.isNotEmpty) {
      LocalStore.applyHasReadSeq(conversationID, messages, Config.userID);
      await LocalStore.putMessages(conversationID, messages);
    }

    return messages;
  }

  /// Revoke a message (by seq)
  Future revokeMessage({
    required String conversationID,
    required int seq,
    String? operationID,
  }) async {
    await HttpClient.post('/msg/revoke_msg', data: {
      'userID': Config.userID,
      'conversationID': conversationID,
      'seq': seq,
    });
  }

  /// Delete a message from local storage
  Future deleteMessageFromLocalStorage({
    required String conversationID,
    required String clientMsgID,
    String? operationID,
  }) async {
    // Local-only deletion, not supported via HTTP
  }

  /// Delete a message from local and server (by seq)
  Future<dynamic> deleteMessageFromLocalAndSvr({
    required String conversationID,
    required int seq,
    String? operationID,
  }) async {
    await HttpClient.post('/msg/delete_msgs', data: {
      'userID': Config.userID,
      'conversationID': conversationID,
      'seqs': [seq],
    });
  }

  /// Search local messages (via server search)
  Future<SearchResult> searchLocalMessages({
    String? conversationID,
    List<String> keywordList = const [],
    int keywordListMatchType = 0,
    List<String> senderUserIDList = const [],
    List<int> messageTypeList = const [],
    int searchTimePosition = 0,
    int searchTimePeriod = 0,
    int pageIndex = 1,
    int count = 40,
    String? operationID,
  }) async {
    final data = await HttpClient.post('/msg/search_msg', data: {
      'userID': Config.userID,
      'conversationID': conversationID,
      'keywordList': keywordList,
      'messageTypeList': messageTypeList,
      'pageIndex': pageIndex,
      'count': count,
    }, showErrorToast: false);
    if (data == null) return SearchResult(totalCount: 0, searchResultItems: []);
    return SearchResult.fromJson(data);
  }

  /// Set app badge
  Future setAppBadge(int count, {String? operationID}) async {
    // No-op for non-mobile
  }

  /// 从消息中解析 conversationID
  String _resolveConversationID(Message msg, String? userID, String? groupID) {
    final recvID = userID ?? msg.recvID ?? '';
    final gid = groupID ?? msg.groupID ?? '';
    final st = msg.sessionType;
    if (st == ConversationType.single && recvID.isNotEmpty) {
      final ids = [msg.sendID ?? Config.userID, recvID]..sort();
      return 'si_${ids[0]}_${ids[1]}';
    }
    if ((st == ConversationType.group || st == ConversationType.superGroup) && gid.isNotEmpty) {
      return 'sg_$gid';
    }
    return '';
  }
}
