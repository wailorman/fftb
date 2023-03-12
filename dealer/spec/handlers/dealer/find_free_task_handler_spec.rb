require 'rails_helper'

RSpec.describe Dealer::FindFreeTaskHandler, type: :handler do
  let(:performer) { Performer.local }
  let!(:task) { create(:task, :convert) }

  let(:request) { Fftb::FindFreeTaskRequest.new(authorization: performer.name) }

  describe '#call' do
    subject(:response) { described_class.new(request, nil).call }

    context 'when free task exists' do
      it { expect(response).to be_kind_of(Fftb::Task) }
      it { expect(response.id).to eq(task.id) }
    end

    context 'when no free tasks' do
      let(:task) { nil }

      it { expect(response).to be_kind_of(Twirp::Error) }
      it { expect(response.code).to eq(:not_found) }
    end
  end
end
