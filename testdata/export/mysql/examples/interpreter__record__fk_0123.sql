CREATE SCHEMA `auth`;

CREATE TABLE `auth`.`users` (
  `id` int,
  `tenant_id` int,
  PRIMARY KEY (`id`, `tenant_id`)
);

CREATE TABLE `posts` (
  `id` int PRIMARY KEY,
  `user_id` int,
  `tenant_id` int,
  `content` text
);

ALTER TABLE `posts` ADD FOREIGN KEY (`user_id`, `tenant_id`) REFERENCES `auth`.`users` (`id`, `tenant_id`);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `auth`.`users` (`id`, `tenant_id`)
VALUES
  (1, 100),
  (2, 100);
INSERT INTO `posts` (`id`, `user_id`, `tenant_id`, `content`)
VALUES
  (1, 1, 100, 'Hello'),
  (2, 999, 100, 'Invalid user');

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;