import 'package:get/get.dart';

import 'app_routes.dart';
import '../pages/splash/splash_page.dart';
import '../pages/login/login_page.dart';
import '../pages/login/register_page.dart';
import '../pages/login/forget_password_page.dart';
import '../pages/home/home_page.dart';
import '../pages/chat/chat_page.dart';
import '../pages/chat/chat_settings_page.dart';
import '../pages/contacts/friend_apply_page.dart';
import '../pages/contacts/add_friend_page.dart';
import '../pages/contacts/user_profile_page.dart';
import '../pages/contacts/select_contacts_page.dart';
import '../pages/contacts/send_verification_page.dart';
import '../pages/group/create_group_page.dart';
import '../pages/group/group_list_page.dart';
import '../pages/group/group_profile_page.dart';
import '../pages/group/group_setup_page.dart';
import '../pages/group/group_requests_page.dart';
import '../pages/search/global_search_page.dart';
import '../pages/mine/my_info_page.dart';
import '../pages/mine/edit_my_info_page.dart';
import '../pages/mine/settings_page.dart';
import '../pages/mine/server_config_page.dart';
import '../pages/mine/change_password_page.dart';
import '../pages/mine/blacklist_page.dart';
import '../pages/mine/about_us_page.dart';
import '../pages/mine/language_setup_page.dart';
import '../pages/mine/my_qrcode_page.dart';

class AppPages {
  static final routes = [
    GetPage(name: AppRoutes.splash, page: () => const SplashPage()),
    GetPage(name: AppRoutes.login, page: () => const LoginPage()),
    GetPage(name: AppRoutes.register, page: () => const RegisterPage()),
    GetPage(name: AppRoutes.forgetPassword, page: () => const ForgetPasswordPage()),
    GetPage(name: AppRoutes.home, page: () => const HomePage()),
    GetPage(name: AppRoutes.chat, page: () => const ChatPage()),
    GetPage(name: AppRoutes.chatSettings, page: () => const ChatSettingsPage()),
    GetPage(name: AppRoutes.friendApply, page: () => const FriendApplyPage()),
    GetPage(name: AppRoutes.addFriend, page: () => const AddFriendPage()),
    GetPage(name: AppRoutes.userProfile, page: () => const UserProfilePage()),
    GetPage(name: AppRoutes.selectContacts, page: () => const SelectContactsPage()),
    GetPage(name: AppRoutes.sendVerification, page: () => const SendVerificationPage()),
    GetPage(name: AppRoutes.createGroup, page: () => const CreateGroupPage()),
    GetPage(name: AppRoutes.groupList, page: () => const GroupListPage()),
    GetPage(name: AppRoutes.groupProfile, page: () => const GroupProfilePage()),
    GetPage(name: AppRoutes.groupSetup, page: () => const GroupSetupPage()),
    GetPage(name: AppRoutes.groupRequests, page: () => const GroupRequestsPage()),
    GetPage(name: AppRoutes.globalSearch, page: () => const GlobalSearchPage()),
    GetPage(name: AppRoutes.myInfo, page: () => const MyInfoPage()),
    GetPage(name: AppRoutes.editMyInfo, page: () => const EditMyInfoPage()),
    GetPage(name: AppRoutes.settings, page: () => const SettingsPage()),
    GetPage(name: AppRoutes.serverConfig, page: () => const ServerConfigPage()),
    GetPage(name: AppRoutes.changePassword, page: () => const ChangePasswordPage()),
    GetPage(name: AppRoutes.blacklist, page: () => const BlacklistPage()),
    GetPage(name: AppRoutes.aboutUs, page: () => const AboutUsPage()),
    GetPage(name: AppRoutes.languageSetup, page: () => const LanguageSetupPage()),
    GetPage(name: AppRoutes.myQrcode, page: () => const MyQrcodePage()),
  ];
}
