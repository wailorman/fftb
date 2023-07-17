# == Schema Information
#
# Table name: file_selection_items
#
#  id                :uuid             not null, primary key
#  mime_type         :string
#  rclone_path       :string
#  removed           :boolean          default(FALSE), not null
#  size              :bigint
#  created_at        :datetime         not null
#  updated_at        :datetime         not null
#  file_selection_id :uuid             not null
#
# Indexes
#
#  index_file_selection_items_on_file_selection_id  (file_selection_id)
#
# Foreign Keys
#
#  fk_rails_...  (file_selection_id => file_selections.id)
#
class FileSelectionItem < ApplicationRecord
  belongs_to :file_selection, counter_cache: :items_count
  has_many :tasks, inverse_of: :file_selection_item, dependent: :restrict_with_error

  default_scope { not_removed }
  scope :not_removed, -> { where(removed: false) }

  def video?
    mime_type.match?(/^video\//)
  end

  def audio?
    mime_type.match?(/^audio\//)
  end
end
