require 'rails_helper'

RSpec.describe OptsHelper do
  describe '#array_to_string' do
    let(:arr) { ['-i', 'abc 123', '-c:v', 'h264'] }

    subject { described_class.array_to_string(arr) }

    it { expect(subject).to eq("-i \"abc 123\" \n-c:v h264") }

    context 'when empty elements' do
      let(:arr) { ['-c:v', '', 'h264'] }

      it { expect(subject).to eq('-c:v h264') }
    end
  end

  describe '#string_to_array' do
    let(:str) { '-i "abc 123" -c:v h264' }

    subject { described_class.string_to_array(str) }

    it { expect(subject).to eq(['-i', 'abc 123', '-c:v', 'h264']) }

    context 'when newlines' do
      let(:str) { "-c:v h264\n-c:a aac" }

      it { expect(subject).to eq(['-c:v', 'h264', '-c:a', 'aac']) }
    end
  end
end
