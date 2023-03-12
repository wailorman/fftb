class Task < ApplicationRecord
  extend Enumerize

  OCCUPATION_TTL = 2.minutes

  enumerize :state, in: %i[published finished failed], scope: true

  has_many :task_failures
  belongs_to :occupied_by, class_name: 'Performer', optional: true

  validates :current_progress, comparison: { greater_than_or_equal_to: 0, less_than_or_equal_to: 1 }
  validates :current_step, inclusion: { in: %w[downloading_input processing uploading_output] }, if: :current_step

  scope :not_occupied, -> { where('occupied_at < ? OR occupied_at IS NULL', Time.current - OCCUPATION_TTL) }
  scope :not_occupied_for, -> (performer) {
    where(
      'occupied_at < :current_timestamp OR (occupied_at > :current_timestamp AND occupied_by_id = :performer_id) OR occupied_at IS NULL',
      current_timestamp: Time.current - OCCUPATION_TTL,
      performer_id: performer.id
    )
  }
  scope :not_failed_by, -> (performer) {
    where <<~SQL.squish
      NOT EXISTS (
        SELECT 1 FROM task_failures
        WHERE task_failures.task_id = tasks.id AND task_failures.performer_id = '#{performer.id}'
      )
    SQL
  }

  def deoccupy
    self.occupied_at = nil
    self.occupied_by = nil
  end
end
