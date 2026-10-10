ActiveRecord::Schema.define(version: 2023_11_20_101010) do
  create_table "products", id: :uuid, default: -> { "gen_random_uuid()" }, force: :cascade do |t|
    t.string "name", limit: 120, null: false
    t.decimal "price", precision: 10, scale: 2, default: "0.0", null: false
    t.jsonb "properties", default: {}
    t.string "tags", default: [], array: true
    t.timestamps
  end

  create_table "line_items", force: :cascade do |t|
    t.uuid "product_id", null: false
    t.integer "order_id", null: false
    t.integer "quantity", default: 1, null: false
    t.index ["order_id", "product_id"], name: "index_line_items_on_order_id_and_product_id", unique: true
  end

  create_table "orders", force: :cascade do |t|
    t.string "status", default: "cart"
    t.references "customer", foreign_key: true
    t.datetime "placed_at", precision: 6
  end

  create_table "customers", force: :cascade do |t|
    t.string "name"
  end

  add_foreign_key "line_items", "orders", on_delete: :cascade
  add_foreign_key "line_items", "products", column: "product_id", primary_key: "id"
end
