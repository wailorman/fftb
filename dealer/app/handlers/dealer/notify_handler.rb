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

    if task.short_type.convert?
      task.payload.current_bitrate = req.convertProgress.bitrate
      task.payload.current_fps = req.convertProgress.fps
      task.payload.current_frame = req.convertProgress.frame
      task.payload.current_speed = req.convertProgress.speed
      task.payload.current_time = req.convertProgress.time
    end

    return Twirp::Error.unknown(task.full_messages.join(', ')) unless task.save

    Fftb::Empty.new
  end
end
