class Dealer::FindFreeTaskHandler < ApplicationHandler
  include ::Dealer::AuthorizePerformer

  before_execute :authorize_performer

  def execute
    run = FindFreeTask.run(performer: current_performer)

    return Twirp::Error.invalid_argument(run.errors.full_messages.join(';')) unless run.valid?
    return Twirp::Error.not_found('Free task not found') unless run.result

    Rpc::TaskPresenter.new(run.result).call
  end
end
