# == Schema Information
#
# Table name: convert_task_payloads
#
#  id                   :uuid             not null, primary key
#  current_bitrate      :float
#  current_fps          :float
#  current_frame        :bigint
#  current_speed        :float
#  current_time         :bigint
#  input_rclone_path    :string
#  opts                 :string           not null, is an Array
#  output_rclone_path   :string
#  created_at           :datetime         not null
#  updated_at           :datetime         not null
#  media_meta_report_id :uuid
#
# Indexes
#
#  index_convert_task_payloads_on_media_meta_report_id  (media_meta_report_id)
#
# Foreign Keys
#
#  fk_rails_...  (media_meta_report_id => media_meta_reports.id)
#
class Payloads::Tasks::Convert < ApplicationRecord
  self.table_name = :convert_task_payloads

  belongs_to :media_meta_report, optional: true

  validates :opts,
            :input_rclone_path,
            :output_rclone_path, presence: true

  def string_opts=(val)
    self.opts = OptsHelper.string_to_array(val)
  end

  def string_opts
    OptsHelper.array_to_string(opts)
  end
end
