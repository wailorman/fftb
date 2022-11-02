class CreateStorageClaims < ActiveRecord::Migration[7.0]
  def change
    create_table :storage_claims, id: :uuid do |t|
      t.string :kind
      t.string :provider
      t.string :path

      t.timestamps
    end
  end
end
