CREATE TABLE `users` (
  `id` int,
  `name` varchar(5),
  `quantity` int
);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `users` (`id`, `name`, `quantity`)
VALUES
  (1, 'Alice', 10),
  (2, 'Bob', 20),
  (3, 'Christopher', 30),
  (4, 'Dave', 40),
  (5, 'Elizabeth', 50);

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;