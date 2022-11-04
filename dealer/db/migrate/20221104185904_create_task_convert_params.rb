class CreateTaskConvertParams < ActiveRecord::Migration[7.0]
  def change
    create_table :convert_params, id: :uuid do |t|
      t.string :video_codec
      t.string :hw_accel
      t.string :video_bit_rate
      t.integer :video_quality
      t.string :preset
      t.string :scale
      t.integer :keyframe_interval
      t.string :muxer

      t.timestamps
    end

    add_reference :tasks, :convert_params, foreign_key: true, type: :uuid
    remove_column :tasks, :params
  end
end
