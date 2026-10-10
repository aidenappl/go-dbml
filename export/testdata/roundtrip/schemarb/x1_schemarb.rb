ActiveRecord::Schema[7.1].define(version: 2024_05_01_000000) do
  enable_extension "plpgsql"
  create_enum "mood", ["it's ok", "sad"]
  create_table "accts", comment: "it's \"quoted\"" do |t|
    t.string "name", default: "O'Brien \\ x", null: false, comment: "a\nb"
    t.enum "mood", enum_type: "mood", default: "it's ok"
    t.virtual "full", type: :string, as: "name || 'x'", stored: true
    t.references "parent", foreign_key: { to_table: :accts, on_delete: :nullify }
    t.check_constraint "length(name) > 0", name: "name_len"
    t.timestamps
  end
end
