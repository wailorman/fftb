namespace :tasks do
  task generate: :environment do
    input_storage_claim = StorageClaim.new(kind: :s3, provider: :yandex, path: 'aerial_shot_of_a_lighthouse.mp4')

    task =
      Task.create!(kind: :convert_v1,
                   state: :published,
                   input_storage_claim: input_storage_claim,
                   params: {
                     video_codec: :h264,
                     hw_accel: nil,
                     video_bit_rate: nil,
                     video_quality: 28,
                     preset: 'fast',
                     scale: nil,
                     keyframe_interval: nil,
                     muxer: 'mp4',
                     position: nil
                   })

    puts "Created task #{task.id}"
  end
end
