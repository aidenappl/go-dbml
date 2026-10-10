CREATE TABLE `users` (
  `id` int PRIMARY KEY,
  `email` varchar(255) UNIQUE
);

CREATE TABLE `products` (
  `id` int PRIMARY KEY,
  `sku` varchar(255) UNIQUE
);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `users` (`id`, `email`)
VALUES
  (1, 'alice@example.com'),
  (2, 'alice@example.com');
INSERT INTO `products` (`id`, `sku`)
VALUES
  (1, 'PROD-001'),
  (2, 'PROD-001');

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;