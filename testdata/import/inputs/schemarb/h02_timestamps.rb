ActiveRecord::Schema.define(version: 1) do
  create_table "users" do |t|
    t.string "name"
    t.timestamps
  end

  create_table "accounts" do |t|
    t.integer "user_id"
  end
end
