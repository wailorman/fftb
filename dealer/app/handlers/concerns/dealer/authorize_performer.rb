module Dealer
  module AuthorizePerformer
    extend ActiveSupport::Concern

    private

    def authorize_performer_task
      Twirp::Error.permission_denied('performer mismatch') if task.occupied_by != current_performer
      Twirp::Error.unknown('task cancelled') if task.state.cancelled?
    end

    def current_performer
      return nil unless @authorization_payload

      @current_performer ||= Performer.find_or_create_by!(name: @authorization_payload.dig(0, 'worker_name'))
    end

    def authorize_performer
      return Twirp::Error.permission_denied('missing token') if req.authorization.blank?

      @authorization_payload = JWT.decode(
        req.authorization,
        Rails.application.config.application_options[:secret],
        true,
        { algorithm: 'HS256' }
      )

      if @authorization_payload.dig(0, 'worker_name').blank?
        return Twirp::Error.permission_denied('missing worker_name in token')
      end

      nil
    rescue JWT::ImmatureSignature, JWT::VerificationError
      Twirp::Error.permission_denied('invalid signature')
    end
  end
end
