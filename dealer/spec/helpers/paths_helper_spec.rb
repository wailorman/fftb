require 'rails_helper'

RSpec.describe PathsHelper do
  describe '#generalize_paths' do
    let(:paths) { nil }
    subject { described_class.generalize_paths(paths) }

    context 'when passed nil' do
      let(:paths) { nil }

      it { expect(subject).to eq(nil) }
    end

    context 'when passed []' do
      let(:paths) { [] }

      it { expect(subject).to eq(nil) }
    end

    context 'when passed [nil]' do
      let(:paths) { [nil] }

      it { expect(subject).to eq(nil) }
    end

    context 'when passed 1 path' do
      let(:paths) { ['storage:/r/example/1.mov'] }

      it { expect(subject).to eq('storage:/r/example/') }
    end

    context 'when passed 3 similar paths' do
      let(:paths) do
        [
          'storage:/r/example/movies/first/1.mov',
          'storage:/r/example/movies/second/1.mov',
          'storage:/r/example/movies/1.mov'
        ]
      end

      it { expect(subject).to eq('storage:/r/example/movies/') }
    end

    context 'when passed 3 similar short paths' do
      let(:paths) do
        [
          'storage:/r/1.mov',
          'storage:/r/2.mov',
          'storage:/r/3.mov'
        ]
      end

      it { expect(subject).to eq('storage:/r/') }
    end

    context 'when passed 3 similar very short paths' do
      let(:paths) do
        [
          'storage:/1.mov',
          'storage:/2.mov',
          'storage:/3.mov'
        ]
      end

      it { expect(subject).to eq(nil) }
    end

    context 'when passed 2 paths from different remotes' do
      let(:paths) do
        [
          'storage:/r/example/movies/first/1.mov',
          'cloud:/r/example/movies/second/1.mov'
        ]
      end

      it { expect(subject).to eq(nil) }
    end
  end

  describe '#common_path' do
    let(:path_a) { nil }
    let(:path_b) { nil }
    subject { described_class.common_path(path_a, path_b) }

    context 'when passed nil' do
      let(:path_a) { nil }
      let(:path_b) { nil }

      it { expect(subject).to eq(nil) }
    end

    context 'when passed 2 similar paths' do
      let(:path_a) { 'storage:/r/example/movies/first/1.mov' }
      let(:path_b) { 'storage:/r/example/movies/second/1.mov' }

      it { expect(subject).to eq('storage:/r/example/movies/') }
    end

    context 'when passed 2 similar short paths' do
      let(:path_a) { 'storage:/r/1.mov' }
      let(:path_b) { 'storage:/r/2.mov' }

      it { expect(subject).to eq('storage:/r/') }
    end

    context 'when passed 2 similar very short paths' do
      let(:path_a) { 'storage:/1.mov' }
      let(:path_b) { 'storage:/2.mov' }

      it { expect(subject).to eq(nil) }
    end

    context 'when passed 2 paths from different remotes' do
      let(:path_a) { 'storage:/r/example/movies/first/1.mov' }
      let(:path_b) { 'cloud:/r/example/movies/second/1.mov' }

      it { expect(subject).to eq(nil) }
    end
  end
end
