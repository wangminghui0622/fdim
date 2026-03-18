import 'message.dart';
import 'user_info.dart';

class SearchResult {
  int? totalCount;
  List<SearchResultItem>? searchResultItems;
  List<FriendInfo>? findResultItems;

  SearchResult({this.totalCount, this.searchResultItems, this.findResultItems});

  SearchResult.fromJson(Map<String, dynamic> json) {
    totalCount = json['totalCount'];
    if (json['searchResultItems'] is List) {
      searchResultItems = (json['searchResultItems'] as List)
          .map((e) => SearchResultItem.fromJson(e))
          .toList();
    }
    if (json['findResultItems'] is List) {
      findResultItems = (json['findResultItems'] as List)
          .map((e) => FriendInfo.fromJson(e))
          .toList();
    }
  }
}

class SearchResultItem {
  String? conversationID;
  int? messageCount;
  List<Message>? messageList;

  SearchResultItem({this.conversationID, this.messageCount, this.messageList});

  SearchResultItem.fromJson(Map<String, dynamic> json) {
    conversationID = json['conversationID'];
    messageCount = json['messageCount'];
    if (json['messageList'] is List) {
      messageList = (json['messageList'] as List)
          .map((e) => Message.fromJson(e))
          .toList();
    }
  }
}

class SearchParams {
  String? conversationID;
  List<String>? clientMsgIDList;

  SearchParams({this.conversationID, this.clientMsgIDList});

  Map<String, dynamic> toJson() => {
        'conversationID': conversationID,
        'clientMsgIDList': clientMsgIDList,
      };
}

class SearchFriendsInfo {
  FriendInfo? friendInfo;
  int? relationship;

  SearchFriendsInfo({this.friendInfo, this.relationship});

  SearchFriendsInfo.fromJson(Map<String, dynamic> json) {
    friendInfo = json['friendInfo'] != null
        ? FriendInfo.fromJson(json['friendInfo'])
        : FriendInfo.fromJson(json);
    relationship = json['relationship'];
  }
}
