import '../models/message.dart';
import '../models/user_info.dart';

class OnAdvancedMsgListener {
  Function(Message msg)? onMsgDeleted;
  Function(RevokedInfo info)? onNewRecvMessageRevoked;
  Function(List<ReadReceiptInfo> list)? onRecvC2CReadReceipt;
  Function(List<ReadReceiptInfo> list)? onRecvGroupReadReceipt;
  Function(Message msg)? onRecvNewMessage;
  Function(Message msg)? onRecvOfflineNewMessage;
  Function(Message msg)? onRecvOnlineOnlyMessage;

  String id;

  OnAdvancedMsgListener({
    this.onMsgDeleted,
    this.onNewRecvMessageRevoked,
    this.onRecvC2CReadReceipt,
    this.onRecvGroupReadReceipt,
    this.onRecvNewMessage,
    this.onRecvOfflineNewMessage,
    this.onRecvOnlineOnlyMessage,
  }) : id = "id_${DateTime.now().microsecondsSinceEpoch}";

  void msgDeleted(Message msg) {
    onMsgDeleted?.call(msg);
  }

  void newRecvMessageRevoked(RevokedInfo info) {
    onNewRecvMessageRevoked?.call(info);
  }

  void recvC2CReadReceipt(List<ReadReceiptInfo> list) {
    onRecvC2CReadReceipt?.call(list);
  }

  void recvGroupReadReceipt(List<ReadReceiptInfo> list) {
    onRecvGroupReadReceipt?.call(list);
  }

  void recvNewMessage(Message msg) {
    onRecvNewMessage?.call(msg);
  }

  void recvOfflineNewMessage(Message msg) {
    onRecvOfflineNewMessage?.call(msg);
  }

  void recvOnlineOnlyMessage(Message msg) {
    onRecvOnlineOnlyMessage?.call(msg);
  }
}
