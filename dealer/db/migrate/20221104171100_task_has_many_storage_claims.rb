class TaskHasManyStorageClaims < ActiveRecord::Migration[7.0]
  def change
    remove_reference :tasks, :input_storage_claim
    remove_reference :tasks, :output_storage_claim

    add_column :storage_claims, :type, :string
    add_reference :storage_claims, :task, foreign_key: true, type: :uuid
  end
end
