class CreateTasks < ActiveRecord::Migration[7.0]
  def change
    create_table :tasks, id: :uuid do |t|
      t.string :kind
      t.jsonb :params
      t.string :state
      t.datetime :occupied_until

      t.timestamps
    end
  end
end
