CREATE TABLE `products` (
  `id` integer PRIMARY KEY,
  `name` nvarchar(255),
  `price` decimal(10,2),
  `in_stock` bit
);

CREATE TABLE `orders` (
  `id` integer PRIMARY KEY,
  `product_id` integer,
  `quantity` integer
);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `products` (`id`, `name`, `price`, `in_stock`)
VALUES
  (1, 'Laptop', 999.99, TRUE),
  (2, 'Mouse', 29.99, TRUE),
  (3, 'Keyboard', 79.99, FALSE),
  (4, 'Monitor', 149.99, TRUE);
INSERT INTO `orders` (`id`, `product_id`, `quantity`)
VALUES
  (1, 1, 2),
  (2, 2, 5),
  (3, 1, 1);

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;