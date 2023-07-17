class MediaMetaDenormalizationService
  attr_reader :media_meta_report

  def initialize(media_meta_report)
    @media_meta_report = media_meta_report
  end

  def perform
    return if json_data.blank?

    media_meta_report.video_codec = video_stream[:codec_name]
    media_meta_report.video_codec_long = video_stream[:codec_long_name]
    media_meta_report.video_bitrate = video_stream[:bit_rate].to_i
    media_meta_report.audio_codec = audio_stream[:codec_name]
    media_meta_report.audio_codec_long = audio_stream[:codec_long_name]
    media_meta_report.audio_bitrate = audio_stream[:bit_rate].to_i
    media_meta_report.bitrate = format_meta[:bit_rate].to_i
    media_meta_report.duration = format_meta[:duration].to_f
    media_meta_report.size = format_meta[:size].to_i
    media_meta_report.resolution_w = video_stream[:width]
    media_meta_report.resolution_h = video_stream[:height]
    media_meta_report.pix_fmt = video_stream[:pix_fmt]

    if format_meta.dig(:tags, :creation_time)
      media_meta_report.created_at_by_meta = Time.zone.parse(format_meta.dig(:tags, :creation_time))
    end
  end

  private def video_stream
    @video_stream ||=
      json_data[:streams].detect { |stream| stream[:codec_type] == 'video' } || {}
  end

  private def audio_stream
    @audio_stream ||=
      json_data[:streams].detect { |stream| stream[:codec_type] == 'audio' } || {}
  end

  private def format_meta
    json_data[:format] || {}
  end

  private def json_data
    @json_data ||= media_meta_report.data.deep_symbolize_keys
  end
end
