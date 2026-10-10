CREATE TABLE `products` (
  `id` int,
  `price` decimal(5,2),
  `rate` numeric(10,3)
);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `products` (`id`, `price`, `rate`)
VALUES
  (1, 99.99, 1.234),
  (2, 12345.67, 12345.678),
  (3, 99.999, 1234567.89);

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;