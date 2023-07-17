module PathsHelper
  class << self
    def generalize_paths(paths = [])
      return nil if Array(paths).compact.blank?

      first_parsed = Rclone.parse(paths.first)

      return nil if Pathname.new(first_parsed[:path]).parent.to_s == '/'

      if paths.size == 1
        return Rclone.join_path(
          "#{first_parsed[:remote]}:",
          Pathname.new(first_parsed[:path]).parent,
          '/'
        )
      end

      Array(paths[1..].compact).inject(paths[0]) do |prev, cur|
        common = common_path(prev, cur)
        break nil if common.blank?

        common
      end
    end

    def common_path(path_a, path_b)
      return nil if path_a.blank? || path_b.blank?

      parsed_a = Rclone.parse(path_a)
      parsed_b = Rclone.parse(path_b)

      return nil if parsed_a[:remote] != parsed_b[:remote]

      path_length = Pathname.new(parsed_b[:path]).each_filename.to_a.size

      common = Array.new(path_length).inject(Pathname.new(parsed_b[:path])) do |short, _cur|
        break nil if short.parent.to_s == '/'
        break short.parent if parsed_a[:path].starts_with?(short.parent.to_s)

        short.parent
      end

      return nil if common.blank?

      Rclone.join_path("#{parsed_a[:remote]}:", common, '/')
    end
  end
end
