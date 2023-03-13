class Dealer::NotifyHandler < ApplicationHandler
  include ::Dealer::SetTask
  include ::Dealer::AuthorizePerformer

  before_execute :authorize_performer
  before_execute :authorize_performer_task

  def execute
    task.current_step = req.step.to_s.downcase
    task.current_progress = req.progress
    task.occupied_at = Time.current
    task.occupied_by = current_performer

    return Twirp::Error.unknown(task.full_messages.join(', ')) unless task.save

    Fftb::Empty.new
  end
end
