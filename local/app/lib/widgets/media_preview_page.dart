import 'package:flutter/material.dart';
import 'package:photo_view/photo_view.dart';
import 'package:video_player/video_player.dart';
import 'package:chewie/chewie.dart';

/// Media preview page for images and videos
class MediaPreviewPage extends StatefulWidget {
  final String url;
  final bool isVideo;
  final String? heroTag;

  const MediaPreviewPage({
    super.key,
    required this.url,
    this.isVideo = false,
    this.heroTag,
  });

  @override
  State<MediaPreviewPage> createState() => _MediaPreviewPageState();
}

class _MediaPreviewPageState extends State<MediaPreviewPage> {
  VideoPlayerController? _videoController;
  ChewieController? _chewieController;
  bool _isLoading = true;
  String? _error;

  @override
  void initState() {
    super.initState();
    if (widget.isVideo) {
      _initVideoPlayer();
    } else {
      _isLoading = false;
    }
  }

  Future<void> _initVideoPlayer() async {
    try {
      final uri = Uri.parse(widget.url);
      _videoController = VideoPlayerController.networkUrl(uri);
      await _videoController!.initialize();

      _chewieController = ChewieController(
        videoPlayerController: _videoController!,
        autoPlay: true,
        looping: false,
        aspectRatio: _videoController!.value.aspectRatio,
        placeholder: Container(
          color: Colors.black,
          child: const Center(
            child: CircularProgressIndicator(color: Colors.white),
          ),
        ),
        errorBuilder: (context, errorMessage) {
          return Center(
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                const Icon(Icons.error, color: Colors.white, size: 48),
                const SizedBox(height: 16),
                Text('视频加载失败', style: const TextStyle(color: Colors.white)),
              ],
            ),
          );
        },
      );

      setState(() {
        _isLoading = false;
      });
    } catch (e) {
      setState(() {
        _isLoading = false;
        _error = e.toString();
      });
    }
  }

  @override
  void dispose() {
    _videoController?.dispose();
    _chewieController?.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.black,
      appBar: AppBar(
        backgroundColor: Colors.transparent,
        elevation: 0,
        iconTheme: const IconThemeData(color: Colors.white),
      ),
      body: _buildBody(),
    );
  }

  Widget _buildBody() {
    if (_isLoading) {
      return const Center(
        child: CircularProgressIndicator(color: Colors.white),
      );
    }

    if (_error != null) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Icon(Icons.error, color: Colors.white, size: 48),
            const SizedBox(height: 16),
            const Text('加载失败', style: TextStyle(color: Colors.white)),
          ],
        ),
      );
    }

    if (widget.isVideo) {
      return Center(
        child: _chewieController != null
            ? Chewie(controller: _chewieController!)
            : const Text('视频加载失败', style: TextStyle(color: Colors.white)),
      );
    }

    // Image preview with zoom capability
    return GestureDetector(
      onTap: () => Navigator.pop(context),
      child: PhotoView(
        imageProvider: NetworkImage(widget.url),
        backgroundDecoration: const BoxDecoration(color: Colors.black),
        minScale: PhotoViewComputedScale.contained,
        maxScale: PhotoViewComputedScale.covered * 3,
        initialScale: PhotoViewComputedScale.contained,
        heroAttributes: widget.heroTag != null
            ? PhotoViewHeroAttributes(tag: widget.heroTag!)
            : null,
        loadingBuilder: (context, event) => Center(
          child: CircularProgressIndicator(
            value: event?.expectedTotalBytes != null
                ? event!.cumulativeBytesLoaded / event.expectedTotalBytes!
                : null,
            color: Colors.white,
          ),
        ),
        errorBuilder: (context, error, stackTrace) => const Center(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Icon(Icons.broken_image, color: Colors.white, size: 48),
              SizedBox(height: 16),
              Text('图片加载失败', style: TextStyle(color: Colors.white)),
            ],
          ),
        ),
      ),
    );
  }
}
