import '../models/signaling_info.dart';
import 'http_client.dart';

class RtcApi {
  RtcApi._();

  static Future<SignalingCertificate> getTokenForRTC(
      String room, String identity) async {
    final data = await HttpClient.post(
      '/user/rtc/get_token',
      data: {'room': room, 'identity': identity},
    );
    return SignalingCertificate.fromJson(data as Map<String, dynamic>);
  }
}
