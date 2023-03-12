class InitTables < ActiveRecord::Migration[7.0]
  def change
    create_table :performers, id: :uuid do |t|
      t.string :name
      t.string :token

      t.timestamps
    end

    create_table :tasks, id: :uuid do |t|
      t.string :type, null: false
      t.string :state, default: 'published'
      t.datetime :occupied_at
      t.references :occupied_by, foreign_key: { to_table: :performers }, type: :uuid
      t.string :current_step
      t.float :current_progress, default: 0, null: false

      t.timestamps
    end

    create_table :convert_task_payloads, id: :uuid do |t|
      t.references :task, foreign_key: { to_table: :tasks }, type: :uuid, null: false, index: {unique: true}
      t.string :opts, array: true, null: false
      t.string :input_rclone_path, null: false
      t.string :output_rclone_path, null: false

      t.timestamps
    end

    create_table :task_failures, id: :uuid do |t|
      t.references :task, foreign_key: { to_table: :tasks }, type: :uuid, null: false
      t.references :performer, foreign_key: { to_table: :performers }, type: :uuid, null: false
      t.string :reason

      t.timestamps
    end
  end
end
