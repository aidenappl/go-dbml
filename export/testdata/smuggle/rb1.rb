ActiveRecord::Schema[7.1].define(version: 2024_01_01_000000) do
  create_table "t", force: :cascade do |t|
    t.string "a", default: -> { "concat('x\\'', ')); DROP TABLE victim; CREATE TABLE pwned (i int); -- ')" }
    t.string "b", comment: "x'); DROP TABLE victim; --"
    t.check_constraint "a <> 'y\\')); DROP TABLE victim; CREATE TABLE pwned2 (i int); -- '", name: "ck"
  end
end
