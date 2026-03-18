import '../http_client.dart';
import '../config.dart';
import '../models/conversation_info.dart';

class ConversationApi {
  static Future<List<ConversationInfo>> getSortedConversationList({
    int pageNumber = 1,
    int showNumber = 50,
  }) async {
    final data =
        await HttpClient.post('/conversation/get_sorted_conversation_list',
            data: {
          'userID': Config.userID,
          'pagination': {'pageNumber': pageNumber, 'showNumber': showNumber},
        });
    if (data == null) return [];
    final list = data['conversationElems'] ?? data['conversations'] ?? data;
    if (list is List) {
      return list
          .map((e) {
            final conv = e['conversation'] ?? e;
            return ConversationInfo.fromJson(conv as Map<String, dynamic>);
          })
          .toList();
    }
    return [];
  }

  static Future<ConversationInfo?> getConversation(
      String conversationID) async {
    final data = await HttpClient.post('/conversation/get_conversation', data: {
      'ownerUserID': Config.userID,
      'conversationID': conversationID,
    });
    if (data == null) return null;
    final conv = data['conversation'] ?? data;
    if (conv is Map<String, dynamic>) {
      return ConversationInfo.fromJson(conv);
    }
    return null;
  }

  static Future<void> setConversations({
    required List<String> conversationIDs,
    bool? isPinned,
    int? recvMsgOpt,
    bool? isPrivateChat,
  }) async {
    final conversation = <String, dynamic>{};
    if (isPinned != null) conversation['isPinned'] = isPinned;
    if (recvMsgOpt != null) conversation['recvMsgOpt'] = recvMsgOpt;
    if (isPrivateChat != null) conversation['isPrivateChat'] = isPrivateChat;

    await HttpClient.post('/conversation/set_conversations', data: {
      'userID': Config.userID,
      'conversationIDs': conversationIDs,
      'conversation': conversation,
    });
  }

  static Future<void> deleteConversations(
      List<String> conversationIDs) async {
    await HttpClient.post('/conversation/delete_conversations', data: {
      'userID': Config.userID,
      'conversationIDs': conversationIDs,
    });
  }

  static Future<void> markConversationAsRead(String conversationID) async {
    await HttpClient.post('/msg/mark_conversation_as_read', data: {
      'userID': Config.userID,
      'conversationID': conversationID,
    }, showErrorToast: false);
  }
}
