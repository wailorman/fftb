class AddTaskStatuses < ActiveRecord::Migration[7.0]
  def change
    add_column :tasks, :current_step, :string
    add_column :tasks, :current_progress, :float, default: 0, null: false
  end
end
