CREATE TABLE `data` (
  `id` int,
  `value` int
);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `data` (`id`, `value`)
VALUES
  (1, -100),
  (2, -999);

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;