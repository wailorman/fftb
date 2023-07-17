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
class Orders::Convert < Order
  validate :ensure_payload_has_correct_type

  def set_default_values
    super()

    unless payload
      self.payload = Payloads::Orders::Convert.new
      payload.set_default_values
    end
  end

  def build_payload(params = {})
    self.payload = Payloads::Orders::Convert.new(params.except(:type))
  end

  private def ensure_payload_has_correct_type
    errors.add(:base, "Incorrect payload type: `#{payload.class.name}`") if payload.class != Payloads::Orders::Convert
  end
end
