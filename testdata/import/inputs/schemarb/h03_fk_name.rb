ActiveRecord::Schema.define(version: 2) do
  create_table "teams" do |t|
    t.string "name"
  end
  create_table "members" do |t|
    t.integer "team_id"
  end
  add_foreign_key "members", "teams", name: "fk_custom"
end
