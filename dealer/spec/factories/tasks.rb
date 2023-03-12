FactoryBot.define do
  factory :convert_task, class: 'Tasks::Convert' do
    state { :published }

    convert_task_payload { build(:convert_task_payload) }

    trait :occupied do
      occupied_at { Time.current }
      after(:build) do |task|
        task.occupied_by = create(:performer) if task.occupied_by.blank?
      end
    end
  end
end
