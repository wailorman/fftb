module Rpc
  class StorageClaimPresenter < ::Rpc::BasePresenter
    alias claim object

    def call
      raise "Unknown access_type type: `#{options[:access_type]}`" unless options[:access_type].in?([:get, :put])

      Fftb::StorageClaim.new(
        id: claim.id,
        url: signer.send(options[:access_type], expires_in: StorageClaim::DEFAULT_URL_TTL),
        purpose: purpose,
        name: claim.name
      )
    end

    private

    def signer
      @signer ||= S3UrlSignService.new(claim)
    end

    def purpose
      case claim.purpose.to_s
      when 'none'
        Fftb::StorageClaimPurpose::NONE
      when 'convert_input'
        Fftb::StorageClaimPurpose::CONVERT_INPUT
      when 'convert_output'
        Fftb::StorageClaimPurpose::CONVERT_OUTPUT
      else
        raise "Unknown storage claim purpose: `#{claim.purpose}`"
      end
    end
  end
end
