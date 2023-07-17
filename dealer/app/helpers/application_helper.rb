module ApplicationHelper
  def human_errors(errors)
    errors.full_messages.join('; ')
  end

  def link_to_path(rclone_path, &block)
    parsed = Rclone.parse(rclone_path)

    ActionController::Base.helpers.link_to(remote_files_path(remote_name: parsed[:remote], path: parsed[:path])) do
      capture(&block)
    end
  end
end
