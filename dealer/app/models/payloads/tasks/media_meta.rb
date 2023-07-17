# == Schema Information
#
# Table name: media_meta_task_payloads
#
#  id                   :uuid             not null, primary key
#  input_rclone_path    :string
#  output_rclone_path   :string
#  created_at           :datetime         not null
#  updated_at           :datetime         not null
#  media_meta_report_id :uuid
#
# Indexes
#
#  index_media_meta_task_payloads_on_media_meta_report_id  (media_meta_report_id) UNIQUE
#
# Foreign Keys
#
#  fk_rails_...  (media_meta_report_id => media_meta_reports.id)
#
class Payloads::Tasks::MediaMeta < ApplicationRecord
  self.table_name = :media_meta_task_payloads

  # belongs_to :media_meta_task, class_name: 'Tasks::MediaMeta', foreign_key: :task_id
  belongs_to :media_meta_report, optional: true

  validates :input_rclone_path,
            :output_rclone_path, presence: true
end
