package media

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wailorman/ffchunker/pkg/ctxlog"
	"github.com/wailorman/ffchunker/pkg/files"
	ffmpegModels "github.com/wailorman/goffmpeg/models"
	"github.com/wailorman/goffmpeg/transcoder"
)

// Converter _
type Converter struct {
	infoGetter InfoGetter
}

// NewConverter _
func NewConverter(infoGetter InfoGetter) *Converter {
	return &Converter{
		infoGetter: infoGetter,
	}
}

const (
	// NvencHardwareAccelerationType _
	NvencHardwareAccelerationType = "nvenc"
	// VideoToolboxHardwareAccelerationType _
	VideoToolboxHardwareAccelerationType = "videotoolbox"
)

const (
	// HevcCodecType _
	HevcCodecType = "hevc"
	// H264CodecType _
	H264CodecType = "h264"
)

// ConvertProgress _
type ConvertProgress struct {
	FramesProcessed string
	CurrentTime     string
	CurrentBitrate  string
	Progress        float64
	Speed           string
	FPS             int
	File            files.Filer
}

// ConverterTask _
type ConverterTask struct {
	InFile               files.Filer
	OutFile              files.Filer
	VideoCodec           string
	HardwareAcceleration string
	VideoBitRate         string
}

// RecursiveConverterTask _
type RecursiveConverterTask struct {
	InPath               files.Pather
	OutPath              files.Pather
	VideoCodec           string
	HardwareAcceleration string
	VideoBitRate         string
}

// RecursiveConvert _
func (c *Converter) RecursiveConvert(task RecursiveConverterTask) (
	progressChan chan ConvertProgress,
	doneChan chan bool,
	errChan chan error,
) {
	progressChan = make(chan ConvertProgress)
	doneChan = make(chan bool)
	errChan = make(chan error)

	log := ctxlog.New(ctxlog.DefaultContext + ".recursive-converter")

	go func() {
		defer close(progressChan)
		defer close(doneChan)
		defer close(errChan)

		var err error

		if err = task.OutPath.Create(); err != nil {
			errChan <- errors.Wrap(err, "Creating output path error")
			return
		}

		allFiles, err := task.InPath.Files()

		if err != nil {
			errChan <- errors.Wrap(err, "Listing files in input path error")
			return
		}

		videoFiles := make([]*files.File, 0)

		for _, file := range allFiles {
			fileLog := log.WithField("file_path", file.FullPath())

			mediaInfo, err := c.infoGetter.GetMediaInfo(file)

			if err != nil {
				fileLog.WithField("error", err).Debug("Getting media info error. Ignoring file")
				continue
			}

			if isVideo(mediaInfo) {
				videoFiles = append(videoFiles, file)

				fileLog.Debug("Found video file")
			} else {
				fileLog.Debug("File is not video")
			}
		}

		log.WithField("found_videos_count", len(videoFiles)).
			Debug("Video files filtering done")

		for fileIndex, file := range videoFiles {
			err := c.proceedRecursiveConvert(file, task, progressChan)

			if err != nil {
				errChan <- errors.Wrap(err, "Processing recursive conversion task error")
				return
			}

			log.WithFields(logrus.Fields{
				"file_index":  fileIndex + 1,
				"total_files": len(videoFiles),
				"file_path":   file.FullPath(),
			}).Info("File converted")
		}
	}()

	return progressChan, doneChan, errChan
}

func (c *Converter) proceedRecursiveConvert(file files.Filer, task RecursiveConverterTask, progressChan chan ConvertProgress) error {
	_progressChan, _doneChan, _errChan := c.Convert(ConverterTask{
		InFile:               file,
		OutFile:              file.NewWithSuffix("_out"),
		VideoCodec:           task.VideoCodec,
		HardwareAcceleration: task.HardwareAcceleration,
		VideoBitRate:         task.VideoBitRate,
	})

	for {
		select {
		case progress := <-_progressChan:
			progressChan <- progress
		case err := <-_errChan:
			return err
		case <-_doneChan:
			return nil
		}
	}
}

// Convert _
func (c *Converter) Convert(task ConverterTask) (
	progressChan chan ConvertProgress,
	doneChan chan bool,
	errChan chan error,
) {
	progressChan = make(chan ConvertProgress)
	doneChan = make(chan bool)
	errChan = make(chan error)

	log := ctxlog.New(ctxlog.DefaultContext + ".converter")

	go func() {

		defer close(progressChan)
		defer close(doneChan)
		defer close(errChan)

		// Create new instance of transcoder
		trans := new(transcoder.Transcoder)

		// Initialize transcoder passing the input file path and output file path
		err := trans.Initialize(
			task.InFile.FullPath(),
			task.OutFile.FullPath(),
		)

		if err != nil {
			errChan <- errors.Wrap(err, "Transcoder initializing error")
			return
		}

		metadata, err := c.infoGetter.GetMediaInfo(task.InFile)

		if err != nil {
			errChan <- errors.Wrap(err, "Getting file metadata")
			return
		}

		log.WithFields(
			logrus.Fields{
				"input_file":            task.InFile.FullPath(),
				"output_file":           task.InFile.FullPath(),
				"video_codec":           task.VideoCodec,
				"hardware_acceleration": task.HardwareAcceleration,
				"bit_rate":              task.VideoBitRate,
			},
		).Debug("Received converter task")

		if !isVideo(metadata) {
			errChan <- errors.Wrap(err, "Input file is not video")
			return
		}

		log.WithField("input_video_codec", getVideoCodec(metadata)).
			Debug("Detected input video codec")

		if task.VideoCodec != HevcCodecType && task.VideoCodec != H264CodecType {
			errChan <- fmt.Errorf("Unsupported video codec passed: `%s`. Only hevc & h264 is supported", task.VideoCodec)
			return
		}

		switch task.HardwareAcceleration {
		case NvencHardwareAccelerationType:
			nvencDecoder := c.nvencDecoder(metadata)

			if nvencDecoder != "" {
				trans.MediaFile().SetInputVideoCodec(nvencDecoder)
			} else {
				log.Debug("Input video codec is not supported by NVENC. Using CPU as decoder")
			}

			nvencEncoder := c.nvencEncoder(task.VideoCodec)

			if nvencEncoder != "" {
				trans.MediaFile().SetHardwareAcceleration("cuvid")
				trans.MediaFile().SetVideoCodec(nvencEncoder)
			} else {
				log.Warn("Encoding codec is not supported by NVENC. Using CPU as encoder")
			}

			log.WithFields(logrus.Fields{
				"decoder": nvencDecoder,
				"encoder": nvencEncoder,
			}).Debug("Using HW acceleration options")

		case VideoToolboxHardwareAccelerationType:
			vtbEncoder := c.videoToolboxEncoder(task.VideoCodec)

			if vtbEncoder != "" {
				trans.MediaFile().SetHardwareAcceleration("videotoolbox")
				trans.MediaFile().SetVideoCodec(vtbEncoder)
			} else {
				log.Warn("Encoding codec is not supported by VideoToolbox. Using CPU as encoder")
			}

			log.WithFields(logrus.Fields{
				"encoder": vtbEncoder,
			}).Debug("Using HW acceleration options")

		default:
			var codec string

			if task.VideoCodec == HevcCodecType {
				codec = "libx265"
			} else if task.VideoCodec == H264CodecType {
				codec = "libx264"
			}

			log.WithField("encoder", codec).
				Debug("Hardware acceleration not passed or unknown. Using CPU as encoder & decoder")

			trans.MediaFile().SetVideoCodec(codec)
		}

		if task.HardwareAcceleration != VideoToolboxHardwareAccelerationType {
			trans.MediaFile().SetPreset("slow")
		}

		trans.MediaFile().SetHideBanner(true)
		trans.MediaFile().SetVsync(true)
		trans.MediaFile().SetVideoBitRate(task.VideoBitRate)
		trans.MediaFile().SetAudioCodec("aac")

		if task.VideoCodec == HevcCodecType && task.HardwareAcceleration != VideoToolboxHardwareAccelerationType {
			trans.MediaFile().SetVideoTag("hvc1")
		}

		done := trans.Run(true)

		log.Debug("Converting started")

		progress := trans.Output()

		for {
			select {
			case progressMessage := <-progress:
				progressChan <- ConvertProgress{
					FramesProcessed: progressMessage.FramesProcessed,
					CurrentTime:     progressMessage.CurrentTime,
					CurrentBitrate:  progressMessage.CurrentBitrate,
					Progress:        progressMessage.Progress,
					Speed:           progressMessage.Speed,
					FPS:             progressMessage.FPS,
					File:            task.InFile,
				}
			case err := <-done:
				if err != nil {
					errChan <- err
					return
				}

				doneChan <- true
				return
			}
		}
	}()

	return progressChan, doneChan, errChan
}

func (c *Converter) nvencDecoder(metadata ffmpegModels.Metadata) string {
	if !isVideo(metadata) {
		return ""
	}

	if len(metadata.Streams) == 0 {
		return ""
	}

	switch metadata.Streams[0].CodecName {
	case "hevc":
		return "hevc_cuvid"
	case "h264":
		return "h264_cuvid"
	default:
		return ""
	}
}

func (c *Converter) videoToolboxDecoder(metadata ffmpegModels.Metadata) string {
	if !isVideo(metadata) {
		return ""
	}

	if len(metadata.Streams) == 0 {
		return ""
	}

	switch metadata.Streams[0].CodecName {
	case "hevc":
		return "hevc_videotoolbox"
	case "h264":
		return "h264_videotoolbox"
	default:
		return ""
	}
}

func (c *Converter) nvencEncoder(codec string) string {
	switch codec {
	case HevcCodecType:
		return "hevc_nvenc"
	case H264CodecType:
		return "h264_nvenc"
	default:
		return ""
	}
}

func (c *Converter) videoToolboxEncoder(codec string) string {
	switch codec {
	case HevcCodecType:
		return "h264_videotoolbox"
	case H264CodecType:
		return "hevc_videotoolbox"
	default:
		return ""
	}
}

func isVideo(metadata ffmpegModels.Metadata) bool {
	if len(metadata.Streams) == 0 {
		return false
	}

	if metadata.Streams[0].BitRate == "" {
		return false
	}

	return true
}

func getVideoCodec(metadata ffmpegModels.Metadata) string {
	if !isVideo(metadata) {
		return ""
	}

	return metadata.Streams[0].CodecName
}
