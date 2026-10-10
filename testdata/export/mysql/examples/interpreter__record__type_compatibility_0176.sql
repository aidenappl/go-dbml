CREATE TABLE `data` (
  `id` int,
  `active` boolean
);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `data` (`id`, `active`)
VALUES
  (1, TRUE),
  (2, FALSE),
  (3, TRUE),
  (4, FALSE),
  (5, TRUE),
  (6, FALSE),
  (7, TRUE),
  (8, FALSE),
  (9, TRUE),
  (10, FALSE),
  (11, TRUE),
  (12, FALSE),
  (13, TRUE),
  (14, FALSE);

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;