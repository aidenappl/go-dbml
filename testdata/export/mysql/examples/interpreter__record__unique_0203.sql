CREATE TABLE `user_profiles` (
  `user_id` int,
  `profile_type` varchar(255),
  `data` text
);

CREATE UNIQUE INDEX `user_profiles_index_0` ON `user_profiles` (`user_id`, `profile_type`);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `user_profiles` (`user_id`, `profile_type`, `data`)
VALUES
  (1, 'work', 'Software Engineer'),
  (1, 'work', 'Updated job title');

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;