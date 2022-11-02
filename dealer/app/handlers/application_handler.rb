class ApplicationHandler
  attr_reader :req, :env

  def initialize(req, env)
    @req = req
    @env = env
  end

  def call
    execute
  rescue => e
    filtered_backtrace = e.backtrace.select { |line| Rails.root.to_s.in?(line) }
    Rails.logger.error "#{e}\n#{filtered_backtrace.join("\n")}"
    raise e
  end

  private

  def execute
    raise NotImplementedError
  end
end
