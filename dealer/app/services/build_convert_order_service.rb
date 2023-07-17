class BuildConvertOrderService < ApplicationService
  include ActiveModel::Validations

  delegate :file_selection, to: :order

  attr_reader :order

  validates :file_selection, presence: true
  validate :ensure_order_persisted

  def initialize(order)
    super()
    @order = order
  end

  def perform
    return false if invalid?

    Order.transaction do
      order.tasks = file_selection.items.not_removed.map do |item|
        task = build_task(item)
        task.save!
        task
      rescue ActiveRecord::RecordInvalid => e
        handle_task_invalid(item, e)
        raise e
      end

    rescue ActiveRecord::RecordInvalid => e
      errors.add(:base, e.record.errors.full_messages.join(', ')) unless e.record.kind_of?(Task)
    end

    errors.empty?
  end

  private def handle_task_invalid(file_selection_item, exception)
    errors.add(
      :base,
      format(
        'Failed to create task for file `%{path}`: %{errors}',
        path: file_selection_item.rclone_path,
        errors: exception.record.errors.full_messages.join(', ')
      )
    )
  end

  private def build_task(file_selection_item)
    id = SecureRandom.uuid

    input_parsed = Rclone.parse(file_selection_item.rclone_path)

    t = Tasks::Convert
          .where(order_id: order.id, file_selection_item_id: file_selection_item.id)
          .first_or_initialize

    t.id ||= id
    t.file_selection_item = file_selection_item
    t.order = order
    t.state = 'created' if t.new_record?
    t.payload = Payloads::Tasks::Convert.new

    t.payload.media_meta_report =
      MediaMetaReport.find_by(rclone_path: file_selection_item.rclone_path,
                              size: file_selection_item.size)

    t.payload.string_opts = build_opts(file_selection_item)

    t.payload.input_rclone_path = file_selection_item.rclone_path

    t.payload.output_rclone_path =
      if common_path
        common_parsed = Rclone.parse(common_path)
        exclusive_path = Pathname.new(Pathname.new(input_parsed[:path]).relative_path_from(common_parsed[:path]).to_s).parent.to_s

        Rclone.join_path(order.payload.output_rclone_path, exclusive_path, '/')
      else
        Rclone.join_path(order.payload.output_rclone_path, t.id, '/')
      end

    t
  end

  private def build_opts(file_selection_item)
    input_parsed = Rclone.parse(file_selection_item.rclone_path)

    input_path = "input/#{Pathname.new(input_parsed[:path]).basename}"

    if file_selection_item.video?
      opts_template = order.payload.video_opts
      output_path = format('output/%{name}.%{ext}', name: basename_no_ext(input_parsed[:path]),
                                                    ext: order.payload.video_muxer)
    elsif file_selection_item.audio?
      opts_template = order.payload.audio_opts
      output_path = format('output/%{name}.%{ext}', name: basename_no_ext(input_parsed[:path]),
                                                    ext: order.payload.audio_muxer)
    else
      raise "Unsupported mime_type: `#{file_selection_item.mime_type}`"
    end

    format(
      opts_template,
      input_path: input_path,
      output_path: output_path,
      basename: basename_no_ext(input_parsed[:path])
    )
  end

  private def basename_no_ext(file_path)
    Pathname.new(file_path).basename.to_s.gsub(Pathname.new(file_path).extname.to_s, '')
  end

  private def common_path
    return @common_path if defined?(@common_path)

    @common_path = PathsHelper.generalize_paths(file_selection.items.pluck(:rclone_path))
  end

  private def ensure_order_persisted
    errors.add(:order, 'Order have to be persisted') unless order.persisted?
  end
end
