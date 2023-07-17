# == Schema Information
#
# Table name: convert_order_payloads
#
#  id                 :uuid             not null, primary key
#  audio_muxer        :string
#  audio_opts         :string
#  output_rclone_path :string
#  video_muxer        :string
#  video_opts         :string
#  created_at         :datetime         not null
#  updated_at         :datetime         not null
#
class Payloads::Orders::Convert < ApplicationRecord
  self.table_name = :convert_order_payloads

  validates :output_rclone_path, presence: true

  def set_default_values
    defaults = config.dig(:order_defaults, :convert)

    self.id ||= SecureRandom.uuid
    self.video_muxer ||= defaults[:video_muxer]
    self.video_opts ||= defaults[:video_opts]
    self.audio_muxer ||= defaults[:audio_muxer]
    self.audio_opts ||= defaults[:audio_opts]

    self.output_rclone_path ||= Rclone.join_path(
      defaults[:output_location],
      Time.zone.today.to_s,
      id,
      '/'
    )
  end
end
