class ConvertTaskPayload < ApplicationRecord
  belongs_to :convert_task, class_name: 'Tasks::Convert', foreign_key: :task_id

  validates :opts,
            :input_rclone_path,
            :output_rclone_path, presence: true
end
