import 'package:flutter/material.dart';
import 'package:flutter_easyloading/flutter_easyloading.dart';

import '../../core/config.dart';
import '../../core/local_store.dart';
import '../../sdk/flutter_openim_sdk.dart';

class FriendApplyPage extends StatefulWidget {
  const FriendApplyPage({super.key});

  @override
  State<FriendApplyPage> createState() => _FriendApplyPageState();
}

class _FriendApplyPageState extends State<FriendApplyPage> {
  List<FriendApplicationInfo> _applies = [];
  bool _loading = true;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    try {
      final list = await OpenIM.iMManager.friendshipManager
          .getFriendApplicationListAsRecipient();
      if (mounted) setState(() { _applies = list; _loading = false; });
    } catch (_) {
      if (mounted) setState(() => _loading = false);
    }
  }

  Future<void> _accept(FriendApplicationInfo apply) async {
    EasyLoading.show();
    try {
      await OpenIM.iMManager.friendshipManager.acceptFriendApplication(
          userID: apply.fromUserID ?? '');
      EasyLoading.showToast('已同意');

      // 好友通过后在本地创建"来自 A"的消息，触发会话创建
      _createLocalFirstMessage(apply);

      _load();
    } catch (e) {
      EasyLoading.showToast('操作失败');
    }
  }

  /// 好友通过后，在 B 本地创建两条消息并创建会话：
  /// 1) A 的申请语句（来自 A，左侧）
  /// 2) 默认问候（来自 B，右侧）
  Future<void> _createLocalFirstMessage(FriendApplicationInfo apply) async {
    try {
      final friendUserID = apply.fromUserID ?? '';
      if (friendUserID.isEmpty) return;

      final ids = [Config.userID, friendUserID]..sort();
      final convID = 'si_${ids[0]}_${ids[1]}';
      final now = DateTime.now().millisecondsSinceEpoch;
      int unread = 0;
      Message? latestMsg;

      final friendNickname = (apply.fromNickname ?? '').isNotEmpty
          ? apply.fromNickname!
          : friendUserID;

      // 消息1: A 的申请语句（如果有）— 来自 A，B 的聊天页左侧
      if ((apply.reqMsg ?? '').isNotEmpty) {
        final msg1 = Message(
          clientMsgID: '${friendUserID}_${now}_req',
          sendID: friendUserID,
          recvID: Config.userID,
          senderNickname: friendNickname,
          senderFaceUrl: apply.fromFaceURL,
          sessionType: ConversationType.single,
          contentType: MessageType.text,
          createTime: now,
          sendTime: now,
          status: MessageStatus.succeeded,
          isRead: false,
          textElem: TextElem(content: apply.reqMsg!),
        );
        await LocalStore.putMessage(convID, msg1);
        latestMsg = msg1;
        unread++;
      }

      // 消息2: 默认问候 — 来自 B（自己），B 的聊天页右侧
      final msg2 = Message(
        clientMsgID: '${Config.userID}_${now}_greet',
        sendID: Config.userID,
        recvID: friendUserID,
        senderNickname: Config.nickname.isNotEmpty ? Config.nickname : Config.userID,
        senderFaceUrl: Config.faceURL,
        sessionType: ConversationType.single,
        contentType: MessageType.text,
        createTime: now + 1,
        sendTime: now + 1,
        status: MessageStatus.succeeded,
        isRead: true,
        textElem: TextElem(content: '我们已经是好友了，可以开始聊天了'),
      );
      await LocalStore.putMessage(convID, msg2);
      latestMsg = msg2;

      // 创建或更新会话（用 A 的昵称）
      final existingConv = LocalStore.getConversation(convID);
      final conv = existingConv ?? ConversationInfo(
        conversationID: convID,
        conversationType: ConversationType.single,
        userID: friendUserID,
        showName: friendNickname,
        faceURL: apply.fromFaceURL,
      );
      conv.latestMsg = latestMsg;
      conv.latestMsgSendTime = now + 1;
      conv.unreadCount = (existingConv?.unreadCount ?? 0) + unread;
      await LocalStore.putConversation(conv);

      // 通知 UI 刷新
      OpenIM.iMManager.conversationManager.listener
          .conversationChanged([conv]);
      OpenIM.iMManager.conversationManager.listener
          .totalUnreadMessageCountChanged(LocalStore.getTotalUnreadCount());

      debugPrint('[FriendApply] 2 local messages created for $friendUserID in $convID');
    } catch (e) {
      debugPrint('[FriendApply] Failed to create local messages: $e');
    }
  }

  Future<void> _reject(FriendApplicationInfo apply) async {
    EasyLoading.show();
    try {
      await OpenIM.iMManager.friendshipManager.refuseFriendApplication(
          userID: apply.fromUserID ?? '');
      EasyLoading.showToast('已拒绝');
      _load();
    } catch (e) {
      EasyLoading.showToast('操作失败');
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('好友申请')),
      body: _loading
          ? const Center(child: CircularProgressIndicator())
          : _applies.isEmpty
              ? const Center(child: Text('暂无申请', style: TextStyle(color: Colors.grey)))
              : RefreshIndicator(
                  onRefresh: _load,
                  child: ListView.separated(
                    itemCount: _applies.length,
                    separatorBuilder: (_, __) => const Divider(height: 1),
                    itemBuilder: (_, i) => _buildItem(_applies[i]),
                  ),
                ),
    );
  }

  Widget _buildItem(FriendApplicationInfo apply) {
    return ListTile(
      leading: CircleAvatar(
        backgroundColor: Colors.grey[300],
        backgroundImage: (apply.fromFaceURL ?? '').isNotEmpty
            ? NetworkImage(apply.fromFaceURL!) : null,
        child: (apply.fromFaceURL ?? '').isEmpty
            ? Text((apply.fromNickname ?? '').isNotEmpty ? apply.fromNickname![0].toUpperCase() : '?',
                style: const TextStyle(color: Colors.white))
            : null,
      ),
      title: Text((apply.fromNickname ?? '').isNotEmpty ? apply.fromNickname! : (apply.fromUserID ?? '')),
      subtitle: Text((apply.reqMsg ?? '').isNotEmpty ? apply.reqMsg! : '请求添加你为好友',
          style: const TextStyle(fontSize: 13, color: Colors.grey)),
      trailing: apply.handleResult == 0
          ? Row(mainAxisSize: MainAxisSize.min, children: [
              TextButton(onPressed: () => _reject(apply), child: const Text('拒绝')),
              const SizedBox(width: 4),
              ElevatedButton(
                onPressed: () => _accept(apply),
                style: ElevatedButton.styleFrom(
                  backgroundColor: const Color(0xFF1B72EC),
                  foregroundColor: Colors.white,
                  minimumSize: const Size(60, 32),
                ),
                child: const Text('同意'),
              ),
            ])
          : Text(apply.handleResult == 1 ? '已同意' : '已拒绝',
              style: const TextStyle(color: Colors.grey, fontSize: 13)),
    );
  }
}
