require 'rails_helper'

RSpec.describe MediaMetaDenormalizationService do
  let(:json_data) { JSON.parse(File.read('spec/examples/ffprobe.json')) }
  let!(:media_meta_report) { create(:media_meta_report, data: json_data) }

  describe '#perform' do
    describe 'denormalized data' do
      before { described_class.new(media_meta_report).perform }
      subject { media_meta_report }

      it 'writes media meta from json' do
        expect(subject.video_codec).to eq('hevc')
        expect(subject.video_bitrate).to eq(3_750_443)
        expect(subject.duration).to eq(565.823333)
        expect(subject.size).to eq(285_899_757)
        expect(subject.resolution_w).to eq(1920)
        expect(subject.resolution_h).to eq(1080)
        expect(subject.pix_fmt).to eq('yuv420p')
        expect(subject.created_at_by_meta).to eq(Time.zone.parse('2018-09-08T10:05:41.000000Z'))
      end
    end
  end
end
