module Rpc
  class SegmentPresenter < ::Rpc::BasePresenter
    alias task object

    def call
      Fftb::Segment.new(
        type: segment_type,
        id: task.id,
        convertParams: ConvertParamsPresenter.new(task.convert_params).call
      )
    end

    private

    def segment_type
      Fftb::SegmentType::CONVERT_V1
    end
  end
end
