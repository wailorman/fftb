class Dealer::QuitSegmentHandler < ApplicationHandler
  include ::Dealer::SetTask
  include ::Dealer::AuthorizePerformer

  def execute
    task.state = :published
    task.deoccupy

    return Twirp::Error.unknown(task.full_messages.join(', ')) unless task.save

    Fftb::Empty.new
  end
end
