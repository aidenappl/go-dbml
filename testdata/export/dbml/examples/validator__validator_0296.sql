Table "t1" {
  "c1" int
}

Table "t2" {
  "c2" int
}

Dep: "t2"."c2" -> "t1"."c1"
