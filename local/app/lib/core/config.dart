import 'package:shared_preferences/shared_preferences.dart';

class Config {
  Config._();

  static late SharedPreferences _prefs;

  // 服务端地址（根据实际部署修改）
  static String apiUrl = 'http://42.192.129.44:10002';
  static String wsUrl = 'ws://42.192.129.44:10001';

  static const int platformID = 5; // H5=5, Android=2, iOS=1

  static Future<void> init() async {
    _prefs = await SharedPreferences.getInstance();
    // 从本地读取自定义服务器地址
    apiUrl = _prefs.getString('apiUrl') ?? apiUrl;
    wsUrl = _prefs.getString('wsUrl') ?? wsUrl;
  }

  // 登录凭证（imToken 用于 IM 层，chatToken 用于 Chat 层）
  static String get token => _prefs.getString('token') ?? '';
  static set token(String v) => _prefs.setString('token', v);

  static String get chatToken => _prefs.getString('chatToken') ?? '';
  static set chatToken(String v) => _prefs.setString('chatToken', v);

  static String get userID => _prefs.getString('userID') ?? '';
  static set userID(String v) => _prefs.setString('userID', v);

  static String get nickname => _prefs.getString('nickname') ?? '';
  static set nickname(String v) => _prefs.setString('nickname', v);

  static String get faceURL => _prefs.getString('faceURL') ?? '';
  static set faceURL(String v) => _prefs.setString('faceURL', v);

  static bool get isLoggedIn => token.isNotEmpty && userID.isNotEmpty;

  static Future<void> clearLogin() async {
    await _prefs.remove('token');
    await _prefs.remove('chatToken');
    await _prefs.remove('userID');
    await _prefs.remove('nickname');
    await _prefs.remove('faceURL');
  }

  static Future<void> setServer(String api, String ws) async {
    apiUrl = api;
    wsUrl = ws;
    await _prefs.setString('apiUrl', api);
    await _prefs.setString('wsUrl', ws);
  }
}
