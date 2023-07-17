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

ActiveRecord::Schema[7.0].define(version: 2022_10_27_205416) do
  # These are extensions that must be enabled in order to support this database
  enable_extension "pgcrypto"
  enable_extension "plpgsql"

  create_table "convert_order_payloads", id: :uuid, default: -> { "gen_random_uuid()" }, force: :cascade do |t|
    t.string "video_muxer"
    t.string "video_opts"
    t.string "audio_muxer"
    t.string "audio_opts"
    t.string "output_rclone_path"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
  end

  create_table "convert_task_payloads", id: :uuid, default: -> { "gen_random_uuid()" }, force: :cascade do |t|
    t.string "opts", null: false, array: true
    t.string "input_rclone_path"
    t.string "output_rclone_path"
    t.bigint "current_time"
    t.float "current_fps"
    t.bigint "current_frame"
    t.float "current_speed"
    t.float "current_bitrate"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.uuid "media_meta_report_id"
    t.index ["media_meta_report_id"], name: "index_convert_task_payloads_on_media_meta_report_id"
  end

  create_table "file_selection_items", id: :uuid, default: -> { "gen_random_uuid()" }, force: :cascade do |t|
    t.uuid "file_selection_id", null: false
    t.string "rclone_path"
    t.string "mime_type"
    t.bigint "size"
    t.boolean "removed", default: false, null: false
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["file_selection_id"], name: "index_file_selection_items_on_file_selection_id"
  end

  create_table "file_selections", id: :uuid, default: -> { "gen_random_uuid()" }, force: :cascade do |t|
    t.string "root_rclone_path"
    t.boolean "reached_max_depth", default: false, null: false
    t.integer "items_count"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
  end

  create_table "media_meta_reports", id: :uuid, default: -> { "gen_random_uuid()" }, force: :cascade do |t|
    t.string "rclone_path"
    t.jsonb "data"
    t.string "video_codec"
    t.string "video_codec_long"
    t.bigint "video_bitrate"
    t.string "audio_codec"
    t.string "audio_codec_long"
    t.bigint "audio_bitrate"
    t.bigint "bitrate"
    t.float "duration"
    t.bigint "size"
    t.integer "resolution_w"
    t.integer "resolution_h"
    t.string "pix_fmt"
    t.datetime "created_at_by_meta"
    t.datetime "created_at_by_name"
    t.datetime "created_at_by_mtime"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["rclone_path", "size"], name: "index_media_meta_reports_on_rclone_path_and_size"
    t.index ["rclone_path"], name: "index_media_meta_reports_on_rclone_path"
    t.index ["size"], name: "index_media_meta_reports_on_size"
  end

  create_table "media_meta_task_payloads", id: :uuid, default: -> { "gen_random_uuid()" }, force: :cascade do |t|
    t.uuid "media_meta_report_id"
    t.string "input_rclone_path"
    t.string "output_rclone_path"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["media_meta_report_id"], name: "index_media_meta_task_payloads_on_media_meta_report_id", unique: true
  end

  create_table "orders", id: :uuid, default: -> { "gen_random_uuid()" }, force: :cascade do |t|
    t.string "type"
    t.string "state", default: "created"
    t.string "payload_type", null: false
    t.uuid "payload_id", null: false
    t.uuid "file_selection_id"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["file_selection_id"], name: "index_orders_on_file_selection_id"
  end

  create_table "performers", id: :uuid, default: -> { "gen_random_uuid()" }, force: :cascade do |t|
    t.string "name"
    t.string "token"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
  end

  create_table "task_failures", id: :uuid, default: -> { "gen_random_uuid()" }, force: :cascade do |t|
    t.uuid "task_id", null: false
    t.uuid "performer_id"
    t.string "reason"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["performer_id"], name: "index_task_failures_on_performer_id"
    t.index ["task_id"], name: "index_task_failures_on_task_id"
  end

  create_table "tasks", id: :uuid, default: -> { "gen_random_uuid()" }, force: :cascade do |t|
    t.string "type", null: false
    t.string "payload_type", null: false
    t.uuid "payload_id", null: false
    t.string "state", default: "published"
    t.datetime "occupied_at"
    t.uuid "occupied_by_id"
    t.string "current_step"
    t.float "current_progress", default: 0.0, null: false
    t.boolean "result_verified", default: false
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.uuid "file_selection_item_id"
    t.uuid "order_id"
    t.index ["file_selection_item_id"], name: "index_tasks_on_file_selection_item_id"
    t.index ["occupied_by_id"], name: "index_tasks_on_occupied_by_id"
    t.index ["order_id"], name: "index_tasks_on_order_id"
  end

  add_foreign_key "convert_task_payloads", "media_meta_reports"
  add_foreign_key "file_selection_items", "file_selections"
  add_foreign_key "media_meta_task_payloads", "media_meta_reports"
  add_foreign_key "orders", "file_selections"
  add_foreign_key "task_failures", "performers"
  add_foreign_key "task_failures", "tasks"
  add_foreign_key "tasks", "file_selection_items"
  add_foreign_key "tasks", "orders"
  add_foreign_key "tasks", "performers", column: "occupied_by_id"
end
