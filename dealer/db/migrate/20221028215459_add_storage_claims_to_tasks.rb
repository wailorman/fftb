class AddStorageClaimsToTasks < ActiveRecord::Migration[7.0]
  def change
    add_reference :tasks, :input_storage_claim, foreign_key: { to_table: :storage_claims }, type: :uuid
    add_reference :tasks, :output_storage_claim, foreign_key: { to_table: :storage_claims }, type: :uuid
  end
end
