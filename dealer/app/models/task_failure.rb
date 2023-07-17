# == Schema Information
#
# Table name: task_failures
#
#  id           :uuid             not null, primary key
#  reason       :string
#  created_at   :datetime         not null
#  updated_at   :datetime         not null
#  performer_id :uuid
#  task_id      :uuid             not null
#
# Indexes
#
#  index_task_failures_on_performer_id  (performer_id)
#  index_task_failures_on_task_id       (task_id)
#
# Foreign Keys
#
#  fk_rails_...  (performer_id => performers.id)
#  fk_rails_...  (task_id => tasks.id)
#
class TaskFailure < ApplicationRecord
  belongs_to :task
  belongs_to :performer, optional: true

  validates :reason, presence: true
end
