require 'rails_helper'

RSpec.describe Dealer::GetAllInputStorageClaimsHandler, type: :handler do
  let(:performer) { Performer.local }
  let!(:task) { create(:task, :convert, occupied_by: performer, input_storage_claims: [build(:input_storage_claim)]) }

  let(:request) do
    Fftb::StorageClaimRequest.new(authorization: performer.name,
                                  segmentId: task.id,
                                  purpose: Fftb::StorageClaimPurpose::CONVERT_INPUT)
  end

  describe '#call' do
    subject(:response) { described_class.new(request, nil).call }

    it { expect(response).to be_kind_of(Fftb::StorageClaimList) }
    it { expect(response.storageClaims.map(&:id)).to contain_exactly(*task.input_storage_claims.pluck(:id)) }
  end
end
