class Dealer::AllocateOutputStorageClaimHandler < ApplicationHandler
  include ::Dealer::SetTask
  include ::Dealer::AuthorizePerformer

  def execute
    # TODO: use default provider
    task.output_storage_claim = StorageClaim.new(kind: :s3, provider: :yandex, path: "claims/#{SecureRandom.uuid}")

    return Twirp::Error.unknown(task.full_messages.join(', ')) unless task.save

    signer = S3UrlSignService.new(task.output_storage_claim)

    Fftb::StorageClaim.new(url: signer.put(expires_in: StorageClaim::DEFAULT_URL_TTL))
  end
end

