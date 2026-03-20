import 'package:just_audio/just_audio.dart';
import 'package:flutter/material.dart';

/// Audio player manager for voice messages
class AudioPlayerManager {
  static final AudioPlayerManager _instance = AudioPlayerManager._internal();
  factory AudioPlayerManager() => _instance;
  AudioPlayerManager._internal();

  final _player = AudioPlayer();
  String? _currentPlayingUrl;
  final _playingStateNotifier = ValueNotifier<String?>(null);

  ValueNotifier<String?> get playingStateNotifier => _playingStateNotifier;

  bool isPlaying(String url) => _currentPlayingUrl == url && _player.playing;

  Future<void> playOrPause(String url) async {
    try {
      if (_currentPlayingUrl == url) {
        // Same audio, toggle play/pause
        if (_player.playing) {
          await _player.pause();
          _playingStateNotifier.value = null;
        } else {
          await _player.play();
          _playingStateNotifier.value = url;
        }
      } else {
        // Different audio, play new one
        await _player.stop();
        await _player.setUrl(url);
        await _player.play();
        _currentPlayingUrl = url;
        _playingStateNotifier.value = url;

        // Listen for completion
        _player.playerStateStream.listen((state) {
          if (state.processingState == ProcessingState.completed) {
            _playingStateNotifier.value = null;
            _currentPlayingUrl = null;
          }
        });
      }
    } catch (e) {
      debugPrint('[AudioPlayer] Error: $e');
      _playingStateNotifier.value = null;
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
