FactoryBot.define do
  factory :task do
    state { :published }

    trait :convert do
      kind { :convert_v1 }
      association :convert_params, factory: :convert_params
    end

    trait :occupied do
      occupied_at { Time.current }
      after(:build) do |task|
        task.occupied_by = create(:performer) if task.occupied_by.blank?
      end
    end
  end
end
