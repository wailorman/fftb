module Dealer
  class FindFreeSegmentHandler < ApplicationHandler
    def execute
      run = FindFreeTask.run(performer: Performer.local)

      return Twirp::Error.invalid_argument(run.errors.full_messages.join(';')) unless run.valid?
      return Twirp::Error.not_found('Free task not found') unless run.result

      Rpc::SegmentPresenter.new(run.result).call
    end
  end
end
