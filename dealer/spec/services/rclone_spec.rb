require 'rails_helper'

RSpec.describe Rclone do
  describe '#join_path' do
    it { expect(Rclone.join_path('storage:/abc', '.')).to eq('storage:/abc') }
    it { expect(Rclone.join_path('storage:/abc', '/')).to eq('storage:/abc/') }
  end
end
