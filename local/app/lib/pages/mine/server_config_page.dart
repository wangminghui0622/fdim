import 'package:flutter/material.dart';
import 'package:flutter_easyloading/flutter_easyloading.dart';

import '../../core/config.dart';
import '../../core/http_client.dart';

class ServerConfigPage extends StatefulWidget {
  const ServerConfigPage({super.key});

  @override
  State<ServerConfigPage> createState() => _ServerConfigPageState();
}

class _ServerConfigPageState extends State<ServerConfigPage> {
  final _apiController = TextEditingController();
  final _wsController = TextEditingController();

  @override
  void initState() {
    super.initState();
    _apiController.text = Config.apiUrl;
    _wsController.text = Config.wsUrl;
  }

  @override
  void dispose() {
    _apiController.dispose();
    _wsController.dispose();
    super.dispose();
  }

  Future<void> _save() async {
    final api = _apiController.text.trim();
    final ws = _wsController.text.trim();
    if (api.isEmpty || ws.isEmpty) {
      EasyLoading.showToast('请填写完整');
      return;
    }
    await Config.setServer(api, ws);
    HttpClient.updateBaseUrl(api);
    EasyLoading.showToast('保存成功，重启生效');
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('服务器配置')),
      body: Padding(
        padding: const EdgeInsets.all(24),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            const Text('API 地址',
                style: TextStyle(fontWeight: FontWeight.w500)),
            const SizedBox(height: 8),
            TextField(
              controller: _apiController,
              decoration: InputDecoration(
                hintText: 'http://ip:port',
                border: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(12)),
              ),
            ),
            const SizedBox(height: 20),
            const Text('WebSocket 地址',
                style: TextStyle(fontWeight: FontWeight.w500)),
            const SizedBox(height: 8),
            TextField(
              controller: _wsController,
              decoration: InputDecoration(
                hintText: 'ws://ip:port',
                border: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(12)),
              ),
            ),
            const SizedBox(height: 32),
            SizedBox(
              height: 48,
              child: ElevatedButton(
                onPressed: _save,
                style: ElevatedButton.styleFrom(
                  backgroundColor: const Color(0xFF1B72EC),
                  foregroundColor: Colors.white,
                  shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(12)),
                ),
                child: const Text('保存', style: TextStyle(fontSize: 16)),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
