module Rpc
  class TaskPresenter < ::Rpc::BasePresenter
    alias task object

    def call
      Fftb::Task.new(
        type: task_type,
        id: task.id,
        convertParams: convert_params,
        mediaMetaParams: media_meta_params
      )
    end

    private def task_type
      case task.type
      when 'Tasks::Convert'
        Fftb::Task::TaskType::CONVERT_V1
      when 'Tasks::MediaMeta'
        Fftb::Task::TaskType::MEDIA_META_V1
      else
        raise NotImplementedError
      end
    end

    private def convert_params
      return nil if task.type != 'Tasks::Convert'

      Fftb::ConvertTaskParams.new(
        inputRclonePath: task.payload.input_rclone_path,
        outputRclonePath: task.payload.output_rclone_path,
        opts: task.payload.opts
      )
    end

    private def media_meta_params
      return nil if task.type != 'Tasks::MediaMeta'

      Fftb::MediaMetaTaskParams.new(
        inputRclonePath: task.payload.input_rclone_path,
        outputRclonePath: task.payload.output_rclone_path
      )
    end
  end
end
