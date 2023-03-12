class TaskFailure < ApplicationRecord
  belongs_to :task
  belongs_to :performer

  validates :reason, presence: true
end
