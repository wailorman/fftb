class Dealer::GetAllInputStorageClaimsHandler < ApplicationHandler
  include ::Dealer::SetTask
  include ::Dealer::AuthorizePerformer

  def execute
    return Twirp::Error.not_found('storage claims not found') unless task.input_storage_claims.exists?

    claims = task.input_storage_claims.map do |claim|
      Rpc::StorageClaimPresenter.new(claim, access_type: :get).call
    end

    Fftb::StorageClaimList.new(storageClaims: claims)
  end
end
