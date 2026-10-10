ActiveRecord::Schema.define(version: 20200101000000) do
  create_table "countries", id: false, force: :cascade do |t|
    t.string "code", limit: 2, null: false
    t.string "name", null: false, comment: "English name"
  end

  create_table "regions", id: :serial, force: :cascade, comment: "Admin regions" do |t|
    t.string "country_code", limit: 2, null: false
    t.string "name", null: false
    t.float "area_km2"
    t.date "founded_on"
    t.index ["country_code", "name"], unique: true
  end

  create_table "settings", id: false, force: :cascade do |t|
    t.string "key", null: false
    t.text "value"
  end

  add_foreign_key "regions", "countries", column: "country_code", primary_key: "code", on_delete: :restrict
end
