CREATE TABLE `tasks` (
  `id` int PRIMARY KEY,
  `name` varchar(255),
  `priority` ENUM ('low', 'medium', 'high')
);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `tasks` (`id`, `name`, `priority`)
VALUES
  (1, 'Task 1', 'high'),
  (2, 'Task 2', 'invalid_priority');

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;