class OnMsgSendProgressListener {
  Function(String msgID, int progress)? onProgress;

  OnMsgSendProgressListener({this.onProgress});

  void progress(String msgID, int progress) {
    onProgress?.call(msgID, progress);
  }
}
