class Dealer::AllocateOutputStorageClaimHandler < ApplicationHandler
  include ::Dealer::SetTask
  include ::Dealer::AuthorizePerformer

  def execute
    # TODO: use default provider

    storage_claim_id = SecureRandom.uuid

    claim =
      OutputStorageClaim.new(id: storage_claim_id,
                             task: task,
                             kind: :s3,
                             provider: :yandex,
                             purpose: req.purpose.downcase.to_sym,
                             name: req.name,
                             path: "claims/#{storage_claim_id}/#{req.name}")

    return Twirp::Error.unknown(claim.full_messages.join(', ')) unless claim.save

    ::Rpc::StorageClaimPresenter.new(claim, access_type: :put).call
  end
end

