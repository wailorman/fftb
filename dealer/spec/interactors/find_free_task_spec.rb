require 'rails_helper'

RSpec.describe FindFreeTask, type: :interactor do
  describe '#run' do
    let!(:task) { create(:convert_task) }
    let!(:performer) { create(:performer) }

    subject { described_class.run!(performer: performer) }

    context 'when task not occupied by anyone' do
      it { expect(subject).to eq(task) }
    end

    context 'when task is occupied by another performer' do
      let!(:another_performer) { create(:performer) }
      let!(:task) { create(:convert_task, :occupied, occupied_by: another_performer) }

      it { expect(subject).to eq(nil) }
    end

    context 'when occupation of task has been expired' do
      let!(:another_performer) { create(:performer) }
      let!(:task) { create(:convert_task, :occupied, occupied_by: another_performer, occupied_at: 2.years.ago) }

      it { expect(subject).to eq(task) }
    end

    context 'when task was failed by current performer' do
      let!(:task_failure) { create(:task_failure, task: task, performer: performer) }

      it { expect(subject).to eq(nil) }
    end
  end
end
