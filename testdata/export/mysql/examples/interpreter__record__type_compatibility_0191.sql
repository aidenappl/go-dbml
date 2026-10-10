CREATE TABLE `users` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `status` varchar(255) NOT NULL DEFAULT 'active'
);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `users` (`id`, `status`)
VALUES
  (NULL, NULL),
  (NULL, 'inactive');

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;