Table "a" {
  "x" int
}

Table "b" {
  "y" int
}

Table "c" {
  "z" int
}

Dep {
  "a" -> "c"
  "b" -> "c"

  note: 'built nightly'
}

Dep {
  "a"."x" -> "c"."z"
  "b"."y" -> "c"."z"
}
