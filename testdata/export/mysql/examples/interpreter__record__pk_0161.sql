CREATE TABLE `users` (
  `id` int AUTO_INCREMENT,
  `tenant_id` int,
  PRIMARY KEY (`id`, `tenant_id`)
);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `users` (`tenant_id`)
VALUES
  (NULL),
  (1);

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;