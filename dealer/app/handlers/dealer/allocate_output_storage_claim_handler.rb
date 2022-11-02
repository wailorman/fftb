class Dealer::AllocateOutputStorageClaimHandler < ApplicationHandler
  before_execute :set_task
  before_execute :authorize_performer

  attr_accessor :task

  def execute
    # TODO: use default provider
    task.output_storage_claim = StorageClaim.new(kind: :s3, provider: :yandex, path: "claims/#{SecureRandom.uuid}")

    return Twirp::Error.unknown(task.full_messages.join(', ')) unless task.save

    signer = S3UrlSignService.new(task.output_storage_claim)

    Fftb::StorageClaim.new(url: signer.put(expires_in: StorageClaim::DEFAULT_URL_TTL))
  end

  private

  def set_task
    self.task = Task.find(req.segmentId)
  end

  def authorize_performer
    Twirp::Error.permission_denied('performer mismatch') if task.occupied_by != current_performer
  end
end

