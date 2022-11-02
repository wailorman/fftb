class AddOccupiedByToTasks < ActiveRecord::Migration[7.0]
  def change
    add_reference :tasks, :occupied_by, foreign_key: { to_table: :performers }, type: :uuid
  end
end
