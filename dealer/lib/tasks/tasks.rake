namespace :tasks do
  task generate: :environment do
    # task =
      # Tasks::Convert.create!(
      #   convert_task_payload: ConvertTaskPayload.new(
      #     opts: %w[
      #       -i input/RE7-1_.mp4
      #       -c:v h264
      #       -b:v 5M
      #       -c:a aac
      #       output/RE7-1_out.mp4
      #     ],
      #     input_rclone_path: 'novus_smb:/r/records/movies/_fftb_test/RE7-1_.mp4',
      #     output_rclone_path: "novus_smb:/r/records/fftb_test/output/#{SecureRandom.hex(2)}/"
      #   )
      # )

    # task =
      # Tasks::MediaMeta.create!(
      #   media_meta_task_payload: MediaMetaTaskPayload.new(
      #     input_rclone_path: 'novus_smb:/r/records/movies/_fftb_test/RE7-1_.mp4',
      #     output_rclone_path: "novus_smb:/r/records/fftb_test/output/#{SecureRandom.hex(2)}/"
      #   )
      # )
      Tasks::MediaMeta.create!(
        payload: MediaMetaTaskPayload.new(
          input_rclone_path: 'novus_smb:/r/records/movies/Серёжа/Composit_0127_095859.mpg',
          output_rclone_path: "novus_smb:/r/records/fftb_test/output/#{SecureRandom.hex(2)}/"
        )
      )

    # [20, 22, 24, 26, 28].map do |crf|
    #   Tasks::Convert.create!(
    #     convert_task_payload: ConvertTaskPayload.new(
    #       opts: [
    #         '-i', 'input/SnowRunner __  KSIVA_ p_61a0d 09-01-2022 21-29-13.mp4',
    #         '-c:v', 'h264',
    #         '-b:v', '5M',
    #         '-c:a', 'aac',
    #         'output/SnowRunner __  KSIVA_ p_61a0d 09-01-2022 21-29-13.mp4'
    #       ],
    #       input_rclone_path: 'novus_smb:/r/records/_STORE/S/SnowRunner 2022/2022-01-09/W/SnowRunner __  KSIVA_ p_61a0d 09-01-2022 21-29-13.mp4',
    #       output_rclone_path: "novus_smb:/r/records/fftb_test/SR/#{crf}/"
    #     )
    #   )
    # end

    # puts "Created task #{task.id}"
  end
end
