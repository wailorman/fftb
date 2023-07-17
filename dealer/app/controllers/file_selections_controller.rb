class FileSelectionsController < ApplicationController
  before_action :set_file_selection, only: %i[show destroy]
  before_action :set_file_selection_item, only: %i[hide_item show_item]

  helper_method :media_meta_for_item

  def create
    service = SearchFilesService.new(
      rclone_path: params[:path],
      mime_type: /^(video|audio)\/.+/,
      file_name: /(^[^.]).+/
    )

    unless service.perform
      redirect_back alert: "Failed to search files: #{service.errors.full_messages.join(', ')}"
      return
    end

    @file_selection = service.file_selection
    session[:current_file_selection_id] = @file_selection.id

    redirect_to file_selection_path(id: @file_selection.id)
  end

  def show
    @items =
      FileSelectionItem
        .unscoped
        .where(file_selection: @file_selection)
        .order(rclone_path: :asc)

    meta_tuples = @items.map do |item|
      [item.rclone_path, item.size]
    end

    @media_meta_reports =
      MediaMetaReport
        .by_path_and_size(meta_tuples)
        .order(created_at: :asc)
        .index_by { |media_meta_report| [media_meta_report.rclone_path, media_meta_report.size] }
  end

  def hide_item
  end

  def show_item
  end

  def destroy
    @file_selection.destroy
    redirect_to remote_files_path(remote_name: @remote_name, path: @remote_path)
  end

  private def set_file_selection
    @file_selection = FileSelection.find(params[:id])

    parsed_path = Rclone.parse(@file_selection.root_rclone_path)
    @remote_name = parsed_path[:remote]
    @remote_path = parsed_path[:path]
  end

  private def media_meta_for_item(item)
    @media_meta_reports[[item.rclone_path, item.size]]
  end
end
