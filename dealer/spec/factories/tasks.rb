FactoryBot.define do
  factory :task do
    trait :occupied do
      occupied_at { Time.current }
      after(:build) do |task|
        task.occupied_by = create(:performer) if task.occupied_by.blank?
      end
    end
  end
end
