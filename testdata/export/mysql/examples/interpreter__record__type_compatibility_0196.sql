CREATE SCHEMA `auth`;

CREATE SCHEMA `app`;

CREATE TABLE `auth`.`users` (
  `id` int PRIMARY KEY,
  `role` ENUM ('admin', 'user'),
  `status` ENUM ('active', 'inactive')
);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `auth`.`users` (`id`, `role`, `status`)
VALUES
  (1, , CAST('active' AS status)),
  (2, , );

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;