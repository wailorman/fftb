module Rpc
  class BasePresenter
    attr_reader :object,
                :options

    def initialize(object, options = {})
      @object = object
      @options = options
    end

    def call
      raise NotImplementedError
    end
  end
end
