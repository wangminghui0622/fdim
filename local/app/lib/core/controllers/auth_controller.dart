import 'dart:convert';
import 'package:crypto/crypto.dart';
import 'package:flutter/foundation.dart';
import 'package:get/get.dart';

import '../config.dart';
import '../http_client.dart';
import '../../sdk/flutter_openim_sdk.dart';
import '../../routes/app_routes.dart';
import 'im_controller.dart';

class AuthController extends GetxController {
  final userInfo = Rx<UserInfo?>(null);
  final isLoggedIn = false.obs;

  @override
  void onInit() {
    super.onInit();
    HttpClient.init();
    _checkLogin();
  }

  static String _md5(String input) {
    return md5.convert(utf8.encode(input)).toString();
  }

  void _checkLogin() {
    if (Config.isLoggedIn) {
      isLoggedIn.value = true;
      _initSDKAndLogin();
    }
  }

  /// Initialize SDK + login via SDK after tokens are available
  Future<void> _initSDKAndLogin() async {
    try {
      final imCtrl = Get.find<IMController>();
      await imCtrl.initOpenIM();
      final info = await imCtrl.loginSDK();
      userInfo.value = info;
      Config.nickname = info.nickname ?? '';
      Config.faceURL = info.faceURL ?? '';
      // Load initial friend apply count for tabbar badge
      imCtrl.loadFriendApplyCount();
    } catch (e) {
      debugPrint('[Auth] SDK init/login error: $e');
    }
  }

  /// 登录（对齐官方 /account/login）
  Future<void> login({
    String? areaCode,
    String? phoneNumber,
    String? account,
    String? email,
    String? password,
    String? verificationCode,
  }) async {
    final data = await HttpClient.post('/account/login', data: {
      'areaCode': areaCode,
      'phoneNumber': phoneNumber,
      'account': account,
      'email': email,
      'password': password != null ? _md5(password) : null,
      'verifyCode': verificationCode,
      'platform': Config.platformID,
    });

    final imToken = (data['imToken'] ?? data['token'] ?? '') as String;
    final chatToken = (data['chatToken'] ?? '') as String;
    final userID = (data['userID'] ?? '') as String;

    if (imToken.isEmpty || userID.isEmpty) {
      throw ApiException(500, '登录返回数据异常');
    }

    Config.token = imToken;
    Config.chatToken = chatToken;
    Config.userID = userID;
    isLoggedIn.value = true;

    // Initialize SDK and login via SDK
    await _initSDKAndLogin();
  }

  /// 注册（对齐官方 /account/register）
  Future<void> register({
    String? areaCode,
    String? phoneNumber,
    String? email,
    String? account,
    required String password,
    required String nickname,
    String? faceURL,
    int gender = 1,
    int birth = 0,
    String? verificationCode,
    String? invitationCode,
  }) async {
    final data = await HttpClient.post('/account/register', data: {
      'verifyCode': verificationCode ?? '',
      'platform': Config.platformID,
      'invitationCode': invitationCode,
      'autoLogin': true,
      'user': {
        'nickname': nickname,
        'faceURL': faceURL ?? '',
        'areaCode': areaCode,
        'phoneNumber': phoneNumber,
        'email': email,
        'account': account,
        'password': _md5(password),
        'birth': birth,
        'gender': gender,
      },
    });

    final imToken = (data['imToken'] ?? data['token'] ?? '') as String;
    final chatToken = (data['chatToken'] ?? '') as String;
    final userID = (data['userID'] ?? '') as String;

    if (imToken.isNotEmpty && userID.isNotEmpty) {
      Config.token = imToken;
      Config.chatToken = chatToken;
      Config.userID = userID;
      isLoggedIn.value = true;
      await _initSDKAndLogin();
    } else {
      await login(phoneNumber: phoneNumber, password: password);
    }
  }

  /// 登出
  Future<void> logout() async {
    try {
      final imCtrl = Get.find<IMController>();
      await imCtrl.logoutSDK();
    } catch (_) {}
    await Config.clearLogin();
    isLoggedIn.value = false;
    userInfo.value = null;
    Get.offAllNamed(AppRoutes.login);
  }

  /// 修改密码
  Future<void> changePassword({
    required String currentPassword,
    required String newPassword,
  }) async {
    await HttpClient.post('/account/password/change', data: {
      'userID': Config.userID,
      'currentPassword': _md5(currentPassword),
      'newPassword': _md5(newPassword),
      'platform': Config.platformID,
    });
  }

  /// 重置密码
  Future<void> resetPassword({
    String? areaCode,
    String? phoneNumber,
    String? email,
    required String newPassword,
    required String verificationCode,
  }) async {
    await HttpClient.post('/account/password/reset', data: {
      'areaCode': areaCode,
      'phoneNumber': phoneNumber,
      'email': email,
      'password': _md5(newPassword),
      'verifyCode': verificationCode,
      'platform': Config.platformID,
    });
  }

  /// 发送验证码
  Future<void> sendVerifyCode({
    String? areaCode,
    String? phoneNumber,
    String? email,
    required int usedFor,
    String? invitationCode,
  }) async {
    await HttpClient.post('/account/code/send', data: {
      'areaCode': areaCode,
      'phoneNumber': phoneNumber,
      'email': email,
      'usedFor': usedFor,
      'invitationCode': invitationCode,
    });
  }

  /// 校验验证码
  Future<void> checkVerifyCode({
    String? areaCode,
    String? phoneNumber,
    String? email,
    required String verificationCode,
    required int usedFor,
  }) async {
    await HttpClient.post('/account/code/verify', data: {
      'areaCode': areaCode,
      'phoneNumber': phoneNumber,
      'email': email,
      'verifyCode': verificationCode,
      'usedFor': usedFor,
    });
  }

  /// 更新用户信息（chat层 + IM层 + SDK）
  Future<void> updateUserInfo({
    String? nickname,
    String? faceURL,
    int? gender,
    int? birth,
    String? email,
    String? account,
  }) async {
    // Chat layer
    final chatData = <String, dynamic>{'userID': Config.userID};
    if (nickname != null) chatData['nickname'] = nickname;
    if (faceURL != null) chatData['faceURL'] = faceURL;
    if (gender != null) chatData['gender'] = gender;
    if (birth != null) chatData['birth'] = birth;
    if (email != null) chatData['email'] = email;
    if (account != null) chatData['account'] = account;

    await HttpClient.post('/user/update', data: {
      ...chatData,
      'platform': Config.platformID,
    });

    // IM layer via SDK
    await OpenIM.iMManager.userManager.setSelfInfo(
      nickname: nickname,
      faceURL: faceURL,
    );

    // 更新本地响应式状态，触发 Obx UI 刷新
    final current = userInfo.value;
    if (current != null) {
      if (nickname != null) current.nickname = nickname;
      if (faceURL != null) current.faceURL = faceURL;
      if (gender != null) current.gender = gender;
      if (email != null) current.email = email;
      userInfo.refresh();
    }
  }
}
