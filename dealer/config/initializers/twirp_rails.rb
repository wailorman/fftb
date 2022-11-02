require_relative '../../lib/pb/pkg/distributed/remote/pb/fftb_twirp.rb'

Twirp::Rails.configuration do |c|
  # Modify the path below if you locates handlers under the different directory.
  c.handlers_path = Rails.root.join('app', 'controllers', 'rpc')
end
