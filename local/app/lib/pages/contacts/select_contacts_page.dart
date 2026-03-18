import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../core/apis/friend_api.dart';
import '../../core/models/user_info.dart';

class SelectContactsPage extends StatefulWidget {
  const SelectContactsPage({super.key});

  @override
  State<SelectContactsPage> createState() => _SelectContactsPageState();
}

class _SelectContactsPageState extends State<SelectContactsPage> {
  List<FriendInfo> _friends = [];
  final Set<String> _selected = {};
  bool _loading = true;
  final _searchCtrl = TextEditingController();
  String _keyword = '';

  @override
  void initState() {
    super.initState();
    _load();
  }

  @override
  void dispose() {
    _searchCtrl.dispose();
    super.dispose();
  }

  Future<void> _load() async {
    try {
      final list = await FriendApi.getFriendList();
      if (mounted) setState(() { _friends = list; _loading = false; });
    } catch (_) {
      if (mounted) setState(() => _loading = false);
    }
  }

  List<FriendInfo> get _filtered {
    if (_keyword.isEmpty) return _friends;
    return _friends.where((f) {
      final name = f.showName.toLowerCase();
      final uid = f.friendUserID.toLowerCase();
      final kw = _keyword.toLowerCase();
      return name.contains(kw) || uid.contains(kw);
    }).toList();
  }

  void _confirm() {
    Get.back(result: _selected.toList());
  }

  @override
  Widget build(BuildContext context) {
    final filtered = _filtered;
    return Scaffold(
      appBar: AppBar(
        title: const Text('选择联系人'),
        actions: [
          TextButton(
            onPressed: _selected.isEmpty ? null : _confirm,
            child: Text('确定(${_selected.length})',
                style: const TextStyle(color: Colors.white)),
          ),
        ],
      ),
      body: Column(
        children: [
          Padding(
            padding: const EdgeInsets.all(12),
            child: TextField(
              controller: _searchCtrl,
              decoration: InputDecoration(
                hintText: '搜索',
                prefixIcon: const Icon(Icons.search),
                border: OutlineInputBorder(borderRadius: BorderRadius.circular(24)),
                contentPadding: const EdgeInsets.symmetric(horizontal: 16),
                isDense: true,
              ),
              onChanged: (v) => setState(() => _keyword = v),
            ),
          ),
          Expanded(
            child: _loading
                ? const Center(child: CircularProgressIndicator())
                : ListView.builder(
                    itemCount: filtered.length,
                    itemBuilder: (_, i) {
                      final f = filtered[i];
                      final uid = f.friendUserID;
                      final checked = _selected.contains(uid);
                      return CheckboxListTile(
                        value: checked,
                        onChanged: (v) {
                          setState(() {
                            if (v == true) {
                              _selected.add(uid);
                            } else {
                              _selected.remove(uid);
                            }
                          });
                        },
                        secondary: CircleAvatar(
                          backgroundColor: Colors.grey[300],
                          backgroundImage: f.faceURL.isNotEmpty
                              ? NetworkImage(f.faceURL)
                              : null,
                          child: f.faceURL.isEmpty
                              ? Text(f.showName.isNotEmpty
                                  ? f.showName[0].toUpperCase()
                                  : '?')
                              : null,
                        ),
                        title: Text(f.showName),
                      );
                    },
                  ),
          ),
        ],
      ),
    );
  }
}
