CREATE TABLE `products` (
  `id` int PRIMARY KEY,
  `created_at` timestamp,
  `name` varchar(255),
  `price` decimal
);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `products` (`id`, `name`, `price`)
VALUES
  (1, 'Widget', 9.99),
  (2, 'Gadget', 19.99);

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;