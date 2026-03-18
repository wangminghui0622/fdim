class OnConnectListener {
  Function(int? code, String? errorMsg)? onConnectFailed;
  Function()? onConnectSuccess;
  Function()? onConnecting;
  Function()? onKickedOffline;
  Function()? onUserTokenExpired;
  Function()? onUserTokenInvalid;

  OnConnectListener({
    this.onConnectFailed,
    this.onConnectSuccess,
    this.onConnecting,
    this.onKickedOffline,
    this.onUserTokenExpired,
    this.onUserTokenInvalid,
  });

  void connectFailed(int? code, String? errorMsg) {
    onConnectFailed?.call(code, errorMsg);
  }

  void connectSuccess() {
    onConnectSuccess?.call();
  }

  void connecting() {
    onConnecting?.call();
  }

  void kickedOffline() {
    onKickedOffline?.call();
  }

  void userTokenExpired() {
    onUserTokenExpired?.call();
  }

  void userTokenInvalid() {
    onUserTokenInvalid?.call();
  }
}
