class DealerHandler
  include Twirp::Rails::Helpers

  bind Fftb::DealerService

  def get_all_input_storage_claims(req, env)
    Dealer::GetAllInputStorageClaimsHandler.new(req, env).call
  end

  def allocate_output_storage_claim(req, env)
    Dealer::AllocateOutputStorageClaimHandler.new(req, env).call
  end

  def finish_segment(req, env)
    Dealer::FinishSegmentHandler.new(req, env).call
  end

  def quit_segment(req, env)
    Dealer::QuitSegmentHandler.new(req, env).call
  end

  def fail_segment(req, env)
    Dealer::FailSegmentHandler.new(req, env).call
  end

  def find_free_segment(req, env)
    Dealer::FindFreeSegmentHandler.new(req, env).call
  end

  def notify(req, env)
    Dealer::NotifyHandler.new(req, env).call
  end
end
