# This file is auto-generated from the current state of the database.

ActiveRecord::Schema[7.1].define(version: 2024_05_01_120000) do
  # These are extensions that must be enabled in order to support this database
  enable_extension "plpgsql"

  create_table "categories", force: :cascade do |t|
    t.string "name", null: false
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
  end

  create_table "people", force: :cascade do |t|
    t.string "email", default: "", null: false
    t.bigint "category_id"
    t.index ["email"], name: "index_people_on_email", unique: true
  end

  create_table "posts", id: :uuid, force: :cascade do |t|
    t.references :person
    t.string "title", limit: 200
    t.text "body"
    t.bigint "author_id"
    t.decimal "price", precision: 10, scale: 2
    t.check_constraint "char_length(title) > 0", name: "title_present"
  end

  create_table "comments", force: :cascade do |t|
    t.bigint "post_id", null: false
    t.bigint "person_id"
  end

  add_foreign_key "comments", "posts", on_delete: :cascade
  add_foreign_key "posts", "people", column: "author_id", primary_key: "id"
  add_check_constraint "comments", "post_id > 0", name: "positive_post"
end
