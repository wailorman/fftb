package models

import (
	"strconv"
)

// Ffmpeg _
type Ffmpeg struct {
	FfmpegBinPath  string
	FfprobeBinPath string
}

// Metadata _
type Metadata struct {
	Streams []Streams `json:"streams"`
	Format  Format    `json:"format"`
}

// Streams _
type Streams struct {
	Index              int         `json:"index"`
	ID                 string      `json:"id"`
	CodecName          string      `json:"codec_name"`
	CodecLongName      string      `json:"codec_long_name"`
	Profile            string      `json:"profile"`
	CodecType          string      `json:"codec_type"`
	CodecTimeBase      string      `json:"codec_time_base"`
	CodecTagString     string      `json:"codec_tag_string"`
	CodecTag           string      `json:"codec_tag"`
	Width              int         `json:"width"`
	Height             int         `json:"height"`
	CodedWidth         int         `json:"coded_width"`
	CodedHeight        int         `json:"coded_height"`
	HasBFrames         int         `json:"has_b_frames"`
	SampleAspectRatio  string      `json:"sample_aspect_ratio"`
	DisplayAspectRatio string      `json:"display_aspect_ratio"`
	PixFmt             string      `json:"pix_fmt"`
	Level              int         `json:"level"`
	ChromaLocation     string      `json:"chroma_location"`
	Refs               int         `json:"refs"`
	QuarterSample      string      `json:"quarter_sample"`
	DivxPacked         string      `json:"divx_packed"`
	RFrameRrate        string      `json:"r_frame_rate"`
	AvgFrameRate       string      `json:"avg_frame_rate"`
	TimeBase           string      `json:"time_base"`
	DurationTs         int         `json:"duration_ts"`
	Duration           string      `json:"duration"`
	Disposition        Disposition `json:"disposition"`
	BitRate            string      `json:"bit_rate"`

	DurationFloat float64
}

// Framer _
type Framer interface {
	GetType() string
}

// AudioFrame _
type AudioFrame struct {
	MediaType               string `json:"media_type"`
	StreamIndex             int    `json:"stream_index"`
	KeyFrame                int    `json:"key_frame"`
	PktPts                  int    `json:"pkt_pts"`
	PktPtsTime              string `json:"pkt_pts_time"`
	PktDts                  int    `json:"pkt_dts"`
	PktDtsTime              string `json:"pkt_dts_time"`
	BestEffortTimestamp     int    `json:"best_effort_timestamp"`
	BestEffortTimestampTime string `json:"best_effort_timestamp_time"`
	PktDuration             int    `json:"pkt_duration"`
	PktDurationTime         string `json:"pkt_duration_time"`
	PktPos                  string `json:"pkt_pos"`
	PktSize                 string `json:"pkt_size"`
	SampleFmt               string `json:"sample_fmt"`
	NbSamples               int    `json:"nb_samples"`
	Channels                int    `json:"channels"`
	ChannelLayout           string `json:"channel_layout"`
}

// GetType _
func (f *AudioFrame) GetType() string {
	return "AudioFrame"
}

// PktPtsTimeFloat _
func (f *AudioFrame) PktPtsTimeFloat() (float64, error) {
	return strconv.ParseFloat(f.PktPtsTime, 64)
}

// PktDtsTimeFloat _
func (f *AudioFrame) PktDtsTimeFloat() (float64, error) {
	return strconv.ParseFloat(f.PktDtsTime, 64)
}

// BestEffortTimestampTimeFloat _
func (f *AudioFrame) BestEffortTimestampTimeFloat() (float64, error) {
	return strconv.ParseFloat(f.BestEffortTimestampTime, 64)
}

// PktDurationTimeFloat _
func (f *AudioFrame) PktDurationTimeFloat() (float64, error) {
	return strconv.ParseFloat(f.PktDurationTime, 64)
}

// PktPosInt _
func (f *AudioFrame) PktPosInt() (int64, error) {
	return strconv.ParseInt(f.PktPos, 10, 64)
}

// PktSizeInt _
func (f *AudioFrame) PktSizeInt() (int64, error) {
	return strconv.ParseInt(f.PktSize, 10, 64)
}

// VideoFrame _
type VideoFrame struct {
	MediaType               string `json:"media_type"`
	StreamIndex             int    `json:"stream_index"`
	KeyFrame                int    `json:"key_frame"`
	PktPts                  int    `json:"pkt_pts"`
	PktPtsTime              string `json:"pkt_pts_time"`
	PktDts                  int    `json:"pkt_dts"`
	PktDtsTime              string `json:"pkt_dts_time"`
	BestEffortTimestamp     int    `json:"best_effort_timestamp"`
	BestEffortTimestampTime string `json:"best_effort_timestamp_time"`
	PktDuration             int    `json:"pkt_duration"`
	PktDurationTime         string `json:"pkt_duration_time"`
	PktPos                  string `json:"pkt_pos"`
	PktSize                 string `json:"pkt_size"`
	Width                   int    `json:"width"`
	Height                  int    `json:"height"`
	PixFmt                  string `json:"pix_fmt"`
	SampleAspectRatio       string `json:"sample_aspect_ratio"`
	PictType                string `json:"pict_type"`
	CodedPictureNumber      int    `json:"coded_picture_number"`
	DisplayPictureNumber    int    `json:"display_picture_number"`
	InterlacedFrame         int    `json:"interlaced_frame"`
	TopFieldFirst           int    `json:"top_field_first"`
	RepeatPict              int    `json:"repeat_pict"`
	ChromaLocation          string `json:"chroma_location"`
}

// GetType _
func (f *VideoFrame) GetType() string {
	return "VideoFrame"
}

// PktPtsTimeFloat _
func (f *VideoFrame) PktPtsTimeFloat() (float64, error) {
	return strconv.ParseFloat(f.PktPtsTime, 64)
}

// PktDtsTimeFloat _
func (f *VideoFrame) PktDtsTimeFloat() (float64, error) {
	return strconv.ParseFloat(f.PktDtsTime, 64)
}

// BestEffortTimestampTimeFloat _
func (f *VideoFrame) BestEffortTimestampTimeFloat() (float64, error) {
	return strconv.ParseFloat(f.BestEffortTimestampTime, 64)
}

// PktDurationTimeFloat _
func (f *VideoFrame) PktDurationTimeFloat() (float64, error) {
	return strconv.ParseFloat(f.PktDurationTime, 64)
}

// PktPosInt _
func (f *VideoFrame) PktPosInt() (int64, error) {
	return strconv.ParseInt(f.PktPos, 10, 64)
}

// PktSizeInt _
func (f *VideoFrame) PktSizeInt() (int64, error) {
	return strconv.ParseInt(f.PktSize, 10, 64)
}

// Disposition _
type Disposition struct {
	Default         int `json:"default"`
	Dub             int `json:"dub"`
	Original        int `json:"original"`
	Comment         int `json:"comment"`
	Lyrics          int `json:"lyrics"`
	Karaoke         int `json:"karaoke"`
	Forced          int `json:"forced"`
	HearingImpaired int `json:"hearing_impaired"`
	VisualImpaired  int `json:"visual_impaired"`
	CleanEffects    int `json:"clean_effects"`
}

// Format _
type Format struct {
	Filename       string `json:"filename"`
	NbStreams      int    `json:"nb_streams"`
	NbPrograms     int    `json:"nb_programs"`
	FormatName     string `json:"format_name"`
	FormatLongName string `json:"format_long_name"`
	Duration       string `json:"duration"`
	Size           string `json:"size"`
	BitRate        string `json:"bit_rate"`
	ProbeScore     int    `json:"probe_score"`
	Tags           Tags   `json:"tags"`
}

// Progress _
type Progress struct {
	FramesProcessed string
	CurrentTime     string
	CurrentBitrate  string
	Progress        float64
	Speed           string
	FPS             float64
}

// Tags _
type Tags struct {
	Encoder string `json:"ENCODER"`
}
