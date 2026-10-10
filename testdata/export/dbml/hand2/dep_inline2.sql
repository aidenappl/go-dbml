Table "src" {
  "v" int
}

Table "mid" {
  "v" int
}

Table "out" {
  "w" int
}

Dep: "src" -> "mid"

Dep: "src"."v" -> "mid"."v"

Dep: "out" -> "mid"
