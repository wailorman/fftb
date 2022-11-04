module Rpc
  class ConvertParamsPresenter < ::Rpc::BasePresenter
    alias convert_params object

    def call
      return nil unless convert_params

      Fftb::ConvertSegmentParams.new(
        convert_params
          .attributes
          .except(*%w[id created_at updated_at])
          .transform_keys { |key| key.camelize(:lower) }
          .transform_keys(&:to_sym)
      )
    end
  end
end
