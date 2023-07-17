# == Schema Information
#
# Table name: orders
#
#  id                :uuid             not null, primary key
#  payload_type      :string           not null
#  state             :string           default("created")
#  type              :string
#  created_at        :datetime         not null
#  updated_at        :datetime         not null
#  file_selection_id :uuid
#  payload_id        :uuid             not null
#
# Indexes
#
#  index_orders_on_file_selection_id  (file_selection_id)
#
# Foreign Keys
#
#  fk_rails_...  (file_selection_id => file_selections.id)
#
class Order < ApplicationRecord
  extend Enumerize

  belongs_to :payload, polymorphic: true, dependent: :destroy
  belongs_to :file_selection, optional: true

  has_many :tasks, inverse_of: :order, dependent: :destroy

  enumerize :state, in: %i[created published cancelled], scope: true

  accepts_nested_attributes_for :payload,
                                :tasks

  def short_type
    ActiveSupport::StringInquirer.new(type.split('::').last.underscore)
  end

  def set_default_values
  end

  def build_payload(*args)
    raise NotImplementedError
  end

  def publish
    self.state = :published
    tasks.each { |t| t.state = :published }
  end

  def cancel
    self.state = :cancelled
    tasks.each { |t| t.state = :cancelled }
  end
end
