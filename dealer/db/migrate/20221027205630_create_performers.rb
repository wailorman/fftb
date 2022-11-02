class CreatePerformers < ActiveRecord::Migration[7.0]
  def change
    create_table :performers, id: :uuid do |t|
      t.string :name
      t.string :token

      t.timestamps
    end
  end
end
