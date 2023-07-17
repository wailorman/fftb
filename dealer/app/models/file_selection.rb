# == Schema Information
#
# Table name: file_selections
#
#  id                :uuid             not null, primary key
#  items_count       :integer
#  reached_max_depth :boolean          default(FALSE), not null
#  root_rclone_path  :string
#  created_at        :datetime         not null
#  updated_at        :datetime         not null
#
class FileSelection < ApplicationRecord
  has_many :items, dependent: :destroy, class_name: 'FileSelectionItem'
  has_many :orders, dependent: :restrict_with_error

  def video_items?
    items.reject(&:removed).any?(&:video?)
  end

  def audio_items?
    items.reject(&:removed).any?(&:audio?)
  end
end
