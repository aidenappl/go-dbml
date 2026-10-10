CREATE TABLE `users` (
  `id` bigserial PRIMARY KEY,
  `name` varchar(255)
);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `users` (`id`, `name`)
VALUES
  (NULL, 'Alice'),
  (NULL, 'Bob');

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;