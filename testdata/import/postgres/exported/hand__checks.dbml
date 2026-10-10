Table "products" {
  "price" decimal [check: `price > 0`]
  "discount" decimal

  Checks {
    `discount >= 0` [name: 'disc_nonneg']
    `discount < price`
  }
}
