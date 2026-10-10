CREATE TABLE `users` (
  `id` int,
  `code` varchar(5),
  `name` char(10)
);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `users` (`id`, `code`, `name`)
VALUES
  (1, 'ABCDE', 'Alice'),
  (2, 'XYZ', 'Bob'),
  (3, '', '');

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;