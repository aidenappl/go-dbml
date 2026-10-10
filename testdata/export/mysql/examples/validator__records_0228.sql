CREATE TABLE `users` (
  `id` int PRIMARY KEY,
  `name` varchar(255)
);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `users` (`id`, `name`)
VALUES
  ;

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;