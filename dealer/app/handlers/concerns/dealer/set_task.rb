module Dealer
  module SetTask
    extend ActiveSupport::Concern

    included do
      attr_accessor :task

      before_execute :set_task
    end

    private

    def set_task
      self.task = Task.find(req.segmentId)
    end
  end
end
