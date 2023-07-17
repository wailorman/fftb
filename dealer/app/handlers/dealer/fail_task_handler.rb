class Dealer::FailTaskHandler < ApplicationHandler
  include ::Dealer::SetTask
  include ::Dealer::AuthorizePerformer

  before_execute :authorize_performer
  before_execute :authorize_performer_task

  def execute
    task.state = :failed
    task.task_failures.build(performer: current_performer, reason: req.failures.join(', '))

    return Twirp::Error.unknown(task.full_messages.join(', ')) unless task.save

    Fftb::Empty.new
  end
end
