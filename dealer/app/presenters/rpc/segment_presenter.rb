module Rpc
  class SegmentPresenter < ::Rpc::BasePresenter
    alias task object

    def call
      Fftb::Segment.new(
        type: segment_type,
        id: task.id,
        convertParams: convert_params
      )
    end

    private

    def segment_type
      Fftb::SegmentType::CONVERT_V1
    end

    def convert_params
      Fftb::ConvertSegmentParams.new(
        videoCodec: task.params['video_codec'],
        hwAccel: task.params['hw_accel'],
        videoBitRate: task.params['video_bit_rate'],
        videoQuality: task.params['video_quality'],
        preset: task.params['preset'],
        scale: task.params['scale'],
        keyframeInterval: task.params['keyframe_interval'],
        muxer: task.params['muxer'],
        position: task.params['position']
      )
    end
  end
end
