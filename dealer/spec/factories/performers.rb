# == Schema Information
#
# Table name: performers
#
#  id         :uuid             not null, primary key
#  name       :string
#  token      :string
#  created_at :datetime         not null
#  updated_at :datetime         not null
#
FactoryBot.define do
  factory :performer do
  end
end
