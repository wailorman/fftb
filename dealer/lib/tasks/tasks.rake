namespace :tasks do
  task generate: :environment do
    task =
      Tasks::Convert.create!(
        convert_task_payload: ConvertTaskPayload.new(
          opts: %w[
            -i input/RE7-1_.mp4
            -c:v h264
            -b:v 5M
            -c:a aac
            output/RE7-1_out.mp4
          ],
          input_rclone_path: 'novus_smb:/r/records/movies/_fftb_test/RE7-1_.mp4',
          output_rclone_path: "novus_smb:/r/records/fftb_test/output/#{SecureRandom.hex(2)}/"
        )
      )

    puts "Created task #{task.id}"
  end
end
