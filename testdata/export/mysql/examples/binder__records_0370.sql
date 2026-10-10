CREATE SCHEMA `auth`;

CREATE TABLE `auth`.`users` (
  `id` int,
  `email` varchar(255)
);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `auth`.`users` (`id`, `email`)
VALUES
  (1, 'alice@example.com');

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;