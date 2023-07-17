module Remotes
  class FilesController < ApplicationController
    helper_method :media_meta_report_for_entry

    def index
      @remote = params[:remote_name]

      unless @remote
        redirect_to remotes_path, alert: 'Missing remote name'
        return
      end

      @entries =
        Rclone.ls("#{params[:remote_name]}:#{File.join('/', params[:path])}")
              .reject { |entry| ignore_entry?(entry) }
              .sort_by { |entry| [entry[:is_dir] ? 0 : 1, entry[:name]] }

      meta_tuples = @entries.map do |entry|
        [entry[:full_path], entry[:size]]
      end

      @media_meta_reports =
        MediaMetaReport.by_path_and_size(meta_tuples)
                       .order(created_at: :asc)
                       .index_by { |media_meta_report| [media_meta_report.rclone_path, media_meta_report.size] }
    end

    private def ignore_entry?(entry)
      entry[:name].match?(/DS_Store/) || \
      (entry[:name].start_with?('.') && entry[:size] == 4.kilobytes)
    end

    private def media_meta_report_for_entry(entry)
      @media_meta_reports[[entry[:full_path], entry[:size]]]
    end
  end
end
