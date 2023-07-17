# == Schema Information
#
# Table name: tasks
#
#  id                     :uuid             not null, primary key
#  current_progress       :float            default(0.0), not null
#  current_step           :string
#  occupied_at            :datetime
#  payload_type           :string           not null
#  result_verified        :boolean          default(FALSE)
#  state                  :string           default("published")
#  type                   :string           not null
#  created_at             :datetime         not null
#  updated_at             :datetime         not null
#  file_selection_item_id :uuid
#  occupied_by_id         :uuid
#  order_id               :uuid
#  payload_id             :uuid             not null
#
# Indexes
#
#  index_tasks_on_file_selection_item_id  (file_selection_item_id)
#  index_tasks_on_occupied_by_id          (occupied_by_id)
#  index_tasks_on_order_id                (order_id)
#
# Foreign Keys
#
#  fk_rails_...  (file_selection_item_id => file_selection_items.id)
#  fk_rails_...  (occupied_by_id => performers.id)
#  fk_rails_...  (order_id => orders.id)
#

class Tasks::MediaMeta < Task
  # belongs_to :payload, polymorphic: true, dependent: :destroy, class_name: 'Payloads::Tasks::MediaMeta'

  # validates :payload, presence: true
  # accepts_nested_attributes_for :payload

  validates :payload_type, inclusion: { in: ['Payloads::Tasks::MediaMeta'] }

  def verify_result
    found =
      Rclone.ls(payload.output_rclone_path)
            .reject { |entry| entry[:name].start_with?('.') }
            .sort_by { |entry| entry[:mod_time] }
            .reverse
            .detect { |entry| entry[:mime_type] == 'application/json' }

    unless found
      mark_failed('Result json not found in output path')
      return
    end

    path = File.join(payload.output_rclone_path, found[:name])
    content = Rclone.read(path)
    meta = JSON.parse(content)

    payload
      .update!(media_meta_report: ::MediaMetaReport
                                    .create!(rclone_path: payload.input_rclone_path,
                                             data: meta))

    update!(result_verified: true)

    save!
  rescue => e # rubocop:disable Style/RescueStandardError
    mark_failed("Failed to pull result: #{e.message}")
    # TODO: sentry
  end
end
