class TaskDecorator < ApplicationDecorator
  delegate_all

  def human_progress
    return '' if !task.occupied? || !task.state.published?

    [
      human_extended_progress,
      (human_mpps || human_fps)
    ].compact.join(' ')
  end

  def human_step
    return unless task.occupied?
    return unless task.state.published?

    task.current_step
  end

  def human_fps
    return unless task.short_type.convert?
    return unless task.current_step == 'processing'
    return if task.payload.current_fps.blank?

    "@#{task.payload.current_fps.to_i} FPS"
  end

  def human_mpps
    return unless task.short_type.convert?
    return unless task.current_step == 'processing'
    return if task.payload.media_meta_report.blank?

    fps = task.payload.current_fps
    return if task.payload.current_fps.blank?

    pixels_per_frame = task.payload.media_meta_report.pixels_per_frame
    return if pixels_per_frame.blank?

    mpps = fps * (pixels_per_frame.to_f / 1_000_000)

    "@#{mpps.round} MPPS"
  end

  private def human_basic_progress
    "#{(task.current_progress || 0) * 100}%"
  end

  private def human_extended_progress
    return human_extended_convert_progress if task.short_type.convert?

    human_basic_progress
  end

  private def human_extended_convert_progress
    return human_basic_progress unless task.short_type.convert?
    return human_basic_progress unless task.current_step == 'processing'
    return human_basic_progress unless task.payload.media_meta_report&.duration&.present?

    current_time = task.payload.current_time.presence || 0
    meta_duration = task.payload.media_meta_report.duration * 1000

    format(
      '%{progress}%%',
      progress: (current_time / meta_duration * 100).to_i
    )
  end
end
