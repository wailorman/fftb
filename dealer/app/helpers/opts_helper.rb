module OptsHelper
  class << self
    def array_to_string(arr = [])
      arr
        .map(&:strip)
        .map.with_index do |opt, i|
          if opt.match?(/\s/)
            opt.inspect
          elsif opt.starts_with?('-') && i != 0
            "\n#{opt}"
          else
            opt
          end
        end
        .reject(&:blank?)
        .join(' ')
    end

    def string_to_array(str = '')
      str
        .split(/\s(?=(?:[^"]|"[^"]*")*$)/).map do |opt|
          opt.gsub(/^"|"$/, '').strip
        end
        .reject(&:blank?)
    end
  end
end
