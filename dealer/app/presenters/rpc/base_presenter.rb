module Rpc
  class BasePresenter
    attr_reader :object

    def initialize(object)
      @object = object
    end

    def call
      raise NotImplementedError
    end
  end
end
