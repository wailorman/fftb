class Performer < ApplicationRecord
  has_many :tasks, dependent: :nullify

  def self.local
    @local ||= find_by!(name: :local)
  end
end
