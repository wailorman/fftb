class ExecOperation
  attr_reader :log_prefix

  def initialize(params = {})
    @log_prefix = params[:log_prefix] || 'exec'
  end

  def perform(*args)
    command = [args].flatten.compact.join(' ')
    log("command: #{command}")
    stdout, stderr, status = Open3.capture3(command)

    unless status.success?
      lines(stdout).each { |line| log("stdout: #{line}") }
      lines(stderr).each { |line| log("stderr: #{line}") }
    end

    log("status: #{status.to_i}")

    [stdout, stderr, status]
  end

  private def log(str)
    Rails.logger.debug { "<#{operation_id}> #{log_prefix}: #{str}" }
  end

  private def quote(str)
    "\"#{str}\""
  end

  private def operation_id
    @operation_id ||= SecureRandom.hex[0..8]
  end

  private def lines(str)
    str.to_s.split("\n")
  end
end
