class ApplicationRecord < ActiveRecord::Base
  default_scope { order(created_at: :asc, id: :asc) }

  primary_abstract_class

  def config
    @config ||= Rails.application.config.application_options
  end
end
