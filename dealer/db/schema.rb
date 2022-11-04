# This file is auto-generated from the current state of the database. Instead
# of editing this file, please use the migrations feature of Active Record to
# incrementally modify your database, and then regenerate this schema definition.
#
# This file is the source Rails uses to define your schema when running `bin/rails
# db:schema:load`. When creating a new database, `bin/rails db:schema:load` tends to
# be faster and is potentially less error prone than running all of your
# migrations from scratch. Old migrations may fail to apply correctly if those
# migrations use external dependencies or application code.
#
# It's strongly recommended that you check this file into your version control system.

ActiveRecord::Schema[7.0].define(version: 2022_11_04_185904) do
  # These are extensions that must be enabled in order to support this database
  enable_extension "pgcrypto"
  enable_extension "plpgsql"

  create_table "convert_params", id: :uuid, default: -> { "gen_random_uuid()" }, force: :cascade do |t|
    t.string "video_codec"
    t.string "hw_accel"
    t.string "video_bit_rate"
    t.integer "video_quality"
    t.string "preset"
    t.string "scale"
    t.integer "keyframe_interval"
    t.string "muxer"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
  end

  create_table "performers", id: :uuid, default: -> { "gen_random_uuid()" }, force: :cascade do |t|
    t.string "name"
    t.string "token"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
  end

  create_table "storage_claims", id: :uuid, default: -> { "gen_random_uuid()" }, force: :cascade do |t|
    t.string "kind"
    t.string "provider"
    t.string "path"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.string "name"
    t.string "purpose", default: "none", null: false
    t.string "type"
    t.uuid "task_id"
    t.index ["task_id"], name: "index_storage_claims_on_task_id"
  end

  create_table "tasks", id: :uuid, default: -> { "gen_random_uuid()" }, force: :cascade do |t|
    t.string "kind"
    t.string "state", default: "published"
    t.datetime "occupied_at"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.uuid "occupied_by_id"
    t.string "current_step"
    t.float "current_progress", default: 0.0, null: false
    t.string "failure"
    t.uuid "convert_params_id"
    t.index ["convert_params_id"], name: "index_tasks_on_convert_params_id"
    t.index ["occupied_by_id"], name: "index_tasks_on_occupied_by_id"
  end

  add_foreign_key "storage_claims", "tasks"
  add_foreign_key "tasks", "convert_params", column: "convert_params_id"
  add_foreign_key "tasks", "performers", column: "occupied_by_id"
end
