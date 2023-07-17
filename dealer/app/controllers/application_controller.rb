class ApplicationController < ActionController::Base
  include ApplicationHelper

  helper_method :current_file_selection

  def config
    @config ||= Rails.application.config.application_options
  end

  def current_file_selection
    return nil if session[:current_file_selection_id].blank?
    return @current_file_selection if defined?(@current_file_selection)

    @current_file_selection = FileSelection.find_by(id: session[:current_file_selection_id])

    session[:current_file_selection_id] = nil unless @current_file_selection

    @current_file_selection
  end
end
