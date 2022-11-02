class RenameOccupiedUntil < ActiveRecord::Migration[7.0]
  def change
    rename_column :tasks, :occupied_until, :occupied_at
  end
end
