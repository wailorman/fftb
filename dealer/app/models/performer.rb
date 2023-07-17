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
class Performer < ApplicationRecord
  has_many :tasks, dependent: :nullify

  def self.local
    @local ||= find_by!(name: :local)
  end
end
