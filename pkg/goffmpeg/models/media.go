package models

import (
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
)

// Mediafile _
type Mediafile struct {
	unsafe                bool
	aspect                string
	resolution            string
	videoBitRate          string
	videoBitRateTolerance int
	videoMaxBitRate       int
	videoMinBitrate       int
	videoCodec            string
	vframes               int
	frameRate             int
	audioRate             int
	maxKeyframe           int
	minKeyframe           int
	keyframeInterval      int
	audioCodec            string
	audioBitrate          string
	audioChannels         int
	audioVariableBitrate  bool
	bufferSize            int
	threadset             bool
	threads               int
	preset                string
	tune                  string
	audioProfile          string
	videoProfile          string
	target                string
	duration              string
	durationInput         string
	seekTime              string
	qscale                uint32
	crf                   uint32
	strict                int
	muxDelay              string
	seekUsingTsInput      bool
	seekTimeInput         string
	inputPath             string
	inputPipe             bool
	inputPipeReader       *io.PipeReader
	inputPipeWriter       *io.PipeWriter
	outputPipe            bool
	outputPipeReader      *io.PipeReader
	outputPipeWriter      *io.PipeWriter
	movFlags              string
	hideBanner            bool
	outputPath            string
	outputFormat          string
	inputFormat           string
	copyTs                bool
	nativeFramerateInput  bool
	inputInitialOffset    string
	rtmpLive              string
	hlsPlaylistType       string
	hlsListSize           int
	hlsSegmentDuration    int
	hlsMasterPlaylistName string
	hlsSegmentFilename    string
	httpMethod            string
	httpKeepAlive         bool
	hwaccel               string
	vsync                 bool
	inputVideoCodec       string
	streamIds             map[int]string
	metadata              Metadata
	videoFilter           string
	audioFilter           string
	skipVideo             bool
	skipAudio             bool
	compressionLevel      int
	mapMetadata           string
	tags                  map[string]string
	encryptionKey         string
	movflags              string
	bframe                int
	pixFmt                string
	fileSizeLimit         string
	videoTag              string
	// https://superuser.com/a/1296511
	nvencRateControl string
	// https://stackoverflow.com/questions/49686244/ffmpeg-too-many-packets-buffered-for-output-stream-01
	maxMuxingQueueSize       int
	useWallclockAsTimestamps bool
	constantQuantization     int
	nvencTargetQuality       int
	libx265Params            *Libx265Params
	mapFlag                  string
	segmentTime              int
	resetTimestamps          bool
}

// Libx265Params _
type Libx265Params struct {
	CRF uint32
}

/*** SETTERS ***/

// SetUnsafe _
func (m *Mediafile) SetUnsafe(v bool) {
	m.unsafe = v
}

// SetAudioFilter _
func (m *Mediafile) SetAudioFilter(v string) {
	m.audioFilter = v
}

// SetVideoFilter _
func (m *Mediafile) SetVideoFilter(v string) {
	m.videoFilter = v
}

// SetFilter Deprecated: Use SetVideoFilter instead.
func (m *Mediafile) SetFilter(v string) {
	m.SetVideoFilter(v)
}

// SetAspect _
func (m *Mediafile) SetAspect(v string) {
	m.aspect = v
}

// SetResolution _
func (m *Mediafile) SetResolution(v string) {
	m.resolution = v
}

// SetVideoBitRate _
func (m *Mediafile) SetVideoBitRate(v string) {
	m.videoBitRate = v
}

// SetVideoBitRateTolerance _
func (m *Mediafile) SetVideoBitRateTolerance(v int) {
	m.videoBitRateTolerance = v
}

// SetVideoMaxBitrate _
func (m *Mediafile) SetVideoMaxBitrate(v int) {
	m.videoMaxBitRate = v
}

// SetVideoMinBitRate _
func (m *Mediafile) SetVideoMinBitRate(v int) {
	m.videoMinBitrate = v
}

// SetVideoCodec _
func (m *Mediafile) SetVideoCodec(v string) {
	m.videoCodec = v
}

// SetVframes _
func (m *Mediafile) SetVframes(v int) {
	m.vframes = v
}

// SetFrameRate _
func (m *Mediafile) SetFrameRate(v int) {
	m.frameRate = v
}

// SetAudioRate _
func (m *Mediafile) SetAudioRate(v int) {
	m.audioRate = v
}

// SetAudioVariableBitrate _
func (m *Mediafile) SetAudioVariableBitrate() {
	m.audioVariableBitrate = true
}

// SetMaxKeyFrame _
func (m *Mediafile) SetMaxKeyFrame(v int) {
	m.maxKeyframe = v
}

// SetMinKeyFrame _
func (m *Mediafile) SetMinKeyFrame(v int) {
	m.minKeyframe = v
}

// SetKeyframeInterval _
func (m *Mediafile) SetKeyframeInterval(v int) {
	m.keyframeInterval = v
}

// SetAudioCodec _
func (m *Mediafile) SetAudioCodec(v string) {
	m.audioCodec = v
}

// SetAudioBitRate _
func (m *Mediafile) SetAudioBitRate(v string) {
	m.audioBitrate = v
}

// SetAudioChannels _
func (m *Mediafile) SetAudioChannels(v int) {
	m.audioChannels = v
}

// SetPixFmt _
func (m *Mediafile) SetPixFmt(v string) {
	m.pixFmt = v
}

// SetBufferSize _
func (m *Mediafile) SetBufferSize(v int) {
	m.bufferSize = v
}

// SetThreads _
func (m *Mediafile) SetThreads(v int) {
	m.threadset = true
	m.threads = v
}

// SetPreset _
func (m *Mediafile) SetPreset(v string) {
	m.preset = v
}

// SetTune _
func (m *Mediafile) SetTune(v string) {
	m.tune = v
}

// SetAudioProfile _
func (m *Mediafile) SetAudioProfile(v string) {
	m.audioProfile = v
}

// SetVideoProfile _
func (m *Mediafile) SetVideoProfile(v string) {
	m.videoProfile = v
}

// SetDuration _
func (m *Mediafile) SetDuration(v string) {
	m.duration = v
}

// SetDurationInput _
func (m *Mediafile) SetDurationInput(v string) {
	m.durationInput = v
}

// SetSeekTime _
func (m *Mediafile) SetSeekTime(v string) {
	m.seekTime = v
}

// SetSeekTimeInput _
func (m *Mediafile) SetSeekTimeInput(v string) {
	m.seekTimeInput = v
}

// SetQScale Q Scale must be integer between 1 to 31 - HTTPs://trac.ffmpeg.org/wiki/Encode/MPEG-4
func (m *Mediafile) SetQScale(v uint32) {
	m.qscale = v
}

// SetNvencRateControl _
func (m *Mediafile) SetNvencRateControl(v string) {
	m.nvencRateControl = v
}

// SetMaxMuxingQueueSize _
func (m *Mediafile) SetMaxMuxingQueueSize(v int) {
	m.maxMuxingQueueSize = v
}

// SetConstantQuantization _
func (m *Mediafile) SetConstantQuantization(v int) {
	m.constantQuantization = v
}

// SetNvencTargetQuality _
func (m *Mediafile) SetNvencTargetQuality(v int) {
	m.nvencTargetQuality = v
}

// SetCRF _
func (m *Mediafile) SetCRF(v uint32) {
	m.crf = v
}

// SetLibx265Params _
func (m *Mediafile) SetLibx265Params(v *Libx265Params) {
	m.libx265Params = v
}

// SetStrict _
func (m *Mediafile) SetStrict(v int) {
	m.strict = v
}

// SetSeekUsingTsInput _
func (m *Mediafile) SetSeekUsingTsInput(val bool) {
	m.seekUsingTsInput = val
}

// SetCopyTs _
func (m *Mediafile) SetCopyTs(val bool) {
	m.copyTs = val
}

// SetInputPath _
func (m *Mediafile) SetInputPath(val string) {
	m.inputPath = val
}

// SetInputPipe _
func (m *Mediafile) SetInputPipe(val bool) {
	m.inputPipe = val
}

// SetInputPipeReader _
func (m *Mediafile) SetInputPipeReader(r *io.PipeReader) {
	m.inputPipeReader = r
}

// SetInputPipeWriter _
func (m *Mediafile) SetInputPipeWriter(w *io.PipeWriter) {
	m.inputPipeWriter = w
}

// SetOutputPipe _
func (m *Mediafile) SetOutputPipe(val bool) {
	m.outputPipe = val
}

// SetOutputPipeReader _
func (m *Mediafile) SetOutputPipeReader(r *io.PipeReader) {
	m.outputPipeReader = r
}

// SetOutputPipeWriter _
func (m *Mediafile) SetOutputPipeWriter(w *io.PipeWriter) {
	m.outputPipeWriter = w
}

// SetMovFlags _
func (m *Mediafile) SetMovFlags(val string) {
	m.movFlags = val
}

// SetMap _
func (m *Mediafile) SetMap(val string) {
	m.mapFlag = val
}

// SetSegmentTime _
func (m *Mediafile) SetSegmentTime(val int) {
	m.segmentTime = val
}

// SetResetTimestamps _
func (m *Mediafile) SetResetTimestamps(val bool) {
	m.resetTimestamps = val
}

// SetHideBanner _
func (m *Mediafile) SetHideBanner(val bool) {
	m.hideBanner = val
}

// SetFileSizeLimit _
func (m *Mediafile) SetFileSizeLimit(val string) {
	m.fileSizeLimit = val
}

// SetMuxDelay _
func (m *Mediafile) SetMuxDelay(val string) {
	m.muxDelay = val
}

// SetOutputPath _
func (m *Mediafile) SetOutputPath(val string) {
	m.outputPath = val
}

// SetInputFormat _
func (m *Mediafile) SetInputFormat(val string) {
	m.inputFormat = val
}

// SetOutputFormat _
func (m *Mediafile) SetOutputFormat(val string) {
	m.outputFormat = val
}

// SetNativeFramerateInput _
func (m *Mediafile) SetNativeFramerateInput(val bool) {
	m.nativeFramerateInput = val
}

// SetRtmpLive _
func (m *Mediafile) SetRtmpLive(val string) {
	m.rtmpLive = val
}

// SetHlsListSize _
func (m *Mediafile) SetHlsListSize(val int) {
	m.hlsListSize = val
}

// SetHlsSegmentDuration _
func (m *Mediafile) SetHlsSegmentDuration(val int) {
	m.hlsSegmentDuration = val
}

// SetHlsPlaylistType _
func (m *Mediafile) SetHlsPlaylistType(val string) {
	m.hlsPlaylistType = val
}

// SetHlsMasterPlaylistName _
func (m *Mediafile) SetHlsMasterPlaylistName(val string) {
	m.hlsMasterPlaylistName = val
}

// SetHlsSegmentFilename _
func (m *Mediafile) SetHlsSegmentFilename(val string) {
	m.hlsSegmentFilename = val
}

// SetHTTPMethod _
func (m *Mediafile) SetHTTPMethod(val string) {
	m.httpMethod = val
}

// SetHTTPKeepAlive _
func (m *Mediafile) SetHTTPKeepAlive(val bool) {
	m.httpKeepAlive = val
}

// SetHardwareAcceleration _
func (m *Mediafile) SetHardwareAcceleration(val string) {
	m.hwaccel = val
}

// SetVsync _
func (m *Mediafile) SetVsync(val bool) {
	m.vsync = val
}

// SetInputVideoCodec _
func (m *Mediafile) SetInputVideoCodec(val string) {
	m.inputVideoCodec = val
}

// SetInputInitialOffset _
func (m *Mediafile) SetInputInitialOffset(val string) {
	m.inputInitialOffset = val
}

// SetStreamIds _
func (m *Mediafile) SetStreamIds(val map[int]string) {
	m.streamIds = val
}

// SetSkipVideo _
func (m *Mediafile) SetSkipVideo(val bool) {
	m.skipVideo = val
}

// SetSkipAudio _
func (m *Mediafile) SetSkipAudio(val bool) {
	m.skipAudio = val
}

// SetMetadata _
func (m *Mediafile) SetMetadata(v Metadata) {
	m.metadata = v
}

// SetCompressionLevel _
func (m *Mediafile) SetCompressionLevel(val int) {
	m.compressionLevel = val
}

// SetMapMetadata _
func (m *Mediafile) SetMapMetadata(val string) {
	m.mapMetadata = val
}

// SetTags _
func (m *Mediafile) SetTags(val map[string]string) {
	m.tags = val
}

// SetVideoTag _
func (m *Mediafile) SetVideoTag(val string) {
	m.videoTag = val
}

// SetBframe _
func (m *Mediafile) SetBframe(v int) {
	m.bframe = v
}

/*** GETTERS ***/

// Filter Deprecated: Use VideoFilter instead.
func (m *Mediafile) Filter() string {
	return m.VideoFilter()
}

// VideoFilter _
func (m *Mediafile) VideoFilter() string {
	return m.videoFilter
}

// Unsafe _
func (m *Mediafile) Unsafe() bool {
	return m.unsafe
}

// AudioFilter _
func (m *Mediafile) AudioFilter() string {
	return m.audioFilter
}

// Aspect _
func (m *Mediafile) Aspect() string {
	return m.aspect
}

// Resolution _
func (m *Mediafile) Resolution() string {
	return m.resolution
}

// VideoBitrate _
func (m *Mediafile) VideoBitrate() string {
	return m.videoBitRate
}

// VideoBitRateTolerance _
func (m *Mediafile) VideoBitRateTolerance() int {
	return m.videoBitRateTolerance
}

// VideoMaxBitRate _
func (m *Mediafile) VideoMaxBitRate() int {
	return m.videoMaxBitRate
}

// VideoMinBitRate _
func (m *Mediafile) VideoMinBitRate() int {
	return m.videoMinBitrate
}

// VideoCodec _
func (m *Mediafile) VideoCodec() string {
	return m.videoCodec
}

// Vframes _
func (m *Mediafile) Vframes() int {
	return m.vframes
}

// FrameRate _
func (m *Mediafile) FrameRate() int {
	return m.frameRate
}

// GetPixFmt _
func (m *Mediafile) GetPixFmt() string {
	return m.pixFmt
}

// AudioRate _
func (m *Mediafile) AudioRate() int {
	return m.audioRate
}

// MaxKeyFrame _
func (m *Mediafile) MaxKeyFrame() int {
	return m.maxKeyframe
}

// MinKeyFrame _
func (m *Mediafile) MinKeyFrame() int {
	return m.minKeyframe
}

// KeyFrameInterval _
func (m *Mediafile) KeyFrameInterval() int {
	return m.keyframeInterval
}

// AudioCodec _
func (m *Mediafile) AudioCodec() string {
	return m.audioCodec
}

// AudioBitrate _
func (m *Mediafile) AudioBitrate() string {
	return m.audioBitrate
}

// AudioChannels _
func (m *Mediafile) AudioChannels() int {
	return m.audioChannels
}

// BufferSize _
func (m *Mediafile) BufferSize() int {
	return m.bufferSize
}

// Threads _
func (m *Mediafile) Threads() int {
	return m.threads
}

// Target _
func (m *Mediafile) Target() string {
	return m.target
}

// Duration _
func (m *Mediafile) Duration() string {
	return m.duration
}

// DurationInput _
func (m *Mediafile) DurationInput() string {
	return m.durationInput
}

// SeekTime _
func (m *Mediafile) SeekTime() string {
	return m.seekTime
}

// Preset _
func (m *Mediafile) Preset() string {
	return m.preset
}

// AudioProfile _
func (m *Mediafile) AudioProfile() string {
	return m.audioProfile
}

// VideoProfile _
func (m *Mediafile) VideoProfile() string {
	return m.videoProfile
}

// Tune _
func (m *Mediafile) Tune() string {
	return m.tune
}

// SeekTimeInput _
func (m *Mediafile) SeekTimeInput() string {
	return m.seekTimeInput
}

// QScale _
func (m *Mediafile) QScale() uint32 {
	return m.qscale
}

// NvencRateControl _
func (m *Mediafile) NvencRateControl() string {
	return m.nvencRateControl
}

// MaxMuxingQueueSize _
func (m *Mediafile) MaxMuxingQueueSize() int {
	return m.maxMuxingQueueSize
}

// ConstantQuantization _
func (m *Mediafile) ConstantQuantization() int {
	return m.constantQuantization
}

// NvencTargetQuality _
func (m *Mediafile) NvencTargetQuality() int {
	return m.nvencTargetQuality
}

// CRF _
func (m *Mediafile) CRF() uint32 {
	return m.crf
}

// Libx265Params _
func (m *Mediafile) Libx265Params() *Libx265Params {
	return m.libx265Params
}

// Strict _
func (m *Mediafile) Strict() int {
	return m.strict
}

// MuxDelay _
func (m *Mediafile) MuxDelay() string {
	return m.muxDelay
}

// SeekUsingTsInput _
func (m *Mediafile) SeekUsingTsInput() bool {
	return m.seekUsingTsInput
}

// CopyTs _
func (m *Mediafile) CopyTs() bool {
	return m.copyTs
}

// InputPath _
func (m *Mediafile) InputPath() string {
	return m.inputPath
}

// InputPipe _
func (m *Mediafile) InputPipe() bool {
	return m.inputPipe
}

// InputPipeReader _
func (m *Mediafile) InputPipeReader() *io.PipeReader {
	return m.inputPipeReader
}

// InputPipeWriter _
func (m *Mediafile) InputPipeWriter() *io.PipeWriter {
	return m.inputPipeWriter
}

// OutputPipe _
func (m *Mediafile) OutputPipe() bool {
	return m.outputPipe
}

// OutputPipeReader _
func (m *Mediafile) OutputPipeReader() *io.PipeReader {
	return m.outputPipeReader
}

// OutputPipeWriter _
func (m *Mediafile) OutputPipeWriter() *io.PipeWriter {
	return m.outputPipeWriter
}

// MovFlags _
func (m *Mediafile) MovFlags() string {
	return m.movFlags
}

// Map _
func (m *Mediafile) Map() string {
	return m.mapFlag
}

// SegmentTime _
func (m *Mediafile) SegmentTime() int {
	return m.segmentTime
}

// ResetTimestamps _
func (m *Mediafile) ResetTimestamps() bool {
	return m.resetTimestamps
}

// HideBanner _
func (m *Mediafile) HideBanner() bool {
	return m.hideBanner
}

// FileSizeLimit _
func (m *Mediafile) FileSizeLimit() string {
	return m.fileSizeLimit
}

// OutputPath _
func (m *Mediafile) OutputPath() string {
	return m.outputPath
}

// InputFormat _
func (m *Mediafile) InputFormat() string {
	return m.inputFormat
}

// OutputFormat _
func (m *Mediafile) OutputFormat() string {
	return m.outputFormat
}

// NativeFramerateInput _
func (m *Mediafile) NativeFramerateInput() bool {
	return m.nativeFramerateInput
}

// RtmpLive _
func (m *Mediafile) RtmpLive() string {
	return m.rtmpLive
}

// HlsListSize _
func (m *Mediafile) HlsListSize() int {
	return m.hlsListSize
}

// HlsSegmentDuration _
func (m *Mediafile) HlsSegmentDuration() int {
	return m.hlsSegmentDuration
}

// HlsMasterPlaylistName _
func (m *Mediafile) HlsMasterPlaylistName() string {
	return m.hlsMasterPlaylistName
}

// HlsSegmentFilename _
func (m *Mediafile) HlsSegmentFilename() string {
	return m.hlsSegmentFilename
}

// HlsPlaylistType _
func (m *Mediafile) HlsPlaylistType() string {
	return m.hlsPlaylistType
}

// InputInitialOffset _
func (m *Mediafile) InputInitialOffset() string {
	return m.inputInitialOffset
}

// HTTPMethod _
func (m *Mediafile) HTTPMethod() string {
	return m.httpMethod
}

// HTTPKeepAlive _
func (m *Mediafile) HTTPKeepAlive() bool {
	return m.httpKeepAlive
}

// HardwareAcceleration _
func (m *Mediafile) HardwareAcceleration() string {
	return m.hwaccel
}

// StreamIds _
func (m *Mediafile) StreamIds() map[int]string {
	return m.streamIds
}

// SkipVideo _
func (m *Mediafile) SkipVideo() bool {
	return m.skipVideo
}

// SkipAudio _
func (m *Mediafile) SkipAudio() bool {
	return m.skipAudio
}

// Metadata _
func (m *Mediafile) Metadata() Metadata {
	return m.metadata
}

// CompressionLevel _
func (m *Mediafile) CompressionLevel() int {
	return m.compressionLevel
}

// MapMetadata _
func (m *Mediafile) MapMetadata() string {
	return m.mapMetadata
}

// Tags _
func (m *Mediafile) Tags() map[string]string {
	return m.tags
}

// VideoTag _
func (m *Mediafile) VideoTag() string {
	return m.videoTag
}

// SetEncryptionKey _
func (m *Mediafile) SetEncryptionKey(v string) {
	m.encryptionKey = v
}

// EncryptionKey _
func (m *Mediafile) EncryptionKey() string {
	return m.encryptionKey
}

/** OPTS **/

// ToStrCommand _
func (m *Mediafile) ToStrCommand() []string {
	var strCommand []string

	opts := []string{
		"Unsafe",
		"SeekTimeInput",
		"SeekUsingTsInput",
		"NativeFramerateInput",
		"DurationInput",
		"RtmpLive",
		"InputInitialOffset",
		"HardwareAcceleration",
		"Vsync",
		"InputVideoCodec",
		"InputFormat",
		"InputPath",
		"InputPipe",
		"Map",
		"HideBanner",
		"FileSizeLimit",
		"Aspect",
		"Resolution",
		"FrameRate",
		"AudioRate",
		"VideoCodec",
		"Vframes",
		"VideoBitRate",
		"VideoBitRateTolerance",
		"VideoMaxBitRate",
		"VideoMinBitRate",
		"VideoProfile",
		"SkipVideo",
		"AudioCodec",
		"AudioBitRate",
		"AudioChannels",
		"AudioProfile",
		"SkipAudio",
		"CRF",
		"Libx265Params",
		"QScale",
		"NvencRateControl",
		"MaxMuxingQueueSize",
		"ConstantQuantization",
		"NvencTargetQuality",
		"Strict",
		"BufferSize",
		"MuxDelay",
		"Threads",
		"KeyframeInterval",
		"Preset",
		"PixFmt",
		"Tune",
		"Target",
		"SeekTime",
		"Duration",
		"CopyTs",
		"StreamIds",
		"MovFlags",
		"OutputFormat",
		"OutputPipe",
		"HlsListSize",
		"HlsSegmentDuration",
		"HlsPlaylistType",
		"HlsMasterPlaylistName",
		"SegmentTime",
		"HlsSegmentFilename",
		"AudioFilter",
		"VideoFilter",
		"HTTPMethod",
		"HTTPKeepAlive",
		"CompressionLevel",
		"MapMetadata",
		"Tags",
		"VideoTag",
		"EncryptionKey",
		"OutputPath",
		"Bframe",
		"MovFlags",
	}

	for _, name := range opts {
		opt := reflect.ValueOf(m).MethodByName(fmt.Sprintf("Obtain%s", name))
		if (opt != reflect.Value{}) {
			result := opt.Call([]reflect.Value{})

			if val, ok := result[0].Interface().([]string); ok {
				strCommand = append(strCommand, val...)
			}
		}
	}

	return strCommand
}

// ObtainUnsafe _
func (m *Mediafile) ObtainUnsafe() []string {
	if m.unsafe {
		return []string{"-safe", "0"}
	}

	return nil
}

// ObtainAudioFilter _
func (m *Mediafile) ObtainAudioFilter() []string {
	if m.audioFilter != "" {
		return []string{"-af", m.audioFilter}
	}

	return nil
}

// ObtainVideoFilter _
func (m *Mediafile) ObtainVideoFilter() []string {
	if m.videoFilter != "" {
		return []string{"-vf", m.videoFilter}
	}

	return nil
}

// ObtainAspect _
func (m *Mediafile) ObtainAspect() []string {
	// Set aspect
	if m.resolution != "" {
		resolution := strings.Split(m.resolution, "x")
		if len(resolution) != 0 {
			width, _ := strconv.ParseFloat(resolution[0], 64)
			height, _ := strconv.ParseFloat(resolution[1], 64)
			return []string{"-aspect", fmt.Sprintf("%f", width/height)}
		}
	}

	if m.aspect != "" {
		return []string{"-aspect", m.aspect}
	}

	return nil
}

// ObtainHardwareAcceleration _
func (m *Mediafile) ObtainHardwareAcceleration() []string {
	if m.hwaccel != "" {
		return []string{"-hwaccel", m.hwaccel}
	}

	return nil
}

// ObtainVsync _
func (m *Mediafile) ObtainVsync() []string {
	if m.vsync {
		return []string{"-vsync", "0"}
	}

	return nil
}

// ObtainInputVideoCodec _
func (m *Mediafile) ObtainInputVideoCodec() []string {
	if m.inputVideoCodec != "" {
		return []string{"-c:v", m.inputVideoCodec}
	}

	return nil
}

// ObtainInputPath _
func (m *Mediafile) ObtainInputPath() []string {
	if m.inputPath != "" {
		return []string{"-i", m.inputPath}
	}

	return nil
}

// ObtainInputPipe _
func (m *Mediafile) ObtainInputPipe() []string {
	if m.inputPipe {
		return []string{"-i", "pipe:0"}
	}

	return nil
}

// ObtainOutputPipe _
func (m *Mediafile) ObtainOutputPipe() []string {
	if m.outputPipe {
		return []string{"pipe:1"}
	}

	return nil
}

// ObtainMovFlags _
func (m *Mediafile) ObtainMovFlags() []string {
	if m.movFlags != "" {
		return []string{"-movflags", m.movFlags}
	}

	return nil
}

// ObtainMap _
func (m *Mediafile) ObtainMap() []string {
	if m.mapFlag != "" {
		return []string{"-map", m.mapFlag}
	}

	return nil
}

// ObtainSegmentTime _
func (m *Mediafile) ObtainSegmentTime() []string {
	if m.segmentTime != 0 {
		return []string{"-segment_time", fmt.Sprintf("%d", m.segmentTime)}
	}

	return nil
}

// ObtainResetTimestamps _
func (m *Mediafile) ObtainResetTimestamps() []string {
	if m.resetTimestamps {
		return []string{"-reset_timestamps", "1"}
	}

	return nil
}

// ObtainHideBanner _
func (m *Mediafile) ObtainHideBanner() []string {
	if m.hideBanner {
		return []string{"-hide_banner"}
	}

	return nil
}

// ObtainFileSizeLimit _
func (m *Mediafile) ObtainFileSizeLimit() []string {
	if m.fileSizeLimit != "" {
		return []string{"-fs", m.fileSizeLimit}
	}

	return nil
}

// ObtainNativeFramerateInput _
func (m *Mediafile) ObtainNativeFramerateInput() []string {
	if m.nativeFramerateInput {
		return []string{"-re"}
	}

	return nil
}

// ObtainOutputPath _
func (m *Mediafile) ObtainOutputPath() []string {
	if m.outputPath != "" {
		return []string{m.outputPath}
	}

	return nil
}

// ObtainVideoCodec _
func (m *Mediafile) ObtainVideoCodec() []string {
	if m.videoCodec != "" {
		return []string{"-c:v", m.videoCodec}
	}

	return nil
}

// ObtainVframes _
func (m *Mediafile) ObtainVframes() []string {
	if m.vframes != 0 {
		return []string{"-vframes", fmt.Sprintf("%d", m.vframes)}
	}

	return nil
}

// ObtainFrameRate _
func (m *Mediafile) ObtainFrameRate() []string {
	if m.frameRate != 0 {
		return []string{"-r", fmt.Sprintf("%d", m.frameRate)}
	}

	return nil
}

// ObtainAudioRate _
func (m *Mediafile) ObtainAudioRate() []string {
	if m.audioRate != 0 {
		return []string{"-ar", fmt.Sprintf("%d", m.audioRate)}
	}

	return nil
}

// ObtainResolution _
func (m *Mediafile) ObtainResolution() []string {
	if m.resolution != "" {
		return []string{"-s", m.resolution}
	}

	return nil
}

// ObtainVideoBitRate _
func (m *Mediafile) ObtainVideoBitRate() []string {
	if m.videoBitRate != "" {
		return []string{"-b:v", m.videoBitRate}
	}

	return nil
}

// ObtainAudioCodec _
func (m *Mediafile) ObtainAudioCodec() []string {
	if m.audioCodec != "" {
		return []string{"-c:a", m.audioCodec}
	}

	return nil
}

// ObtainAudioBitRate _
func (m *Mediafile) ObtainAudioBitRate() []string {
	switch {
	case !m.audioVariableBitrate && m.audioBitrate != "":
		return []string{"-b:a", m.audioBitrate}
	case m.audioVariableBitrate && m.audioBitrate != "":
		return []string{"-q:a", m.audioBitrate}
	case m.audioVariableBitrate:
		return []string{"-q:a", "0"}
	default:
		return nil
	}
}

// ObtainAudioChannels _
func (m *Mediafile) ObtainAudioChannels() []string {
	if m.audioChannels != 0 {
		return []string{"-ac", fmt.Sprintf("%d", m.audioChannels)}
	}

	return nil
}

// ObtainVideoMaxBitRate _
func (m *Mediafile) ObtainVideoMaxBitRate() []string {
	if m.videoMaxBitRate != 0 {
		return []string{"-maxrate", fmt.Sprintf("%dk", m.videoMaxBitRate)}
	}

	return nil
}

// ObtainVideoMinBitRate _
func (m *Mediafile) ObtainVideoMinBitRate() []string {
	if m.videoMinBitrate != 0 {
		return []string{"-minrate", fmt.Sprintf("%dk", m.videoMinBitrate)}
	}

	return nil
}

// ObtainBufferSize _
func (m *Mediafile) ObtainBufferSize() []string {
	if m.bufferSize != 0 {
		return []string{"-bufsize", fmt.Sprintf("%dk", m.bufferSize)}
	}

	return nil
}

// ObtainVideoBitRateTolerance _
func (m *Mediafile) ObtainVideoBitRateTolerance() []string {
	if m.videoBitRateTolerance != 0 {
		return []string{"-bt", fmt.Sprintf("%dk", m.videoBitRateTolerance)}
	}

	return nil
}

// ObtainThreads _
func (m *Mediafile) ObtainThreads() []string {
	if m.threadset {
		return []string{"-threads", fmt.Sprintf("%d", m.threads)}
	}

	return nil
}

// ObtainTarget _
func (m *Mediafile) ObtainTarget() []string {
	if m.target != "" {
		return []string{"-target", m.target}
	}

	return nil
}

// ObtainDuration _
func (m *Mediafile) ObtainDuration() []string {
	if m.duration != "" {
		return []string{"-t", m.duration}
	}

	return nil
}

// ObtainDurationInput _
func (m *Mediafile) ObtainDurationInput() []string {
	if m.durationInput != "" {
		return []string{"-t", m.durationInput}
	}

	return nil
}

// ObtainKeyframeInterval _
func (m *Mediafile) ObtainKeyframeInterval() []string {
	if m.keyframeInterval != 0 {
		return []string{"-g", fmt.Sprintf("%d", m.keyframeInterval)}
	}

	return nil
}

// ObtainSeekTime _
func (m *Mediafile) ObtainSeekTime() []string {
	if m.seekTime != "" {
		return []string{"-ss", m.seekTime}
	}

	return nil
}

// ObtainSeekTimeInput _
func (m *Mediafile) ObtainSeekTimeInput() []string {
	if m.seekTimeInput != "" {
		return []string{"-ss", m.seekTimeInput}
	}

	return nil
}

// ObtainPreset _
func (m *Mediafile) ObtainPreset() []string {
	if m.preset != "" {
		return []string{"-preset", m.preset}
	}

	return nil
}

// ObtainTune _
func (m *Mediafile) ObtainTune() []string {
	if m.tune != "" {
		return []string{"-tune", m.tune}
	}

	return nil
}

// ObtainCRF _
func (m *Mediafile) ObtainCRF() []string {
	if m.crf != 0 {
		return []string{"-crf", fmt.Sprintf("%d", m.crf)}
	}

	return nil
}

// ObtainLibx265Params _
func (m *Mediafile) ObtainLibx265Params() []string {
	if m.libx265Params != nil {
		flags := make([]string, 0)

		if m.libx265Params.CRF > 0 {
			flags = append(flags, fmt.Sprintf("crf=%d", m.libx265Params.CRF))
		}

		if len(flags) > 0 {
			flags = append([]string{"-x265-params"}, flags...)
		}

		return flags
	}

	return nil
}

// ObtainQScale _
func (m *Mediafile) ObtainQScale() []string {
	if m.qscale != 0 {
		return []string{"-qscale", fmt.Sprintf("%d", m.qscale)}
	}

	return nil
}

// ObtainNvencRateControl _
func (m *Mediafile) ObtainNvencRateControl() []string {
	if m.nvencRateControl != "" {
		return []string{"-rc", m.nvencRateControl}
	}

	return nil
}

// ObtainMaxMuxingQueueSize _
func (m *Mediafile) ObtainMaxMuxingQueueSize() []string {
	if m.maxMuxingQueueSize != 0 {
		return []string{"-max_muxing_queue_size", fmt.Sprintf("%d", m.maxMuxingQueueSize)}
	}

	return nil
}

// ObtainConstantQuantization _
func (m *Mediafile) ObtainConstantQuantization() []string {
	if m.constantQuantization > 0 {
		return []string{"-qp", fmt.Sprintf("%d", m.constantQuantization)}
	}

	return nil
}

// ObtainNvencTargetQuality _
func (m *Mediafile) ObtainNvencTargetQuality() []string {
	if m.nvencTargetQuality > 0 {
		return []string{"-cq", fmt.Sprintf("%d", m.nvencTargetQuality)}
	}

	return nil
}

// ObtainStrict _
func (m *Mediafile) ObtainStrict() []string {
	if m.strict != 0 {
		return []string{"-strict", fmt.Sprintf("%d", m.strict)}
	}

	return nil
}

// ObtainVideoProfile _
func (m *Mediafile) ObtainVideoProfile() []string {
	if m.videoProfile != "" {
		return []string{"-profile:v", m.videoProfile}
	}

	return nil
}

// ObtainAudioProfile _
func (m *Mediafile) ObtainAudioProfile() []string {
	if m.audioProfile != "" {
		return []string{"-profile:a", m.audioProfile}
	}

	return nil
}

// ObtainCopyTs _
func (m *Mediafile) ObtainCopyTs() []string {
	if m.copyTs {
		return []string{"-copyts"}
	}

	return nil
}

// ObtainInputFormat _
func (m *Mediafile) ObtainInputFormat() []string {
	if m.inputFormat != "" {
		return []string{"-f", m.inputFormat}
	}

	return nil
}

// ObtainOutputFormat _
func (m *Mediafile) ObtainOutputFormat() []string {
	if m.outputFormat != "" {
		return []string{"-f", m.outputFormat}
	}

	return nil
}

// ObtainMuxDelay _
func (m *Mediafile) ObtainMuxDelay() []string {
	if m.muxDelay != "" {
		return []string{"-muxdelay", m.muxDelay}
	}

	return nil
}

// ObtainSeekUsingTsInput _
func (m *Mediafile) ObtainSeekUsingTsInput() []string {
	if m.seekUsingTsInput {
		return []string{"-seek_timestamp", "1"}
	}

	return nil
}

// ObtainRtmpLive _
func (m *Mediafile) ObtainRtmpLive() []string {
	if m.rtmpLive != "" {
		return []string{"-rtmp_live", m.rtmpLive}
	}

	return nil
}

// ObtainHlsPlaylistType _
func (m *Mediafile) ObtainHlsPlaylistType() []string {
	if m.hlsPlaylistType != "" {
		return []string{"-hls_playlist_type", m.hlsPlaylistType}
	}

	return nil
}

// ObtainInputInitialOffset _
func (m *Mediafile) ObtainInputInitialOffset() []string {
	if m.inputInitialOffset != "" {
		return []string{"-itsoffset", m.inputInitialOffset}
	}

	return nil
}

// ObtainHlsListSize _
func (m *Mediafile) ObtainHlsListSize() []string {
	if m.hlsListSize > 0 {
		return []string{"-hls_list_size", fmt.Sprintf("%d", m.hlsListSize)}
	}

	return nil
}

// ObtainHlsSegmentDuration _
func (m *Mediafile) ObtainHlsSegmentDuration() []string {
	if m.hlsSegmentDuration != 0 {
		return []string{"-hls_time", fmt.Sprintf("%d", m.hlsSegmentDuration)}
	}

	return nil
}

// ObtainHlsMasterPlaylistName _
func (m *Mediafile) ObtainHlsMasterPlaylistName() []string {
	if m.hlsMasterPlaylistName != "" {
		return []string{"-master_pl_name", fmt.Sprintf("%s", m.hlsMasterPlaylistName)}
	}

	return nil
}

// ObtainHlsSegmentFilename _
func (m *Mediafile) ObtainHlsSegmentFilename() []string {
	if m.hlsSegmentFilename != "" {
		return []string{"-hls_segment_filename", fmt.Sprintf("%s", m.hlsSegmentFilename)}
	}

	return nil
}

// ObtainHTTPMethod _
func (m *Mediafile) ObtainHTTPMethod() []string {
	if m.httpMethod != "" {
		return []string{"-method", m.httpMethod}
	}

	return nil
}

// ObtainPixFmt _
func (m *Mediafile) ObtainPixFmt() []string {
	if m.pixFmt != "" {
		return []string{"-pix_fmt", m.pixFmt}
	}

	return nil
}

// ObtainHTTPKeepAlive _
func (m *Mediafile) ObtainHTTPKeepAlive() []string {
	if m.httpKeepAlive {
		return []string{"-multiple_requests", "1"}
	}

	return nil
}

// ObtainSkipVideo _
func (m *Mediafile) ObtainSkipVideo() []string {
	if m.skipVideo {
		return []string{"-vn"}
	}

	return nil
}

// ObtainSkipAudio _
func (m *Mediafile) ObtainSkipAudio() []string {
	if m.skipAudio {
		return []string{"-an"}
	}

	return nil
}

// ObtainStreamIds _
func (m *Mediafile) ObtainStreamIds() []string {
	if m.streamIds != nil && len(m.streamIds) != 0 {
		result := []string{}
		for i, val := range m.streamIds {
			result = append(result, []string{"-streamid", fmt.Sprintf("%d:%s", i, val)}...)
		}
		return result
	}

	return nil
}

// ObtainCompressionLevel _
func (m *Mediafile) ObtainCompressionLevel() []string {
	if m.compressionLevel != 0 {
		return []string{"-compression_level", fmt.Sprintf("%d", m.compressionLevel)}
	}

	return nil
}

// ObtainMapMetadata _
func (m *Mediafile) ObtainMapMetadata() []string {
	if m.mapMetadata != "" {
		return []string{"-map_metadata", m.mapMetadata}
	}

	return nil
}

// ObtainEncryptionKey _
func (m *Mediafile) ObtainEncryptionKey() []string {
	if m.encryptionKey != "" {
		return []string{"-hls_key_info_file", m.encryptionKey}
	}

	return nil
}

// ObtainBframe _
func (m *Mediafile) ObtainBframe() []string {
	if m.bframe != 0 {
		return []string{"-bf", fmt.Sprintf("%d", m.bframe)}
	}

	return nil
}

// ObtainTags _
func (m *Mediafile) ObtainTags() []string {
	if m.tags != nil && len(m.tags) != 0 {
		result := []string{}
		for key, val := range m.tags {
			result = append(result, []string{"-metadata", fmt.Sprintf("%s=%s", key, val)}...)
		}
		return result
	}

	return nil
}

// ObtainVideoTag _
func (m *Mediafile) ObtainVideoTag() []string {
	if m.videoTag != "" {
		return []string{"-tag:v", m.videoTag}
	}

	return []string{}
}
