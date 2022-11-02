class DealerHandler
  include Twirp::Rails::Helpers

  bind Fftb::DealerService

  def get_input_storage_claim(req, env)
    Dealer::GetInputStorageClaimHandler.new(req, env).call
  end

  def allocate_output_storage_claim(req, env)
    Dealer::AllocateOutputStorageClaimHandler.new(req, env).call
  end

  def finish_segment(req, env)
    Dealer::FinishSegmentHandler.new(req, env).call
  end

  def quit_segment(req, _env)
    Rails.logger.info "quit_segment req: #{req}"

    Fftb::Empty.new
  end

  def fail_segment(req, _env)
    Dealer::FailSegmentHandler.new(req, env).call
  end

  def find_free_segment(req, env)
    Dealer::FindFreeSegmentHandler.new(req, env).call
  end

  def notify(req, env)
    Dealer::NotifyHandler.new(req, env).call
  end
end
