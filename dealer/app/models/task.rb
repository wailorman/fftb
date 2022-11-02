class Task < ApplicationRecord
  OCCUPATION_TTL = 2.minutes

  enum state: { published: 'published',
                finished: 'finished',
                failed: 'failed' }

  enum kind: { convert_v1: Fftb::SegmentType::CONVERT_V1 }

  belongs_to :occupied_by, class_name: 'Performer', optional: true
  belongs_to :input_storage_claim, class_name: 'StorageClaim', optional: true
  belongs_to :output_storage_claim, class_name: 'StorageClaim', optional: true

  validates :current_progress, comparison: { greater_than_or_equal_to: 0, less_than_or_equal_to: 1 }
  validates :current_step, inclusion: { in: %w[uploading_input downloading_input processing
                                               uploading_output downloading_output] }, if: :current_step

  scope :not_occupied, -> { where('occupied_at < ? OR occupied_at IS NULL', Time.current - OCCUPATION_TTL) }
  scope :not_occupied_for, -> (performer) {
    where(
      'occupied_at < :current_timestamp OR (occupied_at > :current_timestamp AND occupied_by_id = :performer_id) OR occupied_at IS NULL',
      current_timestamp: Time.current - OCCUPATION_TTL,
      performer_id: performer.id
    )
  }

  def deoccupy
    self.occupied_at = nil
    self.occupied_by = nil
  end
end
