class ApplicationHandler
  attr_reader :req, :env

  delegate :callbacks, to: :class

  def initialize(req, env)
    @req = req
    @env = env
  end

  def call
    callbacks.each do |method_name|
      callback_result = send(method_name.to_sym)

      return callback_result if callback_result.is_a?(Twirp::Error)
    end

    execute
  rescue ActiveRecord::RecordNotFound => e
    handle_not_found(e)
  rescue => e # rubocop:disable Style/RescueStandardError
    log_exception(e)
    raise e
  end

  def self.before_execute(method_name)
    @callbacks ||= []
    @callbacks << method_name
  end

  def self.callbacks
    @callbacks ||= []
  end

  private

  def execute
    raise NotImplementedError
  end

  def empty_response
    Fftb::Empty.new
  end

  def current_performer
    @current_performer ||= Performer.local
  end

  def log_exception(e)
    filtered_backtrace = e.backtrace.select { |line| Rails.root.to_s.in?(line) }
    Rails.logger.error "#{e}\n#{filtered_backtrace.join("\n")}"
  end

  def handle_not_found(error)
    if error.present? && error.respond_to?(:model)
      model_name = error.model.underscore

      return Twirp::Error.not_found("#{model_name} not found")
    end

    log_exception(error)
  end
end
