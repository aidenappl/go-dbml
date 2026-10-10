CREATE TABLE `user_settings` (
  `user_id` int,
  `category` varchar(255),
  `value` varchar(255)
);

CREATE UNIQUE INDEX `user_settings_index_0` ON `user_settings` (`user_id`, `category`);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `user_settings` (`user_id`, `category`, `value`)
VALUES
  (1, NULL, 'default'),
  (1, NULL, 'another default'),
  (1, 'theme', 'dark');

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;