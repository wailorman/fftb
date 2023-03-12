class DealerHandler
  include Twirp::Rails::Helpers

  bind Fftb::DealerService

  def finish_task(req, env)
    Dealer::FinishTaskHandler.new(req, env).call
  end

  def quit_task(req, env)
    Dealer::QuitTaskHandler.new(req, env).call
  end

  def fail_task(req, env)
    Dealer::FailTaskHandler.new(req, env).call
  end

  def find_free_task(req, env)
    Dealer::FindFreeTaskHandler.new(req, env).call
  end

  def notify(req, env)
    Dealer::NotifyHandler.new(req, env).call
  end
end
