require 'rails_helper'

RSpec.describe BuildConvertOrderService do
  let(:files) do
    %w[
      storage:/records/archive/2023-07-09/CAM1/video.avi
      storage:/records/archive/2023-07-09/video.avi
      storage:/records/archive/2023-07-09/MIC1/audio.ogg
      storage:/records/archive/2023-07-10/CAM1/video.avi
    ]
  end
  let!(:file_selection) do
    create(
      :file_selection,
      items: files.map do |file_path|
        build(
          :file_selection_item,
          rclone_path: file_path,
          mime_type: file_path.match?(/audio/) ? 'audio/ogg' : 'video/avi'
        )
      end
    )
  end
  let(:file_selection_item) { file_selection.items.last }
  let!(:order) do
    create(
      :convert_order,
      file_selection: file_selection,
      payload: build(
        :convert_order_payload,
        output_rclone_path: 'storage:/records/output/'
      )
    )
  end

  subject(:instance) { described_class.new(order) }

  describe '#build_task' do
    subject(:task) { instance.send(:build_task, file_selection_item) }
    subject(:opts) { task.payload.opts }
    subject(:input_path) { opts.detect { |o| o.starts_with?('input/') } }
    subject(:output_path) { opts.detect { |o| o.starts_with?('output/') } }

    it { expect(task.state).to eq('created') }
    it { expect(task.payload.input_rclone_path).to eq(files.last) }
    it { expect(task.payload.output_rclone_path).to eq('storage:/records/output/2023-07-10/CAM1/') }

    it { expect(input_path).to eq('input/video.avi') }
    it { expect(output_path).to eq('output/video.mp4') }

    it do
      expect(opts).to eq([
                           '-i', input_path,
                           '-c:v', 'h264',
                           '-b:v', '5M',
                           '-c:a', 'copy',
                           output_path
                         ])
    end

    context 'when no common paths' do
      let(:files) do
        %w[
          storage:/records1/archive/2023-07-09/CAM1/video.avi
          storage:/records2/archive/2023-07-09/CAM2/video.avi
          storage:/records3/archive/2023-07-09/MIC1/audio.ogg
          storage:/records4/archive/2023-07-10/CAM1/video.avi
        ]
      end

      it { expect(task.payload.output_rclone_path).to eq("storage:/records/output/#{task.id}/") }
    end

    context 'when audio' do
      let(:file_selection_item) do
        index = files.find_index { |f| f.match?(/audio/) }
        file_selection.items[index]
      end

      it { expect(input_path).to eq('input/audio.ogg') }
      it { expect(output_path).to eq('output/audio.m4a') }

      it do
        expect(opts).to eq([
                             '-i', input_path,
                             '-c:a', 'aac',
                             output_path
                           ])
      end
    end

    describe 'media_meta_report' do
      subject(:media_meta_report) { task.payload.media_meta_report }

      context 'when has no matching report' do
        it { expect(media_meta_report).to be_blank }
      end

      context 'when media report with matching path & size present' do
        let!(:existing_media_meta_report) do
          create(:media_meta_report, rclone_path: file_selection_item.rclone_path,
                                     size: file_selection_item.size)
        end

        it { expect(media_meta_report).to eq(existing_media_meta_report) }
      end
    end

  end

  describe '#perform' do
    subject(:result) { instance.perform }
    subject(:errors) { instance.errors }

    it { expect(result).to eq(true) }
    it { expect(errors.full_messages).to be_empty }

    it do
      expect { instance.perform }
        .to change { Task.where(order_id: order.id).count }.by(files.size)
    end

    context 'when called twice' do
      describe 'not recreates tasks' do
        it do
          instance.perform
          expect { instance.perform }
            .to_not change { Task.where(order_id: order.id).count }
        end

        it do
          instance.perform
          expect { instance.perform }
            .to_not change { Task.where(order_id: order.id).pluck(:id) }
        end
      end
    end

    context 'when used paths on multiple levels' do
      let(:files) do
        [
          'storage:/r/records/movies/RKN2/Telegram Desktop.mov',
          'storage:/r/records/movies/RKN2/raw/2018-10-28 13.24.10.mov'
        ]
      end

      it do
        instance.perform
        output_paths = order.tasks.map { |t| t.payload.output_rclone_path }
        expect(output_paths).to contain_exactly(
          'storage:/records/output/',
          'storage:/records/output/raw/'
        )
      end
    end
  end
end
