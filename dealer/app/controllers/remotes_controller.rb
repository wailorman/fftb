class RemotesController < ApplicationController
  def index
    @remotes = Rclone.remotes
  end
end
