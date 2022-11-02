class ApplicationHandler
  include ActiveSupport::Callbacks
  define_callbacks :execute, terminator: ->(_target, result_lambda) { result_lambda.call.is_a?(Twirp::Error) }

  attr_reader :req, :env

  def initialize(req, env)
    @req = req
    @env = env
  end

  def call
    run_callbacks(:execute) do
      execute
    end
  rescue ActiveRecord::NotFound => e
    handle_not_found(e)
  rescue => e
    log_exception(e)
    raise e
  end

  def self.before_execute(method_name)
    set_callback :execute, :before, method_name
  end

  private

  def execute
    raise NotImplementedError
  end

  def current_performer
    @current_performer ||= Performer.local
  end

  def log_exception(e)
    filtered_backtrace = e.backtrace.select { |line| Rails.root.to_s.in?(line) }
    Rails.logger.error "#{e}\n#{filtered_backtrace.join("\n")}"
  end

  def handle_not_found(error)
    if error && error.respond_to?(:model)
      model_name = error.model.underscore

      return Twirp::Error.not_found("#{model_name} not found")
    end

    log_exception(error)
  end
end
