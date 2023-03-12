FactoryBot.define do
  factory :task_failure do
    # assosication :task
    # assosication :performer
    reason { 'failed!' }
  end
end
