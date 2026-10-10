ActiveRecord::Schema[7.1].define(version: 2024_09_09_090909) do
  create_table "accounts", charset: "utf8mb4", collation: "utf8mb4_0900_ai_ci", force: :cascade do |t|
    t.string "name", null: false
    t.integer "plan", limit: 1, default: 0, null: false, unsigned: true
    t.bigint "owner_id"
    t.text "notes", size: :medium
    t.datetime "created_at", precision: nil, null: false
    t.datetime "updated_at", precision: nil, null: false
    t.index ["owner_id"], name: "fk_rails_owner"
  end

  create_table "owners", charset: "utf8mb4", force: :cascade do |t|
    t.string "email"
    t.boolean "admin", default: false
  end

  add_foreign_key "accounts", "owners", on_update: :cascade
end
