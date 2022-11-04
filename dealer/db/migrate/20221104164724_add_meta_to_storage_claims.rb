class AddMetaToStorageClaims < ActiveRecord::Migration[7.0]
  def change
    add_column :storage_claims, :name, :string
    add_column :storage_claims, :purpose, :string, default: :none, null: false
  end
end
