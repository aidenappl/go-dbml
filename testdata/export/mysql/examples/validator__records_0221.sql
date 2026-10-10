CREATE TABLE `data` (
  `id` int,
  `value` decimal
);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `data` (`id`, `value`)
VALUES
  (1, 10000000000),
  (2, 0.0000314),
  (3, 200000000);

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;