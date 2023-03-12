FactoryBot.define do
  factory :convert_task_payload do
    opts do
      %w[
        -i input/example.mp4
        -c:v h264
        -b:v 5M
        -c:a aac
        output/example.mp4
      ]
    end
    input_rclone_path { 'storage_smb:/test/in/example.mp4' }
    output_rclone_path { 'storage_smb:/test/out/example.mp4' }
  end
end
