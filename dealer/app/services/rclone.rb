class Rclone
  class Config
    attr_accessor :binary_path, :config_path

    def binary_path # rubocop:disable Lint/DuplicateMethods
      @binary_path ||= 'rclone'
    end

    def config_path # rubocop:disable Lint/DuplicateMethods
      @config_path ||= '/Users/wailorman/Resilio/wailorman_fftb/rclone.cfg'
    end
  end

  class << self
    def config
      @config ||= Config.new
    end

    def configure
      yield(config)
    end

    # path
    # name
    # size
    # mime_type
    # mod_time
    # is_dir
    # full_path

    def ls(path)
      stdout = exec('lsjson', q(path))

      JSON.parse(stdout)
          .map do |entry|
            entry.transform_keys! { |key| key.underscore.to_sym }
            entry.merge(
              mod_time: entry[:mod_time].presence && Time.zone.parse(entry[:mod_time]),
              full_path: File.join(path, entry[:name])
            )
          end
    end

    def remotes
      exec('listremotes').split("\n").map { |remote| remote[0..-2] }
    end

    def read(path)
      exec('cat', q(path))
    end

    def parse(path)
      splitted = path.split(':', 2)

      if splitted.size == 1
        return {
          remote: nil,
          path: Pathname.new('/').join('/', splitted).to_s
        }
      end

      {
        remote: splitted[0].presence,
        path: Pathname.new('/').join('/', *splitted[1..]).to_s
      }
    end

    def join_path(root, *siblings)
      raise 'Missing root location' unless root

      root_location = parse(root)

      [
        root_location[:remote],
        (':' if root_location[:remote]),
        # Pathname.new('/').join(File.join(*[root_location[:path], siblings].flatten)).expand_path.to_s
        # Pathname.new('/').join(*[root_location[:path], siblings].flatten).to_s
        File.join(*[root_location[:path],
                    siblings].flatten.reject{ |s| s.to_s.starts_with?('.') })
      ].compact.join
    end

    def exec(*options)
      stdout, stderr, status = ExecOperation.new.perform([config.binary_path,
                                                          ("--config #{config.config_path}" if config.config_path),
                                                          options])

      unless status.success?
        stderr_lines = stderr.split("\n")
        raise "Failed to execute rclone: `#{stderr_lines.last}`"
      end

      stdout
    end

    def q(str)
      "\"#{str}\""
    end
  end
end
