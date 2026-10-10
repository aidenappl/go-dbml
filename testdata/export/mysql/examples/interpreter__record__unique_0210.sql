CREATE TABLE `users` (
  `id` int PRIMARY KEY,
  `email` varchar(255) UNIQUE,
  `username` varchar(255) UNIQUE,
  `phone` varchar(255) UNIQUE
);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `users` (`id`, `email`, `username`, `phone`)
VALUES
  (1, 'alice@example.com', 'alice', '111-111'),
  (2, 'bob@example.com', 'bob', '222-222'),
  (3, 'charlie@example.com', 'alice', '333-333'),
  (4, 'dave@example.com', 'dave', '111-111');

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;