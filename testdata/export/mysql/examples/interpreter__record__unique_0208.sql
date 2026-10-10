CREATE TABLE `users` (
  `id` int PRIMARY KEY,
  `email` varchar(255) UNIQUE,
  `phone` varchar(255) UNIQUE
);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `users` (`id`, `email`, `phone`)
VALUES
  (1, 'alice@example.com', NULL),
  (2, 'bob@example.com', NULL),
  (3, NULL, '123-456'),
  (4, NULL, '789-012');

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;