require 'rails_helper'

RSpec.describe Dealer::FindFreeSegmentHandler, type: :handler do
  let(:performer) { Performer.local }
  let!(:task) { create(:task, :convert) }

  let(:request) { Fftb::FindFreeSegmentRequest.new(authorization: performer.name) }

  describe '#call' do
    subject(:response) { described_class.new(request, nil).call }

    context 'when free segment exists' do
      it { expect(response).to be_kind_of(Fftb::Segment) }
      it { expect(response.id).to eq(task.id) }
      it { expect(response.convertParams.videoCodec).to eq(task.convert_params.video_codec) }
    end

    context 'when no free segments' do
      let(:task) { nil }

      it { expect(response).to be_kind_of(Twirp::Error) }
      it { expect(response.code).to eq(:not_found) }
    end
  end
end
