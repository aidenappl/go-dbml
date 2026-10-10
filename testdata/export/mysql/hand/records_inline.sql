CREATE TABLE `posts` (
  `id` int PRIMARY KEY,
  `title` text
);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `posts` (`id`, `title`)
VALUES
  (1, 'Hello'),
  (2, 'World');

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;