class SetDefaultTaskState < ActiveRecord::Migration[7.0]
  def change
    change_column_default :tasks, :state, 'published'
  end
end
