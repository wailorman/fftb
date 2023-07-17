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
class Tasks::Convert < Task
  # has_one :payload, foreign_key: :task_id, dependent: :destroy, class_name: 'Payloads::Task::Convert'

  validates :payload_type, inclusion: { in: ['Payloads::Tasks::Convert'] }

  def verify_result
    update!(result_verified: true)
  end
end
