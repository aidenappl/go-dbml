CREATE TABLE `users` (
  `id` int PRIMARY KEY,
  `email` varchar(255) UNIQUE
);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `users` (`id`, `email`)
VALUES
  (1, 'alice@example.com'),
  (1, 'bob@example.com');

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;