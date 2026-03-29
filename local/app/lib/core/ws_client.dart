import 'dart:async';
import 'dart:convert';

import 'package:flutter/foundation.dart';
import 'package:web_socket_channel/web_socket_channel.dart';

import 'config.dart';

/// Official FDIM WebSocket protocol constants
class WsProtocol {
  // Client → Server request identifiers
  static const int wsGetNewestSeq = 1001;
  static const int wsSendMsg = 1002;
  static const int wsPullMsgBySeq = 1003;
  static const int wsSendSignalMsg = 1004;
  static const int wsSendBinaryMsg = 1005;
  static const int wsSetBackgroundStatus = 1006;

  // Server → Client push identifiers
  static const int wsPushMsg = 2001;
  static const int wsKickOnlineMsg = 2002;
  static const int wsLogoutMsg = 2003;
}

enum WsStatus { connecting, connected, disconnected }

/// Official FDIM WsResp frame
class WsResp {
  final int reqIdentifier;
  final int errCode;
  final String errMsg;
  final String operationID;
  final String msgIncr;
  final dynamic data;

  WsResp({
    this.reqIdentifier = 0,
    this.errCode = 0,
    this.errMsg = '',
    this.operationID = '',
    this.msgIncr = '',
    this.data,
  });

  factory WsResp.fromJson(Map<String, dynamic> json) {
    return WsResp(
      reqIdentifier: json['reqIdentifier'] ?? 0,
      errCode: json['errCode'] ?? 0,
      errMsg: json['errMsg'] ?? '',
      operationID: json['operationID'] ?? '',
      msgIncr: json['msgIncr'] ?? '',
      data: json['data'],
    );
  }
}

class WsClient {
  WsClient._();

  static WebSocketChannel? _channel;
  static final _statusController = StreamController<WsStatus>.broadcast();
  static final _messageController = StreamController<WsResp>.broadcast();
  static Timer? _heartbeatTimer;
  static Timer? _reconnectTimer;
  static bool _manualClose = false;
  static int _msgIncr = 0;

  static Stream<WsStatus> get statusStream => _statusController.stream;
  static Stream<WsResp> get messageStream => _messageController.stream;

  static void connect() {
    if (_channel != null) return;
    _manualClose = false;
    _doConnect();
  }

  static void _doConnect() {
    try {
      final token = Config.token;
      final userID = Config.userID;
      if (token.isEmpty || userID.isEmpty) return;

      final wsUrl =
          '${Config.wsUrl}/ws?sendID=$userID&token=$token&platformID=${Config.platformID}';
      debugPrint('[WS] Connecting to $wsUrl');
      _statusController.add(WsStatus.connecting);

      _channel = WebSocketChannel.connect(Uri.parse(wsUrl));
      _channel!.stream.listen(
        _onMessage,
        onError: _onError,
        onDone: _onDone,
      );

      _statusController.add(WsStatus.connected);
      _startHeartbeat();
    } catch (e) {
      debugPrint('[WS] Connect error: $e');
      _scheduleReconnect();
    }
  }

  static void _onMessage(dynamic raw) {
    try {
      final json = jsonDecode(raw as String) as Map<String, dynamic>;
      final resp = WsResp.fromJson(json);
      if (kDebugMode) {
        debugPrint('[WS] Received reqIdentifier=${resp.reqIdentifier}');
      }
      _messageController.add(resp);
    } catch (e) {
      debugPrint('[WS] Parse error: $e');
    }
  }

  static void _onError(dynamic error) {
    debugPrint('[WS] Error: $error');
    _statusController.add(WsStatus.disconnected);
    _cleanup();
    _scheduleReconnect();
  }

  static void _onDone() {
    debugPrint('[WS] Connection closed');
    _statusController.add(WsStatus.disconnected);
    _cleanup();
    if (!_manualClose) {
      _scheduleReconnect();
    }
  }

  /// Send a WsReq frame (official protocol)
  static void sendReq({
    required int reqIdentifier,
    dynamic data,
    String? operationID,
  }) {
    _msgIncr++;
    final req = {
      'reqIdentifier': reqIdentifier,
      'token': Config.token,
      'sendID': Config.userID,
      'operationID': operationID ?? '',
      'msgIncr': '$_msgIncr',
      'data': data,
    };
    _channel?.sink.add(jsonEncode(req));
  }

  static void _startHeartbeat() {
    _heartbeatTimer?.cancel();
    _heartbeatTimer = Timer.periodic(const Duration(seconds: 30), (_) {
      sendReq(reqIdentifier: WsProtocol.wsSetBackgroundStatus);
    });
  }

  static void _scheduleReconnect() {
    if (_manualClose) return;
    _reconnectTimer?.cancel();
    _reconnectTimer = Timer(const Duration(seconds: 3), () {
      debugPrint('[WS] Reconnecting...');
      _doConnect();
    });
  }

  static void _cleanup() {
    _heartbeatTimer?.cancel();
    _heartbeatTimer = null;
    _channel = null;
  }

  static void disconnect() {
    _manualClose = true;
    _heartbeatTimer?.cancel();
    _reconnectTimer?.cancel();
    _channel?.sink.close();
    _channel = null;
    _statusController.add(WsStatus.disconnected);
  }
}
