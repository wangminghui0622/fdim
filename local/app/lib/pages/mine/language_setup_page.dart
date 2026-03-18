import 'package:flutter/material.dart';
import 'package:get/get.dart';

class LanguageSetupPage extends StatefulWidget {
  const LanguageSetupPage({super.key});

  @override
  State<LanguageSetupPage> createState() => _LanguageSetupPageState();
}

class _LanguageSetupPageState extends State<LanguageSetupPage> {
  String _selected = 'zh_CN';

  final _languages = [
    {'code': 'zh_CN', 'name': '简体中文'},
    {'code': 'en_US', 'name': 'English'},
  ];

  void _select(String code) {
    setState(() => _selected = code);
    final parts = code.split('_');
    if (parts.length == 2) {
      Get.updateLocale(Locale(parts[0], parts[1]));
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('语言设置')),
      body: ListView.separated(
        itemCount: _languages.length,
        separatorBuilder: (_, __) => const Divider(height: 1),
        itemBuilder: (_, i) {
          final lang = _languages[i];
          return ListTile(
            title: Text(lang['name']!),
            trailing: _selected == lang['code']
                ? const Icon(Icons.check, color: Color(0xFF1B72EC))
                : null,
            onTap: () => _select(lang['code']!),
          );
        },
      ),
    );
  }
}
