module Dealer
  class GetInputStorageClaimHandler < ApplicationHandler
    def execute
      task = Task.find_by(id: req.segmentId)

      return Twirp::Error.not_found('Segment not found') unless task
      return Twirp::Error.not_found('StorageClaim not found') unless task.input_storage_claim

      signer = S3UrlSignService.new(task.input_storage_claim)

      Fftb::StorageClaim.new(url: signer.get(expires_in: StorageClaim::DEFAULT_URL_TTL))
    end
  end
end
