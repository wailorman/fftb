class Dealer::NotifyHandler < ApplicationHandler
  def execute
    task = Task.find(req.segmentId)

    return Twirp::Error.not_found('Segment not found') unless task

    task.current_step = req.step.to_s.downcase
    task.current_progress = req.progress
    task.occupied_at = Time.current
    task.occupied_by = current_performer

    return Twirp::Error.unknown(task.full_messages.join(', ')) unless task.save

    Fftb::Empty.new
  end
end
