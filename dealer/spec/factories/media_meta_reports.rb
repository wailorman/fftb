# == Schema Information
#
# Table name: media_meta_reports
#
#  id                  :uuid             not null, primary key
#  audio_bitrate       :bigint
#  audio_codec         :string
#  audio_codec_long    :string
#  bitrate             :bigint
#  created_at_by_meta  :datetime
#  created_at_by_mtime :datetime
#  created_at_by_name  :datetime
#  data                :jsonb
#  duration            :float
#  pix_fmt             :string
#  rclone_path         :string
#  resolution_h        :integer
#  resolution_w        :integer
#  size                :bigint
#  video_bitrate       :bigint
#  video_codec         :string
#  video_codec_long    :string
#  created_at          :datetime         not null
#  updated_at          :datetime         not null
#
# Indexes
#
#  index_media_meta_reports_on_rclone_path           (rclone_path)
#  index_media_meta_reports_on_rclone_path_and_size  (rclone_path,size)
#  index_media_meta_reports_on_size                  (size)
#
FactoryBot.define do
  factory :media_meta_report do
    rclone_path { 'storage_smb:/test/in/example.mp4' }
    data { {} }
  end
end
