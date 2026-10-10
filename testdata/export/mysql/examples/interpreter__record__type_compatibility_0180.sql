CREATE TABLE `data` (
  `id` int,
  `quantity` int,
  `price` decimal
);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `data` (`id`, `quantity`, `price`)
VALUES
  (not_a_number, 100, 100),
  (2, 10, 100),
  (3, 20, also_not_a_number),
  (4, 30, true);

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;