CREATE TABLE `data` (
  `id` int,
  `active` boolean
);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `data` (`id`, `active`)
VALUES
  (1, FALSE),
  (2, FALSE),
  (3, FALSE),
  (4, FALSE),
  (5, FALSE);

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;