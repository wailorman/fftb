class SearchFilesService < ApplicationService
  include ActiveModel::Validations

  attr_reader :rclone_path, :mime_type, :file_selection

  def initialize(opts = {})
    super()
    @rclone_path = opts[:rclone_path]
    @mime_type = opts[:mime_type] || /.+/
    @file_name = opts[:file_name] || /.+/
    @max_depth = opts[:max_depth] || 5
    @max_items = opts[:max_items] || 100
    @file_selection = FileSelection.new(root_rclone_path: opts[:rclone_path])
  end

  def perform
    scan_files(@rclone_path, @max_depth)
    @file_selection.save!

    true
  end

  private def scan_files(path, remaning_depth)
    if remaning_depth <= 0
      @file_selection.reached_max_depth = true
      return []
    end

    Rclone.ls(path).each do |entry|
      break if @file_selection.items.size >= @max_items

      if entry[:is_dir]
        scan_files(entry[:full_path], remaning_depth - 1)
        next
      end

      next unless @mime_type.match?(entry[:mime_type])
      next unless @file_name.match?(entry[:name])

      @file_selection.items <<
        FileSelectionItem.new(rclone_path: entry[:full_path],
                              mime_type: entry[:mime_type],
                              size: entry[:size])
    end
  end
end
