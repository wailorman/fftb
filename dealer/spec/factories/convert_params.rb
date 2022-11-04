FactoryBot.define do
  factory :convert_params do
    video_codec { 'h264' }
    hw_accel { nil }
    video_bit_rate { nil }
    video_quality { 28 }
    preset { 'fast' }
    scale { nil }
    keyframe_interval { nil }
    muxer { 'mp4' }
  end
end
