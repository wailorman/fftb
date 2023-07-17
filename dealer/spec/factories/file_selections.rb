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
FactoryBot.define do
  factory :file_selection do
    
  end
end
