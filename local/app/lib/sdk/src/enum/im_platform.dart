import 'dart:io';

class IMPlatform {
  static const ios = 1;
  static const android = 2;
  static const windows = 3;
  static const osx = 4;
  static const web = 5;
  static const miniWeb = 6;
  static const linux = 7;
  static const aPad = 8;
  static const iPad = 9;

  static int get current {
    if (Platform.isAndroid) return android;
    if (Platform.isIOS) return ios;
    if (Platform.isWindows) return windows;
    if (Platform.isMacOS) return osx;
    if (Platform.isLinux) return linux;
    return web;
  }
}
