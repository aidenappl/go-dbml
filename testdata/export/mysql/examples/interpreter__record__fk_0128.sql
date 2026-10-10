CREATE TABLE `users` (
  `id` int PRIMARY KEY,
  `name` varchar(255)
);

CREATE TABLE `posts` (
  `id` int PRIMARY KEY,
  `user_id` int,
  `title` varchar(255)
);

CREATE TABLE `employees` (
  `id` int PRIMARY KEY,
  `manager_id` int,
  `name` varchar(255)
);

ALTER TABLE `posts` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `employees` ADD FOREIGN KEY (`manager_id`) REFERENCES `employees` (`id`);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `users` (`id`, `name`)
VALUES
  (1, 'Alice');
INSERT INTO `posts` (`id`, `user_id`, `title`)
VALUES
  (1, 1, 'Valid Post'),
  (2, 999, 'Invalid Post');
INSERT INTO `employees` (`id`, `manager_id`, `name`)
VALUES
  (1, NULL, 'CEO'),
  (2, 1, 'Manager'),
  (3, 999, 'Invalid Manager Reference');

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;