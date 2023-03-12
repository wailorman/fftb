module Tasks
  class Convert < Task
    has_one :convert_task_payload, foreign_key: :task_id
  end
end
