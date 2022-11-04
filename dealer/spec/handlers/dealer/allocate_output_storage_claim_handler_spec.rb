require 'rails_helper'

RSpec.describe Dealer::AllocateOutputStorageClaimHandler, type: :handler do
  let(:performer) { Performer.local }
  let!(:task) { create(:task, :convert, occupied_by: performer) }

  let(:request) do
    Fftb::StorageClaimRequest.new(authorization: performer.name,
                                  segmentId: task.id,
                                  purpose: Fftb::StorageClaimPurpose::CONVERT_OUTPUT,
                                  name: 'output.mp4')
  end

  describe '#call' do
    subject(:response) { described_class.new(request, nil).call }

    it { expect(response).to be_kind_of(Fftb::StorageClaim) }
    it { expect(response.id).to eq(OutputStorageClaim.last.id) }
    it do
      expect { response }
        .to change(OutputStorageClaim, :count)
        .by(1)
    end

    describe 'created storage claim' do
      subject(:created_storage_claim) { OutputStorageClaim.find(response.id) }

      it { expect(created_storage_claim.purpose).to eq('convert_output') }
      it { expect(created_storage_claim.name).to eq('output.mp4') }
    end
  end
end
