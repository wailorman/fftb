module Rpc
  class TaskPresenter < ::Rpc::BasePresenter
    alias task object

    def call
      Fftb::Task.new(
        type: task_type,
        id: task.id,
        convertParams: convert_params
      )
    end

    private

    def task_type
      case task.type
      when 'Tasks::Convert'
        Fftb::TaskType::CONVERT_V1
      else
        raise NotImplementedError
      end
    end

    def convert_params
      return nil if task.type != 'Tasks::Convert'

      Fftb::ConvertTaskParams.new(
        inputRclonePath: task.convert_task_payload.input_rclone_path,
        outputRclonePath: task.convert_task_payload.output_rclone_path,
        opts: task.convert_task_payload.opts
      )
    end
  end
end
