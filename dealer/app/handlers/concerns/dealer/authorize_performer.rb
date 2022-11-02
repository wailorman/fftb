module Dealer
  module AuthorizePerformer
    extend ActiveSupport::Concern

    included do
      before_execute :authorize_performer
    end

    private

    def authorize_performer
      Twirp::Error.permission_denied('performer mismatch') if task.occupied_by != current_performer
    end
  end
end
