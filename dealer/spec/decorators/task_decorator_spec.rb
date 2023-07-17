require 'rails_helper'

RSpec.describe TaskDecorator do
  describe '#human_progress' do
    let!(:task) { create(:convert_task) }

    subject { task.decorate.human_progress }

    context 'when Convert' do
      before do
        task.payload.current_bitrate = 1_382_092.0
        task.payload.current_fps = 240
        task.payload.current_frame = 45_621
        task.payload.current_speed = 1.6
        task.payload.current_time = 2_412_423
      end

      context 'when no media meta' do
        it { expect(subject).to eq('@240 FPS') }
      end

      context 'when media meta present' do
        before do
          task.payload.update!(
            media_meta_report: create(:media_meta_report, duration: (2_412_423 * 2))
          )
        end

        it { expect(subject).to eq('50% @240 FPS') }
      end
    end
  end
end
