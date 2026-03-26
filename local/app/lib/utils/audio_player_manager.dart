import 'dart:io';

import 'package:just_audio/just_audio.dart';
import 'package:flutter/material.dart';

/// Audio player manager for voice messages
class AudioPlayerManager {
  static final AudioPlayerManager _instance = AudioPlayerManager._internal();
  factory AudioPlayerManager() => _instance;
  AudioPlayerManager._internal() {
    // 双重保险：stream listener + play().then() 都会清除波浪
    _player.playerStateStream.listen((state) {
      debugPrint('[AudioPlayer] state: playing=${state.playing}, processing=${state.processingState}');
      if (state.processingState == ProcessingState.completed) {
        _clearPlaying();
      }
    });
  }

  void _clearPlaying() {
    debugPrint('[AudioPlayer] clearing wave animation');
    _playingStateNotifier.value = null;
    _currentPlayingUrl = null;
  }

  final _player = AudioPlayer();
  String? _currentPlayingUrl;
  final _playingStateNotifier = ValueNotifier<String?>(null);

  ValueNotifier<String?> get playingStateNotifier => _playingStateNotifier;

  bool isPlaying(String url) => _currentPlayingUrl == url && _player.playing;

  /// 播放/暂停，返回 false 表示播放失败（文件不存在等）
  Future<bool> playOrPause(String url) async {
    try {
      if (_currentPlayingUrl == url) {
        if (_player.playing) {
          // 正在播放 → 暂停
          await _player.pause();
          _playingStateNotifier.value = null;
        } else {
          // 已暂停或播完 → 从头播放
          await _player.seek(Duration.zero);
          _playingStateNotifier.value = url;
          _player.play().then((_) => _clearPlaying());
        }
      } else {
        // 切换到新语音
        await _player.stop();
        await _setAudioSource(url);
        _currentPlayingUrl = url;
        _playingStateNotifier.value = url;
        _player.play().then((_) => _clearPlaying());
      }
      return true;
    } catch (e) {
      debugPrint('[AudioPlayer] Error playing "$url": $e');
      _playingStateNotifier.value = null;
      _currentPlayingUrl = null;
      return false;
    }
  }

  /// 根据路径类型选择正确的加载方式
  Future<void> _setAudioSource(String url) async {
    if (url.startsWith('http://') || url.startsWith('https://')) {
      await _player.setUrl(url);
    } else if (await File(url).exists()) {
      await _player.setFilePath(url);
    } else {
      throw Exception('音频文件不存在: $url');
    }
  }

  Future<void> stop() async {
    await _player.stop();
    _currentPlayingUrl = null;
    _playingStateNotifier.value = null;
  }

  void dispose() {
    _player.dispose();
  }
}
