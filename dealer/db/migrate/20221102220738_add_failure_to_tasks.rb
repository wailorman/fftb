class AddFailureToTasks < ActiveRecord::Migration[7.0]
  def change
    add_column :tasks, :failure, :string
  end
end
