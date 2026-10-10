CREATE TABLE `users` (
  `id` integer PRIMARY KEY,
  `name` varchar(255),
  `active` boolean
);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `users` (`id`, `name`, `active`)
VALUES
  (1, 'Alice', NULL),
  (2, 'Bob', NULL),
  (NULL, 'Charlie', TRUE),
  (3, NULL, FALSE);

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;