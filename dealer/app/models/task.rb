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
class Task < ApplicationRecord
  extend Enumerize

  OCCUPATION_TTL = 2.minutes

  enumerize :state, in: %i[created published cancelled finished failed], scope: true

  with_options inverse_of: :task do
    has_many :task_failures, dependent: :destroy
  end

  belongs_to :payload, polymorphic: true, dependent: :destroy
  belongs_to :order, optional: true
  belongs_to :file_selection_item, optional: true
  belongs_to :occupied_by, class_name: 'Performer', optional: true

  accepts_nested_attributes_for :payload

  after_commit :handle_state_change, if: -> { saved_change_to_state? }

  validates :current_progress, comparison: { greater_than_or_equal_to: 0, less_than_or_equal_to: 1 }
  validates :current_step, inclusion: { in: %w[downloading_input processing uploading_output] }, if: :current_step

  scope :not_occupied, -> { where('occupied_at < ? OR occupied_at IS NULL', Time.current - OCCUPATION_TTL) }
  scope :not_occupied_for, ->(performer) {
    where(
      'occupied_at < :current_timestamp OR (occupied_at > :current_timestamp AND occupied_by_id = :performer_id) OR occupied_at IS NULL',
      current_timestamp: Time.current - OCCUPATION_TTL,
      performer_id: performer.id
    )
  }
  scope :not_failed_by, ->(performer) {
    where <<~SQL.squish
      NOT EXISTS (
        SELECT 1 FROM task_failures
        WHERE task_failures.task_id = tasks.id AND task_failures.performer_id = '#{performer.id}'
      )
    SQL
  }

  def short_type
    ActiveSupport::StringInquirer.new(type.split('::').last.underscore)
  end

  def deoccupy
    self.occupied_at = nil
    self.occupied_by = nil
  end

  def occupied?
    return false if occupied_by.blank? || occupied_at.blank?

    occupied_at > OCCUPATION_TTL.ago
  end

  def handle_state_change
    verify_result if state == :finished
  end

  def verify_result
    raise NotImplementedError
  end

  def mark_failed(reason)
    self.state = :failed
    task_failures.build(performer: nil, reason: reason)
    save!
  end

  def safe_occupied_by
    return nil unless occupied?

    occupied_by
  end

  def build_payload(params)
    self.payload = params[:class].constantize.new(params.except(:type)) if params[:class]
  end
end
