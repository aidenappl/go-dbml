CREATE TABLE `users` (
  `created_at` timestamp,
  `updated_at` timestamp,
  `id` int PRIMARY KEY,
  `name` varchar(255),
  `email` varchar(255)
);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `users` (`created_at`, `updated_at`, `id`, `name`, `email`)
VALUES
  ('2024-01-01T00:00:00', '2024-01-01T00:00:00', 1, 'Alice', 'alice@example.com'),
  ('2024-01-02T00:00:00+08:00', '2024-01-02T00:00:00.100+07:00', 2, 'Bob', 'bob@example.com');

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;