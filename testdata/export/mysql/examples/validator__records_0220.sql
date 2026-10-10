CREATE TABLE `users` (
  `id` int PRIMARY KEY,
  `created_at` timestamp
);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `users` (`id`, `created_at`)
VALUES
  (1, now()),
  (2, uuid_generate_v4());

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;