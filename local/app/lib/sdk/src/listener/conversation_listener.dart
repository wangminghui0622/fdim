import 'package:flutter/foundation.dart';
import '../models/conversation_info.dart';
import '../models/input_status_changed_data.dart';

class OnConversationListener {
  Function(List<ConversationInfo> list)? onConversationChanged;
  Function(List<ConversationInfo> list)? onNewConversation;
  Function(int count)? onTotalUnreadMessageCountChanged;
  Function(bool? reinstalled)? onSyncServerStart;
  Function(int? progress)? onSyncServerProgress;
  Function(bool? reinstalled)? onSyncServerFinish;
  Function(bool? reinstalled)? onSyncServerFailed;
  ValueChanged<InputStatusChangedData>? onInputStatusChanged;

  OnConversationListener({
    this.onConversationChanged,
    this.onNewConversation,
    this.onTotalUnreadMessageCountChanged,
    this.onSyncServerStart,
    this.onSyncServerProgress,
    this.onSyncServerFinish,
    this.onSyncServerFailed,
    this.onInputStatusChanged,
  });

  void conversationChanged(List<ConversationInfo> list) {
    onConversationChanged?.call(list);
  }

  void newConversation(List<ConversationInfo> list) {
    onNewConversation?.call(list);
  }

  void totalUnreadMessageCountChanged(int count) {
    onTotalUnreadMessageCountChanged?.call(count);
  }

  void syncServerStart(bool? reinstalled) {
    onSyncServerStart?.call(reinstalled);
  }

  void syncServerProgress(int? progress) {
    onSyncServerProgress?.call(progress);
  }

  void syncServerFailed(bool? reinstalled) {
    onSyncServerFailed?.call(reinstalled);
  }

  void syncServerFinish(bool? reinstalled) {
    onSyncServerFinish?.call(reinstalled);
  }

  void conversationUserInputStatusChanged(InputStatusChangedData data) {
    onInputStatusChanged?.call(data);
  }
}
