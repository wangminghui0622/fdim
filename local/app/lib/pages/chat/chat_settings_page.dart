import 'package:flutter/material.dart';
import 'package:flutter_easyloading/flutter_easyloading.dart';
import 'package:get/get.dart';

import '../../core/apis/group_api.dart';
import '../../core/controllers/im_controller.dart';
import '../../core/models/group_info.dart';
import '../../models/signaling_info.dart';
import '../../sdk/flutter_openim_sdk.dart' hide GroupInfo;
import '../../routes/app_routes.dart';
import '../../widgets/group_avatar_widget.dart';

class ChatSettingsPage extends StatefulWidget {
  const ChatSettingsPage({super.key});

  @override
  State<ChatSettingsPage> createState() => _ChatSettingsPageState();
}

class _ChatSettingsPageState extends State<ChatSettingsPage> {
  late final String _conversationID;
  late final String _userID;
  late final String _groupID;
  late final int _sessionType;
  bool _pinned = false;
  int _recvMsgOpt = 0; // 0: normal, 1: no notify, 2: no recv
  GroupInfo? _groupInfo;
  List<GroupMemberInfo> _members = [];

  @override
  void initState() {
    super.initState();
    _conversationID = Get.arguments?['conversationID'] ?? '';
    _userID = Get.arguments?['userID'] ?? '';
    _groupID = Get.arguments?['groupID'] ?? '';
    _sessionType = Get.arguments?['sessionType'] ?? 1;
    _loadSettings();
  }

  Future<void> _loadSettings() async {
    try {
      final conv = await OpenIM.iMManager.conversationManager
          .getOneConversation(sessionType: _sessionType, sourceID: _groupID.isNotEmpty ? _groupID : _userID);
      if (mounted) {
        setState(() {
          _pinned = conv.isPinned ?? false;
          _recvMsgOpt = conv.recvMsgOpt ?? 0;
        });
      }
    } catch (_) {}

    if (_groupID.isNotEmpty) {
      try {
        final info = await GroupApi.getGroupsInfo(_groupID);
        final members = await GroupApi.getGroupMemberList(_groupID, showNumber: 10);
        if (mounted) {
          setState(() {
            _groupInfo = info;
            _members = members;
          });
        }
      } catch (_) {}
    }
  }

  Future<void> _togglePin(bool v) async {
    try {
      await OpenIM.iMManager.conversationManager.pinConversation(
        conversationID: _conversationID,
        isPinned: v,
      );
      setState(() => _pinned = v);
    } catch (_) {
      EasyLoading.showToast('操作失败');
    }
  }

  Future<void> _setRecvMsgOpt(int opt) async {
    try {
      await OpenIM.iMManager.conversationManager.setConversation(
        _conversationID,
        ConversationReq(recvMsgOpt: opt),
      );
      setState(() => _recvMsgOpt = opt);
    } catch (_) {
      EasyLoading.showToast('操作失败');
    }
  }

  Future<void> _clearChat() async {
    final confirm = await Get.dialog<bool>(AlertDialog(
      title: const Text('清空聊天记录'),
      content: const Text('确定要清空聊天记录吗？此操作不可恢复。'),
      actions: [
        TextButton(onPressed: () => Get.back(result: false), child: const Text('取消')),
        TextButton(onPressed: () => Get.back(result: true), child: const Text('确定')),
      ],
    ));
    if (confirm != true) return;
    try {
      await OpenIM.iMManager.conversationManager.clearConversationAndDeleteAllMsg(
        conversationID: _conversationID,
      );
      EasyLoading.showToast('已清空');
    } catch (_) {
      EasyLoading.showToast('操作失败');
    }
  }

  void _startCall(CallType callType) {
    if (_userID.isEmpty) return;
    final imCtrl = Get.find<IMController>();
    if (imCtrl.isRtcBusy) {
      EasyLoading.showToast('当前正在通话中');
      return;
    }
    imCtrl.call(callType: callType, inviteeUserIDList: [_userID]);
  }

  @override
  Widget build(BuildContext context) {
    final isGroup = _groupID.isNotEmpty;
    return Scaffold(
      appBar: AppBar(title: const Text('聊天设置')),
      body: ListView(
        children: [
          if (isGroup && _groupInfo != null) ...[
            ListTile(
              leading: _members.isNotEmpty
                  ? GroupAvatarWidget(
                      size: 48,
                      faceURLs: _members.take(4).map((m) => m.faceURL).toList(),
                      names: _members.take(4).map((m) => m.nickname).toList(),
                    )
                  : null,
              title: Text(_groupInfo!.groupName,
                  style: const TextStyle(fontWeight: FontWeight.bold)),
              subtitle: Text('${_groupInfo!.memberCount} 名成员'),
              trailing: const Icon(Icons.chevron_right),
              onTap: () => Get.toNamed(AppRoutes.groupProfile,
                  arguments: {'groupID': _groupID}),
            ),
            if (_members.isNotEmpty)
              SizedBox(
                height: 80,
                child: ListView.builder(
                  scrollDirection: Axis.horizontal,
                  padding: const EdgeInsets.symmetric(horizontal: 12),
                  itemCount: _members.length,
                  itemBuilder: (_, i) {
                    final m = _members[i];
                    return Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 6),
                      child: Column(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          CircleAvatar(
                            radius: 22,
                            backgroundColor: Colors.grey[300],
                            backgroundImage: m.faceURL.isNotEmpty
                                ? NetworkImage(m.faceURL)
                                : null,
                            child: m.faceURL.isEmpty
                                ? Text(m.nickname.isNotEmpty
                                    ? m.nickname[0].toUpperCase()
                                    : '?',
                                    style: const TextStyle(fontSize: 12))
                                : null,
                          ),
                          const SizedBox(height: 4),
                          SizedBox(
                            width: 48,
                            child: Text(
                              m.nickname,
                              maxLines: 1,
                              overflow: TextOverflow.ellipsis,
                              textAlign: TextAlign.center,
                              style: const TextStyle(fontSize: 11),
                            ),
                          ),
                        ],
                      ),
                    );
                  },
                ),
              ),
            const Divider(),
          ],
          SwitchListTile(
            title: const Text('置顶聊天'),
            value: _pinned,
            onChanged: _togglePin,
          ),
          ListTile(
            title: const Text('消息免打扰'),
            trailing: DropdownButton<int>(
              value: _recvMsgOpt,
              underline: const SizedBox(),
              items: const [
                DropdownMenuItem(value: 0, child: Text('正常')),
                DropdownMenuItem(value: 1, child: Text('不提醒')),
                DropdownMenuItem(value: 2, child: Text('屏蔽')),
              ],
              onChanged: (v) {
                if (v != null) _setRecvMsgOpt(v);
              },
            ),
          ),
          if (!isGroup && _userID.isNotEmpty) ...[          const Divider(),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
              child: Row(
                children: [
                  Expanded(
                    child: OutlinedButton.icon(
                      icon: const Icon(Icons.phone_outlined),
                      label: const Text('语音通话'),
                      onPressed: () => _startCall(CallType.audio),
                    ),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: OutlinedButton.icon(
                      icon: const Icon(Icons.videocam_outlined),
                      label: const Text('视频通话'),
                      onPressed: () => _startCall(CallType.video),
                    ),
                  ),
                ],
              ),
            ),
          ],
          const Divider(),
          ListTile(
            title: const Text('清空聊天记录', style: TextStyle(color: Colors.red)),
            onTap: _clearChat,
          ),
        ],
      ),
    );
  }
}
