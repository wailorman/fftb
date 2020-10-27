# [](https://github.com/wailorman/fftb/compare/v0.8.0...v) (2020-10-27)



# [0.8.0](https://github.com/wailorman/fftb/compare/v0.7.1...v0.8.0) (2020-10-27)


### Features

* **split:** split video by keyframes ([369efed](https://github.com/wailorman/fftb/commit/369efed37e2bfc7c98fb9bfc3479043dc9737286))



## [0.7.1](https://github.com/wailorman/fftb/compare/v0.7.0...v0.7.1) (2020-10-24)


### Bug Fixes

* **convert:** passing hwaccel cli option ([b6e0c54](https://github.com/wailorman/fftb/commit/b6e0c540116b1a077e5133ec980bfcd1382c2a32))
* **convert/hevc:** set hvc1 tag at codec level ([3ae7029](https://github.com/wailorman/fftb/commit/3ae70293d979b5a52a28fad6adf6e73633b05620))
* **convert/vtb:** exit with error when quality option passed to vtb ([18f0141](https://github.com/wailorman/fftb/commit/18f01419e047d509f392caff203af8c14bec7e3b))



# [0.7.0](https://github.com/wailorman/fftb/compare/v0.6.2...v0.7.0) (2020-10-24)


### Bug Fixes

* **convert:** use ffmpeg -max_muxing_queue_size option ([88769b2](https://github.com/wailorman/fftb/commit/88769b2b34a4540dbbafe84915e163c55aa318e3))


### Features

* **converter:** add keyframe_interval option ([c680159](https://github.com/wailorman/fftb/commit/c68015966544892a244e92f5ddf13f9a7cc34397))



## [0.6.2](https://github.com/wailorman/fftb/compare/v0.6.1...v0.6.2) (2020-10-24)



## [0.6.1](https://github.com/wailorman/fftb/compare/v0.6.0...v0.6.1) (2020-10-24)


### Bug Fixes

* **converter:** video_quality yaml name ([75583db](https://github.com/wailorman/fftb/commit/75583db65192df427265a00e05acb0a8b80d3a30))



# [0.6.0](https://github.com/wailorman/fftb/compare/v0.5.1...v0.6.0) (2020-10-24)


### Features

* **converter:** add nvenc video quality mode ([f83f195](https://github.com/wailorman/fftb/commit/f83f195942ce193b159b06f3fb1af0c8b715ba24))
* **converter:** CRF support for CPU encoding ([6f0bc91](https://github.com/wailorman/fftb/commit/6f0bc91f3df743c8dcf505a235586ee603e3b006))
* **converter:** use copy audio codec ([e7c8cff](https://github.com/wailorman/fftb/commit/e7c8cff768af0ce777cfc7addd75999839d05f4d))
* **converter:** use float as progress meter ([4d4ed69](https://github.com/wailorman/fftb/commit/4d4ed6955275aba141b72777780e7ac475619f81))



## [0.5.1](https://github.com/wailorman/fftb/compare/v0.5.0...v0.5.1) (2020-10-24)


### Bug Fixes

* **chunker:** respect original filename ([f2e1961](https://github.com/wailorman/fftb/commit/f2e1961f91b8827bc085d159ef203f83e4f448cb))
* **files/path:** create directory recursively ([2c6f75d](https://github.com/wailorman/fftb/commit/2c6f75d8784de37f3c291ef64a0ca12fb5206723))



# [0.5.0](https://github.com/wailorman/fftb/compare/v0.4.0...v0.5.0) (2020-10-24)


### Bug Fixes

* **converter:** panic-safe stopping conversion again ([fde8d63](https://github.com/wailorman/fftb/commit/fde8d63e04fc3ab3cbb2ea72100f700fe797d8b4))
* **converters:** panic-safe stopping conversion ([2034e9f](https://github.com/wailorman/fftb/commit/2034e9fa76eb802a0ca5d0b36e2f251f1806deba))


### Features

* **converter:** add yaml config & dry-run option ([6269d51](https://github.com/wailorman/fftb/commit/6269d519d10dc06a415348e2dd1949f994985ab3))
* **converter:** overwrite existing files ([442e62b](https://github.com/wailorman/fftb/commit/442e62bdcdd221c1eb30edd9da725f9e13d90082))
* **recursive converter:** add ids to batch tasks ([1c972dd](https://github.com/wailorman/fftb/commit/1c972dd400890a6478c66a322461f1727c09722a))
* add configurable log level ([76b6534](https://github.com/wailorman/fftb/commit/76b6534414e404fdd5860470c0e641d13aef5d53))



# [0.4.0](https://github.com/wailorman/fftb/compare/v0.1.0...v0.4.0) (2020-10-24)


### Bug Fixes

* **convert:** disable preset for videotoolbox ([9db173b](https://github.com/wailorman/fftb/commit/9db173b18df3d3bca2c5e68700154b9f4bb947cf))
* **convert:** disable tagging for videotoolbox ([863cf0b](https://github.com/wailorman/fftb/commit/863cf0b57541469b824f7c639d0d6b6cded8494e))
* **converter:** hide banner ([a05feb1](https://github.com/wailorman/fftb/commit/a05feb183d8788b5ea598e613501a72d7e83bf86))


### Features

* **convert:** batch converting ([67a33f7](https://github.com/wailorman/fftb/commit/67a33f794d0caf7fdb50a983d1bf99991a2371f3))
* **convert:** scaling option ([c29152a](https://github.com/wailorman/fftb/commit/c29152a05a7186d8ade4654dec1afa05714c862c))
* **converter:** hvc1 tag for hevc ([9d9ea2d](https://github.com/wailorman/fftb/commit/9d9ea2da69247fb0aac7a99a29b7f1744ca31889))
* add converter ([0548742](https://github.com/wailorman/fftb/commit/05487429fdcd4c329631f5d84445b60ca89e7b10))
* Use goffmpeg instead of exec() ([9627348](https://github.com/wailorman/fftb/commit/9627348ebbd11762954c24b327987101cff075f2))



# 0.1.0 (2020-05-23)



