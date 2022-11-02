class HandlerGenerator < Rails::Generators::NamedBase
  source_root File.expand_path('templates', __dir__)

  def create_handler_file
    create_file "app/handlers/#{file_path}_handler.rb", <<~FILE
      class #{class_name}Handler < ApplicationHandler
        # include ::Dealer::SetTask
        # include ::Dealer::AuthorizePerformer

        def execute
          # Available methods:
          # current_performer

          raise NotImplementedError
        end
      end
    FILE
  end
end
