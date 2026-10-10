CREATE TABLE `users` (
  `extra` varchar(255),
  `id` int PRIMARY KEY,
  `name` varchar(255)
);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `users` (`extra`, `id`, `name`)
VALUES
  ('extra_value', 1, 'Alice'),
  ('extra_value2', 2, 'Bob');

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;