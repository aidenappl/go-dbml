Table "A"."a" {
  "AB" integer
  "BA" integer

  Indexes {
    (AB, BA) [pk]
  }
}

Table "B"."b" {
  "BC" integer
  "CB" integer

  Indexes {
    (BC, CB) [pk]
  }
}

Table "C"."c" {
  "CD" integer [pk]
  "DC" integer
}

Table "D"."d" {
  "DE" integer [pk]
  "ED" integer
}

Table "E"."e" {
  "EF" integer
  "FE" integer
  "DE" integer
  "ED" integer

  Indexes {
    (EF, FE) [pk]
  }
}

Table "G"."g" {
  "GH" integer
  "HG" integer
  "EH" integer
  "HE" integer

  Indexes {
    (GH, HG) [pk]
  }
}

Table "t1" {
  "a" int [pk]
  "b" int [unique]
}

Table "t2" {
  "a" int [pk]
  "b" int [unique]
}

Table "t1_t2" {
  "a" int
}

Ref:"D"."d"."DE" <> "C"."c"."CD"

Ref:"A"."a".("AB", "BA") <> "B"."b".("BC", "CB")

Ref:"E"."e".("EF", "FE") <> "G"."g".("GH", "HG")

Ref:"t1"."a" <> "t2"."a"

Ref:"t1"."b" <> "t2"."b"

Ref:"schema"."image"."id" <> "schema"."content_item"."id"

Ref:"schema1"."customers"."id" <> "schema2"."orders"."id"

Table "schema"."image" {
  "id" integer [pk]
  "url" varchar
}

Table "schema"."content_item" {
  "id" integer [pk]
  "heading" varchar
  "description" varchar
}

Table "schema"."footer_item" {
  "id" integer [pk]
  "left" varchar
  "centre" varchar
  "right" varchar
}

Table "schema1"."customers" {
  "id" integer [pk]
  "full_name" varchar
}

Table "schema2"."orders" {
  "id" integer [pk]
  "total_price" integer
}
