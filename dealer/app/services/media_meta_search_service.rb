class MediaMetaSearchService < ApplicationService
  attr_reader :tuples

  def initialize(tuples = [])
    super()
    @tuples = tuples.map { |t| t.slice(:full_path, :size) }
  end

  def perform
  end
end
