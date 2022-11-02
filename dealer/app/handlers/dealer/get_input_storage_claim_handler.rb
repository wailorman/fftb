class Dealer::GetInputStorageClaimHandler < ApplicationHandler
  before_execute :set_task
  before_execute :authorize_performer

  attr_accessor :task

  def execute
    return Twirp::Error.not_found('StorageClaim not found') unless task.input_storage_claim

    signer = S3UrlSignService.new(task.input_storage_claim)

    Fftb::StorageClaim.new(url: signer.get(expires_in: StorageClaim::DEFAULT_URL_TTL))
  end

  private

  def set_task
    self.task = Task.find(req.segmentId)
  end

  def authorize_performer
    Twirp::Error.permission_denied('performer mismatch') if task.occupied_by != current_performer
  end
end
