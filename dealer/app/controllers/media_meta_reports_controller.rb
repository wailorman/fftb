class MediaMetaReportsController < ApplicationController
  before_action :set_media_meta_report, except: [:create]

  def show; end

  def create
    Task.transaction do
      paths.each do |path|
        payload = Payloads::Tasks::MediaMeta.new(
          input_rclone_path: path,
          output_rclone_path: Rclone.join_path(config[:media_meta_location], Date.today.to_s, SecureRandom.hex(8), '/')
        )

        Tasks::MediaMeta.create!(payload: payload)
        # payload.validate!
      end

      redirect_back fallback_location: remotes_path, alert: 'Analyze enqueued'
    end
  end

  private def set_media_meta_report
    @media_meta_report = MediaMetaReport.find(params[:id])
  end

  private def permitted_params
    params.permit(:paths, :path, :file_selection_id)
  end

  private def paths
    @paths ||=
      Array(
        permitted_params[:path].presence ||
        permitted_params[:paths].presence ||
        FileSelection.find_by(id: permitted_params[:file_selection_id]).items.pluck(:rclone_path)
      )
  end
end
