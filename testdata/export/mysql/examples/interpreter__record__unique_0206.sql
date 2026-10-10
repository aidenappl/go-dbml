CREATE TABLE `users` (
  `id` int PRIMARY KEY,
  `email` varchar(255) UNIQUE,
  `username` varchar(255) UNIQUE
);

CREATE TABLE `products` (
  `id` int PRIMARY KEY,
  `sku` varchar(255) UNIQUE
);

CREATE TABLE `accounts` (
  `id` int PRIMARY KEY,
  `account_number` int UNIQUE
);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `users` (`id`, `email`, `username`)
VALUES
  (1, 'alice@example.com', 'alice'),
  (2, 'bob@example.com', 'bob');
INSERT INTO `products` (`id`, `sku`)
VALUES
  (1, 'PROD-001'),
  (2, 'PROD-002');
INSERT INTO `accounts` (`id`, `account_number`)
VALUES
  (1, 0),
  (2, -1),
  (3, 1000);

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;