Table "users" {
  "id" integer
  "wealth" integer
  "debt" integer

  Checks {
    `debt + wealth >= 0` [name: 'chk_positive_money']
  }
}
