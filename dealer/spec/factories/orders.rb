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
FactoryBot.define do
  factory :order do
    type { 'Orders::Convert' }
    payload { build(:convert_order_payload) }
  end

  factory :convert_order, class: 'Orders::Convert' do
    type { 'Orders::Convert' }

    after(:build) do |order| # rubocop:disable Style/SymbolProc
      order.set_default_values
    end
  end

  factory :convert_order_payload, class: 'Payloads::Orders::Convert' do
    after(:build) do |order| # rubocop:disable Style/SymbolProc
      order.set_default_values
    end
  end
end
