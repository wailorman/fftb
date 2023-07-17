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
FactoryBot.define do
  factory :convert_task, class: 'Tasks::Convert' do
    state { :published }

    payload { build(:convert_task_payload) }

    trait :occupied do
      occupied_at { Time.current }
      after(:build) do |task|
        task.occupied_by = create(:performer) if task.occupied_by.blank?
      end
    end
  end

  factory :convert_task_payload, class: 'Payloads::Tasks::Convert' do
    opts do
      %w[
        -i input/example.mp4
        -c:v h264
        -b:v 5M
        -c:a aac
        output/example.mp4
      ]
    end
    input_rclone_path { 'storage_smb:/test/in/example.mp4' }
    output_rclone_path { 'storage_smb:/test/out/example.mp4' }
  end

  factory :media_meta_task_payload, class: 'Payloads::Tasks::MediaMeta' do
    input_rclone_path { 'storage_smb:/test/in/example.mp4' }
    output_rclone_path { 'storage_smb:/test/out/example.mp4' }
  end
end
