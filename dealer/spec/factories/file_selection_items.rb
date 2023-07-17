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
FactoryBot.define do
  factory :file_selection_item do
    file_selection { nil }
    rclone_path { "MyString" }
  end
end
