class InitTables < ActiveRecord::Migration[7.0]
  def change
    create_table :performers, id: :uuid do |t|
      t.string :name
      t.string :token

      t.timestamps
    end

    create_table :tasks, id: :uuid do |t|
      t.string :type, null: false
      t.string :payload_type, null: false
      t.uuid :payload_id, null: false
      t.string :state, default: 'published'
      t.datetime :occupied_at
      t.references :occupied_by, foreign_key: { to_table: :performers }, type: :uuid
      t.string :current_step
      t.float :current_progress, default: 0, null: false
      t.boolean :result_verified, default: false

      t.timestamps
    end

    create_table :convert_task_payloads, id: :uuid do |t|
      t.string :opts, array: true, null: false
      t.string :input_rclone_path
      t.string :output_rclone_path

      t.bigint :current_time
      t.float :current_fps
      t.bigint :current_frame
      t.float :current_speed
      t.float :current_bitrate

      t.timestamps
    end

    create_table :media_meta_reports, id: :uuid do |t|
      t.string :rclone_path, index: true
      t.jsonb :data
      t.string :video_codec
      t.string :video_codec_long
      t.bigint :video_bitrate
      t.string :audio_codec
      t.string :audio_codec_long
      t.bigint :audio_bitrate
      t.bigint :bitrate
      t.float :duration
      t.bigint :size, index: true
      t.integer :resolution_w
      t.integer :resolution_h
      t.string :pix_fmt
      t.datetime :created_at_by_meta
      t.datetime :created_at_by_name
      t.datetime :created_at_by_mtime

      t.timestamps
    end

    add_index :media_meta_reports, %i[rclone_path size]
    add_reference :convert_task_payloads, :media_meta_report, foreign_key: true, type: :uuid, null: true

    create_table :media_meta_task_payloads, id: :uuid do |t|
      t.references :media_meta_report, foreign_key: { to_table: :media_meta_reports }, type: :uuid, index: { unique: true }
      t.string :input_rclone_path
      t.string :output_rclone_path

      t.timestamps
    end

    create_table :task_failures, id: :uuid do |t|
      t.references :task, foreign_key: { to_table: :tasks }, type: :uuid, null: false
      t.references :performer, foreign_key: { to_table: :performers }, type: :uuid
      t.string :reason

      t.timestamps
    end

    create_table :file_selections, id: :uuid do |t|
      t.string :root_rclone_path
      t.boolean :reached_max_depth, default: false, null: false
      t.integer :items_count
      t.timestamps
    end

    create_table :file_selection_items, id: :uuid do |t|
      t.references :file_selection, null: false, foreign_key: true, type: :uuid
      t.string :rclone_path
      t.string :mime_type
      t.bigint :size
      t.boolean :removed, default: false, null: false

      t.timestamps
    end

    add_reference :tasks, :file_selection_item, null: true, foreign_key: true, type: :uuid

    create_table :orders, id: :uuid do |t|
      t.string :type
      t.string :state, default: 'created'
      t.string :payload_type, null: false
      t.uuid :payload_id, null: false
      t.references :file_selection, null: true, foreign_key: true, type: :uuid

      t.timestamps
    end

    add_reference :tasks, :order, null: true, foreign_key: true, type: :uuid

    create_table :convert_order_payloads, id: :uuid do |t|
      t.string :video_muxer
      t.string :video_opts
      t.string :audio_muxer
      t.string :audio_opts
      t.string :output_rclone_path

      t.timestamps
    end
  end
end
