import 'package:flutter/material.dart';

class AboutUsPage extends StatelessWidget {
  const AboutUsPage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('关于我们')),
      body: ListView(
        padding: const EdgeInsets.all(24),
        children: [
          const SizedBox(height: 40),
          Center(
            child: Container(
              width: 80,
              height: 80,
              decoration: BoxDecoration(
                color: const Color(0xFF1B72EC),
                borderRadius: BorderRadius.circular(20),
              ),
              child: const Icon(Icons.chat_bubble_rounded,
                  color: Colors.white, size: 40),
            ),
          ),
          const SizedBox(height: 16),
          const Center(
            child: Text('OpenIM',
                style: TextStyle(fontSize: 22, fontWeight: FontWeight.bold)),
          ),
          const SizedBox(height: 4),
          const Center(
            child: Text('v1.0.0',
                style: TextStyle(fontSize: 14, color: Colors.grey)),
          ),
          const SizedBox(height: 40),
          const ListTile(
            leading: Icon(Icons.language),
            title: Text('官方网站'),
            subtitle: Text('https://www.openim.io'),
          ),
          const Divider(height: 1),
          const ListTile(
            leading: Icon(Icons.code),
            title: Text('开源地址'),
            subtitle: Text('https://github.com/openimsdk'),
          ),
          const Divider(height: 1),
          const ListTile(
            leading: Icon(Icons.description),
            title: Text('开源协议'),
            subtitle: Text('Apache License 2.0'),
          ),
        ],
      ),
    );
  }
}
